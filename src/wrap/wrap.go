package wrap

import (
	"io"

	"github.com/koshatul/goembed/src/goembed"
)

// Wrapper is the interface that provides a method for handling the data.
type Wrapper interface {
	Name() string
	AddFile(filename string, file goembed.File) error
	Render(w io.Writer) error
}
