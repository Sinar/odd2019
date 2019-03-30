package cmd

import "fmt"

// ApplicationID is the primary lookup key for Applications
type ApplicationID string

// ApplicationTracking is to keep a look up on which items to be checked for any refresher
type ApplicationTracking struct {
	// ID - Application ID; used to look up
	// Form URLs --> any Borang related to this Appllication; zero or more ..
	IDs []ApplicationID
}

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

func fetchApplicationPage() {

}

// FetchAll will Extract from authority + label; all 15 pages of the information
func FetchAll(authorityToScrape string, snapshotLabel string) {
	// Raw structure like .. ./raw/<snapshotLabel>/<uniqueSearchID>
	// NOTE: Descructive action will override the data; ensure it is git diff ..
	fmt.Println("ACTION: FetchAll")
	// Figure out when was the last successful run and if not exist; create it!
	// Also will reset if passed a flag of some sort?
	// var currentDateLabel = time.Now().Format("20060102") // "20190316"
	// var volumePrefix = "." // When in CodeFresh, it will be relative .. so that we can have the persistence
	// var uniqueSearchID = mapAuthorityToDirectory(authorityToScrape)

	// Step #1: Extract into marshal structure
	// Store into metadata structure like ./data/<uniqueSearchID>/
	// e.g. ./data/selangor-mbpj-1003/tracking.yaml; append only new unique items;
	//	sorted by ApplicationID
	// marked the successful / completed into archive? <-- Done in another step
}

func getLatestComparison(uniqueSearchID string) (absolutePathToComparison string) {

	absolutePathToComparison = "./data/"

	return absolutePathToComparison
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

}

// ExtractAll parses the raw HTML collected under the snapshotLabel
// 	mostly is run once at the  start to kick off the process? Unless overridden
func ExtractAll(authorityToScrape string, snapshotLabel string) {

	// Save into the tracking metadata portion of it ..
}

// ExtractNew parses the raw HTML files for the new ranges
func ExtractNew(authorityToScrape string) {
	// Metadata structure like ./data/<uniqueSearchID>-<snapshotDiffLabels>
	// e.g. ./data/selangor-mbpj-1003-20190330_20190317
	fmt.Println("ACTION: ExtractNew")

}

// extract out the fields related to Application; which is what?
func extractAllApplicationDetails() {

}

func saveApplicationDetails(applicationID string) {
	// Marshal into the strcuture
	// Including the list of Borangs relatd to it
}
