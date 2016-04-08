package main

import "testing"

func TestParseDeclaration(t *testing.T) {
	Parse("int foo; void bar;")
	Parse("int a, b, c;")
	Parse("int a[100];")
}
