package run

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
)

type TryoutDockerRunner struct {
	TempDir    string
	TempPrefix string
}

func (r *TryoutDockerRunner) Run(lang, code, testCode string) (*JUnitReport, error) {
	// Create temp directory
	tempDir, err := ioutil.TempDir(r.TempDir, r.TempPrefix)
	if err != nil {
		return nil, err
	}
	defer os.RemoveAll(tempDir)

	if lang == "go" {
		return runGo(code, testCode, tempDir)
	}
	if lang == "groovy" {
		return runGroovy(code, testCode, tempDir)
	}
	if lang == "java" {
		return runJava(code, testCode, tempDir)
	}
	if lang == "javascript" {
		return runJavascript(code, testCode, tempDir)
	}
	if lang == "python" {
		return runPython(code, testCode, tempDir)
	}
	if lang == "scala" {
		return runScala(code, testCode, tempDir)
	}
	return nil, fmt.Errorf("Unknown language: %s", lang)
}

type DockerRunContext struct {
	Image   string
	Volumes map[string]string
}

func runDocker(ctx *DockerRunContext) error {
	args := []string{"run", "--rm"}

	for localDir, contDir := range ctx.Volumes {
		args = append(args, "-v")
		args = append(args, fmt.Sprintf("%s:%s", localDir, contDir))
	}
	args = append(args, ctx.Image)

	return exec.Command("docker", args...).Run()
}

// run go test
func runGo(code, testCode string, dir string) (*JUnitReport, error) {
	var err error

	srcDir := filepath.Join(dir, "src")
	reportDir := filepath.Join(dir, "report")
	if err = createDirs(srcDir, reportDir); err != nil {
		return nil, err
	}

	err = createCodeFiles(map[string]string{
		filepath.Join(srcDir, "example.go"):      code,
		filepath.Join(srcDir, "example_test.go"): testCode,
	})
	if err != nil {
		return nil, err
	}

	// Run docker
	err = runDocker(&DockerRunContext{
		Image: "kodability/go:1.9",
		Volumes: map[string]string{
			srcDir:    "/test/src",
			reportDir: "/test/report",
		},
	})
	if err != nil {
		return nil, err
	}

	// read report file
	return readJunitReportFromDir(reportDir)
}

// run groovy test
func runGroovy(code, testCode string, dir string) (*JUnitReport, error) {
	var err error

	srcDir := filepath.Join(dir, "src")
	reportDir := filepath.Join(dir, "report")
	if err = createDirs(srcDir, reportDir); err != nil {
		return nil, err
	}

	err = createCodeFiles(map[string]string{
		filepath.Join(srcDir, "Example.groovy"):     code,
		filepath.Join(srcDir, "TestExample.groovy"): testCode,
	})
	if err != nil {
		return nil, err
	}

	// Run docker
	err = runDocker(&DockerRunContext{
		Image: "kodability/groovy:2",
		Volumes: map[string]string{
			srcDir:    "/test/src",
			reportDir: "/test/report",
		},
	})
	if err != nil {
		return nil, err
	}

	// read report file
	return readJunitReportFromDir(reportDir)
}

// run java test
func runJava(code, testCode string, dir string) (*JUnitReport, error) {
	var err error

	srcDir := filepath.Join(dir, "src")
	reportDir := filepath.Join(dir, "report")
	if err = createDirs(srcDir, reportDir); err != nil {
		return nil, err
	}

	err = createCodeFiles(map[string]string{
		filepath.Join(srcDir, "Example.java"):     code,
		filepath.Join(srcDir, "TestExample.java"): testCode,
	})
	if err != nil {
		return nil, err
	}

	// Run docker
	err = runDocker(&DockerRunContext{
		Image: "kodability/java:8",
		Volumes: map[string]string{
			srcDir:    "/test/src",
			reportDir: "/test/report",
		},
	})
	if err != nil {
		return nil, err
	}

	// read report file
	return readJunitReportFromDir(reportDir)
}

// run javascript test
func runJavascript(code, testCode string, dir string) (*JUnitReport, error) {
	var err error

	srcDir := filepath.Join(dir, "src")
	testDir := filepath.Join(dir, "test")
	reportDir := filepath.Join(dir, "report")
	if err = createDirs(srcDir, testDir, reportDir); err != nil {
		return nil, err
	}

	err = createCodeFiles(map[string]string{
		filepath.Join(srcDir, "example.js"):      code,
		filepath.Join(testDir, "exampleSpec.js"): testCode,
	})
	if err != nil {
		return nil, err
	}

	// Run docker
	err = runDocker(&DockerRunContext{
		Image: "kodability/javascript:9",
		Volumes: map[string]string{
			srcDir:    "/test/src",
			testDir:   "/test/test",
			reportDir: "/test/report",
		},
	})
	if err != nil {
		return nil, err
	}

	// read report file
	return readJunitReportFromDir(reportDir)
}

// run python test
func runPython(code, testCode string, dir string) (*JUnitReport, error) {
	var err error

	srcDir := filepath.Join(dir, "src")
	reportDir := filepath.Join(dir, "report")
	if err = createDirs(srcDir, reportDir); err != nil {
		return nil, err
	}

	err = createCodeFiles(map[string]string{
		filepath.Join(srcDir, "example.py"):      code,
		filepath.Join(srcDir, "example_test.py"): testCode,
	})
	if err != nil {
		return nil, err
	}

	// Run docker
	err = runDocker(&DockerRunContext{
		Image: "kodability/python:3",
		Volumes: map[string]string{
			srcDir:    "/test/src",
			reportDir: "/test/report",
		},
	})
	if err != nil {
		return nil, err
	}

	// read report file
	return readJunitReportFromDir(reportDir)
}

// run scala test
func runScala(code, testCode string, dir string) (*JUnitReport, error) {
	var err error

	srcDir := filepath.Join(dir, "src")
	reportDir := filepath.Join(dir, "report")
	if err = createDirs(srcDir, reportDir); err != nil {
		return nil, err
	}

	err = createCodeFiles(map[string]string{
		filepath.Join(srcDir, "Example.scala"):     code,
		filepath.Join(srcDir, "TestExample.scala"): testCode,
	})
	if err != nil {
		return nil, err
	}

	// Run docker
	err = runDocker(&DockerRunContext{
		Image: "kodability/scala:2.12",
		Volumes: map[string]string{
			srcDir:    "/test/src",
			reportDir: "/test/report",
		},
	})
	if err != nil {
		return nil, err
	}

	// read report file
	return readJunitReportFromDir(reportDir)
}
