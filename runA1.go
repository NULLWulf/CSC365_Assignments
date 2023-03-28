package main

import (
	"github.com/julienschmidt/httprouter"
	"log"
	"net/http"
	"os"
)

// RUNA1 Runtime for Assignment 1 program
//func RUNA1() {
//	log.Printf("Starting Assignment 1")
//	router := httprouter.New()                          // Create HTTP router
//	router.GET("/", homepage())                           // Services index.html
//	router.GET("/random", returnRandomBusinessListJson) // Handler for random Businesses list
//	router.GET("/relatable", getRelatableBusinesses)    // Handler for Relatable Businesses
//	readBusinessesJson()                                // Review Business JSON data and filter
//	readReviewsJsonScannner()                           // Rview Review JSON and generate term count tables
//	removeNullReviewsCalculateFrequency()               // REmove Businesses with no reviews and calculate term document frequency
//	calculatetfIdf()                                    // Iterate through businesses and calculate td-idf for terms
//	sortTfIdf()                                         // Sort tf-idf map within the Businneses
//	addMostRelevantTermsKeyMap()                        // Add top x percent of relatable terms to global key map
//	err := http.ListenAndServe(":7500", router)
//	if err != nil {
//		log.Fatal("Problem starting service: " + err.Error())
//	}
//	log.Println(serviceName + " running")
//	log.Println("Finished")
//}

// RUN2_1 Runtime for Assignment 2 program, the loader for the service
func RUN2_1() {
	ReadBusinessJSON2()
	k := KmediodsDS{}
	k.BuildFromPSD()
	err := k.saveKMDStoDisc("kmed.bin")
	if err != nil {
		return
	}
	log.Printf("Kmediods data structure saved to disc")
	log.Printf("Program finished")
	os.Exit(0)
}

func RUN2_2() {
	var kmed = &KmediodsDS{}
	err := kmed.loadKMDStoDisc("kmed.bin")
	log.Printf("Kmediods data structure loaded from disc")
	if err != nil {
		log.Fatal("Problem loading kmediods data structure: " + err.Error())
	}
	log.Printf("Starting Assignment 2 Application")
	router := httprouter.New() // Create HTTP router
	router.GET("/", func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		homepage2(w, r)
	}) // Services index.html
	router.GET("/random", func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		returnRandomBusinessListJson(w, r, kmed)
	}) // Handler for random Businesses list
	router.GET("/clustered", func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		getRelatableCluster(w, r, kmed)
	}) // Handler for Relatable Businesses
	log.Printf("Listening on port 7500")
	err = http.ListenAndServe(":7500", router)
	if err != nil {
		log.Fatal("Problem starting service: " + err.Error())
	}
}
