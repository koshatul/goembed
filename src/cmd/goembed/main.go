package main

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var rootCmd = &cobra.Command{
	Use:  "goembed",
	Run:  mainCommand,
	Args: cobra.MinimumNArgs(1),
}

func init() {
	cobra.OnInitialize(configInit)

	rootCmd.PersistentFlags().BoolP("debug", "d", false, "Debug output")
	viper.BindPFlag("debug", rootCmd.PersistentFlags().Lookup("debug"))
	viper.BindEnv("debug", "DEBUG")

	rootCmd.PersistentFlags().StringP("file", "f", "-", "Output file, or '-' for STDOUT")
	viper.BindPFlag("file", rootCmd.PersistentFlags().Lookup("file"))
	viper.BindEnv("file", "OUTPUT_FILE")

	rootCmd.PersistentFlags().StringP("package", "p", "", "golang package name for file (default: based on output file directory)")
	viper.BindPFlag("package", rootCmd.PersistentFlags().Lookup("package"))
	viper.BindEnv("package", "PACKAGE_NAME")

	rootCmd.PersistentFlags().StringP("compression", "c", "snappy", "Compression to use, options are 'deflate', 'gzip', 'lzw', 'snappy', 'snappystream', 'zlib' or 'none'")
	viper.BindPFlag("compression", rootCmd.PersistentFlags().Lookup("compression"))
	viper.BindEnv("compression", "COMPRESSION")

	rootCmd.PersistentFlags().StringP("wrapper", "w", "none", "Wrapper to use, options are 'none' or 'afero'")
	viper.BindPFlag("wrapper", rootCmd.PersistentFlags().Lookup("wrapper"))
	viper.BindEnv("wrapper", "WRAPPER")
}

func main() {
	rootCmd.Execute()
}
