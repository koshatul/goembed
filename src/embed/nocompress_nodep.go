package embed

import (
	"encoding/base64"
	"fmt"
	"io"
	"strings"

	"github.com/dave/jennifer/jen"
	"github.com/sirupsen/logrus"
)

// NoCompressNoDepBuilder is a Builder compatible struct that uses no compression and uses no dependencies.
type NoCompressNoDepBuilder struct {
	file     *jen.File
	files    map[string]string
	children map[string]*jen.Code
}

// NewNoCompressNoDepBuilder returns a Builder compatible class that uses no compression
func NewNoCompressNoDepBuilder(packageName string) Builder {
	f := jen.NewFile(packageName)
	f.HeaderComment("Code generated - DO NOT EDIT.")
	f.Line()
	f.Type().Id("assetFileData").Struct(
		jen.Id("name").String(),
		jen.Id("data").Index().Byte(),
		jen.Id("dir").Bool(),
	)
	f.Func().Params(jen.Id("a").Op("*").Id("assetFileData")).Id("Children").Params().Params(jen.Index().Op("*").Id("assetFileData")).Block(
		jen.Id("o").Op(":=").Index().Op("*").Id("assetFileData").Values(),
		jen.For(
			jen.Id("f").Op(",").Id("v").Op(":=").Range().Id("fileData").Block(
				jen.If(
					jen.Op("!").Qual("strings", "EqualFold").Params(jen.Id("a").Dot("name"), jen.Id("f")),
					jen.Qual("strings", "HasPrefix").Params(jen.Id("a").Dot("name"), jen.Id("f")),
				).Block(
					jen.Id("ft").Op(":=").Id("f").Index(jen.Len(jen.Id("a").Dot("name")), jen.Null()),
					jen.If(
						jen.Op("!").Qual("strings", "Contains").Params(jen.Id("ft"), jen.Lit("/")).Block(
							jen.Id("o").Op("=").Append(jen.Id("o"), jen.Id("v")),
						),
					),
				),
			),
		),
		jen.Return(jen.Id("o")),
	)

	f.Type().Id("Fs").Struct()
	f.Func().Params(jen.Id("a").Id("Fs")).Id("Open").Params(
		jen.Id("name").String(),
	).Params(
		jen.Qual("net/http", "File"),
		jen.Error(),
	).Block(
		jen.If(
			jen.Id("v").Op(",").Id("ok").Op(":=").Id("fileData").Index(jen.Id("name")),
			jen.Id("ok"),
		).Block(
			jen.Return(
				jen.Op("&").Id("assetFile").Values(
					jen.Id("Reader").Op(":").Qual("bytes", "NewReader").Params(jen.Id("v").Dot("data")),
					jen.Id("assetFileData").Op(":").Id("v"),
				),
			),
		),
		jen.Return(
			jen.Nil(),
			jen.Qual("os", "ErrNotExist"),
		),
	)

	f.Type().Id("assetFileInfo").Struct(
		jen.Id("f").Op("*").Id("assetFile"),
	)

	f.Func().Params(jen.Id("a").Id("assetFileInfo")).Id("Name").Params().Params(jen.String()).Block(jen.Return(jen.Id("a").Dot("f").Dot("name")))
	f.Func().Params(jen.Id("a").Id("assetFileInfo")).Id("Size").Params().Params(jen.Int64()).Block(jen.Return(jen.Id("int64").Call(jen.Id("len").Call(jen.Id("a").Dot("f").Dot("data")))))
	f.Func().Params(jen.Id("a").Id("assetFileInfo")).Id("Mode").Params().Params(jen.Qual("os", "FileMode")).Block(jen.Return(jen.Lit(0444)))
	f.Func().Params(jen.Id("a").Id("assetFileInfo")).Id("ModTime").Params().Params(jen.Qual("time", "Time")).Block(jen.Return(jen.Qual("time", "Time").Block()))
	f.Func().Params(jen.Id("a").Id("assetFileInfo")).Id("IsDir").Params().Params(jen.Bool()).Block(jen.Return(jen.Id("a").Dot("f").Dot("dir")))
	f.Func().Params(jen.Id("a").Id("assetFileInfo")).Id("Sys").Params().Params(jen.Interface()).Block(jen.Return(jen.Nil()))

	f.Type().Id("assetFile").Struct(
		jen.Op("*").Qual("bytes", "Reader"),
		jen.Op("*").Id("assetFileData"),
	)

	f.Func().Params(jen.Id("a").Op("*").Id("assetFile")).Id("Stat").Params().Params(jen.Qual("os", "FileInfo"), jen.Error()).Block(
		jen.Return(jen.Id("assetFileInfo").Values(jen.Id("f").Op(":").Id("a")), jen.Nil()),
	)

	f.Func().Params(jen.Id("a").Op("*").Id("assetFile")).Id("Readdir").Params(jen.Id("count").Int()).Params(jen.Index().Qual("os", "FileInfo"), jen.Error()).Block(
		jen.If(jen.Id("a").Dot("dir")).Block(
			jen.Id("fl").Op(":=").Index().Qual("os", "FileInfo").Block(),
			jen.For(jen.Id("_").Op(",").Id("ok").Op(":=").Range().Id("a").Dot("children")).Block(
				jen.Id("d").Op(":=").Op("&").Id("assetFile").Values(jen.Id("assetFileData").Op(":").Id("c")),
				jen.Id("fl").Op("=").Append(jen.Id("fl"), jen.Op("&").Id("assetFileInfo").Values(jen.Id("f").Op(":").Id("d"))),
			),
			jen.Return(jen.Id("fl"), jen.Nil()),
		),
		jen.Return(jen.Nil(), jen.Nil()),
	)

	f.Func().Params(jen.Id("a").Op("*").Id("assetFile")).Id("Close").Params().Params(jen.Error()).Block(
		jen.Return(jen.Nil()),
	)

	// f.Comment("Fs is the filesystem containing the assets embedded in this package.").Line().Var().Id("Fs").Id("afero.Fs")

	return &NoCompressNoDepBuilder{
		file:     f,
		files:    map[string]string{},
		children: map[string]*jen.Code{},
	}
}

func (b *NoCompressNoDepBuilder) addDir(dir string) {
	for name, _ := range b.files {
		logrus.Infof("Name: %s; Dir: %s", name, dir)
		if strings.EqualFold(dir, name) {
			return
		}
	}

	if strings.EqualFold(dir, "") {
		return
	}
	b64filename := base64.RawStdEncoding.EncodeToString([]byte(dir))
	fileid := fmt.Sprintf("dir%s", b64filename)
	b.files[dir] = fileid

	b.file.Var().Id(fileid).Op("*").Id("assetFileData").Op("=").Op("&").Id("assetFileData").Values(
		jen.Id("name").Op(":").Lit(dir),
		jen.Id("dir").Op(":").Lit(true),
	)

	logrus.WithField("compression", "none").Debugf("Added directory %s to asset list", dir)
}

// AddFile adds a file to the embedded package.
func (b *NoCompressNoDepBuilder) AddFile(filename string, file io.Reader) error {
	v := []jen.Code{}

	b.addDir("/")

	sp := strings.Split(filename, "/")

	for len(sp) > 0 {
		sp = sp[:len(sp)-1]
		b.addDir(strings.Join(sp, "/"))
	}

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

	b.file.Var().Id(fileid).Op("*").Id("assetFileData").Op("=").Op("&").Id("assetFileData").Values(
		jen.Id("name").Op(":").Lit(filename),
		jen.Id("dir").Op(":").Lit(false),
		jen.Id("data").Op(":").Index().Byte().Values(v...),
	)
	// b.file.Var().Id(fileid).Op("=").Index().Byte().Values(
	// 	v...,
	// )

	return nil
}

// Render writes the generated Go code to the supplied io.Writer, returning an
// error on failure to write
func (b *NoCompressNoDepBuilder) Render(w io.Writer) error {
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
