%{
package main

import (
    "os"
    "go/scanner"
    "go/token"
    "fmt"
)

%}

%union {
  token Token

  expression Expression
  expressions []Expression

  declarator Declarator
  declarators []Declarator

  statement Statement
  statements []Statement

  parameters []ParameterDeclaration
  parameter_declaration ParameterDeclaration
}

%type<expression> external_declaration declaration function_definition function_prototype identifier_expression identifier
%type<expression> expression add_expression mult_expression assign_expression primary_expression logical_or_expression logical_and_expression equal_expression relation_expression unary_expression optional_expression postfix_expression
%type<expressions> program
%type<statements> statements declarations
%type<statement> statement compound_statement
%type<declarator> declarator
%type<declarators> declarators
%type<parameters> parameters optional_parameters
%type<parameter_declaration> parameter_declaration
%type<token> unary_op
%token<token> NUMBER IDENT TYPE IF LOGICAL_OR LOGICAL_AND RETURN EQL NEQ GEQ LEQ ELSE WHILE FOR

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

declarations
  : declaration
  {
    $$ = []Statement{ $1 }
  }
  | declarations declaration
  {
    $$ = append($1, $2)
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
  : identifier_expression
  {
    $$ = Declarator{ identifier: $1 }
  }
  | identifier_expression '[' NUMBER ']'
  {
    $$ = Declarator{ identifier: $1, size: $3.lit }
  }

external_declaration
  : declaration
  | function_prototype
  | function_definition

function_prototype
  : TYPE identifier_expression '(' optional_parameters ')' ';'
  {
    $$ = FunctionPrototype{ typeName: $1.lit, identifier: $2, parameters: $4 }
  }

function_definition
  : TYPE identifier_expression '(' optional_parameters ')' compound_statement
  {
    $$ = FunctionDefinition{ typeName: $1.lit, identifier: $2, parameters: $4, statement: $6 }
  }

identifier_expression
  : identifier
  | '*' identifier_expression
  {
    $$ = PointerExpression{ expression: $2 }
  }

identifier
  : IDENT
  {
    $$ = IdentifierExpression{ name: $1.lit }
  }

optional_parameters
  : { $$ = []ParameterDeclaration{} }
  | parameters

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
  : TYPE identifier_expression
  {
    $$ = ParameterDeclaration{ typeName: $1.lit, identifier: $2 }
  }

compound_statement
  : '{' '}'
  {
    $$ = CompoundStatement{}
  }
  | '{' declarations '}'
  {
    $$ = CompoundStatement{ declarations: $2 }
  }
  | '{' statements '}'
  {
    $$ =  CompoundStatement{ statements: $2 }
  }
  | '{' declarations statements '}'
  {
    $$ =  CompoundStatement{ declarations: $2, statements: $3 }
  }

statements
  : statement
  {
    $$ = []Statement{ $1 }
  }
  | statements statement
  {
    $$ = append($1, $2)
  }

statement
  : ';'
  {
    $$ = ExpressionStatement{}
  }
  | expression ';'
  {
    $$ = ExpressionStatement{ expression: $1 }
  }
  | compound_statement
  | IF '(' expression ')' statement
  {
    $$ = IfStatement{ expression: $3, trueStatement: $5 }
  }
  | IF '(' expression ')' statement ELSE statement
  {
    $$ = IfStatement{ expression: $3, trueStatement: $5, falseStatement: $7 }
  }
  | WHILE '(' expression ')' statement
  {
    $$ = WhileStatement{ condition: $3, statement: $5 }
  }
  | FOR '(' optional_expression ';' optional_expression ';' optional_expression ')' statement
  {
    $$ = ForStatement{ init: $3, condition: $5, loop: $7, statement: $9 }
  }
  | RETURN optional_expression ';'
  {
    $$ = ReturnStatement{ expression: $1 }
  }

optional_expression: { $$ = nil }
  | expression

expression
  : assign_expression

assign_expression
  : logical_or_expression
  | logical_or_expression '=' logical_or_expression
  {
    $$ = AssignExpression{ left: $1, right: $3 }
  }

logical_or_expression
  : logical_and_expression
  | logical_and_expression LOGICAL_OR logical_and_expression
  {
    $$ = BinOpExpression{ left: $1, operator: $2.lit, right: $3}
  }

logical_and_expression
  : equal_expression
  | equal_expression LOGICAL_AND equal_expression
  {
    $$ = BinOpExpression{ left: $1, operator: $2.lit, right: $3}
  }

equal_expression
  : relation_expression
  | relation_expression EQL relation_expression
  {
    $$ = BinOpExpression{ left: $1, operator: $2.lit, right: $3}
  }
  | relation_expression NEQ relation_expression
  {
    $$ = BinOpExpression{ left: $1, operator: $2.lit, right: $3}
  }

relation_expression
  : add_expression
  | add_expression '>' add_expression
  {
    $$ = BinOpExpression{ left: $1, operator: ">", right: $3}
  }
  | add_expression '<' add_expression
  {
    $$ = BinOpExpression{ left: $1, operator: "<", right: $3}
  }
  | add_expression GEQ add_expression
  {
    $$ = BinOpExpression{ left: $1, operator: $2.lit, right: $3}
  }
  | add_expression LEQ add_expression
  {
    $$ = BinOpExpression{ left: $1, operator: $2.lit, right: $3}
  }

add_expression
  : mult_expression
  | add_expression '+' mult_expression
  {
    $$ = BinOpExpression{ left: $1, operator: "+", right: $3 }
  }
  | add_expression '-' mult_expression
  {
    $$ = BinOpExpression{ left: $1, operator: "-", right: $3 }
  }

mult_expression
  : unary_expression
  | mult_expression '*' primary_expression
  {
    $$ = BinOpExpression{ left: $1, operator: "*", right: $3 }
  }
  | mult_expression '/' primary_expression
  {
    $$ = BinOpExpression{ left: $1, operator: "/", right: $3 }
  }

unary_expression
  : postfix_expression
  | unary_op unary_expression
  {
    $$ = UnaryExpression{ operator: $1.lit, expression: $2 }
  }

unary_op
  : '-' { $$ = Token{ lit: "-" } }
  | '&' { $$ = Token{ lit: "&" } }
  | '*' { $$ = Token{ lit: "*" } }

postfix_expression
  : primary_expression
  | postfix_expression '[' expression ']'
  {
    $$ = ArrayReferenceExpression{ target: $1, index: $3 }
  }
  | IDENT '(' optional_expression ')'
  {
    $$ = FunctionCallExpression{ identifier: $1.lit, argument: $3  }
  }

primary_expression
  : NUMBER
  {
    $$ = NumberExpression{ lit: $1.lit }
  }
  | identifier
  | '(' expression ')'
  {
    $$ = $2
  }

%%

type Lexer struct {
    scanner.Scanner
    result Expression
}

var tokenMap = map[token.Token]int {
  token.LOR: LOGICAL_OR,
  token.LAND: LOGICAL_AND,
  token.IF: IF,
  token.ELSE: ELSE,
  token.RETURN: RETURN,
  token.EQL: EQL,
  token.NEQ: NEQ,
  token.GEQ: GEQ,
  token.LEQ: LEQ,
  token.FOR: FOR,
}

func identToNumber(lit string) int {
  switch lit {
  case "int", "void":
    return TYPE
  case "while":
    return WHILE
  default:
    return IDENT
  }
}

func (l *Lexer) Lex(lval *yySymType) int {
  pos, tok, lit := l.Scan()
  token_number := int(tok)

  if len(os.Getenv("DEBUG")) > 0 {
    fmt.Println(tok, lit)
  }

  if tokenMap[tok] > 0 {
    return tokenMap[tok]
  }

  switch tok {
  case token.EOF:
    return -1
  case token.INT:
    token_number = NUMBER
  case token.ADD, token.SUB, token.MUL, token.QUO, token.AND,
    token.COMMA, token.SEMICOLON,
    token.ASSIGN,
    token.GTR, token.LSS,
    token.LBRACK, token.RBRACK,
    token.LBRACE, token.RBRACE,
    token.LPAREN, token.RPAREN:
    // newline
    if tok.String() == ";" && lit != ";" {
      // read next
      return l.Lex(lval)
    }
    token_number = int(tok.String()[0])
  case token.IDENT:
    token_number = identToNumber(lit)
  default:
    return -1
  }

  lval.token = Token{ tok: tok, lit: lit, pos: pos }

  return token_number
}

func (l *Lexer) Error(e string) {
  panic(e)
}
