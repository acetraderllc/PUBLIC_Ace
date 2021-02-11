package AWS_ImageOPS

import (

	// -=-= NATIVE LIBS

//	"bytes"
	"flag"
//	"io/ioutil"
//	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	// -=-= COMMON / Personal LIBS
	. "TerryCommon"

	// -=-= 3RD Party Libs

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/aws/credentials/stscreds"
	"github.com/aws/aws-sdk-go/aws/awserr"
)

var SNAPSHOT_RESULTS *ec2.DescribeSnapshotsOutput // Holds search results when querying snapshots
var VOLUME_RESULTS *ec2.DescribeVolumesOutput     // Holds search results when querying VOLUMES

type IMG_item struct {
 	A_id 	string
 	A_name	string
 	A_DESC	string
 	A_date	string
}






// Gets the NAME of an instance when passed TAG OBJs
func GET_INSTANCE_NAME(tobj []*ec2.Tag) string {

	result := ""
	
	for _, t := range tobj {
		keyname := *t.Key
		keyval := *t.Value
		if keyname == "Name" {
			result = keyval
			break
		}
	} //end of for

	return result
} //end of func

func GET_INSTANCE_INFO(inst *ec2.Instance) (string, string, string, string) {
	
	a_id := *inst.InstanceId
	a_name := GET_INSTANCE_NAME(inst.Tags)
	a_status := *inst.State.Name
	a_type := *inst.InstanceType

	return a_id, a_name, a_status, a_type
}




func GET_IMAGE_INFO(IMAGEID string, IN_SESS *session.Session) (string, string) {

	svc := ec2.New(IN_SESS)
	input := &ec2.DescribeImagesInput{
		ImageIds: []*string{
			aws.String(IMAGEID),
		},
	}

	result, err := svc.DescribeImages(input)

	if err != nil {
		R.Println(" ERROR Trying to get Image Name!!")
		return "", ""
	}

	t_name := ""
	t_status := ""

	if *result.Images[0].Name != "" {
		t_name = *result.Images[0].Name
	}

	if *result.Images[0].State != "" {
		t_status = *result.Images[0].State
	}

	return t_name, t_status

} //end of function

func UPDATE_IMAGE_PERMISSION(IN_SESS *session.Session, IMAGE_ID string, OWNER_ACC_ID string) bool {

	var success = false

	M.Print(" *** Granting ")
	C.Print(OWNER_ACC_ID)
	M.Print(" Permission to access ")
	W.Println(IMAGE_ID)

	svc := ec2.New(IN_SESS)
	input := &ec2.ModifyImageAttributeInput{
		ImageId: aws.String(IMAGE_ID),
		LaunchPermission: &ec2.LaunchPermissionModifications{
			Add: []*ec2.LaunchPermission{
				{
					UserId: aws.String(OWNER_ACC_ID),
				},
			},
		},
	}

	_, err := svc.ModifyImageAttribute(input)
	if err != nil {
		R.Print("ERROR, Cannot modify permissions on ")
		W.Print(IMAGE_ID)
		R.Print(" to allow account ")
		Y.Println(OWNER_ACC_ID)
		err.Error()
	} else {

		G.Print(" PERMISSIONS UPDATED SUCCESSFULLY!! ")
		success = true

	}

	W.Println("")
	return success
} //end of function

var MAX_COPY_TIME = 30 // maximum time in minutes that this will run before timing out and existing the program

func DO_IMAGE_COPY(IN_SESS *session.Session, SOURCE_IMAGE_ID string, SOURCE_REGION string, DEST_IMAGE_NAME string, DEST_REGION string, DO_ENCRYPT bool, THIS_THREAD int) (bool, string) {

	Y.Println("\n ** IN PRIMARY COPY ROUTINE ** \n")
	//1. Convert the DEST_REGION to the REGION ID
	_, SOURCE_REGION_ID := SEARCH_REGIONS(SOURCE_REGION)

	svc := ec2.New(IN_SESS)
	input := &ec2.CopyImageInput{
		Name:          aws.String(DEST_IMAGE_NAME),
		Description:   aws.String(DEST_IMAGE_NAME),
		SourceImageId: aws.String(SOURCE_IMAGE_ID),
		SourceRegion:  aws.String(SOURCE_REGION_ID), // Source Region will always be us-east-1
		Encrypted:     aws.Bool(DO_ENCRYPT),
	}

	result, err := svc.CopyImage(input)
	if err != nil {
		R.Println(" ** ERROR, Cannot copy to DEST region in ELG account for some reason!! **")
		Y.Println(err.Error())

		return false, ""
	}

	//3. This is the image ID of the DESTINATION image..
	DEST_IMAGE_ID := *result.ImageId

	//3b. Sleep for a bit to let AWS get started copying
	Sleep(25, true)

	//5. Now lets go into a loop and check for status of the copy
	CTIME := 0

	GOOD_COPY := false

	for {
		_, STATUS := GET_IMAGE_INFO(DEST_IMAGE_ID, IN_SESS)

		W.Print("\n         ** Copy Progress: ")
		M.Print(STATUS)
		W.Print(" on Thread: ")
		G.Print(THIS_THREAD)

		if strings.Contains(STATUS, "pending") {

			W.Print(" ** Waiting for ")
			Y.Print(DEST_IMAGE_ID)
			W.Print(" To finish COPYING TO ")
			G.Print(DEST_REGION, " ** \n\n")

			// Sleep 60 seconds before checking again
			Sleep(60, false)

			//5. When status is no longer "PENDING" we exit this loop. Means it is AVAILABLE
		} else if strings.Contains(STATUS, "available") {

			GOOD_COPY = true
			break

			//6. Safety.. Lets exit if this is going longer than MAX_TIME
		} else if CTIME >= MAX_COPY_TIME {

			R.Println("\n\n ERROR: Copy took WAY too long.. Exiting with ERROR!")
			os.Exit(129)

		}

		CTIME++

	} //end of for loop

	//7. Safety.. if Copy never copmleted we exit
	if GOOD_COPY == false {
		R.Println(" ** ERROR: Copy failed .. Exitting!!")
		return false, ""
	}

	Y.Println(" *** COPY COMPLETE!!! ***\n\n")
	return true, DEST_IMAGE_ID
} //end of function



type ARK_item struct {
	ID          string
	Name        string
	Date        string
	IMG_OBJ     *ec2.Image
	RELEASE_TAG string
	PENDING_TAG string
	ACTION      string
}

/*
	Searches Ec2 Instance TAGS
*/
func SEARCH_EC2_TAGS(TAGOBJ ec2.Instance, input_key string, input_value string) (bool, string) {

	for _, TAG := range TAGOBJ.Tags {

		temp_tagname := *TAG.Key
		temp_tagval := *TAG.Value

		//1. first search pattern is to see if the passed keyname and keyval exist
		// if so, we return tthe keyval

		if temp_tagname == input_key && temp_tagval == input_value {
			return true, temp_tagval

			//2. If input VALUE is blank, we want to RETURN the value of whatever is input_key
		} else if input_value == "" {

			if temp_tagname == input_key {

				return true, temp_tagval
			}
		}

	} //end of for

	//4. Otherwise, nothing is found, we return false
	return false, ""

} // end of function

/*
	Searches AMI IMage Tags
*/
func SEARCH_IMAGE_TAGS(TAGOBJ *ec2.Image, input_key string, input_value string) (bool, string) {

	for _, TAG := range TAGOBJ.Tags {

		temp_tagname := *TAG.Key
		temp_tagval := *TAG.Value

		//1. first search pattern is to see if the passed keyname and keyval exist
		// if so, we return tthe keyval

		if temp_tagname == input_key && temp_tagval == input_value {
			return true, temp_tagval

			//2. If input VALUE is blank, we want to RETURN the value of whatever is input_key
		} else if input_value == "" {

			if temp_tagname == input_key {

				return true, temp_tagval
			}
		}

	} //end of for

	//4. Otherwise, nothing is found, we return false
	return false, ""

} // end of function

/*
   //1. First lets QUEUE up the RELEASED items
   QUEUE_EXISTING_RELEASED_ITEMS()

   //2. Process RELEASED item queue.. This archives and deletes items. The current RELEASED item gets archived..
   ARCHIVE_EXISTING_RELEASED_ITEMS()

   //3. Now lets queue up the NEWLY built items (things compiled today)
   QUEUE_NEWLY_BUILT_ITEMS()

   //4. Now lets take the most recent BUILT items (things where released tag is FALSE)
   // And RELEASE THEM. this gives us a NEW item with RELEASED tag set to whatever value is in the PENDING tag
   PERFORM_NEW_RELEASE_LOGIC()
*/

/*
	Searches the Tags on EC2 Volumes
*/
func SEARCH_VOLUME_TAGS(TAGOBJ ec2.Volume, input_key string, input_value string) (bool, string) {

	for _, TAG := range TAGOBJ.Tags {

		temp_tagname := *TAG.Key
		temp_tagval := *TAG.Value

		//1. first search pattern is to see if the passed keyname and keyval exist
		// if so, we return tthe keyval

		if temp_tagname == input_key && temp_tagval == input_value {
			return true, temp_tagval

			//2. If input VALUE is blank, we want to RETURN the value of whatever is input_key
		} else if input_value == "" {

			if temp_tagname == input_key {

				return true, temp_tagval
			}
		}

	} //end of for

	//4. Otherwise, nothing is found, we return false
	return false, ""

} // end of function

/*
	Searches the Tags on EC2 SNAPSHOTS
*/
func SEARCH_SNAPSHOT_TAGS(TAGOBJ ec2.Snapshot, input_key string, input_value string) (bool, string) {

	for _, TAG := range TAGOBJ.Tags {

		temp_tagname := *TAG.Key
		temp_tagval := *TAG.Value

		//1. first search pattern is to see if the passed keyname and keyval exist
		// if so, we return tthe keyval

		if temp_tagname == input_key && temp_tagval == input_value {
			return true, temp_tagval

			//2. If input VALUE is blank, we want to RETURN the value of whatever is input_key
		} else if input_value == "" {

			if temp_tagname == input_key {

				return true, temp_tagval
			}
		}

	} //end of for

	//4. Otherwise, nothing is found, we return false
	return false, ""

} // end of function


// This gets teh value of ANY Tag on an instance.. returns its valuer
// making this case insensitive 
func GET_KEYVAL(tobj []*ec2.Tag, find_keyname string) string {

	result := ""
	find_keyname = strings.ToLower(find_keyname)

	for _, t := range tobj {
		t_key := strings.ToLower(*t.Key)

		t_val := *t.Value
		if t_key == find_keyname {
			result = t_val
			break
		}
	} //end of for

	return result
} //end of func



var AMI_SEARCH_RESULTS *ec2.DescribeImagesOutput  // Holds search results when querying AMI's
func GET_AMI_COUNT() (int, string) {

	tcount := 0

	for _, IMJ := range AMI_SEARCH_RESULTS.Images {
		// This is just here so the loop completes (we need to make use of the IMJ object)
		if IMJ.State == nil {

		}
		tcount++
	}

	tstring := ShowNumber(tcount)

	Y.Print("\n   Total AMI's found: ")
	W.Println(tstring, "\n")

	return tcount, tstring
}
