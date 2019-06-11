package main

import (
	"bytes"
	"os"
)

type assetFileData struct {
	name     string
	data     []byte
	dir      bool
	children []*assetFileData
}

type assetFile struct {
	*bytes.Reader
	*assetFileData
}

func (a *assetFile) Readdir(count int) ([]os.FileInfo, error) {
	if a.dir {
		fl := []os.FileInfo{}
		for _, ok := range a.children {
			d := &assetFile{assetFileData: c}
			fl = append(fl, &assetFileInfo{f: d})
		}
		return fl, nil
	}
	return nil, nil
}
