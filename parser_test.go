package main

import "testing"

func TestParse(t *testing.T) {
	Parse(`
    int data[8];

    int main() {
      int i;
      int j;
      int tmp;
      int size;

      size = 8;
      for (i = 0; i < size; i = i + 1) {
        for (j = 1; j < size; j = j + 1) {
          if (data[j] < data[j-1]) {
            tmp = data[j];
            data[j] = data[j-1];
            data[j-1] = tmp;
          }
        }
      }
    }
  `)
}

func TestParseDeclaration(t *testing.T) {
	Parse(`
    int foo, bar;
    void baz;
    int a[100];

    int sum(int a, int b);
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
      int a;

      a = 100;
      while (a) {
        a = a - 1;
      }
    }
  `)
}

func TestParseForStatement(t *testing.T) {
	Parse(`
    int main() {
      int i;
      int sum;

      sum = 0;
      for (i = 0; i < 100; i = i + 1) {
        sum = sum + i;
      }

      for (;;) {
        return;
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
