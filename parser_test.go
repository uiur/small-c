package main

import "testing"

func TestParseDeclaration(t *testing.T) {
	Parse("int foo;")
	Parse("int a, b, c;")
	Parse("int a[100];")
	Parse("void bar;")
}
