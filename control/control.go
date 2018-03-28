package control

import (
	"fmt"
	"time"

	"github.com/jefferickson/two-americas/model"
	"github.com/jefferickson/two-americas/twitter"
	"github.com/jefferickson/two-americas/util"
)

type Controller struct {
	CountyListingsFilename string
	TopicsFilename         string
	TwitterMaxPer15Min     int
}

// Run is where it all happens
func (settings *Controller) Run() {
	// import datasets
	fmt.Println("Importing CSVs...")
	topics := util.ImportTopics(settings.TopicsFilename)
	countyListings := util.ImportCountyListings(settings.CountyListingsFilename)

	// add topics X counties to mongodb if they don't already exist
	fmt.Println("Adding county/topic pairs to DB if they do not exist...")
	model.InsertTopicCountyPairsToDB(topics, countyListings)

	// start timer so that we don't hit the Twitter throttle
	tick := time.Tick(time.Duration(15/float64(settings.TwitterMaxPer15Min)*60) * time.Second)

	// start off the process!
	fmt.Println("Starting tweet extraction...")
	for {
		start := time.Now()

		// look up one topic/county to run (based on last run date)
		countyTopicToRun := model.FindStaleTopicCounty()
		if countyTopicToRun == nil {
			fmt.Println("Something didn't work. Trying again.")
			continue
		}
		fmt.Println("Running:", countyTopicToRun, "at", start)

		// spawn go routine to fetch the data and insert into the db
		countyListingToRun := countyListings[countyTopicToRun.GeoID]
		go twitter.FetchCountyTopicAndInsertToDB(&countyListingToRun, countyTopicToRun)

		// sleep until tick
		<-tick

		// print time diagnostics
		loop_time := time.Since(start)
		runs_per_15m := 15 * 60 / float64(loop_time/time.Second)
		fmt.Println("Loop time:", loop_time)
		fmt.Println("Runs per 15m: about", runs_per_15m)
	}
}
