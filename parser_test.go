package main

import "testing"

func TestParseDeclaration(t *testing.T) {
	Parse("int foo; void bar;")
	Parse("int a, b, c;")
	Parse("int a[100];")
}

func TestParseFunctionDefinition(t *testing.T) {
	Parse("int foo() {} \n")
	Parse(`
    int foo() {
      a = 1 + 2;
    }
  `)
}
