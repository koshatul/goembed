package main

import (
	"log"
	"net/http"

	"github.com/koshatul/goembed/examples/webserver/assets"
)

func main() {
	fileserver := http.FileServer(assets.Fs)
	http.Handle("/", http.StripPrefix("/", fileserver))
	log.Println("Listening on :8080")
	http.ListenAndServe(":8080", nil)
}
