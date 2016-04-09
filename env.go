package main

import "fmt"

type Env struct {
	Table    map[string]*Symbol
	Level    int
	Children []*Env
	Parent   *Env
}

func (env *Env) CreateChild() *Env {
	newEnv := &Env{Parent: env, Level: env.Level + 1}
	env.Children = append(env.Children, newEnv)
	return newEnv
}

func (env *Env) Add(symbol *Symbol) error {
	if env.Table == nil {
		env.Table = map[string]*Symbol{}
	}

	name := symbol.Name
	if symbol.Kind != "proto" && env.Table[name] != nil {
		return fmt.Errorf("`%s` is already defined", name)
	}

	if symbol.Level == 0 {
		symbol.Level = env.Level
	}

	env.Table[name] = symbol
	return nil
}

type Symbol struct {
	Name  string
	Level int
	Kind  string
	Type  SymbolType
}

type SymbolType interface{}
type BasicType struct {
	Name string
}
type ArrayType struct {
	Value SymbolType
	Size  int
}
type PointerType struct {
	Value SymbolType
}
type FunctionType struct {
	Return SymbolType
	Args   []SymbolType
}
