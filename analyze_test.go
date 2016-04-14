package main

import (
	"reflect"
	"testing"
)

func TestAnalyze(t *testing.T) {
	env := &Env{}
	statements, _ := Parse(`
		int sum(int a, int b) {
			return a + b;
		}
	`)

	errs := Analyze(statements, env)
	if len(errs) != 0 {
		t.Errorf("expect no error, but got: %v", errs)
	}
}

func TestAnalyzeDeclaration(t *testing.T) {
	{
		statements, err := Parse("int a, b, c;\n")

		if err != nil {
			t.Errorf("parse error: %v", err)
			return
		}

		declaration := statements[0].(*Declaration)

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
		declaration := statements[0].(*Declaration)

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
		declaration := statements[0].(*Declaration)

		errs := analyzeDeclaration(declaration, &Env{})

		if len(errs) == 0 {
			t.Errorf("should return an error when variables are double defined: %v", errs)
			return
		}
	}
}

func TestAnalyzeFunctionDefinition(t *testing.T) {
	{
		statements, _ := Parse(`
	    int foo(int a, int b) {
	      return a + b;
	    }
	  `)

		fn := statements[0].(*FunctionDefinition)
		env := &Env{}
		analyzeFunctionDefinition(fn, env)

		symbol := env.Table["foo"]
		if symbol == nil {
			t.Errorf("env should have `foo` as symbol: %v", env)
			return
		}

		symbolType, ok := symbol.Type.(FunctionType)
		if !ok {
			t.Errorf("symbol type should be FunctionType: %v", symbol)
			return
		}

		returnIsInt := symbolType.Return.(BasicType).Name == "int"
		if !returnIsInt {
			t.Errorf("expect return type to be int, but got %v", symbolType)
		}

		argsHaveTwoInt := len(symbolType.Args) == 2 && symbolType.Args[0].String() == "int"

		if !argsHaveTwoInt {
			t.Errorf("expect args to be (int, int): %v", symbolType.Args)
		}
	}

	{
		statements, _ := Parse(`
	    int foo(int a, int a) {
				int b;

	      return a + b;
	    }
	  `)

		fn := statements[0].(*FunctionDefinition)
		errs := analyzeFunctionDefinition(fn, &Env{})

		if len(errs) != 1 {
			t.Errorf("should return `parameter already defined` error: %v", errs)
		}
	}
}

func TestAnalyzeCompoundStatement(t *testing.T) {
	statements, _ := Parse(`
		int main() {
			int a;
			int a;
		}
	`)

	def := statements[0].(*FunctionDefinition)
	compoundStatement := def.Statement.(*CompoundStatement)
	errs := analyzeCompoundStatement(compoundStatement, &Env{})

	if len(errs) != 1 {
		t.Errorf("should have 1 error: %v", errs)
	}
}

func TestAnalyzeExpression(t *testing.T) {
	{
		env := &Env{}
		env.Add(&Symbol{Name: "foo", Kind: "var"})

		errs := analyzeExpression(&IdentifierExpression{Name: "foo"}, env)
		if len(errs) > 0 {
			t.Errorf("expect no error, got %v", errs)
		}

		errs = analyzeExpression(&IdentifierExpression{Name: "bar"}, env)
		if len(errs) != 1 {
			t.Errorf("expect reference error, got %v", errs)
		}
	}

	{
		env := &Env{}
		env.Add(&Symbol{Name: "foo", Kind: "fun"})

		errs := analyzeExpression(&IdentifierExpression{Name: "foo"}, env)
		if len(errs) != 1 {
			t.Errorf("expect not variable error, got %v", errs)
		}
	}

	{
		e := &FunctionCallExpression{
			Identifier: &IdentifierExpression{
				Name: "foo",
			},
			Argument: &ExpressionList{},
		}

		env := &Env{}
		env.Add(&Symbol{Name: "foo", Kind: "fun"})

		errs := analyzeExpression(e, env)

		if len(errs) != 0 {
			t.Errorf("expect no error, but got: %v", errs)
		}

		env.Table["foo"] = &Symbol{Name: "foo", Kind: "var"}
		errs = analyzeExpression(e, env)

		if len(errs) != 1 {
			t.Errorf("expect not function error, got %v", errs)
		}
	}

	{
		env := &Env{}
		errs := analyzeExpression(&UnaryExpression{
			Operator: "&",
			Value: &NumberExpression{
				Value: "10",
			},
		}, env)

		if len(errs) != 1 {
			t.Errorf("expect memory reference error, but got %v", errs)
		}
	}

	{
		env := &Env{}
		env.Add(&Symbol{
			Name: "f",
			Kind: "fun",
		})

		errs := analyzeExpression(&BinOpExpression{
			Operator: "=",
			Left: &NumberExpression{ Value: "1" },
			Right: &NumberExpression{ Value: "2" },
		}, env)

		if len(errs) != 1 {
			t.Errorf("expect assignment error, but got: %v", errs)
		}
	}
}
