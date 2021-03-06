package shrink

import (
	"github.com/dave/jennifer/jen"
	"github.com/koshatul/goembed/goembed"
)

// Shrinker is the interface that provides a compression method for the data.
type Shrinker interface {
	Name() string
	Compress(file goembed.File) ([]jen.Code, error)
	Header() []jen.Code
	Decompressor() []jen.Code
	IsStream() bool
	IsReaderWithError() bool
	Reader(params ...jen.Code) jen.Code
	ReaderWithError(params ...jen.Code) jen.Code
}
