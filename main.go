package main

import (
	"fmt"
	"go/scanner"
	"go/token"
	"os"
)

func main() {
	expressions := Parse(os.Args[1])
	fmt.Printf("%#v\n", expressions)
}

func Parse(src string) Expression {
	fset := token.NewFileSet()
	file := fset.AddFile("", fset.Base(), len(src))

	l := new(Lexer)
	l.Init(file, []byte(src), nil, scanner.ScanComments)
	yyParse(l)

	return l.result
}
