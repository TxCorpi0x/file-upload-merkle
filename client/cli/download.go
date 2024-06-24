package cli

import (
	"encoding/hex"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/spf13/cobra"

	httpclient "github.com/TxCorpi0x/file-upload-merkle/client/http"
	"github.com/TxCorpi0x/file-upload-merkle/conf"
)

type Downloader interface {
	DownloadFileAt(index int, destination *os.File) error
}

var _ Downloader = (*httpclient.HttpDownloader)(nil)

var downloadCmd = &cobra.Command{
	Use:   "download",
	Short: "Download a file by index, from the server, and verify its integrity",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 1 {
			fmt.Println("Please enter one index of an uploaded file to download")

			return
		}

		index, err := strconv.Atoi(args[0])
		if err != nil || index < 1 {
			fmt.Println("The index must be a number starting from 1")

			return
		}

		rootHash, err := os.ReadFile(conf.EnvStr("MERKLE_ROOT_FILENAME", defaultMerkleRootFilename))
		if err != nil {
			fmt.Println("Merkle Root hash is missing or unreadable:", err)

			return
		}

		rootHash, err = hex.DecodeString(string(rootHash))
		if err != nil {
			fmt.Println("Error parsing root hash from file:", err)
		}

		downloader := httpclient.NewHttpDownloader(
			&http.Client{Timeout: time.Second * 30},
			conf.EnvStr("SERVER_URL", defaultServerURL),
			rootHash,
		)

		if err := downloader.DownloadFileAt(index, os.Stdout); err != nil {
			fmt.Println(err)

			return
		}
	},
}
