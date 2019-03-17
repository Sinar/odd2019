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
	flag.Parse()

	if *actionPtr == "update" {
		cmd.BasicCollyFromRaw()
		return
	} else if *actionPtr == "diff" {
		cmd.FindNewRequests()
		return
	}

	fmt.Println("INVALID ACTION: ", *actionPtr)
	fmt.Println("VALID: update, new, diff")
}
