package main

import (
	"testing"

	"github.com/k0kubun/pp"
)

func TestCompileIR(t *testing.T) {
	{
		statements := ast(`
      int main() {
        int a;
        a = 1 + 2;
      }
    `)

		ir := CompileIR(statements)

		if len(ir.Functions) != 1 {
			t.Errorf("expect len(functions) == 1, got %v", len(ir.Functions))
		}
	}

	{
		statements := ast(`
      int main() {
        int a;

        if (a > 0) {
          a = 1;
        } else {
          a = 2;
        }
      }
    `)

		CompileIR(statements)
	}

	{
		statements := ast(`
      int main() {
        int a;
        while (a > 0) {
          a = a - 1;
        }
      }
    `)

		CompileIR(statements)
	}

	{
		statements := ast(`
      int main() {
        int a;
        int *p;
        p = &a;
      }
    `)

		CompileIR(statements)
	}

	{
		statements := ast(`
      int main() {
        int a[10];
        int *p;
        p = a;
        *(p + 1) = 1;
      }
    `)

		CompileIR(statements)
	}
}

func TestCompileIRStatement(t *testing.T) {
	// int a;
	// int *p;
	symbolP := &Symbol{Name: "p", Type: Pointer(Int())}
	symbolA := &Symbol{Name: "a", Type: Int()}

	{
		// *p = a;
		//
		// tmp = a
		// *p = tmp
		s := &ExpressionStatement{
			Value: &BinaryExpression{
				Operator: "=",
				Left: &UnaryExpression{
					Operator: "*",
					Value:    &IdentifierExpression{Symbol: symbolP},
				},
				Right: &IdentifierExpression{Symbol: symbolA},
			},
		}

		ir := compileIRStatement(s)
		compoundStatement, ok := ir.(*IRCompoundStatement)
		if !ok {
			t.Errorf("expect *IRCompoundStatement, but got %v", ir)
			return
		}

		_, ok = compoundStatement.Statements[2].(*IRWriteStatement)
		if !ok {
			t.Errorf("expect WriteStatement %v", ir)
		}
	}

	{
		// a = *p;
		s := &ExpressionStatement{
			Value: &BinaryExpression{
				Operator: "=",
				Left:     &IdentifierExpression{Symbol: symbolA},
				Right: &UnaryExpression{
					Operator: "*",
					Value:    &IdentifierExpression{Symbol: symbolP},
				},
			},
		}

		ir := compileIRStatement(s)
		compoundStatement, ok := ir.(*IRCompoundStatement)
		if !ok {
			t.Errorf("expect *IRCompoundStatement, but got %v", ir)
			return
		}

		if len(compoundStatement.Statements) == 0 {
			t.Error(compoundStatement)
		}
	}
}

func TestCompileIRExpression(t *testing.T) {
	// 0 || 1
	e := &BinaryExpression{
		Operator: "||",
		Left:     &NumberExpression{Value: "0"},
		Right:    &NumberExpression{Value: "1"},
	}

	ir, decls, before := compileIRExpression(e)
	if len(before) == 0 || len(decls) == 0 {
		t.Errorf("expect decls and statements, got %v %v", before, decls)
	}
	pp.Println(ir, decls, before)
}
