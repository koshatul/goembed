// Code generated - DO NOT EDIT.

package assets

import (
	snappy "github.com/golang/snappy"
	afero "github.com/spf13/afero"
	"os"
)

// Fs is the filesystem containing the assets embedded in this package.
var Fs afero.Fs

func decode(input []byte) []byte {
	o, _ := snappy.Decode(nil, input)
	return o
}

var fileL2luZGV4Lmh0bWw = []byte{56, 60, 60, 104, 116, 109, 108, 62, 10, 32, 32, 32, 32, 60, 98, 111, 100, 121, 9, 11, 1, 1, 32, 84, 101, 115, 116, 32, 70, 105, 108, 101, 9, 29, 56, 47, 98, 111, 100, 121, 62, 10, 60, 47, 104, 116, 109, 108, 62, 10}
var fileL3MxL3MyL2luZGV4Lmh0bWw = []byte{56, 60, 60, 104, 116, 109, 108, 62, 10, 32, 32, 32, 32, 60, 98, 111, 100, 121, 9, 11, 1, 1, 32, 84, 101, 115, 116, 32, 70, 105, 108, 101, 9, 29, 56, 47, 98, 111, 100, 121, 62, 10, 60, 47, 104, 116, 109, 108, 62, 10}

func init() {
	Fs = afero.NewMemMapFs()
	afero.WriteFile(Fs, "/index.html", decode(fileL2luZGV4Lmh0bWw), os.ModePerm)
	afero.WriteFile(Fs, "/s1/s2/index.html", decode(fileL3MxL3MyL2luZGV4Lmh0bWw), os.ModePerm)
}
