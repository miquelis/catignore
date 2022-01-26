package createzip

import (
	"archive/zip"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"

	formatpath "github.com/miquelis/catignore/formatPath"
)

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
