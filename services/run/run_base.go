package run

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/kodability/apiserver/models"
)

// JUnitTestcaseResult represents the result of each testcase.
type JUnitTestcaseResult struct {
	Name  string
	Time  float64
	Error string
}

// JUnitReport represents the result of each tryout.
type JUnitReport struct {
	Tests       int
	Errors      int
	Failures    int
	Timestamp   time.Time
	ElapsedTime float64
	TestResults []JUnitTestcaseResult
}

func (r JUnitReport) ToTryoutResult() models.TryoutResult {
	var errorNames, failureNames []string
	for _, result := range r.TestResults {
		if result.Error != "" {
			errorNames = append(errorNames, result.Name)
		}
	}
	return models.TryoutResult{
		TestCount:    r.Tests,
		ErrorCount:   r.Errors,
		ErrorNames:   strings.Join(errorNames, ","),
		FailureCount: r.Failures,
		FailureNames: strings.Join(failureNames, ","),
		ElapsedTime:  r.ElapsedTime,
	}
}

// Read JUnit report xml file
func readJunitReport(xmlFile string) (*JUnitReport, error) {
	// TODO:
	return nil, errors.New("readJunitReport() is not implemented")
}

// Find and read JUnit report file from the directory
func readJunitReportFromDir(dir string) (*JUnitReport, error) {
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	for _, f := range files {
		if filepath.Ext(f.Name()) == ".xml" {
			return readJunitReport(f.Name())
		}
	}
	return nil, fmt.Errorf("Junit report file not found in : %s", dir)
}

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
