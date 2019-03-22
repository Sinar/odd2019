package cmd

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gocolly/colly"
)

func loadMetaData() {
	// IN yaml format
	// Tells us the last unique ID that was processed/seen
}

func saveMetaData() {
	// In yaml format
	// saves the first unique ID seen; assuming this is called once it is successful!
}

func saveData() {
	//IN yaml format

}

func extractDataFromPage() {

}

func extractAllData(pagesToExtract []string) {
	// Loop in the whole identified folder ..
	// and run extractDataFromPage
	c := colly.NewCollector(
		colly.UserAgent("Sinar Project :P"),
		colly.Async(true),
	)
	// As per in: https://github.com/gocolly/colly/issues/260
	// can register local file: transport with the absolute path?
	t := &http.Transport{}
	t.RegisterProtocol("file", http.NewFileTransport(http.Dir("/")))
	c.WithTransport(t)

	c.OnScraped(func(r *colly.Response) {
		// Next time; queue the further extraction of item?
		// spew.Println(r.StatusCode)
		fmt.Println("DONE!  CODE:", r.StatusCode)
	})

	// Example finalURL will be like below:
	// finalURL := "file:///Users/leow/GOMOD/odd2019/scrapers/OSCv3/raw/20190322/selangor-mbpj-1003/_osc_Carian_Proj3.cfm_CurrentPage_10_Maxrows_15_Cari_AgensiKod_1003_Pilih_3.html"
	for _, url := range pagesToExtract {
		finalURL := fmt.Sprintf("file://%s", url)
		// DEBUG
		// fmt.Println("FILE: ", finalURL)
		verr := c.Visit(finalURL)
		if verr != nil {
			// panic(verr)
			fmt.Println("ERR:", verr.Error())
		}
	}

	c.Wait() // Barrier for aync; so we can go as fast as opossible ..
	// Sort the final based on the BIL as Int
	// Then iterate until the last observed item is matched!
}

// FindNewRequests will look for the changes since the last time run and offer a pull request
func FindNewRequests(authorityToScrape string) {
	fmt.Println("ACTION: FindNewRequests")
	// Figure out when was the last successful run and if not exist; create it!
	// Also will reset if passed a flag of some sort?
	// var currentDateLabel = time.Now().Format("20060102") // "20190316"
	var currentDateLabel = "20190322"
	var previousDateLabel = "20190317"
	var uniqueSearchID = mapAuthorityToDirectory(authorityToScrape)
	var volumePrefix = "." // When in CodeFresh, it will be relative .. so that we can have the persistence
	// NOTE: Won't work on Windoze :(
	var absoluteRawDataPath = fmt.Sprintf("%s/raw/%s/%s", volumePrefix, currentDateLabel, uniqueSearchID)

	dir, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	// DEBUG
	// fmt.Println("PWD: ", dir, " DIR:", absoluteRawDataPath)
	absoluteRawDataPath = fmt.Sprintf("%s/%s", dir, absoluteRawDataPath)
	// Iterate through the raw data .. and append the name to the map; we know it is always 15 pages max
	pages := []string{}
	fi, rerr := ioutil.ReadDir(absoluteRawDataPath)
	if rerr != nil {
		panic(rerr)
	}
	for _, f := range fi {
		if !f.IsDir() {
			path := filepath.Join(absoluteRawDataPath, "/", f.Name())
			// DEBUG
			// fmt.Println("FILE: ", path)
			pages = append(pages, path)
		}
	}
	// fmt.Println("=============================++******")
	// sort.Strings(pages)
	// fmt.Println(pages)
	// filepath.Walk(absoluteRawDataPath, func(path string, info os.FileInfo, err error) error {
	// 	if !info.IsDir() {
	// 		fmt.Println("PATh: ", path)
	// 	}
	// 	return nil
	// })
	extractAllData(pages)
	// If in Codefresh; do a branch, git add + commit?
	fmt.Println("Now compare against the previous: ", previousDateLabel)
}
