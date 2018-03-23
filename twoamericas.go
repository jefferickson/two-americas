package main

import (
	"github.com/jefferickson/two-americas/control"
)

const (
	// datasets
	countyListing = "datasets/county_listing.csv"
	topics        = "datasets/topics.csv"

	// Twitter throlling
	maxPer15Min = 180
)

func main() {
	controller := control.Controller{
		CountyListingsFilename: countyListing,
		TopicsFilename:         topics,
		TwitterMaxPer15Min:     maxPer15Min,
	}

	controller.Run()
}
