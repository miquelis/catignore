/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package main

import (
	"log"
	"os"
	"path/filepath"

	createzip "github.com/miquelis/catignore/createZip"
)

func main() {
	rootDir, err := os.Getwd()

	if err != nil {
		log.Fatal(err)
	}

	path := filepath.Join(rootDir, ".catignore")

	msg, err := createzip.CreateZipFile(path, filepath.Join(rootDir, "tmp", "functions"))

	if err != nil {
		log.Fatal(err)
	}

	log.Println(msg)

}
