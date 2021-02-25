/*	APIceCream - A bunch of SUPER lean REST and JSON helper functions to make REST dev in GO ..FUN!

	v1.54	Feb 22, 2021	- Removed SCrape Master.. moved a better version to TERRY_COMMON
	v1.52	Feb 03, 2020	- Added JSON_DOWNLOAD. Downloads a json Byte object you can read or iterate
	v1.50	Jan 06, 2020	- Rebranded just making this the generic name for this Module
	v1.45	Dec 27, 2019	- Added parameter parse support. sends them in an array you can access
	v1.40	Feb 02, 2019	- Now  This supports SSL and non-SSL (just http) mode. 
				              				If you need SSL, specifiy a PEM and KEY file on the command line
							  				with SSL PEM and Key file specified
							  
	v1.30	Nov 05, 2018	- Added support for Spinnig up a micro webserver
	v1.24	Sep 24, 2018	- Initial Release

*/

package APIceCream

import (

	// 1. Native LIBS
		"net/http"
		"unicode"
		"net/url"
		"os"
		"flag"
		"strings"
		"io/ioutil"

	// 2. 3RD Party Libs
		"github.com/buger/jsonparser"

	// 3. Proprietary (Personal) Libs/Modules
		. "dev.azure.com/acetraderllc/shared/_git/PUBLIC_Ace.git/GOGO_Gadgets"
)



/*
	###
	###		Globals!!
	###
*/

var DEFAULT_LISTEN_PORT = "8888" // This is the default "listen" port for all webservices.. if using ports below 1024 you must run as root

var JSON_CLEAN_FLAG	   = false	// When set to TRUE, we return "MINIFIED" json  (it omits whatever spaces and indents you used get stripped out )

var SSL_ENABLE_FLAG     = false	// If set to true.. we listen in SSL mode.. meaning we need a PEM file and a KEY specified below
var SSL_CERT_PEM_FILE   = ""		// Full path to wherever this CERT/ pem file is			(this is just example: /opt/SSL_CERTS/biolab_COMBINED.pem  )
var SSL_KEY_FILE        = ""		// Full path to wherever the KEY file of this cert is (this is just example: /opt/SSL_CERTS/biolab_ALPHA-SSL.key )


func MinifyJSON(str string) string {
	return strings.Map(func(r rune) rune {
		 if unicode.IsSpace(r) {
			  return -1
		 }
		 
		 return r
	}, str)
} //end of 


/* This is meant to be passed a keymap of url.Values
Returns true and expects a POINTER to hold the value that it finds
If keyname is found in the map, its VALUE is returned
NOTE: this lets you send URL vars of ANY case

FIND_URL_key_using_POINTER("data", URL_MAP, &result):
*/
func FIND_URL_key(keyname string, UMAP url.Values, myresult *string) bool {

	//1. this allows you to specify case INSENSITIVE keyNames
	kUpper := strings.ToUpper(keyname)
	klow := strings.ToLower(keyname)
	firstLetter := UpperFirst(keyname)

	if keyValue_ARRAY, ok := UMAP[kUpper]; ok {

		*myresult = keyValue_ARRAY[0]

		return true
	}

	if keyValue_ARRAY, ok := UMAP[klow]; ok {

		*myresult = keyValue_ARRAY[0]

		return true
	}

	if keyValue_ARRAY, ok := UMAP[firstLetter]; ok {

		*myresult = keyValue_ARRAY[0]

		return true
	}

	return false
}

/*

	FIND_URL_key
	Takes in a key to search for
	this is similar to find_url_key but doesnt use a pointer
*/
func Simple_FIND_URL_key(keyname string, UMAP url.Values) (bool, string) {

	//1. this allows you to specify case INSENSITIVE keyNames
	kUpper := strings.ToUpper(keyname)
	klow := strings.ToLower(keyname)
	firstLetter := UpperFirst(keyname)

	if keyValue_ARRAY, ok := UMAP[kUpper]; ok {

		return true, keyValue_ARRAY[0]
	}

	if keyValue_ARRAY, ok := UMAP[klow]; ok {

		return true, keyValue_ARRAY[0]

	}

	if keyValue_ARRAY, ok := UMAP[firstLetter]; ok {

		return true, keyValue_ARRAY[0]

	}

	return false, ""
} //end of func



// This takes in some json (text blurb) and parameters and returns the value you ask for
func Extract_JSON(inputBODY string, jsargs ...string) (string, bool) {

	var result_FLAG = false
	var output_Value = ""

	inputBYTES := []byte(inputBODY)

	singleVAL_BYTES, _, _, err2 := jsonparser.Get(inputBYTES, jsargs...)
	output_Value = string(singleVAL_BYTES)

	if err2 != nil {
		R.Println("\n Extract_JSON PARSE Error, invalid Input!!")
		result_FLAG = false
	} else {
		result_FLAG = true
	}

	return output_Value, result_FLAG

} //end


// This is the main way we read from JSON API's
// Returns the VALUE you specified with jsargs, teh WHOLE text json output you are walking and TRUE if successful (false if error)
func PULL_JSON(dest_URL string, jsargs ...string) (bool, string, string, []byte) {

	//1. Open an http connection to the JSON API URL
	var emptyBYTE []byte
	jresp, err := http.Get(dest_URL)
	if err != nil {
		R.Println("\n** ERROR: Cant Read from the JSON API Url: ", dest_URL)
		

		return false, "", "", emptyBYTE
	}

	//2. Now we have the json response in byte format.. lets put the whole thing in a string
	JSON_BYTE_OBJ, _ := ioutil.ReadAll(jresp.Body)
	FULL_JSON_RESPONSE := string(JSON_BYTE_OBJ)
	SINGLE_JSON_VALUE := ""

	singleVAL_BYTES, _, _, err2 := jsonparser.Get(JSON_BYTE_OBJ, jsargs...)
	SINGLE_JSON_VALUE = string(singleVAL_BYTES)

	if err2 != nil {
		R.Println("\n ERROR IN PULL_JSON: ", err2)
		Y.Print(" DEST URL: ")
		W.Println(dest_URL)
		return false, "", "", emptyBYTE
	}

	//3. Else return our results

	return true, SINGLE_JSON_VALUE, FULL_JSON_RESPONSE, JSON_BYTE_OBJ

} //end of func


// Alias for PULL_JSON
func GET_JSON(dest_URL string, jsargs ...string) (bool, string, string, []byte) {

	bool_val, str_01, str_02, byte_val := PULL_JSON(dest_URL, jsargs...)
	return bool_val, str_01, str_02, byte_val
} //end of func

// Alias for PULL_JSON
func READ_JSON(dest_URL string, jsargs ...string) (bool, string, string, []byte) {

	bool_val, str_01, str_02, byte_val := PULL_JSON(dest_URL, jsargs...)
	return bool_val, str_01, str_02, byte_val
} //end of func



/*
 This downloads json from an url and returns a JSON_BYTE object
 You can read from this object with jsonparser
 You can also iterate elements in this byte object with jsonparser
*/
var MAX_JSON_RETRY = 3		// Maximum times we retry on any error
func JSON_DOWNLOAD(API_URL string) (bool, []byte, string) {


	//1. Setup the JSON_BYTE_OBJ
	var JSON_BYTE_OBJ []byte
	var WAS_SUCCESS = false

	for x := 0; x < MAX_JSON_RETRY; x++ {
		var HAD_PROBLEM = false

		var err2 error
		//2. Now do the json call
		jresp, err := http.Get(API_URL)

		if err != nil {
			M.Println(" ** Error in JSON_Download: ", err)
			HAD_PROBLEM = true
		} else {

			JSON_BYTE_OBJ, err2 = ioutil.ReadAll(jresp.Body)

			if err2 != nil {
				M.Println(" ** ErrorTWO in JSON_Download: ", err)
				HAD_PROBLEM = true
			}		
		}

		//3. Check for an error, if so .. we retry
		if HAD_PROBLEM {
			M.Println("")
			M.Println(" ** Retrying", x, "of", MAX_JSON_RETRY)
			M.Println("")
			Sleep(5, true)

		//4. If NO error, exit the loop and set fail_flag to FALSE
		} else {
			WAS_SUCCESS = true
			break
		}
	} //end of for

	//3. Some arbitraray error handling:
	response_TEXT := string(JSON_BYTE_OBJ)

	if response_TEXT == "" || len(response_TEXT) <= 3 {
		WAS_SUCCESS = false

		R.Print(" ** Error in JSON_Download: ")
		Y.Print(" Not enough Text Returned: ")
		W.Println(response_TEXT)
	}



	//4. Finally return our byte OBJ , a failure flag and string version of JSON response
	return WAS_SUCCESS, JSON_BYTE_OBJ, response_TEXT

} //end of func




/*

 Sends a response to HTTP requests. Defaults to a generic JSON response

*/
func SEND_Response_2_Client(JSON_Response string, ClientResponseWriter *http.ResponseWriter) {

	//M.Print("\n\n   *** Sending Response Back to HTTP Client ...\n")

	//1. Set the outgoing headers for this message type:
	(*ClientResponseWriter).Header().Set("Content-Type", "application/json") 	// ---- Set THIS header for sending back an actual json object
	(*ClientResponseWriter).Header().Set("Access-Control-Allow-Origin", "*") 	//  --- always need this regardless

	//2. This is a default message.. i keep this for a way to debug (meaning the REST stuff works but the JSON is just fucked up)
	defaultMessage := `
    [
		{
			"ERROR (default Resp #1)" : "No parameters specified",
		  	"ERROR (default Resp #2)" : "or.. Invalid API endpoint requested"
		}
    ]
`
	//2b. If no json provided, we give a default response
	if JSON_Response == "" {

		JSON_Response = defaultMessage

	}

	//3. Now lets check the JSON_CLEAN_FLAG .. if it is true, we strup out all the spaces and tabs
	if JSON_CLEAN_FLAG {
		JSON_Response = MinifyJSON(JSON_Response)
	}

	/*
			5. OLD  Method..uses NewEncoder
			err := json.NewEncoder(*ClientResponseWriter).Encode(JSON_Response)

			5c. This method just pushes the already formatted outgoing JSON directly on the ClientResponseWriter

			 rawIn := json.RawMessage(JSON_Response)
			 bytes, err2 := rawIn.MarshalJSON()
			if err2 != nil {
			 }	
			byteObj := []byte(bytes)
	*/


	//4. This formats the string and dumps it on the HTTP response writer
	byteObj := []byte(JSON_Response)
	_, err := (*ClientResponseWriter).Write(byteObj)
	if err != nil {
		R.Print("*** ERROR ***: Weird Error responding to the Client!!!")
		Y.Println(err.Error())
	}

} //end of SEND_Response_2_Client



/*
  This is the 'http service / routing handler' definition .. that is called... from here we call other "routeCOmmand(s)"
  which form the business logic. These are Controllers for the routes that are used/passed
  anyhting after the route/entry point is a "variable" that is consumed by the routeCommand
  You must ALWAYS include ?data= to trigger this.. the service will error out otherwise
*/
type GENERIC_API_ENDPOINT_HANDLER func([]URL_PARAMS) string
/*
	1. EXAMPLE of an ENDPOINT HANLDER.. Should look like this:

			func HANDLER_for_API_endpoint (URL_input []string) {

				//1. Do some business logic based on the URL_input that was passed.. ie: localhost:8020/getuser/username123

				OUTGOING_JSON = `
				[
					{
						"HelloWorld":"This is valid JSON sending back to client!!"
					}
				]
				`

				return OUTGOING_JSON
			}



	//2. Now when you startup your program, great a service endpoint TO this handler like this:

	CREATE_SERVICE_ENDPOINT("/getuser", HANDLER_for_API_endpoint) 

	//3. Now start your listener service:
	
	Start_LISTENER_SERVICE_Engine()

*/

type URL_PARAMS struct {
	KEY		string
	Value	string
}

/*
 This parses the input prameters passed.. Supports the following formats

	?name=Mary&job=librarian

	or

	/name/Valerie/age/37
 
*/
func PARAM_PARSER(input string, PARAMLINE string) []URL_PARAMS {

	var RESULTS []URL_PARAMS


	//1. Lets determine what format the params are in
	// if a ? is specified, we have parameters specified on the get line
	if strings.Contains(input, "=") {

		//2. If there are multiple params, we need to split on &
		if strings.Contains(input, "&") {
			ms := strings.Split(input, "&")

			for _, x := range ms {			

				sp := strings.Split(x, "=")
				var v URL_PARAMS
				v.KEY = sp[0]
				v.Value = sp[1]

				if v.KEY != "" && v.Value != "" {
					RESULTS = append(RESULTS, v)
				}
			} //end of for

		//2b. Otherwise there is a SINGLE parameter here.. Short and sweet
		} else {
			sp := strings.Split(input, "=")
			var v URL_PARAMS
			v.KEY = sp[0]
			v.Value = sp[1]

			if v.KEY != "" && v.Value != "" {
				RESULTS = append(RESULTS, v)
			}
		}
	}

	/*
		** DISABLING passing URL vars this way.. it is not reliable (and difficult to parse the router from the variables that follow it (hte mUX router reads all of them)
		BECAUSE THIS FORMAT IS INCONSISTANT 
	 3. Otherwise they have specified the parameter as a / values.. which means the FIRST value is the KEY (second param is the value).. 	
	  You can have multiples of these

		} else {

			//4. Split on /
			pp := strings.Split(PARAMLINE, "/")

			//4b, Iterate through them
			for x := 0; x < len(pp); x++ {

				item_key := pp[x]

				R.Println("KEY IS: ", item_key)

				if item_key == "" {
					x++
					continue
				}

				// a little error handling
				if strings.Contains(item_key, "favicon") || item_key == "" {
					continue
				}

				var v URL_PARAMS
				v.KEY = item_key

				//4c. The value of this item will be the very next element
				valnum := x+1

				//4d. But we need to make sure this exists in the array (or we will get an index out of range error)
				if valnum < len(pp) {
					v.Value = pp[valnum]
					x++

				//4e. IF no value is avalable, dont save this parameter item
				} else {
					continue
				}

				//5. Otherwise, save this parameter into our array			
				RESULTS = append(RESULTS, v)
				
			}
		}
	*/

	return RESULTS
} //end of func



func (RouteEntry_Handler_FUNC_to_USE GENERIC_API_ENDPOINT_HANDLER) ServeHTTP(outgoingResponseObj http.ResponseWriter, inReq *http.Request) {

	//1. Remember,the outgoingResponseObj is passed in by the MUX... we have to send a pointer to this object forward
	//  MUST ALWAYS DO THIS so the response gets sent
	
	//2. This splits the URL path and vars on the forward slash	
	REQ_PARAMS := PARAM_PARSER(inReq.URL.RawQuery, inReq.URL.Path)
	for _, p := range REQ_PARAMS {
		W.Print("\n **DEBUG PARAMS** KEY: ")
		Y.Print(p.KEY)
		W.Print(" | Value: ")
		G.Println(p.Value)
	} //end of for
	
	//3. Send it to the appropriate "controller/handler"   .. Every handler should respond with formatted JSON
	JSON_response := RouteEntry_Handler_FUNC_to_USE(REQ_PARAMS)

	//4. Finally we send a response TO the client.. meaning we are done.HandleIncomingData
	SEND_Response_2_Client(JSON_response, &outgoingResponseObj)

} //end of http service / routing handler

// Gets a value from URL params. Takes in a list of URL Params
func GET_VALUE(KEY string, inputVARS []URL_PARAMS) string {
	//1. FIrst we iterate through the Url Params looking for the one that matches they KEY speicfied
	for _, x := range inputVARS {
		if x.KEY == KEY {
			return x.Value
		}
	} //end of for


	//2. otherwise if we find nothing, return nothing
	return ""

} // end of func

// Alias to GET_VALUE
func FIND_VALUE(KEY string, inputVARS []URL_PARAMS) string {
	return GET_VALUE(KEY, inputVARS)
}
func GET_KEY(KEY string, inputVARS []URL_PARAMS) string {
	return GET_VALUE(KEY, inputVARS)
}

// This is mostly for debug, just shows all values in an URL_PARAMS list
// returns  a json formatted string
func SHOW_ALL_PARAMS(inputVARS []URL_PARAMS) string {

	var JSON_OUTPUT = ``
	for _, x := range inputVARS {
		JSON_OUTPUT += "     " + x.KEY + `:` + x.Value + `,
`		
	} //end of 

	return JSON_OUTPUT	
} //edn of SHOW ALL


type API_ENDPOINT_DEFINITION struct {
	API_Endpoint string
	API_Handler  GENERIC_API_ENDPOINT_HANDLER
}

var ALL_SERVICE_ENDPOINTS []API_ENDPOINT_DEFINITION

/*
	Call this to create a new REST API Endpoint /THISADDY 
*/
func CREATE_SERVICE_ENDPOINT(api_endpoint string, handler_name GENERIC_API_ENDPOINT_HANDLER) {

	var etmp API_ENDPOINT_DEFINITION
	etmp.API_Endpoint = api_endpoint
	etmp.API_Handler = handler_name

	ALL_SERVICE_ENDPOINTS = append(ALL_SERVICE_ENDPOINTS, etmp)
} //end func


func Start_LISTENER_SERVICE_Engine() { 

	//1. Error handling
	slen := len(ALL_SERVICE_ENDPOINTS)

	// C.Println(" Hey Slen is: ", slen)
	if slen == 0 {
		R.Println(" ** No Entries in the ALL_SERVICE_ENDPOINTS list")
		return
	}

	//1. First lets create the SERVICE object

	SERVICE_MUX_OBJ := http.NewServeMux()

	//1b. Show a hostname prefix
	hname := "http://localhost:" + DEFAULT_LISTEN_PORT
	if SSL_ENABLE_FLAG {
		hname = "https://localhost:" + DEFAULT_LISTEN_PORT
	}


	//2. Now lets create a loop that goes through the DEFINITIONS that were passed in
	// and creates http handlers handlers
	for x, ep := range ALL_SERVICE_ENDPOINTS {

		if x == 0 {
			C.Println("")
		}

		C.Print(" ** Configuring REST Endpoint for: ")
		M.Println(hname + ep.API_Endpoint)
		SERVICE_MUX_OBJ.Handle(ep.API_Endpoint, GENERIC_API_ENDPOINT_HANDLER(ep.API_Handler))

		// SERVICE_MUX_OBJ.Handle(ep.API_Endpoint, ep.API_Handler)

	} //end of for


	//2. Formatting defualt port this way so we can use it in the MUX/ listener call
	colon_PORT := ":" + DEFAULT_LISTEN_PORT

	//3. If we are running in default mode.. we WONT be using SSL
	if SSL_ENABLE_FLAG == false {

			G.Println(" REST API is now Ready!")			
			
			err := http.ListenAndServe(colon_PORT, SERVICE_MUX_OBJ)

			if err != nil {
				R.Print("Start_LISTENER_SERVICE_Engine ERROR!! Trying to use PORT: ")
				C.Print(DEFAULT_LISTEN_PORT, "\n\n")
				M.Println(err.Error())
			}

	//5. Otherwise we start the service in SSL Mode
	} else {

		C.Print(" REST API")
		Y.Print(", USING SSL")
		C.Println(" is now Listening")

		err := http.ListenAndServeTLS(colon_PORT, SSL_CERT_PEM_FILE, SSL_KEY_FILE, SERVICE_MUX_OBJ)

		if err != nil {
			R.Print("Start_LISTENER_SERVICE_Engine IN SSL MODE ERROR!! Trying to use PORT: ")
			C.Print(DEFAULT_LISTEN_PORT, "\n\n")
			M.Println(err.Error())
		}
	}

} //end of Start_SERVICE




// This is a mini/micro Web Server 
func MINI_WEB(WEBROOT string, listenPort string, START_MESSAGE string) {

	listenPort = ":" + listenPort

	webURL := "http://localhost" + listenPort

	MESSAGE := "Micro WEB Server is RUNNING!"

	if START_MESSAGE != "" {
		MESSAGE = START_MESSAGE
	}

	C.Println(" - - - -", MESSAGE)
	lastChar := WEBROOT[len(WEBROOT)-1:]
	if lastChar != "/" {
		// C.Println("\n - - - INFO: Appending / (forward slash) to WEBROOT\n")
		WEBROOT = WEBROOT + "/"
	}	

	M.Print("\n     * Listening On: ")
	G.Println(webURL)
	CLIPBOARD_COPY(webURL)
	Y.Println("     (Saved to clipboard), Just paste in your browser!")

	W.Println("\n\n  .....CTRL -C to Exit!")
	W.Println("")
	W.Println("")

	http.Handle("/", http.FileServer(http.Dir(WEBROOT)))
	http.ListenAndServe(listenPort, nil)


} //end of

// Alias for MINI_WEB
func MINI_WEBROOT(WEBROOT string, listenPort string, START_MESSAGE string) {
	MINI_WEB(WEBROOT, listenPort, START_MESSAGE) 
}
// Alias for MINI_WEB
func MINI_WEB_ROOT(WEBROOT string, listenPort string, START_MESSAGE string) {
	MINI_WEB(WEBROOT, listenPort, START_MESSAGE) 
}





func init() {
  
	flag.StringVar(&DEFAULT_LISTEN_PORT, "port", DEFAULT_LISTEN_PORT,  "       The port the listener should listen on")
	flag.BoolVar(&SSL_ENABLE_FLAG,       "enableSSL", SSL_ENABLE_FLAG, "       Enables SSL support. Requires --cert and --key")
	flag.StringVar(&SSL_CERT_PEM_FILE,   "cert", SSL_CERT_PEM_FILE,    "       Full path to the CERT / PEM file for SSL (requires --enableSSL )")
	flag.StringVar(&SSL_KEY_FILE,        "key", SSL_KEY_FILE,          "       Full path to the KEY FILE for the SSL CERT (requires --enableSSL )")


	//2. Error handling for the SSL stuff
	if SSL_ENABLE_FLAG {
		if SSL_CERT_PEM_FILE == "" || SSL_KEY_FILE == "" {

			R.Println("ERROR: You MUST specify the CERT and KEY file with --cert and --key")
			R.Println("       To run in SSL mode ")

			os.Exit(-9)
		}
	}

}
