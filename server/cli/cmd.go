package cli

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/spf13/cobra"

	"github.com/TxCorpi0x/file-upload-merkle/conf"
	"github.com/TxCorpi0x/file-upload-merkle/server"
	"github.com/TxCorpi0x/file-upload-merkle/storage"
)

const (
	defaultPort = 8080
)

var Cmd = &cobra.Command{
	Use:   "server",
	Short: "The fxmerkle server exposes a HTTP API for verifiable files upload & download",
	Run: func(cmd *cobra.Command, args []string) {
		//repository := storage.NewInMemoryStorage()
		repository := storage.NewInMemoryStorage()

		r := mux.NewRouter()
		r.HandleFunc("/upload", server.NewUploadHandler(repository))
		r.HandleFunc("/download/{index}", server.NewDownloadHandler(repository))
		r.HandleFunc("/proof/{index}", server.NewProofHandler(repository))

		port := conf.EnvInt("PORT", defaultPort)
		log.Println("fxmerkle server started on port", port)
		if err := http.ListenAndServe(fmt.Sprintf(":%d", port), r); err != nil {
			log.Fatal(err)
		}
	},
}
