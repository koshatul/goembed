package embed

import (
	"encoding/base64"
	"fmt"
	"io"

	"github.com/dave/jennifer/jen"
	"github.com/sirupsen/logrus"
)

// NoCompressBuilder is a Builder compatible struct that uses no compression
type NoCompressBuilder struct {
	file  *jen.File
	files map[string]string
}

// NewNoCompressBuilder returns a Builder compatible class that uses no compression
func NewNoCompressBuilder(packageName string) Builder {
	f := jen.NewFile(packageName)
	f.HeaderComment("This file is generated - do not edit.")
	f.Line()
	f.ImportName("github.com/spf13/afero", "afero")
	f.Comment("Fs is the filesystem containing the assets embedded in this package.").Line().Var().Id("Fs").Id("afero.Fs")

	return &NoCompressBuilder{
		file:  f,
		files: map[string]string{},
	}
}

// AddFile adds a file to the embedded package.
func (b *NoCompressBuilder) AddFile(filename string, file io.Reader) error {
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
func (b *NoCompressBuilder) Render(w io.Writer) error {
	v := []jen.Code{
		jen.Id("Fs").Op("=").Qual("github.com/spf13/afero", "NewMemMapFs").Call(),
	}

	for filename, file := range b.files {
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

	b.file.Func().Id("init").Params().Block(
		v...,
	)

	return b.file.Render(w)
}
