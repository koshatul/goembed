package shrink

import (
	"bytes"
	"compress/lzw"
	"io"

	"github.com/dave/jennifer/jen"
	"github.com/koshatul/goembed/src/goembed"
	"github.com/sirupsen/logrus"
)

// LzwStreamShrinker is a Shrinker compatible struct that uses lzw compression
type LzwStreamShrinker struct {
}

// NewLzwStreamShrinker returns a Shrinker compatible class that uses lzw compression
func NewLzwStreamShrinker() Shrinker {
	return &LzwStreamShrinker{}
}

// IsStream returns true if the shrinker works on streams instead of byte slices
func (b *LzwStreamShrinker) IsStream() bool {
	return true
}

// IsReaderWithError returns true if the shrinker uses a reader that also can error.
func (b *LzwStreamShrinker) IsReaderWithError() bool {
	return false
}

// Compress returns a byte array of compressed file data
func (b *LzwStreamShrinker) Compress(file goembed.File) ([]jen.Code, error) {
	v := []jen.Code{}

	cmpOut := new(bytes.Buffer)
	rawIn := lzw.NewWriter(cmpOut, lzw.LSB, 8)
	n, err := io.Copy(rawIn, file)
	rawIn.Close()
	if err != nil {
		return nil, err
	}
	logrus.WithField("compression", "lzw").Debugf("Copied %d bytes into compressor", n)

	buf := make([]byte, 1)
	for {
		_, err := cmpOut.Read(buf)
		if err != nil {
			break
		}
		v = append(v, jen.Lit(int(buf[0])))
	}

	logrus.WithField("compression", "lzw").Debugf("Wrote %d bytes to static asset", len(v))

	return v, nil
}

// Decompressor returns the body code for the `decode(input)` function
func (b *LzwStreamShrinker) Decompressor() []jen.Code {
	return nil // Not used for streams.
}

// Reader returns the stream handler for the byte stream used when returning `Open()`
func (b *LzwStreamShrinker) Reader(params ...jen.Code) jen.Code {
	return jen.Qual("compress/lzw", "NewReader").Call(jen.Qual("bytes", "NewReader").Params(params...), jen.Qual("compress/lzw", "LSB"), jen.Lit(8))
}

// ReaderWithError returns the stream handler for the byte stream used when returning `Open()` but also returns an error
func (b *LzwStreamShrinker) ReaderWithError(params ...jen.Code) jen.Code {
	return nil
}

// Header returns additional code that is inserted in the body
func (b *LzwStreamShrinker) Header() []jen.Code {
	return []jen.Code{}
}
