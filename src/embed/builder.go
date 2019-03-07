package embed

import "io"

// Builder is the interface that provides a Builder for the data.
type Builder interface {
	AddFile(filename string, file io.Reader) error
	Render(w io.Writer) error
}
