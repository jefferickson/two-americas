package util

import (
	"encoding/csv"
	"log"
	"os"
	"strconv"

	"github.com/jefferickson/two-americas/model"
)

// ImportTopics imports topics CSV
func ImportTopics(inputFile string) []model.Topic {
	csvfile, err := os.Open(inputFile)
	if err != nil {
		log.Panic(err)
	}
	defer csvfile.Close()

	reader := csv.NewReader(csvfile)
	reader.FieldsPerRecord = 1

	rawCSVdata, err := reader.ReadAll()
	if err != nil {
		log.Panic(err)
	}

	topics := []model.Topic{}
	for _, row := range rawCSVdata {
		topics = append(topics, model.Topic(row[0]))
	}

	return topics
}

// ImportCountyListings imports county listing CSV
func ImportCountyListings(inputFile string) model.CountyListings {
	csvfile, err := os.Open(inputFile)
	if err != nil {
		log.Panic(err)
	}
	defer csvfile.Close()

	reader := csv.NewReader(csvfile)
	reader.FieldsPerRecord = 0

	rawCSVdata, err := reader.ReadAll()
	if err != nil {
		log.Panic(err)
	}

	counties := make(model.CountyListings)
	for _, row := range rawCSVdata {
		GID, err := strconv.Atoi(row[0])
		if err != nil {
			panic(err)
		}

		MaxRadius, err := strconv.ParseFloat(row[6], 64)
		if err != nil {
			panic(err)
		}

		Lon, err := strconv.ParseFloat(row[7], 64)
		if err != nil {
			panic(err)
		}

		Lat, err := strconv.ParseFloat(row[8], 64)
		if err != nil {
			panic(err)
		}

		counties[row[1]] = model.CountyGeo{
			GID:       GID,
			GeoID:     row[1],
			State:     row[2],
			County:    row[3],
			Name:      row[4],
			Lsad:      row[5],
			MaxRadius: MaxRadius,
			Lon:       Lon,
			Lat:       Lat,
		}
	}

	return counties
}
