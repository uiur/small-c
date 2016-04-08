package main

import "go/token"

type Expression interface{}
type Token struct {
	tok token.Token
	lit string
	pos token.Pos
}

type NumExpr struct {
	lit string
}

type UnaryExpression struct {
	operator   string
	expression Expression
}

type BinOpExpression struct {
	left     Expression
	operator string
	right    Expression
}

type Declarator struct {
	identifier string
	size       string
}

type Declaration struct {
	varType     string
	declarators []Declarator
}

type FunctionDefinition struct {
	typeName   string
	identifier string
	statement  Statement
}

type Statement interface{}
type CompoundStatement struct {
	declarations []Statement
	statements   []Statement
}

type ExpressionStatement struct {
	expression Expression
}

type IfStatement struct {
	expression     Expression
	trueStatement  Statement
	falseStatement Statement
}

type WhileStatement struct {
	condition Expression
	statement Statement
}

type ForStatement struct {
	init      Expression
	condition Expression
	loop      Expression
	statement Statement
}

type ReturnStatement struct {
	expression Expression
}

type AssignExpression struct {
	left  Expression
	right Expression
}

type ParameterDeclaration struct {
	typeName   string
	identifier string
}
