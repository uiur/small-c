package main

import (
	"testing"
)

func TestTypeOfExpression(t *testing.T) {
	{
		env := &Env{}

		identifier := &IdentifierExpression{Name: "a"}
		env.Register(identifier, &Symbol{
			Kind: "var",
			Type: BasicType{Name: "int"},
		})

		expression := &BinaryExpression{
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
		expression := &BinaryExpression{
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
        int *i;
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

func TestCheckTypeOfFunctionCallExpression(t *testing.T) {
	{
		statements := ast(`
      int sum(int a, int b) {
        return a + b;
      }

      int main() {
        return sum(1, 2);
      }
    `)

		err := CheckType(statements)
		if err != nil {
			t.Error(err)
		}
	}

	{
		statements := ast(`
      int sum(int a, int b) {
        return a + b;
      }

      int main() {
        return sum(1);
      }
    `)

		err := CheckType(statements)
		if err == nil {
			t.Error("expect argument error, but nil")
		}
	}

	{
		statements := ast(`
      int sum(int a, int b) {
        return a + b;
      }

      int main() {
        int *a;
        return sum(1, a);
      }
    `)

		err := CheckType(statements)
		if err == nil {
			t.Error("expect argument type mismatch error, but nil")
		}
	}
}

func TestCheckTypeOfReturn(t *testing.T) {
	{
		statements := ast(`
      int main() {
        return 1;
      }
    `)

		err := CheckType(statements)
		if err != nil {
			t.Errorf("expect no error, got %v", err)
		}
	}

	{
		statements := ast(`
      int main() {
        int *a;
        return a;
      }
    `)

		err := CheckType(statements)
		if err == nil {
			t.Error("expect return type mismatch error, but nil")
		}

	}

	{
		statements := ast(`
      void dance() {
        return;
      }
    `)

		err := CheckType(statements)
		if err != nil {
			t.Errorf("expect no error, got %v", err)
		}
	}
}

func TestCheckTypeOfDeclaration(t *testing.T) {
	statements := ast(`
    int a, b, c;
    int ary[10];
    void v;

    int main() {
      a + b;
    }
  `)

	err := CheckType(statements)
	if err != nil {
		t.Errorf("expect no error, got %v", err)
	}
}

func TestCheckVoidType(t *testing.T) {
	{
		statements := ast(`
      int main() {
        void a;
        a + 0;
      }
    `)

		err := CheckType(statements)
		if err == nil {
			t.Errorf("expect void error, but nil")
		}
	}

	{
		statements := ast(`
      void *a;
      int main() {
        ;
      }
    `)

		err := CheckType(statements)
		if err == nil {
			t.Errorf("expect void error, but nil")
		}
	}

	{
		statements := ast(`
      int f(void a) {
      }
    `)

		err := CheckType(statements)
		if err == nil {
			t.Errorf("expect void error, but nil")
		}
	}
}

func TestTypeSize(t *testing.T) {
	if Int().ByteSize() != 4 {
		t.Errorf("expect size of int == 4, got %v", Int().ByteSize())
	}

	arrayType := ArrayType{Value: Int(), Size: 4}
	expected := 4 * 4
	if arrayType.ByteSize() != expected {
		t.Errorf("expect size of array[4] == %v, got %v", expected, arrayType.ByteSize())
	}
}
