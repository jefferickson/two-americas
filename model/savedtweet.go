package model

import (
	"strings"
	"time"

	"github.com/ChimeraCoder/anaconda"
	"github.com/globalsign/mgo"
)

type SavedTweet struct {
	TweetID                string
	Status                 anaconda.Tweet
	SearchResponseMetadata anaconda.SearchMetadata
	Counter                Counter
	CountyListing          CountyGeo
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
			return false, nil
		} else {
			return false, err
		}
	}

	return true, nil
}
