package main

import (
  "fmt"
  "strings"
)

func CalculateOffset(ir *IRProgram) {
  for _, f := range ir.Functions {
    calculateOffsetFunction(f)
  }
}

func calculateOffsetFunction(ir *IRFunctionDefinition) {
  offset := 0

  for i := len(ir.Parameters)-1; i >= 0; i-- {
    p := ir.Parameters[i]
    p.Var.Offset = offset
    offset -= 4
  }

  minOffset := calculateOffsetStatement(ir.Body, offset)
  ir.VarSize = -minOffset
}

func calculateOffsetStatement(statement IRStatement, base int) int {
  offset := base
  minOffset := 0

  switch s := statement.(type) {
  case *IRCompoundStatement:
    for _, d := range s.Declarations {
      d.Var.Offset = offset
      offset -= 4
    }

    minOffset = offset
    for _, s := range s.Statements {
      statementOffset := calculateOffsetStatement(s, offset)

      if statementOffset < minOffset {
        minOffset = statementOffset
      }
    }
  }

  return minOffset
}

// Compile takes ir program as input and returns mips code
func Compile(program *IRProgram) string {
  CalculateOffset(program)

  code := ""
  code += ".text\n.globl main\n"
  for _, f := range program.Functions {
    code += "\n" + strings.Join(compileFunction(f), "\n") + "\n"
  }

  return code
}

func compileFunction(function *IRFunctionDefinition) []string {
  size := function.VarSize + 4 * 2 // arguments + local vars + $ra + $fp

  var code []string
  code = append(
    code,
    fmt.Sprintf("%s:", function.Var.Name),
    fmt.Sprintf("addi $sp, $sp, %d", -size),
    "sw $ra, 4($sp)",
    "sw $fp, 0($sp)",
    fmt.Sprintf("addi $fp, $sp, %d", size - 4),
  )

  for i := len(function.Parameters)-1; i >= 0; i-- {
    p := function.Parameters[i]
    code = append(code, fmt.Sprintf("sw $a%d, %d($fp)", i, p.Var.Offset))
  }

  code = append(code, compileStatement(function.Body)...)

  code = append(
    code,
    "lw $fp, 0($sp)",
    "lw $ra, 4($sp)",
    fmt.Sprintf("addi $sp, $sp, %d", size),
    "jr $ra",
  )

  return code
}

func compileStatement(statement IRStatement) []string {
  var code []string

  switch s := statement.(type) {
  case *IRCompoundStatement:
    for _, statement := range s.Statements {
      code = append(code, compileStatement(statement)...)
    }

  case *IRAssignmentStatement:
    code = append(code, assignExpression("$t0", s.Expression)...)
    code = append(code, fmt.Sprintf("sw $t0, %d($fp)", s.Var.Offset))

  case *IRCallStatement:
    for i, v := range s.Vars {
      code = append(code, fmt.Sprintf("lw $a%d, %d($fp)", i, v.Offset))
    }

    code = append(code, fmt.Sprintf("jal %s", s.Func.Name))
    code = append(code, fmt.Sprintf("sw $v0, %d($fp)", s.Dest.Offset))

  case *IRReturnStatement:
    code = append(code, lw("$v0", s.Var))

  case *IRWriteStatement:
  case *IRReadStatement:
  case *IRLabelStatement:
  case *IRIfStatement:
  case *IRGotoStatement:
  case *IRPrintStatement:
    return []string{
      "li $v0, 1",
      lw("$a0", s.Var),
      "syscall",
    }

  }

  return code
}

func assignExpression(register string, expression IRExpression) []string {
  var code []string

  switch e := expression.(type) {
  case *IRNumberExpression:
    code = append(code, fmt.Sprintf("li %s, %d", register, e.Value))

  case *IRBinaryExpression:
    switch e.Operator {
    case "+":
      code = append(code, assignExpression("$t1", e.Left)...)
      code = append(code, assignExpression("$t2", e.Right)...)
      code = append(code, fmt.Sprintf("add %s, $t1, $t2", register))

    default:
      panic("implement!")
    }

  case *IRVariableExpression:
    code = append(code, lw(register, e.Var))

  case *IRAddressExpression:
  }

  return code
}

func lw(register string, symbol *Symbol) string {
  return fmt.Sprintf("lw %s, %d($fp)", register, symbol.Offset)
}
