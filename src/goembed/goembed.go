package goembed

import (
	"io"
	"os"
)

type File struct {
	io.Reader

	Name string
	Stat os.FileInfo
}

func NewFile(filename string, stat os.FileInfo, reader io.Reader) File {
	return File{
		Reader: reader,
		Name:   filename,
		Stat:   stat,
	}
}
