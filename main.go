package main

import (
	"fmt"
	//"github.com/gorilla/handlers"
	"github.com/rs/cors"
	"github.com/gorilla/mux"
	"log"
	"net/http"

	"grabjobs-by-location/controllers"
	"grabjobs-by-location/misc"
	//"strconv"
)

func handleRequests() {
	
	myRouter := mux.NewRouter().StrictSlash(true)
	
	myRouter.HandleFunc("/near_by", controllers.SearchLatitude).Methods("GET")
	myRouter.HandleFunc("/all_data", controllers.GetLocationData).Methods("GET")
	myRouter.HandleFunc("/test", controllers.Test).Methods("GET")
	myRouter.HandleFunc("/calculateRadius", controllers.CalculateRadius).Methods("GET")
	
	//handle cors acces
	c := cors.New(cors.Options{
        AllowedOrigins: []string{"*"},
        AllowCredentials: true,
    })
	handler := c.Handler(myRouter)

	log.Fatal(http.ListenAndServe(":10000", handler))
}



func main() {
	fmt.Println("GrabsJobs - By Location")
	
	processDb()
	handleRequests()

}

func processDb () {
	csvToDb.PrepareDatabase()
}


