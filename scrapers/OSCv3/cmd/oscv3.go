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
		panic(lerr)
	} else {
		// spew.Dump(fi)
		if !fi.IsDir() {
			panic(fmt.Errorf("NOT DIR!! %s", absoluteRawDataPath))
		}
		fmt.Println("Directory ", absoluteRawDataPath, " EXISTS!!")
		// When the folder already exist for the day, no need to proceed
		fmt.Println("Skipping... ")
	}

	return false
}

func mapAuthorityToDirectory(authorityID string) string {
	var directoryName string
	switch authorityID {
	case "1003":
		fmt.Println("MBPJ!!")
		directoryName = fmt.Sprintf("selangor-mbpj-%s", authorityID)
	case "1007":
		fmt.Println("MPSJ!!")
		directoryName = fmt.Sprintf("selangor-mpsj-%s", authorityID)
	case "0212":
		fmt.Println("KULIM!!")
		directoryName = fmt.Sprintf("penang-kulim-%s", authorityID)
	case "9999":
		fmt.Println("DBKL!!")
		directoryName = fmt.Sprintf("kl-dbkl-%s", authorityID)
	default:
		fmt.Println("INVALID AUTHORITY: ", authorityID, " Maybe 1003 for MBPJ?? Or 9999 for DBKL? Or 0212 for Kulim?")
		panic("BAD AUTHORITY!!!")
	}

	return directoryName
}

// BasicCollyFromRaw is meant to read data from fixtures and extract out data ..
func BasicCollyFromRaw(authorityToScrape string) {
	fmt.Println("ACTION: BasicCollyFromRaw for Authority - ", authorityToScrape)
	// From the collection; can run another round while tweaking the strcuture
	// Removing the extra cost of network and being blocked ..
	var currentDateLabel = time.Now().Format("20060102") // "20190316"
	var uniqueSearchID = mapAuthorityToDirectory(authorityToScrape)
	var volumePrefix = "." // When in CodeFresh, it will be relative .. so that we can have the persistence
	// NOTE: Won't work on Windoze :(
	var absoluteRawDataPath = fmt.Sprintf("%s/raw/%s/%s", volumePrefix, currentDateLabel, uniqueSearchID)
	// Start with the seedURL to kick things off
	// TODO: Use the helper sling to make it better for API caller?
	seedURL := fmt.Sprintf("http://www.epbt.gov.my/osc/Carian_Proj3.cfm?CurrentPage=1&Maxrows=2&Cari=&AgensiKod=%s&Pilih=3", authorityToScrape)

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
		// Make it wait longer .. crazy 300 secs! This is needed for MPSJ; DBKL cannot totally!
		//c.SetRequestTimeout(300 * time.Second)
		//c.WithTransport(&http.Transport{
		//	Proxy: http.ProxyFromEnvironment,
		//	DialContext: (&net.Dialer{
		//		Timeout:   300 * time.Second,
		//		KeepAlive: 30 * time.Second,
		//		DualStack: true,
		//	}).DialContext,
		//	MaxIdleConns:          100,
		//	IdleConnTimeout:       90 * time.Second,
		//	TLSHandshakeTimeout:   10 * time.Second,
		//	ExpectContinueTimeout: 1 * time.Second,
		//	ResponseHeaderTimeout: 300 * time.Second,
		//})

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

		verr := c.Visit(seedURL)
		if verr != nil {
			panic(verr)
		}
	}
}
