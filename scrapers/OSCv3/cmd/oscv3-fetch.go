package cmd

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/gocolly/colly"
	"github.com/gocolly/colly/queue"
	"github.com/y0ssar1an/q"
	"gopkg.in/yaml.v2"
)

// ApplicationDetails has in-depth details of the Application; more than summary ApplicationRecord
type ApplicationDetails struct {
	// ID - Application ID; used to look up
	// Form URLs --> any Borang related to this Appllication; zero or more ..
	AR ApplicationRecord
	// Agensi      string
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
	URL              string
	TarikhPermohonan string
	JenisPemohonan   string
	TarikhLulus      string
	Status           string
}

// NOTE: All raw data here will be stored under the following pattern
// ./raw/<uniqueSearchID>/<ApplicationID>/

func extractApplicationDetailsData(appDetails *ApplicationDetails, pagesToExtract []string) {
	fmt.Println("START ==> extractApplicationDetailsData =================")
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

	// Pattern for Application Details
	c.OnHTML("body > table > tbody > tr > td > table:nth-child(2) > tbody > tr:nth-child(3) > td > table > tbody > tr:nth-child(2) > td > table > tbody", func(e *colly.HTMLElement) {

		// DEBUG
		// q.Q("-- START DETAILS MATCH ---")
		// q.Q(e.Text)
		// q.Q("== END MATCH ====")

		// Every row of data ..
		e.DOM.Children().Each(func(i int, s *goquery.Selection) {
			// Reset items per row; like key back to empty
			colKey := ""

			// Pattern is key : Value
			// Each column
			s.Children().Each(func(j int, c *goquery.Selection) {
				// Reset columnVal
				colValue := ""
				if j == 0 {
					colKey = strings.TrimSpace(c.Text())
					// DEBUG
					// q.Q("KEY:", colKey)
				} else if j == 1 {
					// Col separator .. skip ..
				} else if j == 2 {
					colValue = strings.TrimSpace(c.Text())
					// DEBUG
					// q.Q("VALUE:", colValue)
					// based on columnKey
					switch colKey {
					case "Agensi":
						// Do nothing
						// appDetails.Agensi = colValue
						// Reset colKey for next round
						colKey = ""

					case "No. Rujukan Projek":
						appDetails.Rujukan = colValue
						// Reset colKey for next round
						colKey = ""

					case "Nama Projek":
						// Do nothing if exist
						if appDetails.AR.Projek == "" {
							appDetails.AR.Projek = colValue
						}
						// Reset colKey for next round
						colKey = ""
					case "Untuk Tetuan":
						appDetails.Tetuan = colValue
						// Reset colKey for next round
						colKey = ""

					case "Kategori Projek":
						appDetails.Kategori = colValue
						// Reset colKey for next round
						colKey = ""

					case "Jenis Pemajuan":
						appDetails.JenisPemaju = colValue
						// Reset colKey for next round
						colKey = ""

					case "Rancangan Tempatan (RT)":
						appDetails.RT = colValue
						// Reset colKey for next round
						colKey = ""

					case "No. Lot":
						// Do nothing if exist
						if appDetails.AR.Lot == "" {
							appDetails.AR.Lot = colValue
						}
						// Reset colKey for next round
						colKey = ""
					case "Mukim":
						// Do nothing if exist
						if appDetails.AR.Mukim == "" {
							appDetails.AR.Mukim = colValue
						}
						// Reset colKey for next round
						colKey = ""
					default:
						// Do nothing
						// Reset colKey for next round
						colKey = ""
					}
				} else {
					q.Q("UNKNOWN COL:", j, " VAL: ", strings.TrimSpace(c.Text()))
				}
			})
		})
	})

	// Pattern for Form Summary
	c.OnHTML("body > table > tbody > tr > td > table:nth-child(2) > tbody > tr:nth-child(3) > td > table > tbody > tr:nth-child(3) > td > table > tbody", func(e *colly.HTMLElement) {

		// DEBUG
		// q.Q("-- START FORM MATCH ---")
		// q.Q(e.Text)
		// q.Q("== END MATCH ====")
		formRecords := make([]FormRecord, 0)

		// Every row of data ..
		e.DOM.Children().Each(func(i int, s *goquery.Selection) {
			// Initialize for this row iteration ..
			formRecord := FormRecord{}
			//formBil := 0
			// Pattern is key : Value
			// Each Form Column
			s.Children().Each(func(j int, fc *goquery.Selection) {
				// DEBUG
				// q.Q("-- Start FORM Child -- ", j)
				// q.Q(fc.Html())
				// q.Q(strings.TrimSpace(fc.Text()))
				// q.Q("== END FORM Child =====")
				if j == 0 {
					// TODO: Better check??
					// q.Q("BIL: ", fc.Text())
					formRecord.bil, _ = strconv.Atoi(strings.Trim(fc.Text(), "."))
					// TODO: Try to extract out the URL matching the pattern "Borang_info.cfm"
					fc.Find("a").Map(func(_ int, m *goquery.Selection) string {
						href, _ := m.Attr("href")
						idURL, err := url.Parse(href)
						if err != nil {
							panic(err)
						}
						formRecord.URL = href
						//DEBUG
						// q.Q("URL :>>> ", formRecord.URL)
						// fmt.Println(appRecord.URL)
						return idURL.Query().Get("ID")
					})
				} else if j == 1 {
					// Application date
					formRecord.TarikhPermohonan = strings.Trim(fc.Text(), "\n\t ")
				} else if j == 2 {
					// Application Type; including SLA?
					formRecord.JenisPemohonan = strings.Trim(fc.Text(), "\n\t ")
				} else if j == 3 {
					// Approval date; this will be dynamic; start out as empty?
				} else if j == 4 {
					// Status
					formRecord.Status = strings.Trim(fc.Text(), "\n\t ")
				} else {
					// Something is wrong!
					panic(fmt.Errorf("Too many columns?? Col: %d", j))
				}
			})
			// Got out; append the data; if it is not header
			// Don;t append if it is a weird structure instead??
			// Or if bil is not positive int??
			q.Q(formRecord)
			// TODO: Remove this stupid hack  ..
			if formRecord.TarikhPermohonan != "Tarikh Permohonan" {
				// formRecord.bil = 1
				formRecords = append(formRecords, formRecord)
			}
			//formBil++
		})

		// Attach all found records ..
		appDetails.FormRecords = formRecords
	})

	c.OnScraped(func(r *colly.Response) {
		// Next time; queue the further extraction of item?
		// spew.Println(r.StatusCode)
		// DEBUG
		// fmt.Println("DONE!  CODE:", r.StatusCode)
		q.Q(appDetails)
	})

	// Example finalURL will be like below:
	// finalURL := "file:///Users/leow/GOMOD/odd2019/scrapers/OSCv3/raw/selangor-mbpj-1003/AR_770177/_osc_Proj1_Info.cfm_Name_770177_S_S.html"
	for _, URL := range pagesToExtract {
		finalURL := fmt.Sprintf("file://%s", URL)
		// DEBUG
		// fmt.Println("FILE: ", finalURL)
		verr := c.Visit(finalURL)
		if verr != nil {
			// panic(verr)
			fmt.Println("ERR:", verr.Error())
		}
		// DEBUG
		// break
	}

	c.Wait() // Barrier for aync; so we can go as fast as opossible ..
	// Sort the final based on the bil as Int
	// Then iterate until the last observed item is matched!
	// Whats doe sit look like?
	// DEBUG
	// fmt.Println(strings.TrimSpace(strings.Join(idOrder, ",")))

}

func fetchApplicationPage(absoluteRawDataPath string, pageURL string) {
	// DEBUG
	//fmt.Println("Inside fetchApplicationPage ..")
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
		rePattern := regexp.MustCompile(`http://www.epbt.gov.my/osc/Borang_info.+$`)
		// // Only those with result page we grab
		if rePattern.Match([]byte(e.Request.AbsoluteURL(link))) {
			err := queue.AddURL(e.Request.AbsoluteURL(link))
			if err != nil {
				panic(err)
			}
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
		err := r.Save(fmt.Sprintf("%s/%s.html", absoluteRawDataPath, r.FileName()))
		if err != nil {
			panic(err)
		}
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

func isApplicationPageActive(pageURL string) bool {
	fmt.Println("pageURL: ", pageURL)
	fmt.Println("Inside isApplicationPageActive ...")
	var foundAgensi string
	c := colly.NewCollector(
		colly.UserAgent("Sinar Project :P"),
		colly.Async(true),
	)

	c.OnError(func(r *colly.Response, e error) {
		fmt.Println("DIE!!!")
		q.Q(r.Request.URL, r.StatusCode)
		panic(e)
	})

	c.OnHTML("body > table > tbody > tr > td > table:nth-child(2) > tbody > tr:nth-child(3) > td > table > tbody > tr:nth-child(2) > td > table > tbody > tr:nth-child(1)", func(e *colly.HTMLElement) {

		rawRow := strings.Split(e.Text, ":")
		if len(rawRow) > 0 {
			foundAgensi = strings.TrimSpace(rawRow[1])
			//q.Q(foundAgensi)
		}

	})
	// Look for Agensi; missing consider it NOT PageActive

	// If has Agensi; but is the special MBPJ; ignore it! NOT PageActive

	c.OnResponse(func(r *colly.Response) {
		q.Q("Visited", r.Request.URL)
		// DEBUG
		//fmt.Println(r.StatusCode)
		//spew.Dump(r.Body)
	})
	// If regexp "Majlis Bandaraya Petaling Jaya"; stop!!!

	verr := c.Visit(fmt.Sprintf("http://www.epbt.gov.my/osc/%s", pageURL))
	if verr != nil {
		panic(verr)
	}

	c.Wait()
	// Default is NOT active and it will stop!!

	matched, rerr := regexp.MatchString("^Majlis Bandaraya Petaling Jaya.*$", foundAgensi)
	if rerr != nil {
		panic(rerr)
	}
	if matched {
		q.Q("Skipping as ", foundAgensi, " contains MBPJ!")
	} else if foundAgensi == "" {
		q.Q("Skipping as NO Agensi Found!")
	} else {
		// Anything else  is OK!!
		q.Q("Agensi is OK: ", foundAgensi)
		return true
	}
	// Got here; so failed!
	q.Q("is NOT ACTIVE!!!", pageURL)
	return false
}

func returnOldestMBPJTracked() (int, int) {
	var startingID, endingID int

	// Load MBPJ's tracking
	// First one will be the end; last one to be the start
	b, rerr := ioutil.ReadFile("./data/selangor-mbpj-1003/tracking.yml")
	if rerr != nil {
		panic(rerr)
	}

	var appTracking ApplicationTracking
	umerr := yaml.Unmarshal(b, &appTracking)
	if umerr != nil {
		panic(umerr)
	}

	var appID ApplicationID
	var cerr error
	for _, appID = range appTracking.IDs {
		if endingID == 0 {
			endingID, cerr = strconv.Atoi(string(appID))
			if cerr != nil {
				panic(cerr)
			}
			if endingID == 0 {
				panic(fmt.Errorf("FAILED conversion!! ORIG: %s", appID))
			}
			// DEBUG
			//fmt.Println("Ending ID: ", endingID)
		}
	}
	// Leftover at the bottom is the startingID as the list is actually reverse sorted!
	//fmt.Println("LEFTOVER: APPID: ", appID)
	startingID, cerr = strconv.Atoi(string(appID))
	if cerr != nil {
		panic(cerr)
	}
	if startingID == 0 {
		panic(fmt.Errorf("FAILED conversion!! ORIG: %s", appID))
	}
	// DEBUG
	//fmt.Println("Starting ID: ", startingID)
	// For now; we only collect last 100; for testing; otherwise it is pretty long ..
	// Trackibgn is reverse sorted; so; starting is ending - 100!!
	startingID = endingID - 100
	// maybe don;t even  need to convert! int needed for the freeform exploration; lesve it there for now ..
	//return 555555, 555560
	return startingID, endingID
}

// FetchMissing will use the minimum tracking number in authorityToScrape and run till no more valid Application Page is
// found and extracted out
func FetchMissing() {
	fmt.Println("ACTION: FetchMissing")

	volumePrefix := "." // When in CodeFresh, it will be relative .. so that we can have the persistence

	// Take the earliest in the set
	// ignore the ones in authorityToScrape tracking and fetch for the rest; put it in missing category?

	// Store raw here below: <== NO need as we skip and get to the ApplicationRecords directly!
	// 	./raw/<currentdatelabel>/malaysia-missing-0000/

	// Store tracking here below: <== maybe .. maybe NOT
	// ./data/malaysia-missing-0000-<currentdatelabel>/new.yml
	//currentDateLabel := time.Now().Format("20060102") // "20190316"
	//saveData("malaysia-missing-0000", currentDateLabel, nil)

	// Part one, take the authorityToScrape's tracking and find the smallest ID there ..
	applicationIDs := []ApplicationID{}
	// TODO: Above  maybe do as array of APplicationDetails instead; then have nicer reuse; also have the URL ..
	// We take the oldest and iterarte from there; until we have a bunmch of applicationIDs
	startingID, endingID := returnOldestMBPJTracked()
	if startingID == 0 || endingID == 0 {
		panic(fmt.Sprintf("INVALID STATE!! ABORTING!!!"))
	}

	fmt.Println("ANALYZE FROM ", startingID, " TO ", endingID)
	for i := startingID; i <= endingID; i++ {
		//fmt.Println("ID: ", i)
		// DEBUG
		//for i := startingID; i <= startingID+10; i++ {
		var currentAppID, pageURL string
		// Left pad fill to 6 digits to form pageURL?
		currentAppID = fmt.Sprintf("%06d", i)
		// pageURL is relative; to osc; http://www.epbt.gov.my/osc/Proj1_Info.cfm?Name=773399&S=S
		pageURL = fmt.Sprintf("Proj1_Info.cfm?Name=%s&S=S", currentAppID)
		if !isApplicationPageActive(pageURL) {
			// Hit missing page; bail out NOW!!!
			break
		}
		//panic("DEBUG!!!")
		rawApplicationDetailsPath := fmt.Sprintf("%s/raw/malaysia-notmbpj-0000/AR_%s", volumePrefix, currentAppID)
		proceedScraping := rawDataFolderSetup(rawApplicationDetailsPath)

		// DEBUG; for testing can just go ahead ..
		//proceedScraping = true
		// If the data already exist; you can skip it; actually! This will prevent getting new items; possibly put in another function?
		if proceedScraping {
			// Now fetch the page!
			fetchApplicationPage(rawApplicationDetailsPath, pageURL)
			// If evrything is OK, append the page
			applicationIDs = append(applicationIDs, ApplicationID(currentAppID))
		}
	}

	// DEBUG!
	//spew.Dump(applicationIDs)
	//return

	// Next part; we go through the raw files; from our special catchall MALAYSIA!!
	uniqueSearchID := mapAuthorityToDirectory("0000")
	absoluteNewDataPath := fmt.Sprintf("%s/data/%s", volumePrefix, uniqueSearchID)
	rawDataFolderSetup(absoluteNewDataPath)

	// Below is built as we iterate
	for _, appID := range applicationIDs {
		// Look into the folder which we will construct based on the appID
		pages := []string{}
		// This is the relative path to the ApplicationRecord directory
		absoluteRawDataPath := fmt.Sprintf("%s/raw/%s/AR_%s", volumePrefix, uniqueSearchID, appID)
		// This get us the absolute unix path
		dir, err := os.Getwd()
		if err != nil {
			panic(err)
		}
		// DEBUG
		//fmt.Println("PWD: ", dir, " DIR:", absoluteRawDataPath)
		absoluteRawDataPath = fmt.Sprintf("%s/%s", dir, absoluteRawDataPath)
		// Read all the raw HTML files in this directory
		fi, rerr := ioutil.ReadDir(absoluteRawDataPath)
		if rerr != nil {
			panic(rerr)
		}
		for _, f := range fi {
			if !f.IsDir() {
				// We only want non-directory files ..
				path := filepath.Join(absoluteRawDataPath, "/", f.Name())
				// DEBUG
				//fmt.Println("FILE: ", path)
				pages = append(pages, path)
			}
		}

		// Extract the Snapshot data from newest pages
		appDetails := &ApplicationDetails{
			//AR: singleRecord, // No record exist; maybe to extract fully from details itself; good enough?
			AR: ApplicationRecord{
				ID: string(appID),
			},
		}

		extractApplicationDetailsData(appDetails, pages)

		saveApplicationDetails(uniqueSearchID, appDetails)
	}

}

// ExtractAll parses the raw HTML collected under the snapshotLabel
// 	mostly is run once at the  start to kick off the process? Unless overridden
func ExtractAll(authorityToScrape string) {
	fmt.Println("ACTION: ExtractAll")

	// Save into the tracking metadata portion of it ..

	// Step #1: Open from metadata tracking.yml to determine the ApplicationID
	var volumePrefix = "." // When in CodeFresh, it will be relative .. so that we can have the persistence
	var uniqueSearchID = mapAuthorityToDirectory(authorityToScrape)

	appTracking := ApplicationTracking{}
	b, rerr := ioutil.ReadFile(fmt.Sprintf("%s/data/%s/tracking.yml", volumePrefix, uniqueSearchID))
	if rerr != nil {
		panic(rerr)
	}
	err := yaml.Unmarshal(b, &appTracking)
	if err != nil {
		panic(err)
	}
	fmt.Println("LABEL: ", appTracking.Label)
	for _, arID := range appTracking.IDs {
		// Iterate through the raw data .. and append the name to the map
		pages := []string{}
		// This is the relative path to the ApplicationRecord directory
		absoluteRawDataPath := fmt.Sprintf("%s/raw/%s/%s", volumePrefix, uniqueSearchID, fmt.Sprintf("AR_%s", arID))
		// This get us the absolute unix path
		dir, err := os.Getwd()
		if err != nil {
			panic(err)
		}
		// DEBUG
		fmt.Println("PWD: ", dir, " DIR:", absoluteRawDataPath)
		absoluteRawDataPath = fmt.Sprintf("%s/%s", dir, absoluteRawDataPath)
		// Read all the raw HTML files in this directory
		fi, rerr := ioutil.ReadDir(absoluteRawDataPath)
		if rerr != nil {
			panic(rerr)
		}
		for _, f := range fi {
			if !f.IsDir() {
				// We only want non-directory files ..
				path := filepath.Join(absoluteRawDataPath, "/", f.Name())
				// DEBUG
				// fmt.Println("FILE: ", path)
				pages = append(pages, path)
			}
		}
		// Extract the Snapshot data from newest pages
		appDetails := &ApplicationDetails{}
		extractApplicationDetailsData(appDetails, pages)

		// TODO: can perist data now ..
		// saveApplicationDetails(uniqueSearchID, appDetails)

		// TODO: Remove later after tested
		//break
	}

}

// ExtractNew parses the raw HTML files for the new ranges
func ExtractNew(authorityToScrape string) {
	// Metadata structure like ./data/<uniqueSearchID>-<snapshotDiffLabels>
	// e.g. ./data/selangor-mbpj-1003-20190330_20190317
	fmt.Println("ACTION: ExtractNew")

	// Step #1: Open from metadata new.yml to determine the ApplicationID
	var volumePrefix = "." // When in CodeFresh, it will be relative .. so that we can have the persistence
	var uniqueSearchID = mapAuthorityToDirectory(authorityToScrape)

	// TODO: Load the new.yml
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
	for _, singleRecord := range newDiff.AR {
		// Iterate through the raw data .. and append the name to the map
		pages := []string{}
		// This is the relative path to the ApplicationRecord directory
		absoluteRawDataPath := fmt.Sprintf("%s/raw/%s/%s", volumePrefix, uniqueSearchID, fmt.Sprintf("AR_%s", singleRecord.ID))
		// This get us the absolute unix path
		dir, err := os.Getwd()
		if err != nil {
			panic(err)
		}
		// DEBUG
		fmt.Println("PWD: ", dir, " DIR:", absoluteRawDataPath)
		absoluteRawDataPath = fmt.Sprintf("%s/%s", dir, absoluteRawDataPath)
		// Read all the raw HTML files in this directory
		fi, rerr := ioutil.ReadDir(absoluteRawDataPath)
		if rerr != nil {
			panic(rerr)
		}
		for _, f := range fi {
			if !f.IsDir() {
				// We only want non-directory files ..
				path := filepath.Join(absoluteRawDataPath, "/", f.Name())
				// DEBUG
				fmt.Println("FILE: ", path)
				pages = append(pages, path)
			}
		}
		// Extract the Snapshot data from newest pages
		appDetails := &ApplicationDetails{
			AR: singleRecord,
		}
		extractApplicationDetailsData(appDetails, pages)

		saveApplicationDetails(uniqueSearchID, appDetails)

		// DEBUG for single dataset ..
		// break
	}

}

func saveApplicationDetails(uniqueSearchID string, ad *ApplicationDetails) {
	fmt.Println("oscv3-fetch: saveApplicationDetails")
	// In ./data/<uniqueSearchID>/AR_<applicationID>/details.yml
	// DEBUG
	// spew.Dump(ad)
	appDetails := []ApplicationDetails{*ad}
	if len(appDetails) == 0 {
		fmt.Println("NOTHING to DO .. skipping!!")
		return
	}
	// Assume gets this far; just persist it!!
	// Get those bytes out
	b, err := yaml.Marshal(appDetails)
	if err != nil {
		panic(err)
	}

	// DEBUG
	// spew.Println(string(b))

	// Open file and persist it into the format
	// Metadata structure like ./data/<uniqueSearchID>/AR_<appID>/summary.yml
	var absoluteNewDataPath = fmt.Sprintf("./data/%s/AR_%s", uniqueSearchID, ad.AR.ID)
	rawDataFolderSetup(absoluteNewDataPath)
	nerr := ioutil.WriteFile(fmt.Sprintf("%s/details.yml", absoluteNewDataPath), b, 0744)
	if nerr != nil {
		panic(nerr)
	}

}

func saveApplicationRecordSummary(uniqueSearchID string, ar *ApplicationRecord) {
	// In ./data/<uniqueSearchID>/<applicationID>/summary.yml
	appRecords := []ApplicationRecord{*ar}
	//IN yaml format
	if len(appRecords) == 0 {
		// Nothing to be done ..
		fmt.Println("NOTHING to DO .. skipping!!")
		return
	}

	// Assume gets this far; just persist it!!
	// Get those bytes out
	b, err := yaml.Marshal(appRecords)
	if err != nil {
		panic(err)
	}

	// DEBUG
	// spew.Println(string(b))

	// Open file and persist it into the format
	// Metadata structure like ./data/<uniqueSearchID>/AR_<appID>/summary.yml
	var absoluteNewDataPath = fmt.Sprintf("./data/%s/AR_%s", uniqueSearchID, ar.ID)
	rawDataFolderSetup(absoluteNewDataPath)
	nerr := ioutil.WriteFile(fmt.Sprintf("%s/summary.yml", absoluteNewDataPath), b, 0744)
	if nerr != nil {
		panic(nerr)
	}
}
