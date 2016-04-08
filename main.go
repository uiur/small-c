package main

import (
	"fmt"
	"go/scanner"
	"go/token"
	"io/ioutil"
	"os"

	"github.com/k0kubun/pp"
)

func main() {
	data, _ := ioutil.ReadAll(os.Stdin)
	statements, err := Parse(string(data))
	if err != nil {
		fmt.Println(err)
		return
	}

	for i, statement := range statements {
		statements[i] = Walk(statement)
	}

	pp.Print(statements)
}

func Parse(src string) ([]Statement, error) {
	fset := token.NewFileSet()
	file := fset.AddFile("", fset.Base(), len(src))

	l := new(Lexer)
	l.Init(file, []byte(src), nil, scanner.ScanComments)
	yyErrorVerbose = true

	fail := yyParse(l)
	if fail == 1 {
		return nil, l.err
	}

	return l.result, nil
}

// Iterate over statement nodes and replace syntax sugar
func Walk(statement Statement) Statement {
	switch s := statement.(type) {
	case FunctionDefinition:
		for i, p := range s.Parameters {
			s.Parameters[i] = WalkExpression(p)
		}

		s.Statement = Walk(s.Statement)

		return s

	case CompoundStatement:
		for i, st := range s.Statements {
			s.Statements[i] = Walk(st)
		}

		for i, d := range s.Declarations {
			s.Declarations[i] = Walk(d)
		}

		return s

	case ForStatement:
		// for (init; cond; loop) s
		// => init; while (cond) { s; loop; }
		return CompoundStatement{
			Statements: []Statement{
				ExpressionStatement{Value: s.Init},
				WhileStatement{
					pos:       s.Pos(),
					Condition: s.Condition,
					Statement: CompoundStatement{
						Statements: []Statement{
							s.Statement,
							ExpressionStatement{Value: s.Loop},
						},
					},
				},
			},
		}

	case IfStatement:
		s.Condition = WalkExpression(s.Condition)
		s.TrueStatement = Walk(s.TrueStatement)
		s.FalseStatement = Walk(s.FalseStatement)

		return s

	case ReturnStatement:
		s.Value = WalkExpression(s.Value)
		return s

	case ExpressionStatement:
		s.Value = WalkExpression(s.Value)
		return s
	}

	return statement
}

func WalkExpression(expression Expression) Expression {
	switch e := expression.(type) {
	case BinOpExpression:
		e.Left = WalkExpression(e.Left)
		e.Right = WalkExpression(e.Right)

		return e

	case UnaryExpression:
		e.Value = WalkExpression(e.Value)

		if e.Operator == "-" {
			return BinOpExpression{
				Left:     NumberExpression{pos: e.Pos(), Value: "0"},
				Operator: "-",
				Right:    e.Value,
			}
		}

		return e

	case ArrayReferenceExpression:
		// a[100]  =>  *(a + 100)
		e.Target = WalkExpression(e.Target)
		e.Index = WalkExpression(e.Index)

		return UnaryExpression{
			Operator: "*",
			Value: BinOpExpression{
				Left:     e.Target,
				Operator: "+",
				Right:    e.Index,
			},
		}
	}

	return expression
}
