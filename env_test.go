package main

import "testing"

func TestType(t *testing.T) {
	data := [][]string{
		{PointerType{Value: BasicType{Name: "int"}}.String(), "int*"},
		{
			ArrayType{Value: BasicType{Name: "int"}, Size: 4}.String(),
			"int[4]",
		},
		{
			FunctionType{
				Return: BasicType{Name: "int"},
				Args:   []SymbolType{BasicType{Name: "int"}, BasicType{Name: "int"}},
			}.String(),
			"(int, int) -> int",
		},
	}

	for _, pair := range data {
		if pair[0] != pair[1] {
			t.Errorf("expect `%v`, got `%v`", pair[1], pair[0])
		}
	}
}

func TestCreateChild(t *testing.T) {
	env := &Env{}

	child := env.CreateChild()
	if !(len(env.Children) > 0 && env.Children[0] == child && child.Level == env.Level+1) {
		t.Errorf("the return value should be a child: parent: %v, child: %v", env, child)
	}
}

func TestAdd(t *testing.T) {
	env := &Env{}
	err := env.Add(&Symbol{
		Name: "foo",
		Kind: "var",
	})

	if err != nil {
		t.Errorf("expect err == nil, but %v", err)
	}

	env.Add(&Symbol{
		Name: "bar",
		Kind: "var",
	})

	err = env.Add(&Symbol{
		Name: "bar",
		Kind: "var",
	})

	if err == nil {
		t.Errorf("should return already defined error, but err == nil")
		return
	}

	env.Add(&Symbol{Name: "f", Kind: "proto"})
	err = env.Add(&Symbol{Name: "f", Kind: "proto"})

	if err != nil {
		t.Errorf("kind `proto` can be defined double, but got \"%v\"", err)
	}
}

func TestRegister(t *testing.T) {
	env := &Env{}
	identifier := &IdentifierExpression{Name: "foo"}

	err := env.Register(identifier, &Symbol{
		Kind: "var",
	})

	if !(err == nil && identifier.Symbol.Name == "foo") {
		t.Errorf("expect identifier.Symbol == `foo`, but: %v", identifier.Symbol)
	}
}
