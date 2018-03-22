package main

import (
	"encoding/json"
	"fmt"
	"net/url"
	"os"

	"github.com/ChimeraCoder/anaconda"
)

func main() {
	api := twitterConnect()

	v := url.Values{}
	v.Set("lang", "en")
	v.Set("geocode", "32.514730,-93.748001,10mi")

	searchResult, err := api.GetSearch("Trump", nil)
	if err != nil {
		fmt.Println("Could not execute search.")
		return
	}

	for _, tweet := range searchResult.Statuses {
		tweetJSON, err := json.MarshalIndent(tweet, "", "\t")
		if err != nil {
			fmt.Println(err)
			return
		}

		os.Stdout.Write(tweetJSON)
	}
}

func twitterConnect() *anaconda.TwitterApi {
	return anaconda.NewTwitterApiWithCredentials(
		os.Getenv("TWITTER_ACCESS_TOKEN"),
		os.Getenv("TWITTER_TOKEN_SECRET"),
		os.Getenv("TWITTER_CONSUMER_KEY"),
		os.Getenv("TWITTER_CONSUMER_SECRET"))
}
