package main

import (
	"log"
	"os"
	"path/filepath"
	"runtime"

	createzip "github.com/miquelis/catignore/createZip"
	formatpath "github.com/miquelis/catignore/formatPath"
)

func init() {
	if runtime.GOOS == "windows" {
		os.Setenv("PATH_SEPARATOR", "\\")
	} else {
		os.Setenv("PATH_SEPARATOR", "/")
	}
}

func main() {

	rootDir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	path := filepath.Join(rootDir, ".catignore")

	if err := formatpath.CheckCatIgnore(path); err != nil {
		log.Fatal(err)
	}

	lines, err := formatpath.ReadCatIgnore(path)

	if err != nil {
		log.Fatal(err)
	}

	fileName, err := createzip.ZipCatIgnore(rootDir, lines, filepath.Join(rootDir, "tmp", "functions"))

	if err != nil {
		log.Fatal(err)
	}

	log.Printf("%s file is created successfully", fileName)

}
