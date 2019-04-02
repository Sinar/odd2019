package main

import (
	"flag"
	"fmt"

	"github.com/sinar/odd2019/scrapers/OSCv3/cmd"
)

func main() {
	fmt.Println("Welcome to GOMOD OSCv3!!")
	// TODO: Use github.com/mitchellh/cli for cli
	// For now just use the simple flag package?
	actionPtr := flag.String("action", "update", "What action to run: default is update, you can call: diff")
	authorityPtr := flag.String("authority", "0212", "Which Local Authority to scrape? MBPJ - 1003, Kulim - 0212")
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
		specificLabel := "20190317"
		cmd.FindAllApplications(*authorityPtr, forceRefresh, specificLabel)
		return
	}

	fmt.Println("INVALID ACTION: ", *actionPtr)
	fmt.Println("VALID: update, new, diff")
}
