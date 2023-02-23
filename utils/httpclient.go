package utils

import (
	"io"
	"net/http"
	"os"
)

// DownloadFileTo Download file and save to local file
func DownloadFileTo(uri string, savePath string) (err error) {
	res, err := http.Get(uri)
	if err != nil {
		return err
	}
	f, err := os.Create(savePath)
	if err != nil {
		return err
	}
	_, err = io.Copy(f, res.Body)

	return err
}
