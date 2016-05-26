package main

import (
	"fmt"
)

// Parse returns ast
func Parse(src string) ([]Statement, error) {
	l := new(Lexer)
	l.Init(src)
	yyErrorVerbose = true

	fail := yyParse(l)
	if fail == 1 {
		err := fmt.Errorf("%d:%d: %s", l.pos.Line, l.pos.Column, l.errMessage)

		return nil, err
	}

	return l.result, nil
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

		var statements []Statement
		if s.Init != nil {
			statements = append(statements, &ExpressionStatement{Value: WalkExpression(s.Init)})
		}

		body := Walk(s.Statement)
		whileBody := []Statement{body}
		if s.Loop != nil {
			whileBody = append(whileBody, &ExpressionStatement{Value: WalkExpression(s.Loop)})
		}


		var condition Expression
		if s.Condition != nil {
			condition = WalkExpression(s.Condition)
		} else {
			condition = &NumberExpression{Value: "1"}
		}

		statements = append(statements,
			&WhileStatement{
				pos:       s.Pos(),
				Condition: condition,
				Statement: &CompoundStatement{
					Statements: whileBody,
				},
			},
		)

		return &CompoundStatement{
			Statements: statements,
		}

	case *WhileStatement:
		s.Condition = WalkExpression(s.Condition)
		s.Statement = Walk(s.Statement)

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

	case *FunctionCallExpression:
		e.Argument = WalkExpression(e.Argument)

		return e

	case *BinaryExpression:
		e.Left = WalkExpression(e.Left)
		e.Right = WalkExpression(e.Right)

		return e

	case *UnaryExpression:
		e.Value = WalkExpression(e.Value)

		if e.Operator == "-" {
			return &BinaryExpression{
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
			Value: &BinaryExpression{
				Left:     e.Target,
				Operator: "+",
				Right:    e.Index,
			},
		}
	}

	return expression
}
