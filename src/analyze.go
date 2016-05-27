package main

import (
	"errors"
	"fmt"
)

// Analyze ast and register variables to env
func Analyze(statements []Statement, env *Env) []error {
	var errs []error
	for _, statement := range statements {
		errs = append(errs, analyzeStatement(statement, env)...)
	}

	return errs
}

func analyzeStatement(statement Statement, env *Env) []error {
	var errs []error

	switch s := statement.(type) {
	case *FunctionDefinition:
		errs = analyzeFunctionDefinition(s, env)

	case *Declaration:
		errs = analyzeDeclaration(s, env)

	case *CompoundStatement:
		errs = analyzeCompoundStatement(s, env)

	case *IfStatement:
		errs = append(errs, analyzeExpression(s.Condition, env)...)
		errs = append(errs, analyzeStatement(s.TrueStatement, env)...)
		errs = append(errs, analyzeStatement(s.FalseStatement, env)...)

	case *WhileStatement:
		// ForStatement is converted to WhileStatement
		errs = analyzeExpression(s.Condition, env)
		errs = append(errs, analyzeStatement(s.Statement, env)...)

	case *ExpressionStatement:
		errs = analyzeExpression(s.Value, env)

	case *ReturnStatement:
		// Set current function symbol to check type
		s.FunctionSymbol = env.Get("#func")
		errs = analyzeExpression(s.Value, env)

	}

	return errs
}

func analyzeFunctionDefinition(s *FunctionDefinition, env *Env) []error {
	errs := []error{}

	identifier := findIdentifierExpression(s.Identifier)
	argTypes := []SymbolType{}

	for _, p := range s.Parameters {
		parameter, ok := p.(*ParameterDeclaration)
		if ok {
			argType := BasicType{Name: parameter.TypeName}
			argTypes = append(argTypes, composeType(parameter.Identifier, argType))
		}
	}

	returnType := composeType(s.Identifier, BasicType{Name: s.TypeName})
	symbolType := FunctionType{Return: returnType, Args: argTypes}

	kind := ""
	if s.Statement != nil {
		kind = "fun"
	} else {
		kind = "proto"
	}

	err := env.Register(identifier, &Symbol{
		Kind: kind,
		Type: symbolType,
	})

	if err != nil {
		errs = append(errs, SemanticError{
			Pos: s.Pos(),
			Err: err,
		})
	}

	if s.Statement != nil {
		paramEnv := env.CreateChild()
		// Set special symbol to analyze function type
		paramEnv.Add(&Symbol{
			Name: "#func",
			Type: symbolType,
		})

		for _, p := range s.Parameters {
			parameter, ok := p.(*ParameterDeclaration)

			if ok {
				identifier := findIdentifierExpression(parameter.Identifier)
				argType := composeType(parameter.Identifier, BasicType{Name: parameter.TypeName})

				err := paramEnv.Register(identifier, &Symbol{
					Kind: "parm",
					Type: argType,
				})

				if err != nil {
					errs = append(errs, SemanticError{
						Pos: parameter.Pos(),
						Err: fmt.Errorf("parameter `%s` is already defined", identifier.Name),
					})
				}
			}
		}

		errs = append(errs, analyzeStatement(s.Statement, paramEnv)...)
	}

	return errs
}

func analyzeDeclaration(s *Declaration, env *Env) []error {
	errs := []error{}
	for _, declarator := range s.Declarators {
		symbolType := composeType(declarator.Identifier, BasicType{Name: s.VarType})
		if declarator.Size > 0 {
			symbolType = ArrayType{Value: symbolType, Size: declarator.Size}
		}

		identifier := findIdentifierExpression(declarator.Identifier)
		err := env.Register(identifier, &Symbol{
			Kind: "var",
			Type: symbolType,
		})

		if err != nil {
			errs = append(errs, SemanticError{
				Pos: declarator.Pos(),
				Err: err,
			})
		}
	}

	return errs
}

func analyzeCompoundStatement(s *CompoundStatement, env *Env) []error {
	var errs []error
	newEnv := env.CreateChild()
	for _, declaration := range s.Declarations {
		errs = append(errs, analyzeStatement(declaration, newEnv)...)
	}

	for _, statement := range s.Statements {
		errs = append(errs, analyzeStatement(statement, newEnv)...)
	}

	return errs
}

func analyzeExpression(expression Expression, env *Env) []error {
	var errs []error

	switch e := expression.(type) {
	case *IdentifierExpression:
		symbol := env.Get(e.Name)

		if symbol == nil {
			errs = append(errs, SemanticError{
				Pos: e.Pos(),
				Err: fmt.Errorf("reference error: `%v` is undefined", e.Name),
			})
		} else {
			if !symbol.IsVariable() {
				errs = append(errs, SemanticError{
					Pos: e.Pos(),
					Err: fmt.Errorf("`%v` is not variable", e.Name),
				})
			} else {
				e.Symbol = symbol
			}
		}

	case *ExpressionList:
		for _, value := range e.Values {
			errs = append(errs, analyzeExpression(value, env)...)
		}

	case *BinaryExpression:
		errs = append(errs, analyzeExpression(e.Left, env)...)
		errs = append(errs, analyzeExpression(e.Right, env)...)

		if e.Operator == "=" {
			leftIsAssignable := true

			switch left := e.Left.(type) {
			case *IdentifierExpression:
				expressionErrs := analyzeExpression(left, env)
				errs = append(errs, expressionErrs...)

				if len(errs) == 0 {
					_, isArrayType := left.Symbol.Type.(ArrayType)
					if !left.Symbol.IsVariable() || isArrayType {
						leftIsAssignable = false
					}
				}

			case *UnaryExpression:
				if left.Operator != "*" {
					leftIsAssignable = false
				}

			default:
				leftIsAssignable = false
			}

			if !leftIsAssignable {
				errs = append(errs, SemanticError{
					Pos: e.Left.Pos(),
					Err: errors.New("expression is not assignable"),
				})
			}
		}

	case *UnaryExpression:
		if e.Operator == "&" {
			switch v := e.Value.(type) {
			case *IdentifierExpression:
			default:
				errs = append(errs, SemanticError{
					Pos: v.Pos(),
					Err: errors.New("the operand of `&` must be on memory"),
				})
			}
		}

		return append(errs, analyzeExpression(e.Value, env)...)

	case *ArrayReferenceExpression:
		errs = append(errs, analyzeExpression(e.Target, env)...)
		errs = append(errs, analyzeExpression(e.Index, env)...)

	case *FunctionCallExpression:
		identifier := findIdentifierExpression(e.Identifier)
		symbol := env.Get(identifier.Name)
		if symbol == nil {
			return []error{
				SemanticError{
					Pos: identifier.Pos(),
					Err: fmt.Errorf("unknown function `%v` call", identifier.Name),
				},
			}
		}

		if !(symbol.Kind == "fun" || symbol.Kind == "proto") {
			return []error{
				SemanticError{
					Pos: identifier.Pos(),
					Err: fmt.Errorf("`%v` is not a function", identifier.Name),
				},
			}
		}

		identifier.Symbol = symbol
		return analyzeExpression(e.Argument, env)
	}

	return errs
}

func findIdentifierExpression(expression Expression) *IdentifierExpression {
	switch e := expression.(type) {
	case *IdentifierExpression:
		return e
	case *UnaryExpression:
		return findIdentifierExpression(e.Value)
	}

	panic("IdentifierExpression not found")
}

func composeType(identifier Expression, symbolType SymbolType) SymbolType {
	switch e := identifier.(type) {
	case *UnaryExpression:
		if e.Operator == "*" {
			return PointerType{Value: composeType(e.Value, symbolType)}
		}
	case *IdentifierExpression:
		return symbolType
	}

	return nil
}
