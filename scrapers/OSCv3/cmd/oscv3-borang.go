package cmd

import (
	"fmt"
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

// FormNum ==> Different Type of Forms; same across Local Authorities?
// NoForm=Form1 - Borang Rekod Permohonan Serentak
// NoForm=Form2 - Borang Rekod Permohonan Kebenaran Merancang (KM)
// NoForm=Form3 - Borang Rekod Pelan Kejuruteraan (PJ)
// NoForm=Form4 - Borang Rekod Pelan Bangunan (PB)
// NoForm=Form5 - Borang Rekod Permohonan dan Pengeluaran CFO
// NoForm=Form6 - Borang Rekod Sistem Pengeluaran Perakuan Siap Dan Pematuhan

// FormDetails  will have the time record as approvals are given and coming in?
type FormDetails struct {
	ID      string
	FormNum string // see above for detailed definitions
}

// Have a command to sweep and  update from data fields all "new" and ALL data? concurrent ..
// This will update the trackoing field to be  consistent?

func extractFormDetailsData(formDetails *FormDetails, pagesToExtract []string) {
	fmt.Println("START ==> extractFormDetailsData =================")

}

// Default use the  tracking to pull in the new items ..

func fetchFormPage(absoluteRawDataPath string, pageURL string) {

}

//  Pull the data from details ..

func loadApplicationApprovalForms(uniqueSearchID string) []FormDetails {
	var borangs []FormDetails

	// Load Application Details from file; scenario sinlge form
	ad := ApplicationDetails{FormRecords: []FormRecord{{URL: "Borang_info.cfm?ID=260530&NoForm=Form3"}}}
	// TOOD: Scenario for multiple forms ..
	for _, form := range ad.FormRecords {
		fmt.Println("Fetch:", form.URL)
		// Split up the ID and FormNum
	}
	return borangs
}

// loadApplicationDetailsFromFile will load Active Application Details per authority mapping
func loadApplicationDetailsFromFile(authorityToScrape string) {
	// Metadata structure like ./data/<uniqueSearchID>-<snapshotDiffLabels>
	// e.g. ./data/selangor-mbpj-1003-20190330_20190317
	fmt.Println("ACTION: loadApplicationDetailsFromFile")
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
	for _, id := range appTracking.IDs {
		absoluteRawDataPath := fmt.Sprintf("%s/raw/%s/AR_%s", volumePrefix, uniqueSearchID, id)
		rawDataFolderSetup(absoluteRawDataPath)
		//fetchApplicationPage(absoluteRawDataPath, ar.URL)
		// ftech and see if got new approval date
		// Those ith approvsl date can be ignored and marked for removal
		// Those still there waiting; put back in the list ..
		// ALTERNATIVELY: We can assume  it is ok; the state-sync in another function
	}

	// Once loaded; can update tracking file about the newest state ..

}

// syncTracking will pull latest raw data and  check their Application date/status?
// Try to use existing data to use a simple  regexp instead? possible?
func syncTracking(authorityToScrape string) {
	// Load up ApplicationTracking metedata tate for use in checking ..
	// Assumes the URL is S=S; until we can prove otherwise ..
	// Can refactor previous function ..
}

// ExtractFormNew (fetches if needed); and parses the raw HTML files for Form Details Info
func ExtractFormNew(authorityToScrape string) {

}
