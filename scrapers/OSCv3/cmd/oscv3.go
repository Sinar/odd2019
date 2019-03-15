package main

import (
	"fmt"

	"github.com/gocolly/colly"
	"github.com/y0ssar1an/q"
)

func main() {
	fmt.Println("Welcome to GOMOD OSCv3!!")
	BasicCollyFromRaw()
}

// BasicCollyFromRaw is meant to read data from fixtures and extract out data ..
func BasicCollyFromRaw() {
	// From the collection; can run another round while tweaking the strcuture
	// Removing the extra cost of network and being blocked ..

	c := colly.NewCollector(
		colly.UserAgent("Sinar Project :P"),
	)
	c.OnError(func(r *colly.Response, e error) {
		fmt.Println("DIE!!!")
		q.Q(r.Request.URL, r.StatusCode)
		panic(e)
	})
}
