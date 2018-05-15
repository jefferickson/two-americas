package main

import (
	"encoding/csv"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/globalsign/mgo/bson"
	"github.com/jefferickson/two-americas/model"
)

const (
	TWEETS_OUTPUT_FILE      = "tmp_datasets/tweets.tsv"
	TWEETS_URLS_OUTPUT_FILE = "tmp_datasets/urls.tsv"
)

func main() {
	extractTweets()
}

func extractTweets() {
	// Set filters, selectors (can be empty for all tweets)
	filter := bson.M{}
	selector := bson.M{
		"tweetid":                         1,
		"status.user.screenname":          1,
		"status.createdat":                1,
		"status.entities.urls":            1,
		"counter.geoid":                   1,
		"counter.topic":                   1,
		"countylisting.state":             1,
		"countylisting.name":              1,
		"countylisting.lon":               1,
		"countylisting.lat":               1,
		"status.fulltext":                 1,
		"status.retweetedstatus.fulltext": 1,
	}

	// tweet file
	tweetFile, err := os.Create(TWEETS_OUTPUT_FILE)
	if err != nil {
		fmt.Println("Could not create file", err)
		return
	}
	defer tweetFile.Close()

	// tweet writer
	tweetWriter := csv.NewWriter(tweetFile)
	tweetWriter.Comma = '\t'
	defer tweetWriter.Flush()

	// tweet URL file
	tweetURLFile, err := os.Create(TWEETS_URLS_OUTPUT_FILE)
	if err != nil {
		fmt.Println("Count not create file", err)
		return
	}
	defer tweetURLFile.Close()

	// tweet urls writer
	tweetURLWriter := csv.NewWriter(tweetURLFile)
	tweetURLWriter.Comma = '\t'
	defer tweetURLWriter.Flush()

	// Create channel and get tweets concurrently
	tweetChan := make(chan model.SavedTweet)
	go model.IterTweets(filter, selector, tweetChan)

	// Loop over tweets and write out
	for tweet := range tweetChan {
		// Write out tweet record
		err := writeTweetRec(&tweet, tweetWriter)
		if err != nil {
			fmt.Println("Cannot write line", err)
			return
		}

		// Write out tweet URL record(s)
		err = writeTweetURLRec(&tweet, tweetURLWriter)
		if err != nil {
			fmt.Println("Cannot write line", err)
			return
		}
	}
}

func writeTweetRec(tweet *model.SavedTweet, writer *csv.Writer) error {
	// delete newlines, carriage returns within tweets
	cleanTweet := func(t string) string {
		t = strings.Replace(t, "\n", " ", -1)
		t = strings.Replace(t, "\r", " ", -1)
		return t
	}

	// clean tweet or retweet
	var tweetToPrint, retweet string
	if tweet.Status.RetweetedStatus == nil {
		tweetToPrint = cleanTweet(tweet.Status.FullText)
		retweet = "false"
	} else {
		tweetToPrint = cleanTweet(tweet.Status.RetweetedStatus.FullText)
		retweet = "true"
	}

	// what to extract
	tweetRec := []string{
		tweet.TweetID,
		tweet.Status.User.ScreenName,
		tweet.Status.CreatedAt,
		tweet.Counter.GeoID,
		tweet.Counter.Topic,
		tweet.CountyListing.State,
		tweet.CountyListing.Name,
		strconv.FormatFloat(tweet.CountyListing.Lon, 'f', -1, 64),
		strconv.FormatFloat(tweet.CountyListing.Lat, 'f', -1, 64),
		retweet,
		tweetToPrint,
	}

	// write out
	err := writer.Write(tweetRec)
	return err
}

func writeTweetURLRec(tweet *model.SavedTweet, writer *csv.Writer) error {
	// for each URL, extract the domain and write out record
	for _, url := range tweet.Status.Entities.Urls {
		re := regexp.MustCompile(`^(?:https?:\/\/)?(?:[^@\/\n]+@)?(?:www\.)?([^:\/\n]+)`)
		domainFound := re.FindString(url.Expanded_url)

		tweetURLRec := []string{
			tweet.TweetID,
			domainFound,
		}

		//write out
		err := writer.Write(tweetURLRec)
		if err != nil {
			return err
		}
	}

	return nil
}
