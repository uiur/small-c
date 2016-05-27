package main

import (
	"testing"
)

func TestExtractVarsFromExpression(t *testing.T) {
	symbol := &Symbol{Name: "foo"}

	{
		// foo + 42
		vars := extractVarsFromExpression(
			&IRBinaryExpression{
				Operator: "+",
				Left:     &IRVariableExpression{Var: symbol},
				Right:    &IRNumberExpression{Value: 42},
			},
		)

		expected := len(vars) == 1 && vars[0] == symbol
		if !expected {
			t.Errorf("expect vars of `foo + 42` to be `foo`, got %v", vars)
		}
	}
}
