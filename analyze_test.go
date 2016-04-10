package main

import (
	"reflect"
	"testing"
)

func TestAnalyzeDeclaration(t *testing.T) {
	{
		statements, err := Parse("int a, b, c;\n")

		if err != nil {
			t.Errorf("parse error: %v", err)
			return
		}

		declaration := statements[0].(Declaration)

		env := &Env{}
		analyzeDeclaration(declaration, env)

		if len(env.Table) != 3 {
			t.Errorf("env.Table should have 3 vars, but %v", env.Table)
		}

		symbol := env.Table["a"]
		if !(symbol != nil && symbol.Name == "a" && symbol.Kind == "var") {
			t.Errorf("symbol should be a variable, got %v", symbol)
		}
	}

	{
		statements, _ := Parse("int a[10], b;\n")
		declaration := statements[0].(Declaration)

		env := &Env{}
		analyzeDeclaration(declaration, env)

		symbol := env.Table["a"]

		isArrayType := symbol != nil && reflect.TypeOf(symbol.Type).Name() == "ArrayType"
		correctSize := symbol != nil && symbol.Type.(ArrayType).Size == 10
		if !(isArrayType && correctSize) {
			t.Errorf("expect `a` to be an array: %v", symbol)
		}
	}

	{
		statements, _ := Parse("int a, b, a;\n")
		declaration := statements[0].(Declaration)

		errs := analyzeDeclaration(declaration, &Env{})

		if len(errs) == 0 {
			t.Errorf("should return an error when variables are double defined: %v", errs)
			return
		}
	}
}
