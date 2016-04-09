package main

import "testing"

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
