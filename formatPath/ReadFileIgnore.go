package formatpath

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

// ReadCatIgnore - Reads the .ignore file and adds the path of the files
func ReadCatIgnore(path string) ([]string, error) {

	var pathIgnore []string

	var rootDir string = filepath.Dir(path)

	readFile, err := os.Open(path)

	if err != nil {
		return nil, err
	}

	fileScanner := bufio.NewScanner(readFile)
	fileScanner.Split(bufio.ScanLines)

	for fileScanner.Scan() {

		if fileScanner.Text() != "" {
			line := filepath.Join(rootDir, fileScanner.Text())

			extension := strings.SplitAfter(line, "*")
			var extensionInt = len(extension)

			if extensionInt > 1 && extension[1] != "" {

				re, _ := regexp.Compile(fmt.Sprintf(`([a-zA-Z0-9\s_\\.\-\(\):])+(%s)$`, extension[1]))

				files, err := ListFiles(rootDir)

				if err != nil {
					return nil, err
				}

				for _, path := range files {
					if re.MatchString(path) {
						pathIgnore = append(pathIgnore, path)
					}
				}

			} else {
				file, err := SearchFiles(line)

				if err != nil {
					return nil, err
				}
				pathIgnore = append(pathIgnore, file...)
			}

		}
	}

	pathIgnore = DefaultIgnore(pathIgnore, rootDir)

	readFile.Close()

	return pathIgnore, nil

}

func DefaultIgnore(pathIgnore []string, rootDir string) []string {

	listIgnore := []string{
		".catignore",
		".gcloudignore",
		".gitignore",
		"__debug_bin",
		"LICENSE",
		"README.md",
	}

	fileDir, _ := SearchFiles(".vscode")
	listIgnore = append(listIgnore, fileDir...)

	for _, ignoreFile := range listIgnore {
		pathIgnore = append(pathIgnore, filepath.Join(rootDir, ignoreFile))
	}

	return pathIgnore
}
