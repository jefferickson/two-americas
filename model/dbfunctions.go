package model

import (
	"os"
	"time"

	"github.com/globalsign/mgo"
)

const (
	counterCollection = "counters"
	tweetsCollection  = "tweets"
	counterInit       = 0
)

var dateTimeLastRunInit, _ = time.Parse(time.RFC3339, "2006-03-21T00:00:00Z")

type mongoSettings struct {
	host string
	db   string
}

var settings = mongoSettings{
	os.Getenv("MONGODB_HOST"),
	os.Getenv("MONGODB_DB"),
}

// ConnectToDB returns a db session. Don't forget to close it!
func connectToDB() *mgo.Session {
	session, err := mgo.Dial(settings.host)
	if err != nil {
		panic(err)
	}

	return session
}
