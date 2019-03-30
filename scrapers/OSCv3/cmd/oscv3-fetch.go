package cmd

// NOTE: All raw data here will be stored under the following pattern
// ./raw/<uniqueSearchID>/<ApplicationID>/

// FetchAll will Extract from authority + label; all 15 pages of the information
func FetchAll() {
	// Raw structure like .. ./raw/<snapshotLabel>/<uniqueSearchID>
	// NOTE: Descructive action will override the data; ensure it is git diff ..

	// Step #1: Extract into marshal structure
	// Store into metadata structure like ./data/<uniqueSearchID>/
	// e.g. ./data/selangor-mbpj-1003/tracking.yaml; append only new unique items;
	//	sorted by ApplicationID
	// marked the successful / completed into archive? <-- Done in another step
}

// FetchNew will only Extract the New items per authority mapping
func FetchNew() {
	// Metadata structure like ./data/<uniqueSearchID>-<snapshotDiffLabels>
	// e.g. ./data/selangor-mbpj-1003-20190330_20190317

}
