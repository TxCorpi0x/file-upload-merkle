package http

import "errors"

var (
	errFailedDownload  = errors.New("failed to download file")
	errFailedProveHash = errors.New("failed to prove hash")
	errFailedUpload    = errors.New("failed to upload files")
)
