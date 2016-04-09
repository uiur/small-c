package main

type Env struct {
	Table  map[string]*Symbol
	Parent *Env
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
