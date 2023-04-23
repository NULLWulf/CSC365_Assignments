package main

import (
	"github.com/julienschmidt/httprouter"
	"log"
	"net/http"
	"os"
)

// RUN2_1 Runtime for Assignment 2 program, the loader for the service
func RUN2_1() {
	ReadBusinessJSON2()
	k := KmediodsDS{}
	k.BuildFromPSD()
	k.BuildGraphFromKM()
	err := k.saveKMDStoDisc("kmed.bin")
	if err != nil {
		return
	}
	log.Printf("Kmediods data structure saved to disc")
	log.Printf("Program finished")
	os.Exit(0)
}

func RUN2_2() {
	//log.Printf("Starting Assignment 2 Application")
	//var kmed = &KmediodsDS{}
	//err := kmed.loadKMDStoDisc("kmed.bin")
	log.Printf("Kmediods data structure loaded from disc")
	graph, err := deserializeGraph()
	log.Printf("Graph data structure loaded from disc")
	log.Printf("Kmediods data structure loaded from disc")
	if err != nil {
		log.Fatal("Problem loading kmediods data structure: " + err.Error())
	}
	log.Printf("Starting Assignment 3 Application")
	router := httprouter.New() // Create HTTP router
	router.GET("/", func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		homepage2(w, r)
	}) // Services index.html
	router.GET("/random", func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		returnRandomBusinessListJsonFromGraph(w, r, graph)
	}) // Handler for random Businesses list
	router.GET("/random", func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		performDjstrika(w, r, graph)
	}) // Handler for random Businesses list

	log.Printf("Listening on port 7500")
	err = http.ListenAndServe(":7500", router)
	if err != nil {
		log.Fatal("Problem starting service: " + err.Error())
	}
}
