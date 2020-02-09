package cmd

import (
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "kaleidoscope",
	Short: "kaleidoscope",
	Long:  `kaleidoscope`,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		logrus.Fatal(err)
	}
}