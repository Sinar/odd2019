package cmd

import (
	"fmt"
	"io/ioutil"
	"regexp"

	"github.com/gocolly/colly"
	"github.com/gocolly/colly/queue"
	"github.com/y0ssar1an/q"
	"gopkg.in/yaml.v2"
)

// ApplicationDetails has in-depth details of the Application; more than summary ApplicationRecord
type ApplicationDetails struct {
	// ID - Application ID; used to look up
	// Form URLs --> any Borang related to this Appllication; zero or more ..
	ID          ApplicationID
	Agensi      string
	Rujukan     string
	Tetuan      string
	Kategori    string
	JenisPemaju string
	RT          string
	FormRecords []FormRecord
}

// FormRecord holds details on the forms related to this Application
type FormRecord struct {
	// Bil
	// Tarikh Permohonan
	// Borang3 : Jenis Permohonan : Piagam
	// Tkh Lulus Jabatan Teknikal PBT
	// Status / Tkh Lulus / Tempoh
	bil              int
	TarikhPermohonan string
	JenisPemohonan   string
	TarikhLulus      string
	Status           string
}

// NOTE: All raw data here will be stored under the following pattern
// ./raw/<uniqueSearchID>/<ApplicationID>/

func extractApplicationDetailsData(appSnapshot *ApplicationSnapshot, pagesToExtract []string) {

}

func fetchApplicationPage(absoluteRawDataPath string, pageURL string) {
	// URL is partial; add on the needed full hostname?
	pageURL = fmt.Sprintf("http://www.epbt.gov.my/osc/%s", pageURL)
	// Extra checks will make sur egot not http/https??
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
		rePattern := regexp.MustCompile("http://www\\.epbt\\.gov\\.my/osc/Borang_info.+$")
		// // Only those with result page we grab
		if rePattern.Match([]byte(e.Request.AbsoluteURL(link))) {
			queue.AddURL(e.Request.AbsoluteURL(link))
			q.Q("FOUND BORANG: ", e.Request.AbsoluteURL(link))
		}
	})

	c.OnError(func(r *colly.Response, e error) {
		fmt.Println("DIE!!!")
		q.Q(r.Request.URL, r.StatusCode)
		panic(e)
	})

	c.OnScraped(func(r *colly.Response) {
		// // Define the result page Collector; which will just save the file
		// d := colly.NewCollector()
		// d.OnScraped(func(r *colly.Response) {
		fmt.Println("FINISH: ", r.Request.URL, "<================")
		r.Save(fmt.Sprintf("%s/%s.html", absoluteRawDataPath, r.FileName()))
		q.Q("FILE: ", r.FileName())
		q.Q("SAVED ==================>")
		// 	//fmt.Println(r.Headers)
		// })
		// // Kick off the queue once all the pages are collected already ..
		// queue.Run(d)
	})

	verr := c.Visit(pageURL)
	if verr != nil {
		panic(verr)
	}

}

func saveApplicationDetails(ad *ApplicationDetails) {
	// In ./data/<uniqueSearchID>/<applicationID>/details.yml

}

func saveApplicationRecordSummary(ar *ApplicationRecord) {
	// In ./data/<uniqueSearchID>/<applicationID>/summary.yml

}

// FetchAll will Extract from authority + label; all 15 pages of the information
func FetchAll(authorityToScrape string, forceRefresh bool, specificLabel string) {
	// Raw structure like .. ./raw/<snapshotLabel>/<uniqueSearchID>
	// NOTE: Descructive action will override the data; ensure it is git diff ..
	fmt.Println("ACTION: FetchAll")
	// Figure out when was the last successful run and if not exist; create it!
	// Also will reset if passed a flag of some sort?
	// var currentDateLabel = time.Now().Format("20060102") // "20190316"
	// var volumePrefix = "." // When in CodeFresh, it will be relative .. so that we can have the persistence
	// var uniqueSearchID = mapAuthorityToDirectory(authorityToScrape)

	// Step #1: Open from metadata tracking.yml to determine the ApplicationID
	// Open metadata structure like ./data/<uniqueSearchID>/
	// e.g. ./data/selangor-mbpj-1003/tracking.yaml; append only new unique items;
	//	sorted by ApplicationID
	// marked the successful / completed into archive? <-- Done in another step
	// Get the snapshot from memory..
	// Open up the raw data specified by Label
	var volumePrefix = "." // When in CodeFresh, it will be relative .. so that we can have the persistence
	var uniqueSearchID = mapAuthorityToDirectory(authorityToScrape)
	specifiedSnapshot := extractDataFromSnapshot(volumePrefix, specificLabel, uniqueSearchID)

	for _, singleRecord := range specifiedSnapshot.appRecords {
		// DEBUG
		// fmt.Println("Fetch URL: ", singleRecord.URL)
		// Also fetch the URLs into ./raw/<uniqueSearchID>/<applicationID>/
		absoluteRawDataPath := fmt.Sprintf("%s/raw/%s/AR_%s", volumePrefix, uniqueSearchID, singleRecord.ID)
		rawDataFolderSetup(absoluteRawDataPath)
		fetchApplicationPage(absoluteRawDataPath, singleRecord.URL)
	}

}

// FetchNew will only Extract the New items per authority mapping
func FetchNew(authorityToScrape string) {
	// Metadata structure like ./data/<uniqueSearchID>-<snapshotDiffLabels>
	// e.g. ./data/selangor-mbpj-1003-20190330_20190317
	fmt.Println("ACTION: FetchNew")
	// Figure out when was the last successful run and if not exist; create it!
	// Also will reset if passed a flag of some sort?
	// var currentDateLabel = time.Now().Format("20060102") // "20190316"
	// var volumePrefix = "." // When in CodeFresh, it will be relative .. so that we can have the persistence
	// var uniqueSearchID = mapAuthorityToDirectory(authorityToScrape)

	// Step #1: Open from metadata new.yml to determine the ApplicationID

	var volumePrefix = "." // When in CodeFresh, it will be relative .. so that we can have the persistence
	var uniqueSearchID = mapAuthorityToDirectory(authorityToScrape)
	var absoluteNewDataPath = fmt.Sprintf("./data/%s", uniqueSearchID)
	rawDataFolderSetup(absoluteNewDataPath)

	newDiff := NewDiff{}
	b, rerr := ioutil.ReadFile(fmt.Sprintf("%s/data/%s/new.yml", volumePrefix, uniqueSearchID))
	if rerr != nil {
		panic(rerr)
	}
	err := yaml.Unmarshal(b, &newDiff)
	if err != nil {
		panic(err)
	}
	fmt.Println("LABEL: ", newDiff.Label)
	for _, ar := range newDiff.AR {
		// DEBUG
		// fmt.Println("Fetch URL: ", ar.URL)
		absoluteRawDataPath := fmt.Sprintf("%s/raw/%s/AR_%s", volumePrefix, uniqueSearchID, ar.ID)
		rawDataFolderSetup(absoluteRawDataPath)
		fetchApplicationPage(absoluteRawDataPath, ar.URL)
	}

}

// ExtractAll parses the raw HTML collected under the snapshotLabel
// 	mostly is run once at the  start to kick off the process? Unless overridden
func ExtractAll(authorityToScrape string, forceRefresh bool, specificLabel string) {
	fmt.Println("ACTION: ExtractAll")

	// Save into the tracking metadata portion of it ..

	// Step #1: Open from metadata tracking.yml to determine the ApplicationID
	var volumePrefix = "." // When in CodeFresh, it will be relative .. so that we can have the persistence
	var uniqueSearchID = mapAuthorityToDirectory(authorityToScrape)
	var absoluteNewDataPath = fmt.Sprintf("./data/%s", uniqueSearchID)
	rawDataFolderSetup(absoluteNewDataPath)

	newDiff := NewDiff{}
	b, rerr := ioutil.ReadFile(fmt.Sprintf("%s/data/%s/tracking.yml", volumePrefix, uniqueSearchID))
	if rerr != nil {
		panic(rerr)
	}
	err := yaml.Unmarshal(b, &newDiff)
	if err != nil {
		panic(err)
	}
	fmt.Println("LABEL: ", newDiff.Label)
	for _, ar := range newDiff.AR {
		fmt.Println("Fetch URL: ", ar.URL)
	}

}

// ExtractNew parses the raw HTML files for the new ranges
func ExtractNew(authorityToScrape string) {
	// Metadata structure like ./data/<uniqueSearchID>-<snapshotDiffLabels>
	// e.g. ./data/selangor-mbpj-1003-20190330_20190317
	fmt.Println("ACTION: ExtractNew")

	// TODO: Load the tracking.yml

	// TODO: Persist after appending the new items; of maybe just append direct
}

// extract out the fields related to Application; which is what?
func extractAllApplicationDetails() {

}
