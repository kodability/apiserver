package controllers

import "os"
import "io/ioutil"

// createTempDir create a temporary directory.
func createTempDir(prefix string) (name string, err error) {
	return ioutil.TempDir("", prefix)
}

// deleteDir removes dir and any children it contains.
func deleteDir(dir string) error {
	return os.RemoveAll(dir)
}
