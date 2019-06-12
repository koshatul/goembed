package main

import (
	"log"
	"net/http"

	"github.com/koshatul/goembed/examples/webserver/assets"
	"github.com/spf13/afero"
)

func main() {
	var fileserver http.Handler
	if v, ok := assets.Fs.(http.FileSystem); ok {
		fileserver = http.FileServer(v)
	} else if v, ok := assets.Fs.(afero.Fs); ok {
		httpFs := afero.NewHttpFs(v)
		fileserver = http.FileServer(httpFs)
	}
	http.Handle("/", http.StripPrefix("/", fileserver))
	log.Println("Listening on :8080")
	http.ListenAndServe(":8080", nil)
}
