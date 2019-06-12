package shrink

import (
	"github.com/dave/jennifer/jen"
	"github.com/koshatul/goembed/src/goembed"
	"github.com/sirupsen/logrus"
)

// NoShrinker is a Shrinker compatible struct that uses no compression
type NoShrinker struct {
}

// NewNoShrinker returns a Shrinker compatible class that uses no compression
func NewNoShrinker() Shrinker {
	return &NoShrinker{}
}

// Name returns a simple name for this module
func (b *NoShrinker) Name() string {
	return "none"
}

// IsStream returns true if the shrinker works on streams instead of byte slices
func (b *NoShrinker) IsStream() bool {
	return false
}

// IsReaderWithError returns true if the shrinker uses a reader that also can error.
func (b *NoShrinker) IsReaderWithError() bool {
	return false
}

// Compress returns a byte array of compressed file data
func (b *NoShrinker) Compress(file goembed.File) ([]jen.Code, error) {
	v := []jen.Code{}

	buf := make([]byte, 1)
	for {
		_, err := file.Read(buf)
		if err != nil {
			break
		}
		v = append(v, jen.Lit(int(buf[0])))
	}

	logrus.WithField("compression", "none").Debugf("Wrote %d bytes to static asset", len(v))

	return v, nil
}

// Decompressor returns the body code for the `decode(input)` function
func (b *NoShrinker) Decompressor() []jen.Code {
	return nil // Not used when there's no compression.
}

// Reader returns the stream handler for the byte stream used when returning `Open()`
func (b *NoShrinker) Reader(params ...jen.Code) jen.Code {
	return jen.Qual("bytes", "NewReader").Params(params...)
}

// ReaderWithError returns the stream handler for the byte stream used when returning `Open()` but also returns an error
func (b *NoShrinker) ReaderWithError(params ...jen.Code) jen.Code {
	return nil
}

// Header returns additional code that is inserted in the body
func (b *NoShrinker) Header() []jen.Code {
	return []jen.Code{}
}
