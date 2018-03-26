package model

import (
	"fmt"
	"strings"
	"time"

	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
)

type Counter struct {
	GeoIDTopic       string
	GeoID            string
	Topic            string
	Count            int
	DatetimeInserted time.Time
}

// IncrementCounter increments the counter of a county topic
func (ct *Counter) IncrementCounter() error {
	// connect to db
	session := connectToDB()
	defer session.Close()

	// collection counters
	c := session.DB(settings.db).C(counterCollection)

	// find and increment
	err := c.Update(bson.M{strings.ToLower("GeoIDTopic"): ct.GeoIDTopic},
		bson.M{"$inc": bson.M{"count": 1}})

	return err
}

// InsertTopicCountyPairsToDB adds (if dne) topic x county pairs to DB
func InsertTopicCountyPairsToDB(topics []Topic, counties CountyListings) {
	// connect to db
	session := connectToDB()
	defer session.Close()

	// collection counters
	c := session.DB(settings.db).C(counterCollection)

	// Index for key
	keyIndex := mgo.Index{
		Key:        []string{strings.ToLower("GeoIDTopic")},
		Unique:     true,
		DropDups:   true,
		Background: false,
		Sparse:     true,
	}
	err := c.EnsureIndex(keyIndex)
	if err != nil {
		panic(err)
	}

	// iterate over counties, topics and insert (if needed)
	for _, topic := range topics {
		for geoID := range counties {
			record := Counter{
				GeoIDTopic:       geoID + ":" + string(topic),
				GeoID:            geoID,
				Topic:            string(topic),
				Count:            counterInit,
				DatetimeInserted: time.Now(),
			}

			err := c.Insert(&record)
			if err != nil {
				if !mgo.IsDup(err) {
					panic(err)
				}
			} else {
				fmt.Println("Inserted new county/topic: " + record.GeoIDTopic)
			}
		}
	}
}

// FindTopicCountyWithLowCount will return a (random) county/topic that has the
// lowest number of runs so far
func FindTopicCountyWithLowCount() *Counter {
	// connect to db
	session := connectToDB()
	defer session.Close()

	// collection counters
	c := session.DB(settings.db).C(counterCollection)

	// sort ascending by count and then pick one from bottom
	var results []Counter
	err := c.Find(bson.M{}).Sort("count").Limit(1).All(&results)
	if err != nil {
		fmt.Println(err)
		return nil
	}

	return &results[0]
}
