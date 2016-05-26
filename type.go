package main

import (
	"fmt"
	"reflect"
	"strings"
)

type SymbolType interface {
	String() string

	ByteSize() int
}

type BasicType struct {
	Name string
}

func (t BasicType) ByteSize() int {
	if t.Name == "int" {
		return 4
	}

	return 0
}

func (t BasicType) String() string {
	return t.Name
}

type PointerType struct {
	Value SymbolType
}

func (t PointerType) ByteSize() int {
	return Int().ByteSize()
}

func (t PointerType) String() string {
	return t.Value.String() + "*"
}

type ArrayType struct {
	Value SymbolType
	Size  int
}

func (t ArrayType) ByteSize() int {
	return t.Value.ByteSize() * t.Size
}

func (t ArrayType) String() string {
	return fmt.Sprintf("%s[%d]", t.Value.String(), t.Size)
}

type FunctionType struct {
	Return SymbolType
	Args   []SymbolType
}

func (t FunctionType) ByteSize() int {
	return 0
}

func (t FunctionType) String() string {
	args := []string{}
	for _, a := range t.Args {
		args = append(args, a.String())
	}

	return "(" + strings.Join(args, ", ") + ")" + " -> " + t.Return.String()
}

func Int() SymbolType {
	return BasicType{Name: "int"}
}

func Void() SymbolType {
	return BasicType{Name: "void"}
}

func Pointer(symbolType SymbolType) SymbolType {
	return PointerType{Value: symbolType}
}

// CheckType checks that ast is well-typed
// statements must be analyzed (should have symbol information)
func CheckType(statements []Statement) error {
	for _, s := range statements {
		err := CheckTypeOfStatement(s)
		if err != nil {
			return err
		}
	}

	return nil
}

func CheckTypeOfStatement(statement Statement) error {
	if statement == nil {
		return nil
	}

	switch s := statement.(type) {
	case *Declaration:
		for _, declarator := range s.Declarators {
			identifier := findIdentifierExpression(declarator.Identifier)
			t := identifier.Symbol.Type
			isVoid := strings.Contains(t.String(), "void")

			if isVoid && t.String() != "void" {
				return SemanticError{
					Pos: s.Pos(),
					Err: fmt.Errorf("type error: void pointer is not allowed"),
				}
			}
		}

		return nil

	case *FunctionDefinition:
		identifier := findIdentifierExpression(s.Identifier)
		funcType, _ := identifier.Symbol.Type.(FunctionType)
		for _, argType := range funcType.Args {
			isVoid := strings.Contains(argType.String(), "void")
			if isVoid {
				return SemanticError{
					Pos: s.Pos(),
					Err: fmt.Errorf("type error: void pointer is not allowed"),
				}
			}
		}

		return CheckTypeOfStatement(s.Statement)

	case *ExpressionStatement:
		_, err := typeOfExpression(s.Value)
		return err

	case *CompoundStatement:
		return CheckType(s.Statements)

	case *IfStatement:
		err := checkTypeOfCondition(s.Condition)
		if err != nil {
			return err
		}

		return CheckType(s.Statements())

	case *WhileStatement:
		err := checkTypeOfCondition(s.Condition)
		if err != nil {
			return err
		}

		return CheckType(s.Statements())

	case *ReturnStatement:
		var valueType SymbolType

		if s.Value == nil {
			valueType = Void()
		} else {
			t, err := typeOfExpression(s.Value)
			valueType = t

			if err != nil {
				return err
			}
		}

		functionType := s.FunctionSymbol.Type.(FunctionType)
		if valueType.String() != functionType.Return.String() {
			return SemanticError{
				Pos: s.Pos(),
				Err: fmt.Errorf("type error: must return %v, not %v", functionType.Return, valueType),
			}
		}

		return nil
	}

	return SemanticError{
		Pos: statement.Pos(),
		Err: fmt.Errorf("type error: statement %v", statement),
	}
}

func typeOfExpression(expression Expression) (SymbolType, error) {
	switch e := expression.(type) {
	case *NumberExpression:
		return BasicType{Name: "int"}, nil

	case *IdentifierExpression:
		switch t := e.Symbol.Type.(type) {
		case ArrayType:
			return Pointer(t.Value), nil

		default:
			return e.Symbol.Type, nil

		}

	case *ExpressionList:
		var lastType SymbolType

		for _, value := range e.Values {
			symbolType, err := typeOfExpression(value)

			if err != nil {
				return nil, err
			}

			lastType = symbolType
		}

		return lastType, nil

	case *UnaryExpression:
		valueType, err := typeOfExpression(e.Value)

		if err != nil {
			return nil, err
		}

		switch e.Operator {
		case "&":
			if valueType.String() == "int" {
				return Pointer(valueType), nil
			}

		case "*":
			switch t := valueType.(type) {
			case PointerType:
				return t.Value, nil

			default:
				return nil, SemanticError{
					Pos: e.Value.Pos(),
					Err: fmt.Errorf("type error: expect pointer type: %v", e.Value),
				}
			}
		}

	case *BinaryExpression:
		return typeOfBinaryExpression(e)

	case *FunctionCallExpression:
		var args []Expression
		switch arg := e.Argument.(type) {
		case *ExpressionList:
			args = arg.Values
		default:
			if arg != nil {
				args = []Expression{arg}
			}
		}

		identifier := findIdentifierExpression(e.Identifier)
		funcType := identifier.Symbol.Type.(FunctionType)

		if len(args) != len(funcType.Args) {
			return nil, SemanticError{
				Pos: e.Pos(),
				Err: fmt.Errorf("function `%v` must be called with %v arguments, not %v", identifier.Name, len(funcType.Args), len(args)),
			}
		}

		for i, arg := range args {
			argType, err := typeOfExpression(arg)
			if err != nil {
				return nil, err
			}

			if argType.String() != funcType.Args[i].String() {
				return nil, SemanticError{
					Pos: arg.Pos(),
					Err: fmt.Errorf("type error: argument type mismatch: %v", argType.String()),
				}
			}
		}

		return funcType.Return, nil
	}

	return nil, SemanticError{
		Pos: expression.Pos(),
		Err: fmt.Errorf("type error: expression %s", reflect.TypeOf(expression)),
	}
}

func typeOfBinaryExpression(e *BinaryExpression) (SymbolType, error) {
	leftType, leftErr := typeOfExpression(e.Left)
	if leftErr != nil {
		return nil, leftErr
	}

	rightType, rightErr := typeOfExpression(e.Right)
	if rightErr != nil {
		return nil, rightErr
	}

	if e.IsArithmetic() {
		if leftType.String() == "int" && rightType.String() == "int" {
			return BasicType{Name: "int"}, nil
		}

		switch e.Operator {
		case "+":
			// int* + int, int + int* -> int*
			if (leftType.String() == "int*" && rightType.String() == "int") || (leftType.String() == "int" && rightType.String() == "int*") {
				return Pointer(Int()), nil
			}

			// int** + int, int + int** -> int**
			if (leftType.String() == "int**" && rightType.String() == "int") || (leftType.String() == "int" && rightType.String() == "int**") {
				return Pointer(Pointer(Int())), nil
			}

		case "-":
			if leftType.String() == "int*" && rightType.String() == "int" {
				return Pointer(Int()), nil
			}

			if leftType.String() == "int**" && rightType.String() == "int" {
				return Pointer(Pointer(Int())), nil
			}
		}
	}

	if e.IsAssignment() {
		if leftType.String() == rightType.String() {
			return leftType, nil
		}
	}

	if e.IsLogical() {
		if leftType.String() == "int" && rightType.String() == "int" {
			return Int(), nil
		}
	}

	if e.IsEqual() {
		if leftType.String() == rightType.String() {
			return Int(), nil
		}
	}

	return nil, SemanticError{
		Pos: e.Pos(),
		Err: fmt.Errorf("type error: %v %v %v", leftType.String(), e.Operator, rightType.String()),
	}
}

func checkTypeOfCondition(condition Expression) error {
	t, err := typeOfExpression(condition)
	if err != nil {
		return err
	}

	if t.String() != "int" {
		return SemanticError{
			Pos: condition.Pos(),
			Err: fmt.Errorf("type error: condition must be int, not `%v`", t),
		}
	}

	return nil
}
