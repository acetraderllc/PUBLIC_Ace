/*
Thready - GoLang Ro-routine Threading wrapper
---------------------------------------------------------------------------------------
	v1.25	-	Dec 29, 2019	-	Converted this to a GO MOdule that lives in Github/Bitbucket. No more local stuff!
	v1.23	- 	Nov 05, 2016	-	Initial Release

*/

package Thready

import (
	
	// - - - STANDARD Library - - - -
		"context"
		"time"
		"flag"
		"strings"
		"os"

	// - - - 3RD Party MODULES - - - -  /*
		"go.mongodb.org/mongo-driver/mongo"
		"go.mongodb.org/mongo-driver/mongo/options"	
		"go.mongodb.org/mongo-driver/x/bsonx"	
		. "bitbucket.org/cowboytarik/SHARED/Terry_COMMON"
*/
)



func THREAD_MASTER(CHILD_FUNC_NAME string) {


	//C.Println(" MAX BARS IS: ", MAX_BARS)


	C.Println(" *** Starting MULTI-THREADED Signal Generation Routine ***")
	C.Print(" *** ")
	W.Print("Total ANALYSIS Records:  ")
	G.Println(ShowNum(MASTER_TOTAL))
	M.Println("")

	var TOTAL_ITEMS_TO_PROCESS = MASTER_TOTAL
	if MAX_THREADS > TOTAL_ITEMS_TO_PROCESS {
		MAX_THREADS = TOTAL_ITEMS_TO_PROCESS
	}

	var WAIT_GROUP sync.WaitGroup
	WAIT_GROUP.Add(MAX_THREADS)

	var THREAD_COUNTER = 0

	startTime, startOBJ := GET_CURRENT_TIME()	

	//2. Now lets iterate through this list.. there may be dups
	for thread_INDEX, analOBJ := range ANALYSIS_LIST {

		//2b. Bump thread_INDEX (so we show threads at 1 instead of 0)
		thread_INDEX++
		thread_INDEX_TEXT := ShowNumber(thread_INDEX)

		go CHILD_SIG_GEN_Routine(analOBJ, &WAIT_GROUP, thread_INDEX_TEXT)
		THREAD_COUNTER++
		TOTAL_ITEMS_TO_PROCESS--
		Sleep(1, false)

		//3. Error Hanlding if we are on the LAST item to process
		if TOTAL_ITEMS_TO_PROCESS == 0 {
			WAIT_GROUP.Wait()
			break
		}

	}
} //end of MAIN






var TDB_HOST
func init() {

	flag.StringVar(&TDB_HOST, "tdbhost", TDB_HOST, " Mongo DB HOST")
}

