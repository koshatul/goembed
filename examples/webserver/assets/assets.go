// Code generated - DO NOT EDIT.

package assets

import (
	"bytes"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

type assetFileData struct {
	name string
	data []byte
	dir  bool
}

func (a *assetFileData) Children() []*assetFileData {
	o := []*assetFileData{}
	for f, v := range fileData {
		if !strings.EqualFold(a.name, f) && strings.HasPrefix(a.name, f) {
			log.Printf("f:'%s' a.name:'%s'", f, a.name)
			log.Printf("f(%d:%d)", len(a.name), len(f))
			ft := f[len(a.name):len(f)]
			if !strings.Contains(ft, "/") {
				o = append(o, v)
			}
		}
	}
	return o
}

type Fs struct{}

func (a Fs) Open(name string) (http.File, error) {
	if v, ok := fileData[name]; ok {
		return &assetFile{Reader: bytes.NewReader(v.data), assetFileData: v}, nil
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
		for _, c := range a.Children() {
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

var dirLw *assetFileData = &assetFileData{name: "/", dir: true}
var fileL2luZGV4Lmh0bWw *assetFileData = &assetFileData{name: "/index.html", dir: false, data: []byte{60, 104, 116, 109, 108, 62, 10, 32, 32, 32, 32, 60, 98, 111, 100, 121, 62, 10, 32, 32, 32, 32, 32, 32, 32, 32, 84, 101, 115, 116, 32, 70, 105, 108, 101, 10, 32, 32, 32, 32, 60, 47, 98, 111, 100, 121, 62, 10, 60, 47, 104, 116, 109, 108, 62, 10}}
var dirL3MxL3My *assetFileData = &assetFileData{name: "/s1/s2", dir: true}
var dirL3Mx *assetFileData = &assetFileData{name: "/s1", dir: true}
var fileL3MxL3MyL2luZGV4Lmh0bWw *assetFileData = &assetFileData{name: "/s1/s2/index.html", dir: false, data: []byte{60, 104, 116, 109, 108, 62, 10, 32, 32, 32, 32, 60, 98, 111, 100, 121, 62, 10, 32, 32, 32, 32, 32, 32, 32, 32, 84, 101, 115, 116, 32, 70, 105, 108, 101, 10, 32, 32, 32, 32, 60, 47, 98, 111, 100, 121, 62, 10, 60, 47, 104, 116, 109, 108, 62, 10}}
var fileData = map[string]*assetFileData{"/s1": dirL3Mx, "/s1/s2/index.html": fileL3MxL3MyL2luZGV4Lmh0bWw, "/": dirLw, "/index.html": fileL2luZGV4Lmh0bWw, "/s1/s2": dirL3MxL3My}
