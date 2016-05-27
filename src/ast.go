package main

import (
	"text/scanner"
)

type Token struct {
	lit string
	pos scanner.Position
}

type Node interface {
	Pos() scanner.Position
}

type Expression interface {
	Node
}

type ExpressionList struct {
	Values []Expression
}

func (e *ExpressionList) Pos() scanner.Position {
	first := e.Values[0]
	return first.Pos()
}

type NumberExpression struct {
	pos   scanner.Position
	Value string
}

func (e *NumberExpression) Pos() scanner.Position { return e.pos }

type IdentifierExpression struct {
	pos    scanner.Position
	Name   string
	Symbol *Symbol
}

func (e *IdentifierExpression) Pos() scanner.Position { return e.pos }

type UnaryExpression struct {
	pos      scanner.Position
	Operator string
	Value    Expression
}

func (e *UnaryExpression) Pos() scanner.Position { return e.pos }

type BinaryExpression struct {
	Left     Expression
	Operator string
	Right    Expression
}

func (e *BinaryExpression) Pos() scanner.Position {
	return e.Left.Pos()
}

func (e *BinaryExpression) IsAssignment() bool {
	return e.Operator == "="
}

func (e *BinaryExpression) IsArithmetic() bool {
	return e.Operator == "+" || e.Operator == "-" || e.Operator == "/" || e.Operator == "*"
}

func (e *BinaryExpression) IsLogical() bool {
	return e.Operator == "&&" || e.Operator == "||"
}

func (e *BinaryExpression) IsEqual() bool {
	switch e.Operator {
	case "==", "!=", ">=", ">", "<=", "<":
		return true
	}

	return false
}

type FunctionCallExpression struct {
	Identifier Expression
	Argument   Expression
}

func (e *FunctionCallExpression) Pos() scanner.Position {
	identifier := e.Identifier.(*IdentifierExpression)
	return identifier.Pos()
}

type ArrayReferenceExpression struct {
	Target Expression
	Index  Expression
}

func (e *ArrayReferenceExpression) Pos() scanner.Position {
	return e.Target.Pos()
}

type PointerExpression struct {
	pos   scanner.Position
	Value Expression
}

func (e *PointerExpression) Pos() scanner.Position { return e.pos }

type Declarator struct {
	Identifier Expression
	Size       int
}

func (e *Declarator) Pos() scanner.Position {
	switch identifier := e.Identifier.(type) {
	case *IdentifierExpression:
		return identifier.Pos()

	case *UnaryExpression:
		return identifier.Pos()
	}

	panic("unexpected identifier")
}

type Declaration struct {
	pos         scanner.Position
	VarType     string
	Declarators []*Declarator
}

func (e *Declaration) Pos() scanner.Position { return e.pos }

type FunctionDefinition struct {
	pos        scanner.Position
	TypeName   string
	Identifier Expression
	Parameters []Expression
	Statement  Statement
}

func (e *FunctionDefinition) Pos() scanner.Position { return e.pos }

type Statement interface {
	Node
}

type CompoundStatement struct {
	pos          scanner.Position
	Declarations []Statement
	Statements   []Statement
}

func (e *CompoundStatement) Pos() scanner.Position { return e.pos }

type ExpressionStatement struct {
	Value Expression
}

func (e *ExpressionStatement) Pos() scanner.Position {
	return e.Value.Pos()
}

type IfStatement struct {
	pos            scanner.Position
	Condition      Expression
	TrueStatement  Statement
	FalseStatement Statement
}

func (e *IfStatement) Pos() scanner.Position { return e.pos }
func (e *IfStatement) Statements() []Statement {
	return []Statement{e.TrueStatement, e.FalseStatement}
}

type WhileStatement struct {
	pos       scanner.Position
	Condition Expression
	Statement Statement
}

func (e *WhileStatement) Pos() scanner.Position { return e.pos }
func (e *WhileStatement) Statements() []Statement {
	return []Statement{e.Statement}
}

type ForStatement struct {
	pos       scanner.Position
	Init      Expression
	Condition Expression
	Loop      Expression
	Statement Statement
}

func (e *ForStatement) Pos() scanner.Position { return e.pos }

type ReturnStatement struct {
	pos            scanner.Position
	Value          Expression
	FunctionSymbol *Symbol
}

func (e *ReturnStatement) Pos() scanner.Position { return e.pos }

type ParameterDeclaration struct {
	pos        scanner.Position
	TypeName   string
	Identifier Expression
}

func (e *ParameterDeclaration) Pos() scanner.Position { return e.pos }
