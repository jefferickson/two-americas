package model

type CountyListings map[string]CountyGeo

type CountyGeo struct {
	GID       int
	GeoID     string
	State     string
	County    string
	Name      string
	Lsad      string
	MaxRadius float64
	Lon       float64
	Lat       float64
}
