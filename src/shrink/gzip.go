package shrink

import (
	"bytes"
	"compress/gzip"
	"io"

	"github.com/dave/jennifer/jen"
	"github.com/koshatul/goembed/src/goembed"
	"github.com/sirupsen/logrus"
)

// GzipStreamShrinker is a Shrinker compatible struct that uses gzip compression
type GzipStreamShrinker struct {
}

// NewGzipStreamShrinker returns a Shrinker compatible class that uses gzip compression
func NewGzipStreamShrinker() Shrinker {
	return &GzipStreamShrinker{}
}

// Name returns a simple name for this module
func (b *GzipStreamShrinker) Name() string {
	return "gzip"
}

// IsStream returns true if the shrinker works on streams instead of byte slices
func (b *GzipStreamShrinker) IsStream() bool {
	return true
}

// IsReaderWithError returns true if the shrinker uses a reader that also can error.
func (b *GzipStreamShrinker) IsReaderWithError() bool {
	return true
}

// Compress returns a byte array of compressed file data
func (b *GzipStreamShrinker) Compress(file goembed.File) ([]jen.Code, error) {
	v := []jen.Code{}

	cmpOut := new(bytes.Buffer)
	rawIn := gzip.NewWriter(cmpOut)
	n, err := io.Copy(rawIn, file)
	rawIn.Close()
	if err != nil {
		return nil, err
	}
	logrus.WithField("compression", "gzip").Debugf("Copied %d bytes into compressor", n)

	buf := make([]byte, 1)
	for {
		_, err := cmpOut.Read(buf)
		if err != nil {
			break
		}
		v = append(v, jen.Lit(int(buf[0])))
	}

	logrus.WithField("compression", "gzip").Debugf("Wrote %d bytes to static asset", len(v))

	return v, nil
}

// Decompressor returns the body code for the `decode(input)` function
func (b *GzipStreamShrinker) Decompressor() []jen.Code {
	return nil // Not used for streams.
}

// Reader returns the stream handler for the byte stream used when returning `Open()`
func (b *GzipStreamShrinker) Reader(params ...jen.Code) jen.Code {
	return nil
}

// ReaderWithError returns the stream handler for the byte stream used when returning `Open()`
func (b *GzipStreamShrinker) ReaderWithError(params ...jen.Code) jen.Code {
	return jen.Id("reader").Call(jen.Qual("bytes", "NewReader").Params(params...))
}

// Header returns additional code that is inserted in the body
func (b *GzipStreamShrinker) Header() []jen.Code {
	return []jen.Code{
		jen.Func().Id("reader").Params(jen.Id("i").Qual("io", "Reader")).Params(jen.Qual("io", "Reader"), jen.Error()).Block(
			jen.Return(jen.Qual("compress/gzip", "NewReader").Params(jen.Id("i"))),
		),
	}
}
