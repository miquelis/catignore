package createzip

import (
	"archive/zip"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	formatpath "github.com/miquelis/catignore/formatPath"
)

func init() {
	if runtime.GOOS == "windows" {
		os.Setenv("PATH_SEPARATOR", "\\")
	} else {
		os.Setenv("PATH_SEPARATOR", "/")
	}
}

// CreateZipFile - Start creating the .zip file
func CreateZipFile(path, outputPath string) (string, error) {

	var rootDir string = filepath.Dir(path)

	if err := formatpath.CheckCatIgnore(path, rootDir); err != nil {
		return "", err
	}

	pathIgnore, err := formatpath.ReadCatIgnore(path)

	if err != nil {
		return "", err
	}

	fileName, err := ZipCatIgnore(rootDir, pathIgnore, outputPath)

	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%s file is created successfully", fileName), nil
}

// ZipCatIgnore - create the zipped file
func ZipCatIgnore(rootDir string, pathIgnore []string, outputFilePath string) (string, error) {

	// create directory if it doesn't exist
	if os.MkdirAll(filepath.Dir(outputFilePath), 0666) != nil {
		log.Fatal("Unable to create directory for tagfile!")
	}

	var fileName string = outputFilePath + ".zip"

	zipfile, err := os.Create(fileName)
	if err != nil {
		return "", err
	}
	defer zipfile.Close()

	archive := zip.NewWriter(zipfile)
	defer archive.Close()

	filterFiles, err := formatpath.FilterFiles(rootDir, pathIgnore)

	if err != nil {
		return "", err
	}

	for _, source := range filterFiles {

		if err := filepath.Walk(source, func(path string, info os.FileInfo, err error) error {

			if err != nil {
				return err
			}

			if info.IsDir() {
				return nil
			}

			header, err := zip.FileInfoHeader(info)

			if err != nil {
				return err
			}

			fileName := strings.Split(path, rootDir+os.Getenv("PATH_SEPARATOR"))[1]

			header.Name = fileName
			header.Method = zip.Deflate

			writer, err := archive.CreateHeader(header)

			if err != nil {
				return err
			}

			file, err := os.Open(path)

			if err != nil {
				return err

			}
			defer file.Close()

			_, err = io.Copy(writer, file)

			return err
		}); err != nil {
			return "", err
		}
	}

	if err = archive.Flush(); err != nil {
		return "", err
	}

	return fileName, nil
}
