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

var EARTH_RADIUS_KM = 6371.01 // Earth's radius in km
//var EARTH_RADIUS_MI = 3958.762079 // Earth's radius in miles
var MAX_LAT = 3.142 / 2         // 90 degrees
var MIN_LAT = -MAX_LAT          // -90 degrees
var MAX_LON = 3.142             // 180 degrees
var MIN_LON = -MAX_LON          // -180 degrees
var FULL_CIRCLE_RAD = 3.142 * 2 // Full cirle (360 degrees) in radians
var radius = EARTH_RADIUS_KM

func Test(res http.ResponseWriter, req *http.Request) {
	json.NewEncoder(res).Encode(models.Response{
		Message: "true", Success: true,
	})
}
func SearchLatitude(res http.ResponseWriter, req *http.Request) {
	vars := req.URL.Query()

	lat := vars.Get("lat")
	long := vars.Get("lng")
	title := vars.Get("title")
	radius := vars.Get("radius")
	fmt.Print("Requesting data on " + lat)

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
		var _filter1 []models.LocationData

		if lat != "" {

			//CONVERT STRINGS To FLOAT
			_lat, e := strconv.ParseFloat(lat, 64)
			_long, e := strconv.ParseFloat(long, 64)
			_radius, e := strconv.ParseFloat(radius, 64)

			if e != nil {
				json.NewEncoder(res).Encode(models.Response{
					Success: false, Error: e, Message: "Conversion error",
				})
				return
			}

			var bounds = calculateRadiusBounds(_lat, _radius, _long, false)
			//fmt.Print("filtering latitude")
			//fmt.Println(lat)
			_filter1 = goterators.Filter(filteredResult, func(i models.LocationData) bool {
				_currentLatitude, e := strconv.ParseFloat(i.Latitude, 64)
				_currentLongitude, e := strconv.ParseFloat(i.Longitude, 64)
				//fmt.Println(e)
				if e != nil {
					json.NewEncoder(res).Encode(models.Response{
						Success: false, Error: e, Message: "Conversion error",
					})
					
				}

				/*fmt.Println("------------------------------------")
				fmt.Println(_currentLatitude, " min : ", bounds.MinLatitude, " max : ", bounds.MaxLatitude, _lat)
				fmt.Println(_currentLongitude, " min : ", bounds.MinLongitude, " max : ", bounds.MaxLongitude, _long)
				fmt.Println(_currentLatitude >= bounds.MinLatitude)
				fmt.Println(_currentLatitude <= bounds.MaxLatitude)
				fmt.Println(_currentLongitude <= bounds.MaxLongitude)
				fmt.Println(_currentLongitude >= bounds.MinLongitude)
				fmt.Println("------------------------------------")*/

				if _currentLatitude >= bounds.MinLatitude && _currentLatitude <= bounds.MaxLatitude {

					/*fmt.Println("------------------------------------")
					fmt.Println("Latitude this")
					fmt.Println(_currentLatitude, bounds.MinLatitude, bounds.MaxLatitude)
					fmt.Println("------------------------------------")*/

					if _currentLongitude >= bounds.MinLongitude && _currentLongitude <= bounds.MaxLongitude {
						return true
					} else {
						/*fmt.Println("------------------------------------")
						fmt.Println("Passed Latitude test but not this")
						fmt.Println(_currentLongitude, bounds.MinLongitude, bounds.MaxLongitude)
						fmt.Println("------------------------------------")
						*/
						return false
					}
				} else {

					return false
				}
			})
			//fmt.Println(_filter1)
			dataLength := len(_filter1)

			json.NewEncoder(res).Encode(models.Response{
				Success: true, Message: "Data retrieved successfully", Data: _filter1, DataLength: dataLength,
			})
			return
		}
		dataLength := len(filteredResult)

		json.NewEncoder(res).Encode(models.Response{
			Success: true, Message: "Data retrieved successfully", Data: filteredResult, DataLength: dataLength,
		})
		return
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

func calculateRadiusBounds(lat float64, distance float64, long float64, inRadians bool) models.RadiusLocation {
	//Adding check for radians or degree and ensurig degreee is used

	if !inRadians {
		lat = degreesToRadian(lat)
		long = degreesToRadian(long)
	}
	if distance > 0 {

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

		if !inRadians {
			return models.RadiusLocation{
				MaxLongitude: radianToDegrees(maxLon), MaxLatitude: radianToDegrees(maxLat), MinLongitude: radianToDegrees(minLon), MinLatitude: radianToDegrees(minLat),
			}
		} else {
			return models.RadiusLocation{
				MaxLongitude: long, MaxLatitude: lat, MinLongitude: long, MinLatitude: lat,
			}
		}

	} else {
		return models.RadiusLocation{
			MaxLongitude: long, MaxLatitude: lat, MinLongitude: long, MinLatitude: lat,
		}
	}
}

func CalculateRadius(res http.ResponseWriter, req *http.Request) {
	vars := req.URL.Query()
	fmt.Println(vars)
	lat := vars.Get("lat")
	long := vars.Get("lng")
	distance := vars.Get("distance")
	radius := vars.Get("radius")
	_lat, e := strconv.ParseFloat(lat, 64)
	_long, e := strconv.ParseFloat(long, 64)

	_radius, e := strconv.ParseFloat(radius, 64)

	fmt.Print(e)
	fmt.Print(distance)

	json.NewEncoder(res).Encode(calculateRadiusBounds(_lat, _radius, _long, false))
}

func radianToDegrees(value float64) float64 {
	var RAD2DEG = 180 / 3.142 // radians to degrees conversion

	return value * RAD2DEG
}
func degreesToRadian(value float64) float64 {
	var DEG2RAD = 3.142 / 180 // degrees to radian conversion

	return value * DEG2RAD
}
