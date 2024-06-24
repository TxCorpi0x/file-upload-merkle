package cli

import (
	"log"

	"github.com/spf13/cobra"
)

func init() {
	Cmd.AddCommand(uploadCmd)
	Cmd.AddCommand(downloadCmd)
}

const (
	defaultServerURL          = "http://localhost:8080"
	defaultMerkleRootFilename = ".runtime/merkleroot"
)

var Cmd = &cobra.Command{
	Use:   "client",
	Short: "The fxmerkle client can upload & download files and verify their integrity",
	Run: func(cmd *cobra.Command, args []string) {
		if err := cmd.Help(); err != nil {
			log.Fatal(err)
		}
	},
}
