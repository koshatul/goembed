package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/koshatul/goembed/src/shrink"

	"github.com/koshatul/goembed/src/goembed"
	"github.com/koshatul/goembed/src/wrap"
	"github.com/sirupsen/logrus"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func processDirectory(e wrap.Wrapper, absPath string) error {
	logrus.Infof("Processing directory: %s", absPath)
	convertFs := afero.NewBasePathFs(afero.NewOsFs(), fmt.Sprintf("%s/", absPath))

	return afero.Walk(convertFs, "/", func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() {
			logrus.WithField("file", path).Infof("Adding file: %s", path)
			f, err := convertFs.Open(path)
			if err != nil {
				logrus.WithField("file", path).Errorf("Unable to open file: %s", err)
				return err
			}
			s, err := f.Stat()
			if err != nil {
				logrus.WithField("file", path).Errorf("Unable to get file stat: %s", err)
				return err
			}
			err = e.AddFile(path, goembed.NewFile(path, s, f))
			if err != nil {
				logrus.WithField("file", path).Errorf("Unable to add file: %s", err)
				return err
			}
		}
		return nil
	})
}

func processFile(e wrap.Wrapper, absPath string) error {
	path := filepath.Base(absPath)
	logrus.WithField("file", path).Infof("Processing file: %s", path)
	f, err := os.Open(absPath)
	if err != nil {
		logrus.WithField("file", path).Errorf("Unable to open file: %s", err)
		return err
	}
	s, err := f.Stat()
	if err != nil {
		logrus.WithField("file", path).Errorf("Unable to get file stat: %s", err)
		return err
	}
	err = e.AddFile(path, goembed.NewFile(path, s, f))
	if err != nil {
		logrus.WithField("file", path).Errorf("Unable to add file: %s", err)
		return err
	}

	return nil
}

func processPath(e wrap.Wrapper, path string) error {
	absPath, err := filepath.Abs(path)
	if err != nil {
		logrus.Errorf("Unable to process path: %s", err)
		return err
	}

	absStat, err := os.Stat(absPath)
	if err != nil {
		logrus.Errorf("Unable to process path: %s", err)
		return err
	}

	if absStat.IsDir() {
		err := processDirectory(e, absPath)
		if err != nil {
			logrus.Errorf("Unable to process directory: %s", err)
			return err
		}
	} else {
		err := processFile(e, absPath)
		if err != nil {
			logrus.Errorf("Unable to process file: %s", err)
			return err
		}
	}

	return nil
}

func mainCommand(cmd *cobra.Command, args []string) {
	switch viper.GetString("file") {
	case "-":
		// If STDOUT, user *must* specify a package name
		if viper.GetString("package") == "" {
			logrus.Error("If using STDOUT then --package <name> must be specified")
			os.Exit(1)
		}
	default:
		if viper.GetString("package") == "" {
			absPath, err := filepath.Abs(viper.GetString("file"))
			if err != nil {
				logrus.Fatal(err)
			}
			packageName := filepath.Base(filepath.Dir(absPath))
			viper.Set("package", packageName)
		}
	}

	e := wrap.NewNoDepWrapper(viper.GetString("package"), shrink.NewSnappyShrinker())
	// var e embed.Builder
	// switch strings.ToLower(viper.GetString("compression")) {
	// case "none", "nocompress":
	// 	e = embed.NewNoCompressBuilder(viper.GetString("package"))
	// case "none_nodep", "nocompress_nodep":
	// 	e = embed.NewNoCompressNoDepBuilder(viper.GetString("package"))
	// case "deflate":
	// 	e = embed.NewDeflateBuilder(viper.GetString("package"))
	// case "gzip":
	// 	e = embed.NewGzipBuilder(viper.GetString("package"))
	// case "lzw":
	// 	e = embed.NewLzwBuilder(viper.GetString("package"))
	// case "snappy":
	// 	e = embed.NewSnappyBuilder(viper.GetString("package"))
	// case "zlib":
	// 	e = embed.NewZlibBuilder(viper.GetString("package"))
	// default:
	// 	logrus.Errorf("Invalid compression type: %s", strings.ToLower(viper.GetString("compression")))
	// 	cmd.Help()
	// 	return
	// }

	for _, path := range args {
		err := processPath(e, path)
		if err != nil {
			os.Exit(1)
		}
	}

	var out io.Writer
	switch viper.GetString("file") {
	case "-":
		out = os.Stdout
	default:
		f, err := os.Create(viper.GetString("file"))
		if err != nil {
			logrus.Fatal(err)
		}
		out = f
		defer f.Close()
	}

	err := e.Render(out)
	if err != nil {
		logrus.Fatal(err)
	}

}
