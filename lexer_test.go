package main

import (
  "testing"
)

func TestLex(t *testing.T) {
  {
    l := new(Lexer)
    l.Init(`42`)

    var sym yySymType
    result := l.Lex(&sym)

    if result != NUMBER {
      t.Errorf("expect NUMBER, got %v", result)
    }
  }
}
