package main

import "go/token"

type Expression interface{}
type Token struct {
	tok token.Token
	lit string
	pos token.Pos
}

type NumberExpression struct {
	Value string
}

type IdentifierExpression struct {
	Name string
}

type UnaryExpression struct {
	Operator   string
	Expression Expression
}

type BinOpExpression struct {
	Left     Expression
	Operator string
	Right    Expression
}

type Declarator struct {
	Identifier Expression
	Size       string
}

type Declaration struct {
	VarType     string
	Declarators []Declarator
}

type FunctionDefinition struct {
	TypeName   string
	Identifier Expression
	Parameters []ParameterDeclaration
	Statement  Statement
}

type FunctionPrototype struct {
	TypeName   string
	Identifier Expression
	Parameters []ParameterDeclaration
}

type Statement interface{}
type CompoundStatement struct {
	Declarations []Statement
	Statements   []Statement
}

type ExpressionStatement struct {
	Value Expression
}

type IfStatement struct {
	Condition      Expression
	TrueStatement  Statement
	FalseStatement Statement
}

type WhileStatement struct {
	Condition Expression
	Statement Statement
}

type ForStatement struct {
	Init      Expression
	Condition Expression
	Loop      Expression
	Statement Statement
}

type ReturnStatement struct {
	Value Expression
}

type AssignExpression struct {
	Left  Expression
	Right Expression
}

type FunctionCallExpression struct {
	Identifier string
	Argument   Expression
}

type ArrayReferenceExpression struct {
	Target Expression
	Index  Expression
}

type ParameterDeclaration struct {
	TypeName   string
	Identifier Expression
}

type PointerExpression struct {
	Value Expression
}
