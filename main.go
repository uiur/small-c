package main

import (
	"fmt"
	"go/scanner"
	"go/token"
	"io/ioutil"
	"os"

	"github.com/k0kubun/pp"
)

func main() {
	data, _ := ioutil.ReadAll(os.Stdin)
	statements, err := Parse(string(data))
	if err != nil {
		fmt.Println(err)
		return
	}

	pp.Print(statements)
}

func Parse(src string) ([]Statement, error) {
	fset := token.NewFileSet()
	file := fset.AddFile("", fset.Base(), len(src))

	l := new(Lexer)
	l.Init(file, []byte(src), nil, scanner.ScanComments)

	fail := yyParse(l)
	if fail == 1 {
		return nil, l.err
	}

	return l.result, nil
}
