package controllers

import (
	"io/ioutil"
	"os"
	"path/filepath"
)

func runTest(lang, code, testCode string) error {
	var err error

	// Create temp directory
	tempDir, err := ioutil.TempDir("/tmp", "kodability-")
	if err != nil {
		return err
	}
	defer os.RemoveAll(tempDir)

	if lang == "groovy" {
		srcDir := filepath.Join(tempDir, "src")
		if err = os.Mkdir(srcDir, 0777); err != nil {
			return err
		}

		reportDir := filepath.Join(tempDir, "report")
		if err = os.Mkdir(reportDir, 0777); err != nil {
			return err
		}

		err = ioutil.WriteFile(filepath.Join(srcDir, "Example.groovy"), []byte(code), 0777)
		if err != nil {
			return err
		}
		err = ioutil.WriteFile(filepath.Join(srcDir, "TestExample.groovy"), []byte(testCode), 0777)
		if err != nil {
			return err
		}

		// TODO: process docker run report
		// Run Docker
		// cmd := exec.Command("docker", "run", "--rm",
		// 	"-v", fmt.Sprintf("%s:/test/src", srcDir),
		// 	"-v", fmt.Sprintf("%s:/test/report", reportDir),
		// 	"kodability/groovy:2",
		// )
		// err = cmd.Run()
		// if err != nil {
		// 	return err
		// }

		return nil
	}
	return nil
}
