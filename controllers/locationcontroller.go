package controllers

import(
	"fmt"
	"net/http"
	//"strconv"
	//"github.com/gorilla/mux"
	//"grabjobs-by-location/misc"
)

func SearchLatitude (res http.ResponseWriter, req *http.Request) {
	vars := req.URL.Query()
	fmt.Println(vars)
	lat := vars.Get("lat")
	long := vars.Get("lng")
	title := vars.Get("title")
	radius := vars.Get("radius")
	fmt.Fprint(res, "Requesting data on " + lat)
	fmt.Fprint(res, "Requesting data on " + long)
	fmt.Fprint(res, "Requesting data on " + title)
	fmt.Fprint(res, "Requesting data on " + radius)
}