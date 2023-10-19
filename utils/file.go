package utils

import (
	"fmt"
	"io"
	"log"
	"os"
)

// IsDir Path is directory
func IsDir(fileAddr string) bool {
	s, err := os.Stat(fileAddr)
	if err != nil {
		log.Println(err)
		return false
	}
	return s.IsDir()
}

// CreateDir creates a directory
func CreateDir(dir string) bool {
	err := os.MkdirAll(dir, os.ModePerm)
	if err != nil {
		log.Println(err)
		return false
	}
	return true
}

// IsExists Path is exists
func IsExists(path string) bool {
	_, err := os.Stat(path)
	if err == nil {
		return true
	}
	return false
}

// Copy Path to Path
func Copy(srcFile, dstFile string) error {
	srcStat, err := os.Stat(srcFile)
	if err != nil {
		return err
	}

	if !srcStat.Mode().IsRegular() {
		return fmt.Errorf("%s is not a regular file", srcFile)
	}

	source, err := os.Open(srcFile)
	if err != nil {
		return err
	}
	defer source.Close()

	dest, err := os.Create(dstFile)
	if err != nil {
		return err
	}
	defer dest.Close()

	_, err = io.Copy(dest, source)
	if err != nil {
		return err
	}
	return nil
}

// Remove Path
func Remove(path string) bool {
	if !IsExists(path) {
		return false
	}
	err := os.RemoveAll(path)
	if err != nil {
		log.Println(err)
		return false
	}
	return true
}

// Clear Path 清空路径下的所有文件
func Clear(dir string) bool {
	if IsDir(dir) {
		err := os.RemoveAll(dir)
		if err != nil {
			log.Println(err)
			return false
		}
		return CreateDir(dir)
	}
	log.Printf("%s is not a directory\n", dir)
	return false
}
