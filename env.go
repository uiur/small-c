package main

import (
	"fmt"
	"go/token"
	"strings"
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
	if symbol.Kind != "proto" && env.Table[name] != nil {
		return fmt.Errorf("`%s` is already defined", name)
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

type Symbol struct {
	Name  string
	Level int
	Kind  string
	Type  SymbolType
}

type SymbolType interface {
	String() string
}

type BasicType struct {
	Name string
}

func (t BasicType) String() string {
	return t.Name
}

type PointerType struct {
	Value SymbolType
}

func (t PointerType) String() string {
	return t.Value.String() + "*"
}

type ArrayType struct {
	Value SymbolType
	Size  int
}

func (t ArrayType) String() string {
	return fmt.Sprintf("%s[%d]", t.Value.String(), t.Size)
}

type FunctionType struct {
	Return SymbolType
	Args   []SymbolType
}

func (t FunctionType) String() string {
	args := []string{}
	for _, a := range t.Args {
		args = append(args, a.String())
	}

	return "(" + strings.Join(args, ", ") + ")" + " -> " + t.Return.String()
}

type SemanticError struct {
	error
	Pos token.Pos
	Err error
}

func (e SemanticError) Error() string {
	return e.Err.Error()
}
