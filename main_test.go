package main

import (
	"io/ioutil"
	"path/filepath"
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

    code, errs := CompileSource(string(src))
    for _, err := range errs {
      t.Error(err)
    }

    if len(code) == 0 {
      t.Error("expect code to be present")
    }
	}
}
