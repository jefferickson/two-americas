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
	selector := bson.M{"status.fulltext": 1,
		"tweetid": 1,
		"counter": 1}

	// Fetch all tweets
	allTweets := model.AllTweets(filter, selector)

	// Loop over tweets and write out
	for _, tweet := range allTweets {
		fmt.Println(tweet.TweetID, "\t",
			tweet.Counter.GeoID, "\t",
			tweet.Counter.Topic, "\t",
			tweet.Status.FullText,
		)
	}
}
