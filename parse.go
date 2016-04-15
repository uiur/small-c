package main

import (
	"fmt"
	"go/scanner"
	"go/token"
	"strings"
)

// Parse returns ast
func Parse(src string) ([]Statement, error) {
	fset := token.NewFileSet()
	file := fset.AddFile("", fset.Base(), len(src))

	l := new(Lexer)
	l.Init(file, []byte(src), nil, scanner.ScanComments)
	yyErrorVerbose = true

	fail := yyParse(l)
	if fail == 1 {
		lineNumber, columnNumber := posToLineInfo(src, int(l.pos))
		err := fmt.Errorf("%d:%d: %s", lineNumber, columnNumber, l.errMessage)

		return nil, err
	}

	return l.result, nil
}

func posToLineInfo(src string, pos int) (int, int) {
	if pos < 0 {
		panic("pos must be positive")
	}

	lineNumber := strings.Count(src[:pos], "\n") + 1

	lines := strings.Split(src[:pos], "\n")
	lastLine := lines[len(lines)-1]
	columnNumber := len(lastLine)

	return lineNumber, columnNumber
}

// Walk iterates over statement nodes and replace syntax sugar
func Walk(statement Statement) Statement {
	switch s := statement.(type) {
	case *FunctionDefinition:
		for i, p := range s.Parameters {
			s.Parameters[i] = WalkExpression(p)
		}

		s.Statement = Walk(s.Statement)

		return s

	case *CompoundStatement:
		for i, st := range s.Statements {
			s.Statements[i] = Walk(st)
		}

		for i, d := range s.Declarations {
			s.Declarations[i] = Walk(d)
		}

		return s

	case *ForStatement:
		// for (init; cond; loop) s
		// => init; while (cond) { s; loop; }
		body := Walk(s.Statement)
		return &CompoundStatement{
			Statements: []Statement{
				&ExpressionStatement{Value: s.Init},
				&WhileStatement{
					pos:       s.Pos(),
					Condition: s.Condition,
					Statement: &CompoundStatement{
						Statements: []Statement{
							body,
							&ExpressionStatement{Value: s.Loop},
						},
					},
				},
			},
		}

	case *IfStatement:
		s.Condition = WalkExpression(s.Condition)
		s.TrueStatement = Walk(s.TrueStatement)
		s.FalseStatement = Walk(s.FalseStatement)

		return s

	case *ReturnStatement:
		s.Value = WalkExpression(s.Value)
		return s

	case *ExpressionStatement:
		s.Value = WalkExpression(s.Value)
		return s
	}

	return statement
}

func WalkExpression(expression Expression) Expression {
	switch e := expression.(type) {
	case *ExpressionList:
		for i, value := range e.Values {
			e.Values[i] = WalkExpression(value)
		}

		return e

	case *BinOpExpression:
		e.Left = WalkExpression(e.Left)
		e.Right = WalkExpression(e.Right)

		return e

	case *UnaryExpression:
		e.Value = WalkExpression(e.Value)

		if e.Operator == "-" {
			return &BinOpExpression{
				Left:     &NumberExpression{pos: e.Pos(), Value: "0"},
				Operator: "-",
				Right:    e.Value,
			}
		} else if e.Operator == "&" {
			// &(*e) -> e
			switch value := e.Value.(type) {
			case *UnaryExpression:
				if value.Operator == "*" {
					return value.Value
				}
			}
		}

		return e

	case *ArrayReferenceExpression:
		// a[100]  =>  *(a + 100)
		e.Target = WalkExpression(e.Target)
		e.Index = WalkExpression(e.Index)

		return &UnaryExpression{
			Operator: "*",
			Value: &BinOpExpression{
				Left:     e.Target,
				Operator: "+",
				Right:    e.Index,
			},
		}
	}

	return expression
}
