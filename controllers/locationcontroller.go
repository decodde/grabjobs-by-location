package controllers

import (
	"encoding/json"
	"fmt"
	"math"
	"net/http"
	//"bytes"
	"regexp"
	"strconv"
	//"strconv"
	//"github.com/gorilla/mux"
	"github.com/ledongthuc/goterators"
	"grabjobs-by-location/misc"
	"grabjobs-by-location/models"
)

func Test(res http.ResponseWriter, req *http.Request) {
	json.NewEncoder(res).Encode(models.Response{
		Message: "true", Success: true,
	})
}
func SearchLatitude(res http.ResponseWriter, req *http.Request) {
	vars := req.URL.Query()
	fmt.Println(vars)
	lat := vars.Get("lat")
	long := vars.Get("lng")
	title := vars.Get("title")
	radius := vars.Get("radius")
	fmt.Print("Requesting data on " + radius)

	//getting db from memory
	var a = csvToDb.FetchDB()

	//title takes precedence, so check if title has been added as a filter
	//if yes, filter off by title; before lat and long
	//or we could take it as far as considering the order of the filter keys :)...v2 probably
	var filteredResult []models.LocationData
	r, err := regexp.Compile(title)
	fmt.Println(err)
	fmt.Println(r)
	if r == nil {
		json.NewEncoder(res).Encode(models.Response{
			Message: "error", Success: false,
		})
	} else {

		if title != "" {
			filteredResult = goterators.Filter(a, func(i models.LocationData) bool {
				_find := r.FindStringSubmatch(i.Title) //title,i.Title)
				if _find != nil {
					/*json.NewEncoder(res).Encode(models.Response{
						Success : false, Message : "Error whilst performing search",Data: filteredResult,
					})*/
					return true
				}
				return false
			})
			//fmt.Print(filteredResult)
		} else {
			filteredResult = a
		}
		var _filter1, _filter2 []models.LocationData

		if lat != "nil" {
			fmt.Print("filtering latitude")
			fmt.Println(lat)
			_filter1 = goterators.Filter(filteredResult, func(i models.LocationData) bool {
				if lat <= i.Latitude {
					fmt.Println(i.Latitude, lat)
					return true
				} else {
					fmt.Println("Here checking each latitude")
					fmt.Println(i.Latitude, lat)
					return false
				}
			})
			fmt.Println(_filter1)
			fmt.Println(_filter2)

		}
		if long != "nil" {
			fmt.Print("filtering longitude")
			fmt.Print(long)
			_filter2 = goterators.Filter(_filter1, func(i models.LocationData) bool {

				if long <= i.Longitude {
					fmt.Println(i.Longitude, long)
					return true
				} else {
					fmt.Println("Here checking each longitude")
					fmt.Println(i.Longitude, long)
					return false
				}
			})
		}
		dataLength := len(_filter2)

		json.NewEncoder(res).Encode(models.Response{
			Success: true, Message: "Data retrieved successfully", Data: _filter2, DataLength: dataLength,
		})
	}

}

func GetLocationData(res http.ResponseWriter, req *http.Request) {
	//json.
	var a = csvToDb.FetchDB()
	for i, v := range a {
		fmt.Print(i, " : ", v.Title)
	}
	json.NewEncoder(res).Encode(a)

}

func calculateRadius(lat float64,  distance float64, long float64) models.RadiusLocation {
	//Assuming data is in radians and not degrees
	

	var EARTH_RADIUS_KM = 6371.01     // Earth's radius in km
	//var EARTH_RADIUS_MI = 3958.762079 // Earth's radius in miles
	var MAX_LAT = 3.142 / 2           // 90 degrees
	var MIN_LAT = -MAX_LAT            // -90 degrees
	var MAX_LON = 3.142               // 180 degrees
	var MIN_LON = -MAX_LON            // -180 degrees
	var FULL_CIRCLE_RAD = 3.142 * 2   // Full cirle (360 degrees) in radians
	var radius = EARTH_RADIUS_KM
	var radDist float64 = distance / radius
	var minLat = lat - float64(radDist)
	var maxLat = lat + float64(radDist)
	var minLon float64
	var deltaLon float64
	var maxLon float64
	if minLat > MIN_LAT && maxLat < MAX_LAT {
		deltaLon = math.Asin(math.Sin(radDist) / math.Cos(lat))
		minLon = long - deltaLon

		if minLon < MIN_LON {
			minLon += FULL_CIRCLE_RAD
		}
		maxLon = long + deltaLon
		if maxLon > MAX_LON {
			maxLon -= FULL_CIRCLE_RAD
		}
	} else {
		fmt.Println("Here cs of shit")
		minLat = math.Max(minLat, MIN_LAT)
		maxLat = math.Min(maxLat, MAX_LAT)
		minLon = MIN_LON
		maxLon = MAX_LON
	}
	return models.RadiusLocation{
		MaxLongitude: maxLon, MaxLatitude: maxLat, MinLongitude: minLon, MinLatitude: minLat,
	}
}

func CalculateRadius (res http.ResponseWriter, req *http.Request){
	vars := req.URL.Query()
	fmt.Println(vars)
	lat := vars.Get("lat")
	long := vars.Get("lng")
	distance := vars.Get("distance")
	radius := vars.Get("radius")
	_lat,e := strconv.ParseFloat(lat,64)
	_long,e := strconv.ParseFloat(long,64)
	
	_radius,e := strconv.ParseFloat(radius,64)
	
	fmt.Print(e)
	fmt.Print(distance)

	json.NewEncoder(res).Encode(calculateRadius(_lat,_radius,_long))
}
