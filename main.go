package main

import (
	"fmt"

	"io/ioutil"
	"os"

	"github.com/k0kubun/pp"
)

func main() {
	var src string

	if len(os.Args) > 1 {
		filename := os.Args[1]
		data, err := ioutil.ReadFile(filename)
		if err != nil {
			fmt.Println(err)
		}

		src = string(data)
	} else {
		data, _ := ioutil.ReadAll(os.Stdin)
		src = string(data)
	}

	statements, err := Parse(src)
	if err != nil {
		fmt.Println(err)
		return
	}

	for i, statement := range statements {
		statements[i] = Walk(statement)
	}

	pp.Println(statements)
}

func Exit(src string, errs []error) {
	for _, err := range errs {
		fmt.Fprintln(os.Stderr, err)
	}

	os.Exit(1)
}
