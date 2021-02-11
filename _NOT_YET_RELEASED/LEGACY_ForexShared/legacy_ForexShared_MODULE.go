/*	ForexShared Common Functions 55	
	
*/
package ForexShared

import (
	. "TerryCommon"

	"sort"
	"strings"
	"time"
	//	"unicode"
	//"errors"
	"math"
	"os"
	//"fmt"
	"strconv"
	//"sync"

	. "MONGO_Common"

	"gopkg.in/mgo.v2/bson"

)



var TOTAL_MINUTES_DAY = 1440


type PRICE_DATA_RECORDS struct {
	SYMBOL 										string  `bson:"SYMBOL"` // The symbol for the security
	Date   										string  `bson:"Date"`   // just an easier to view DATE... pretty date
	Weekday										string  `bson:"Weekday"`   // i like have day of week here
	Price  										float64 `bson:"Price"`  // Numeric price (based on the CLOSE)
	INTERVAL 									string  `bson:"INTERVAL"` // The INTERVAL (1min, 15min or hourly)
	
	Open     									float64 `bson:"Open"`     // Numeric OPEN price of bar
	High     									float64 `bson:"High"`     // Numeric HIGH of bar
	Low      									float64 `bson:"Low"`      // Numeric LOW of bar
	Close    									float64 `bson:"Close"`    // Numeric Close of bar (same as PRICE above)
	BarColor 									string  `bson:"BarColor"` // the "bar color": Red or Green based on open > close or vice versa

	ACTUAL_VOLUME								float64  `bson:"ACTUAL_VOLUME"` // Actual Volume as reported in the import file
	MAGIC_VOLUME								int  `bson:"MAGIC_VOLUME"` // Approx volume of bar (using the magic method i made up)	

	DEBUG_DAYS_of_data_sofar				    int  `bson:"DEBUG_DAYS_of_data_sofar"`    // Numeric price (based on the CLOSE)
	DEBUG_YEARS_of_data_sofar				float64 `bson:"DEBUG_YEARS_of_data_sofar"`    // Numeric price (based on the CLOSE)
	
	/*
		JUST the Month/Day/Year of this bar.. 
		helpful for writing validation test (ie to make sure we have at least 1440 minute bars in a day 
		(or 24 hour bars in a day if we are doing an hourly interval
	*/
	DEBUG_Mon_Day_Year						string  `bson:"DEBUG_Mon_Day_Year"`  
	
	/* 
		Name of the file from which we imported this record
		Helpful for doing incremental additions to the price Data base.. 
		Note we are using a unique index so we dont have to worry about duplicate entries
	*/
	DEBUG_SOURCE_FILE							string  `bson:"DEBUG_SOURCE_FILE"` 

	DATE_OBJ 									time.Time 	`bson:"DATE_OBJ"`

}




type VALID_SYM_REC struct {
	SYMBOL             		string 
	REACTION_ID					string
	LOOK_AHEAD_DATEOBJ		time.Time
	LOOK_PRETTY			string
}




type CHANGE_SUMMARY_OBJ struct {

	CHANGE_Actual_vs_Prev      string  `bson:"CHANGE_Actual_vs_Prev"`  // Type of Change between Actual and PREV (Increase or Decrease)
	CHANGE_Forecast_vs_Prev   	string  `bson:"CHANGE_Forecast_vs_Prev"`    // Type of Change between Forecast vs PREV (Increase or Decrease)
	CHANGE_Actual_vs_Forecast	string  `bson:"CHANGE_Actual_vs_Forecast"`  // Type of Change between Forecast vs Actual 

	PERCENT_Actual_vs_Prev 	   float64 `bson:"PERCENT_Actual_vs_Prev"` // NUMERIC value in Percentage of Actual vs Prev
	PERCENT_Forecast_vs_Prev 	float64 `bson:"PERCENT_Forecast_vs_Prev"` 	 // NUMERIC value in Percentage of Forecast vs Prev
	PERCENT_Actual_vs_Forecast	float64 `bson:"PERCENT_Actual_vs_Forecast"` // NUMERIC Percentage of Forecast vs Actual

	Prev   						 string `bson:"Prev"`   // This is the PREVIOUS value for this news item (as scraped)
	Actual 						 string `bson:"Actual"` // This is the ACTUAL value for news item  (as scraped)
	Forecast   					 string `bson:"Forecast"`   // This is the N_Forecast value for news item  (as scraped)	
}


type NEWS_META_OBJ struct {

			/* Found i needed this to assist with ensuring the correct time and timezone being stored
					as it is reported from the SITE i am scraping from
					... also for peace of mind. This is usually displayed as GMT something or other
			*/
	Source_TIMEZONE	 string    `bson:"Source_TIMEZONE"`

	Time_RETRIEVED		 string    `bson:"Time_RETRIEVED"`

	SANITIZED_Report 	 string 	  `bson:"SANITIZED_Report"` // "CLEANED" report Name.. We MAY use this instead of originla

}

// This is the format of NEWS "documents/records" stored in the mongo database
type NEWS_RECORDS struct {
	SOURCE          			 string `bson:"SOURCE"`           // This is the news source we scrape from (invest or dailyfx)
	Date    		      		 string `bson:"Date"`          // Pretty date of the time of this report		
	Report          			 string `bson:"Report"`           // Original Name of report (as scraped)
	Currency		 				 string `bson:"Currency"` 			// The Currency the report is associated with

	Weekday       				 string    `bson:"Weekday"`       // Day of week
	DateOBJ				   	 time.Time `bson:"DateOBJ"`    // This is the Mongo Date object for searching	
	
	Prev   						 string `bson:"Prev"`   // This is the PREVIOUS value for this news item (as scraped)
	Actual 						 string `bson:"Actual"` // This is the ACTUAL value for news item  (as scraped)
	Forecast   					 string `bson:"Forecast"`   // This is the N_Forecast value for news item  (as scraped)

	CHANGE_SUMMARY	 			 CHANGE_SUMMARY_OBJ	`bson:"CHANGE_SUMMARY"`   // The various Change metrics

	META                     NEWS_META_OBJ 		`bson:"META"`   // Contains misc META data for the News Record
	
}


type REACTION_META_OBJ struct {
	SYMBOL 					string	`bson:"SYMBOL"`
	INTERVAL					string	`bson:"INTERVAL"`

	BASE_PRICE 				float64	`bson:"BASE_PRICE"`
	HIGHEST    				float64	`bson:"HIGHEST"`
	LOWEST     				float64	`bson:"LOWEST"`

	HIGHEST_BAR_TIME 		string	`bson:"HIGHEST_BAR_TIME"`		// Time / date object of HIGHEST bar
	LOWEST_BAR_TIME		string	`bson:"LOWEST_BAR_TIME"`		// Time / date object of LOWEST bar

	BARS_TO_HIGH			int		`bson:"BARS_TO_HIGH"`	 // Bars until we reach the highest point in a price window
	BARS_TO_LOW				int		`bson:"BARS_TO_LOW"`  // Bars until we reach the LOWEST point in a price window
	
	ABOVE_SUMMARY			string	`bson:"ABOVE_SUMMARY"`
	BELOW_SUMMARY			string	`bson:"BELOW_SUMMARY"`

	PERCENT_ABOVE 			float64	`bson:"PERCENT_ABOVE"`
	PERCENT_BELOW 			float64	`bson:"PERCENT_BELOW"`

	PIPS_ABOVE				int		`bson:"PIPS_ABOVE"`
	PIPS_BELOW				int		`bson:"PIPS_BELOW"`
	PIP_DIFFER				int		`bson:"PIP_DIFFER"`

	ABOVE_VOLUME			int		`bson:"ABOVE_VOLUME"`
	BELOW_VOLUME			int		`bson:"BELOW_VOLUME"`
	TOTAL_VOLUME			int		`bson:"TOTAL_VOLUME"`
	
	DEBUG_Window_START				string `bson:"DEBUG_Window_START"`
	DEBUG_Window_END					string `bson:"DEBUG_Window_END"`
	DEBUG_BAR_SUMMARY					string `bson:"DEBUG_BAR_SUMMARY"`
	DEBUG_TOTAL_BARS_in_WINDOW		int `bson:"DEBUG_TOTAL_BARS_in_WINDOW"`
}


/* This is the PRAE Reaction Collection Record
This is an evaluation of price movement (in pips) based on the time of the report for a specific currency
*/

type PRAE_REACTION_RECORDS struct {
	REACTION_ID					string		`bson:"REACTION_ID"`		// a unique ID generated to identify this reaction
	SYMBOL             		string 		`bson:"SYMBOL"`	
	
	Report             		string 		`bson:"Report"`            // Name of the Report
	Currency      				string		`bson:"Currency"`     // Currency SYMBOL associated with this report (as reported by the Source
	NEWS_SOURCE					string		`bson:"NEWS_SOURCE"`
	Date      			 		string 		`bson:"Date"`       // The DATE of the report (prettified)
	

	Weekday            		string 	 	`bson:"Weekday"`           //Weekday

	SANITIZED_Report   		string 		`bson:"SANITIZED_Report"`  // SANITIZED_Report of the Report name
		
	DATE_OBJ				 		time.Time 				`bson:"DATE_OBJ"` 		// The mongo DATE / TIME object (ver of 	
	
	CHANGE_SUMMARY	 			CHANGE_SUMMARY_OBJ	`bson:"CHANGE_SUMMARY"`   // The various Change metrics of the news REPORT

	REACTION_SUMMARY			REACTION_META_OBJ	   `bson:"REACTION_SUMMARY"`   // This is the close price at the time th
}



/*
	#####
	#####		Threading routines for PRAE ENGINE
	#####
*/








/*
	#############
	#############		Common functions Start HERE
	#############
*/

// Returns the number of decimal places in a float
func GET_NUM_DECIMAL_PLACES(v float64) int {
    s := strconv.FormatFloat(v, 'f', -1, 64)
    i := strings.IndexByte(s, '.')
    if i > -1 {
        return len(s) - i - 1
    }
    return 0
} //edn of func





func GET_SECURITY_TYPE(msymbol string) string {

	var temp_sym = strings.ToUpper(msymbol)

	var FOUND = ""

	for _, mx := range (SYMBOL_MATRIX) {

		if mx.SYMBOL == temp_sym {
			FOUND = mx.TYPE
			break
		}
	}

	return FOUND
} //end of


// THis returns an appropriately formatted date string as needed by the multiple news sources 
func FORMAT_SCRAPER_date(inDate time.Time, dtype string) string {

	//1. extract the month day and year
	// i_month := time.Date(inDate)

	tYear, month, tDay := inDate.Date()
	
	tMonth := int(month)

	var monText = strconv.Itoa(tMonth)
	var dayText = strconv.Itoa(tDay)

	var yearText = strconv.Itoa(tYear)

	if tMonth < 10 {
		monText = "0" + monText
	}

	if tDay < 10 {
		dayText = "0" + dayText
	}

	// Now depending onwhat dtype is, we return a string for DailyFX or.. InVESTING.com

	if dtype == "DAILYFX" {

		// DailyFX expects dates to look like this:
		// 2017/0108
		result := yearText + "/" + monText + dayText
		return result

	} 
	
	if dtype == "INVESTING" {

		// Invest expects dates to look like this:
		// dateFrom=2017-01-03

		result := yearText + "-" + monText + "-" + dayText
		return result
	}



	if dtype == "BRIEFING" {
		//BRIEFING expects a date to look like this:  2017/12/06

		result := yearText + "/" + monText + "/" + dayText
		return result
	}	
	

	// forexfactory
	if dtype == "FOREXFACTORY" {
		/*
			forexfactory expects a date to look like this:  

			dec08.2017

		*/

		month_NAME_TEXT := "jan"

		if (monText == "02") {
			month_NAME_TEXT = "feb"

		} else if (monText == "03") {
			month_NAME_TEXT = "mar"

		} else if (monText == "04") {
			month_NAME_TEXT = "apr"

		} else if (monText == "05") {
			month_NAME_TEXT = "may"

		} else if (monText == "06") {
			month_NAME_TEXT = "jun"

		} else if (monText == "07") {
			month_NAME_TEXT = "jul"

		} else if (monText == "08") {
			month_NAME_TEXT = "aug"

		} else if (monText == "09") {
			month_NAME_TEXT = "sep"

		} else if (monText == "10") {
			month_NAME_TEXT = "oct"

		} else if (monText == "11") {
			month_NAME_TEXT = "nov"

		} else if (monText == "12") {
			month_NAME_TEXT = "dec"

		}


		result := month_NAME_TEXT + dayText + "." + yearText
		return result
	}	
	

	return ""

} //end of FUNCTION



func PULL_NEWS_REPORTS(REPORTS_AFTER_DATE string) []NEWS_RECORDS {

	c := DBSession.DB(DBName).C("NewsReports")

	/*
		For some reason, sortCriteria doesnt work on time.Time (i assume becasue it isnt a string)
		
	*/
	W.Println("\n   *** PULLING NEWS REPORTS ")

	// c.Find(bson.M{"name": "Ale"}).Sort("-timestamp").All(&results)
	//1. If REPORTS_AFTER_DATE is NOT blank, pull reports after this date ONLY

	DEBUG_LIMIT := 0

	var TEMP_LIST []NEWS_RECORDS
	if REPORTS_AFTER_DATE != "" {
		after_DATEOBJ := CONVERT_DATE(REPORTS_AFTER_DATE)
		c.Find(bson.M{
			"DateOBJ": bson.M{
				"$gte": after_DATEOBJ,
			},
		}).Limit(DEBUG_LIMIT).All(&TEMP_LIST)

	//2. Otherwise, pull everything
	} else {
		c.Find(nil).Limit(DEBUG_LIMIT).All(&TEMP_LIST)
	}

	RETURNED_LIST := TEMP_LIST

	rlen := len(RETURNED_LIST)

	M.Print(" Found a Total of: ")
	G.Print(ShowNum(rlen))
	M.Println(" News Reports! ..")

	if rlen == 0 {

		R.Println(" ERROR: For some reason I cant find any news reports.. I have to stop")
		os.Exit(1)

	}
	
	// // This sorts teh slice
	// W.Println("   *** Now SORTING News Reports")
	sort.Slice(RETURNED_LIST, func(i, j int) bool { return (RETURNED_LIST)[i].DateOBJ.Before((RETURNED_LIST)[j].DateOBJ)})	

	return RETURNED_LIST
	
} //end of func
 

/*
-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-= FUNCTIONS BELOW HERE -=-=-=-=-=-=-=-=-=-=-=-=
*/


type HOLI_OBJ struct {
	Date string
	DESC string
}

var HOLIDAY_LIST []HOLI_OBJ

var INDY = "Independence Day"
var LABOR = "Labor Day"
var MEM = "Memorial DAY"
var PRES = "Presdents Day"
var MARTY = "Martin Luther King Day"
var NEW = "New Years Day"

var THANK = "Thanks Giving"
var THANKS = "Thanks Giving"
var TANK = "Thanks Giving"
var TANKS = "Thanks Giving"

var CHRIS = "Christmas day"
var GOOF = "Good Friday"
var EAST = "Easter"
var VETS = "Veterans Day"

var COLUM = "Columbus Day"
var COLUMB = "Columbus Day"
var BOXD = "Boxing Day"
var THANK_BLACK_FRIDAY = "Day After Thanksgiving BLACK FRIDAY"
var WASH = "Washingtons Birthday"

func INIT_HOLIDAY_LIST() {

	/*

		Last Revised 11/27/2018

	 SOURCE(s):
	 	https://www.sifma.org/wp-content/uploads/2017/06/misc-us-historical-holiday-market-recommendations-sifma.pdf
	 	https://www.calendar-365.com/holidays/2016.html


		HOLIDAY_LIST = append(HOLIDAY_LIST, HOLI_OBJ{Date: "XXXXXXXX", DESC: NEW})
		HOLIDAY_LIST = append(HOLIDAY_LIST, HOLI_OBJ{Date: "XXXXXXXX", DESC: MARTY})
		HOLIDAY_LIST = append(HOLIDAY_LIST, HOLI_OBJ{Date: "XXXXXXXX", DESC: PRES})
		HOLIDAY_LIST = append(HOLIDAY_LIST, HOLI_OBJ{Date: "XXXXXXXX", DESC: MEM})
		HOLIDAY_LIST = append(HOLIDAY_LIST, HOLI_OBJ{Date: "XXXXXXXX", DESC: INDY})
		HOLIDAY_LIST = append(HOLIDAY_LIST, HOLI_OBJ{Date: "XXXXXXXX", DESC: LABOR})
		HOLIDAY_LIST = append(HOLIDAY_LIST, HOLI_OBJ{Date: "XXXXXXXX", DESC: COLUM})
		HOLIDAY_LIST = append(HOLIDAY_LIST, HOLI_OBJ{Date: "XXXXXXXX", DESC: VETS})
		HOLIDAY_LIST = append(HOLIDAY_LIST, HOLI_OBJ{Date: "XXXXXXXX", DESC: THANK})
		HOLIDAY_LIST = append(HOLIDAY_LIST, HOLI_OBJ{Date: "XXXXXXXX", DESC: CHRIS})

	*/

	HOLIDAY_LIST = append(HOLIDAY_LIST, HOLI_OBJ{Date: "01/01/2003", DESC: NEW})
	HOLIDAY_LIST = append(HOLIDAY_LIST, HOLI_OBJ{Date: "01/17/2003", DESC: MARTY})
	HOLIDAY_LIST = append(HOLIDAY_LIST, HOLI_OBJ{Date: "01/20/2003", DESC: MARTY})	
	HOLIDAY_LIST = append(HOLIDAY_LIST, HOLI_OBJ{Date: "02/14/2003", DESC: PRES})
	HOLIDAY_LIST = append(HOLIDAY_LIST, HOLI_OBJ{Date: "02/17/2003", DESC: PRES})
	HOLIDAY_LIST = append(HOLIDAY_LIST, HOLI_OBJ{Date: "04/17/2003", DESC: GOOF})
	HOLIDAY_LIST = append(HOLIDAY_LIST, HOLI_OBJ{Date: "04/18/2003", DESC: GOOF})
	HOLIDAY_LIST = append(HOLIDAY_LIST, HOLI_OBJ{Date: "05/23/2003", DESC: MEM})	
	HOLIDAY_LIST = append(HOLIDAY_LIST, HOLI_OBJ{Date: "05/26/2003", DESC: MEM})
	HOLIDAY_LIST = append(HOLIDAY_LIST, HOLI_OBJ{Date: "07/03/2003", DESC: INDY})
	HOLIDAY_LIST = append(HOLIDAY_LIST, HOLI_OBJ{Date: "07/04/2003", DESC: INDY})
	HOLIDAY_LIST = append(HOLIDAY_LIST, HOLI_OBJ{Date: "08/29/2003", DESC: LABOR})
	HOLIDAY_LIST = append(HOLIDAY_LIST, HOLI_OBJ{Date: "09/01/2003", DESC: LABOR})
	HOLIDAY_LIST = append(HOLIDAY_LIST, HOLI_OBJ{Date: "10/10/2003", DESC: COLUM})	
	HOLIDAY_LIST = append(HOLIDAY_LIST, HOLI_OBJ{Date: "10/13/2003", DESC: COLUM})
	HOLIDAY_LIST = append(HOLIDAY_LIST, HOLI_OBJ{Date: "11/10/2003", DESC: VETS})	
	HOLIDAY_LIST = append(HOLIDAY_LIST, HOLI_OBJ{Date: "11/11/2003", DESC: VETS})
	HOLIDAY_LIST = append(HOLIDAY_LIST, HOLI_OBJ{Date: "11/26/2003", DESC: THANK})
	HOLIDAY_LIST = append(HOLIDAY_LIST, HOLI_OBJ{Date: "11/27/2003", DESC: THANK})
	HOLIDAY_LIST = append(HOLIDAY_LIST, HOLI_OBJ{Date: "11/28/2003", DESC: THANK})
	HOLIDAY_LIST = append(HOLIDAY_LIST, HOLI_OBJ{Date: "12/24/2003", DESC: CHRIS})
	HOLIDAY_LIST = append(HOLIDAY_LIST, HOLI_OBJ{Date: "12/25/2003", DESC: CHRIS})
	HOLIDAY_LIST = append(HOLIDAY_LIST, HOLI_OBJ{Date: "12/31/2003", DESC: CHRIS})


HOLIDAY_LIST = append(HOLIDAY_LIST, HOLI_OBJ{Date: "01/01/2004", DESC: NEW})
HOLIDAY_LIST = append(HOLIDAY_LIST, HOLI_OBJ{Date: "01/02/2004", DESC: NEW})
HOLIDAY_LIST = append(HOLIDAY_LIST, HOLI_OBJ{Date: "01/16/2004", DESC: MARTY})
HOLIDAY_LIST = append(HOLIDAY_LIST, HOLI_OBJ{Date: "01/19/2004", DESC: MARTY})
HOLIDAY_LIST = append(HOLIDAY_LIST, HOLI_OBJ{Date: "02/13/2004", DESC: PRES})
HOLIDAY_LIST = append(HOLIDAY_LIST, HOLI_OBJ{Date: "02/16/2004", DESC: PRES})
HOLIDAY_LIST = append(HOLIDAY_LIST, HOLI_OBJ{Date: "04/08/2004", DESC: GOOF})
HOLIDAY_LIST = append(HOLIDAY_LIST, HOLI_OBJ{Date: "04/09/2004", DESC: GOOF})
HOLIDAY_LIST = append(HOLIDAY_LIST, HOLI_OBJ{Date: "05/28/2004", DESC: MEM})
HOLIDAY_LIST = append(HOLIDAY_LIST, HOLI_OBJ{Date: "05/31/2004", DESC: MEM})
HOLIDAY_LIST = append(HOLIDAY_LIST, HOLI_OBJ{Date: "06/11/2004", DESC: PRES})  // Ronald reagan passed away.. day of mourning/observance
HOLIDAY_LIST = append(HOLIDAY_LIST, HOLI_OBJ{Date: "07/02/2004", DESC: INDY})
HOLIDAY_LIST = append(HOLIDAY_LIST, HOLI_OBJ{Date: "07/05/2004", DESC: INDY})
HOLIDAY_LIST = append(HOLIDAY_LIST, HOLI_OBJ{Date: "09/03/2004", DESC: LABOR})
HOLIDAY_LIST = append(HOLIDAY_LIST, HOLI_OBJ{Date: "09/06/2004", DESC: LABOR})
HOLIDAY_LIST = append(HOLIDAY_LIST, HOLI_OBJ{Date: "10/08/2004", DESC: COLUM})
HOLIDAY_LIST = append(HOLIDAY_LIST, HOLI_OBJ{Date: "10/11/2004", DESC: COLUM})
HOLIDAY_LIST = append(HOLIDAY_LIST, HOLI_OBJ{Date: "11/11/2004", DESC: VETS})
HOLIDAY_LIST = append(HOLIDAY_LIST, HOLI_OBJ{Date: "11/24/2004", DESC: THANKS})
HOLIDAY_LIST = append(HOLIDAY_LIST, HOLI_OBJ{Date: "11/25/2004", DESC: THANKS})
HOLIDAY_LIST = append(HOLIDAY_LIST, HOLI_OBJ{Date: "11/26/2004", DESC: THANKS})
HOLIDAY_LIST = append(HOLIDAY_LIST, HOLI_OBJ{Date: "12/23/2004", DESC: CHRIS})
HOLIDAY_LIST = append(HOLIDAY_LIST, HOLI_OBJ{Date: "12/24/2004", DESC: CHRIS})
HOLIDAY_LIST = append(HOLIDAY_LIST, HOLI_OBJ{Date: "12/31/2004", DESC: CHRIS})

	HOLIDAY_LIST = append(HOLIDAY_LIST, HOLI_OBJ{Date: "01/01/2005", DESC: NEW})
	HOLIDAY_LIST = append(HOLIDAY_LIST, HOLI_OBJ{Date: "01/14/2005", DESC: MARTY})
	HOLIDAY_LIST = append(HOLIDAY_LIST, HOLI_OBJ{Date: "01/17/2005", DESC: MARTY})
	HOLIDAY_LIST = append(HOLIDAY_LIST, HOLI_OBJ{Date: "02/18/2005", DESC: PRES})
	HOLIDAY_LIST = append(HOLIDAY_LIST, HOLI_OBJ{Date: "02/21/2005", DESC: PRES})
	HOLIDAY_LIST = append(HOLIDAY_LIST, HOLI_OBJ{Date: "03/24/2005", DESC: EAST})
	HOLIDAY_LIST = append(HOLIDAY_LIST, HOLI_OBJ{Date: "03/25/2005", DESC: EAST})
	HOLIDAY_LIST = append(HOLIDAY_LIST, HOLI_OBJ{Date: "05/27/2005", DESC: MEM})
	HOLIDAY_LIST = append(HOLIDAY_LIST, HOLI_OBJ{Date: "05/30/2005", DESC: MEM})
	HOLIDAY_LIST = append(HOLIDAY_LIST, HOLI_OBJ{Date: "07/01/2005", DESC: INDY})
	HOLIDAY_LIST = append(HOLIDAY_LIST, HOLI_OBJ{Date: "07/04/2005", DESC: INDY})
	HOLIDAY_LIST = append(HOLIDAY_LIST, HOLI_OBJ{Date: "09/02/2005", DESC: LABOR})
	HOLIDAY_LIST = append(HOLIDAY_LIST, HOLI_OBJ{Date: "09/05/2005", DESC: LABOR})
	HOLIDAY_LIST = append(HOLIDAY_LIST, HOLI_OBJ{Date: "10/07/2005", DESC: COLUM})
	HOLIDAY_LIST = append(HOLIDAY_LIST, HOLI_OBJ{Date: "10/10/2005", DESC: COLUM})
	HOLIDAY_LIST = append(HOLIDAY_LIST, HOLI_OBJ{Date: "11/10/2005", DESC: VETS})
	HOLIDAY_LIST = append(HOLIDAY_LIST, HOLI_OBJ{Date: "11/11/2005", DESC: VETS})
	HOLIDAY_LIST = append(HOLIDAY_LIST, HOLI_OBJ{Date: "11/23/2005", DESC: THANK})
	HOLIDAY_LIST = append(HOLIDAY_LIST, HOLI_OBJ{Date: "11/24/2005", DESC: THANK})
	HOLIDAY_LIST = append(HOLIDAY_LIST, HOLI_OBJ{Date: "11/25/2005", DESC: THANK})
	HOLIDAY_LIST = append(HOLIDAY_LIST, HOLI_OBJ{Date: "12/23/2005", DESC: CHRIS})
	HOLIDAY_LIST = append(HOLIDAY_LIST, HOLI_OBJ{Date: "12/26/2005", DESC: CHRIS})
	HOLIDAY_LIST = append(HOLIDAY_LIST, HOLI_OBJ{Date: "12/30/2005", DESC: CHRIS})

HOLIDAY_LIST = append(HOLIDAY_LIST, HOLI_OBJ{Date: "01/02/2006", DESC: NEW})
HOLIDAY_LIST = append(HOLIDAY_LIST, HOLI_OBJ{Date: "01/13/2006", DESC: MARTY})
HOLIDAY_LIST = append(HOLIDAY_LIST, HOLI_OBJ{Date: "01/16/2006", DESC: MARTY})
HOLIDAY_LIST = append(HOLIDAY_LIST, HOLI_OBJ{Date: "02/17/2006", DESC: PRES})
HOLIDAY_LIST = append(HOLIDAY_LIST, HOLI_OBJ{Date: "02/20/2006", DESC: PRES})
HOLIDAY_LIST = append(HOLIDAY_LIST, HOLI_OBJ{Date: "04/13/2006", DESC: GOOF})
HOLIDAY_LIST = append(HOLIDAY_LIST, HOLI_OBJ{Date: "04/14/2006", DESC: GOOF})
HOLIDAY_LIST = append(HOLIDAY_LIST, HOLI_OBJ{Date: "05/26/2006", DESC: MEM})
HOLIDAY_LIST = append(HOLIDAY_LIST, HOLI_OBJ{Date: "05/29/2006", DESC: MEM})
HOLIDAY_LIST = append(HOLIDAY_LIST, HOLI_OBJ{Date: "07/03/2006", DESC: INDY})
HOLIDAY_LIST = append(HOLIDAY_LIST, HOLI_OBJ{Date: "07/04/2006", DESC: INDY})
HOLIDAY_LIST = append(HOLIDAY_LIST, HOLI_OBJ{Date: "09/01/2006", DESC: LABOR})
HOLIDAY_LIST = append(HOLIDAY_LIST, HOLI_OBJ{Date: "09/04/2006", DESC: LABOR})
HOLIDAY_LIST = append(HOLIDAY_LIST, HOLI_OBJ{Date: "10/06/2006", DESC: COLUM})
HOLIDAY_LIST = append(HOLIDAY_LIST, HOLI_OBJ{Date: "10/09/2006", DESC: COLUM})
HOLIDAY_LIST = append(HOLIDAY_LIST, HOLI_OBJ{Date: "11/22/2006", DESC: THANK})
HOLIDAY_LIST = append(HOLIDAY_LIST, HOLI_OBJ{Date: "11/23/2006", DESC: THANK})
HOLIDAY_LIST = append(HOLIDAY_LIST, HOLI_OBJ{Date: "11/24/2006", DESC: THANK})
HOLIDAY_LIST = append(HOLIDAY_LIST, HOLI_OBJ{Date: "12/22/2006", DESC: CHRIS})
HOLIDAY_LIST = append(HOLIDAY_LIST, HOLI_OBJ{Date: "12/25/2006", DESC: CHRIS})
HOLIDAY_LIST = append(HOLIDAY_LIST, HOLI_OBJ{Date: "12/29/2006", DESC: NEW})

	HOLIDAY_LIST = append(HOLIDAY_LIST, HOLI_OBJ{Date: "01/01/2007", DESC: NEW})
	HOLIDAY_LIST = append(HOLIDAY_LIST, HOLI_OBJ{Date: "01/12/2007", DESC: MARTY})
	HOLIDAY_LIST = append(HOLIDAY_LIST, HOLI_OBJ{Date: "01/15/2007", DESC: MARTY})
	HOLIDAY_LIST = append(HOLIDAY_LIST, HOLI_OBJ{Date: "02/16/2007", DESC: PRES})
	HOLIDAY_LIST = append(HOLIDAY_LIST, HOLI_OBJ{Date: "02/19/2007", DESC: PRES})
	HOLIDAY_LIST = append(HOLIDAY_LIST, HOLI_OBJ{Date: "04/06/2007", DESC: GOOF})	
	HOLIDAY_LIST = append(HOLIDAY_LIST, HOLI_OBJ{Date: "05/25/2007", DESC: MEM})
	HOLIDAY_LIST = append(HOLIDAY_LIST, HOLI_OBJ{Date: "05/28/2007", DESC: MEM})
	HOLIDAY_LIST = append(HOLIDAY_LIST, HOLI_OBJ{Date: "07/03/2007", DESC: INDY})
	HOLIDAY_LIST = append(HOLIDAY_LIST, HOLI_OBJ{Date: "07/04/2007", DESC: INDY})
	HOLIDAY_LIST = append(HOLIDAY_LIST, HOLI_OBJ{Date: "08/31/2007", DESC: LABOR})
	HOLIDAY_LIST = append(HOLIDAY_LIST, HOLI_OBJ{Date: "09/03/2007", DESC: LABOR})
	HOLIDAY_LIST = append(HOLIDAY_LIST, HOLI_OBJ{Date: "10/05/2007", DESC: COLUM})
	HOLIDAY_LIST = append(HOLIDAY_LIST, HOLI_OBJ{Date: "10/08/2007", DESC: COLUM})
	HOLIDAY_LIST = append(HOLIDAY_LIST, HOLI_OBJ{Date: "11/09/2007", DESC: VETS})
	HOLIDAY_LIST = append(HOLIDAY_LIST, HOLI_OBJ{Date: "11/12/2007", DESC: VETS})
	HOLIDAY_LIST = append(HOLIDAY_LIST, HOLI_OBJ{Date: "11/21/2007", DESC: THANK})
	HOLIDAY_LIST = append(HOLIDAY_LIST, HOLI_OBJ{Date: "11/22/2007", DESC: THANK})
	HOLIDAY_LIST = append(HOLIDAY_LIST, HOLI_OBJ{Date: "11/23/2007", DESC: THANK})
	HOLIDAY_LIST = append(HOLIDAY_LIST, HOLI_OBJ{Date: "12/24/2007", DESC: CHRIS})
	HOLIDAY_LIST = append(HOLIDAY_LIST, HOLI_OBJ{Date: "12/25/2007", DESC: CHRIS})
	HOLIDAY_LIST = append(HOLIDAY_LIST, HOLI_OBJ{Date: "12/31/2007", DESC: NEW})

HOLIDAY_LIST = append(HOLIDAY_LIST, HOLI_OBJ{Date: "01/01/2008", DESC: NEW})
HOLIDAY_LIST = append(HOLIDAY_LIST, HOLI_OBJ{Date: "01/18/2008", DESC: MARTY})
HOLIDAY_LIST = append(HOLIDAY_LIST, HOLI_OBJ{Date: "01/21/2008", DESC: MARTY})
HOLIDAY_LIST = append(HOLIDAY_LIST, HOLI_OBJ{Date: "02/15/2008", DESC: PRES})
HOLIDAY_LIST = append(HOLIDAY_LIST, HOLI_OBJ{Date: "02/18/2008", DESC: PRES})
HOLIDAY_LIST = append(HOLIDAY_LIST, HOLI_OBJ{Date: "03/20/2008", DESC: GOOF})
HOLIDAY_LIST = append(HOLIDAY_LIST, HOLI_OBJ{Date: "03/21/2008", DESC: GOOF})
HOLIDAY_LIST = append(HOLIDAY_LIST, HOLI_OBJ{Date: "05/23/2008", DESC: MEM})
HOLIDAY_LIST = append(HOLIDAY_LIST, HOLI_OBJ{Date: "05/26/2008", DESC: MEM})
HOLIDAY_LIST = append(HOLIDAY_LIST, HOLI_OBJ{Date: "07/03/2008", DESC: INDY})
HOLIDAY_LIST = append(HOLIDAY_LIST, HOLI_OBJ{Date: "07/04/2008", DESC: INDY})
HOLIDAY_LIST = append(HOLIDAY_LIST, HOLI_OBJ{Date: "08/29/2008", DESC: LABOR})
HOLIDAY_LIST = append(HOLIDAY_LIST, HOLI_OBJ{Date: "09/01/2008", DESC: LABOR})
HOLIDAY_LIST = append(HOLIDAY_LIST, HOLI_OBJ{Date: "10/10/2008", DESC: COLUM})
HOLIDAY_LIST = append(HOLIDAY_LIST, HOLI_OBJ{Date: "10/13/2008", DESC: COLUM})
HOLIDAY_LIST = append(HOLIDAY_LIST, HOLI_OBJ{Date: "11/26/2008", DESC: VETS})
HOLIDAY_LIST = append(HOLIDAY_LIST, HOLI_OBJ{Date: "11/27/2008", DESC: THANK})
HOLIDAY_LIST = append(HOLIDAY_LIST, HOLI_OBJ{Date: "11/28/2008", DESC: THANK})
HOLIDAY_LIST = append(HOLIDAY_LIST, HOLI_OBJ{Date: "12/24/2008", DESC: CHRIS})
HOLIDAY_LIST = append(HOLIDAY_LIST, HOLI_OBJ{Date: "12/25/2008", DESC: CHRIS})
HOLIDAY_LIST = append(HOLIDAY_LIST, HOLI_OBJ{Date: "12/31/2008", DESC: NEW})

	HOLIDAY_LIST = append(HOLIDAY_LIST, HOLI_OBJ{Date: "01/01/2009", DESC: NEW})
	HOLIDAY_LIST = append(HOLIDAY_LIST, HOLI_OBJ{Date: "01/16/2009", DESC: MARTY})
	HOLIDAY_LIST = append(HOLIDAY_LIST, HOLI_OBJ{Date: "01/19/2009", DESC: MARTY})
	HOLIDAY_LIST = append(HOLIDAY_LIST, HOLI_OBJ{Date: "02/13/2009", DESC: PRES})
	HOLIDAY_LIST = append(HOLIDAY_LIST, HOLI_OBJ{Date: "02/16/2009", DESC: PRES})
	HOLIDAY_LIST = append(HOLIDAY_LIST, HOLI_OBJ{Date: "04/09/2009", DESC: GOOF})
	HOLIDAY_LIST = append(HOLIDAY_LIST, HOLI_OBJ{Date: "04/10/2009", DESC: GOOF})
	HOLIDAY_LIST = append(HOLIDAY_LIST, HOLI_OBJ{Date: "04/13/2009", DESC: EAST})
	HOLIDAY_LIST = append(HOLIDAY_LIST, HOLI_OBJ{Date: "05/22/2009", DESC: MEM})
	HOLIDAY_LIST = append(HOLIDAY_LIST, HOLI_OBJ{Date: "05/25/2009", DESC: MEM})
	HOLIDAY_LIST = append(HOLIDAY_LIST, HOLI_OBJ{Date: "07/03/2009", DESC: INDY})
	HOLIDAY_LIST = append(HOLIDAY_LIST, HOLI_OBJ{Date: "09/07/2009", DESC: LABOR})
	HOLIDAY_LIST = append(HOLIDAY_LIST, HOLI_OBJ{Date: "10/12/2009", DESC: COLUM})
	HOLIDAY_LIST = append(HOLIDAY_LIST, HOLI_OBJ{Date: "11/11/2009", DESC: VETS})
	HOLIDAY_LIST = append(HOLIDAY_LIST, HOLI_OBJ{Date: "11/26/2009", DESC: THANK})
	HOLIDAY_LIST = append(HOLIDAY_LIST, HOLI_OBJ{Date: "11/27/2009", DESC: THANK})
	HOLIDAY_LIST = append(HOLIDAY_LIST, HOLI_OBJ{Date: "12/24/2009", DESC: CHRIS})
	HOLIDAY_LIST = append(HOLIDAY_LIST, HOLI_OBJ{Date: "12/25/2009", DESC: CHRIS})
	HOLIDAY_LIST = append(HOLIDAY_LIST, HOLI_OBJ{Date: "12/31/2009", DESC: NEW})


HOLIDAY_LIST = append(HOLIDAY_LIST, HOLI_OBJ{Date: "01/01/2010", DESC: NEW})
HOLIDAY_LIST = append(HOLIDAY_LIST, HOLI_OBJ{Date: "01/18/2010", DESC: MARTY})
HOLIDAY_LIST = append(HOLIDAY_LIST, HOLI_OBJ{Date: "02/15/2010", DESC: PRES})
HOLIDAY_LIST = append(HOLIDAY_LIST, HOLI_OBJ{Date: "04/02/2010", DESC: GOOF})
HOLIDAY_LIST = append(HOLIDAY_LIST, HOLI_OBJ{Date: "05/28/2010", DESC: MEM})
HOLIDAY_LIST = append(HOLIDAY_LIST, HOLI_OBJ{Date: "05/31/2010", DESC: MEM})
HOLIDAY_LIST = append(HOLIDAY_LIST, HOLI_OBJ{Date: "07/05/2010", DESC: INDY})
HOLIDAY_LIST = append(HOLIDAY_LIST, HOLI_OBJ{Date: "09/06/2010", DESC: LABOR})
HOLIDAY_LIST = append(HOLIDAY_LIST, HOLI_OBJ{Date: "10/11/2010", DESC: COLUM})
HOLIDAY_LIST = append(HOLIDAY_LIST, HOLI_OBJ{Date: "11/11/2010", DESC: VETS})
HOLIDAY_LIST = append(HOLIDAY_LIST, HOLI_OBJ{Date: "11/25/2010", DESC: THANK})
HOLIDAY_LIST = append(HOLIDAY_LIST, HOLI_OBJ{Date: "11/26/2010", DESC: THANK})
HOLIDAY_LIST = append(HOLIDAY_LIST, HOLI_OBJ{Date: "12/23/2010", DESC: CHRIS})
HOLIDAY_LIST = append(HOLIDAY_LIST, HOLI_OBJ{Date: "12/24/2010", DESC: CHRIS})

	HOLIDAY_LIST = append(HOLIDAY_LIST, HOLI_OBJ{Date: "01/17/2011", DESC: MARTY})	
	HOLIDAY_LIST = append(HOLIDAY_LIST, HOLI_OBJ{Date: "02/21/2011", DESC: PRES})
	HOLIDAY_LIST = append(HOLIDAY_LIST, HOLI_OBJ{Date: "04/21/2011", DESC: GOOF})
	HOLIDAY_LIST = append(HOLIDAY_LIST, HOLI_OBJ{Date: "04/22/2011", DESC: GOOF})	
	HOLIDAY_LIST = append(HOLIDAY_LIST, HOLI_OBJ{Date: "05/02/2011", DESC: "Early May Bank Holiday"})	
	HOLIDAY_LIST = append(HOLIDAY_LIST, HOLI_OBJ{Date: "05/27/2011", DESC: MEM})
	HOLIDAY_LIST = append(HOLIDAY_LIST, HOLI_OBJ{Date: "05/30/2011", DESC: MEM})
	HOLIDAY_LIST = append(HOLIDAY_LIST, HOLI_OBJ{Date: "07/04/2011", DESC: INDY})
	HOLIDAY_LIST = append(HOLIDAY_LIST, HOLI_OBJ{Date: "08/29/2011", DESC: "Summer Bank Holiday"})	
	HOLIDAY_LIST = append(HOLIDAY_LIST, HOLI_OBJ{Date: "09/05/2011", DESC: LABOR})
	HOLIDAY_LIST = append(HOLIDAY_LIST, HOLI_OBJ{Date: "10/10/2011", DESC: COLUM})	
	HOLIDAY_LIST = append(HOLIDAY_LIST, HOLI_OBJ{Date: "11/11/2011", DESC: VETS})	
	HOLIDAY_LIST = append(HOLIDAY_LIST, HOLI_OBJ{Date: "11/24/2011", DESC: THANK})
	HOLIDAY_LIST = append(HOLIDAY_LIST, HOLI_OBJ{Date: "11/25/2011", DESC: THANK_BLACK_FRIDAY})
	HOLIDAY_LIST = append(HOLIDAY_LIST, HOLI_OBJ{Date: "12/23/2011", DESC: CHRIS})
	HOLIDAY_LIST = append(HOLIDAY_LIST, HOLI_OBJ{Date: "12/26/2011", DESC: CHRIS})	
	HOLIDAY_LIST = append(HOLIDAY_LIST, HOLI_OBJ{Date: "12/30/2011", DESC: NEW})	

HOLIDAY_LIST = append(HOLIDAY_LIST, HOLI_OBJ{Date: "01/02/2012", DESC: NEW})
HOLIDAY_LIST = append(HOLIDAY_LIST, HOLI_OBJ{Date: "01/16/2012", DESC: MARTY})
HOLIDAY_LIST = append(HOLIDAY_LIST, HOLI_OBJ{Date: "02/20/2012", DESC: PRES})
HOLIDAY_LIST = append(HOLIDAY_LIST, HOLI_OBJ{Date: "04/06/2012", DESC: GOOF})
HOLIDAY_LIST = append(HOLIDAY_LIST, HOLI_OBJ{Date: "05/25/2012", DESC: MEM})
HOLIDAY_LIST = append(HOLIDAY_LIST, HOLI_OBJ{Date: "05/28/2012", DESC: MEM})
HOLIDAY_LIST = append(HOLIDAY_LIST, HOLI_OBJ{Date: "07/04/2012", DESC: INDY})
HOLIDAY_LIST = append(HOLIDAY_LIST, HOLI_OBJ{Date: "09/03/2012", DESC: LABOR})
HOLIDAY_LIST = append(HOLIDAY_LIST, HOLI_OBJ{Date: "10/08/2012", DESC: COLUM})	
HOLIDAY_LIST = append(HOLIDAY_LIST, HOLI_OBJ{Date: "10/29/2012", DESC: "Hurricane Sandy"})	
HOLIDAY_LIST = append(HOLIDAY_LIST, HOLI_OBJ{Date: "10/30/2012", DESC: "Hurricane Sandy"})	
HOLIDAY_LIST = append(HOLIDAY_LIST, HOLI_OBJ{Date: "11/12/2012", DESC: VETS})	
HOLIDAY_LIST = append(HOLIDAY_LIST, HOLI_OBJ{Date: "11/22/2012", DESC: THANK})
HOLIDAY_LIST = append(HOLIDAY_LIST, HOLI_OBJ{Date: "11/23/2012", DESC: THANK})
HOLIDAY_LIST = append(HOLIDAY_LIST, HOLI_OBJ{Date: "12/24/2012", DESC: CHRIS})	
HOLIDAY_LIST = append(HOLIDAY_LIST, HOLI_OBJ{Date: "12/25/2012", DESC: CHRIS})	
HOLIDAY_LIST = append(HOLIDAY_LIST, HOLI_OBJ{Date: "12/31/2012", DESC: NEW})	

	HOLIDAY_LIST = append(HOLIDAY_LIST, HOLI_OBJ{Date: "01/01/2013", DESC: NEW})
	HOLIDAY_LIST = append(HOLIDAY_LIST, HOLI_OBJ{Date: "01/21/2013", DESC: MARTY})	
	HOLIDAY_LIST = append(HOLIDAY_LIST, HOLI_OBJ{Date: "02/18/2013", DESC: PRES})
	HOLIDAY_LIST = append(HOLIDAY_LIST, HOLI_OBJ{Date: "03/28/2013", DESC: GOOF})
	HOLIDAY_LIST = append(HOLIDAY_LIST, HOLI_OBJ{Date: "03/29/2013", DESC: GOOF})		
	HOLIDAY_LIST = append(HOLIDAY_LIST, HOLI_OBJ{Date: "05/24/2013", DESC: MEM})
	HOLIDAY_LIST = append(HOLIDAY_LIST, HOLI_OBJ{Date: "05/27/2013", DESC: MEM})
	HOLIDAY_LIST = append(HOLIDAY_LIST, HOLI_OBJ{Date: "07/04/2013", DESC: INDY})	
	HOLIDAY_LIST = append(HOLIDAY_LIST, HOLI_OBJ{Date: "09/02/2013", DESC: LABOR})
	HOLIDAY_LIST = append(HOLIDAY_LIST, HOLI_OBJ{Date: "10/14/2013", DESC: COLUM})
	HOLIDAY_LIST = append(HOLIDAY_LIST, HOLI_OBJ{Date: "11/11/2013", DESC: VETS})
	HOLIDAY_LIST = append(HOLIDAY_LIST, HOLI_OBJ{Date: "11/28/2013", DESC: THANK})	
	HOLIDAY_LIST = append(HOLIDAY_LIST, HOLI_OBJ{Date: "11/29/2013", DESC: THANK})	
	HOLIDAY_LIST = append(HOLIDAY_LIST, HOLI_OBJ{Date: "12/24/2013", DESC: CHRIS})
	HOLIDAY_LIST = append(HOLIDAY_LIST, HOLI_OBJ{Date: "12/25/2013", DESC: CHRIS})
	HOLIDAY_LIST = append(HOLIDAY_LIST, HOLI_OBJ{Date: "12/31/2013", DESC: NEW})
	

HOLIDAY_LIST = append(HOLIDAY_LIST, HOLI_OBJ{Date: "01/01/2014", DESC: NEW})
HOLIDAY_LIST = append(HOLIDAY_LIST, HOLI_OBJ{Date: "01/20/2014", DESC: MARTY})
HOLIDAY_LIST = append(HOLIDAY_LIST, HOLI_OBJ{Date: "02/17/2014", DESC: PRES})
HOLIDAY_LIST = append(HOLIDAY_LIST, HOLI_OBJ{Date: "04/17/2014", DESC: GOOF})	
HOLIDAY_LIST = append(HOLIDAY_LIST, HOLI_OBJ{Date: "04/18/2014", DESC: GOOF})	
HOLIDAY_LIST = append(HOLIDAY_LIST, HOLI_OBJ{Date: "05/23/2014", DESC: MEM})	
HOLIDAY_LIST = append(HOLIDAY_LIST, HOLI_OBJ{Date: "05/26/2014", DESC: MEM})	
HOLIDAY_LIST = append(HOLIDAY_LIST, HOLI_OBJ{Date: "07/03/2014", DESC: INDY})
HOLIDAY_LIST = append(HOLIDAY_LIST, HOLI_OBJ{Date: "07/04/2014", DESC: INDY})
HOLIDAY_LIST = append(HOLIDAY_LIST, HOLI_OBJ{Date: "09/01/2014", DESC: LABOR})	
HOLIDAY_LIST = append(HOLIDAY_LIST, HOLI_OBJ{Date: "10/13/2014", DESC: COLUM})		
HOLIDAY_LIST = append(HOLIDAY_LIST, HOLI_OBJ{Date: "11/11/2014", DESC: VETS})
HOLIDAY_LIST = append(HOLIDAY_LIST, HOLI_OBJ{Date: "11/27/2014", DESC: THANK})
HOLIDAY_LIST = append(HOLIDAY_LIST, HOLI_OBJ{Date: "11/28/2014", DESC: THANK})
HOLIDAY_LIST = append(HOLIDAY_LIST, HOLI_OBJ{Date: "12/24/2014", DESC: CHRIS})		
HOLIDAY_LIST = append(HOLIDAY_LIST, HOLI_OBJ{Date: "12/25/2014", DESC: CHRIS})				
HOLIDAY_LIST = append(HOLIDAY_LIST, HOLI_OBJ{Date: "12/31/2014", DESC: NEW})				
	
	HOLIDAY_LIST = append(HOLIDAY_LIST, HOLI_OBJ{Date: "01/01/2015", DESC: NEW})
	HOLIDAY_LIST = append(HOLIDAY_LIST, HOLI_OBJ{Date: "01/19/2015", DESC: MARTY})	
	HOLIDAY_LIST = append(HOLIDAY_LIST, HOLI_OBJ{Date: "02/16/2015", DESC: PRES})
	HOLIDAY_LIST = append(HOLIDAY_LIST, HOLI_OBJ{Date: "04/02/2015", DESC: GOOF})
	HOLIDAY_LIST = append(HOLIDAY_LIST, HOLI_OBJ{Date: "04/03/2015", DESC: GOOF})
	HOLIDAY_LIST = append(HOLIDAY_LIST, HOLI_OBJ{Date: "05/22/2015", DESC: MEM})
	HOLIDAY_LIST = append(HOLIDAY_LIST, HOLI_OBJ{Date: "05/25/2015", DESC: MEM})	
	HOLIDAY_LIST = append(HOLIDAY_LIST, HOLI_OBJ{Date: "07/02/2015", DESC: INDY})
	HOLIDAY_LIST = append(HOLIDAY_LIST, HOLI_OBJ{Date: "07/03/2015", DESC: INDY})
	HOLIDAY_LIST = append(HOLIDAY_LIST, HOLI_OBJ{Date: "09/07/2015", DESC: LABOR})
	HOLIDAY_LIST = append(HOLIDAY_LIST, HOLI_OBJ{Date: "10/12/2015", DESC: COLUM})
	HOLIDAY_LIST = append(HOLIDAY_LIST, HOLI_OBJ{Date: "11/11/2015", DESC: VETS})
	HOLIDAY_LIST = append(HOLIDAY_LIST, HOLI_OBJ{Date: "11/26/2015", DESC: THANK})
	HOLIDAY_LIST = append(HOLIDAY_LIST, HOLI_OBJ{Date: "11/27/2015", DESC: THANK_BLACK_FRIDAY})
	HOLIDAY_LIST = append(HOLIDAY_LIST, HOLI_OBJ{Date: "12/24/2015", DESC: CHRIS})
	HOLIDAY_LIST = append(HOLIDAY_LIST, HOLI_OBJ{Date: "12/25/2015", DESC: CHRIS})
	HOLIDAY_LIST = append(HOLIDAY_LIST, HOLI_OBJ{Date: "12/31/2015", DESC: NEW})			

HOLIDAY_LIST = append(HOLIDAY_LIST, HOLI_OBJ{Date: "01/01/2016", DESC: NEW})
HOLIDAY_LIST = append(HOLIDAY_LIST, HOLI_OBJ{Date: "01/18/2016", DESC: MARTY})
HOLIDAY_LIST = append(HOLIDAY_LIST, HOLI_OBJ{Date: "02/15/2016", DESC: PRES})
HOLIDAY_LIST = append(HOLIDAY_LIST, HOLI_OBJ{Date: "03/24/2016", DESC: GOOF})
HOLIDAY_LIST = append(HOLIDAY_LIST, HOLI_OBJ{Date: "03/25/2016", DESC: GOOF})
HOLIDAY_LIST = append(HOLIDAY_LIST, HOLI_OBJ{Date: "05/27/2016", DESC: MEM})
HOLIDAY_LIST = append(HOLIDAY_LIST, HOLI_OBJ{Date: "05/30/2016", DESC: MEM})
HOLIDAY_LIST = append(HOLIDAY_LIST, HOLI_OBJ{Date: "07/01/2016", DESC: INDY})	
HOLIDAY_LIST = append(HOLIDAY_LIST, HOLI_OBJ{Date: "07/04/2016", DESC: INDY})
HOLIDAY_LIST = append(HOLIDAY_LIST, HOLI_OBJ{Date: "09/05/2016", DESC: LABOR})
HOLIDAY_LIST = append(HOLIDAY_LIST, HOLI_OBJ{Date: "10/10/2016", DESC: COLUM})
HOLIDAY_LIST = append(HOLIDAY_LIST, HOLI_OBJ{Date: "11/11/2016", DESC: VETS})
HOLIDAY_LIST = append(HOLIDAY_LIST, HOLI_OBJ{Date: "11/24/2016", DESC: THANK})
HOLIDAY_LIST = append(HOLIDAY_LIST, HOLI_OBJ{Date: "11/25/2016", DESC: THANK_BLACK_FRIDAY})
HOLIDAY_LIST = append(HOLIDAY_LIST, HOLI_OBJ{Date: "12/23/2016", DESC: CHRIS})
HOLIDAY_LIST = append(HOLIDAY_LIST, HOLI_OBJ{Date: "12/26/2016", DESC: CHRIS})							
HOLIDAY_LIST = append(HOLIDAY_LIST, HOLI_OBJ{Date: "12/30/2016", DESC: NEW})

	HOLIDAY_LIST = append(HOLIDAY_LIST, HOLI_OBJ{Date: "01/02/2017", DESC: NEW})
	HOLIDAY_LIST = append(HOLIDAY_LIST, HOLI_OBJ{Date: "01/16/2017", DESC: MARTY})	
	HOLIDAY_LIST = append(HOLIDAY_LIST, HOLI_OBJ{Date: "02/20/2017", DESC: PRES})
	HOLIDAY_LIST = append(HOLIDAY_LIST, HOLI_OBJ{Date: "04/13/2017", DESC: GOOF})
	HOLIDAY_LIST = append(HOLIDAY_LIST, HOLI_OBJ{Date: "04/14/2017", DESC: GOOF})		
	HOLIDAY_LIST = append(HOLIDAY_LIST, HOLI_OBJ{Date: "05/26/2017", DESC: MEM})
	HOLIDAY_LIST = append(HOLIDAY_LIST, HOLI_OBJ{Date: "05/29/2017", DESC: MEM})
	HOLIDAY_LIST = append(HOLIDAY_LIST, HOLI_OBJ{Date: "07/03/2017", DESC: INDY})	
	HOLIDAY_LIST = append(HOLIDAY_LIST, HOLI_OBJ{Date: "07/04/2017", DESC: INDY})
	HOLIDAY_LIST = append(HOLIDAY_LIST, HOLI_OBJ{Date: "09/04/2017", DESC: LABOR})
	HOLIDAY_LIST = append(HOLIDAY_LIST, HOLI_OBJ{Date: "10/09/2017", DESC: COLUM})	
	HOLIDAY_LIST = append(HOLIDAY_LIST, HOLI_OBJ{Date: "11/23/2017", DESC: THANK})
	HOLIDAY_LIST = append(HOLIDAY_LIST, HOLI_OBJ{Date: "11/24/2017", DESC: THANK_BLACK_FRIDAY})
	HOLIDAY_LIST = append(HOLIDAY_LIST, HOLI_OBJ{Date: "12/22/2017", DESC: CHRIS})
	HOLIDAY_LIST = append(HOLIDAY_LIST, HOLI_OBJ{Date: "12/25/2017", DESC: CHRIS})
	HOLIDAY_LIST = append(HOLIDAY_LIST, HOLI_OBJ{Date: "12/29/2017", DESC: NEW})
	
HOLIDAY_LIST = append(HOLIDAY_LIST, HOLI_OBJ{Date: "01/01/2018", DESC: NEW})
HOLIDAY_LIST = append(HOLIDAY_LIST, HOLI_OBJ{Date: "01/12/2018", DESC: MARTY})
HOLIDAY_LIST = append(HOLIDAY_LIST, HOLI_OBJ{Date: "01/15/2018", DESC: MARTY})
HOLIDAY_LIST = append(HOLIDAY_LIST, HOLI_OBJ{Date: "02/12/2018", DESC: PRES})
HOLIDAY_LIST = append(HOLIDAY_LIST, HOLI_OBJ{Date: "02/19/2018", DESC: PRES})
HOLIDAY_LIST = append(HOLIDAY_LIST, HOLI_OBJ{Date: "03/30/2018", DESC: GOOF})
HOLIDAY_LIST = append(HOLIDAY_LIST, HOLI_OBJ{Date: "04/02/2018", DESC: EAST})
HOLIDAY_LIST = append(HOLIDAY_LIST, HOLI_OBJ{Date: "05/28/2018", DESC: MEM})
HOLIDAY_LIST = append(HOLIDAY_LIST, HOLI_OBJ{Date: "07/04/2018", DESC: INDY})
HOLIDAY_LIST = append(HOLIDAY_LIST, HOLI_OBJ{Date: "09/03/2018", DESC: LABOR})			
HOLIDAY_LIST = append(HOLIDAY_LIST, HOLI_OBJ{Date: "10/08/2018", DESC: COLUM})
HOLIDAY_LIST = append(HOLIDAY_LIST, HOLI_OBJ{Date: "11/09/2018", DESC: VETS})	
HOLIDAY_LIST = append(HOLIDAY_LIST, HOLI_OBJ{Date: "11/11/2018", DESC: VETS})
HOLIDAY_LIST = append(HOLIDAY_LIST, HOLI_OBJ{Date: "11/12/2018", DESC: VETS})
HOLIDAY_LIST = append(HOLIDAY_LIST, HOLI_OBJ{Date: "11/22/2018", DESC: THANK})	
HOLIDAY_LIST = append(HOLIDAY_LIST, HOLI_OBJ{Date: "11/23/2018", DESC: THANK_BLACK_FRIDAY})			
HOLIDAY_LIST = append(HOLIDAY_LIST, HOLI_OBJ{Date: "12/24/2018", DESC: CHRIS})
HOLIDAY_LIST = append(HOLIDAY_LIST, HOLI_OBJ{Date: "12/25/2018", DESC: CHRIS})
HOLIDAY_LIST = append(HOLIDAY_LIST, HOLI_OBJ{Date: "12/31/2018", DESC: NEW})
						

//	G.Println("   *** HOLIDAY LIST INITIALIZED ***")


} //e nd of func



var LOOK_AHEAD_MINUTES = 180
var LOOK_AHEAD_TEXT = ""


var MAX_LOOK_AHEAD_BARS = 500

type INTERVAL_DEF_OBJ struct {
	INTERVAL 			string
	LOOK_AHEAD_BY  	int			// this is always in hours..	
}

var INTERVAL_MATRIX []INTERVAL_DEF_OBJ

func init_INTERVAL_MATRIX() {
	INTERVAL_MATRIX = append(INTERVAL_MATRIX, INTERVAL_DEF_OBJ{INTERVAL: "1min", LOOK_AHEAD_BY: 3})

	INTERVAL_MATRIX = append(INTERVAL_MATRIX, INTERVAL_DEF_OBJ{INTERVAL: "15min", LOOK_AHEAD_BY: 3})			// For 15 min blocks we use 12 bars (which is 180 mins)

	INTERVAL_MATRIX = append(INTERVAL_MATRIX, INTERVAL_DEF_OBJ{INTERVAL: "HOURLY", LOOK_AHEAD_BY: 3})			// For HOULRY we use 3 bars

	INTERVAL_MATRIX = append(INTERVAL_MATRIX, INTERVAL_DEF_OBJ{INTERVAL: "4hour", LOOK_AHEAD_BY: 8})			// For 4 hour we use 2 bars

	INTERVAL_MATRIX = append(INTERVAL_MATRIX, INTERVAL_DEF_OBJ{INTERVAL: "DAILY", LOOK_AHEAD_BY: 24})			// for DAILY we use one bar

} //end of func




func DATE_EXTRACT(dtype string, dobj time.Time) string {

	var result = ""

	if dtype == "day" {
		dtemp := dobj.Day()
		result = strconv.Itoa(dtemp)

		if dtemp < 10 {
			result = "0" + result
		}

		return result
	}

	if dtype == "month" || dtype == "mon" {
		dtemp := int(dobj.Month())
		result = strconv.Itoa(dtemp)

		if dtemp < 10 {
			result = "0" + result
		}

		return result

	}

	if dtype == "year" {
		dtemp := dobj.Year()
		result = strconv.Itoa(dtemp)

		return result
	}

	return result
}



var STANDARD_5_DEC = "0.00001"
var JPY_3_DEC = "0.001"
var ONE_DEC = "0.1"
var FULL_POINT = "1.0"
var TWO_DEC = "0.01"
var THREE_DEC = "0.001"


// var PERC_0 = 0
// var PERC_1 = 1
// var PERC_2 = 2
// var PERC_3 = 3
// var PERC_5 = 5


type SYMBOL_META_RECORD struct {

	SYMBOL				string
	BROKER				string
	TYPE				string
	TICK_INCREMENT		string		// This needs to be string format because of the issue with precision math
	GMT_HOURS			string		// Keeping the GMT hours here for reference ... we translate THIS hour range to EST\
	MARKET_HOURS 		string		// Market hours are in EST / New York Time

}

/* Rules of 00:00_23:00

	1. Never Trade on Saturday or Sunday
	2. Friday there should be a 3 hour window before trading at close end of day


	*** INVALUABLE TOOL(s):
		https://www.worldtimebuddy.com/
		https://www.fortrade.com/wp-content/uploads/legal/Indices-Trading-Conditions.pdf?v=1542591130794
		https://www.easymarkets.com/int/trade/hours/

	FORMAT:
		Start-of-TradeWindow_End-of_TradeWindow | Secondary-Start-of-TradeWindow_End-of_TradeWindow

	* GMT Hours accurate as of 11/23/2018
	* Market_HOURS (EST/NYT) accurate as of 11/23/2018
	
*/

var CURR_SYM = ""

// func PIP_CALC(increment_string string, highnum float64, lownum float64) float64 {

// 	diff := (highnum - lownum)
// 	inc_num, err := strconv.ParseFloat(increment_string, 64)

// 	if err != nil {
// 		R.Println("")
// 		R.Println(" ERROR in PIP_CALC: ")
// 		Y.Println(err.Error())
// 		os.Exit(1)		
// 	}	

// 	totalpips := (diff / inc_num)

// 	totalpips = FIX_FLOAT_PRECISION(totalpips, 1)

// 	return totalpips

// }// end of PIP_CALC calc

// This takes in a tick string like 0.00001, a source number float and a multiplier 
// and returns the sum of source_NUM + tick_string .. and optionally with tick_string * multiX
func TICK_MATH(tick_string string, source_NUM float64) float64 {

	//1. Convert the tick_string to a float
	inc_num, err := strconv.ParseFloat(tick_string, 64)

	if err != nil {
		R.Println("")
		R.Println(" ERROR in TICK_MATH: ")
		Y.Println(err.Error())
		os.Exit(1)		
	}

	//2. Now do a mathematic addition to this to get our result
	//tmulti := 0.00 + float64(multiX)

	result := (source_NUM + inc_num)

	numdecs := GET_NUM_DECIMAL_PLACES(inc_num)

	//4. Now fix float precision
	finalval := FIX_FLOAT_PRECISION(result, numdecs)

	//3. Finally lets format this to have the proper number of decimals based on the number in TICK_STRING
	// Y.Println("\n Curr SYM: ", CURR_SYM)
	// G.Println("          S: ", source_NUM)
	// M.Println("  FINAL VAL: ", finalval )
	// W.Print(" DECIMALS count: ")
	// C.Print(tick_string)
	// W.Print(" is ")
	// G.Print("", numdecs,"\n\n" )
	



	return finalval
}// end of func




var SYMBOL_MATRIX []SYMBOL_META_RECORD


func init_SYMBOL_METADATA_MATRIX() {


	SYMBOL_MATRIX = append(SYMBOL_MATRIX, SYMBOL_META_RECORD{
				   SYMBOL: "AUDJPY",
				   BROKER: "NADEX",
		 TICK_INCREMENT: JPY_3_DEC,
			MARKET_HOURS: "00:00_23:00",
					  TYPE: "FOREX",
	})
	
	SYMBOL_MATRIX = append(SYMBOL_MATRIX, SYMBOL_META_RECORD{
					   SYMBOL: "AUDUSD",
					   BROKER: "NADEX",
			 TICK_INCREMENT: STANDARD_5_DEC,
				MARKET_HOURS: "00:00_23:00",
				        TYPE: "FOREX",
	})

	SYMBOL_MATRIX = append(SYMBOL_MATRIX, SYMBOL_META_RECORD{
					   SYMBOL: "AUSIDXAUD",
					   BROKER: "OANDA",					   
			 TICK_INCREMENT: FULL_POINT,
			      GMT_HOURS: "00:00_06:30|07:15_21:00",	   
				MARKET_HOURS: "19:00_01:30|02:15_16:00",
				        TYPE: "INDEX",		
	})

	SYMBOL_MATRIX = append(SYMBOL_MATRIX, SYMBOL_META_RECORD{
					   SYMBOL: "BRENTCMDUSD",
					   BROKER: "NADEX",					   
			 TICK_INCREMENT: THREE_DEC,
				   GMT_HOURS: "01:00_22:00",
				MARKET_HOURS: "20:00_17:00",
				        TYPE: "COMMOD",		
	})

	SYMBOL_MATRIX = append(SYMBOL_MATRIX, SYMBOL_META_RECORD{
					   SYMBOL: "BUNDTREUR",
					   BROKER: "OANDA",					   
			 TICK_INCREMENT: TWO_DEC,
				   GMT_HOURS: "XXX_XXXX",
				MARKET_HOURS: "XXX_XXXX",
				        TYPE: "BOND",		
	})	

	SYMBOL_MATRIX = append(SYMBOL_MATRIX, SYMBOL_META_RECORD{
					   SYMBOL: "DEUIDXEUR",
					   BROKER: "NADEX",					   
			 TICK_INCREMENT: TWO_DEC,
			  	   GMT_HOURS: "07:00_21:00",
				MARKET_HOURS: "02:00_16:00",
				        TYPE: "INDEX",		
	})	

	SYMBOL_MATRIX = append(SYMBOL_MATRIX, SYMBOL_META_RECORD{
					   SYMBOL: "EURAUD",
					   BROKER: "OANDA",					   
			 TICK_INCREMENT: STANDARD_5_DEC,
				MARKET_HOURS: "00:00_23:00",
				        TYPE: "FOREX",		
	})	

	SYMBOL_MATRIX = append(SYMBOL_MATRIX, SYMBOL_META_RECORD{
					   SYMBOL: "EURCHF",
					   BROKER: "OANDA",					   
			 TICK_INCREMENT: STANDARD_5_DEC,
				MARKET_HOURS: "00:00_23:00",
				        TYPE: "FOREX",		
	})	

	SYMBOL_MATRIX = append(SYMBOL_MATRIX, SYMBOL_META_RECORD{
					   SYMBOL: "EURGBP",
					   BROKER: "NADEX",					   
			 TICK_INCREMENT: STANDARD_5_DEC,
				MARKET_HOURS: "00:00_23:00",
				        TYPE: "FOREX",		
	})	
	SYMBOL_MATRIX = append(SYMBOL_MATRIX, SYMBOL_META_RECORD{
					   SYMBOL: "EURJPY",
					   BROKER: "NADEX",					   
			 TICK_INCREMENT: JPY_3_DEC,
				MARKET_HOURS: "00:00_23:00",
				        TYPE: "FOREX",		
	})	
	SYMBOL_MATRIX = append(SYMBOL_MATRIX, SYMBOL_META_RECORD{
					   SYMBOL: "EURUSD",
					   BROKER: "NADEX",					   
			 TICK_INCREMENT: STANDARD_5_DEC,
				MARKET_HOURS: "00:00_23:00",
				        TYPE: "FOREX",		
	})	
	SYMBOL_MATRIX = append(SYMBOL_MATRIX, SYMBOL_META_RECORD{
					   SYMBOL: "EUSIDXEUR",
					   BROKER: "OANDA",					   
			 TICK_INCREMENT: TWO_DEC,
			  	   GMT_HOURS: "07:30_21:00",
			   MARKET_HOURS: "02:30_16:00",
				        TYPE: "INDEX",		
	})	
	SYMBOL_MATRIX = append(SYMBOL_MATRIX, SYMBOL_META_RECORD{
					   SYMBOL: "FRAIDXEUR",
					   BROKER: "OANDA",					   
			 TICK_INCREMENT: TWO_DEC,
			  	   GMT_HOURS: "07:00_21:00",
			   MARKET_HOURS: "02:00_16:00",
				        TYPE: "INDEX",		
	})	

	SYMBOL_MATRIX = append(SYMBOL_MATRIX, SYMBOL_META_RECORD{
					   SYMBOL: "GASCMDUSD",
					   BROKER: "NADEX",					   
			 TICK_INCREMENT: THREE_DEC,
				MARKET_HOURS: "XXXX_XXXX",
				        TYPE: "COMMOD",		
	})		

	SYMBOL_MATRIX = append(SYMBOL_MATRIX, SYMBOL_META_RECORD{
					   SYMBOL: "GBPCHF",
					   BROKER: "OANDA",					   
			 TICK_INCREMENT: STANDARD_5_DEC,
				MARKET_HOURS: "00:00_23:00",
				        TYPE: "FOREX",		
	})	
	SYMBOL_MATRIX = append(SYMBOL_MATRIX, SYMBOL_META_RECORD{
					   SYMBOL: "GBPJPY",
					   BROKER: "NADEX",					   
			 TICK_INCREMENT: JPY_3_DEC,
				MARKET_HOURS: "00:00_23:00",
				        TYPE: "FOREX",		
	})	
	SYMBOL_MATRIX = append(SYMBOL_MATRIX, SYMBOL_META_RECORD{
					   SYMBOL: "GBPUSD",
					   BROKER: "NADEX",					   
			 TICK_INCREMENT: STANDARD_5_DEC,
				MARKET_HOURS: "00:00_23:00",
				        TYPE: "FOREX",		
	})	
	SYMBOL_MATRIX = append(SYMBOL_MATRIX, SYMBOL_META_RECORD{
						 SYMBOL: "GBRIDXGBP",
					   BROKER: "NADEX",						 
			  TICK_INCREMENT: TWO_DEC,
			  		 GMT_HOURS: "08:00_21:00",
			 	 MARKET_HOURS: "03:00_16:00",
				        TYPE: "INDEX",		  
	})	
	SYMBOL_MATRIX = append(SYMBOL_MATRIX, SYMBOL_META_RECORD{
					   SYMBOL: "HKGIDXHKD",
					   BROKER: "OANDA",					   
			 TICK_INCREMENT: TWO_DEC,
			  	   GMT_HOURS: "01:15_04:00|05:00_08:15",
				MARKET_HOURS: "20:15_23:00|00:00_03:15",
				        TYPE: "INDEX",		
	})	
	SYMBOL_MATRIX = append(SYMBOL_MATRIX, SYMBOL_META_RECORD{
					   SYMBOL: "JPNIDXJPY",
					   BROKER: "NADEX",					   
			 TICK_INCREMENT: TWO_DEC,
		      MARKET_HOURS: "00:00_23:00",
				        TYPE: "INDEX",		
	})	
	SYMBOL_MATRIX = append(SYMBOL_MATRIX, SYMBOL_META_RECORD{
					   SYMBOL: "LIGHTCMDUSD",
					   BROKER: "NADEX",					   
			 TICK_INCREMENT: TWO_DEC,
				MARKET_HOURS: "00:00_23:00",
				        TYPE: "COMMOD",		
	})	

	SYMBOL_MATRIX = append(SYMBOL_MATRIX, SYMBOL_META_RECORD{
					   SYMBOL: "NZDUSD",
					   BROKER: "OANDA",					   
			 TICK_INCREMENT: STANDARD_5_DEC,
				MARKET_HOURS: "00:00_23:00",
				        TYPE: "FOREX",		
	})				


	SYMBOL_MATRIX = append(SYMBOL_MATRIX, SYMBOL_META_RECORD{
					   SYMBOL: "SUGARCMDUSD",
					   BROKER: "OANDA",					   
			 TICK_INCREMENT: STANDARD_5_DEC,
				MARKET_HOURS: "XXXX_XXXX",
				        TYPE: "COMMOD",		
	})	
	
	
	SYMBOL_MATRIX = append(SYMBOL_MATRIX, SYMBOL_META_RECORD{
					   SYMBOL: "UKGILTTRGBP",
					   BROKER: "OANDA",					   
			 TICK_INCREMENT: THREE_DEC,
		      MARKET_HOURS: "00:00_23:00",
				        TYPE: "BOND",		
	})			

	SYMBOL_MATRIX = append(SYMBOL_MATRIX, SYMBOL_META_RECORD{
					   SYMBOL: "USA30IDXUSD",
					   BROKER: "NADEX",					   
			 TICK_INCREMENT: TWO_DEC,
		      MARKET_HOURS: "00:00_23:00",
				        TYPE: "INDEX",		
	})				

	SYMBOL_MATRIX = append(SYMBOL_MATRIX, SYMBOL_META_RECORD{
					   SYMBOL: "USA500IDXUSD",
					   BROKER: "NADEX",					   
			 TICK_INCREMENT: TWO_DEC,
		      MARKET_HOURS: "00:00_23:00",
				        TYPE: "INDEX",		
	})				


	SYMBOL_MATRIX = append(SYMBOL_MATRIX, SYMBOL_META_RECORD{
					   SYMBOL: "USATECHIDXUSD",
					   BROKER: "NADEX",					   
			 TICK_INCREMENT: TWO_DEC,
			   MARKET_HOURS: "00:00_23:00",
				        TYPE: "INDEX",		
	})				
	SYMBOL_MATRIX = append(SYMBOL_MATRIX, SYMBOL_META_RECORD{
					   SYMBOL: "USDCAD",
					   BROKER: "NADEX",					   
			 TICK_INCREMENT: STANDARD_5_DEC,
				MARKET_HOURS: "00:00_23:00",
				        TYPE: "FOREX",		
	})				
	SYMBOL_MATRIX = append(SYMBOL_MATRIX, SYMBOL_META_RECORD{
					   SYMBOL: "USDCHF",
					   BROKER: "NADEX",					   
			 TICK_INCREMENT: STANDARD_5_DEC,
				MARKET_HOURS: "00:00_23:00",
				        TYPE: "FOREX",		
	})				
	SYMBOL_MATRIX = append(SYMBOL_MATRIX, SYMBOL_META_RECORD{
					   SYMBOL: "USDJPY",
					   BROKER: "NADEX",					   
			 TICK_INCREMENT: JPY_3_DEC,
				MARKET_HOURS: "00:00_23:00",
				        TYPE: "FOREX",		
	})				

	SYMBOL_MATRIX = append(SYMBOL_MATRIX, SYMBOL_META_RECORD{
					   SYMBOL: "XAGUSD",
					   BROKER: "NADEX",					   
			 TICK_INCREMENT: THREE_DEC,
				MARKET_HOURS: "00:00_23:00",
				        TYPE: "COMMOD",		
	})				
	SYMBOL_MATRIX = append(SYMBOL_MATRIX, SYMBOL_META_RECORD{
					   SYMBOL: "XAUUSD",
					   BROKER: "NADEX",					   
			 TICK_INCREMENT: THREE_DEC,
			 	MARKET_HOURS: "00:00_23:00",
				        TYPE: "COMMOD",		 
	})																		
	
} //end of func


// type TZ_REC struct {
// 	START_TIME    time.Time
// 	END_TIME      time.Time
// }

/* This translates (primarly GMT --> EST times) in the SYMBOL_MATRIX and checks that the
   input_DATE is within the range of those dates. 
	It takes care of 
*/
func MARKET_HOUR_VALIDATION(symbol string, input_DATE time.Time) (bool, time.Time) {

	//1. Find the symbol we need to "translate" for

	var TSYM SYMBOL_META_RECORD

	var found = false
	for _, tobj := range SYMBOL_MATRIX {

		// W.Println("OBJ SYM: **" + tobj.SYMBOL)
		// G.Println("  insym: **" + symbol )
		if tobj.SYMBOL == symbol {
			TSYM = tobj
			found = true
			break
		}

	} //end of for

	if found == false {

		R.Println(" *** ERROR: Cant find a record in SYMBOL META!!! **" + symbol)
		os.Exit(1)
		//return false, input_DATE
	}


	//2. Now we will create a temporary list of start times and end times appropriate to this 
	//symbol

	// var TIMES []TZ_REC

	//3. Lets see if we have multiple Time ranges in the MARKET_HOURS part of this
	// Multiple times are seperated by the | pipe
	var flist []string

	
	if strings.Contains(TSYM.MARKET_HOURS, "|") {
		d := strings.Split(TSYM.MARKET_HOURS, "|")

		for _, dobj := range d {
			flist = append(flist, dobj)
		}

	//3b. Otherwise, there is just one time range in the matrix, lets just save it
	} else {

		flist = append(flist, TSYM.MARKET_HOURS)
	}


	// badTime := time.Now()

	// return true, badTime

	//4. Now we have a list of time ranges. (or just one) .. Lets iterate through them and 
	// create date objects

	/*	4b. Start Time is left of the _ .. end time is RIGHT of it
		Now while the M/d/y doesnt really matter.. Lets create a time.Time object
	 	We are using 10/15/2010 as a static time
	*/

	// Lets Get the YEAR/MON and DAY from the Date Object that was passed in
	cmon := time.Month(input_DATE.Month())
	cday := input_DATE.Day()
	cyear := input_DATE.Year()

	for _, r := range flist {

		//3b ERROR HANDLING:
		if strings.Contains(r, "_") == false {
			W.Println(" Found a BROKE one!!", r)
			os.Exit(1)
		}
		sobj := strings.Split(r, "_")



		START_STRING := sobj[0]

		tsplit := strings.Split(START_STRING, ":")
		start_HOUR := tsplit[0]
		start_MIN := tsplit[1]

		END_STRING := sobj[1]
		qsplit := strings.Split(END_STRING, ":")
		end_HOUR := qsplit[0]
		end_MIN := qsplit[1]
		
		s_hour, _ := strconv.Atoi(start_HOUR)
		s_min, _ :=  strconv.Atoi(start_MIN)

		e_hour, _ := strconv.Atoi(end_HOUR)
		e_min, _ := strconv.Atoi(end_MIN)
		
		start_dateOBJ := time.Date(cyear, cmon, cday, s_hour, s_min, 0, 0, GMT_LOCATION_OBJ)
	
		//5. Because a date might "flow into the next day" i am incrementing the start_dateOBJ until
		// i reach the end date objct HH:MM
		end_dateOBJ := start_dateOBJ

		var found = false

		for x := 0; x < TOTAL_MINUTES_DAY; x++ {

			//5b. Lets check the HOUR and MINUTE of this new date object. IF they equal the
			//e_hour and e_min .. we know to exit and use this dateOBJ
			T_HOUR := end_dateOBJ.Hour()
			T_MIN := end_dateOBJ.Minute()

			// C.Println("T_HOUR/Mins is: ", T_HOUR, T_MIN)
			// Y.Println("ehour and min is: ", e_hour, e_min)
			// W.Println("")
			if T_HOUR == e_hour && T_MIN == e_min {
				found = true
				break
			}

			end_dateOBJ, _ = DateMath(end_dateOBJ, "add", 1, "min")

			
		} //end of FOR




		//5c. Error handling.. If found is false, means we went through the ENTIRE loop. 
		// This should NEVER happen and I want to hard stop to figure out why:
		if found == false {
			R.Println("ERROR in IS_VALID_TO_TRADE calculation.")
			R.Println(" We didnt find the final end_DATEOBJ for some reason")
			C.Println(" SYMBOL: ", symbol)
			Y.Print(" EVALUATIING: ")
			R.Println(start_dateOBJ)
			Y.Print(" thru ")
			W.Println(end_dateOBJ)
			os.Exit(0)
		}

		//5f. Short Circuit.. If found is true, we found a valid Object in the symbol matrix. Lets check this
		// against the input_DATE		
		if input_DATE.After(start_dateOBJ) && input_DATE.Before(end_dateOBJ) {
			return true, end_dateOBJ
		}

	}//end

	//7. Default is to return false... means the input_DATE is not in the valid trade window for the security

	return false, input_DATE

} //end of func


var SHOW_DEBUG_DATA = false

//This makes sure we are trading within a valid window    DONT CALL DIRECTLY.. use CAN_TRADE_NOW()
func is_WEEKEND(DATEOBJ time.Time) bool {

	//1. Get the current DAY
	curr_DAY := strings.ToLower(DATEOBJ.Weekday().String())

	//R.Println ("CURRDAY is: ", curr_DAY)

	//2. We NEVER trade on saturday or sunday.. no matter what currency
	if strings.Contains(curr_DAY, "saturday") || strings.Contains(curr_DAY, "sunday") {

		if SHOW_DEBUG_DATA {
			R.Println("*** ")
			R.Print("*** ")
			W.Print(" WINDOW CHECK FAIL: This is a WEEKEND ")
			M.Println(curr_DAY)
			R.Println("*** ")
		}

		return true
	}

	//3. Otherwise, we return true

	return false

} //end of function


// This takes a specific date XX/XX/XXXX and returns true or false if its a holiday
func IS_HOLIDAY(input_date time.Time) bool {

	i_mon := DATE_EXTRACT("mon", input_date)
	i_day := DATE_EXTRACT("day", input_date)
	i_year := DATE_EXTRACT("year", input_date)

	// R.Println(" i_mon is: ", i_mon, i_day, i_year)

	if i_mon == "" || i_day == "" || i_year == "" {
		return false
	}


	date_TEMP := i_mon + "/" + i_day + "/" + i_year

	for _, h := range HOLIDAY_LIST {
		// Y.Println("Date is ", h.Date)
		if date_TEMP == h.Date {

			if SHOW_DEBUG_DATA {
				Y.Print(" ** HOLIDAY DETECTED: ")
				M.Println(h.Date, "-", h.DESC)
			}

			return true

		}
	}

	return false
} //end of isHoliday

// This determines if we can trade in this window.. It returns true/false as well as the calculated "LOOK _AHEAD DATE" based
// on Look_AHEAD_MINUTES
func ENOUGH_TIME_REMAINING_IN_TRADE_WINDOW(startWINDOW time.Time, market_CLOSE_DATE time.Time) (bool, time.Time) {
	
	//1. Take the startWindow and add LOOK_AHEAD_MINUTES to it.. 
	look_AHEAD_DATE, _ := DateMath(startWINDOW, "add", LOOK_AHEAD_MINUTES, "mins")

	//2. Next, we have a FRIDAY exception.. End of day friday is 15:30 EST .. All markets (including forex) close
	// at this time... So if it is friday and look_ahead is AFTER 15:30 .. this is returned as false

/*
	spretty, swee := SHOW_PRETTY_DATE(startWINDOW)
	


	Y.Println("   ENOUGH START: ", spretty, " | ", swee)
	G.Println("Look AHEAD DATE: ", lpretty, " | ", lpretty, "\n")
	*/

	_, lweekday := SHOW_PRETTY_DATE(look_AHEAD_DATE, false)
	lweekday = strings.ToLower(lweekday)

	if strings.Contains(lweekday, "friday") {


		cYear := look_AHEAD_DATE.Year()
		cMonthObj := time.Month(look_AHEAD_DATE.Month())
		cDay := look_AHEAD_DATE.Day()

		friday_OVERRIDE_DATE := time.Date(cYear, cMonthObj, cDay, 15, 30, 0, 0, GMT_LOCATION_OBJ) // We always do UTC/GMT

/*
		fpretty, fwee := SHOW_PRETTY_DATE(friday_OVERRIDE_DATE)

		M.Println(" Fpretty: ", fpretty, fwee)
*/

		// Now if our original lookAhead DATE is AFTER the friday_OVERRIDE Date, we return false
		if look_AHEAD_DATE.After(friday_OVERRIDE_DATE) {

			/*
			G.Println(" LOOK AHEAD DATE IS AFTER OVERRIDE DATE!!!!")
			PressAny()
			*/
			return false, startWINDOW

		}

	//3. Otherwise, (on any other day but friday) this is the behavior we process
	} else  {

		// startd, _ := SHOW_PRETTY_DATE(startWINDOW)
		// lookd, _ := SHOW_PRETTY_DATE(look_AHEAD_DATE)

		// C.Println("     STARTD IS: ", startd)
		// W.Println(" LOOK AHEAD IS: ", lookd, "\n")
		//2. Now lets make sure that the look_AHEAD_Date is LESS than the market_CLOSE_DATE

		//3. Otehrwise, we look ahead.. if the look ahead date is after market_Close date, we return false
		if look_AHEAD_DATE.After(market_CLOSE_DATE) {
			return false, startWINDOW
		}


	}



	//5. Default to true unless we had eception above

	return true, look_AHEAD_DATE

} //end of func

/* This checks to see if we can trade now:   (PRAE ENGINE and any simulators use this)
	- Returns FALSE if this is SATURDAY or SUNDAY
	- Returns FALSE if this is a HOLIDAY
	- REturns FALSE if the particular security is not available to trade at this TIME
	

	// Accepts a time.Time date OBJECT
*/
func CAN_TRADE_NOW(SYMBOL string, dateOBJ time.Time) (bool, time.Time, string) {

	// If this is a weekend (SATURDAY OR SUNDAY)
	if is_WEEKEND(dateOBJ) {

		return false, dateOBJ, ""
	
	}

	//2. Check if this is a holiday
	if IS_HOLIDAY(dateOBJ) {
		
		return false, dateOBJ, ""
	
	}
			
	//3.  We need to check the SYMBOL against the time of day... 
	// If it is a valida time to trade this security, lets make sure we have at least 3 hours 
	// before "close" .. .. This will ensure we have a window to analyze price movement
	tready, market_CLOSE_TIME := MARKET_HOUR_VALIDATION(SYMBOL, dateOBJ)
	
	/*if tready {
		G.Println(" TREADY is ", tready)
		PressAny()
	} else {
		M.Println("TREADY is ", tready )
	}
	*/
	
	// If Market_HOUR is false, we return false
	if tready == false {
		return false, dateOBJ, ""

	} 
	
	// Finally if there is NOT enough time remaining in trade Window, we return false

	eresult, LOOK_AHEAD_DATEOBJ := ENOUGH_TIME_REMAINING_IN_TRADE_WINDOW(dateOBJ, market_CLOSE_TIME)
	if eresult  == false {
		return false, dateOBJ, ""
	}

	look_pretty, _ := SHOW_PRETTY_DATE(LOOK_AHEAD_DATEOBJ, false)	
	// Default is true,, and we return the calculated LOOK_AHEAD_DATE
	return true, LOOK_AHEAD_DATEOBJ, look_pretty

} //end of func












func NumDecPlaces(v float64) int {
	s := strconv.FormatFloat(v, 'f', -1, 64)
	i := strings.IndexByte(s, '.')
	if i > -1 {
		return len(s) - i - 1
	}
	return 0
}



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



var start_TEXT = ""
var end_TEXT = ""

/*
	Takes in a 01/01/2001 type date and returns the appropriate time.Time date objects
*/
func GET_DATE_RANGE_Objects(STARTDATE string, ENDDATE string) (time.Time, time.Time) {

	//1. Split the input date
	sd := strings.Split(STARTDATE, "/")

	s_month, _ := strconv.Atoi(sd[0])
	s_day, _ := strconv.Atoi(sd[1])
	s_year, _ := strconv.Atoi(sd[2])

	//2. Safety.. if there is a space, we split on it
	if strings.Contains(sd[2], " ") {

		es := strings.Split(sd[2], " ")
		s_year, _ = strconv.Atoi(es[0])
	}
	e_month := s_month
	e_day := s_day
	e_year := s_year

	//2. If they SPECIFIED a particular end date.. lets use THAT instead of what we default to (which is the start date)
	if ENDDATE != "" {
		ed := strings.Split(ENDDATE, "/")

		e_month, _ = strconv.Atoi(ed[0])
		e_day, _ = strconv.Atoi(ed[1])
		e_year, _ = strconv.Atoi(ed[2])

		//3. Safety.. if there is a space, we split on it
		if strings.Contains(ed[2], " ") {

			es := strings.Split(ed[2], " ")
			e_year, _ = strconv.Atoi(es[0])
		}
	}

	s_MonthObj := time.Month(s_month)
	e_MonthObj := time.Month(e_month)

	startObj := time.Date(s_year, s_MonthObj, s_day, 0, 0, 0, 0, time.UTC)
	endObj := time.Date(e_year, e_MonthObj, e_day, 23, 59, 0, 0, time.UTC)

	return startObj, endObj
}






/*
 Pulls ALL the symbols from the PRICE_DATA table
*/
func PULL_ALL_SYMBOLS() []string {

	var finalResults []string

	var mongo = DBSession.DB("PRICE_DATA").C("PRICES_1min")

	mongo.Find(nil).Distinct("SYMBOL", &finalResults)		

	//2. Now sort this list of strings alphabetically 
	sort.Strings(finalResults)

	C.Println("\n ***")
	Y.Print("     Pulling ALL ")
	Y.Print("symbols... ")
	G.Print(len(finalResults))
	Y.Println(" were FOUND!!")
	C.Println(" ***\n")

	return finalResults

}






// Initializes when this module is attached
func init() {

	//1. INitialize the Financial Market Holiday List:
	// We dont trade on or PRAE engine on these days

	INIT_HOLIDAY_LIST()

	//2. Init the INTERVAL matrix
	init_INTERVAL_MATRIX()

	//3. Init the SYMBOL INCREMENT matrix .. This includes the tick increment used for thesymbol and its trading times
	init_SYMBOL_METADATA_MATRIX()

}
