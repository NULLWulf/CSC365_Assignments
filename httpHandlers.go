package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

// homepage Serves homepage (index.html)
func homepage(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	fmt.Println("Serving homepage")
	http.ServeFile(writer, request, "./html/homepage.html")
}

// Gets a random list of up 10 businesses and returns to front end
// in this case the list is appended to a drop down menu.
func returnRandomBusinessListJson(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	businesses := getRandomBusinessList(10)
	fmt.Printf("Serving random businesses: %d", len(businesses))

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
