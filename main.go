package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"

	"grabjobs-by-location/controllers"
	"grabjobs-by-location/misc"
	//"strconv"
)

var db = ""

// Existing code from above
func handleRequests() {
	// creates a new instance of a mux router
	
	myRouter := mux.NewRouter().StrictSlash(true)
	// replace http.HandleFunc with myRouter.HandleFunc
	//myRouter.HandleFunc("/",)
	myRouter.HandleFunc("/near_by", controllers.SearchLatitude).Methods("GET")
	//myRouter.HandleFunc("/all", returnAllArticles)
	// finally, instead of passing in nil, we want
	// to pass in our newly created router as the second
	// argument
	log.Fatal(http.ListenAndServe(":10000", myRouter))
}



func main() {
	fmt.Println("GrabsJobs - By Location")
	/*Articles = []Article{
	    Article{Title: "Hello", Desc: "Article Description", Content: "Article Content"},
	    Article{Title: "Hello 2", Desc: "Article Description", Content: "Article Content"},
	}*/
	processDb()
	//fmt.Println("------------------------------------")
	//fmt.Println(db)
	//fmt.Println("------------------------------------")
	handleRequests()

}

func processDb () {
	csvToDb.PrepareDatabase()
}


