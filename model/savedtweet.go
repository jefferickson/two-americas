package model

import (
	"fmt"
	"strings"
	"time"

	"github.com/ChimeraCoder/anaconda"
	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
	"gopkg.in/cheggaaa/pb.v1"
)

type SavedTweet struct {
	TweetID                string
	Status                 anaconda.Tweet
	SearchResponseMetadata anaconda.SearchMetadata
	Counter                Counter
	CountyListing          CountyGeo
	SeenIn                 map[string]int
	DatetimeInserted       time.Time
}

// InsertTweet inserts a SavedTweet object
func (st *SavedTweet) InsertTweet() (bool, error) {
	// connect to db
	session := connectToDB()
	defer session.Close()

	// tweets collection
	c := session.DB(settings.db).C(tweetsCollection)

	// Index for key
	tweetIndex := mgo.Index{
		Key:        []string{strings.ToLower("TweetID")},
		Unique:     true,
		DropDups:   true,
		Background: false,
		Sparse:     true,
	}
	err := c.EnsureIndex(tweetIndex)
	if err != nil {
		return false, err
	}

	// insert
	err = c.Insert(st)
	if err != nil {
		if mgo.IsDup(err) {
			// update SeenIn
			c.Upsert(bson.M{strings.ToLower("TweetID"): st.TweetID},
				bson.M{"$inc": bson.M{strings.ToLower("SeenIn") + "." + strings.ToLower(st.Counter.GeoID): 1}})

			return false, nil
		}

		return false, err
	}

	return true, nil
}

// AllTweets returns all tweets at once. Only use this if you have a lot of RAM
func AllTweets(filter, selector bson.M) []SavedTweet {
	// connect to db
	session := connectToDB()
	defer session.Close()

	// collection tweets
	c := session.DB(settings.db).C(tweetsCollection)

	// get all results
	var results []SavedTweet
	err := c.Find(filter).Select(selector).All(&results)
	if err != nil {
		fmt.Println(err)
		return nil
	}

	return results
}

// IterTweets yields one tweet at a time from the DB over a channel
func IterTweets(filter, selector bson.M, yield chan<- SavedTweet) {
	// connect to db
	session := connectToDB()
	defer session.Close()

	// collection tweets
	c := session.DB(settings.db).C(tweetsCollection)

	// set up the query
	query := c.Find(filter).Select(selector)

	// count docs for progress bar
	nDocuments, err := query.Count()
	if err != nil {
		fmt.Println(err)
	}
	bar := pb.StartNew(nDocuments)

	// now iter over the docs, reporting progress
	iter := query.Iter()
	result := SavedTweet{}
	for iter.Next(&result) {
		yield <- result
		bar.Increment()
	}
	if err := iter.Err(); err != nil {
		fmt.Println(err)
	}
	bar.Finish()

	close(yield)
}
