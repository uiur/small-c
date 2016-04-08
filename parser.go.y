%{
package main

import (
    "go/scanner"
    "go/token"
    "fmt"
)

type Expression interface{}
type Token struct {
    tok token.Token
    lit string
    pos token.Pos
}

type NumExpr struct {
    lit string
}
type BinOpExpression struct {
    left     Expression
    operator rune
    right    Expression
}

type Declarator struct {
  identifier string
  size string
}

type Declaration struct {
  varType string
  declarators []Declarator
}

type FunctionDefinition struct {
  typeName string
  identifier string
  statements []Statement
}

type Statement struct {
  expression Expression
}

type AssignExpression struct {
  left Expression
  right Expression
}

type ParameterDeclaration struct {
  typeName string
  identifier string
}

%}

%union {
  token Token

  expr Expression
  expressions []Expression

  declarator Declarator
  declarators []Declarator

  statement Statement
  statements []Statement

  parameters []ParameterDeclaration
  parameter_declaration ParameterDeclaration
}

%type<expr> external_declaration declaration function_definition expression assign_expression primary_expression
%type<expressions> program
%type<statements> statements compound_statement
%type<statement> statement
%type<declarator> declarator
%type<declarators> declarators
%type<parameters> parameters
%type<parameter_declaration> parameter_declaration
%token<token> NUMBER IDENT TYPE

%left '+'
%left '*'
%left '='

%%

program
  : external_declaration
  {
    $$ = []Expression{$1}
    yylex.(*Lexer).result = $$
  }
  | program external_declaration
  {
    $$ = append($1, $2)
    yylex.(*Lexer).result = $$
  }

declaration
  : TYPE declarators ';'
  {
    $$ = Declaration{ varType: $1.lit, declarators: $2 }
  }

declarators
  : declarator
  {
    $$ = []Declarator{ $1 }
  }
  | declarators ',' declarator
  {
    $$ = append($1, $3)
  }

declarator
  : IDENT
  {
    $$ = Declarator{ identifier: $1.lit }
  }
  | IDENT '[' NUMBER ']'
  {
    $$ = Declarator{ identifier: $1.lit, size: $3.lit }
  }

external_declaration
  : declaration
  | function_definition

function_definition
  : TYPE IDENT '(' ')' compound_statement
  {
    $$ = FunctionDefinition{ typeName: $1.lit, identifier: $2.lit, statements: $5 }
  }
  | TYPE IDENT '(' parameters ')' compound_statement
  {
    $$ = FunctionDefinition{ typeName: $1.lit, identifier: $2.lit, statements: $6 }
  }

parameters
  : parameter_declaration
  {
    $$ = []ParameterDeclaration{ $1 }
  }
  | parameters ',' parameter_declaration
  {
    $$ = append($1, $3)
  }

parameter_declaration
  : TYPE IDENT
  {
    $$ = ParameterDeclaration{ typeName: $1.lit, identifier: $2.lit }
  }

compound_statement
  : '{' '}'
  {
    $$ = []Statement{}
  }
  | '{' statements '}'
  {
    $$ =  $2
  }

statements
  : statement
  {
    $$ = []Statement{$1}
  }
  | statements statement
  {
    $$ = append($1, $2)
  }

statement
  : ';'
  {
    $$ = Statement{}
  }
  | assign_expression ';'
  {
    $$ = Statement{ expression: $1 }
  }

assign_expression
  : expression '=' expression
  {
    $$ = AssignExpression{ left: $1, right: $3 }
  }
  ;

expression
  : primary_expression
  | expression '+' expression
  {
    $$ = BinOpExpression{ left: $1, operator: '+', right: $3 }
  }
  | expression '*' expression
  {
    $$ = BinOpExpression{ left: $1, operator: '*', right: $3 }
  }

primary_expression
  : NUMBER
  {
    $$ = NumExpr{ lit: $1.lit }
  }
  | IDENT
  {
    $$ = NumExpr{ lit: $1.lit }
  }

%%

type Lexer struct {
    scanner.Scanner
    result Expression
}

func (l *Lexer) Lex(lval *yySymType) int {
  pos, tok, lit := l.Scan()
  token_number := int(tok)

  fmt.Println(tok, lit)

  switch tok {
  case token.EOF:
    return -1
  case token.INT:
    token_number = NUMBER
  case token.ADD, token.MUL,
    token.COMMA, token.SEMICOLON,
    token.ASSIGN,
    token.LBRACK, token.RBRACK,
    token.LBRACE, token.RBRACE,
    token.LPAREN, token.RPAREN:
    // eof
    if tok.String() == ";" && lit != ";" {
      return -1
    }
    token_number = int(tok.String()[0])
  case token.IDENT:
    if lit == "int" || lit == "void" {
      token_number = TYPE
    } else {
      token_number = IDENT
    }
  default:
    return -1
  }

  lval.token = Token{ tok: tok, lit: lit, pos: pos }

  return token_number
}

func (l *Lexer) Error(e string) {
  panic(e)
}
