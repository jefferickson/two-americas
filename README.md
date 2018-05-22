## "Two Americas" Tweet Collector

#### Author: Jeff Erickson `<jeff@erick.so>`
#### Date: 2018-05-21

### Overview

A system to collect tweets from every county in the contiguous United States on a selection of politically-relevant search topics.

This system was used to collect data for my Masters capstone project at the Department of Mathematics and Statistics, Hunter College, The City University of New York.

### Method of Search

The purpose of this system is to collect tweets systematically from around the United States. The search terms themselves are defined in advance: see `datasets/topics.csv`.

The Twitter API allows a user to specify a "search circle" to search within. A latitude, longitude, and radius must be provided. A "minimum bounding search circle" for each county has been defined in order to approximate each county: see `datasets/county_listing.csv`.

### Data Storage

Data is stored in a MongoDB NoSQL database.

### Setup and Usage

1. Get dependencies:

```
go get github.com/ChimeraCoder/anaconda
go get github.com/globalsign/mgo
go get gopkg.in/cheggaaa/pb.v1
```

2. Get the repo and install:

```
go get github.com/jefferickson/two-americas
go install github.com/jefferickson/two-americas
```

3. Get Twitter credentials [here](https://twitter.com/settings/applications). Copy the file `.env.example` into `.env` and provide these credentials. Also include the connection details to MongoDB.

4. Source the `.env` and run the service:

```
source .env
$GOPATH/bin/two-americas
```

The service will run as long as you let it.

5. When ready, extract your data: A modifiable utility script is provided (see `extract/toTSV.go`) to extract the data into a ["tidy" format](https://vita.had.co.nz/papers/tidy-data.pdf). Once you have specified the fields you want to export, run:

```
go run $GOPATH/src/github.com/jefferickson/two-americas/extract/toTSV.go
```
