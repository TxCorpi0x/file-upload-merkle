package main

import (
	"log"
	"os"

	clientcli "github.com/TxCorpi0x/file-upload-merkle/client/cli"
	servercli "github.com/TxCorpi0x/file-upload-merkle/server/cli"

	"github.com/spf13/cobra"
)

func main() {
	var rootCmd = &cobra.Command{
		Use:   "fxmerkle",
		Short: "fxmerkle is a tool for verifiable file uploads and downloads",
		Long: `A CLI for managing both fxmerkle client and server,
			and verifying downloaded files with the help of Merkle proofs.`,
		Run: func(cmd *cobra.Command, args []string) {
			if err := cmd.Help(); err != nil {
				log.Fatal(err)
			}
		},
	}

	rootCmd.AddCommand(clientcli.Cmd)
	rootCmd.AddCommand(servercli.Cmd)

	if err := rootCmd.Execute(); err != nil {
		log.Println(err)
		os.Exit(1)
	}
}
