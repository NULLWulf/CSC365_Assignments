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
	BusinessesDPS := ReadBusinessJSON2()
	var KmediodsDSA KmediodsDS
	log.Printf("Starting Kemediods Algorithm")
	KmediodsDSA.PopClusters(BusinessesDPS, 10)
	log.Printf("Kemediods Algorithm Complete")
	log.Printf("Saving KMediods to file")
	err := KmediodsDSA.saveKMDStoDisc("kmed.bin")
	if err != nil {
		log.Fatal("Problem saving KMediods to disc: " + err.Error())
	}
	// clear km
	KmediodsDSA = KmediodsDS{}
	log.Printf("Loading KMediods from file")
	err = KmediodsDSA.loadKMDStoDisc("kmed.bin")
	if err != nil {
		log.Fatal("Problem loading KMediods from disc: " + err.Error())
	}
	log.Printf("KMediods loaded from file and verified %d", len(KmediodsDSA.Clusters))
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
