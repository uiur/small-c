package main

import (
	"encoding/xml"
	"fmt"
	"go/scanner"
	"go/token"
	"io/ioutil"
	"os"
)

func main() {
	data, _ := ioutil.ReadAll(os.Stdin)
	expressions := Parse(string(data))
	buf, _ := xml.MarshalIndent(expressions, "", "  ")
	fmt.Println(string(buf))
	// fmt.Printf("%#v\n", expressions)
}

func Parse(src string) Expression {
	fset := token.NewFileSet()
	file := fset.AddFile("", fset.Base(), len(src))

	l := new(Lexer)
	l.Init(file, []byte(src), nil, scanner.ScanComments)
	yyParse(l)

	return l.result
}
