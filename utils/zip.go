package utils

import (
	"archive/zip"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

// Zip 将路径下的文件压缩成zip文件，并保持原有的目录结构
func Zip(source, target string) error {
	// 创建目标ZIP文件
	zipFile, err := os.Create(target)
	if err != nil {
		return err
	}
	defer zipFile.Close()

	zipWriter := zip.NewWriter(zipFile)
	defer zipWriter.Close()

	// 遍历源目录下的文件和子目录
	err = filepath.Walk(source, func(filePath string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// 创建ZIP文件中的文件头
		zipHeader, err := zip.FileInfoHeader(info)
		if err != nil {
			return err
		}

		// 设置ZIP文件头的文件名
		zipHeader.Name, _ = filepath.Rel(source, filePath)

		// 如果是目录，确保以"/"结尾
		if info.IsDir() {
			zipHeader.Name += "/"
		}

		// 创建ZIP文件中的文件
		zipEntry, err := zipWriter.CreateHeader(zipHeader)
		if err != nil {
			return err
		}

		// 如果不是目录，将文件内容拷贝到ZIP文件中
		if !info.IsDir() {
			file, err := os.Open(filePath)
			if err != nil {
				return err
			}
			defer file.Close()
			_, err = io.Copy(zipEntry, file)
			if err != nil {
				return err
			}
		}

		return nil
	})

	if err != nil {
		return err
	}

	fmt.Println("目录已成功压缩到", target)
	return nil
}

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
