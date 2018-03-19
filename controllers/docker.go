package controllers

import "os"

// createTempDir create a temporary directory.
func createTempDir(prefix string) (name string, err error) {
	return ioutils.createTempDir("", prefix)
}

// deleteDir removes dir and any children it contains.
func deleteDir(dir string) error {
	return os.RemoveAll(dir)
}
