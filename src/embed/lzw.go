package embed

import (
	"bytes"
	"compress/lzw"
	"encoding/base64"
	"fmt"
	"io"

	"github.com/dave/jennifer/jen"
	"github.com/sirupsen/logrus"
)

// LzwBuilder is a Builder compatible struct that uses lzw compression
type LzwBuilder struct {
	file  *jen.File
	files map[string]string
}

// NewLzwBuilder returns a Builder compatible class that uses lzw compression
func NewLzwBuilder(packageName string) Builder {
	f := jen.NewFile(packageName)
	f.HeaderComment("This file is generated - do not edit.")
	f.Line()
	f.ImportName("github.com/spf13/afero", "afero")
	f.Var().Id("Fs").Id("afero.Fs")

	return &LzwBuilder{
		file:  f,
		files: map[string]string{},
	}
}

// AddFile adds a file to the embedded package.
func (b *LzwBuilder) AddFile(filename string, file io.Reader) error {
	v := []jen.Code{}

	cmpOut := new(bytes.Buffer)
	rawIn := lzw.NewWriter(cmpOut, lzw.LSB, 8)
	n, err := io.Copy(rawIn, file)
	rawIn.Close()
	if err != nil {
		return err
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

	b64filename := base64.RawStdEncoding.EncodeToString([]byte(filename))

	fileid := fmt.Sprintf("file_%s", b64filename)

	b.files[filename] = fileid

	b.file.Var().Id(fileid).Op("=").Index().Byte().Values(
		v...,
	)

	return nil
}

// Render writes the generated Go code to the supplied io.Writer, returning an
// error on failure to write
func (b *LzwBuilder) Render(w io.Writer) error {
	v := []jen.Code{
		jen.Id("Fs").Op("=").Qual("github.com/spf13/afero", "NewMemMapFs").Call(),

		jen.Var().Id("bufIn").Op("*").Qual("bytes", "Buffer"),
		jen.Var().Id("bufOut").Op("*").Qual("bytes", "Buffer"),
		jen.Var().Id("cmpOut").Qual("io", "ReadCloser"),
	}

	for filename, file := range b.files {
		v = append(
			v,
			jen.Id("bufIn").Op("=").Qual("bytes", "NewBuffer").Call(jen.Id(file)),
			jen.Id("cmpOut").Op("=").Qual("compress/lzw", "NewReader").Call(jen.Id("bufIn"), jen.Qual("compress/lzw", "LSB"), jen.Lit(8)),
			jen.Id("bufOut").Op("=").New(jen.Qual("bytes", "Buffer")),
			jen.Qual("io", "Copy").Call(jen.Id("bufOut"), jen.Id("cmpOut")),
			jen.Qual("github.com/spf13/afero", "WriteFile").Call(
				jen.Id("Fs"),
				jen.Lit(filename),
				jen.Id("bufOut").Dot("Bytes").Call(),
				jen.Qual("os", "ModePerm"),
			),
			jen.Id("cmpOut").Dot("Close").Call(),
		)
	}

	b.file.Func().Id("init").Params().Block(
		v...,
	)

	return b.file.Render(w)
}
