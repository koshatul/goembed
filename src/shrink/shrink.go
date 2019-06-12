package shrink

import (
	"github.com/dave/jennifer/jen"
	"github.com/koshatul/goembed/src/goembed"
)

// Shrinker is the interface that provides a compression method for the data.
type Shrinker interface {
	Compress(file goembed.File) ([]jen.Code, error)
	Header() []jen.Code
	Decompressor() []jen.Code
	// AddDir(dir string) error
	// AddFile(filename string, file io.Reader) error
	// Render(w io.Writer) error
}
