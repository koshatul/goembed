package wrap

import (
	"encoding/base64"
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/dave/jennifer/jen"
	"github.com/koshatul/goembed/goembed"
	"github.com/koshatul/goembed/shrink"
	"github.com/sirupsen/logrus"
)

// NoDepWrapper is a Wrapper compatible struct that uses no dependencies.
type NoDepWrapper struct {
	file           *jen.File
	files          map[string]string
	shrinker       shrink.Shrinker
	children       map[string]*jen.Statement
	openSwitchFunc *jen.Statement
	buildTags      []string
}

const noDepFileMode = 0444

// NewNoDepWrapper returns a Wrapper compatible class that uses no dependencies for the file system.
//nolint:funlen // it's long but it's readable
func NewNoDepWrapper(packageName string, shrinker shrink.Shrinker, opts ...Option) Wrapper {
	w := &NoDepWrapper{
		files:    map[string]string{},
		shrinker: shrinker,
		children: map[string]*jen.Statement{},
	}

	for _, opt := range opts {
		opt(w)
	}

	f := jen.NewFile(packageName)

	f.HeaderComment("Code generated - DO NOT EDIT.")

	if len(w.buildTags) > 0 {
		f.HeaderComment(fmt.Sprintf("+build %s", strings.Join(w.buildTags, ",")))
	}

	f.Line()

	f.Add(shrinker.Header()...)

	f.Comment("Fs is the filesystem containing the assets embedded in this package.").Line().Var().Id("Fs").Qual("net/http", "FileSystem").Op("=").Op("&").Id("fs").Values()

	f.Type().Id("assetFileData").Struct(
		jen.Id("name").String(),
		jen.Id("data").Index().Byte(),
		jen.Id("dir").Bool(),
		jen.Id("size").Int64(),
		jen.Id("modtime").Qual("time", "Time"),
		jen.Id("children").Index().Qual("os", "FileInfo"),
	)

	openSwitchFunc := jen.Switch(jen.Id("name"))
	w.openSwitchFunc = openSwitchFunc

	f.Type().Id("fs").Struct()
	f.Func().Params(jen.Id("a").Id("fs")).Id("Open").Params(
		jen.Id("name").String(),
	).Params(
		jen.Qual("net/http", "File"),
		jen.Error(),
	).Block(
		jen.Return(jen.Id("Open").Call(jen.Id("name"))),
	)

	f.Func().Id("Open").Params(
		jen.Id("name").String(),
	).Params(
		jen.Qual("net/http", "File"),
		jen.Error(),
	).Block(
		openSwitchFunc,
		jen.Return(
			jen.Nil(),
			jen.Qual("os", "ErrNotExist"),
		),
	)

	f.Type().Id("assetFileInfo").Struct(
		jen.Id("f").Op("*").Id("assetFile"),
	)

	f.Func().Params(jen.Id("a").Id("assetFileInfo")).Id("Name").Params().Params(jen.String()).Block(jen.Return(jen.Qual("path", "Base").Call(jen.Id("a").Dot("f").Dot("name"))))
	f.Func().Params(jen.Id("a").Id("assetFileInfo")).Id("Size").Params().Params(jen.Int64()).Block(jen.Return(jen.Id("a").Dot("f").Dot("size")))
	f.Func().Params(jen.Id("a").Id("assetFileInfo")).Id("Mode").Params().Params(jen.Qual("os", "FileMode")).Block(jen.Return(jen.Lit(noDepFileMode)))
	f.Func().Params(jen.Id("a").Id("assetFileInfo")).Id("ModTime").Params().Params(jen.Qual("time", "Time")).Block(jen.Return(jen.Id("a").Dot("f").Dot("modtime")))
	f.Func().Params(jen.Id("a").Id("assetFileInfo")).Id("IsDir").Params().Params(jen.Bool()).Block(jen.Return(jen.Id("a").Dot("f").Dot("dir")))
	f.Func().Params(jen.Id("a").Id("assetFileInfo")).Id("Sys").Params().Params(jen.Interface()).Block(jen.Return(jen.Nil()))

	f.Type().Id("assetFile").Struct(
		jen.Qual("io", "Reader"),
		jen.Qual("io", "Seeker"),
		jen.Op("*").Id("assetFileData"),
	)

	f.Func().Params(jen.Id("a").Op("*").Id("assetFile")).Id("Stat").Params().Params(jen.Qual("os", "FileInfo"), jen.Error()).Block(
		jen.Return(jen.Id("assetFileInfo").Values(jen.Id("f").Op(":").Id("a")), jen.Nil()),
	)

	f.Func().Params(jen.Id("a").Op("*").Id("assetFile")).Id("Readdir").Params(jen.Id("count").Int()).Params(jen.Index().Qual("os", "FileInfo"), jen.Error()).Block(
		jen.If(jen.Id("a").Dot("dir")).Block(
			jen.Return(jen.Id("a").Dot("children"), jen.Nil()),
		),
		jen.Return(jen.Nil(), jen.Nil()),
	)

	f.Func().Params(jen.Id("a").Op("*").Id("assetFile")).Id("Close").Params().Params(jen.Error()).Block(
		jen.Return(jen.Nil()),
	)

	decodeFunc := shrinker.Decompressor()
	if decodeFunc != nil {
		f.Func().Id("decode").Params(jen.Id("input").Index().Byte()).Params(jen.Index().Byte()).Block(
			decodeFunc...,
		)
	}

	w.file = f

	return w
}

// Name returns a simple name for this module.
func (b *NoDepWrapper) Name() string {
	return "nodep"
}

func (b *NoDepWrapper) addDir(dir string) {
	for name := range b.files {
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
	b.children[dir] = jen.Null()

	b.file.Var().Id(fileid).Op("*").Id("assetFileData").Op("=").Op("&").Id("assetFileData").Values(
		jen.Id("name").Op(":").Lit(dir),
		jen.Id("dir").Op(":").Lit(true),
		jen.Id("modtime").Op(":").Qual("time", "Unix").Params(jen.Lit(time.Now().Unix()), jen.Lit(0)),
		b.children[dir],
	)

	logrus.WithField("wrapper", "nodep").Debugf("Added directory %s to asset list", dir)
}

// AddFile adds a file to the embedded package.
func (b *NoDepWrapper) AddFile(filename string, file goembed.File) error {
	b.addDir("/")

	sp := strings.Split(filename, "/")

	for len(sp) > 0 {
		sp = sp[:len(sp)-1]
		b.addDir(strings.Join(sp, "/"))
	}

	v, err := b.shrinker.Compress(file)
	if err != nil {
		return err
	}

	b64filename := base64.RawStdEncoding.EncodeToString([]byte(filename))
	fileid := fmt.Sprintf("file%s", b64filename)
	b.files[filename] = fileid

	b.file.Var().Id(fileid).Op("*").Id("assetFileData").Op("=").Op("&").Id("assetFileData").Values(
		jen.Id("name").Op(":").Lit(filename),
		jen.Id("dir").Op(":").Lit(false),
		jen.Id("size").Op(":").Lit(file.Stat.Size()),
		jen.Id("modtime").Op(":").Qual("time", "Unix").Params(jen.Lit(file.Stat.ModTime().Unix()), jen.Lit(0)),
		jen.Id("data").Op(":").Index().Byte().Values(v...),
	)

	logrus.WithField("wrapper", "nodep").Debugf("Added file %s to asset list", filename)

	return nil
}

// Render writes the generated Go code to the supplied io.Writer, returning an.
// error on failure to write.
func (b *NoDepWrapper) Render(w io.Writer) error {
	caseList := []jen.Code{}

	for filename, file := range b.files {
		if b.shrinker.IsReaderWithError() {
			caseList = append(
				caseList,
				jen.Case(jen.Lit(filename)).Block(
					jen.List(jen.Id("r"), jen.Id("err")).Op(":=").Add(b.shrinker.ReaderWithError(jen.Id(file).Dot("data"))),
					jen.Return(
						jen.Op("&").Id("assetFile").Values(
							jen.Id("Reader").Op(":").Id("r"),
							jen.Id("assetFileData").Op(":").Id(file),
						),
						jen.Id("err"),
					),
				),
			)
		} else {
			caseList = append(
				caseList,
				jen.Case(jen.Lit(filename)).Block(
					jen.Return(
						jen.Op("&").Id("assetFile").Values(
							jen.Id("Reader").Op(":").Add(b.shrinker.Reader(jen.Id(file).Dot("data"))),
							jen.Id("assetFileData").Op(":").Id(file),
						),
						jen.Nil(),
					),
				),
			)
		}

		children := []jen.Code{}

		for f, v := range b.files {
			if !strings.EqualFold(filename, f) && strings.HasPrefix(f, filename) {
				ft := strings.TrimLeft(f[len(filename):], "/")
				if !strings.Contains(ft, "/") {
					children = append(children, jen.Op("&").Id("assetFileInfo").Values(jen.Id("f").Op(":").Op("&").Id("assetFile").Values(jen.Id("assetFileData").Op(":").Id(v))))
				}
			}
		}

		if len(children) > 0 {
			b.children[filename].Id("children").Op(":").Index().Qual("os", "FileInfo").Values(
				children...,
			)
		}
	}

	b.openSwitchFunc.Block(
		caseList...,
	)

	return b.file.Render(w)
}
