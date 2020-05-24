package main

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

//nolint:gochecknoglobals // cobra uses globals in main
var rootCmd = &cobra.Command{
	Use:  "goembed",
	Run:  mainCommand,
	Args: cobra.MinimumNArgs(1),
}

//nolint:gochecknoinits // init is used in main for cobra
func init() {
	cobra.OnInitialize(configInit)

	rootCmd.PersistentFlags().BoolP(
		"debug", "d",
		false,
		"Debug output",
	)
	rootCmd.PersistentFlags().StringP(
		"file", "f",
		"-",
		"Output file, or '-' for STDOUT",
	)
	rootCmd.PersistentFlags().StringSliceP(
		"build", "b",
		[]string{},
		"comma sepearted list of build flags",
	)
	rootCmd.PersistentFlags().StringP(
		"package", "p",
		"",
		"golang package name for file (default: based on output file directory)",
	)
	rootCmd.PersistentFlags().StringP(
		"compression", "c",
		"snappy",
		"Compression to use, options are 'deflate', 'gzip', 'lzw', 'snappy', 'snappystream', 'zlib' or 'none'",
	)
	rootCmd.PersistentFlags().StringP(
		"wrapper", "w",
		"none",
		"Wrapper to use, options are 'none' or 'afero'",
	)

	_ = viper.BindPFlag("debug", rootCmd.PersistentFlags().Lookup("debug"))
	_ = viper.BindEnv("debug", "DEBUG")
	_ = viper.BindPFlag("file", rootCmd.PersistentFlags().Lookup("file"))
	_ = viper.BindEnv("file", "OUTPUT_FILE")
	_ = viper.BindPFlag("build", rootCmd.PersistentFlags().Lookup("build"))
	_ = viper.BindEnv("build", "BUILD_FLAGS")
	_ = viper.BindPFlag("package", rootCmd.PersistentFlags().Lookup("package"))
	_ = viper.BindEnv("package", "PACKAGE_NAME")
	_ = viper.BindPFlag("compression", rootCmd.PersistentFlags().Lookup("compression"))
	_ = viper.BindEnv("compression", "COMPRESSION")
	_ = viper.BindPFlag("wrapper", rootCmd.PersistentFlags().Lookup("wrapper"))
	_ = viper.BindEnv("wrapper", "WRAPPER")
}

func main() {
	_ = rootCmd.Execute()
}
