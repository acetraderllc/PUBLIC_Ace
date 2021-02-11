package Vault_OPS

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

)

var VAULT_USER = ""
var VAULT_PASS = ""
var VAULT_URL = "https://XXXX.YYYYYYYYY.com:8200"
var MAX_VAULT_LIFE = "12h"

func VAULT_PARAMS_INIT() {

	//1. Define the common command line params we use for AWS connectivity

	flag.StringVar(&VAULT_USER, "vuser", VAULT_USER, "       Vault USER")
	flag.StringVar(&VAULT_PASS, "vpass", VAULT_PASS, "       Vault PASWORD")
	flag.StringVar(&VAULT_URL, "region", VAULT_URL,         "       URL To the Vault Server")
	flag.StringVar(&MAX_VAULT_LIFE, "life", MAX_VAULT_LIFE, "       The Maximum Vault Session Lifespan")

	//3. If neither ACCESS/SECRET and PROFILE are supplied.. we default to the ENVIRONMENT variables	
	if VAULT_USER == "" && VAULT_PASS == "" {

		VAULT_USER = os.Getenv("VAULT_USER")
		VAULT_PASS = os.Getenv("VAULT_PASS")

		if VAULT_USER == "" {
			R.Println("\n")
			R.Println(" ERROR: You need to specify --vuser and --vpass ... or set the VAULT_USER and VAULT_PASS environment variables")
			os.Exit(-9)
		}
	}
} //end of func


func VAULT_QUERY_ENGINE(VAULT_PROFILE string) (string, string, string, string) {

	Y.Print(" *** Querying VAULT with Profile for: ")
	C.Println(VAULT_PROFILE)

	url := VAULT_URL + "/v1/auth/ldap/login/" + VAULT_USER
	jsonStrParams := []byte(`{"password":"` + VAULT_PASS + `"}`)

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStrParams))
	req.Header.Set("Content-Type", "application/json")

	//1. Error Handling:
	if err != nil {

		M.Println("\n\n ** ERROR: Cannot POST to ", VAULT_URL)
		M.Println("")

		return "", "", "", ""
	}

	client := &http.Client{}
	resp, err2 := client.Do(req)
	if err2 != nil {
		M.Println("\n\n ** ERROR: With CLIENT response from  ", VAULT_URL)
		M.Println("")

		return "", "", "", ""
	}

	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)

	JSON_RESPONSE := string(body)

	msplit := strings.Split(JSON_RESPONSE, "\"")

	mlen := len(msplit)

	// C.Println("JSON RESP", JSON_RESPONSE, " MLEN", mlen)

	if mlen < 24 {
		M.Println("\n\n ** ERROR in response we got back.. not enough data: ")
		M.Println("")
		return "", "", "", ""
	}

	//2. Pull out the client token
	CLIENT_TOKEN := msplit[23]

	//3. Now lets call vault again and get some aws credentials
	DESTOBJ := VAULT_PROFILE + "/sts/sts-AdministratorAccess"

	nexturl := VAULT_URL + "/v1/" + VAULT_PROFILE + "/sts/sts-AdministratorAccess"
	jsonStrParams = []byte(`
        {
          "password":"` + VAULT_PASS + `",
          "write":"` + DESTOBJ + `",
          "ttl":"` + LIFESPAN + `"      
        }    
    `)

	req2, err3 := http.NewRequest("POST", nexturl, bytes.NewBuffer(jsonStrParams))
	req2.Header.Set("Content-Type", "application/json")
	req2.Header.Set("X-Vault-Token", CLIENT_TOKEN)

	if err3 != nil {
		M.Println("\n\n ** ERROR sending SECONDARY post to ", VAULT_URL)
		M.Println("")

		return "", "", "", ""
	}

	client2 := &http.Client{}
	resp2, err4 := client2.Do(req2)
	if err4 != nil {
		M.Println("\n\n ** ERROR with RESPONSE back from VAULT_URL ", VAULT_URL)
		M.Println("")

		return "", "", "", ""
	}

	defer resp2.Body.Close()

	body2, _ := ioutil.ReadAll(resp2.Body)
	FINAL_JSON_RESPONSE := string(body2)

	//4. Now split the json we got back.. The values are static so we can use positional references to the values
	ksplit := strings.Split(FINAL_JSON_RESPONSE, "\"")

	V_ACCESS_KEY := ksplit[17]
	V_ACCESS_SECRET := ksplit[21]
	V_ACCESS_TOKEN := ksplit[25]

	//6. We now have the access, secret and token.. lets put them in a json structure and return them as output
	OUTGOING_JSON := `
    {
         "ACCESS":"` + V_ACCESS_KEY + `",
         "SECRET":"` + V_ACCESS_SECRET + `",
          "TOKEN":"` + V_ACCESS_TOKEN + `",
    },`

	return OUTGOING_JSON, V_ACCESS_KEY, V_ACCESS_SECRET, V_ACCESS_TOKEN
} //end of func



// Kept here as filler/example.. anything you put in this function will start when the module is imported
func init() {

	//1. Startup Stuff (init the command line params etc)

} // end of main

