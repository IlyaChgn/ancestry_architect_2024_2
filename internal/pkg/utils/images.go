package utils

import (
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"time"
)

const (
	filenameLen = 8
)

func WriteFile(file *multipart.FileHeader, staticDirectory, folderName string) (string, error) {
	uploadedFile, err := file.Open()
	if err != nil {
		return "", err
	}

	defer uploadedFile.Close()

	currentTime := time.Now()

	dirName := fmt.Sprintf("%s/%s/%d-%02d-%02d", staticDirectory, folderName,
		currentTime.Year(), currentTime.Month(), currentTime.Day())

	err = os.MkdirAll(dirName, os.ModePerm)
	if err != nil {
		return "", err
	}

	extension := filepath.Ext(file.Filename)
	filename := RandString(filenameLen) + extension
	fullPath := dirName + "/" + filename

	destination, err := os.Create(fullPath)
	if err != nil {
		return "", err
	}

	defer destination.Close()

	if _, err := io.Copy(destination, uploadedFile); err != nil {
		return "", err
	}

	return fullPath, nil
}
