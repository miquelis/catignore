package formatpath

import (
	"errors"
	"fmt"
	"os"
	"strings"
)

// CheckCatIgnore - Checks if .ignore files exist
func CheckCatIgnore(path, rootDir string) error {

	files := strings.Split(path, rootDir+os.Getenv("PATH_SEPARATOR"))

	var fileName = files[len(files)-1]

	var msg string = `File "%s" does not exist! Make sure the file is in the corresponding directory`

	switch fileName {
	case ".catignore":

		msgErro := fmt.Sprintf(msg, fileName)

		if err := CheckFileExist(path, fileName, msgErro); err != nil {
			return err
		}
		return nil
	case ".gcloudignore":

		msgErro := fmt.Sprintf(msg, fileName)

		if err := CheckFileExist(path, fileName, msgErro); err != nil {
			return err
		}
		return nil

	default:
		msg := fmt.Sprintf("check if the %s file is part of the official list of supported files!", fileName)

		return errors.New(msg)
	}

}

//CheckFileExist - Checks if a file exists
func CheckFileExist(path, fileName, msgErro string) error {
	file, err := os.Stat(path)
	if os.IsNotExist(err) || file.Name() != fileName {
		if msgErro != "" {
			return errors.New(msgErro)
		} else {
			msg := fmt.Sprintf(`File or folder named "%s" does not exist! Removing your .ignore file!`, fileName)
			return errors.New(msg)
		}

	}
	return nil
}
