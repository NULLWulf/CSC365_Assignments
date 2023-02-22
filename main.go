package main

import (
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/kardianos/service"
)

type program struct{}

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

func (p *program) run() {
	readBusinessesJson()
	readReviewsJsonScannner()
	removeNullReviewsCalculateFrequency()
	calculatetfIdf()
	sortTfIdf()
	a := addMostRelevantTermsKeyMap()
	log.Printf("%v", len(a))
	// log.Print(getRandomBusiness(10))
	// log.Println(getRandomBusiness(5))
	log.Printf("%v", findRelatableBusinesses("gqMRM12j5mder76t3bPxJw"))

	router := httprouter.New()
	router.ServeFiles("/js/*filepath", http.Dir("js"))
	router.ServeFiles("/css/*filepath", http.Dir("css"))
	router.GET("/", homepage)
	//router.GET('searc/:query', search')
	err := http.ListenAndServe(":8080", router)
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
