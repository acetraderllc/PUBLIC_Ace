/*   Terry_COMMON - Common code snippets i use for Go Development used module code
---------------------------------------------------------------------------------------
NOTE: For Functions or Variables to be globally availble. The MUST start with a capital letter.
	  (This is a GO Thing)

	Feb 22, 2021	v1.90	- Added SCRAPE_TOOL for screen scraping
	Feb 15, 2020	v1.81	- Major Revamp to the GET_CURRENT_TIME and also have --zone TIME_ZONE_FLAG variable available to force Timezone

	Feb 14, 2020	v1.79	- Removed a lot of redundant functions (for date in particular)
							- Added an awesome ADD_LEADING_ZERO  function
							- Remove some more stuff i dont need
							- Working on the SHOW_PRETTY_DATE
							
	Feb 13, 2020	v1.76	- Removed UNEEDED stuff
	Feb 03, 2020	v1.75	- Added badass GET_CALLER_FUNC_NAME (for finding the parent function call)
	Jan 26, 2020	v1.73	- Got rid of more redundant functions that got pushed to APIcebaby
	Jan 23, 2020	v1.67	- Removed redundant fucntions (stuff that is now in APIce)
	Jan 05, 2020	v1.66	- Some cosmetic changes, Updated TerryCOMMON again
	Dec 29, 2019	v1.63	- Updated TerryCOMMON again
	Jun 05, 2014    v1.23   - Initial Rollout

*/

package Terry_COMMON

import (

	//1.  = = = CORE / Standard Library Deps
		"bufio"
		"crypto/md5"
		"encoding/hex"
		"flag"
		"io"
		"io/ioutil"
		"math"
		"math/rand"
		"net/http" 			// Needed for the functions that send JSON back and forth
		"os"
		"os/exec"
		"regexp"
		"runtime"
		"sort"
		"strconv"
		"strings"
		"time"
		"unicode"

	
	//2. = = = 3rd Party DEPENDENCIES
		"github.com/atotto/clipboard"
		"github.com/briandowns/spinner"
		"github.com/dustin/go-humanize"
		"github.com/fatih/color"
		"github.com/PuerkitoBio/goquery"
)

/*
	- - - -
	- - - -
	- - - - START OF GLOBALS WE NEED - - - - - -
	- - - -
	- - - -
*/

var SERIAL_NUM = "" // This is unique execution sid generated everytime a program starts. useful for troubleshooting in jenkins
var SHOW_SERIAL = false 		// If set, we Generate and show a serial number

var OSTYPE=""
var CURRENT_OS = ""	
var GOOS_VALUE = "" 		// Holds the current OS as reported by runtime.GOOS

var DEBUG_MODE = false		// Useful universal flag for enabling DEBUG_MODE code blocks
var ERROR_EXIT_CODE = -9999
var ERROR_CODE = ERROR_EXIT_CODE

var TIME_ZONE_FLAG = ""		// This is flag we use to set/force proper timezones

var PROG_START_TIME string
var PROG_START_TIMEOBJ time.Time

var GLOBAL_CURR_DATE = ""		// Current Actual Date in the Timezone we specified
var GLOBAL_CURR_TIME = ""		// alias to CURR_DATE
var GLOBAL_DATE_OBJ time.Time	// Actual Global Date OBJ


// -=-=-= COMMON COLOR GLOBAL references =-=-=-=-

var R = color.New(color.FgRed, color.Bold)
var G = color.New(color.FgGreen, color.Bold)
var Y = color.New(color.FgYellow, color.Bold)
var B = color.New(color.FgBlue, color.Bold)
var M = color.New(color.FgMagenta, color.Bold)
var C = color.New(color.FgCyan, color.Bold)
var W = color.New(color.FgWhite, color.Bold)

var R2 = color.New(color.FgRed)
var G2 = color.New(color.FgGreen)
var Y2 = color.New(color.FgYellow)
var B2 = color.New(color.FgBlue)
var M2 = color.New(color.FgMagenta)
var C2 = color.New(color.FgCyan)
var W2 = color.New(color.FgWhite)

var R3 = color.New(color.FgRed, color.Underline)
var G3 = color.New(color.FgGreen, color.Underline)
var Y3 = color.New(color.FgYellow, color.Underline)
var B3 = color.New(color.FgBlue, color.Underline)
var M3 = color.New(color.FgMagenta, color.Underline)
var C3 = color.New(color.FgCyan, color.Underline)
var W3 = color.New(color.FgWhite, color.Underline)


/*
	- - - -
	- - - -     End of GLOBALS definitions
	- - - -
*/

// Returns true if the string contains ONLY numbers
func HasOnlyNumbers(s string) bool {
    for _, r := range s {
        if (r < '0' || r > '9') {
            return false
        }
    }
    return true
} //end of func

// This takes in a string and converts it to a float
func CONVERT_FLOAT(input string, precision int) (float64, string) {

	f_NUM, _ := strconv.ParseFloat(input, 64)
	NUM_result := FIX_FLOAT_PRECISION(f_NUM, precision)
	FIXED_text := strconv.FormatFloat(NUM_result, 'f', precision, 64)
	
	return NUM_result, FIXED_text

} //end of func


// This converts a float to a WHOLE number
func CONVERT_FLOAT_TO_WHOLE(infloat float64, CV_PRECISION int) (int, string) {
	entry_FIXED_text := strconv.FormatFloat(infloat, 'f', CV_PRECISION, 64)
	entry_NUM_text := strings.Replace(entry_FIXED_text, ".", "", -1)
	entry_NUM, _ := strconv.Atoi(entry_NUM_text)

	return entry_NUM, entry_NUM_text
}

// Makes a floating point number rounded up and returns integer
func MakeRound(num float64) int {
	return int(num + math.Copysign(0.5, num))
}

func FIX_FLOAT_PRECISION(num float64, precision int) float64 {
	output := math.Pow(10, float64(precision))
	return float64(MakeRound(num*output)) / output
}

// This removes all spaces from a string via unicode
func REMOVE_ALL_SPACES(str string) string {
	return strings.Map(func(r rune) rune {
		if unicode.IsSpace(r) {
			return -1
		}

		return r
	}, str)
}

// This utilizes the HUMANIZE library and shows a HUMAN readable number of the passed variable
func ShowNumber(innum int) string {

	result := humanize.Comma(int64(innum))

	// result = strconv.Itoa(result)
	return result
}

func ShowNum64(innum int64) string {

	result := humanize.Comma(innum)

	// result = strconv.Itoa(result)
	return result
}

// This is an alias to ShowNumber
func ShowNum(innum int) string {

	return ShowNumber(innum)
}

func IS_EVEN(input_NUM int) bool {

	if input_NUM%2 == 0 {
		return true

	}

	return false
}



// MULTI-PURPOSE SCREEN SCRAPE TOOL
// Params: URL, UserAgent
// Returns: bool, GOQUERY_DOC, Text of Response

func SCRAPE_TOOL(EXTRA_ARGS ...string) (bool, *goquery.Document, string) {

	C.Println("")
	C.Println(" *** Calling SCRAPE_TOOL ***")
	C.Println("")

	URL := ""

	// Defaults to CHrome
	USER_AGENT := "Mozilla/5.0 (Macintosh; Intel Mac OS X 11_2_1) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/88.0.4324.182 Safari/537.36"		
	
	//1. Get tvars passed
	for x, VAL := range EXTRA_ARGS {

		//1b. First param is always DB NAME
		if x == 0 {
			URL = VAL
			continue
		}
		if x == 1 {
			USER_AGENT = VAL
			continue
		}
	}

	var GOQUERY_doc *goquery.Document

	//2. Now generate a NewRequest Object with http
	client := &http.Client{}
	req, err := http.NewRequest("GET", URL, nil)
	if err != nil {
		R.Println(" *** ")
		R.Println(" *** ERROR IN SCRAPE TOOL - During http OBJECT Create: ")
		Y.Println(err)
		R.Println(" *** ")
		R.Println("")
		return false, GOQUERY_doc, ""
	}

	//3. Next, Set the User Agent the client will use during the HTTP Pull
	req.Header.Set("User-Agent", USER_AGENT)

	//3b. Now.. actually do the Http Client Request (with the header)
	res, err2 := client.Do(req)
	if err2 != nil {
		R.Println(" *** ")
		R.Println(" *** ERROR IN SCRAPE TOOL - During CLIENT HTTP Pull: ")
		Y.Println(err2)
		R.Println(" *** ")
		R.Println("")
		return false, GOQUERY_doc, ""
	}
	defer res.Body.Close()	

	//4. If we got this far, all is well.. Lets query the body of the response and put it into TEXT mode
	body, err3 := ioutil.ReadAll(res.Body)
	if err3 != nil {
		R.Println(" *** ")
		R.Println(" *** ERROR IN SCRAPE TOOL - IoUtil Body Parse: ")
		Y.Println(err3)
		R.Println(" *** ")
		R.Println("")
		return false, GOQUERY_doc, ""
	}		

	FULL_RESPONSE_TEXT := string(body)

	//5. Now finally, lets create our DOM object using goquery
	doc, err4 := goquery.NewDocumentFromReader(res.Body)
	if err4 != nil {
		R.Println(" *** ")
		R.Println(" *** ERROR IN SCRAPE TOOL - During GOQUERY: ")
		Y.Println(err4)
		R.Println(" *** ")
		R.Println("")
		return false, GOQUERY_doc, ""
	}	

	//6. All Done!! Return the Goquery DOC Object

 	doc.Find("td[class^=snapshot-td2-cp]").Each(func(i int, O *goquery.Selection) {
		W.Print(i, " ")
		C.Println(" TEXT: **" + O.Text() + "** ")
  	})

	PressAny()	

	//2. Value of stuff in boxes
 	doc.Find(".snapshot-td2").Each(func(i int, O *goquery.Selection) {
		Y.Print(i, " ")
		C.Println(" VAL: **" + O.Text() + "** ")
  	})
	PressAny()	

	return true, doc, FULL_RESPONSE_TEXT

} //end of func



/* This allows you to add/subtract a percentage ie 1.05 TO a specified value */
func PERCENT_MATH(input_FLOAT float64, command string, perc float64) float64 {
	/*

		  Percentage calculation formula is:
			- To ADD a percantage to a value:

				 1530.56 * ( (100 + 0.65) / 100 )

			- And to SUBTRACT:

				1530.56 * ( (100 + 0.65) / 100 )
	*/

	if command == "add" || command == "ADD" {

		result := input_FLOAT * ((100 + perc) / 100)
		return result
	}

	if command == "sub" || command == "SUB" || command == "subtract" {

		result := input_FLOAT * ((100 - perc) / 100)
		return result
	}

	return 0.0
}


func IS_ODD(input_NUM int) bool {

	if input_NUM%2 == 0 {

	} else {

		return true
	}

	return false
}

func CLIPBOARD_COPY(instring string) {
	clipboard.WriteAll(instring)
}


// This takes in a string and removes all non alphanumeric chars from it.. and extra spaces
func CLEAN_STRING(input string) string {
	justAlpha, _ := regexp.Compile("[^a-zA-Z0-9_ ]")
	killExtraSpace := regexp.MustCompile(`[\s\p{Zs}]{2,}`)

	PASS_1 := justAlpha.ReplaceAllString(input, "")
	FINAL_PASS := killExtraSpace.ReplaceAllString(PASS_1, " ")

	return FINAL_PASS
} //end of func

// This removes all extra spaces in a string (like takes     THIS... and removes those extras)
func REMOVE_Extra_Spaces(input string) string {

	re_leadclose_whtsp := regexp.MustCompile(`^[\s\p{Zs}]+|[\s\p{Zs}]+$`)
	re_inside_whtsp := regexp.MustCompile(`[\s\p{Zs}]{2,}`)
	final := re_leadclose_whtsp.ReplaceAllString(input, "")
	final = re_inside_whtsp.ReplaceAllString(final, " ")

	return final
}

// Shows the Local TIMEZONE info .. or if --zone is specified, uses that
func SHOW_ZONE_INFO() {

	W.Print(" [ Timezone: ")
	if ZONE_LOCAL == "Local" {
		C.Print(ZONE_LOCAL)
	}
	Y.Print(" " + ZONE_UPPER)
	G.Print(" (", ZONE_HOUR_OFFSET, ") ")
	W.Println(" ]")
} //end of func

var UTC_LOCATION_OBJ, _ = time.LoadLocation("UTC")
var PST_LOCATION_OBJ, _ = time.LoadLocation("America/Los_Angeles")		// aka PST
var CST_LOCATION_OBJ, _ = time.LoadLocation("America/Chicago")			// aka CST
var EST_LOCATION_OBJ, _ = time.LoadLocation("EST")
var DEFAULT_ZONE_LOCATION_OBJ, _ = time.LoadLocation("Local")

var ZONE_HOUR_OFFSET = "" 
var ZONE_UPPER = ""
var ZONE_FULL = ""
var ZONE_LOCAL = ""

// This initializes some pretty necessary timezone defaults. 
func SET_TIMEZONE_DEFAULTS() {

	//1. By Deffault this will be Local
	ZONE_UPPER = "Local"		//DEFAULT_ZONE_LOCATION_OBJ.String()
 
	//1b. If timezone is NOT blank, we set up our zone objects
	if TIME_ZONE_FLAG != "" {

		//1b. Safety to make sure we foce TIME_ZONE_FLAG to lowercase
		TIME_ZONE_FLAG = strings.ToLower(TIME_ZONE_FLAG)	

		switch TIME_ZONE_FLAG {
			case "pst":
				DEFAULT_ZONE_LOCATION_OBJ = PST_LOCATION_OBJ

			case "cst":
				DEFAULT_ZONE_LOCATION_OBJ = CST_LOCATION_OBJ

			case "est":
				DEFAULT_ZONE_LOCATION_OBJ = EST_LOCATION_OBJ

			case "utc":
				DEFAULT_ZONE_LOCATION_OBJ = UTC_LOCATION_OBJ
		
		} //end of switch	
	}

	//2. Also lets get the Timezone and OFFSET info	
	t := time.Now().In(DEFAULT_ZONE_LOCATION_OBJ)
	curr_zone, offset := t.Zone()

	//2b. see if it is negative, we need to convert this temporarily
	hprefix := "+"
	if offset < 0 {
		fixed_off := math.Abs(float64(offset))

		//2b. Convert fixed_off to an integar
		offset, _ = CONVERT_FLOAT_TO_WHOLE(fixed_off, 0)
		hprefix = "-"
	}

	//3. NOw convert the offset seconds to hours
	off_hours := (offset / 60) / 60
	offstring := hprefix + strconv.Itoa(off_hours) + " hours"

	//4. Set ZONE_UPPER .. we use this in other places. 
	ZONE_UPPER = curr_zone
	ZONE_LOCAL = t.Location().String()
	
	ZONE_HOUR_OFFSET = offstring
	ZONE_FULL = " (" + curr_zone + " " + offstring + ")"
	
	//4. Get the current time to set the global time vaiables
	GLOBAL_CURR_DATE, GLOBAL_DATE_OBJ = GET_CURRENT_TIME("full")
	GLOBAL_CURR_TIME = GLOBAL_CURR_DATE


} //end of func


// Gets the current Time (defaults to UTC, if you specify EST you get EST time)
func GET_CURRENT_TIME(EXTRA_ARGS ...string) (string, time.Time) {

	//1. Default ot the local machines time zone
	dateOBJ := time.Now()

	var output_FORMAT = "full"		// Full is the default format we use

	var PASS_FLAG = false

	//2. Now, see if flags were specified. Iterate through them
	for _, VAL := range EXTRA_ARGS {

		VAL = strings.ToLower(VAL)

		switch VAL {
			case "pst":
				dateOBJ = dateOBJ.In(PST_LOCATION_OBJ)
				PASS_FLAG = true

			case "cst":
				dateOBJ = dateOBJ.In(CST_LOCATION_OBJ)
				PASS_FLAG = true

			case "est":
				dateOBJ = dateOBJ.In(EST_LOCATION_OBJ)
				PASS_FLAG = true

			case "utc":
				dateOBJ = dateOBJ.In(UTC_LOCATION_OBJ)
				PASS_FLAG = true
	
		} //end of switch

		//3. If full, short, british or iso specified, set the output format
		if VAL == "short" || VAL == "full" || VAL == "british" || VAL == "iso" || VAL == "justdate" {
			output_FORMAT = VAL
		}

	} //end of for

	//3. Otherwise, if the global TIME_ZONE_FLAG is set, (to cst, est or pst).. We use THAT object

	if TIME_ZONE_FLAG != "" && PASS_FLAG == false {

		switch TIME_ZONE_FLAG {
			case "pst":
				dateOBJ = dateOBJ.In(PST_LOCATION_OBJ)

			case "cst":
				dateOBJ = dateOBJ.In(CST_LOCATION_OBJ)

			case "est":
				dateOBJ = dateOBJ.In(EST_LOCATION_OBJ)

			case "utc":
				dateOBJ = dateOBJ.In(UTC_LOCATION_OBJ)		
		}
	}

	result, _ := SHOW_PRETTY_DATE(dateOBJ, output_FORMAT)

	return result, dateOBJ

} //end of func

// This takes in a number and returns a string with a leading 0
// If the number is already 10 or greater, it returns that same number as is
func ADD_LEADING_0( myNum int) string {

	RESULT := strconv.Itoa(myNum)

	if myNum <= 9 {
		RESULT = "0" + RESULT
	}

	return RESULT
}

// Alias for ADD_LEADING_0
func ADD_LEADING_ZERO(myNum int) string {
	return ADD_LEADING_0(myNum)
}


/* Enhanced Pretty Date function Takes in a time.Time DATE_OBJ and returns a PRETTY formatted based on what you specify
   Returns a STRING and a WEEKDAY
*/
func SHOW_PRETTY_DATE(input_DATE time.Time, EXTRA_ARGS...string) (string, string) {
	var output_FORMAT = ""
	var SHOW_SECONDS = false

	//1. Parse out EXTRA_ARGS
	for _, VAL := range EXTRA_ARGS {

		//1c. If sec or seconds is passed, we also will show the seconds
		if VAL == "sec" || VAL == "seconds" {
			SHOW_SECONDS = true
			continue

		//1e. If short is passed, we show this format: Wednesday, 11/20/2001
		// If full is passed, we show this format: Wednesday, 11/20/2020 @ 13:56
		// if british or iso is passed, we show: 2015-05-30
		} else if VAL != "" {
			output_FORMAT = VAL
			continue
		}

	} // end of for

	//2. From this object, extract the M/D/Y HH:MM
	montemp := int(input_DATE.Month())
	daytemp := input_DATE.Day()

	hourtemp := input_DATE.Hour()
	mintemp := input_DATE.Minute()

	//3. Then, we add leading 0's as needed
	cMon := ADD_LEADING_ZERO(montemp)
	cDay := ADD_LEADING_ZERO(daytemp)
	cHour := ADD_LEADING_ZERO(hourtemp)
	cMin := ADD_LEADING_ZERO(mintemp)
	
	sectemp := input_DATE.Second()
	cSec := ADD_LEADING_ZERO(sectemp)


	//4. Thankfully we dont have to worry about this fuckery with the year!
	cYear := strconv.Itoa(input_DATE.Year())
	weekd := input_DATE.Weekday().String()

	/* 7. Here is the DEFAULT Pretty format that is returned

		09/26/1978 @ 13:58

			or (if SHOW_SECONDS is passed) 

		09/26/1978 @ 13:58:05
	*/
	result_TEXT := cMon + "/" + cDay + "/" + cYear + " @ " + cHour + ":" + cMin
	if SHOW_SECONDS {
		result_TEXT += ":" + cSec
	}

	//8. SHORT Format is:  Wednesday, 11/20/2001
	if output_FORMAT == "short" {

		result_TEXT = weekd + ", " + cMon + "/" + cDay + "/" + cYear

	//9. FULL Format: //Wednesday, 11/20/2020 @ 13:56 EST (-5 Hours)
	} else if output_FORMAT == "full" {
		
		result_TEXT = weekd + ", " + cMon + "/" + cDay + "/" + cYear + " @ " + cHour + ":" + cMin

		if SHOW_SECONDS {
			result_TEXT += ":" + cSec
		}

		result_TEXT += " " + ZONE_FULL
	
	//10. This is the british/iso format: 2020-09-26
	} else if output_FORMAT == "british" || output_FORMAT == "iso" {

		result_TEXT = cYear + "-" + cMon + "-" + cDay

	//11. This is JUSTDATE:  09/26/1988
	} else if output_FORMAT == "justdate" {

		result_TEXT = cMon + "/" + cDay + "/" + cYear
	}

	//12. As a bonus, we always return the weekday as a second variable

	return result_TEXT, weekd
} //end of func



func SHOW_START_and_END_TIME() {


	Y.Println("\n\n ****************************************************** ")

	endTime, endOBJ := GET_CURRENT_TIME()

	W.Print("              Start Time:")
	B.Println(" " + PROG_START_TIME)
	Y.Print("                End Time:")
	M.Println(" " + endTime)
	C.Print("      Total PROGRAM DURATION: ")
	TIME_DIFF := DISPLAY_TIME_DIFF(PROG_START_TIMEOBJ, endOBJ)	

	G.Println(" ", TIME_DIFF)
	C.Println("******************************************************")
}


// This is an alias to SHOW_PRETTY_DATE
func SHOW_PRETTY_TIME(input_DATE time.Time, EXTRA_ARGS...string) (string, string) {

	return SHOW_PRETTY_DATE(input_DATE)

} //end of alias func



/* Converts TEXT date strings passed in any of the following formats:
	
	- YYYY-MM-DD		(ISO / British format)
	- MM/DD/YYYY
	- MM-DD-YYYY

   and then returns a normalized  M/D/Y format...as well as the weekday... and a time.Time DateOBJ
   If you pass "short" or "full" you will receive a date formatted that way (uses SHOW_PRETTY_DATE)
*/
func CONVERT_DATE(inputDate string, EXTRA_ARGS ...string) (string, string, time.Time) {

	var sMon string
	var sDay string
	var sYear string

	var output_FORMAT = ""

	//1. First parameter is always the input date in the proper format
	for _, VAL := range EXTRA_ARGS {

		//1b. If short or full was passed, we format the output date that way
		if VAL != "" {
			output_FORMAT = VAL
			continue
		}
	} //end of for	

	/* 2. Now, Determine which format the data is in. 
		If it has hyphens can be either YYYY-MM-DD   or
		MM-DD-YYYY
	*/
	if strings.Contains(inputDate, "-") {
		sd := strings.Split(inputDate, "-")

		//1b. See if the first element has FOUR chars.. or TWO chars
		if len(sd[0]) == 2 {

			sMon = strings.TrimSpace(sd[0])
			sDay = strings.TrimSpace(sd[1])
			sYear = strings.TrimSpace(sd[2])

		} else if len(sd[0]) == 4 {

			sYear = strings.TrimSpace(sd[0])
			sMon = strings.TrimSpace(sd[1])
			sDay = strings.TrimSpace(sd[2])		
		}
	
	//3. Or, if this is formatted with / forward slash
	} else if strings.Contains(inputDate, "/") {

		sd := strings.Split(inputDate, "/")

		sMon = strings.TrimSpace(sd[0])
		sDay = strings.TrimSpace(sd[1])
		sYear = strings.TrimSpace(sd[2])	
	}

	//4. Slight bugfix, for sYear (if time was appended we need to fix this and get rid of the space)
	if strings.Contains(sYear, " ") {
	
		ksplit := strings.Split(sYear, " ")
		sYear = ksplit[0]
	
	} else if strings.Contains(sDay, " ") {
		
		ksplit := strings.Split(sDay, " ")
		sDay = ksplit[0]

	}
 

	/*5. For dealing with TIME if it is appended to the date,MUST be formatted like this:

		2020-02-07 14:30
		     
			 or

		2020-02-07_14:30

			or

		2020-02-07@14:30			
	*/
	sHour := ""
	sMin := ""
	num_hour := 0
	num_min := 0

	var hd []string

	if strings.Contains(inputDate, ":") {

		
		//6. OR.. Split on the @ (amphersand) if that was provided
		if strings.Contains(inputDate, "@") {
			hd = strings.Split(inputDate, "@")
		
		//7. Split on the _ (underscore) if that was used
		} else if strings.Contains(inputDate, "_") {
			hd = strings.Split(inputDate, "_")

		//8. By default, Split on the SPACE if that was provided
		} else if strings.Contains(inputDate, " ") {
			hd = strings.Split(inputDate, " ")
		}

		//9. now split on the :
		timePart := strings.TrimSpace(hd[1])
		tp := strings.Split(timePart, ":")
		sHour = strings.TrimSpace(tp[0])
		sMin = strings.TrimSpace(tp[1])

		num_hour, _ = strconv.Atoi(sHour)
		num_min, _ = strconv.Atoi(sMin)
	}

	//10. Now, we need to convert the Mon, Day, Year to integers
	num_year, _ := strconv.Atoi(sYear)
	num_month, _ := strconv.Atoi(sMon)
	num_day, _ := strconv.Atoi(sDay)

	//11. Now Create the dateOBJ,  Format is: Y, MonOBJ, D, H, M, S, Nano
	monthObj := time.Month(num_month)	
	dOBJ := time.Date(num_year, monthObj, num_day, num_hour, num_min, 0, 0, DEFAULT_ZONE_LOCATION_OBJ)

	//12. Now pass to show_Pretty_Date with the output format if specified
	OUTPUT, weekday := SHOW_PRETTY_DATE(dOBJ, output_FORMAT)

	return OUTPUT, weekday, dOBJ
} //end of func


// Another Alias for CONVERT_DATE
func CONVERT_TIME(inputDate string) (string, string, time.Time) {
	return CONVERT_DATE(inputDate)
}


/* Takes in two date objects and returns the TIME DIFFERNCE between them in the 5m40s format
 */

func GET_TIME_DIFF(startTime time.Time, endTime time.Time) string {
	diff := endTime.Sub(startTime)
	return diff.String()
}

// Alias for GET_TIME_DIFF
func DISPLAY_TIME_DIFF(startTime time.Time, endTime time.Time) string {
	return GET_TIME_DIFF(startTime, endTime)
}

var SPINNER_SPEED = 100
var SPINNER_CHAR = 4

var spinOBJ = spinner.New(spinner.CharSets[14], 100*time.Millisecond)

func START_Spinner() {

	sduration := time.Duration(SPINNER_SPEED)

	spinOBJ = spinner.New(spinner.CharSets[SPINNER_CHAR], sduration*time.Millisecond)
	spinOBJ.Start()
}

func STOP_Spinner() {

	spinOBJ.Stop()
}


// this acts as a "delimiter" matrix..It returns true if any of the below delimiters exists
func Delimiter_Matrix_CATCHER(r rune) bool {

	/* In this case return true if there is ANY of the following:

	   - HYPHEN
	   _ UNDERSCORE
	   : Colon
	   / Forward Slash
	   | Pipe


	*/
	return r == ':' || r == '-' || r == '/' || r == '_' || r == '|' || r == '=' || r == '&'
}

/*
 	This is an UBER split on delimiter routine i made that
	makes it easy to return delimited values on MULTIPLE delimiters
	(though I am excluding comma! check the delimiter Matrix!
	It also trims spaces before and AFTER each element
	Pass it a string and it will return an array delimiting on:

	   - HYPHEN
	   _ UNDERSCORE
	   : Colon
	   / Forward Slash
	   | Pipe

*/
func UBER_Split(myText string) []string {

	ptempVals := strings.FieldsFunc(myText, Delimiter_Matrix_CATCHER)

	for x := 0; x < len(ptempVals); x++ {

		ptempVals[x] = strings.TrimSpace(ptempVals[x])

	}

	return ptempVals
} //end of function

/* Takes in a date object and adds or subtracts
based on the number and whatever operation you specify
returns a date object
*/
func DateMath(dateObj time.Time, operation string, v_amount int, interval string) (time.Time, string) {

	//dateObj = dateObj.UTC()

	//1. If we are subtracting, we change amount to a negative number
	if operation == "sub" || operation == "subtract" {

		v_amount = -v_amount

	}

	//2. Now we do the add or subtract operattion based on the time.Duration that is interval
	// Default is minute

	timeINT := time.Minute

	if interval == "hour" || interval == "hours" {

		timeINT = time.Hour

	} else if interval == "min" || interval == "mins" || interval == "minute" || interval == "minutes" {

		timeINT = time.Minute

	} else if interval == "sec" || interval == "secs" || interval == "second" || interval == "seconds" {

		timeINT = time.Second

	} else if interval == "day" || interval == "days" {

		timeINT = (time.Hour * 24)

	}

	//3. Finally do the "date math" on the incoming dateObj
	result := dateObj.Add(time.Duration(v_amount) * timeINT)

	prettyDATE, _ := SHOW_PRETTY_DATE(result)

	return result, prettyDATE

} //end of dateMath

// This splits only on PIPE....and trims space before and after
func PIPE_SPLIT(incoming string) []string {

	ptempVals := strings.Split(incoming, "|")

	// Lets go through each element and trim the spaces from the end
	for x := 0; x < len(ptempVals); x++ {

		ptempVals[x] = strings.TrimSpace(ptempVals[x])

	}

	return ptempVals

} //end of PIPE_SPLIT


func PERC_CALC(child float64, parent float64, PERC_PRECISION int) (float64, string) {

	if child == parent {
		return 0.0, "0.0%"
	}


	//2. Hot fix to make sure we always return a positive number
	OLD_child := child
	OLD_parent := parent

	if child > parent {
		child = OLD_parent
		parent = OLD_child
	}

	mperc := (child / parent) * 100
	mperc = 100 - mperc

	percSTRING := strconv.FormatFloat(mperc, 'f', PERC_PRECISION, 64)
	percNUM, _ := strconv.ParseFloat(percSTRING, 64)
	percSTRING = percSTRING + "%"

	return percNUM, percSTRING

} //end of func

func PERCENTAGE_OF_using_FLOAT(currNUM float64, prevNUM float64) (float64, string) {

	percNUM, percSTRING := PERC_CALC(currNUM, prevNUM, 2)

	return percNUM, percSTRING
	// return percNUM, percSTRING
} //end of func

/*
This returns the percentage smallNUM is of parentNUM  (for IFloats / doubles )
*/
func PERCENTAGE_OF_using_INTEGER(currNUM int, prevNUM int) (float64, string) {

	smallFLOAT := float64(currNUM)
	parentFLOAT := float64(prevNUM)

	percNUM, percSTRING := PERC_CALC(smallFLOAT, parentFLOAT, 2)

	return percNUM, percSTRING

	// return percNUM, percSTRING
} //end of func

// Make sthe first character of a string UPPER CASE
func UpperFirst(inString string) string {

	a := []rune(inString)
	a[0] = unicode.ToUpper(a[0])

	return string(a)
}

// this is a simple sleep function
var SLEEP_SILENT = false // if you set this to true, we do NOT display any output
func Sleep(seconds int, showoutput bool) {

	// the INT is the time in seconds

	if showoutput == true {
		secText := ""
		suffix := "seconds"
		sectemp := seconds

		if seconds >= 119 {
			sectemp = seconds / 60
			suffix = "minutes"
		}

		secText = strconv.Itoa(sectemp)

		C.Print(" ** Sleeping for: ", secText, " ", suffix, "...\n")
	}

	duration := time.Duration(seconds) * time.Second
	time.Sleep(duration)

} //end of sleep function


func VERIFICATION_PROMPT(warning_TEXT string, required_input string) {

	M.Println("\n      - - - - - - - - WARNING - - - - - - - - - - - - - -")
	
	for x := 0; x < 3; x++ {
		C.Println("")
		C.Println("      ", warning_TEXT)
		C.Print("       Type: ")
		G.Print(required_input)
		C.Println(" To Continue")
		Y.Print("       RESPONSE: ")
		userResponse := GET_USER_INPUT()

		if strings.Contains(userResponse, required_input) {
			return
		} else {
			R.Println("\n ! ! ! ! ! ! INVALID RESPONSE  ! ! ! ! ! !")
		    M.Println("\n     - - - - - - - - - - - - - - - - - - - - - - - - -")			
		}
	} //end of for
	

	//2. If we get this far without a valid response, we will exit the program without proceeding
	os.Exit(ERROR_EXIT_CODE)


} //end of prompt

func PROMPT(warning_TEXT string, required_input string) {
	VERIFICATION_PROMPT(warning_TEXT, required_input)
}


func GET_USER_INPUT() string {
	reader := bufio.NewReader(os.Stdin)
	userTEMP, _ := reader.ReadString('\n')
	userTEMP = strings.TrimSuffix(userTEMP, "\n")

	Y.Print("\n     You Typed: ")
	W.Print(userTEMP)
	Y.Println("**")

	return userTEMP

} //end of

func GET_INPUT() string {
	return GET_USER_INPUT()
}



/*
 DownloadFile will download a url to a local file. It's efficient because it will
 write as it downloads and not load the whole file into memory.

  Courtesy of: https://golangcode.com/download-a-file-from-a-url/
*/
func DownloadFile(filepath string, url string) error {

	// Create the file
	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()

	// Get the data
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Write the body to file
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return err
	}

	return nil
}

// Opens a file and returns a file object
func OPEN_FILE(path_to_file string) *os.File {

	file_obj, err := os.Open(path_to_file)
	if err != nil {
		R.Print(" ** ERROR ** Cannot open the file: ")
		W.Println(path_to_file)
		Y.Println(err.Error())
	}

	return file_obj

} //end of func

func AppendStringToFile(path string, text string, OVERWRITE_FILE_FIRST bool) bool {

	if OVERWRITE_FILE_FIRST {
		os.Remove(path)
	}

	f, err := os.OpenFile(path, os.O_RDWR|os.O_APPEND|os.O_CREATE, 0666)
	defer f.Close()

	if err != nil {
		etxt := err.Error()

		R.Println("\n\n **** ERROR Opening for APPENDING FILE ***** check AppendStringToFile")
		Y.Println(etxt)
		return false
	}

	// Writes string to file with NEW LINE
	_, err = f.WriteString(text + "\n")
	if err != nil {
		etxt := err.Error()

		R.Println("\n\n **** ERROR WRITING FILE ***** check AppendStringToFile")
		Y.Println(etxt)

		return false
	}

	return true
}
// Alias to AppendStringToFile
func AppendToFile(path string, text string, OVERWRITE_FILE_FIRST bool) bool {
	return AppendStringToFile(path, text, OVERWRITE_FILE_FIRST)
}
// Alias to AppendStringToFile
func APPEND_FILE(path string, text string, OVERWRITE_FILE_FIRST bool) bool {
	return AppendStringToFile(path, text, OVERWRITE_FILE_FIRST)
}
// Alias to AppendStringToFile
func AppendFile(path string, text string, OVERWRITE_FILE_FIRST bool) bool {
	return AppendStringToFile(path, text, OVERWRITE_FILE_FIRST)
}

func ErrorLog(fname string, SUBJECT string, errOBJ error) {

	//1. Always print the error message on the screen
	ctime, _ := GET_CURRENT_TIME()

	R.Println(" *** ", ctime)
	R.Print(" *** ")
	Y.Println(SUBJECT)
	R.Print(" *** ")
	M.Println(errOBJ.Error())
	R.Println(" *** ")

	AppendStringToFile(fname, " *** "+ctime, false)
	AppendStringToFile(fname, SUBJECT, false)
	AppendStringToFile(fname, errOBJ.Error(), false)
	AppendStringToFile(fname, " *** ", false)

} //end of function


// This checks to see if a file or directory exists
func fileExists(path string) bool {
	_, err := os.Stat(path)
	if err == nil {
		return true
	}
	if os.IsNotExist(err) {
		return false
	}
	return true
} // end of fileExist

func DOES_FILE_EXIST(path string) bool {

	result := fileExists(path)

	return result

} //end of func

// Gets a list of files in a directory
func Get_FILE_LIST(dirname string) []string {

	var results []string

	//1. Read from the specified directory
	files, err := ioutil.ReadDir(dirname)
	if err != nil {
		R.Print(" *** ERROR IN GET FILES LIST: ", err)
		return results
	}

	//2. Now lets add these file names to the directory list

	for _, f := range files {

		results = append(results, f.Name())

	} //end of for loop

	//5. Sort the list alphabetically
	sort.Strings(results)

	return results
} //end of get_File_LIST

func GET_MD5_HASH(inval string) string {

	hasher := md5.New()
	hasher.Write([]byte(inval))
	result := hex.EncodeToString(hasher.Sum(nil))

	// color.Green(" Hey the hasval is: " + result)

	return result
}

// Alias for GET_MD5_HASH
func GENERATE_MD5(inval string) string {
	return GET_MD5_HASH(inval)
}

// Returns a randomly generated number within a given range (returns a STRING AND an int)
func GenRandomRange(min int, max int) (int, string) {

	resultNum := rand.Intn(max-min) + min
	resultText := strconv.Itoa(resultNum)

	return resultNum, resultText

} //end of genRandomRange



var PAGE_COUNT = 0
var PAGE_MAX = 5

// This is a basic Paging routine that prompts you to PressAny key
// after x number of items have been shown
func Pager(tmax int) {
	PAGE_MAX = tmax
	PAGE_COUNT++

	if PAGE_COUNT == PAGE_MAX {
		C.Print("   - - PAGER - -")
		PressAny()
		PAGE_COUNT = 0
	}

} //end of Pager


var PRESS_MESSAGE = ""
// Simple PressAny Key function
func PressAny() {

	defaultMessage := "  ...Press Enter to Continue..."

	if PRESS_MESSAGE != "" {

		defaultMessage = PRESS_MESSAGE
	}

	W.Print("\n\n", defaultMessage, "\n")

	//1. New way of doing PAK
	b := make([]byte, 10)
	if _, err := os.Stdin.Read(b); err != nil {
		R.Println("Fatal error in PressAny Key: ", err)
	}

} // end of func

/*
        b := make([]rune, slen)
        for i := range b {
                b[i] = stringTEMPLATE[rand.Intn(len(stringTEMPLATE))]
        }

        result := "mixed-STRING-" + string(b)
        return result

*/

// This takes IN a string and returns a shuffle of the characters contained in it
func SHUFFLE_STRING(input_STRING string) string {

	//1. Get the length of the string
	slen := len(input_STRING)

	stringRUNE := []rune(input_STRING)

	shuffledString_RESULT := make([]rune, slen)

	for i := range shuffledString_RESULT {
		shuffledString_RESULT[i] = stringRUNE[rand.Intn(slen)]
	}
	return string(shuffledString_RESULT)
} // end of genSESSION

// This gets the platform we are running on (mac, linux, windows) and saves it to OSTYPE
func GET_CURRENT_OS_INFO() {

	if runtime.GOOS == "linux" {
		OSTYPE="Linux"

	//2. Otherwise see if this is MAC
	} else if runtime.GOOS == "darwin" {
		OSTYPE="MAC"

	//3. otherwise.. its windows.. it wins by default!!	
	} else if runtime.GOOS == "windows" {
		OSTYPE="Windows"

	//4. If we get this far, means we have some weird unrecognizable OS:
	} else {
		OSTYPE="- - UNKNOWN OS - -"
	}

	//5. Another courtesy Alias
	CURRENT_OS = OSTYPE
	GOOS_VALUE = runtime.GOOS

} //end of getOsType
func Show_TOTAL_PROG_RUNTIME() {

	endTime, endOBJ := GET_CURRENT_TIME()
	DIFF := DISPLAY_TIME_DIFF(PROG_START_TIMEOBJ, endOBJ)
	//12. DISPPLAY Status on this Threads Performance (and metrics)
	R.Println("")
	W.Print("* * * * * * * * * * * *")
	G.Print(" Script Complete ")
	W.Println("* * * * * * * * * * * * ")

	
	C.Println("")
	W.Print("         STARTED:")
	B.Println(" " + PROG_START_TIME)
	Y.Print("        ENDED on:")
	M.Println(" " + endTime)
	C.Print("  Total DURATION: ")
	G.Println(DIFF)

	W.Println("- - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -")

}


// This is a simplified way of executing external commands... just pass the command and its parameters, it returns the output.. (or the error)
func ComExec(command_INPUT string, VERBOSE_MODE bool) (string, []string) {

	parts := strings.Fields(command_INPUT)
	head := parts[0]
	parts = parts[1:len(parts)]

	if VERBOSE_MODE {
		Y.Print("EXECUTING: ")
		justString := strings.Join(parts, " ")
		W.Print(head + " ")
		Y.Println(justString)

	}

	cmd := exec.Command(head, parts...)
	etext := ""
	out, err := cmd.CombinedOutput()
	if err != nil {

		Y.Println("\n     WARNING: Execution Error!", err.Error())
		Y.Println("")
		etext = err.Error()
	}

	otemp := string(out)
	outText := strings.Split(otemp, "\n")

	outText = append(outText, etext)

	if VERBOSE_MODE {
		M.Println(otemp)
	}

	return otemp, outText

} // end of comExec



// These are the characters used to generate the serial
//var sessTEMPLATE = []rune("GRZBJHUFLEKVXMNTQPSOADWYC527183469")

// This generates a serial.. usually used discern between multiple execution runs like in jenkins
func GenSerial(slen int) {

	result := SHUFFLE_STRING("grzbjhuflcekivxmntqpsoadwy527183469")

	part_ONE := result[0:4]
	part_TWO := result[3:slen]

	SERIAL_NUM = part_ONE + "-" + part_TWO

} // end of GenSerial


func SETUP_DEFAULT_COMMAND_LINE_PARAMS() {

	//1. First lets get some default command line params we always use
	flag.StringVar(&TIME_ZONE_FLAG,       "zone", TIME_ZONE_FLAG,         "  If (cst, pst,est) specified we run THAT time zone")
	flag.BoolVar(&DEBUG_MODE,       "debug", DEBUG_MODE,         "  If specified we runin DEBUG MODE")

	//2. These are convenience variables for the timezone stuff
	var USE_PST = false
	var USE_CST = false
	var USE_EST = false
	var USE_UTC = false
	flag.BoolVar(&USE_PST,       "pst", USE_PST,         "  Force Timezone to be PST")
	flag.BoolVar(&USE_CST,       "cst", USE_CST,         "  Force Timezone to be CST")
	flag.BoolVar(&USE_EST,       "est", USE_EST,         "  Force Timezone to be EST")
	flag.BoolVar(&USE_UTC,       "utc", USE_UTC,         "  Force Timezone to be UTC")

	//3. In case we want to show serial numbers in program runs	
	flag.BoolVar(&SHOW_SERIAL,       "serial", SHOW_SERIAL,         "  If specified we Show a RUNTIME serial number")
	flag.BoolVar(&SHOW_SERIAL,       "useserial", SHOW_SERIAL,         "  alias to --serial")
	flag.BoolVar(&SHOW_SERIAL,       "showserial", SHOW_SERIAL,         "  alias to --serial")
	
	//4. And finally, Very important.. This is the final flag.Parse that is run in the program
	flag.Parse()

	//5. If any of the --pst, --cst were specified, we force TIME_ZONE_FLAG to be that
	if USE_PST {
		TIME_ZONE_FLAG = "pst"

	} else if USE_CST {
		TIME_ZONE_FLAG = "cst"

	} else if USE_EST {
		TIME_ZONE_FLAG = "est"

	} else if USE_UTC {
		TIME_ZONE_FLAG = "utc"
	}

} //end of setup default command line params

// ALWAYS CALL THIS in the MAIN of every program.. This is how command line params get initted
// Also make it the LAST one that is called (before AWS_INIT for example)
func MASTER_INIT(PROGNAME string, Version float64) {

	//1. Setup default COmmand line params
	SETUP_DEFAULT_COMMAND_LINE_PARAMS()

	//2. Lets setup some defaults we found we need for proper Timezone and DateMath operations
	SET_TIMEZONE_DEFAULTS()

	//3. And, Always init the random number seeder
	rand.Seed(time.Now().UTC().UnixNano())	

	//3b. And get current OS Data
	GET_CURRENT_OS_INFO()

	//4. Setup the prog start time globals
	PROG_START_TIME, PROG_START_TIMEOBJ = GET_CURRENT_TIME()
	
	W.Println("")
	M.Println(" ____________________________________________________________")
	M.Println("|")
	M.Print("|  ", PROGNAME, " | ")
	C.Print("Ver: ")	
	Y.Println(Version)

	G.Println("")
	G.Println("   ** ENOUGH WITH TILLY and MICHAEL!!! *****")
	G.Println("")

	//5. Show the OS Information
	M.Print("|  Current OS: ")
	W.Println(CURRENT_OS)
	M.Println("|")

	//6. Courtesy display of TimeZone info
	M.Print("| ")
	SHOW_ZONE_INFO()
	M.Println("|")
	
	//7. If specified we show the unique execution serial. This is useful when running from within Jenkins
	if SHOW_SERIAL {
		GenSerial(10)
		M.Print("| ")
		W.Print(" Exec Serial: ")
		Y.Println(SERIAL_NUM)
	}
	

	//8. And final close of the intro header
	M.Println("|____________________________________________________________")
	W.Println("")

} //end of func



// Kept here as filler/example.. anything you put in this function will start when the module is imported
func init() {

	//1. Startup Stuff (init the command line params etc) . We need these Time ZONE Objects




} // end of main
