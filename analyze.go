package main

import "fmt"

// Analyze ast and register variables to env
func Analyze(statements []Statement, env *Env) []error {
	for _, statement := range statements {
		analyzeStatement(statement, env)
	}

	return nil
}

func analyzeStatement(statement Statement, env *Env) []error {
	errs := []error{}

	switch s := statement.(type) {
	case FunctionDefinition:
		analyzeFunctionDefinition(s, env)

	case Declaration:
		analyzeDeclaration(s, env)

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

	return errs
}

func analyzeFunctionDefinition(s FunctionDefinition, env *Env) []error {
	errs := []error{}

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

	err := env.Add(&Symbol{
		Name: name,
		Kind: kind,
		Type: symbolType,
	})

	if err != nil {
		errs = append(errs, err)
	}

	if s.Statement != nil {
		paramEnv := env.CreateChild()

		for _, p := range s.Parameters {
			parameter, ok := p.(ParameterDeclaration)

			if ok {
				name := parseIdentifierName(parameter.Identifier)
				argType := composeType(parameter.Identifier, BasicType{Name: parameter.TypeName})

				err := paramEnv.Add(&Symbol{
					Name: name,
					Kind: "parm",
					Type: argType,
				})

				if err != nil {
					errs = append(errs, fmt.Errorf("parameter `%s` is already defined", name))
				}
			}
		}

		analyzeStatement(s.Statement, paramEnv)
	}

	return errs
}

func analyzeDeclaration(s Declaration, env *Env) []error {
	errs := []error{}
	for _, declarator := range s.Declarators {
		name := parseIdentifierName(declarator.Identifier)

		symbolType := composeType(declarator.Identifier, BasicType{Name: s.VarType})
		if declarator.Size > 0 {
			symbolType = ArrayType{Value: symbolType, Size: declarator.Size}
		}

		err := env.Add(&Symbol{
			Name: name,
			Kind: "var",
			Type: symbolType,
		})

		if err != nil {
			errs = append(errs, err)
		}
	}

	return errs
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
