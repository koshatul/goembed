package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/koshatul/goembed/src/shrink"

	"github.com/koshatul/goembed/src/goembed"
	"github.com/koshatul/goembed/src/wrap"
	"github.com/sirupsen/logrus"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func processDirectory(e wrap.Wrapper, s shrink.Shrinker, absPath string) error {
	logrus.Infof("Processing directory: %s", absPath)
	convertFs := afero.NewBasePathFs(afero.NewOsFs(), fmt.Sprintf("%s/", absPath))

	return afero.Walk(convertFs, "/", func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() {
			logrus.WithFields(logrus.Fields{"file": path, "compression": s.Name(), "wrapper": e.Name()}).Infof("Adding file: %s", path)
			f, err := convertFs.Open(path)
			if err != nil {
				logrus.WithFields(logrus.Fields{"file": path, "compression": s.Name(), "wrapper": e.Name()}).Errorf("Unable to open file: %s", err)
				return err
			}
			st, err := f.Stat()
			if err != nil {
				logrus.WithFields(logrus.Fields{"file": path, "compression": s.Name(), "wrapper": e.Name()}).Errorf("Unable to get file stat: %s", err)
				return err
			}
			err = e.AddFile(path, goembed.NewFile(path, st, f))
			if err != nil {
				logrus.WithFields(logrus.Fields{"file": path, "compression": s.Name(), "wrapper": e.Name()}).Errorf("Unable to add file: %s", err)
				return err
			}
		}
		return nil
	})
}

func processFile(e wrap.Wrapper, s shrink.Shrinker, absPath string) error {
	path := filepath.Base(absPath)
	logrus.WithFields(logrus.Fields{"file": path, "compression": s.Name(), "wrapper": e.Name()}).Infof("Processing file: %s", path)
	f, err := os.Open(absPath)
	if err != nil {
		logrus.WithFields(logrus.Fields{"file": path, "compression": s.Name(), "wrapper": e.Name()}).Errorf("Unable to open file: %s", err)
		return err
	}
	st, err := f.Stat()
	if err != nil {
		logrus.WithFields(logrus.Fields{"file": path, "compression": s.Name(), "wrapper": e.Name()}).Errorf("Unable to get file stat: %s", err)
		return err
	}
	err = e.AddFile(path, goembed.NewFile(path, st, f))
	if err != nil {
		logrus.WithFields(logrus.Fields{"file": path, "compression": s.Name(), "wrapper": e.Name()}).Errorf("Unable to add file: %s", err)
		return err
	}

	return nil
}

func processPath(e wrap.Wrapper, s shrink.Shrinker, path string) error {
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
		err := processDirectory(e, s, absPath)
		if err != nil {
			logrus.Errorf("Unable to process directory: %s", err)
			return err
		}
	} else {
		err := processFile(e, s, absPath)
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

	var s shrink.Shrinker
	switch strings.ToLower(viper.GetString("compression")) {
	case "none", "nocompress":
		s = shrink.NewNoShrinker()
	case "deflate":
		s = shrink.NewDeflateStreamShrinker()
	case "gzip":
		s = shrink.NewGzipStreamShrinker()
	case "lzw":
		s = shrink.NewLzwStreamShrinker()
	case "snappy":
		s = shrink.NewSnappyShrinker()
	case "snappystream":
		s = shrink.NewSnappyStreamShrinker()
	case "zlib":
		s = shrink.NewZlibStreamShrinker()
	default:
		logrus.Errorf("Invalid compression type: %s", strings.ToLower(viper.GetString("compression")))
		cmd.Help()
		return
	}

	var e wrap.Wrapper
	switch strings.ToLower(viper.GetString("wrapper")) {
	case "none", "nodep":
		e = wrap.NewNoDepWrapper(viper.GetString("package"), s)
	case "afero":
		e = wrap.NewAferoWrapper(viper.GetString("package"), s)
	default:
		logrus.Errorf("Invalid wrapper type: %s", strings.ToLower(viper.GetString("wrapper")))
		cmd.Help()
		return
	}

	for _, path := range args {
		err := processPath(e, s, path)
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
