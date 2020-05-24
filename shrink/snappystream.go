package shrink

import (
	"bytes"
	"io"

	"github.com/dave/jennifer/jen"
	"github.com/golang/snappy"
	"github.com/koshatul/goembed/goembed"
	"github.com/sirupsen/logrus"
)

// SnappyStreamShrinker is a Shrinker compatible struct that uses snappy compression.
type SnappyStreamShrinker struct {
}

// NewSnappyStreamShrinker returns a Shrinker compatible class that uses snappy compression.
func NewSnappyStreamShrinker() Shrinker {
	return &SnappyStreamShrinker{}
}

// Name returns a simple name for this module.
func (b *SnappyStreamShrinker) Name() string {
	return "snappy-stream"
}

// IsStream returns true if the shrinker works on streams instead of byte slices.
func (b *SnappyStreamShrinker) IsStream() bool {
	return true
}

// IsReaderWithError returns true if the shrinker uses a reader that also can error.
func (b *SnappyStreamShrinker) IsReaderWithError() bool {
	return false
}

// Compress returns a byte array of compressed file data.
func (b *SnappyStreamShrinker) Compress(file goembed.File) ([]jen.Code, error) {
	v := []jen.Code{}

	bufOut := &bytes.Buffer{}
	cmpOut := snappy.NewBufferedWriter(bufOut)
	w, err := io.Copy(cmpOut, file)

	_ = cmpOut.Close()

	if err != nil {
		return nil, err
	}

	logrus.WithField("compression", "snappy-stream").Debugf("Copied %d bytes into compressor", w)

	buf := make([]byte, 1)

	for {
		if _, err := bufOut.Read(buf); err != nil {
			break
		}

		v = append(v, jen.Lit(int(buf[0])))
	}

	logrus.WithField("compression", "snappy-stream").Debugf("Wrote %d bytes to static asset", len(v))

	return v, nil
}

// Decompressor returns the body code for the `decode(input)` function.
func (b *SnappyStreamShrinker) Decompressor() []jen.Code {
	return nil // Not used for streams.
}

// Reader returns the stream handler for the byte stream used when returning `Open()`.
func (b *SnappyStreamShrinker) Reader(params ...jen.Code) jen.Code {
	return jen.Qual("github.com/golang/snappy", "NewReader").Call(jen.Qual("bytes", "NewReader").Params(params...))
}

// ReaderWithError returns the stream handler for the byte stream used when returning `Open()` but also returns an error.
func (b *SnappyStreamShrinker) ReaderWithError(params ...jen.Code) jen.Code {
	return nil
}

// Header returns additional code that is inserted in the body.
func (b *SnappyStreamShrinker) Header() []jen.Code {
	return []jen.Code{}
}
