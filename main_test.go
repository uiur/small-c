package main

import (
	"io/ioutil"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"testing"
)

func TestCompileExample(t *testing.T) {
	files, err := filepath.Glob("example/*.sc")
	if err != nil {
		panic(err)
	}

	for _, file := range files {
  	src, err := ioutil.ReadFile(file)

    if err != nil {
      t.Error(err)
      return
    }

    code, errs := CompileSource(string(src), true)
    for _, err := range errs {
      t.Error(err)
    }

    if len(code) == 0 {
      t.Error("expect code to be present")
    }
	}
}

func TestSimulateExample(t *testing.T) {
	examples := [](struct {
		Filename string
		Output string
	}){
		{"example/sum.sc", "120"},
		{"example/bubble_sort.sc", "12345678"},
		{"example/quick_sort.sc", "12345678"},
	}

	for _, example := range examples {
		sourceFilename := example.Filename
		filename := regexp.MustCompile("\\.sc$").ReplaceAllString(sourceFilename, ".s")

		{
			err := compileAndSave(sourceFilename)

			if err != nil {
				t.Error(err)
				continue
			}
		}

		byteOut, err := exec.Command("spim", "-file", filename).Output()

		if err != nil {
			t.Error(err)
			continue
		}

		lines := strings.Split(string(byteOut), "\n")
		output := lines[len(lines)-1]
		expected := example.Output

		if output != expected {
			t.Errorf("expect `%v`'s output to be `%v`, got `%v`", filename, expected, output)
		}
	}
}

func compileAndSave(filename string) error {
	src, err := ioutil.ReadFile(filename)
  if err != nil {
    return err
  }

  code, errs := CompileSource(string(src), true)
  for _, err := range errs {
		return err
  }

	dest := regexp.MustCompile("\\.sc$").ReplaceAllString(filename, ".s")
	err = ioutil.WriteFile(dest, []byte(code), 0777)
	if err != nil {
		return err
	}

	return nil
}
