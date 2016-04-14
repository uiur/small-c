package main

import "testing"

func TestTypeOfExpression(t *testing.T) {
  {
    env := &Env{}

    identifier := &IdentifierExpression{ Name: "a" }
    env.Register(identifier, &Symbol{
      Kind: "var",
      Type: BasicType{ Name: "int" },
    })

    expression := &BinOpExpression{
      Operator: "+",
      Left: &NumberExpression{ Value: "42" },
      Right: identifier,
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
            Size: 10,
          },
        },
      },
      Right: &NumberExpression{ Value: "1" },
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
