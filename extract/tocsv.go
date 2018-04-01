package main

import (
	"fmt"

	"github.com/globalsign/mgo/bson"
	"github.com/jefferickson/two-americas/model"
)

func main() {
	extractTweets()
}

func extractTweets() {
	// Set filters, selectors (can be empty for all tweets)
	filter := bson.M{}
	selector := bson.M{"countylisting": 1}

	// Fetch all tweets
	allTweets := model.AllTweets(filter, selector)

	// Loop over tweets and write out
	for _, tweet := range allTweets {
		fmt.Println(tweet.CountyListing.Lon,
			tweet.CountyListing.Lat,
			tweet.CountyListing.State)
	}
}
