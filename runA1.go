package main

import (
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

// RUNA1 Runtime for Assignment 1 program
func RUNA1() {
	log.Printf("Starting Assignment 1")
	router := httprouter.New()                          // Create HTTP router
	router.GET("/", homepage)                           // Services index.html
	router.GET("/random", returnRandomBusinessListJson) // Handler for random Businesses list
	router.GET("/relatable", getRelatableBusinesses)    // Handler for Relatable Businesses
	readBusinessesJson()                                // Review Business JSON data and filter
	readReviewsJsonScannner()                           // Rview Review JSON and generate term count tables
	removeNullReviewsCalculateFrequency()               // REmove Businesses with no reviews and calculate term document frequency
	calculatetfIdf()                                    // Iterate through businesses and calculate td-idf for terms
	sortTfIdf()                                         // Sort tf-idf map within the Businneses
	addMostRelevantTermsKeyMap()                        // Add top x percent of relatable terms to global key map
	err := http.ListenAndServe(":7500", router)
	if err != nil {
		log.Fatal("Problem starting service: " + err.Error())
	}
	log.Println(serviceName + " running")
	log.Println("Finished")
}

// RUN2_1 Runtime for Assignment 2 program, the loader for the service
func RUN2_1() {
	_, BusinessesDP := ReadBusinessJSON2()
	var KmediodsDS KmediodsDS
	KmediodsDS.PopClusters(BusinessesDP, 10)
	log.Printf("Starting Assignment 2 Loader")
	//Map := NewHashMap22()
	//for _, c := range KmediodsDS.Clusters {
	//	for _, b := range c.Points {
	//		Map.Add(string(c.ID), b.BusinessID)
	//	}
	//}
}

// RUN2_2 Runtime for Assignment 2 program, the application for the service
func RUN2_2() {
	log.Printf("Starting Assignment 2 Application")
	router := httprouter.New()                          // Create HTTP router
	router.GET("/", homepage)                           // Services index.html
	router.GET("/random", returnRandomBusinessListJson) // Handler for random Businesses list
	router.GET("/relatable", getRelatableBusinesses)    // Handler for Relatable Businesses
	err := http.ListenAndServe(":7500", router)
	if err != nil {
		log.Fatal("Problem starting service: " + err.Error())
	}

}
