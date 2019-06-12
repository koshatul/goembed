package shrink

import (
	"bytes"
	"io/ioutil"

	"github.com/dave/jennifer/jen"
	"github.com/golang/snappy"
	"github.com/koshatul/goembed/src/goembed"
	"github.com/sirupsen/logrus"
)

// SnappyShrinker is a Builder compatible struct that uses snappy compression
type SnappyShrinker struct {
}

// NewSnappyShrinker returns a Builder compatible class that uses snappy compression
func NewSnappyShrinker() Shrinker {
	return &SnappyShrinker{}
}

// Compress returns a byte array of compressed file data
func (b *SnappyShrinker) Compress(file goembed.File) ([]jen.Code, error) {
	v := []jen.Code{}

	src, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, err
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

	// b64filename := base64.RawStdEncoding.EncodeToString([]byte(filename))

	// fileid := fmt.Sprintf("file%s", b64filename)

	// b.files[filename] = fileid

	// b.file.Var().Id(fileid).Op("=").Index().Byte().Values(
	// 	v...,
	// )

	return v, nil
}

func (b *SnappyShrinker) Decompressor() []jen.Code {
	return []jen.Code{
		//snappy.Decode(nil, fileL2luZGV4Lmh0bWw)
		jen.List(jen.Id("o"), jen.Id("_")).Op(":=").Qual("github.com/golang/snappy", "Decode").Call(jen.Nil(), jen.Id("input")),
		jen.Return(jen.Id("o")),
	}
}

func (b *SnappyShrinker) Header() []jen.Code {
	return []jen.Code{}
}

// // AddFile adds a file to the embedded package.
// func (b *SnappyShrinker) AddFile(filename string, file io.Reader) error {
// 	v := []jen.Code{}

// 	src, err := ioutil.ReadAll(file)
// 	if err != nil {
// 		return err
// 	}
// 	encoded := snappy.Encode(nil, src)
// 	logrus.WithField("compression", "snappy").Debugf("Copied %d bytes into compressor", len(src))
// 	cmpOut := bytes.NewBuffer(encoded)

// 	buf := make([]byte, 1)
// 	for {
// 		_, err := cmpOut.Read(buf)
// 		if err != nil {
// 			break
// 		}
// 		v = append(v, jen.Lit(int(buf[0])))
// 	}

// 	logrus.WithField("compression", "snappy").Debugf("Wrote %d bytes to static asset", len(v))

// 	b64filename := base64.RawStdEncoding.EncodeToString([]byte(filename))

// 	fileid := fmt.Sprintf("file%s", b64filename)

// 	b.files[filename] = fileid

// 	b.file.Var().Id(fileid).Op("=").Index().Byte().Values(
// 		v...,
// 	)

// 	return nil
// }

// // Render writes the generated Go code to the supplied io.Writer, returning an
// // error on failure to write
// func (b *SnappyShrinker) Render(w io.Writer) error {
// 	v := []jen.Code{
// 		jen.Id("Fs").Op("=").Qual("github.com/spf13/afero", "NewMemMapFs").Call(),
// 		jen.Var().Id("o").Index().Byte(),
// 	}

// 	for filename, file := range b.files {
// 		v = append(
// 			v,
// 			jen.List(jen.Id("o"), jen.Id("_")).Op("=").Qual("github.com/golang/snappy", "Decode").Call(jen.Nil(), jen.Id(file)),
// 			jen.Qual("github.com/spf13/afero", "WriteFile").Call(
// 				jen.Id("Fs"),
// 				jen.Lit(filename),
// 				jen.Id("o"),
// 				jen.Qual("os", "ModePerm"),
// 			),
// 		)
// 	}

// 	b.file.Func().Id("init").Params().Block(
// 		v...,
// 	)

// 	return b.file.Render(w)
// }
