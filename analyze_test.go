package main

import "testing"

func TestAnalyzeDeclaration(t *testing.T) {
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
