package cmd

import (
	"fmt"
	"io/ioutil"
	"net/url"
	"os"
	"strings"

	"github.com/gocolly/colly"
	"github.com/y0ssar1an/q"

	"gopkg.in/yaml.v2"

	h2t "github.com/jaytaylor/html2text"
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
	ID                 string
	FormNum            string // see above for detailed definitions
	Status             string // body > table > tbody > tr > td > table:nth-child(2) > tbody > tr:nth-child(3) > td > table > tbody > tr:nth-child(2) > td > table > tbody > tr:nth-child(1)
	JenisPemohonan     string // body > table > tbody > tr > td > table:nth-child(2) > tbody > tr:nth-child(3) > td > table > tbody > tr:nth-child(2) > td > table > tbody > tr:nth-child(8)
	TarikhPermohonan   string // body > table > tbody > tr > td > table:nth-child(2) > tbody > tr:nth-child(3) > td > table > tbody > tr:nth-child(2) > td > table > tbody > tr:nth-child(9)
	StatusTerkiniAT    string // Raw HTML for later processing -  body > table > tbody > tr > td > table:nth-child(2) > tbody > tr:nth-child(3) > td > table > tbody > tr:nth-child(2) > td > table > tbody > tr:nth-child(10)
	TarikhKeputusanOSC string // Committee approval - body > table > tbody > tr > td > table:nth-child(2) > tbody > tr:nth-child(3) > td > table > tbody > tr:nth-child(2) > td > table > tbody > tr:nth-child(11)
}

// Have a command to sweep and  update from data fields all "new" and ALL data? concurrent ..
// This will update the trackoing field to be  consistent?

//func extractFormDetailsData(formDetails *FormDetails, pagesToExtract []string) {
//	fmt.Println("START ==> extractFormDetailsData =================")
//
//}

// Default use the  tracking to pull in the new items ..

func fetchFormPage(fd *FormDetails, pageURL string) error {
	// URL is partial; add on the needed full hostname?
	pageURL = fmt.Sprintf("http://www.epbt.gov.my/osc/%s", pageURL)
	// Extra checks will make sur egot not http/https??
	// Setup the queue that will be to grab at the available pages up till the 15 pages limit
	//queue, _ := queue.New(
	//	2, // Number of consumer threads
	//	&queue.InMemoryQueueStorage{MaxSize: 10000}, // Use default queue storage
	//)
	// Add the seedURL to the queue to be saved; not needed as we will add below
	// queue.AddURL(seedURL)

	// With pre-reqs setup; we can proceed ...
	c := colly.NewCollector(
		colly.UserAgent("Sinar Project :P"),
	)

	c.OnError(func(r *colly.Response, e error) {
		fmt.Println("DIE!!!")
		q.Q(r.Request.URL, r.StatusCode)
		panic(e)
	})

	c.OnHTML("body > table > tbody > tr > td > table:nth-child(2) > tbody > tr:nth-child(3) > td > table > tbody > tr:nth-child(2) > td > table > tbody > tr:nth-child(1)", func(e *colly.HTMLElement) {
		//q.Q(strings.Split(e.Text, ":"))
		rowcontent := strings.Split(e.Text, ":")
		if len(rowcontent) > 1 {
			fd.Status = strings.TrimSpace(rowcontent[1])
		}
	})

	c.OnHTML("body > table > tbody > tr > td > table:nth-child(2) > tbody > tr:nth-child(3) > td > table > tbody > tr:nth-child(2) > td > table > tbody > tr:nth-child(8)", func(e *colly.HTMLElement) {
		rowcontent := strings.Split(e.Text, ":")
		if len(rowcontent) > 1 {
			fd.JenisPemohonan = strings.TrimSpace(rowcontent[1])
		}
	})

	c.OnHTML("body > table > tbody > tr > td > table:nth-child(2) > tbody > tr:nth-child(3) > td > table > tbody > tr:nth-child(2) > td > table > tbody > tr:nth-child(9)", func(e *colly.HTMLElement) {
		rowcontent := strings.Split(e.Text, ":")
		if len(rowcontent) > 1 {
			fd.TarikhPermohonan = strings.TrimSpace(rowcontent[1])
		}
	})

	c.OnHTML("body > table > tbody > tr > td > table:nth-child(2) > tbody > tr:nth-child(3) > td > table > tbody > tr:nth-child(2) > td > table > tbody > tr:nth-child(10)", func(e *colly.HTMLElement) {
		//q.Q(strings.Split(e.Text, ":"))
		statusHTML, _ := e.DOM.Html()
		fd.StatusTerkiniAT = strings.TrimSpace(statusHTML)
		//fd.StatusTerkiniAT = strings.TrimSpace(strings.Split(e.Text, ":")[1])
	})

	c.OnHTML("body > table > tbody > tr > td > table:nth-child(2) > tbody > tr:nth-child(3) > td > table > tbody > tr:nth-child(2) > td > table > tbody > tr:nth-child(11)", func(e *colly.HTMLElement) {
		//q.Q(e.DOM.Html())
		rowcontent := strings.Split(e.Text, ":")
		if len(rowcontent) > 1 {
			fd.TarikhKeputusanOSC = strings.TrimSpace(rowcontent[1])
		}
	})

	c.OnScraped(func(r *colly.Response) {
		// // Define the result page Collector; which will just save the file
		// d := colly.NewCollector()
		// d.OnScraped(func(r *colly.Response) {
		fmt.Println("FINISH: ", r.Request.URL, "<================")
		//r.Save(fmt.Sprintf("%s/%s.html", absoluteRawDataPath, r.FileName()))
		//q.Q("FILE: ", r.FileName())
		//q.Q("SAVED ==================>")
		// 	//fmt.Println(r.Headers)
		// })
		// // Kick off the queue once all the pages are collected already ..
		// queue.Run(d)
	})

	// Kick off the single page scraping; toally ugly!!
	verr := c.Visit(pageURL)
	if verr != nil {
		return verr
	}

	return nil
}

//  Pull the data from details ..
// Assumes  syncTracking ran previously; so only have the new ones? combined?
func loadApplicationApprovalForms(uniqueSearchID string, ad ApplicationDetails) []FormDetails {
	fmt.Println("ACTION: BORANG - FETCH + EXTRACT")
	var borangs []FormDetails

	// Below id the complete one ..
	//ad := ApplicationDetails{FormRecords: []FormRecord{{URL: "Borang_info.cfm?ID=377290&NoForm=Form2"}}}
	// TOOD: Scenario for multiple forms ..
	for _, form := range ad.FormRecords {
		// Split up the ID and FormNum
		idURL, err := url.Parse(form.URL)
		if err != nil {
			panic(err)
		}
		// Build  initial struct for FormDetails
		formDetail := FormDetails{
			ID:      idURL.Query().Get("ID"),
			FormNum: idURL.Query().Get("NoForm"),
		}

		// Guard rail; can skip if the data file already exists?
		// But we want to check the newest so maybe not ...
		//  For now; skip if it exists already ..
		// Pattern is:
		// Metadata structure like ./data/<uniqueSearchID>/AR_<appID>/FR_<formID>_<formNUm>/details.yml
		//var absoluteNewDataPath = fmt.Sprintf("./data/%s/AR_%s/FR_%s_%s", uniqueSearchID, arID, fd.ID, fd.FormNum)
		formDetailsPath := fmt.Sprintf("./data/%s/AR_%s/FR_%s_%s", uniqueSearchID, ad.AR.ID, formDetail.ID, formDetail.FormNum)
		if fileExists(formDetailsPath + "/details.yml") {
			fmt.Println("Data already exist in ", formDetailsPath, " skipping for now ...")
			continue
		}

		fmt.Println("Fetch:", form.URL)
		// ferch and fill in structure; what if got error?
		ffperr := fetchFormPage(&formDetail, form.URL)
		if ffperr != nil {
			panic(ffperr)
		}
		// Append it for use later ..
		borangs = append(borangs, formDetail)
	}
	// DEBUG:
	//q.Q(borangs)
	return borangs
}

// fileExists checks if a file exists and is not a directory before we
// try using it to prevent further errors.
func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

// loadApplicationDetailsFromFile will load Active Application Details per authority mapping
// TODO: Maybe return just ID, form +err?
func loadApplicationDetailsFromFile(uniqueSearchID string) []ApplicationDetails {
	// Metadata structure like ./data/<uniqueSearchID>-<snapshotDiffLabels>
	// e.g. ./data/selangor-mbpj-1003-20190330_20190317
	fmt.Println("ACTION: loadApplicationDetailsFromFile")

	trackedApplicationDetails := []ApplicationDetails{}

	// Step #1: Open from metadata tracking.yml to determine the TrackedApplicationDetails
	var volumePrefix = "." // When in CodeFresh, it will be relative .. so that we can have the persistence

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
		var applicationDetails []ApplicationDetails
		//applicationDetails := []ApplicationDetails{}

		absoluteRawDataPath := fmt.Sprintf("%s/data/%s/AR_%s", volumePrefix, uniqueSearchID, id)
		rawDataFolderSetup(absoluteRawDataPath)
		appDetailsDataPath := absoluteRawDataPath + "/details.yml"
		if !fileExists(appDetailsDataPath) {
			// Guard rail against missing data
			fmt.Println("Missing file: ", appDetailsDataPath, " ; skipping ...")
			continue
		}

		// Read and unmarshal out the data and extract out the FormRecord
		b, rerr := ioutil.ReadFile(appDetailsDataPath)
		if rerr != nil {
			panic(rerr)
		}

		umerr := yaml.Unmarshal(b, &applicationDetails)
		if umerr != nil {
			panic(umerr)
		}
		// Append it out to eb used later
		trackedApplicationDetails = append(trackedApplicationDetails, applicationDetails[0])
	}

	// Once loaded; can update tracking file about the newest state ..
	return trackedApplicationDetails
}

func extractNewFormsDetails(authorityToScrape string) {
	fmt.Println("Inside extractNewFormsDetails ..")
	var uniqueSearchID = mapAuthorityToDirectory(authorityToScrape)

	// Load Application Details from file; scenario sinlge form
	// Test cse below for no date
	//ad := ApplicationDetails{FormRecords: []FormRecord{{URL: "Borang_info.cfm?ID=260530&NoForm=Form3"}}}

	trackedApplicationDetails := loadApplicationDetailsFromFile(uniqueSearchID)

	for _, ad := range trackedApplicationDetails {
		formDetails := loadApplicationApprovalForms(uniqueSearchID, ad)
		if len(formDetails) == 0 {
			fmt.Println("*** NOTE: *** NO FORMS!! Check  it out?")
		}
		// TODO: What to do when no forms; is it unusual? I don;t think so ..
		for _, fd := range formDetails {
			// DEBUG
			//fmt.Println("SAVE ARID: ", ad.AR.ID, " with FORM: ", fd.ID, " TYPE: ", fd.FormNum)
			saveFormDetails(uniqueSearchID, ApplicationID(ad.AR.ID), fd)
		}
		// TODO: Remove after Testing purpose; run one round!
		//break
	}

}

func saveFormDetails(uniqueSearchID string, arID ApplicationID, formDetails FormDetails) {
	fmt.Println("oscv3-borang: saveFormDetails")
	// In ./data/<uniqueSearchID>/AR_<applicationID>/details.yml
	// DEBUG
	// spew.Dump(ad)
	//formDetails := []FormDetails{fd}
	//if len(formDetails) == 0 {
	//	fmt.Println("NOTHING to DO .. skipping!!")
	//	return
	//}
	// Assume gets this far; just persist it!!
	// Get those bytes out
	b, err := yaml.Marshal(formDetails)
	if err != nil {
		panic(err)
	}

	// DEBUG
	// spew.Println(string(b))

	// Open file and persist it into the format
	// Metadata structure like ./data/<uniqueSearchID>/AR_<appID>/FR_<formID>_<formNUm>/details.yml
	var absoluteNewDataPath = fmt.Sprintf("./data/%s/AR_%s/FR_%s_%s", uniqueSearchID, arID, formDetails.ID, formDetails.FormNum)
	rawDataFolderSetup(absoluteNewDataPath)
	q.Q("Persisting to ", absoluteNewDataPath, "/details.yml")
	nerr := ioutil.WriteFile(fmt.Sprintf("%s/details.yml", absoluteNewDataPath), b, 0744)
	if nerr != nil {
		panic(nerr)
	}

}

// ExtractFormNew (fetches if needed); and parses the raw HTML files for Form Details Info
func ExtractFormNew(authorityToScrape string) {
	// Try it out first cut!
	extractNewFormsDetails(authorityToScrape)
}

func convertFormHistoryHTML(formHistoryHTML string) {
	fmt.Println("Inside convertFormHistoryHTML ..")
	//  use https://github.com/jaytaylor/html2text ...
	prettyTable, cerr := h2t.FromString(formHistoryHTML, h2t.Options{PrettyTables: true})
	if cerr != nil {
		panic(cerr)
	}
	fmt.Println("======== IAT Approval Status ===========")
	fmt.Println(prettyTable)
}

// DisplayFormDetails will render out the Form life history
// 	use a few HTML to text libraries ..
func DisplayFormDetails(authorityToScrape string) {
	// Barrier check ; data file MUST exists!
	// then you can unmarshal it; no problem ..
	fmt.Println("Inside DisplayFormDetails ..")
	uniqueSearchID := mapAuthorityToDirectory(authorityToScrape)
	// Test case: scrapers/OSCv3/data/selangor-mbpj-1003/AR_776053/FR_528126_Form4/details.yml
	// Case #1
	//arID := "776053"
	//formID := "528126"
	//formNum := "Form4"
	// Case #2 - scrapers/OSCv3/data/selangor-mbpj-1003/AR_776183/FR_528195_Form4/details.yml
	//arID := "776183"
	//formID := "528195"
	//formNum := "Form4"
	// Case #3 - scrapers/OSCv3/data/selangor-mbpj-1003/AR_776177/FR_528192_Form4/details.yml
	arID := "778616"
	formID := "422667"
	formNum := "Form2"

	pathOfFormDetails := fmt.Sprintf("./data/%s/AR_%s/FR_%s_%s/details.yml", uniqueSearchID, arID, formID, formNum)
	b, rerr := ioutil.ReadFile(pathOfFormDetails)
	if rerr != nil {
		panic(rerr)
	}
	formDetails := FormDetails{}
	umerr := yaml.Unmarshal(b, &formDetails)
	if umerr != nil {
		panic(umerr)
	}
	// Display in a nice looking text way ..
	convertFormHistoryHTML(formDetails.StatusTerkiniAT)
}

// syncTracking will pull latest raw data and  check their Application date/status?
// Try to use existing data to use a simple  regexp instead? possible?
//func syncTracking(authorityToScrape string) {
//	// Load up ApplicationTracking metedata tate for use in checking ..
//	// Assumes the URL is S=S; until we can prove otherwise ..
//	// Can refactor previous function ..
//}
