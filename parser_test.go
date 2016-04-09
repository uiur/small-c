package main

import (
	"strings"
	"testing"
)

func TestParse(t *testing.T) {
	statements, err := Parse(`
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

	if err != nil {
		t.Error(err)
		return
	}

	if len(statements) == 0 {
		t.Errorf("expect len(statements) > 0, actual: %v", len(statements))
	}
}

func TestParseError(t *testing.T) {
	_, err := Parse(`
    wtf this is wtf
  `)

	if !(err != nil && strings.Contains(err.Error(), "syntax error")) {
		t.Errorf("expect syntax error, but success")
	}
}

func TestParseDeclaration(t *testing.T) {
	_, err := Parse(`
    int foo, bar;
    void baz;
    int a[100];
    int *pointer;

    int sum(int a, int b);
    int *foo();
  `)

	if err != nil {
		t.Error(err)
	}
}

func TestParseCompoundStatement(t *testing.T) {
	_, err := Parse(`
    int foo() {
      int a;

      {
        a = a + b;
      };

      return b;
    }
  `)

	if err != nil {
		t.Error(err)
	}
}

func TestParseFunctionDefinition(t *testing.T) {
	_, err := Parse(`
    int sum(int a, int b) {
      return a + b;
    }
  `)

	if err != nil {
		t.Error(err)
	}
}

func TestParseIfStatement(t *testing.T) {
	_, err := Parse(`
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

	if err != nil {
		t.Error(err)
	}
}

func TestParseWhileStatement(t *testing.T) {
	_, err := Parse(`
    int main() {
      int a;

      a = 100;
      while (a) {
        a = a - 1;
      }
    }
  `)

	if err != nil {
		t.Error(err)
	}
}

func TestParseForStatement(t *testing.T) {
	statements, err := Parse(`
    int main() {
      for (i = 0; i < 100; i = i + 1) {
        sum = sum + i;
      }

      for (;;) {
        return;
      }
    }
  `)

	if err != nil {
		t.Error(err)
		return
	}

	switch mainStatements(statements)[0].(type) {
	case ForStatement:
	default:
		t.Error("expected ForStatement")
	}
}

func mainStatements(statements []Statement) []Statement {
	main := statements[0].(FunctionDefinition)

	return main.Statement.(CompoundStatement).Statements
}

func TestParseUnaryExpression(t *testing.T) {
	_, err := Parse(`
    int main() {
      a = &a;
      a = -a;
      a = *a;
    }
  `)

	if err != nil {
		t.Error(err)
	}
}
