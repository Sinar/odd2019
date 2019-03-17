package cmd

import "fmt"

func main() {

	fmt.Println("OSCv3 Find the new requests!")
}

func extractDataFromPage() {

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

// FindNewRequests will look for the changes since the last time run and offer a pull request
func FindNewRequests() {
	// Figure out when was the last successful run and if not exist; create it!
	// Also will reset if passed a flag of some sort?

	// If in Codefresh; do a branch, git add + commit?
}
