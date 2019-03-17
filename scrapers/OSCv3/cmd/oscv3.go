package cmd

import (
	"fmt"
	"os"
	"regexp"
	"time"

	"github.com/gocolly/colly"
	"github.com/gocolly/colly/queue"
	"github.com/y0ssar1an/q"
)

func rawDataFolderSetup(absoluteRawDataPath string) (proceedScraping bool) {
	fi, lerr := os.Stat(absoluteRawDataPath)
	if lerr != nil {
		if os.IsNotExist(lerr) {
			// Create the needed folder as per needed .. all along the chain
			mkerr := os.MkdirAll(absoluteRawDataPath, 0700)
			if mkerr != nil {
				panic(mkerr)
			}
			return true
		}
	} else {
		// spew.Dump(fi)
		if !fi.IsDir() {
			panic(fmt.Errorf("NOT DIR!! %s", absoluteRawDataPath))
		}
		fmt.Println("Directory ", absoluteRawDataPath, " EXISTS!!")
		// When the folder already exist for the day, no need to proceed
		fmt.Println("Skipping... ")
		// Should turn to false once tested .. no need to re-run
		// return true
	}

	return false
}

// BasicCollyFromRaw is meant to read data from fixtures and extract out data ..
func BasicCollyFromRaw() {
	// From the collection; can run another round while tweaking the strcuture
	// Removing the extra cost of network and being blocked ..
	var currentDateLabel = time.Now().Format("20060102") // "20190316"
	var uniqueSearchID = "penang-kulim-0212"
	var volumePrefix = "." // When in CodeFresh, it will be relative .. so that we can have the persistence
	// NOTE: Won't work on Windoze :(
	var absoluteRawDataPath = fmt.Sprintf("%s/raw/%s/%s", volumePrefix, currentDateLabel, uniqueSearchID)

	// go, no go?
	proceedScraping := rawDataFolderSetup(absoluteRawDataPath)

	if proceedScraping {
		// Setup the queue that will be to grab at the available pages up till the 15 pages limit
		queue, _ := queue.New(
			2, // Number of consumer threads
			&queue.InMemoryQueueStorage{MaxSize: 10000}, // Use default queue storage
		)
		// Add the seedURL to the queue to be saved; not needed as we will add below
		// queue.AddURL(seedURL)

		// With pre-reqs setup; we can proceed ...
		c := colly.NewCollector(
			colly.UserAgent("Sinar Project :P"),
		)

		// On every a element which has href attribute print full link
		c.OnHTML("a[href]", func(e *colly.HTMLElement) {
			link := e.Attr("href")
			rePattern := regexp.MustCompile("http://www\\.epbt\\.gov\\.my/osc/Carian_Proj3.+$")
			// Only those with result page we grab
			if rePattern.Match([]byte(e.Request.AbsoluteURL(link))) {
				queue.AddURL(e.Request.AbsoluteURL(link))
			}
		})

		c.OnError(func(r *colly.Response, e error) {
			fmt.Println("DIE!!!")
			q.Q(r.Request.URL, r.StatusCode)
			panic(e)
		})

		c.OnScraped(func(r *colly.Response) {
			// Define the result page Collector; which will just save the file
			d := colly.NewCollector()
			d.OnScraped(func(r *colly.Response) {
				fmt.Println("FINISH: ", r.Request.URL, "<================")
				r.Save(fmt.Sprintf("%s/%s.html", absoluteRawDataPath, r.FileName()))
				q.Q("FILE: ", r.FileName())
				q.Q("SAVED ==================>")
				//fmt.Println(r.Headers)
			})
			// Kick off the queue once all the pages are collected already ..
			queue.Run(d)
		})

		// Start with the seedURL to kick things off
		seedURL := "http://www.epbt.gov.my/osc/Carian_Proj3.cfm?CurrentPage=1&Maxrows=2&Cari=&AgensiKod=0212&Pilih=3"
		verr := c.Visit(seedURL)
		if verr != nil {
			panic(verr)
		}
	}
}
