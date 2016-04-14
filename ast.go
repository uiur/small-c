package main

import "go/token"

type Token struct {
	tok token.Token
	lit string
	pos token.Pos
}

type Node interface {
	Pos() token.Pos
}

type Expression interface {
	Node
}

type ExpressionList struct {
	Values []Expression
}

func (e *ExpressionList) Pos() token.Pos {
	first := e.Values[0]
	return first.Pos()
}

type NumberExpression struct {
	pos   token.Pos
	Value string
}

func (e *NumberExpression) Pos() token.Pos { return e.pos }

type IdentifierExpression struct {
	pos    token.Pos
	Name   string
	Symbol *Symbol
}

func (e *IdentifierExpression) Pos() token.Pos { return e.pos }

type UnaryExpression struct {
	pos      token.Pos
	Operator string
	Value    Expression
}

func (e *UnaryExpression) Pos() token.Pos { return e.pos }

type BinOpExpression struct {
	Left     Expression
	Operator string
	Right    Expression
}

func (e *BinOpExpression) Pos() token.Pos {
	return e.Left.Pos()
}

func (e *BinOpExpression) IsArithmetic() bool {
	return e.Operator == "+" || e.Operator == "-" || e.Operator == "/" || e.Operator == "*"
}

type FunctionCallExpression struct {
	Identifier Expression
	Argument   Expression
}

func (e *FunctionCallExpression) Pos() token.Pos {
	identifier := e.Identifier.(*IdentifierExpression)
	return identifier.Pos()
}

type ArrayReferenceExpression struct {
	Target Expression
	Index  Expression
}

func (e *ArrayReferenceExpression) Pos() token.Pos {
	return e.Target.Pos()
}

type PointerExpression struct {
	pos   token.Pos
	Value Expression
}

func (e *PointerExpression) Pos() token.Pos { return e.pos }

type Declarator struct {
	Identifier Expression
	Size       int
}

func (e *Declarator) Pos() token.Pos {
	switch identifier := e.Identifier.(type) {
	case *IdentifierExpression:
		return identifier.Pos()

	case *UnaryExpression:
		return identifier.Pos()
	}

	return -1
}

type Declaration struct {
	pos         token.Pos
	VarType     string
	Declarators []*Declarator
}

func (e *Declaration) Pos() token.Pos { return e.pos }

type FunctionDefinition struct {
	pos        token.Pos
	TypeName   string
	Identifier Expression
	Parameters []Expression
	Statement  Statement
}

func (e *FunctionDefinition) Pos() token.Pos { return e.pos }

type Statement interface {
	Node
}

type CompoundStatement struct {
	pos          token.Pos
	Declarations []Statement
	Statements   []Statement
}

func (e *CompoundStatement) Pos() token.Pos { return e.pos }

type ExpressionStatement struct {
	Value Expression
}

func (e *ExpressionStatement) Pos() token.Pos {
	return e.Value.Pos()
}

type IfStatement struct {
	pos            token.Pos
	Condition      Expression
	TrueStatement  Statement
	FalseStatement Statement
}

func (e *IfStatement) Pos() token.Pos { return e.pos }

type WhileStatement struct {
	pos       token.Pos
	Condition Expression
	Statement Statement
}

func (e *WhileStatement) Pos() token.Pos { return e.pos }

type ForStatement struct {
	pos       token.Pos
	Init      Expression
	Condition Expression
	Loop      Expression
	Statement Statement
}

func (e *ForStatement) Pos() token.Pos { return e.pos }

type ReturnStatement struct {
	pos   token.Pos
	Value Expression
}

func (e *ReturnStatement) Pos() token.Pos { return e.pos }

type ParameterDeclaration struct {
	pos        token.Pos
	TypeName   string
	Identifier Expression
}

func (e *ParameterDeclaration) Pos() token.Pos { return e.pos }
