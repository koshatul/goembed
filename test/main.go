package main

import (
	"fmt"
	"log"

	"github.com/spf13/afero"
)

func main() {
	data, err := afero.ReadFile(Fs, "/index.html")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%s", data)
}
