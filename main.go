package main

import (
	"archive/zip"
	"bufio"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"
)

var PATH_SEPARATOR string

func init() {
	if runtime.GOOS == "windows" {
		PATH_SEPARATOR = "\\"

	} else {
		PATH_SEPARATOR = "/"
	}
}

func main() {

	rootDir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	path := filepath.Join(rootDir, ".gcloudignore")

	// var path string = "/mnt/c/Users/rapha/Documents/globo/Cloud Functions/slack-bot-incident-early-warning-alerts/src/.catignore"

	// rootDir := filepath.Dir(path)

	if err := CheckCatIgnore(path); err != nil {
		log.Fatal(err)
	}

	lines, err := ReadCatIgnore(path)

	if err != nil {
		log.Fatal(err)
	}

	fileName, err := ZipCatIgnore(rootDir, lines, filepath.Join(rootDir, "tmp", "functions"))

	if err != nil {
		log.Fatal(err)
	}

	log.Printf("%s file is created successfully", fileName)

}

func CheckCatIgnore(path string) error {

	files := strings.Split(path, PATH_SEPARATOR)

	var fileName = files[len(files)-1]

	switch fileName {
	case ".catignore":
		if err := CheckFileExist(path, fileName); err != nil {
			return err
		}
		return nil
	case ".gcloudignore":
		if err := CheckFileExist(path, fileName); err != nil {
			return err
		}
		return nil

	default:
		msg := fmt.Sprintf("check if the %s file is part of the official list of supported files!", fileName)

		if err := errors.New(msg); err != nil {
			return err
		}
	}
	return nil
}

func CheckFileExist(path, fileName string) error {
	file, err := os.Stat(path)
	if os.IsNotExist(err) || file.Name() != fileName {
		msg := fmt.Sprintf("%s does not exist", fileName)

		if err := errors.New(msg); err != nil {
			return err
		}
	}

	return nil
}

func ReadCatIgnore(path string) ([]string, error) {

	var lines []string

	readFile, err := os.Open(path)

	if err != nil {
		return nil, err
	}

	fileScanner := bufio.NewScanner(readFile)
	fileScanner.Split(bufio.ScanLines)

	for fileScanner.Scan() {

		lines = append(lines, fileScanner.Text())
	}

	lines = append(lines, ".catignore")

	if err := CheckFileExist(filepath.Join(filepath.Dir(path), ".git"), ".git"); err == nil {
		lines = append(lines, ".git")
	}

	readFile.Close()

	return lines, nil

}

func ZipCatIgnore(rootDir string, lines []string, outputFilePath string) (string, error) {

	if os.MkdirAll(filepath.Dir(outputFilePath), 0666) != nil {
		panic("Unable to create directory for tagfile!")
	}

	var fileName string = outputFilePath + ".zip"

	zipfile, err := os.Create(fileName)
	if err != nil {
		return "", err
	}
	defer zipfile.Close()

	archive := zip.NewWriter(zipfile)
	defer archive.Close()

	fileList, err := ioutil.ReadDir(rootDir)
	if err != nil {
		return "", err
	}

	filterFiles := FilterFiles(fileList, lines)

	for _, source := range filterFiles {

		filePath := filepath.Join(rootDir, source)

		err = filepath.Walk(filePath, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}

			if info.IsDir() {
				if filePath == path {
					return nil
				}
				path += PATH_SEPARATOR
			}

			header, err := zip.FileInfoHeader(info)

			if err != nil {
				return err
			}

			split := strings.Split(path, rootDir+PATH_SEPARATOR)

			header.Name = split[1]

			header.Method = zip.Deflate

			writer, err := archive.CreateHeader(header)
			if err != nil {
				return err
			}

			if info.IsDir() {

				return nil
			}

			file, err := os.Open(path)
			if err != nil {
				return err
			}
			defer file.Close()
			_, err = io.Copy(writer, file)
			return err
		})
	}

	if err != nil {
		return "", err
	}
	if err = archive.Flush(); err != nil {
		return "", err
	}
	return fileName, nil
}

func FilterFiles(fileList []fs.FileInfo, lines []string) []string {
	var files []string

	for _, file := range fileList {
		files = append(files, file.Name())
	}

	for _, file := range files {

		extension := strings.Split(file, "*")

		if extension[1] != "" {
			var re = regexp.MustCompile(fmt.Sprintf(`(?i)^(.*\.((%s)$))?[^.]*$`, extension[1]))

			// if len(re.FindStringIndex(str)) > 0 {
			// 	fmt.Println(re.FindString(str), "found at index", re.FindStringIndex(str)[0])
			// }
		}

	}

	var filterFiles []string

	for _, file := range files {
		exist := false
		for _, line := range lines {
			if file == line {
				exist = true
				break
			}
		}
		if !exist {
			filterFiles = append(filterFiles, file)
		}
	}

	return filterFiles
}
