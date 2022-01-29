package formatpath

import (
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

// ListFiles - List all files in the main directory
func ListFiles(rootDir string) ([]string, error) {
	files, err := SearchFiles(rootDir)

	if err != nil {
		return nil, err
	}

	files = RemoveGitTmpDir(files)

	return files, nil
}

// FilterFiles - Performs a filter removing all files that were added in .ignore
func FilterFiles(rootDir string, pathIgnore []string) ([]string, error) {

	var filterFiles []string

	files, err := ListFiles(rootDir)

	if err != nil {
		return nil, err
	}

	//remove the paths added in the ignore file
	for _, file := range files {
		exist := false
		for _, line := range pathIgnore {
			if file == line {
				exist = true
				break
			}
		}
		if !exist {
			filterFiles = append(filterFiles, file)
		}
	}

	return filterFiles, nil
}

//SearchFiles - Search the main directory to see if the file exists
func SearchFiles(sourceDir string) ([]string, error) {

	var list []string
	err := filepath.Walk(sourceDir, func(path string, info os.FileInfo, err error) error {

		files := strings.Split(path, os.Getenv("PATH_SEPARATOR"))

		var fileName = files[len(files)-1]

		if err := CheckFileExist(path, fileName, ""); err != nil {
			return err
		}

		if info.IsDir() {
			if sourceDir == path {
				return nil
			}
			path += os.Getenv("PATH_SEPARATOR")
		} else {
			list = append(list, path)
		}

		return err
	})

	if err != nil {
		return nil, err
	}

	return list, nil
}

//RemoveGitTmpDir - Remove all files from hidden .git directory and tmp directory
func RemoveGitTmpDir(listPaths []string) []string {

	var list []string
	re, _ := regexp.Compile(`(?m)(?:^|\W)(.git|tmp)(?:$|\W)`)

	for _, path := range listPaths {

		if !re.MatchString(path) {
			list = append(list, path)
		}
	}

	return list

}

// RemoveObjectFromSlice - Remove the element at index i from slice.
func RemoveObjectFromSlice(index int, slice []string) []string {

	copy(slice[index:], slice[index+1:]) // Shift slice[i+1:] left one index.
	slice[len(slice)-1] = ""             // Erase last element (write zero value).
	slice = slice[:len(slice)-1]         // Truncate slice.

	return slice
}
