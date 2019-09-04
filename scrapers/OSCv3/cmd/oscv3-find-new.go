package cmd

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/gocolly/colly"
	"github.com/y0ssar1an/q"

	"gopkg.in/yaml.v2"
)

// ApplicationRecord has details
type ApplicationRecord struct {
	// bil
	// Nama Projek
	// No. Lot
	// Mukim
	// Link --> URL
	bil    int
	ID     string
	Projek string
	Lot    string
	Mukim  string
	URL    string
}

// ApplicationRecords will be used to sort by bil field
// which will be the smallest number first on the top
type ApplicationRecords []ApplicationRecord

func (ar ApplicationRecords) Len() int {
	return len(ar)
}

func (ar ApplicationRecords) Less(i, j int) bool {
	return ar[i].bil < ar[j].bil
}

func (ar ApplicationRecords) Swap(i, j int) {
	ar[i], ar[j] = ar[j], ar[i]
}

// ApplicationSnapshot shows history
type ApplicationSnapshot struct {
	snapshotLabel string
	appRecords    ApplicationRecords
}

// NewDiff strcuture defined ..
type NewDiff struct {
	Label string
	AR    []ApplicationRecord `yaml:"new"`
}

// ApplicationID is the primary lookup key for Applications
type ApplicationID string

// ApplicationTracking is to keep a look up on which items to be checked for any refresher
type ApplicationTracking struct {
	Label string
	// ID - Application ID; used to look up
	// Form URLs --> any Borang related to this Appllication; zero or more ..
	IDs []ApplicationID `yaml:"tracking"`
}

func extractAllApplicationData(appSnapshot *ApplicationSnapshot, pagesToExtract []string) {
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

	idOrder := make([]string, 0, 100)
	// Gather all the Records ..
	appRecords := make([]ApplicationRecord, 0)

	c.OnHTML("html body table tbody tr td table tbody tr td table tbody tr td table tbody tr", func(e *colly.HTMLElement) {
		// Every row of data ..
		e.DOM.Each(func(i int, s *goquery.Selection) {
			appRecord := ApplicationRecord{}
			//q.Q("ID:", i, " DATA:", s.Text())
			var isValid bool
			// Each column
			s.Children().Each(func(j int, c *goquery.Selection) {
				// To track if need to ignore bad rows (e.g. like header ..)
				isValid = true
				// bil
				// Nama Projek
				// No. Lot
				// Mukim
				// Link
				if j == 0 {
					// q.Q("bil: ", c.Text())
					bil, err := strconv.Atoi(c.Text())
					if err != nil {
						// DEBUG
						// fmt.Println("ERR:", err)
						isValid = false
					}
					appRecord.bil = bil
				} else if j == 1 {
					// q.Q("Projek: ", c.Text())
					appRecord.Projek = strings.TrimSpace(c.Text())
				} else if j == 2 {
					// q.Q("Lot: ", c.Text())
					appRecord.Lot = strings.TrimSpace(c.Text())
				} else if j == 3 {
					// q.Q("Mukim: ", c.Text())
					appRecord.Mukim = strings.TrimSpace(c.Text())
				} else if j == 4 {
					// Name is Unique Identifier
					id := c.Find("a").Map(func(_ int, m *goquery.Selection) string {
						href, _ := m.Attr("href")
						idURL, err := url.Parse(href)
						if err != nil {
							panic(err)
						}
						return idURL.Query().Get("Name")
					})
					// TODO: What if has bad data; should match regexp at least?
					// Or use atoi again? if cannot convert; ignore??
					if len(id) > 0 && isValid {
						idOrder = append(idOrder, strings.Join(id, ""))
						// DEBUG
						// q.Q("ID: ", id)
						// What is S?
						q.Q("TYPE: ", c.Find("a").Map(func(_ int, m *goquery.Selection) string {
							href, _ := m.Attr("href")
							idURL, err := url.Parse(href)
							if err != nil {
								panic(err)
							}
							appRecord.URL = href
							//DEBUG
							// fmt.Println(appRecord.URL)
							return idURL.Query().Get("S")
						}))

						appRecord.ID = strings.Join(id, "")
						appRecords = append(appRecords, appRecord)
					}

				} else {
					q.Q("UNKNOWN:", c)
				}
			})
		})

		// e.ForEachWithBreak("td", func(_ int, row *colly.HTMLElement) bool {
		// 	spew.Dump(row.ChildText("*"))
		// 	// return false
		// 	return true
		// })

	})

	c.OnScraped(func(r *colly.Response) {
		// Next time; queue the further extraction of item?
		// spew.Println(r.StatusCode)
		// DEBUG
		// fmt.Println("DONE!  CODE:", r.StatusCode)
	})

	// Example finalURL will be like below:
	// finalURL := "file:///Users/leow/GOMOD/odd2019/scrapers/OSCv3/raw/20190322/selangor-mbpj-1003/_osc_Carian_Proj3.cfm_CurrentPage_10_Maxrows_15_Cari_AgensiKod_1003_Pilih_3.html"
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
	appSnapshot.appRecords = appRecords
}

func extractDataFromSnapshot(volumePrefix string, snapshotLabel string, uniqueSearchID string) *ApplicationSnapshot {
	var absoluteRawDataPath = fmt.Sprintf("%s/raw/%s/%s", volumePrefix, snapshotLabel, uniqueSearchID)
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
	// Extract the Snapshot data from newest pages
	appSnapshot := &ApplicationSnapshot{
		snapshotLabel: snapshotLabel,
	}

	extractAllApplicationData(appSnapshot, pages)

	// Persist snapshot into YAML?
	// TODO: OUtput as yaml??
	// spew.Dump(appSnapshot)
	sort.Sort(appSnapshot.appRecords)
	// DEBUG
	// for _, singleRecord := range appSnapshot.appRecords {
	// 	fmt.Printf("%s,", singleRecord.ID)
	// }

	return appSnapshot
}

// FindAllApplications will setup the tracking file for this
//	if overridden with options of specific label or just newest
//	process accordingly ..
//  Example: ./data/selangor-mbpj-1003/tracking.yml
func FindAllApplications(authorityToScrape string, forceRefresh bool, specificLabel string) {
	fmt.Println("ACTION: FindAllApplications")

	// if forceRefresh; empty the list
	// else read from the existing structure and append it?
	// build up the list

	// Open up the raw data specified by Label
	var volumePrefix = "." // When in CodeFresh, it will be relative .. so that we can have the persistence
	var uniqueSearchID = mapAuthorityToDirectory(authorityToScrape)

	// Get the snapshot from memory..
	specifiedSnapshot := extractDataFromSnapshot(volumePrefix, specificLabel, uniqueSearchID)

	aid := make([]ApplicationID, 0)
	for _, singleRecord := range specifiedSnapshot.appRecords {
		// DEBUG
		// fmt.Println("Store ID: ", singleRecord.ID)
		aid = append(aid, ApplicationID(singleRecord.ID))
		// Persist the data being tracked so we can grab them later ..
		saveApplicationRecordSummary(uniqueSearchID, &singleRecord)
	}
	// Nothing to be done
	if len(aid) == 0 {
		fmt.Println("Nothing to be done!")
		return
	}

	b, err := yaml.Marshal(ApplicationTracking{
		Label: uniqueSearchID,
		IDs:   aid,
	})
	// Extract into struct --> ApplicationTracking

	// persist the data into the data folder?
	if err != nil {
		panic(err)
	}

	// DEBUG
	// spew.Println(string(b))

	// Open file and persist it into the format
	// Metadata structure like ./data/<uniqueSearchID>/tracking.yml
	// e.g. ./data/selangor-mbpj-1003/tracking.yml
	var absoluteNewDataPath = fmt.Sprintf("./data/%s", uniqueSearchID)
	rawDataFolderSetup(absoluteNewDataPath)
	nerr := ioutil.WriteFile(fmt.Sprintf("%s/tracking.yml", absoluteNewDataPath), b, 0744)
	if nerr != nil {
		panic(nerr)
	}

}

// FindNewRequests will look for the changes since the last time run and offer a pull request
func FindNewRequests(authorityToScrape string) {
	fmt.Println("ACTION: FindNewRequests")
	// Figure out when was the last successful run and if not exist; create it!
	// Also will reset if passed a flag of some sort?
	// var currentDateLabel = time.Now().Format("20060102") // "20190316"
	var volumePrefix = "." // When in CodeFresh, it will be relative .. so that we can have the persistence
	var uniqueSearchID = mapAuthorityToDirectory(authorityToScrape)

	// Refactor  out the currentDate
	var currentDateLabel = "20190904"
	currentSnapshot := extractDataFromSnapshot(volumePrefix, currentDateLabel, uniqueSearchID)
	// If in Codefresh; do a branch, git add + commit?
	// Refactor out the previousDate
	var previousDateLabel = "20190825"
	fmt.Println("Now compare against the previous: ", previousDateLabel)
	previousSnapshot := extractDataFromSnapshot(volumePrefix, previousDateLabel, uniqueSearchID)

	// Now iterate and compare until the first match
	// then check if the diffs are there ..
	var newAppRecords []ApplicationRecord
	var foundOldID bool
	var previousIDIndex int
	firstOldID := previousSnapshot.appRecords[previousIDIndex].ID
	for _, singleRecord := range currentSnapshot.appRecords {
		if foundOldID {
			previousIDIndex++
			if singleRecord.ID != previousSnapshot.appRecords[previousIDIndex].ID {
				// If this part happens; means the data can be considered corrupted; and should be re-run again!!
				fmt.Println("BIL:", singleRecord.bil, " ERR: ID: ", singleRecord.ID, " NOT matching ", previousSnapshot.appRecords[previousIDIndex].ID)
			}
			// DEBUG
			// else {
			// 	fmt.Println("BIL:", singleRecord.bil, " OK: ID: ", singleRecord.ID, " matches ", previousSnapshot.appRecords[previousIDIndex].ID)
			// }
		} else {
			// Show visually which is it ..
			fmt.Printf("%s,", singleRecord.ID)
			if singleRecord.ID == firstOldID {
				foundOldID = true
				fmt.Println("FOUND IT!! ID: ", firstOldID)
			} else {
				// What to do with the new entries?? save it for further use later ..
				// spew.Dump(singleRecord)
				newAppRecords = append(newAppRecords, singleRecord)
				// Also save a copy for summary in later use
				saveApplicationRecordSummary(uniqueSearchID, &singleRecord)
			}
		}
	}
	// Calculate absoluteRawDataPath
	// Persist it; along with Github?
	snapshotDiffLabels := fmt.Sprintf("%s_%s", currentSnapshot.snapshotLabel, previousSnapshot.snapshotLabel)
	saveData(uniqueSearchID, snapshotDiffLabels, newAppRecords)
}

func saveData(uniqueSearchID string, snapshotDiffLabels string, newAppRecords []ApplicationRecord) {

	//IN yaml format
	if len(newAppRecords) == 0 {
		// Nothing to be done ..
		fmt.Println("NOTHING to DO .. skipping!!")
		return
	}

	// Assume gets this far; just persist it!!
	// Get those bytes out
	b, err := yaml.Marshal(NewDiff{
		Label: uniqueSearchID,
		AR:    newAppRecords,
	})
	if err != nil {
		panic(err)
	}

	// If detect env CF_REPO_NAME; we are in Codefresh and data is meant to be stored there?
	// If in debugging mode; save in $TMPDIR?
	// else use the data folder? or should it be raw?
	// data/<uniqueSearchID>/new.yml <-- new data; including the details

	// DEBUG
	// spew.Println(string(b))

	// Open file and persist it into the format
	// Metadata structure like ./data/<uniqueSearchID>-<snapshotDiffLabels>/new.yml
	// e.g. ./data/selangor-mbpj-1003-20190330_20190317/new.yml
	var absoluteNewDataPath = fmt.Sprintf("./data/%s", uniqueSearchID)
	rawDataFolderSetup(absoluteNewDataPath)
	nerr := ioutil.WriteFile(fmt.Sprintf("%s/new.yml", absoluteNewDataPath), b, 0744)
	if nerr != nil {
		panic(nerr)
	}

	// but ALSO keep a copy at .. for easy access
	// ./data/selangor-mbpj-1003/new.yml
	var absoluteRawDataPath = fmt.Sprintf("./data/%s-%s", uniqueSearchID, snapshotDiffLabels)
	// Create parent data for metadata
	rawDataFolderSetup(absoluteRawDataPath)
	werr := ioutil.WriteFile(fmt.Sprintf("%s/new.yml", absoluteRawDataPath), b, 0744)
	if werr != nil {
		panic(werr)
	}

	// data/<uniqueSearchID>/metadata.yml <-- ONLY ApplicationID those open/active?
	// data/<uniqueSearchID>/snapshot.yml <-- current snapshot of data ..
	// data/<uniqueSearchID>/<ApplicationID>/..
	// in debugging mode; no github action?

}
