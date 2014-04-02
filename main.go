package main

import (
	"./pixelwrapper"
	"./reservationservice"
	"github.com/emicklei/go-restful"
	"log"
	"net/http"
)

func main() {
	restful.Add(reservationservice.New())

	http.HandleFunc("/pixel", pixelwrapper.PixelHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))

}
