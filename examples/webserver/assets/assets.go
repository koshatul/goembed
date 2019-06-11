// Code generated - DO NOT EDIT.

package assets

import (
	"bytes"
	"net/http"
	"os"
	"time"
)

type assetFileData struct {
	name     string
	data     []byte
	dir      bool
	children []*assetFileData
}
type Fs struct{}

func (a Fs) Open(name string) (http.File, error) {
	switch name {
	case "/s1":
		return &assetFile{Reader: bytes.NewReader(dirL3Mx.data), assetFileData: dirL3Mx}, nil
	case "/s1/s2/index.html":
		return &assetFile{Reader: bytes.NewReader(fileL3MxL3MyL2luZGV4Lmh0bWw.data), assetFileData: fileL3MxL3MyL2luZGV4Lmh0bWw}, nil
	case "/":
		return &assetFile{Reader: bytes.NewReader(dirLw.data), assetFileData: dirLw}, nil
	case "/index.html":
		return &assetFile{Reader: bytes.NewReader(fileL2luZGV4Lmh0bWw.data), assetFileData: fileL2luZGV4Lmh0bWw}, nil
	case "/s1/s2":
		return &assetFile{Reader: bytes.NewReader(dirL3MxL3My.data), assetFileData: dirL3MxL3My}, nil
	}
	return nil, os.ErrNotExist
}

type assetFileInfo struct {
	f *assetFile
}

func (a assetFileInfo) Name() string {
	return a.f.name
}
func (a assetFileInfo) Size() int64 {
	return int64(len(a.f.data))
}
func (a assetFileInfo) Mode() os.FileMode {
	return 292
}
func (a assetFileInfo) ModTime() time.Time {
	return time.Time{}
}
func (a assetFileInfo) IsDir() bool {
	return a.f.dir
}
func (a assetFileInfo) Sys() interface{} {
	return nil
}

type assetFile struct {
	*bytes.Reader
	*assetFileData
}

func (a *assetFile) Stat() (os.FileInfo, error) {
	return assetFileInfo{f: a}, nil
}
func (a *assetFile) Readdir(count int) ([]os.FileInfo, error) {
	if a.dir {
		fl := []os.FileInfo{}
		for _, c := range a.children {
			d := &assetFile{assetFileData: c}
			fl = append(fl, &assetFileInfo{f: d})
		}
		return fl, nil
	}
	return nil, nil
}
func (a *assetFile) Close() error {
	return nil
}

var dirLw *assetFileData = &assetFileData{name: "/", dir: true, children: []*assetFileData{dirL3Mx, fileL2luZGV4Lmh0bWw}}
var fileL2luZGV4Lmh0bWw *assetFileData = &assetFileData{name: "/index.html", dir: false, data: []byte{60, 104, 116, 109, 108, 62, 10, 32, 32, 32, 32, 60, 98, 111, 100, 121, 62, 10, 32, 32, 32, 32, 32, 32, 32, 32, 84, 101, 115, 116, 32, 70, 105, 108, 101, 10, 32, 32, 32, 32, 60, 47, 98, 111, 100, 121, 62, 10, 60, 47, 104, 116, 109, 108, 62, 10}}
var dirL3MxL3My *assetFileData = &assetFileData{name: "/s1/s2", dir: true, children: []*assetFileData{}}
var dirL3Mx *assetFileData = &assetFileData{name: "/s1", dir: true, children: []*assetFileData{}}
var fileL3MxL3MyL2luZGV4Lmh0bWw *assetFileData = &assetFileData{name: "/s1/s2/index.html", dir: false, data: []byte{60, 104, 116, 109, 108, 62, 10, 32, 32, 32, 32, 60, 98, 111, 100, 121, 62, 10, 32, 32, 32, 32, 32, 32, 32, 32, 84, 101, 115, 116, 32, 70, 105, 108, 101, 10, 32, 32, 32, 32, 60, 47, 98, 111, 100, 121, 62, 10, 60, 47, 104, 116, 109, 108, 62, 10}}
