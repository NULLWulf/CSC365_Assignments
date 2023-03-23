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
	ReadBusinessJSON2()
	// var KmediodsDS KmediodsDS
	// log.Printf("Starting Kemediods Algorithm")
	// KmediodsDS.PopClusters(BusinessesDP, 10)
	// // Save KMedios to file
	// log.Printf("Saving KMedios to file")

	// log.Printf("Starting Assignment 2 Loader")

	// Hash := NewHashMap()
	// // Put clusters into HashMap
	// log.Printf("Putting clusters into HashMap")
	// for _, cluster := range KmediodsDS.Clusters {
	// 	for _, business := range cluster.Points {
	// 		Hash.Add(cluster.Medoid.BusinessID, business.BusinessID)
	// 	}
	// }

	// log.Printf("Saving HashMap to file")
	// err := Hash.SaveToFile("hashmap.json")
	// if err != nil {
	// 	return
	// }
	// Hash2, err := LoadHashMapFromFile("hashmap.json")
	// if err != nil {
	// 	log.Fatal("Problem loading hashmap: " + err.Error())
	// }

	// log.Printf("HashMap loaded from file and verified %d", Hash2.Size)
	// log.Printf("Finished Assignment 2 Loader Finished")
}

// RUN2_2 Runtime for Assignment 2 program, the application for the service
func RUN2_2() {
	log.Printf("Starting Assignment 2 Application")
	router := httprouter.New()                          // Create HTTP router
	router.GET("/", homepage)                           // Services index.html
	router.GET("/random", returnRandomBusinessListJson) // Handler for random Businesses list
	router.GET("/clustered", getRelatableCluster)       // Handler for Relatable Businesses
	err := http.ListenAndServe(":7500", router)
	if err != nil {
		log.Fatal("Problem starting service: " + err.Error())
	}

}
