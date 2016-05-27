package main

import (
	"fmt"
	"text/scanner"
)

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
	found := env.Table[name]
	if found != nil {
		if symbol.IsVariable() {
			if (found.Kind == "proto" || found.Kind == "fun") && symbol.IsGlobal() {
				return fmt.Errorf("function `%v` is already defined", name)
			}
		}

		if found.Kind != "proto" {
			return fmt.Errorf("`%s` is already defined", name)
		}

		if found.Kind == "proto" && (symbol.Kind == "fun" || symbol.Kind == "proto") {
			functionType, _ := found.Type.(FunctionType)
			if symbol.Type.String() != functionType.String() {
				return fmt.Errorf("prototype mismatch error: function `%v`: `%v` != `%v`", name, functionType, symbol.Type)
			}
		}
	}

	if symbol.Level == 0 {
		symbol.Level = env.Level
	}

	env.Table[name] = symbol
	return nil
}

func (env *Env) Register(identifier *IdentifierExpression, symbol *Symbol) error {
	symbol.Name = identifier.Name
	err := env.Add(symbol)

	if err == nil {
		identifier.Symbol = symbol
	}

	return err
}

func (env *Env) Get(name string) *Symbol {
	symbol := env.Table[name]

	if symbol != nil {
		return symbol
	}

	if env.Parent != nil {
		return env.Parent.Get(name)
	}

	return nil
}

type Symbol struct {
	Name   string
	Level  int
	Kind   string
	Type   SymbolType
	Offset int
}

func (symbol *Symbol) IsVariable() bool {
	return symbol.Kind == "var" || symbol.Kind == "parm"
}

func (symbol *Symbol) IsGlobal() bool {
	return symbol.Level == 0
}

func (symbol *Symbol) AddressPointer() string {
	if symbol.IsGlobal() {
		return "$gp"
	} else {
		return "$fp"
	}
}

type SemanticError struct {
	error
	Pos scanner.Position
	Err error
}

func (e SemanticError) Error() string {
	return e.Err.Error()
}
