package main

import (
	"io/ioutil"
	"os/exec"
	"regexp"
	"strings"
	"testing"
)

func TestSimulateExample(t *testing.T) {
	examples := [](struct {
		Filename string
		Output   string
	}){
		{"example/sum.sc", "1"},
		{"example/sum_for.sc", "45"},
		{"example/many_args.sc", "6"},
		{"example/factorial.sc", "24"},
		{"example/fib.sc", "89"},
		{"example/global_var.sc", "11"},
		{"example/if_test.sc", ""},
		{"example/pointer_test.sc", ""},
		{"example/optimize_constant.sc", "1"},
		{"example/bubble_sort.sc", "12345678"},
		{"example/quick_sort.sc", "12345678"},
		{"example/putchar.sc", "hello world"},
	}

	for _, example := range examples {
		sourceFilename := example.Filename
		filename := regexp.MustCompile("\\.sc$").ReplaceAllString(sourceFilename, ".s")

		{
			err := compileAndSave(sourceFilename)

			if err != nil {
				t.Errorf("%v: %v", sourceFilename, err)
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
			t.Errorf("`%v`: expect `%v`, got `%v`", filename, expected, output)
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
