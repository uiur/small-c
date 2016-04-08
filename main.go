package main

import (
	"fmt"
	"go/scanner"
	"go/token"
	"os"
)

func main() {
	Parse(os.Args[1])
}

func Parse(src string) {
	fset := token.NewFileSet()
	file := fset.AddFile("", fset.Base(), len(src))

	l := new(Lexer)
	l.Init(file, []byte(src), nil, scanner.ScanComments)
	yyParse(l)
	fmt.Printf("%#v\n", l.result)
}
