package utils

import (
	"archive/zip"
	"errors"
	"io"
	"os"
	"path/filepath"
)

// ZipCompose Compose source directory to ZIP archive
func ZipCompose(srcDir string, dstPath string) error {
	s, err := os.Stat(srcDir)
	if err != nil {
		return errors.New("The source directory " + srcDir + "not exists.")
	}
	if !s.IsDir() {
		return errors.New("The source directory " + srcDir + "is not a directory.")
	}
	d, _ := os.Create(dstPath)
	defer d.Close()
	w := zip.NewWriter(d)
	defer w.Close()

	var files []string

	// Get file list in directory
	err = filepath.Walk(srcDir, func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() {
			files = append(files, path)
		}
		return nil
	})

	if err != nil {
		return err
	}

	for _, file := range files {
		f, err := os.Open(file)
		info, err := f.Stat()
		if err != nil {
			return err
		}
		header, err := zip.FileInfoHeader(info)
		if err != nil {
			return err
		}
		writer, err := w.CreateHeader(header)
		if err != nil {
			return err
		}
		_, err = io.Copy(writer, f)
		err = f.Close()
		if err != nil {
			return err
		}
	}
	return nil
}
