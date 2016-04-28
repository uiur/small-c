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

  code = append(code, compileStatement(function.Body, function)...)

  code = append(
    code,
    function.Var.Name + "_exit:",
    "lw $fp, 0($sp)",
    "lw $ra, 4($sp)",
    fmt.Sprintf("addi $sp, $sp, %d", size),
    "jr $ra",
  )

  return code
}

func compileStatement(statement IRStatement, function *IRFunctionDefinition) []string {
  var code []string

  switch s := statement.(type) {
  case *IRCompoundStatement:
    for _, statement := range s.Statements {
      code = append(code, compileStatement(statement, function)...)
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
    code = append(code,
      lw("$v0", s.Var),
      fmt.Sprintf("j %s_exit", function.Var.Name),
    )

  case *IRWriteStatement:
    return []string {
      lw("$t0", s.Src),
      lw("$t1", s.Dest),
      "sw $t0, 0($t1)",
    }

  case *IRReadStatement:
    return []string {
      lw("$t0", s.Src),
      "lw $t1, 0($t0)",
      sw("$t1", s.Dest),
    }

  case *IRLabelStatement:
    return append(code, s.Name + ":")

  case *IRIfStatement:
    falseLabel := label("ir_if_false")
    endLabel := label("ir_if_end")

    code = append(code,
      lw("$t0", s.Var),
      fmt.Sprintf("beq $t0, $zero, %s", falseLabel),
    )

    if len(s.TrueLabel) > 0 {
      code = append(code,
        fmt.Sprintf("j %s", s.TrueLabel),
      )
    } else {
      code = append(code,
        fmt.Sprintf("j %s", endLabel),
      )
    }

    code = append(code,
      falseLabel + ":",
    )

    if len(s.FalseLabel) > 0 {
      code = append(code,
        fmt.Sprintf("j %s", s.FalseLabel),
      )
    }

    code = append(code,
      endLabel + ":",
    )

  case *IRGotoStatement:
    code = append(code, jmp(s.Label))

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

    leftRegister := "$t1"
    rightRegister := "$t2"

    code = append(code, assignExpression(leftRegister, e.Left)...)
    code = append(code, assignExpression(rightRegister, e.Right)...)

    operation := assignBinaryOperation(register, e.Operator, leftRegister, rightRegister)
    return append(code, operation...)

  case *IRVariableExpression:
    code = append(code, lw(register, e.Var))

  case *IRAddressExpression:
    return []string {
      fmt.Sprintf("addi %s, $fp, %d", register, e.Var.Offset),
    }
  }

  return code
}

func assignBinaryOperation(register string, operator string, left string, right string) []string {
  inst := operatorToInst[operator]
  if len(inst) > 0 {
    return []string {
      fmt.Sprintf("%s %s, %s, %s", inst, register, left, right),
    }
  }

  switch operator {
  case "==":
    falseLabel := label("beq_true")
    endLabel := label("beq_end")

    return []string{
      fmt.Sprintf("beq $t1, $t2, %s", falseLabel),
      li(register, 0),
      fmt.Sprintf("j %s", endLabel),
      falseLabel + ":",
      li(register, 1),
      endLabel + ":",
    }

  case ">":
    // a > b <=> (a <= b) < 1
    return append(assignBinaryOperation(register, "<=", left, right),
      fmt.Sprintf("slti %s, %s, 1", register, register),
    )

  case "<=":
    // a <= b <=> a - 1 < b
    return []string{
      "addi $t1, $t1, -1",
      fmt.Sprintf("slt %s, %s, %s", register, left, right),
    }

  case ">=":
    // a >= b <=> b <= a
    return assignBinaryOperation(register, "<=", right, left)
  }

  panic("unimplemented operator: " + operator)
}

var operatorToInst = map[string]string{
  "+": "add",
  "-": "sub",
  "*": "mul",
  "/": "div",
  "<": "slt",
}

func jmp(label string) string {
  return fmt.Sprintf("j %s", label)
}

func li(register string, value int) string {
  return fmt.Sprintf("li %s, %d", register, value)
}

func lw(register string, src *Symbol) string {
  return fmt.Sprintf("lw %s, %d($fp)", register, src.Offset)
}

func sw(register string, dest *Symbol) string {
  return fmt.Sprintf("sw %s, %d($fp)", register, dest.Offset)
}
