package main

import (
	"archive/zip"
	"bufio"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
)

func CheckCatIgnore(path, fileName string) error {
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

	// for fileScanner.Scan() {

	// 	lines = append(lines, filepath.Join(filepath.Dir(path), fileScanner.Text()))
	// }

	// lines = append(lines, filepath.Join(filepath.Dir(path), ".catignore"))
	// lines = append(lines, filepath.Join(filepath.Dir(path), ".git/"))

	for fileScanner.Scan() {

		lines = append(lines, fileScanner.Text())
	}

	lines = append(lines, ".catignore")

	if err := CheckCatIgnore(filepath.Join(filepath.Dir(path), ".git"), ".git"); err == nil {
		lines = append(lines, ".git")
	}

	readFile.Close()

	return lines, nil

}

func ZipCatIgnore(rootDir string, lines []string, outputFilePath string) error {

	if os.MkdirAll(filepath.Dir(outputFilePath), 0666) != nil {
		panic("Unable to create directory for tagfile!")
	}

	var fileName string = outputFilePath + ".zip"

	zipfile, err := os.Create(fileName)
	if err != nil {
		return err
	}
	defer zipfile.Close()

	archive := zip.NewWriter(zipfile)
	defer archive.Close()

	// base := filepath.Base()

	fileInfo, err := ioutil.ReadDir(rootDir)
	if err != nil {
		return err
	}
	var files []string

	// for _, file := range fileInfo {
	// 	files = append(files, filepath.Join(rootDir, file.Name()))
	// }

	for _, file := range fileInfo {
		files = append(files, file.Name())
	}

	result := make([]string, 0, 11)
	for _, file := range files {
		exist := false
		for _, line := range lines {
			if file == line {
				exist = true
				break
			}
		}
		if !exist {
			result = append(result, file)
		}
	}

	for _, source := range result {

		// fileInfo, err := os.Stat(source)
		// if err != nil {
		// 	log.Fatal(err)
		// }

		// if fileInfo.IsDir() {

		// 	err = filepath.Walk(source, func(path string, info os.FileInfo, err error) error {
		// 		fmt.Println(path)
		// 		return err
		// 	})

		// 	if err != nil {
		// 		log.Fatal(err)
		// 	}

		// 	// files, err := os.ReadDir(path)

		// 	// fmt.Print(files)

		// }
		// base := filepath.Base(source)
		filePath := filepath.Join(rootDir, source)

		err = filepath.Walk(filePath, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}

			if info.IsDir() {
				if filePath == path {
					return nil
				}
				path += "/"
			}

			header, err := zip.FileInfoHeader(info)

			if err != nil {
				return err
			}
			// filex, err := os.Stat(path)
			// if err != nil {
			// 	return err
			// }

			split := strings.Split(path, rootDir+"/")

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
		return err
	}
	if err = archive.Flush(); err != nil {
		return err
	}
	return nil
}

// file, err := os.Open(source)
// if err != nil {
// 	return err
// }
// var info os.FileInfo

// header, err := zip.FileInfoHeader(info)
// if err != nil {
// 	return err
// }
// 	fh := new(zip.FileHeader)
// 	writer, err := archive.CreateHeader(fh)
// 	if err != nil {
// 		return err
// 	}

// 	defer file.Close()
// 	_, err = io.Copy(writer, file)

// 	if err != nil {
// 		log.Fatal(err)
// 	}
// }

// if err = archive.Flush(); err != nil {
// 	return err
// }
// return nil
// }

// func zipit(source, target string) error {
// 	zipfile, err := os.Create(target)
// 	if err != nil {
// 		return err
// 	}
// 	defer zipfile.Close()

// 	archive := zip.NewWriter(zipfile)
// 	defer archive.Close()

// 	base := filepath.Base(source)

// 	err = filepath.Walk(source, func(path string, info os.FileInfo, err error) error {
// 		if err != nil {
// 			return err
// 		}

// 		if info.IsDir() {
// 			if source == path {
// 				return nil
// 			} else if fmt.Sprintf("%s/oi", source) == path {
// 				return nil
// 			}
// 			path += "/"
// 		}

// 		// if fmt.Sprintf("%s/oi/*", source) == path {
// 		// 	return nil
// 		// }

// 		header, err := zip.FileInfoHeader(info)
// 		if err != nil {
// 			return err
// 		}
// 		header.Name = path[len(base)+1:]
// 		header.Method = zip.Deflate

// 		writer, err := archive.CreateHeader(header)
// 		if err != nil {
// 			return err
// 		}

// 		if info.IsDir() {
// 			return nil
// 		}

// 		file, err := os.Open(path)
// 		if err != nil {
// 			return err
// 		}

// 		defer file.Close()
// 		_, err = io.Copy(writer, file)
// 		return err
// 	})
// 	if err != nil {
// 		return err
// 	}
// 	if err = archive.Flush(); err != nil {
// 		return err
// 	}
// 	return nil
// }

func main() {

	rootDir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	path := filepath.Join(rootDir, ".catignore")

	// var path string = "/mnt/c/Users/rapha/Documents/globo/Cloud Functions/slack-bot-incident-early-warning-alerts/src/.catignore"

	// rootDir := filepath.Dir(path)

	if err := CheckCatIgnore(path, ".catignore"); err != nil {
		log.Fatal(err)
	}

	lines, err := ReadCatIgnore(path)

	if err != nil {
		log.Fatal(err)
	}

	if err := ZipCatIgnore(rootDir, lines, filepath.Join(rootDir, "tmp", "functions")); err != nil {
		log.Fatal(err)
	}

}
