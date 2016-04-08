package main

import (
	"fmt"
	"go/scanner"
	"go/token"
	"os"
)

func main() {
	src := []byte(os.Args[1])
	fset := token.NewFileSet()
	file := fset.AddFile("", fset.Base(), len(src))

	l := new(Lexer)
	l.Init(file, src, nil, scanner.ScanComments)
	yyParse(l)
	fmt.Printf("%#v\n", l.result)
}
