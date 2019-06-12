package wrap

// // Wrapper is the interface that provides a method for handling the data.
// type Wrapper interface {
// 	// AddFile(filename string, file io.Reader) error
// 	// AddDir(dir string) error
// 	AddFile(filename string, file goembed.File) error
// 	Render(w io.Writer) error
// }

import (
	"encoding/base64"
	"fmt"
	"io"

	"github.com/dave/jennifer/jen"
	"github.com/koshatul/goembed/src/goembed"
	"github.com/koshatul/goembed/src/shrink"
)

// AferoWrapper is a Wrapper compatible struct that uses afero for the file system
type AferoWrapper struct {
	file     *jen.File
	files    map[string]string
	shrinker shrink.Shrinker
}

// NewAferoWrapper returns a Wrapper compatible class that uses afero for the file system
func NewAferoWrapper(packageName string, shrinker shrink.Shrinker) Wrapper {
	f := jen.NewFile(packageName)
	f.HeaderComment("Code generated - DO NOT EDIT.")
	f.Line()
	f.Comment("Fs is the filesystem containing the assets embedded in this package.").Line().Var().Id("Fs").Qual("github.com/spf13/afero", "Fs")

	f.Func().Id("decode").Params(jen.Id("input").Index().Byte()).Params(jen.Index().Byte()).Block(
		shrinker.Decompressor()...,
	)

	return &AferoWrapper{
		file:     f,
		files:    map[string]string{},
		shrinker: shrinker,
	}
}

// AddFile adds a file to the embedded package.
func (b *AferoWrapper) AddFile(filename string, file goembed.File) error {
	v, err := b.shrinker.Compress(file)
	if err != nil {
		return err
	}

	// logrus.WithField("compression", "snappy").Debugf("Wrote %d bytes to static asset", len(v))

	b64filename := base64.RawStdEncoding.EncodeToString([]byte(filename))

	fileid := fmt.Sprintf("file%s", b64filename)

	b.files[filename] = fileid

	b.file.Var().Id(fileid).Op("=").Index().Byte().Values(
		v...,
	)

	return nil
}

// Render writes the generated Go code to the supplied io.Writer, returning an
// error on failure to write
func (b *AferoWrapper) Render(w io.Writer) error {
	v := []jen.Code{
		jen.Id("Fs").Op("=").Qual("github.com/spf13/afero", "NewMemMapFs").Call(),
		// jen.Var().Id("o").Index().Byte(),
	}

	for filename, file := range b.files {
		v = append(
			v,
			jen.Qual("github.com/spf13/afero", "WriteFile").Call(
				jen.Id("Fs"),
				jen.Lit(filename),
				jen.Id("decode").Params(jen.Id(file)),
				jen.Qual("os", "ModePerm"),
			),
		)
	}

	b.file.Func().Id("init").Params().Block(
		v...,
	)

	return b.file.Render(w)
}
