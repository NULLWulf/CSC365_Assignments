package main

import (
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

// func returnRandomBusinessJson(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
// 	fmt.Println("Serving random business")
// 	_, err := writer.Write([]byte(Businesses[rand.Intn(len(Businesses)-1)].ToJson()))
// 	if err != nil {
// 		return
// 	}
// }
