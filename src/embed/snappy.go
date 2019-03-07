package embed

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"io"
	"io/ioutil"

	"github.com/dave/jennifer/jen"
	"github.com/golang/snappy"
	"github.com/sirupsen/logrus"
)

// SnappyBuilder is a Builder compatible struct that uses snappy compression
type SnappyBuilder struct {
	file  *jen.File
	files map[string]string
}

// NewSnappyBuilder returns a Builder compatible class that uses snappy compression
func NewSnappyBuilder(packageName string) Builder {
	f := jen.NewFile(packageName)
	f.HeaderComment("This file is generated - do not edit.")
	f.Line()
	f.ImportName("github.com/spf13/afero", "afero")
	f.Var().Id("Fs").Id("afero.Fs")

	return &SnappyBuilder{
		file:  f,
		files: map[string]string{},
	}
}

// AddFile adds a file to the embedded package.
func (b *SnappyBuilder) AddFile(filename string, file io.Reader) error {
	v := []jen.Code{}

	src, err := ioutil.ReadAll(file)
	if err != nil {
		return err
	}
	encoded := snappy.Encode(nil, src)
	logrus.WithField("compression", "snappy").Debugf("Copied %d bytes into compressor", len(src))
	cmpOut := bytes.NewBuffer(encoded)

	buf := make([]byte, 1)
	for {
		_, err := cmpOut.Read(buf)
		if err != nil {
			break
		}
		v = append(v, jen.Lit(int(buf[0])))
	}

	logrus.WithField("compression", "snappy").Debugf("Wrote %d bytes to static asset", len(v))

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
func (b *SnappyBuilder) Render(w io.Writer) error {
	v := []jen.Code{
		jen.Id("Fs").Op("=").Qual("github.com/spf13/afero", "NewMemMapFs").Call(),
		jen.Var().Id("o").Index().Byte(),
	}

	for filename, file := range b.files {
		v = append(
			v,
			jen.List(jen.Id("o"), jen.Id("_")).Op("=").Qual("github.com/golang/snappy", "Decode").Call(jen.Nil(), jen.Id(file)),
			jen.Qual("github.com/spf13/afero", "WriteFile").Call(
				jen.Id("Fs"),
				jen.Lit(filename),
				jen.Id("o"),
				jen.Qual("os", "ModePerm"),
			),
		)
	}

	b.file.Func().Id("init").Params().Block(
		v...,
	)

	return b.file.Render(w)
}
