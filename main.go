package main

import (
	"log"
	"os"
	"path/filepath"
	"runtime"

	createzip "github.com/miquelis/catignore/createZip"
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

	msg, err := createzip.CreateZipFile(path)

	if err != nil {
		log.Fatal(err)
	}

	log.Println(msg)

}
