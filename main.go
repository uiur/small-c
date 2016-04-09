package main

import (
	"fmt"
	"go/scanner"
	"go/token"
	"io/ioutil"
	"os"
	"strings"

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

	prelude, _ := Parse("void print(int i);\n")
	statements = append(prelude, statements...)

	env := &Env{}
	analyze(statements, env)

	pp.Println(statements, env)
}

func analyze(statements []Statement, env *Env) {
	for _, statement := range statements {
		analyzeStatement(statement, env)
	}
}

func analyzeStatement(statement Statement, env *Env) {
	switch s := statement.(type) {
	case FunctionDefinition:
		name := parseIdentifierName(s.Identifier)

		argTypes := []SymbolType{}

		for _, p := range s.Parameters {
			parameter, ok := p.(ParameterDeclaration)
			if ok {
				argType := BasicType{Name: parameter.TypeName}
				argTypes = append(argTypes, composeType(parameter.Identifier, argType))
			}
		}

		returnType := BasicType{Name: s.TypeName}
		symbolType := FunctionType{Return: returnType, Args: argTypes}

		kind := ""
		if s.Statement != nil {
			kind = "fun"
		} else {
			kind = "proto"
		}

		env.Add(&Symbol{
			Name: name,
			Kind: kind,
			Type: symbolType,
		})

		if s.Statement != nil {
			paramEnv := env.CreateChild()

			for _, p := range s.Parameters {
				parameter, ok := p.(ParameterDeclaration)

				if ok {
					name := parseIdentifierName(parameter.Identifier)
					argType := composeType(parameter.Identifier, BasicType{Name: parameter.TypeName})

					paramEnv.Add(&Symbol{
						Name: name,
						Kind: "param",
						Type: argType,
					})
				}
			}

			analyzeStatement(s.Statement, paramEnv)
		}

	case Declaration:
		for _, declarator := range s.Declarators {
			name := parseIdentifierName(declarator.Identifier)

			symbolType := composeType(declarator.Identifier, BasicType{Name: s.VarType})
			if declarator.Size > 0 {
				symbolType = ArrayType{Value: symbolType, Size: declarator.Size}
			}

			env.Add(&Symbol{
				Name: name,
				Kind: "var",
				Type: symbolType,
			})
		}

	case CompoundStatement:
		newEnv := env.CreateChild()
		for _, declaration := range s.Declarations {
			analyzeStatement(declaration, newEnv)
		}

		for _, statement := range s.Statements {
			analyzeStatement(statement, newEnv)
		}

	case IfStatement:
		analyzeStatement(s.TrueStatement, env)
		analyzeStatement(s.FalseStatement, env)

	case WhileStatement:
		analyzeStatement(s.Statement, env)

	}
}

func parseIdentifierName(expression Expression) string {
	switch e := expression.(type) {
	case IdentifierExpression:
		return e.Name
	case UnaryExpression:
		return parseIdentifierName(e.Value)
	}

	return ""
}

func composeType(identifier Expression, symbolType SymbolType) SymbolType {
	switch e := identifier.(type) {
	case UnaryExpression:
		if e.Operator == "*" {
			return PointerType{Value: composeType(e.Value, symbolType)}
		}
	case IdentifierExpression:
		return symbolType
	}

	return nil
}

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
	lineNumber := strings.Count(src[:pos], "\n") + 1
	lines := strings.Split(src, "\n")
	columnNumber := len(lines[lineNumber-1])

	return lineNumber, columnNumber
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
		} else if e.Operator == "&" {
			// &(*e) -> e
			switch value := e.Value.(type) {
			case UnaryExpression:
				if value.Operator == "*" {
					return value.Value
				}
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
