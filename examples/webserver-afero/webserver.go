package main

import (
	"log"
	"net/http"

	"github.com/koshatul/goembed/examples/webserver-afero/assets"
	"github.com/spf13/afero"
)

func main() {
	httpFs := afero.NewHttpFs(assets.Fs)
	fileserver := http.FileServer(httpFs)
	http.Handle("/", http.StripPrefix("/", fileserver))
	log.Println("Listening on :8080")
	http.ListenAndServe(":8080", nil)
}
