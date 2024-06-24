package cli

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/spf13/cobra"

	httpclient "github.com/TxCorpi0x/file-upload-merkle/client/http"
	"github.com/TxCorpi0x/file-upload-merkle/conf"
	"github.com/TxCorpi0x/file-upload-merkle/types"
)

type Uploader interface {
	UploadFilesFrom(filePaths []string) ([]types.UploadedFile, string, error)
}

var _ Uploader = (*httpclient.HttpUploader)(nil)

var uploadCmd = &cobra.Command{
	Use:   "upload",
	Short: "Upload a set of files, or an entire folder, to the server",
	Long:  "E.g. args: <file1> <file2> <file3> | args: <directory>",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			fmt.Println("Please enter a list of file paths or a directory path to upload")

			return
		}

		filePaths, err := argsToFilesToUpload(args)
		if err != nil {
			fmt.Println(err)

			return
		}

		serverURL := conf.EnvStr("SERVER_URL", defaultServerURL)
		uploader := httpclient.NewHttpUploader(&http.Client{Timeout: time.Second * 30}, serverURL)
		uploadedFiles, merkleRoot, err := uploader.UploadFilesFrom(filePaths)
		if err != nil {
			fmt.Println(err)

			return
		}

		for _, f := range uploadedFiles {
			fmt.Printf("Uploaded file at index #%d: %s\n", f.Index, f.Name)
		}

		merkleRootFilename := conf.EnvStr("MERKLE_ROOT_FILENAME", defaultMerkleRootFilename)
		if err = os.WriteFile(merkleRootFilename, []byte(merkleRoot), 0644); err != nil {
			fmt.Printf("Failed to store merkle root: %s\n", err)

			return
		}

		fmt.Println("Merkle Root hash:", merkleRoot)
	},
}

func argsToFilesToUpload(args []string) (filePaths []string, err error) {
	// Check whether the 1st arg is a directory path
	isDirectory, err := isDirectory(args[0])
	if err != nil {
		err = fmt.Errorf("error checking if %s is a directory: %v", args[0], err)

		return
	}

	if isDirectory {
		filePaths, err = listFilesInDirectory(args[0])
		if err != nil {
			err = fmt.Errorf("error listing files inside of %s: %v", args[0], err)

			return
		}
	} else {
		for _, arg := range args {
			if _, err := os.Stat(arg); err == nil {
				filePaths = append(filePaths, arg)
			}
		}
	}

	if len(filePaths) == 0 {
		err = errors.New("none of the files/dir specified can be found")
	}

	return
}

func isDirectory(path string) (bool, error) {
	fileInfo, err := os.Stat(path)
	if err != nil {
		return false, err
	}

	return fileInfo.IsDir(), nil
}

func listFilesInDirectory(directoryPath string) ([]string, error) {
	var files []string

	err := filepath.Walk(directoryPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Exclude directories
		if !info.IsDir() {
			files = append(files, path)
		}

		return nil
	})

	return files, err
}
