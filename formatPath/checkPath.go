package formatpath

import (
	"errors"
	"fmt"
	"os"
	"strings"
)

// CheckCatIgnore - Checks if .ignore files exist
func CheckCatIgnore(path string) error {

	files := strings.Split(path, os.Getenv("PATH_SEPARATOR"))

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

//CheckFileExist - Checks if a file exists
func CheckFileExist(path, fileName string) error {
	file, err := os.Stat(path)
	if os.IsNotExist(err) || file.Name() != fileName {
		msg := fmt.Sprintf(`File or folder named "%s" does not exist! Removing your .ignore file!`, fileName)

		if err := errors.New(msg); err != nil {
			return err
		}
	}

	return nil
}