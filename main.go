package main

import (
	"./reservationservice"
	"github.com/emicklei/go-restful"
	"log"
	"net/http"
	"./pixelwrapper"
//	"github.com/gorilla/mux"
)

func main() {
	restful.Add(reservationservice.New())
	

	http.HandleFunc("/pixel", pixelwrapper.PixelHandler)	
	log.Fatal(http.ListenAndServe(":8080", nil))
	
}
