%{
package main

import (
    "os"
    "go/scanner"
    "go/token"
    "fmt"
    "strconv"
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

  parameter_declaration ParameterDeclaration
}

%type<expression> expression optional_expression identifier_expression identifier
%type<expression> add_expression mult_expression assign_expression primary_expression logical_or_expression logical_and_expression equal_expression relation_expression unary_expression postfix_expression
%type<expressions> parameters optional_parameters
%type<statements> statements declarations optional_statements optional_declarations program
%type<statement> statement compound_statement external_declaration declaration function_definition function_prototype
%type<declarator> declarator
%type<declarators> declarators
%type<parameter_declaration> parameter_declaration
%token<token> NUMBER IDENT TYPE IF LOGICAL_OR LOGICAL_AND RETURN EQL NEQ GEQ LEQ ELSE WHILE FOR '-' '*' '&' '{'

%%

program
  : external_declaration
  {
    $$ = []Statement{$1}
    yylex.(*Lexer).result = $$
  }
  | program external_declaration
  {
    $$ = append($1, $2)
    yylex.(*Lexer).result = $$
  }

external_declaration
  : declaration
  | function_prototype
  | function_definition

declarations
  : declaration
  {
    $$ = []Statement{ $1 }
  }
  | declarations declaration
  {
    $$ = append($1, $2)
  }

declaration
  : TYPE declarators ';'
  {
    $$ = Declaration{ pos: $1.pos, VarType: $1.lit, Declarators: $2 }
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
    $$ = Declarator{ Identifier: $1 }
  }
  | identifier_expression '[' NUMBER ']'
  {
    i, _ := strconv.Atoi($3.lit)
    $$ = Declarator{ Identifier: $1, Size: i }
  }

function_prototype
  : TYPE identifier_expression '(' optional_parameters ')' ';'
  {
    $$ = FunctionDefinition{ pos: $1.pos, TypeName: $1.lit, Identifier: $2, Parameters: $4 }
  }

function_definition
  : TYPE identifier_expression '(' optional_parameters ')' compound_statement
  {
    $$ = FunctionDefinition{ pos: $1.pos, TypeName: $1.lit, Identifier: $2, Parameters: $4, Statement: $6 }
  }

identifier_expression
  : identifier
  | '*' identifier_expression
  {
    $$ = UnaryExpression{ pos: $1.pos, Operator: "*", Value: $2 }
  }

optional_parameters
  : { $$ = nil }
  | parameters

parameters
  : parameter_declaration
  {
    $$ = []Expression{ $1 }
  }
  | parameters ',' parameter_declaration
  {
    $$ = append($1, $3)
  }

parameter_declaration
  : TYPE identifier_expression
  {
    $$ = ParameterDeclaration{ pos: $1.pos, TypeName: $1.lit, Identifier: $2 }
  }

compound_statement
  : '{' optional_declarations optional_statements '}'
  {
    $$ =  CompoundStatement{ pos: $1.pos, Declarations: $2, Statements: $3 }
  }

optional_declarations
  : { $$ = nil }
  | declarations

optional_statements
  : { $$ = nil }
  | statements

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
    $$ = nil
  }
  | expression ';'
  {
    $$ = ExpressionStatement{ Value: $1 }
  }
  | compound_statement
  | IF '(' expression ')' statement
  {
    $$ = IfStatement{ pos: $1.pos, Condition: $3, TrueStatement: $5 }
  }
  | IF '(' expression ')' statement ELSE statement
  {
    $$ = IfStatement{ pos: $1.pos, Condition: $3, TrueStatement: $5, FalseStatement: $7 }
  }
  | WHILE '(' expression ')' statement
  {
    $$ = WhileStatement{ pos: $1.pos, Condition: $3, Statement: $5 }
  }
  | FOR '(' optional_expression ';' optional_expression ';' optional_expression ')' statement
  {
    $$ = ForStatement{ pos: $1.pos, Init: $3, Condition: $5, Loop: $7, Statement: $9 }
  }
  | RETURN optional_expression ';'
  {
    $$ = ReturnStatement{ pos: $1.pos, Value: $2 }
  }

optional_expression: { $$ = nil }
  | expression

expression
  : assign_expression
  {
    $$ = $1
  }
  | expression ',' assign_expression
  {
    switch e := $1.(type) {
    case ExpressionList:
      $$ = ExpressionList{ Values: append(e.Values, $3) }
    default:
      $$ = ExpressionList{ Values: []Expression{$1, $3} }
    }
  }

assign_expression
  : logical_or_expression
  | logical_or_expression '=' logical_or_expression
  {
    $$ = BinOpExpression{ Left: $1, Operator: "=", Right: $3 }
  }

logical_or_expression
  : logical_and_expression
  | logical_and_expression LOGICAL_OR logical_and_expression
  {
    $$ = BinOpExpression{ Left: $1, Operator: $2.lit, Right: $3}
  }

logical_and_expression
  : equal_expression
  | equal_expression LOGICAL_AND equal_expression
  {
    $$ = BinOpExpression{ Left: $1, Operator: $2.lit, Right: $3}
  }

equal_expression
  : relation_expression
  | relation_expression EQL relation_expression
  {
    $$ = BinOpExpression{ Left: $1, Operator: $2.lit, Right: $3}
  }
  | relation_expression NEQ relation_expression
  {
    $$ = BinOpExpression{ Left: $1, Operator: $2.lit, Right: $3}
  }

relation_expression
  : add_expression
  | add_expression '>' add_expression
  {
    $$ = BinOpExpression{ Left: $1, Operator: ">", Right: $3}
  }
  | add_expression '<' add_expression
  {
    $$ = BinOpExpression{ Left: $1, Operator: "<", Right: $3}
  }
  | add_expression GEQ add_expression
  {
    $$ = BinOpExpression{ Left: $1, Operator: $2.lit, Right: $3}
  }
  | add_expression LEQ add_expression
  {
    $$ = BinOpExpression{ Left: $1, Operator: $2.lit, Right: $3}
  }

add_expression
  : mult_expression
  | add_expression '+' mult_expression
  {
    $$ = BinOpExpression{ Left: $1, Operator: "+", Right: $3 }
  }
  | add_expression '-' mult_expression
  {
    $$ = BinOpExpression{ Left: $1, Operator: "-", Right: $3 }
  }

mult_expression
  : unary_expression
  | mult_expression '*' primary_expression
  {
    $$ = BinOpExpression{ Left: $1, Operator: "*", Right: $3 }
  }
  | mult_expression '/' primary_expression
  {
    $$ = BinOpExpression{ Left: $1, Operator: "/", Right: $3 }
  }

unary_expression
  : postfix_expression
  | '-' unary_expression
  {
    $$ = UnaryExpression{ pos: $1.pos, Operator: "-", Value: $2 }
  }
  | '&' unary_expression
  {
    $$ = UnaryExpression{ pos: $1.pos, Operator: "&", Value: $2 }
  }
  | '*' unary_expression
  {
    $$ = UnaryExpression{ pos: $1.pos, Operator: "*", Value: $2 }
  }

postfix_expression
  : primary_expression
  | postfix_expression '[' expression ']'
  {
    $$ = ArrayReferenceExpression{ Target: $1, Index: $3 }
  }
  | identifier '(' optional_expression ')'
  {
    $$ = FunctionCallExpression{ Identifier: $1, Argument: $3  }
  }

primary_expression
  : NUMBER
  {
    $$ = NumberExpression{ pos: $1.pos, Value: $1.lit }
  }
  | identifier
  | '(' expression ')'
  {
    $$ = $2
  }

identifier
  : IDENT
  {
    $$ = IdentifierExpression{ pos: $1.pos, Name: $1.lit }
  }

%%

type Lexer struct {
    scanner.Scanner
    result []Statement
    token Token
    pos token.Pos
    errMessage string
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
  l.token = lval.token

  return token_number
}

func (l *Lexer) Error(e string) {
  l.pos = l.token.pos
  l.errMessage = e
}
