package BaccShared

import (
	
	// ** IMPORT **  Native Libraries
	// "sort"
	// "strings"
	// "time"
	// //	"unicode"
	// //"errors"
	// "math"
	// // "os"
	// "strconv"

	// ** IMPORT **  Personal / Custom  Libraries	
	// . "TerryCommon"

	// ** IMPORT **  3rd Party Libraries
	// "gopkg.in/mgo.v2/bson"
)

// format of the Baccrat HISTORY table of hands

type HAND_HISTORY_OBJ struct {
	SESS_ID						string 		`bson:"SESS_ID"`
	Date							string 		`bson:"Date"`

   Hand_NUM   					string		`bson:"Hand_NUM"`

	WINNER						string		`bson:"WINNER"`
	PLAYER_HAND_PATTERN 		string		`bson:"PLAYER_HAND_PATTERN"`
	Banker_HAND_PATTERN 		string		`bson:"Banker_HAND_PATTERN"`	


	PLAYER_First				string		`bson:"PLAYER_First"`
	PLAYER_Second				string		`bson:"PLAYER_Second"`
	PLAYER_Third				string		`bson:"PLAYER_Third"`
	PLAYER_SCORE				string		`bson:"PLAYER_SCORE"`

	Banker_First 				string		`bson:"Banker_First"`        
	Banker_Second 				string      `bson:"Banker_Second"`  
	Banker_Third 				string		`bson:"Banker_Third"`
	Banker_SCORE				string		`bson:"Banker_SCORE"`

	PLAYER_WINS_SOFAR			int		`bson:"PLAYER_WINS_SOFAR"`
	Banker_WINS_SOFAR			int		`bson:"Banker_WINS_SOFAR"`
	

	PLAYER_NATURALS_SOFAR	int 		`bson:"Player_NATURALS_SOFAR"`
	Banker_NATURALS_SOFAR	int 		`bson:"Banker_NATURALS_SOFAR"`
	

	TIES_SOFAR					int		`bson:"TIES_SOFAR"`
}

var MAX_DECK_CARDS = 52 		// Maximum number of cards in a deck  52 pickup

var CARDS = []string {"A", "2", "3", "4",	"5", "6", "7",	"8", "9", "10", "J", "Q", "K"}
var SUITES = []string {"CLOVER_BLACK", "DIAMOND_RED", "HEART_RED", "SPADE_BLACK"}


type CARD_OBJ struct {
		Num_Suite   string
		Value 		int
}	


// Used for passing functions.. this is a type of function you can pass around
// determine if it is accurate or not
// Returns TRUE if it accurately detected the hand
//type STRATEGY_ANALYZER_LOGIC_DEFINITION func(hands []HISTORY_OBJ, FINAL_HAND HISTORY_OBJ) bool

