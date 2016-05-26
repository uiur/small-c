package main

import (
	"testing"
)

func TestLex(t *testing.T) {
	testLex(t, `42 7 0`, []int{NUMBER, NUMBER, NUMBER})
	testLex(t, `a == 100`, []int{IDENT, EQL, NUMBER})
}

func testLex(t *testing.T, code string, tokens []int) {
	l := new(Lexer)
	l.Init(`a == 100`)

	tokenTypes := []int{IDENT, EQL, NUMBER}
	result := []int{}

	var sym yySymType
	for {
		tokenNumber := l.Lex(&sym)
		if tokenNumber == -1 {
			break
		}
		result = append(result, tokenNumber)
	}

	if len(result) != len(tokenTypes) {
		t.Errorf("expect %v tokens, got %v: %v", len(tokenTypes), len(result), result)
	}

	for i, resultType := range result {
		if resultType != tokenTypes[i] {
			t.Errorf("%v: expect %v, got %v", i, tokenTypes[i], resultType)
		}
	}
}
