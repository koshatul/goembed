package shrink

import (
	"bytes"
	"compress/zlib"
	"io"

	"github.com/dave/jennifer/jen"
	"github.com/koshatul/goembed/src/goembed"
	"github.com/sirupsen/logrus"
)

// ZlibStreamShrinker is a Shrinker compatible struct that uses zlib compression
type ZlibStreamShrinker struct {
}

// NewZlibStreamShrinker returns a Shrinker compatible class that uses zlib compression
func NewZlibStreamShrinker() Shrinker {
	return &ZlibStreamShrinker{}
}

// Name returns a simple name for this module
func (b *ZlibStreamShrinker) Name() string {
	return "zlib"
}

// IsStream returns true if the shrinker works on streams instead of byte slices
func (b *ZlibStreamShrinker) IsStream() bool {
	return true
}

// IsReaderWithError returns true if the shrinker uses a reader that also can error.
func (b *ZlibStreamShrinker) IsReaderWithError() bool {
	return true
}

// Compress returns a byte array of compressed file data
func (b *ZlibStreamShrinker) Compress(file goembed.File) ([]jen.Code, error) {
	v := []jen.Code{}

	cmpOut := new(bytes.Buffer)
	rawIn := zlib.NewWriter(cmpOut)
	n, err := io.Copy(rawIn, file)
	rawIn.Close()
	if err != nil {
		return nil, err
	}
	logrus.WithField("compression", "zlib").Debugf("Copied %d bytes into compressor", n)

	buf := make([]byte, 1)
	for {
		_, err := cmpOut.Read(buf)
		if err != nil {
			break
		}
		v = append(v, jen.Lit(int(buf[0])))
	}

	logrus.WithField("compression", "zlib").Debugf("Wrote %d bytes to static asset", len(v))

	return v, nil
}

// Decompressor returns the body code for the `decode(input)` function
func (b *ZlibStreamShrinker) Decompressor() []jen.Code {
	return nil // Not used for streams.
}

// Reader returns the stream handler for the byte stream used when returning `Open()`
func (b *ZlibStreamShrinker) Reader(params ...jen.Code) jen.Code {
	return nil
}

// ReaderWithError returns the stream handler for the byte stream used when returning `Open()` but also returns an error
func (b *ZlibStreamShrinker) ReaderWithError(params ...jen.Code) jen.Code {
	return jen.Id("reader").Call(jen.Qual("bytes", "NewReader").Params(params...))
}

// Header returns additional code that is inserted in the body
func (b *ZlibStreamShrinker) Header() []jen.Code {
	return []jen.Code{
		jen.Func().Id("reader").Params(jen.Id("i").Qual("io", "Reader")).Params(jen.Qual("io", "Reader"), jen.Error()).Block(
			jen.Return(jen.Qual("compress/zlib", "NewReader").Params(jen.Id("i"))),
		),
	}
}
