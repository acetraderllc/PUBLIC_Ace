/*
        Abba - Abba is great swedish band..Abba is also a GoLANG Swiss-Army-Knife type Multi Tool!

        This Tool aims to easy the following:
        - Installing GO for your platform (if you dont have it already) by way of:
                - Uses HOMEBREW on mac
                - Uses CHOCOLATEY on windows
                - Uses yum on RHEL/Centos/Amazon
        - Make it easy for dealing with Go Modules (GOMOD) particularly on bitbucket and github.. and especially dealing with private modules
        - Makes it easy for you to manage 3rd Party Packages (so you know what packages you are using)
        - Provides an easy scaffold for building NEW go programs
        - Make it easy to deal with cross platform building
        - Make it easy to manage your local development (environment variables etc)

        
*/
package main


import (
        // Standard
        "flag"

        // Personal
        . "bitbucket.org/cowboytarik/GLOBAL/Terry_COMMON"
	. "bitbucket.org/cowboytarik/private/Eagle"

	. "github.com/terrycowboy/IGGY/Royals"

        // 3RD Party

)

// -=-=-=--=  SOME GLobals to make this thing work

var INTERACTIVE_MODE = false
var BUILD_MAC = true
var BUILD_WINDOWS = false
var BUILD_LINUX = false
var OUTPUT = ""

func main() {

        //First, lets setup some command line parameters
        flag.BoolVar(&INTERACTIVE_MODE, "interactive", INTERACTIVE_MODE, "  Run in interactive mode")

        flag.BoolVar(&BUILD_MAC, "mac", BUILD_MAC, "  Compile a binary for MAC")
        flag.BoolVar(&BUILD_WINDOWS, "win", BUILD_WINDOWS, "  Compile a binary for WINDOWS")
        flag.BoolVar(&BUILD_WINDOWS, "windows", BUILD_WINDOWS, "  Alias for --win")
        flag.BoolVar(&BUILD_LINUX, "linux", BUILD_LINUX, "  Compile a binary for Linux")

        flag.StringVar(&OUTPUT, "out", OUTPUT, "  Full path (if applicable) of the output binary you want. Default is name of Current Directory)")
        flag.StringVar(&OUTPUT, "output", OUTPUT, "  Alias for --out")

        MASTER_INIT("Army", 1.26)

	M.Println(" Here it is: ", ALPHA_API_KEY)

        //1. Ok, lets determine what platform we are on.. it some defaults we want to be set based on platform
        ShowPlatform()
        
 	Tonight()
 }

