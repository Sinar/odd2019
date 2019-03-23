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
	"github.com/davecgh/go-spew/spew"
	"github.com/gocolly/colly"
	"github.com/y0ssar1an/q"
)

// ApplicationRecord has details
type ApplicationRecord struct {
	// Bil
	// Nama Projek
	// No. Lot
	// Mukim
	// Link --> url
	bil    int
	id     string
	projek string
	lot    string
	mukim  string
	url    string
}

// ApplicationRecords will be used to sort by BIL field
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

func extractAllData(appSnapshot *ApplicationSnapshot, pagesToExtract []string) {
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
				// Bil
				// Nama Projek
				// No. Lot
				// Mukim
				// Link
				if j == 0 {
					// q.Q("BIL: ", c.Text())
					bil, err := strconv.Atoi(c.Text())
					if err != nil {
						// DEBUG
						// fmt.Println("ERR:", err)
						isValid = false
					}
					appRecord.bil = bil
				} else if j == 1 {
					// q.Q("PROJEK: ", c.Text())
					appRecord.projek = c.Text()
				} else if j == 2 {
					// q.Q("LOT: ", c.Text())
					appRecord.lot = c.Text()
				} else if j == 3 {
					// q.Q("MUKIM: ", c.Text())
					appRecord.mukim = c.Text()
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
							appRecord.url = href
							//DEBUG
							// fmt.Println(appRecord.url)
							return idURL.Query().Get("S")
						}))

						appRecord.id = strings.Join(id, "")
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
	for _, url := range pagesToExtract {
		finalURL := fmt.Sprintf("file://%s", url)
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
	// Sort the final based on the BIL as Int
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

	extractAllData(appSnapshot, pages)

	// Persist snapshot into YAML?
	// TODO: OUtput as yaml??
	// spew.Dump(appSnapshot)
	sort.Sort(appSnapshot.appRecords)
	// DEBUG
	// for _, singleRecord := range appSnapshot.appRecords {
	// 	fmt.Printf("%s,", singleRecord.id)
	// }

	return appSnapshot
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
	var currentDateLabel = "20190322"
	currentSnapshot := extractDataFromSnapshot(volumePrefix, currentDateLabel, uniqueSearchID)
	// If in Codefresh; do a branch, git add + commit?
	// Refactor out the previousDate
	var previousDateLabel = "20190317"
	fmt.Println("Now compare against the previous: ", previousDateLabel)
	previousSnapshot := extractDataFromSnapshot(volumePrefix, previousDateLabel, uniqueSearchID)

	// Now iterate and compare until the first match
	// then check if the diffs are there ..
	var foundOldID bool
	var previousIDIndex int
	firstOldID := previousSnapshot.appRecords[previousIDIndex].id
	for _, singleRecord := range currentSnapshot.appRecords {
		if foundOldID {
			previousIDIndex++
			if singleRecord.id != previousSnapshot.appRecords[previousIDIndex].id {
				fmt.Println("ERR: ID: ", singleRecord.id, " NOT matching ", previousSnapshot.appRecords[previousIDIndex].id)
			}
		} else {
			// Show visually which is it ..
			fmt.Printf("%s,", singleRecord.id)
			if singleRecord.id == firstOldID {
				foundOldID = true
				fmt.Println("FOUND IT!! ID: ", firstOldID)
			} else {
				// What to do with the new entries??
				spew.Dump(singleRecord)
			}
		}
	}
}
