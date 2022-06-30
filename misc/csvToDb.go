package csvToDb

import (
	"fmt"

	"encoding/csv"
	"encoding/json"
	"log"
	"os"

)

type locationData struct {
	Title     string `json:"title"`
	Latitude  string `json:"latitude"`
	Longitude string `json:"longitude"`
}
var LocationData []locationData


func PrepareDatabase () {
	readCSV()
}


func createLocationData(data [][]string) []locationData {
	var locationList []locationData
	for i, line := range data {
		if i > 0 { // omit header line
			var rec locationData
			for j, field := range line {
				if j == 0 {
					rec.Title = field
				} else if j == 1 {
					rec.Latitude = field
				} else if j == 2 {
					var err error
					rec.Longitude = field
					//err = strconv.Atoi(field)
					if err != nil {
						continue
					}
				}
			}
			locationList = append(locationList, rec)
		}
	}
	return locationList
}

func FetchDB () {
	//LocationData
}

func readCSV() {
	fmt.Println("Reading the CSV file ..")
	f, err := os.Open("location_data_2000.csv")
	if err != nil {
		log.Fatal(err)
	}

	// remember to close the file at the end of the program
	defer f.Close()

	// read csv values using csv.Reader
	csvReader := csv.NewReader(f)
	data, err := csvReader.ReadAll()
	//fmt.Println(data)

	_locationData := createLocationData(data)

	// 4. Convert an array of structs to JSON using marshaling functions from the encoding/json package
	jsonData, err := json.MarshalIndent(_locationData, "", "  ")
	//var m []LocationData

	if err != nil {
		log.Fatal(err)
	}
	json.Unmarshal([]byte(jsonData),&LocationData)
	//fmt.Println("Location Data :%+v ",LocationData)
	//return LocationData
	//LocationData = []locationData{jsonData}
	/*for {
	    rec, err := csvReader.Read()
	    if err == io.EOF {
	        break
	    }
	    if err != nil {
	        log.Fatal(err)
	    }
	    // do something with read line
	    fmt.Printf("%+v\n", rec)


	}*/
}