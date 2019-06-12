# goembed

[![Build Status](https://travis-ci.org/koshatul/goembed.svg?branch=master)](https://travis-ci.org/koshatul/goembed)

goembed takes a list of local files and embeds them into golang source files.

## Usage

Install the command line tool first.

	go get github.com/koshatul/goembed/src/...

goembed has sane defaults and will safely generate a go package from a directory or individual list of files.

    $ goembed -f src/webasset/assets.go ./web

Then in your application just import the generated package

If you prefer to use a pseudo-filesystem there is support for using [afero](https://github.com/spf13/afero) as a wrapper.

~~~ go
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
~~~

Then the filesystem should be visible from http://localhost:8080/

## Examples

Examples are kept in the [examples](https://github.com/koshatul/goembed/tree/master/examples) directory, the usage above is the webserver example.

## Similar projects

- [statik](https://github.com/rakyll/statik)
- [packr](https://github.com/gobuffalo/packr)
