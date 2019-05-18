package main

import (
	"flag"
	"fmt"
	"time"

	"github.com/sinar/odd2019/scrapers/OSCv3/cmd"
)

func main() {
	fmt.Println("Welcome to GOMOD OSCv3!!")
	// TODO: Use github.com/mitchellh/cli for cli
	// For now just use the simple flag package?
	actionPtr := flag.String("action", "update", "What action to run: default is update, you can call: diff")
	authorityPtr := flag.String("authority", "1007", "Which Local Authority to scrape? MBPJ - 1003, Kulim - 0212, DBKL - 9999")
	flag.Parse()

	if *actionPtr == "update" {
		cmd.BasicCollyFromRaw(*authorityPtr)
		return
	} else if *actionPtr == "diff" {
		cmd.FindNewRequests(*authorityPtr)
		return
	} else if *actionPtr == "track" {
		// use a specific option label like
		forceRefresh := false
		// specificLabel := "20190413"
		// Set to current for now; state transition/history will be lost!
		specificLabel := time.Now().Format("20060102") // "20190407"
		cmd.FindAllApplications(*authorityPtr, forceRefresh, specificLabel)
		return
	} else if *actionPtr == "fetch" {

		cmd.FetchNew(*authorityPtr)
		return
	} else if *actionPtr == "fetchall" {
		// TODO: Maybe from the track only??
		// use a specific option label like
		forceRefresh := false
		specificLabel := "20190413"
		// Set to current for now; state transition/history will be lost!
		// specificLabel := time.Now().Format("20060102") // "20190407"
		cmd.FetchAll(*authorityPtr, forceRefresh, specificLabel)
		return
	} else if *actionPtr == "extract" {

		cmd.ExtractNew(*authorityPtr)
		return
	} else if *actionPtr == "extractall" {

		cmd.ExtractAll(*authorityPtr)
		return
	} else if *actionPtr == "borang" {
		// Extract new borang; based on  data not active? needs to be activated?
		//  See the status?
		cmd.ExtractFormNew(*authorityPtr)
		return
	}

	fmt.Println("INVALID ACTION: ", *actionPtr)
	fmt.Println("VALID: update, new, diff")
}
