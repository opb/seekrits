package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var Version = "unknown"

func init() {
	RootCmd.AddCommand(versionCmd)
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number of Hugo",
	Long:  `All software has versions. This is Hugo's`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("seekrits - version "+Version)
	},
}