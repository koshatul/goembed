package shrink

import (
	"bytes"
	"io/ioutil"

	"github.com/dave/jennifer/jen"
	"github.com/golang/snappy"
	"github.com/koshatul/goembed/src/goembed"
	"github.com/sirupsen/logrus"
)

// SnappyShrinker is a Shrinker compatible struct that uses snappy compression
type SnappyShrinker struct {
}

// NewSnappyShrinker returns a Shrinker compatible class that uses snappy compression
func NewSnappyShrinker() Shrinker {
	return &SnappyShrinker{}
}

// Name returns a simple name for this module
func (b *SnappyShrinker) Name() string {
	return "snappy"
}

// IsStream returns true if the shrinker works on streams instead of byte slices
func (b *SnappyShrinker) IsStream() bool {
	return false
}

// IsReaderWithError returns true if the shrinker uses a reader that also can error.
func (b *SnappyShrinker) IsReaderWithError() bool {
	return false
}

// Compress returns a byte array of compressed file data
func (b *SnappyShrinker) Compress(file goembed.File) ([]jen.Code, error) {
	v := []jen.Code{}

	src, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, err
	}
	encoded := snappy.Encode(nil, src)
	logrus.WithField("compression", "snappy").Debugf("Copied %d bytes into compressor", len(src))
	cmpOut := bytes.NewBuffer(encoded)

	buf := make([]byte, 1)
	for {
		_, err := cmpOut.Read(buf)
		if err != nil {
			break
		}
		v = append(v, jen.Lit(int(buf[0])))
	}

	logrus.WithField("compression", "snappy").Debugf("Wrote %d bytes to static asset", len(v))

	return v, nil
}

// Decompressor returns the body code for the `decode(input)` function
func (b *SnappyShrinker) Decompressor() []jen.Code {
	return []jen.Code{
		jen.List(jen.Id("o"), jen.Id("_")).Op(":=").Qual("github.com/golang/snappy", "Decode").Call(jen.Nil(), jen.Id("input")),
		jen.Return(jen.Id("o")),
	}
}

// Reader returns the stream handler for the byte stream used when returning `Open()`
func (b *SnappyShrinker) Reader(params ...jen.Code) jen.Code {
	return jen.Qual("bytes", "NewReader").Params(jen.Id("decode").Call(params...))
}

// ReaderWithError returns the stream handler for the byte stream used when returning `Open()` but also returns an error
func (b *SnappyShrinker) ReaderWithError(params ...jen.Code) jen.Code {
	return nil
}

// Header returns additional code that is inserted in the body
func (b *SnappyShrinker) Header() []jen.Code {
	return []jen.Code{}
}
