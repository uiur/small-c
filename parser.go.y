%{
package main

import (
    "go/scanner"
    "go/token"
)

type Expression interface{}
type Token struct {
    token   int
    literal string
}

type NumExpr struct {
    literal string
}
type BinOpExpr struct {
    left     Expression
    operator rune
    right    Expression
}

%}

%union {
  token Token
  expr Expression
}

%type<expr> program
%type<expr> expr
%token<token> NUMBER

%left '+'
%left '*'

%%

program
  : expr
  {
    $$ = $1
    yylex.(*Lexer).result = $$
  }

expr
  : NUMBER
  {
    $$ = NumExpr{ literal: $1.literal }
  }
  | expr '+' expr
  {
    $$ = BinOpExpr{ left: $1, operator: '+', right: $3 }
  }
  | expr '*' expr
  {
    $$ = BinOpExpr{ left: $1, operator: '*', right: $3 }
  }

%%

type Lexer struct {
    scanner.Scanner
    result Expression
}

func (l *Lexer) Lex(lval *yySymType) int {
  _, tok, lit := l.Scan()
  token_number := int(tok)

  switch tok {
  case token.INT:
    token_number = NUMBER
  case token.ADD, token.MUL:
    token_number = int(tok.String()[0])
  default:
    return 0
  }

  lval.token = Token{ token: token_number, literal: lit }

  return token_number
}

func (l *Lexer) Error(e string) {
  panic(e)
}
