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
	staticDirectory = "./uploads"
	quality         = 90
	filenameLen     = 8
	startX          = 0
	startY          = 0
	endX            = 215
	endY            = 295
)

func WriteFile(file *multipart.FileHeader, folderName string) (string, error) {
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
	fullpath := dirName + "/" + filename

	destination, err := os.Create(fullpath)
	if err != nil {
		return "", err
	}

	defer destination.Close()

	if _, err := io.Copy(destination, uploadedFile); err != nil {
		return "", err
	}

	return fullpath, nil
}
