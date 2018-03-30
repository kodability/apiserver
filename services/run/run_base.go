package run

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"

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
func ReadJunitReportFile(xmlFile string, useTestSuites bool) (*JUnitReport, error) {
	// Read XML file to bytes
	bytes, err := ioutil.ReadFile(xmlFile)
	if err != nil {
		return nil, err
	}
	return ReadJunitReportBytes(bytes, useTestSuites)
}

func ReadJunitReportBytes(bytes []byte, useTestSuites bool) (*JUnitReport, error) {
	type JUnitProperty struct {
		Name  string `xml:"name,attr"`
		Value string `xml:"value,attr"`
	}

	type JUnitError struct {
		XMLName xml.Name `xml:"error"`
		Message string   `xml:"message,attr"`
		Type    string   `xml:"type,attr"`
		Output  string   `xml:",chardata"`
	}

	type JUnitFailure struct {
		XMLName xml.Name `xml:"failure"`
		Message string   `xml:"message,attr"`
		Type    string   `xml:"type,attr"`
		Output  string   `xml:",chardata"`
	}

	type JUnitTestcase struct {
		XMLName   xml.Name      `xml:"testcase"`
		ClassName string        `xml:"classname,attr"`
		Name      string        `xml:"name,attr"`
		Time      float64       `xml:"time,attr"`
		Error     *JUnitError   `xml:"error,omitempty"`
		Failure   *JUnitFailure `xml:"failure,omitempty"`
	}

	type JUnitTestSuite struct {
		XMLName    xml.Name        `xml:"testsuite"`
		Tests      int             `xml:"tests,attr"`
		Failures   int             `xml:"failures,attr,omitempty"`
		Errors     int             `xml:"errors,attr,omitempty"`
		Skipped    int             `xml:"skipped,attr,omitempty"`
		Time       float64         `xml:"time,attr,omitempty"`
		Name       string          `xml:"name,attr,omitempty"`
		Properties []JUnitProperty `xml:"properties>property,omitempty"`
		Testcases  []JUnitTestcase `xml:"testcase,omitempty"`
	}

	type JUnitTestSuites struct {
		XMLName xml.Name         `xml:"testsuites"`
		Suites  []JUnitTestSuite `xml:"testsuite,omitempty"`
	}

	// Parse XML
	testSuite := JUnitTestSuite{}
	if useTestSuites == true {
		testSuites := JUnitTestSuites{}
		if err := xml.Unmarshal(bytes, &testSuites); err != nil {
			return nil, err
		}
		if len(testSuites.Suites) >= 1 {
			testSuite = testSuites.Suites[0]
		}
	} else {
		if err := xml.Unmarshal(bytes, &testSuite); err != nil {
			return nil, err
		}
	}

	// Create JUnitReport instance
	testResults := []JUnitTestcaseResult{}
	for _, testcase := range testSuite.Testcases {
		errorMsg := ""
		if testcase.Error != nil {
			errorMsg = testcase.Error.Message
		}
		if testcase.Failure != nil {
			errorMsg = testcase.Failure.Message
		}
		testResults = append(testResults, JUnitTestcaseResult{
			Name:  testcase.Name,
			Time:  testcase.Time,
			Error: errorMsg,
		})
	}
	report := JUnitReport{
		Tests:       testSuite.Tests,
		Errors:      testSuite.Errors,
		Failures:    testSuite.Failures,
		ElapsedTime: testSuite.Time,
		TestResults: testResults,
	}

	return &report, nil
}

// Find and read JUnit report file from the directory
func ReadJunitReportFromDir(dir string, testSuites bool) (*JUnitReport, error) {
	filenames, err := ListFilenames(dir, ".xml")
	if err != nil {
		return nil, err
	}

	if len(filenames) == 0 {
		return nil, fmt.Errorf("Junit report file not found in : %s", dir)
	}

	return ReadJunitReportFile(filenames[0], testSuites)
}

func ListFilenames(dir, suffix string) ([]string, error) {
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	var filenames []string
	for _, f := range files {
		if filepath.Ext(f.Name()) == suffix {
			filenames = append(filenames, f.Name())
		}
	}
	return filenames, nil
}
