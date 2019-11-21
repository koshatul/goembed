package shrink

import (
	"bytes"
	"compress/flate"
	"io"

	"github.com/dave/jennifer/jen"
	"github.com/koshatul/goembed/goembed"
	"github.com/sirupsen/logrus"
)

// DeflateStreamShrinker is a Shrinker compatible struct that uses deflate compression
type DeflateStreamShrinker struct {
}

// NewDeflateStreamShrinker returns a Shrinker compatible class that uses deflate compression
func NewDeflateStreamShrinker() Shrinker {
	return &DeflateStreamShrinker{}
}

// Name returns a simple name for this module
func (b *DeflateStreamShrinker) Name() string {
	return "deflate"
}

// IsStream returns true if the shrinker works on streams instead of byte slices
func (b *DeflateStreamShrinker) IsStream() bool {
	return true
}

// IsReaderWithError returns true if the shrinker uses a reader that also can error.
func (b *DeflateStreamShrinker) IsReaderWithError() bool {
	return false
}

// Compress returns a byte array of compressed file data
func (b *DeflateStreamShrinker) Compress(file goembed.File) ([]jen.Code, error) {
	v := []jen.Code{}

	cmpOut := new(bytes.Buffer)
	rawIn, err := flate.NewWriter(cmpOut, -1)
	if err != nil {
		return nil, err
	}
	n, err := io.Copy(rawIn, file)
	rawIn.Close()
	if err != nil {
		return nil, err
	}
	logrus.WithField("compression", "deflate").Debugf("Copied %d bytes into compressor", n)

	buf := make([]byte, 1)
	for {
		_, err := cmpOut.Read(buf)
		if err != nil {
			break
		}
		v = append(v, jen.Lit(int(buf[0])))
	}

	logrus.WithField("compression", "deflate").Debugf("Wrote %d bytes to static asset", len(v))

	return v, nil
}

// Decompressor returns the body code for the `decode(input)` function
func (b *DeflateStreamShrinker) Decompressor() []jen.Code {
	return nil // Not used for streams.
}

// Reader returns the stream handler for the byte stream used when returning `Open()`
func (b *DeflateStreamShrinker) Reader(params ...jen.Code) jen.Code {
	return jen.Qual("compress/flate", "NewReader").Call(jen.Qual("bytes", "NewReader").Params(params...))
}

// ReaderWithError returns the stream handler for the byte stream used when returning `Open()` but also returns an error
func (b *DeflateStreamShrinker) ReaderWithError(params ...jen.Code) jen.Code {
	return nil
}

// Header returns additional code that is inserted in the body
func (b *DeflateStreamShrinker) Header() []jen.Code {
	return []jen.Code{}
}
