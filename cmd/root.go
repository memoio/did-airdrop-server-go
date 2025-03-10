package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "did server",
	Short: "A did server",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("did server")
	},
}

func Exceute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.AddCommand(ServerCmd)
	rootCmd.AddCommand(versionCmd)
}
