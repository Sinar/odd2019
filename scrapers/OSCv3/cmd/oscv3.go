package main

import (
	"fmt"
	"os"

	"github.com/gocolly/colly"
	"github.com/y0ssar1an/q"
)

func main() {
	fmt.Println("Welcome to GOMOD OSCv3!!")
	BasicCollyFromRaw()
}

func rawDataFolderSetup(absoluteRawDataPath string) {
	fi, lerr := os.Stat(absoluteRawDataPath)
	if lerr != nil {
		if os.IsNotExist(lerr) {
			// Create the needed folder as per needed .. all along the chain
			mkerr := os.MkdirAll(absoluteRawDataPath, 0700)
			if mkerr != nil {
				panic(mkerr)
			}
		} else {
			// SOme other unknown issue
			panic(lerr)
		}
	} else {
		// spew.Dump(fi)
		if !fi.IsDir() {
			panic(fmt.Errorf("NOT DIR!! %s", absoluteRawDataPath))
		}
		fmt.Println("Directory ", absoluteRawDataPath, " EXISTS!!")
	}
}

// BasicCollyFromRaw is meant to read data from fixtures and extract out data ..
func BasicCollyFromRaw() {
	// From the collection; can run another round while tweaking the strcuture
	// Removing the extra cost of network and being blocked ..
	var currentDateLabel = "20190316"
	var uniqueSearchID = "penang-kulim-0212"
	// NOTE: Won't work on Windoze :(
	var absoluteRawDataPath = fmt.Sprintf("raw/%s/%s", currentDateLabel, uniqueSearchID)
	rawDataFolderSetup(absoluteRawDataPath)

	// With pre-reqs setup; we can proceed ...
	c := colly.NewCollector(
		colly.UserAgent("Sinar Project :P"),
	)
	c.OnError(func(r *colly.Response, e error) {
		fmt.Println("DIE!!!")
		q.Q(r.Request.URL, r.StatusCode)
		panic(e)
	})
}
