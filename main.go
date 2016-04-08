package main

import (
	"go/scanner"
	"go/token"
	"io/ioutil"
	"os"

	"github.com/k0kubun/pp"
)

func main() {
	data, _ := ioutil.ReadAll(os.Stdin)
	expressions := Parse(string(data))
	pp.Print(expressions)
}

func Parse(src string) Expression {
	fset := token.NewFileSet()
	file := fset.AddFile("", fset.Base(), len(src))

	l := new(Lexer)
	l.Init(file, []byte(src), nil, scanner.ScanComments)
	yyParse(l)

	return l.result
}
