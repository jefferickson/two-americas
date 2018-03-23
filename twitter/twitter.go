package twitter

import (
	"fmt"
	"net/url"
	"os"
	"strconv"
	"time"

	"github.com/ChimeraCoder/anaconda"
	"github.com/jefferickson/two-americas/model"
)

const maxTweetsPerRun = 100

type twitterCreds struct {
	accessToken    string
	tokenSecret    string
	consumerKey    string
	consumerSecret string
}

var creds = twitterCreds{
	os.Getenv("TWITTER_ACCESS_TOKEN"),
	os.Getenv("TWITTER_TOKEN_SECRET"),
	os.Getenv("TWITTER_CONSUMER_KEY"),
	os.Getenv("TWITTER_CONSUMER_SECRET"),
}

// TwitterConnect creates connection with Twitter
func twitterConnect() *anaconda.TwitterApi {
	return anaconda.NewTwitterApiWithCredentials(
		creds.accessToken,
		creds.tokenSecret,
		creds.consumerKey,
		creds.consumerSecret)
}

// FetchCountyTopicAndInsertToDB will fetch the tweets from the supplied County and Topic
func FetchCountyTopicAndInsertToDB(countyListing *model.CountyGeo, countyTopic *model.Counter) {
	// connect to Twitter API and Mongo
	api := twitterConnect()
	api.ReturnRateLimitError(true)

	// search settings
	v := url.Values{}
	v.Set("lang", "en")
	v.Set("count", strconv.Itoa(maxTweetsPerRun))

	// geocode is lat, long, radius (with units)
	geocode := strconv.FormatFloat(countyListing.Lat, 'f', 6, 64) +
		"," + strconv.FormatFloat(countyListing.Lon, 'f', 6, 64) +
		"," + strconv.FormatFloat(countyListing.MaxRadius, 'f', 1, 64) + "km"
	v.Set("geocode", geocode)

	// run the search
	searchResult, err := api.GetSearch(countyTopic.Topic, v)
	if err != nil {
		fmt.Println("Could not execute search.", err)
		return
	}

	// store each in the db, adding additional info
	nTweets := 0
	for _, tweet := range searchResult.Statuses {
		tweetToSave := model.SavedTweet{
			TweetID:                tweet.IdStr,
			Status:                 tweet,
			SearchResponseMetadata: searchResult.Metadata,
			Counter:                *countyTopic,
			CountyListing:          *countyListing,
			DatetimeInserted:       time.Now(),
		}

		inserted, err := tweetToSave.InsertTweet()
		if err != nil {
			fmt.Println(err)
		} else if inserted {
			nTweets++
		}
	}
	fmt.Println(nTweets, "tweets added.")

	// increment counter
	err = countyTopic.IncrementCounter()
	if err != nil {
		fmt.Println(err)
	}
}
