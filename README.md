# goembed

[![Build Status](https://travis-ci.org/koshatul/goembed.svg?branch=master)](https://travis-ci.org/koshatul/goembed)

goembed takes a list of local files and embeds them into golang source files.

## Compression

Supported compression algorithms:
- deflate
- gzip
- lzw
- snappy
- snappy (stream)
- zlib
- none

## Wrappers

Supported wrappers:
- `nodep` No runtime dependencies, self-contained.
- `afero` using [spf13/afero](https://github.com/spf13/afero).

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

```
Usage:
  goembed [flags]
  goembed [command]

Available Commands:
  help        Help about any command
  version     Print the version

Flags:
  -c, --compression string   Compression to use, options are 'deflate', 'gzip', 'lzw', 'snappy', 'snappystream', 'zlib' or 'none' (default "snappy")
  -d, --debug                Debug output
  -f, --file string          Output file, or '-' for STDOUT (default "-")
  -h, --help                 help for goembed
  -p, --package string       golang package name for file (default: based on output file directory)
      --version              version for goembed
  -w, --wrapper string       Wrapper to use, options are 'none' or 'afero' (default "none")

Use "goembed [command] --help" for more information about a command.
```

## Examples

Examples are kept in the [examples](https://github.com/koshatul/goembed/tree/master/examples) directory, the usage above is the webserver example.

## Similar projects

- [statik](https://github.com/rakyll/statik)
- [packr](https://github.com/gobuffalo/packr)
