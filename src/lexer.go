package main

import (
	"fmt"
	"regexp"
	"strings"
	"text/scanner"
)

type Lexer struct {
	scanner    scanner.Scanner
	result     []Statement
	token      Token
	pos        scanner.Position
	errMessage string
}

func (l *Lexer) Init(code string) {
	l.scanner.Init(strings.NewReader(code))
}

var keywords = map[string]int{
	"int":    TYPE,
	"void":   TYPE,
	"if":     IF,
	"else":   ELSE,
	"return": RETURN,
	"while":  WHILE,
	"for":    FOR,
}

func (l *Lexer) Lex(lval *yySymType) int {
	tok := l.scanner.Scan()

	if tok == scanner.EOF {
		return -1
	}

	lit := l.scanner.TokenText()
	pos := l.scanner.Pos()

	lval.token = Token{lit: lit, pos: pos}
	l.token = lval.token

	if regexp.MustCompile(`^(0|[1-9][0-9]*)$`).MatchString(lit) {
		return NUMBER
	}

	if keywords[lit] != 0 {
		return keywords[lit]
	}

	two := fmt.Sprintf("%c%c", tok, l.scanner.Peek())
	operators := map[string]int{
		"==": EQL,
		"!=": NEQ,
		"<=": LEQ,
		">=": GEQ,
		"&&": LOGICAL_AND,
		"||": LOGICAL_OR,
	}

	if operators[two] != 0 {
		l.scanner.Next()
		lval.token = Token{lit: two, pos: pos}
		l.token = lval.token
		return operators[two]
	}

	if regexp.MustCompile(`^'.*'$`).MatchString(lit) {
		return CHAR
	}

	switch lit {
	case "(", ")", "{", "}", "&", ";", ",", "[", "]", "+", "-", "*", "/", "<", ">", "=":
		return int(tok)

	default:
		return IDENT
	}
}

func (l *Lexer) Error(e string) {
	l.pos = l.token.pos
	l.errMessage = e
}
