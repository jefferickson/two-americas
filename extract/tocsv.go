package main

import (
	"fmt"
	"strings"

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

	// Write out header
	fmt.Print("tweetid", "\t", "geoid", "\t", "topic", "\t", "tweet", "\n")

	// Loop over tweets and write out
	for _, tweet := range allTweets {
		// delete newlines within tweets
		clean_fulltext := strings.Replace(tweet.Status.FullText, "\n", " ", -1)

		// write out
		fmt.Print(tweet.TweetID, "\t",
			tweet.Counter.GeoID, "\t",
			tweet.Counter.Topic, "\t",
			clean_fulltext, "\n",
		)
	}
}
