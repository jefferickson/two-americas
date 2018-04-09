package main

import (
	"encoding/csv"
	"fmt"
	"os"
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
		"counter": 1,
	}

	// Fetch all tweets
	allTweets := model.AllTweets(filter, selector)

	// file
	file, err := os.Create("tmp_datasets/tweets.tsv")
	if err != nil {
		fmt.Println("Could not create file", err)
		return
	}
	defer file.Close()

	// writer
	writer := csv.NewWriter(file)
	writer.Comma = '\t'
	defer writer.Flush()

	// Loop over tweets and write out
	for _, tweet := range allTweets {
		// delete newlines, carriage returns within tweets
		clean_fulltext := strings.Replace(tweet.Status.FullText, "\n", " ", -1)
		clean_fulltext = strings.Replace(clean_fulltext, "\r", " ", -1)

		// what to extract
		tweetRec := []string{
			tweet.TweetID,
			tweet.Counter.GeoID,
			tweet.Counter.Topic,
			clean_fulltext,
		}

		// write out
		err := writer.Write(tweetRec)
		if err != nil {
			fmt.Println("Cannot write line", err)
			return
		}
	}
}
