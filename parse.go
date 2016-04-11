package main

import (
	"fmt"
	"go/scanner"
	"go/token"
	"strings"
)

// Parse returns ast
func Parse(src string) ([]Statement, error) {
	fset := token.NewFileSet()
	file := fset.AddFile("", fset.Base(), len(src))

	l := new(Lexer)
	l.Init(file, []byte(src), nil, scanner.ScanComments)
	yyErrorVerbose = true

	fail := yyParse(l)
	if fail == 1 {
		lineNumber, columnNumber := posToLineInfo(src, int(l.pos))
		err := fmt.Errorf("%d:%d: %s", lineNumber, columnNumber, l.errMessage)

		return nil, err
	}

	return l.result, nil
}

func posToLineInfo(src string, pos int) (int, int) {
	if pos < 0 {
		panic("pos must be positive")
	}

	lineNumber := strings.Count(src[:pos], "\n") + 1

	lines := strings.Split(src[:pos], "\n")
	lastLine := lines[len(lines)-1]
	columnNumber := len(lastLine)

	return lineNumber, columnNumber
}
