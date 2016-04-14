package main

import (
	"fmt"
  "strings"
)

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
