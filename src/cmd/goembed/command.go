package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/koshatul/goembed/src/embed"
	"github.com/sirupsen/logrus"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func mainCommand(cmd *cobra.Command, args []string) {
	var e embed.Builder
	switch strings.ToLower(viper.GetString("compression")) {
	case "none", "nocompress":
		e = embed.NewNoCompressBuilder(viper.GetString("package"))
	case "deflate":
		e = embed.NewDeflateBuilder(viper.GetString("package"))
	case "gzip":
		e = embed.NewGzipBuilder(viper.GetString("package"))
	case "lzw":
		e = embed.NewLzwBuilder(viper.GetString("package"))
	case "snappy":
		e = embed.NewSnappyBuilder(viper.GetString("package"))
	case "zlib":
		e = embed.NewZlibBuilder(viper.GetString("package"))
	default:
		logrus.Errorf("Invalid compression type: %s", strings.ToLower(viper.GetString("compression")))
		cmd.Help()
		return
	}

	for _, path := range args {
		absPath, err := filepath.Abs(path)
		if err != nil {
			log.Fatal(err)
		}

		absStat, err := os.Stat(absPath)
		if err != nil {
			log.Fatal(err)
		}

		if absStat.IsDir() {
			logrus.Infof("Processing directory: %s", absPath)
			convertFs := afero.NewBasePathFs(afero.NewOsFs(), fmt.Sprintf("%s/", absPath))

			afero.Walk(convertFs, "/", func(path string, info os.FileInfo, err error) error {
				if !info.IsDir() {
					logrus.WithField("file", path).Infof("Adding file: %s", path)
					f, err := convertFs.Open(path)
					if err != nil {
						logrus.WithField("file", path).Errorf("Unable to open file: %s", err)
						return err
					}
					err = e.AddFile(path, f)
					if err != nil {
						logrus.WithField("file", path).Errorf("Unable to add file: %s", err)
						return err
					}
				}
				return nil
			})
		} else {
			path := filepath.Base(absPath)
			logrus.WithField("file", path).Infof("Processing file: %s", path)
			f, err := os.Open(absPath)
			if err != nil {
				logrus.WithField("file", path).Errorf("Unable to open file: %s", err)
			}
			err = e.AddFile(path, f)
			if err != nil {
				logrus.WithField("file", path).Errorf("Unable to add file: %s", err)
			}

		}
	}

	var out io.Writer
	switch viper.GetString("file") {
	case "-":
		out = os.Stdout
	default:
		f, err := os.Create(viper.GetString("file"))
		if err != nil {
			log.Fatal(err)
		}
		out = f
		defer f.Close()
	}

	err := e.Render(out)
	if err != nil {
		log.Fatal(err)
	}

}
