package wrap

import (
	"io"

	"github.com/koshatul/goembed/src/goembed"
)

// Wrapper is the interface that provides a method for handling the data.
type Wrapper interface {
	// AddFile(filename string, file io.Reader) error
	// AddDir(dir string) error
	AddFile(filename string, file goembed.File) error
	Render(w io.Writer) error
}
