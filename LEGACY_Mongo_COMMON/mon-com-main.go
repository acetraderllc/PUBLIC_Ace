/*	MONGO_Common Terrys Common Mongo Functions.. 
   ----------------------------------

	Jul 26, 2020	v1.40	- Added some fixes to Mongo Common to prevent the i/o timeout 127.0.0.1 errors
	Feb 14, 2020	v1.33	- Changed my mind, when in a thread be sure to use the much simpler:

							TEMP_SESSION := DBSession.Copy()
							defer TEMP_SESSION.Close()		

							dbOBJ = TEMP_SESSION.DB(DBName).C(COLLECTION_NAME)
							bulk := dbOBJ.Bulk() 


	Feb 13, 2020	v1.30	- Added NEW_THREAD_DB_SESSION for dealing with DB sessions within multi thread routines
	Aug 05, 2019	v1.25	- Convert to Go Modules
	Sep 05, 2018	v1.23	- Initial Rollout

*/

package LEGACY_Mongo_COMMON

import (
	// = = = CORE / Standard Library Deps
	"flag"
	"time"
	"os"
	
	// = = = Personal/Custom Deps
	. "dev.azure.com/acetraderllc/shared/_git/PUBLIC_Ace.git/Terry_COMMON"


	// = = = 3rd Party DEPS
	"gopkg.in/mgo.v2"

)



// **** DATABASE GLOBALS ***
var DBSession *mgo.Session			// Primary DSB Session / Connection object
var DBErr error							// Generic error object

var DBHOST = "localhost"				//THis is the hostname we need to connect to the database

var DBUSER = ""
var DBPASS = ""

var DB_CONNECTION_TIMEOUT = 300		// 300 seconds (5 mins) before we stop trying to connect to mongo

var DBName = "INPUT_NEEDED_MONGO_DB_NAME"							// Name of the mongo database we work with
var DBNAME = "INPUT_NEEDED_MONGO_DB_NAME"			// alt alias for DBName
var COLLECTION_NAME = "INPUT_NEEDED_MONGO_COLLECTION_NAME"				// Name of the COLLECTION within the database we need

var MAX_DB_RETRY = 5


// NOTE: You have MAY have to do this connection stuff in your MAIN section so that defer clsoe gets called at end of the program
// instead of the end of this function
// If you get "Session already closed" errors re evaluate this .. something might be getting Closed early
func DB_INIT(EXTRA_ARGS ...string) {

	//1. Get the extra VARS. this is basically the database name
	for x, VAL := range EXTRA_ARGS {

		//1b. First param is always DB NAME
		if x == 0 {
			DBName = VAL
			DBNAME = VAL
			continue
		}
		
		//2c.Second is always USERNAME
		if x == 1 {
			DBUSER = VAL
			continue
		}

		//3c. Third is always PASS
		if x == 2 {
			DBPASS = VAL
			continue
		}		

	}//end of for

	//2. Fix for timeout
	timeVAL := time.Duration(DB_CONNECTION_TIMEOUT * int(time.Second))
	M.Print(" ** Initiating NEW Connection to Mongo DB: ")
	C.Println(DBHOST)


	//2. Setup the dial info. If credentials are specified, we use them
	// We need this object to establish a session to our MongoDB.
	mongoDBDialInfo := &mgo.DialInfo{
		Addrs:    []string{DBHOST},
		Timeout:  timeVAL,
	}	

	if DBUSER != "" {
		mongoDBDialInfo = &mgo.DialInfo{
			Addrs:    []string{DBHOST},
			Timeout:  timeVAL,
			Database: DBName,
			Username: DBUSER,
			Password: DBPASS,
		}	
	}

	//temp_DBSess, err := mgo.Dial(DBHOST)

	temp_DBSess, err := mgo.DialWithInfo(mongoDBDialInfo)

	if err != nil {
		R.Println("    ERROR conencting to Mongo Server for some reason!!!")
		Y.Println("     -", err)
		W.Println("    ", "- Perhaps Mongo isnt running??")
		M.Println("    ", "- Or Perhaps, credentials are invalid??")
		Y.Println("")

		os.Exit(ERROR_CODE)
	}
	
//	defer temp_DBSess.Close()
	temp_DBSess.SetMode(mgo.Monotonic, true)
	temp_DBSess.SetSocketTimeout(1 * time.Hour)	// Needed to fix those i/o timeout 127.0.0.1 errors weve been seeing

	// Assign to the global DB session object
	DBSession = temp_DBSess.Copy()
} //end of func


func NEW_DB_SESSION() *mgo.Session {
	temp_new_sess := DBSession.Copy()

	return temp_new_sess
}


var DROP_FLAG = false 		// if specified as true, we drop any indexes or collections




func DB_MULTI_INIT(db_hostname string) *mgo.Session  {

	M.Print(" ** Initiating MULTI_DB type Mongo Connect to: ")
	C.Println(db_hostname)

	temp_DBSess, _ := mgo.Dial(db_hostname)
	
	defer temp_DBSess.Close()
	temp_DBSess.SetMode(mgo.Monotonic, true)

	DBSession = temp_DBSess.Copy()
	
	return temp_DBSess.Copy()
	
} //end of func




var INDEX_KEYLIST_TEMP []string

// when called, it adds the index key to the TEMP_INDEX_KEYS
func INDEX_KEY(index_name string)  {

	INDEX_KEYLIST_TEMP = append(INDEX_KEYLIST_TEMP, index_name)

}//end of func

//just an alias to INDEX_KEY
func ADD_KEY(ind_name string) {
	INDEX_KEY(ind_name)
}

//just an alias to INDEX_KEY
func ADD_KEYS(ind_name string) {
	INDEX_KEY(ind_name)
}

// Easy way to create indexes on Tables.. this only creates and index if it DOES NOT already exist (ver EnsureIndex)
// THere is a global that is created (And set to null when we are done) each time this is called
func INDEX_CREATE(index_NAME string, COLLECT_NAME string, isUNIQ bool) bool {

	dbtemp2 := DBSession.DB(DBName).C(COLLECT_NAME)

	Y.Print("\n ** Attempting to create Index '")
	W.Print(index_NAME)
	Y.Print("' on '")
	W.Println(COLLECT_NAME)
		
	// Note, the index defaults to ASCENDING (oldest to newest)
	// If you want DESCENDING, add a - (hyphen) before the below field name(s)

	index2 := mgo.Index{
		Name:       index_NAME,
		Key:        INDEX_KEYLIST_TEMP,
		Unique:     isUNIQ,
		DropDups:   false,
		Background: false,
		Sparse:     false,
	}


	//2. This means we will create the index ONLY if it DOESNT already exist
	err := dbtemp2.EnsureIndex(index2)

	if err != nil {

		fname := "INDEX_ERROR_" + index_NAME + ".log"
		ErrorLog(fname, "Error Creating Index", err)

		return false
	}

	G.Print(" ** Index Creation SUCCESSFUL!!\n\n")


	//5. WE MUST ALWAYS PURGE THE TEMP_KEYS array... CRITICAL
	INDEX_KEYLIST_TEMP = nil

	return true
} //end of index creator



// This is an alias to INDEX_CREATE for backwards compantibility
func INDEX_CREATOR(ind_name string, coll_name string, isUNIQ_FLAG bool) bool {

	return INDEX_CREATE(ind_name, coll_name, isUNIQ_FLAG )
}





/* PURGES ALL RECORDS IN SPECIFIED COLLECTION!!! USE WITH CARE!!
 */
func PURGE_DB(DBCOLL string) {

	C.Println(" ***")
	M.Print("     Safety DROP of the ")
	W.Print("'" + DBCOLL + "'")
	M.Println(" Collection..")
	C.Println(" ***")

	c := DBSession.DB(DBName).C(DBCOLL)

	c.DropCollection()

}




func init() {

	flag.StringVar(&DBHOST,          "dbhost", DBHOST,                "       Hostname of the Mongo DATABASE")

	flag.StringVar(&DBUSER,          "dbuser", DBUSER,                "       DB Username Credentials")
	flag.StringVar(&DBPASS,          "dbpass", DBPASS,                "       DB Passs credentials")

	flag.IntVar(&DB_CONNECTION_TIMEOUT,          "dbtimeout", DB_CONNECTION_TIMEOUT,                "       DB Passs credentials")

	flag.StringVar(&DBName,          "dbname", DBName,              "       Database Name we want to work within")
	flag.StringVar(&DBName,          "db", DBName,                  "       alias for --dbname")
	flag.StringVar(&COLLECTION_NAME, "coll", COLLECTION_NAME,       "       Collection within the Database we want")
	flag.StringVar(&COLLECTION_NAME, "collection", COLLECTION_NAME, "       Alias for --coll")

	flag.BoolVar(&DROP_FLAG,    "drop", DROP_FLAG,         "  If specified we drop thje collection ")
	flag.BoolVar(&DROP_FLAG,    "purge", DROP_FLAG,         "  alias for --drop")

}

