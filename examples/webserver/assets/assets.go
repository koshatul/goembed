package assets

import (
	"os"
	"time"
	"bytes"
	"net/http"
	"log"
)

type assetFileData struct {
	name string
	data []byte
	dir bool
	children []*assetFileData
}

var fileabc *assetFileData = &assetFileData{
	name: "/index.html", 
	data: []byte{60, 104, 116, 109, 108, 62, 10, 32, 32, 32, 32, 60, 98, 111, 100, 121, 62, 10, 32, 32, 32, 32, 32, 32, 32, 32, 84, 101, 115, 116, 32, 70, 105, 108, 101, 10, 32, 32, 32, 32, 60, 47, 98, 111, 100, 121, 62, 10, 60, 47, 104, 116, 109, 108, 62, 10},
	dir: false,
}

var dirabc *assetFileData = &assetFileData{
	name: "/",
	dir: true,
	children: []*assetFileData{
		fileabc,
	},
}

var fileData = map[string]*assetFileData{
	"/": dirabc,
	"/index.html": fileabc,
}

type FS struct {
}

func (a FS) Open(name string) (http.File, error) {
	log.Printf("Open(): %s", name)
	if v, ok := fileData[name]; ok {
		log.Printf("Returning file: %s", name)
		return &assetFile{
			Reader: bytes.NewReader(v.data),
			assetFileData: v,
		}, nil
	}

	log.Printf("File not found: %s", name)
	return nil, os.ErrNotExist
}

type assetFileInfo struct {
	f *assetFile
}

func (a assetFileInfo) Name() string       {
	log.Printf("[%s] ModTime()", a.f.name)
	return a.f.name
}
func (a assetFileInfo) Size() int64        {
	log.Printf("[%s] ModTime()", a.f.name)
	return int64(len(a.f.data))
}
func (a assetFileInfo) Mode() os.FileMode  {        // Read for all
	log.Printf("[%s] ModTime()", a.f.name)
	return 0444
}
func (a assetFileInfo) ModTime() time.Time { // Return anything
	log.Printf("[%s] ModTime()", a.f.name)
	return time.Time{} 
}
func (a assetFileInfo) IsDir() bool        { 
	log.Printf("[%s] IsDir()", a.f.name)
	return a.f.dir
}
func (a assetFileInfo) Sys() interface{}   {
	log.Printf("[%s] ModTime()", a.f.name)
	return nil
}

type assetFile struct {
	*bytes.Reader
	*assetFileData
}

func (a *assetFile) Stat() (os.FileInfo, error) {
	log.Printf("[%s] Stat()", a.name)
    return assetFileInfo{f: a}, nil
}

func (a *assetFile) Readdir(count int) ([]os.FileInfo, error) {
	log.Printf("[%s] Readdir()", a.name)
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
	log.Printf("[%s] Close()", a.name)
	return nil
}
