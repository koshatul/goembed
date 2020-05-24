package main

import (
	"fmt"

	"github.com/spf13/cobra"
)

const (
	version   = "dev"
	buildDate = "notset"
	gitHash   = ""
)

//nolint:gochecknoinits // init is used in main for cobra
func init() {
	rootCmd.AddCommand(versionCmd)

	rootCmd.Version = fmt.Sprintf("%s [%s] (%s)", version, gitHash, buildDate)
}

//nolint:gochecknoglobals // cobra uses globals in main
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version",
	Run:   versionCommand,
}

func versionCommand(cmd *cobra.Command, args []string) {
	fmt.Printf("%s version %s [%s] (%s)\n", rootCmd.Name(), version, gitHash, buildDate)
}
