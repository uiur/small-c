package main

import (
	"testing"
)

func ast(src string) []Statement {
	statements, _ := Parse(src)

	env := &Env{}
	Analyze(statements, env)

  return statements
}

func TestTypeOfExpression(t *testing.T) {
	{
		env := &Env{}

		identifier := &IdentifierExpression{Name: "a"}
		env.Register(identifier, &Symbol{
			Kind: "var",
			Type: BasicType{Name: "int"},
		})

		expression := &BinOpExpression{
			Operator: "+",
			Left:     &NumberExpression{Value: "42"},
			Right:    identifier,
		}

		symbolType, err := typeOfExpression(expression)
		if err != nil {
			t.Errorf("expect no error: %v", err)
		}

		if symbolType.String() != "int" {
			t.Errorf("expect int type, but got: %v", symbolType)
		}
	}

	{
		// int a[10];
		// a - 1;
		expression := &BinOpExpression{
			Operator: "-",
			Left: &IdentifierExpression{
				Name: "a",
				Symbol: &Symbol{
					Name: "a",
					Kind: "var",
					Type: ArrayType{
						Value: Int(),
						Size:  10,
					},
				},
			},
			Right: &NumberExpression{Value: "1"},
		}

		symbolType, err := typeOfExpression(expression)

		if err != nil {
			t.Errorf("expect no error, but got %v", err)
		}

		if symbolType.String() != "int*" {
			t.Errorf("expect int* type, got %v", symbolType)
		}
	}
}

func TestTypeOfPointerExpression(t *testing.T) {
	{
		// pointer reference: &a
		expression := &UnaryExpression{
			Operator: "&",
			Value: &IdentifierExpression{
				Name: "a",
				Symbol: &Symbol{
					Name: "a",
					Kind: "var",
					Type: Int(),
				},
			},
		}

		symbolType, err := typeOfExpression(expression)
		if err != nil {
			t.Errorf("expect no error, but got %v", err)
		}

		if symbolType.String() != "int*" {
			t.Errorf("expect int* type for `&a`, got %v", symbolType)
		}
	}

	{
		// pointer dereference: *p
		expression := &UnaryExpression{
			Operator: "*",
			Value: &IdentifierExpression{
				Name: "p",
				Symbol: &Symbol{
					Name: "p",
					Kind: "var",
					Type: Pointer(Int()),
				},
			},
		}

		symbolType, err := typeOfExpression(expression)
		if err != nil {
			t.Errorf("expect no error, but got %v", err)
		}

		if symbolType.String() != "int" {
			t.Errorf("expect int type, got %v", symbolType)
		}
	}
}

func TestCheckTypeOfIfStatement(t *testing.T) {
	{
    statements := ast(`
      int main() {
        int a;
        int b;

        if (a == b) {
          a + b;
        }
      }
    `)

		err := CheckType(statements)
		if err != nil {
			t.Error(err)
		}
	}

  {
    statements := ast(`
      int main() {
        int a;
        int *b;

        if (a == b) {
          ;
        }
      }
    `)

		err := CheckType(statements)
		if err == nil {
			t.Error("expect error, but nil")
		}
  }
}

func TestCheckTypeOfWhileStatement(t *testing.T) {
  {
    statements := ast(`
      int main() {
        int i;
        while (i > 0) {
          i = i - 1;
        }
      }
    `)

		err := CheckType(statements)
		if err != nil {
			t.Error(err)
		}
  }

  {
    statements := ast(`
      int main() {
        int **i;
        while (i > 0) {
          ;
        }
      }
    `)

		err := CheckType(statements)
		if err == nil {
			t.Error("expect type error in condition, got nil")
		}
  }
}
