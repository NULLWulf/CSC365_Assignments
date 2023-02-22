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

func bizSearch(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	term := request.URL.Query().Get("q")
	fmt.Println("Serving business search: " + term)
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
