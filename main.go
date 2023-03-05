package main

import (
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/kardianos/service"
)

// Program instantation
func main() {
	log.Println(serviceName + " starting...")
	serviceConfig := &service.Config{
		Name:        serviceName,
		DisplayName: serviceName,
		Description: serviceDescription,
	}
	prg := &program{}
	s, err := service.New(prg, serviceConfig)
	if err != nil {
		log.Println("Cannot start: " + err.Error())
	}
	err = s.Run()
	if err != nil {
		log.Println("Cannot start: " + err.Error())
	}
}

// Main program loop
func (p *program) run() {
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

func (p *program) Start(service.Service) error {
	log.Println(serviceName + " started")
	go p.run()
	return nil
}

func (p *program) Stop(service.Service) error {
	log.Println(serviceName + " stopped")
	return nil
}



