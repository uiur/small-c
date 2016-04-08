package main

import "testing"

func TestParseDeclaration(t *testing.T) {
	Parse(`
    int foo, bar;
    void baz;
    int a[100];
  `)
}

func TestParseCompoundStatement(t *testing.T) {
	Parse(`
    int foo() {
      int a;

      {
        a = a + b;
      };

      return b;
    }
  `)
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
      if (a > 0) a = 1;
      if (a >= 0) a = 1;
      if (a < 0) a = 1;
      if (a <= 0) a = 1;
      if (a && b) return 1;
      if (a || b) return 1;
    }
  `)

	Parse(`
    int foo(int a) {
      if (a) {
        return 1;
      } else {
        return 0;
      }
    }
  `)
}

func TestParseWhileStatement(t *testing.T) {
	Parse(`
    int main() {
      a = 100;
      while (a) {
        a = a - 1;
      }
    }
  `)
}

func TestParseUnaryExpression(t *testing.T) {
	Parse(`
    int main() {
      a = &a;
      a = -a;
      a = *a;
    }
  `)
}
