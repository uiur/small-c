package main

import "fmt"

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
		errs = analyzeStatement(s.TrueStatement, env)
		errs = append(errs, analyzeStatement(s.FalseStatement, env)...)

	case *WhileStatement:
		errs = analyzeStatement(s.Statement, env)

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

	returnType := BasicType{Name: s.TypeName}
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

func parseIdentifierName(expression Expression) string {
	switch e := expression.(type) {
	case *IdentifierExpression:
		return e.Name
	case *UnaryExpression:
		return parseIdentifierName(e.Value)
	}

	return ""
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
