package main

import (
	"encoding/json"
	"github.com/julienschmidt/httprouter"
	"log"
	"net/http"
	"strconv"
)

// homepage Serves homepage (index.html)
func homepage(writer http.ResponseWriter, request *http.Request) {
	log.Println("Serving homepage")
	http.ServeFile(writer, request, "./html/homepage.html")
}

// homepage Serves homepage (index.html)
func homepage2(writer http.ResponseWriter, request *http.Request) {
	log.Println("Serving homepage")
	http.ServeFile(writer, request, "./html/homepage2.html")
}

// Gets a random list of up 10 businesses and returns to front end
// in this case the list is appended to a drop-down menu.  Used in Assignment 1 and 2
func returnRandomBusinessListJson(writer http.ResponseWriter, request *http.Request, km *KmediodsDS) {
	log.Printf("returnRandomBusinessListJson")
	count := 10
	bsDps := km.GetRandomDataPoints(count)
	var businesses []Business
	for _, v := range bsDps {
		t := LoadBusinessFromFile(strconv.Itoa(v.FileIndex))
		businesses = append(businesses, t)
	}

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
// and returns 2 reletable businesses to front end.  Used in Assignment 1.
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

// getRelatableCluster takes a business id and returns the business, cluster count, medoid, and teh most
// relatable business in the cluster.  Used in Assignment 2.
func getRelatableCluster(writer http.ResponseWriter, request *http.Request, kmed *KmediodsDS) {
	log.Printf("getRelatableCluster called for: %s", request.URL.Query().Get("file_id"))
	fileId := request.URL.Query().Get("file_id")
	// get the cluster that the business belongs to
	jsonBytes, err := kmed.FindSimilarBuildResponse(fileId)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}
	writer.Header().Set("Content-Type", "application/json")
	writer.Write(jsonBytes)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}
}
