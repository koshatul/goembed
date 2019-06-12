package wrap

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

	f.Add(shrinker.Header()...)

	decodeFunc := shrinker.Decompressor()
	if decodeFunc != nil {
		f.Func().Id("decode").Params(jen.Id("input").Index().Byte()).Params(jen.Index().Byte()).Block(
			decodeFunc...,
		)
	}

	return &AferoWrapper{
		file:     f,
		files:    map[string]string{},
		shrinker: shrinker,
	}
}

// Name returns a simple name for this module
func (b *AferoWrapper) Name() string {
	return "afero"
}

// AddFile adds a file to the embedded package.
func (b *AferoWrapper) AddFile(filename string, file goembed.File) error {
	v, err := b.shrinker.Compress(file)
	if err != nil {
		return err
	}

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

	useDecodeFunc := b.shrinker.Decompressor()
	if b.shrinker.IsReaderWithError() {
		v = append(
			v,
			jen.Var().Id("o").Qual("io", "Reader"),
			jen.Var().Id("err").Error(),
		)
	}

	for filename, file := range b.files {
		if b.shrinker.IsStream() {
			if b.shrinker.IsReaderWithError() {
				v = append(
					v,
					jen.List(jen.Id("o"), jen.Id("err")).Op("=").Add(b.shrinker.ReaderWithError(jen.Id(file))),
					jen.If(jen.Id("err").Op("==").Nil()).Block(
						jen.Qual("github.com/spf13/afero", "WriteReader").Call(
							jen.Id("Fs"),
							jen.Lit(filename),
							jen.Add(jen.Id("o")),
						),
					),
				)
			} else {
				v = append(
					v,
					jen.Qual("github.com/spf13/afero", "WriteReader").Call(
						jen.Id("Fs"),
						jen.Lit(filename),
						jen.Add(b.shrinker.Reader(jen.Id(file))),
					),
				)
			}
		} else if useDecodeFunc != nil {
			v = append(
				v,
				jen.Qual("github.com/spf13/afero", "WriteFile").Call(
					jen.Id("Fs"),
					jen.Lit(filename),
					jen.Id("decode").Params(jen.Id(file)),
					jen.Qual("os", "ModePerm"),
				),
			)
		} else {
			v = append(
				v,
				jen.Qual("github.com/spf13/afero", "WriteFile").Call(
					jen.Id("Fs"),
					jen.Lit(filename),
					jen.Id(file),
					jen.Qual("os", "ModePerm"),
				),
			)

		}
	}

	b.file.Func().Id("init").Params().Block(
		v...,
	)

	return b.file.Render(w)
}
