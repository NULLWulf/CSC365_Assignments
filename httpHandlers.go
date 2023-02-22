package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func homepage(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	fmt.Println("Serving homepage")
	http.ServeFile(writer, request, "./html/homepage.html")
}

func returnRandomBusinessListJson(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	fmt.Println("Serving random business")
	businesses := getRandomBusinessList(10)

	// marshal the array of structs to JSON

	jsonBytes, err := json.Marshal(businesses)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}

	// write the JSON response to the client
	writer.Header().Set("Content-Type", "application/json")
	_, err = writer.Write(jsonBytes)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}
}

func getRelatableBusinesses(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	// get the business_id parameter from the query string
	businessID := request.URL.Query().Get("business_id")
	// get the relatable businesses
	relatableBusinesses := findRelatableBusinesses(businessID)
	// make bizTouple array
	bizTouples := make([]bizTuple, 0)
	for _, b := range relatableBusinesses {
		// append a new bizTuple to the bizTouples array
		bizTouples = append(bizTouples, bizTuple{
			BusinessName: b.Name,
			BusinessID:   b.BusinessID,
		})
	}

	// marshal the relatableBusinesses map to JSON and write to the response
	jsonBytes, err := json.Marshal(bizTouples)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}
	writer.Header().Set("Content-Type", "application/json")
	writer.Write(jsonBytes)
}
