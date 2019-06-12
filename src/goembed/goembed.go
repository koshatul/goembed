package goembed

import (
	"io"
	"os"
)

// File is a wrapper for io.Reader with extra metadata for use while wrapping.
type File struct {
	io.Reader

	Name string
	Stat os.FileInfo
}

// NewFile takes the metadata and reader and returns a File object.
func NewFile(filename string, stat os.FileInfo, reader io.Reader) File {
	return File{
		Reader: reader,
		Name:   filename,
		Stat:   stat,
	}
}
