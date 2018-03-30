package run

import (
	"io/ioutil"
	"os"
)

// create a directory
func createDir(dir string) error {
	return os.Mkdir(dir, 0777)
}

// create multiple directories
func createDirs(dirs ...string) error {
	for _, dir := range dirs {
		if err := createDir(dir); err != nil {
			return err
		}
	}
	return nil
}

// create source file
func createCodeFile(file, code string) error {
	return ioutil.WriteFile(file, []byte(code), 0777)
}

// create multiple source files
func createCodeFiles(fileCodeMap map[string]string) error {
	for file, code := range fileCodeMap {
		if err := createCodeFile(file, code); err != nil {
			return err
		}
	}
	return nil
}
