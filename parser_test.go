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

	Parse(`
    int sum(int a, int b) {
      return a + b;
    }
  `)
}

func TestParseIfStatement(t *testing.T) {
	Parse(`
    int foo(int a) {
      if (a == 0) a = 1;
      if (a != 0) a = 1;
    }
  `)

	Parse(`
    int foo(int a, int b) {
      if (a && b) return 1;
      if (a || b) return 1;
    }
  `)
}
