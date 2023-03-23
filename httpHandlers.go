package main

import (
	"encoding/json"
	"github.com/julienschmidt/httprouter"
	"log"
	"net/http"
	"os"
)

// homepage Serves homepage (index.html)
func homepage(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	log.Println("Serving homepage")
	http.ServeFile(writer, request, "./html/homepage.html")
}

// homepage Serves homepage (index.html)
func homepage2(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	log.Println("Serving homepage")
	http.ServeFile(writer, request, "./html/homepage2.html")
}

// Gets a random list of up 10 businesses and returns to front end
// in this case the list is appended to a drop down menu.
func returnRandomBusinessListJson(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	businesses := getRandomBusinessList(10)
	log.Printf("Serving random businesses: %d", len(businesses))

	jsonBytes, err := json.Marshal(businesses)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}

	writer.Header().Set("Content-Type", "application/json")
	_, err = writer.Write(jsonBytes)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}
}

// Takes businesses id received from front end and calls find ofRelatableBusiness functions
// and returns
func getRelatableBusinesses(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	businessID := request.URL.Query().Get("business_id")
	relatableBusinesses := findRelatableBusinesses(businessID)
	// make bizTouple array of two businesses
	bizTouples := make([]bizTuple, 0)
	for _, b := range relatableBusinesses {
		// append a new bizTuple to the bizTouples array
		bizTouples = append(bizTouples, bizTuple{
			BusinessName: b.Name,
			BusinessID:   b.BusinessID,
		})
	}

	jsonBytes, err := json.Marshal(bizTouples)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}
	writer.Header().Set("Content-Type", "application/json")
	writer.Write(jsonBytes)
}

func getRelatableCluster(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	log.Printf("getRelatableCluster")
	businessID := request.URL.Query().Get("business_id")
	BizMap, _ := LoadHashMapFromFile("hashmap.json")
	a := BizMap.GetKeyList()

	var selectedBiz Business
	file, err := os.ReadFile("fileblock/" + businessID + ".json")
	if err != nil {
		log.Println(err)
	}
	err = json.Unmarshal(file, &selectedBiz)
	log.Printf("fileblock/%s.json: %+v", businessID, selectedBiz)

	// get the cluster of the selected business
	relatab := make(map[string]float32)
	for i := 0; i < len(a); i++ {
		file, err := os.ReadFile("fileblock/" + a[i] + ".json")
		if err != nil {
			log.Println(err)
		}
		var b Business
		err = json.Unmarshal(file, &b)
		if err != nil {
			log.Println(err)
		}
		relatab[b.BusinessID] = distance(
			BusinessDataPoint{
				BusinessID:  selectedBiz.BusinessID,
				ReviewScore: selectedBiz.Stars,
				Latitude:    selectedBiz.Latitude,
				Longitude:   selectedBiz.Longtitude,
			},
			BusinessDataPoint{
				BusinessID:  b.BusinessID,
				ReviewScore: b.Stars,
				Latitude:    b.Latitude,
				Longitude:   b.Longtitude,
			},
		)
	}

	var jsonBytes []byte
	jsonBytes, err = json.Marshal(relatab)
	log.Printf("jsonBytes: %+v", string(jsonBytes))

	writer.Header().Set("Content-Type", "application/json")
	writer.Write(jsonBytes)
}
