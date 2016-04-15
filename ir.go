// TODO:
//   * function call
//   * print
//   * return
package main

import (
  "strconv"
  "fmt"
  "github.com/k0kubun/pp"
)

// intermediate representation
type IRProgram struct {
  Declarations []*IRVariableDeclaration
  Functions []*IRFunctionDefinition
}

type IRStatement interface {}
type IRExpression interface {}

type IRVariableDeclaration struct {
  Var *Symbol
}

type IRFunctionDefinition struct {
  Var *Symbol
  Parameters []*IRVariableDeclaration
  Body IRStatement
}

type IRAssignmentStatement struct {
  Var *Symbol
  Expression IRExpression
}

type IRWriteStatement struct {
  Dest *Symbol
  Src *Symbol
}

type IRReadStatement struct {
  Dest *Symbol
  Src *Symbol
}

type IRLabelStatement struct {
  Name string
}

type IRIfStatement struct {
  Var *Symbol
  TrueLabel string
  FalseLabel string
}

type IRGotoStatement struct {
  Label string
}

type IRCallStatement struct {
  Dest *Symbol
  Func *Symbol
  Vars []*Symbol
}

type IRReturnStatement struct {
  Var *Symbol
}

type IRPrintStatement struct {
  Var *Symbol
}

type IRCompoundStatement struct {
  Declarations []*IRVariableDeclaration
  Statements []IRStatement
}

// IRExpression

type IRVariableExpression struct {
  Var *Symbol
}

type IRNumberExpression struct {
  Value int
}

type IRBinaryExpression struct {
  Operator string
  Left IRExpression
  Right IRExpression
}

type IRAddressExpression struct {
  Var *Symbol
}

var counter = map[string]int {}
func label(name string) string {
  labelName := fmt.Sprintf("%s-%d", name, counter[name])
  counter[name]++

  return labelName
}

func tmpvar() *Symbol {
  return &Symbol{
    Name: label("#tmp"),
    Type: Int(),
  }
}

// CompileIR convert Statements to intermediate representation
func CompileIR(statements []Statement) *IRProgram {
  var decls []*IRVariableDeclaration
  var funcs []*IRFunctionDefinition

  var irStatements []IRStatement
  for _, statement := range statements {
    switch s := statement.(type) {
    case *Declaration:
      symbols := findSymbolsFromDeclaration(s)
      decls = append(decls, IRVariableDeclarations(symbols)...)
    default:
      irStatements = append(irStatements, compileIRStatement(s))
    }
  }

  for _, statement := range irStatements {
    switch s := statement.(type) {
    case *IRFunctionDefinition:
      funcs = append(funcs, s)
    case *IRVariableDeclaration:
      decls = append(decls, s)
    }
  }

  return &IRProgram{
    Declarations: decls,
    Functions: funcs,
  }
}

func compileIRStatement(statement Statement) IRStatement {
  if statement == nil {
    return nil
  }

  switch s := statement.(type) {
  case *FunctionDefinition:
    if s.Statement == nil {
      return nil
    }

    identifier := findIdentifierExpression(s.Identifier)

    var paramSymbols []*Symbol
    for _, p := range s.Parameters {
  		parameter, ok := p.(*ParameterDeclaration)
  		if ok {
        identifier := findIdentifierExpression(parameter.Identifier)
        paramSymbols = append(paramSymbols, identifier.Symbol)
  		}
    }

    return &IRFunctionDefinition{
      Var: identifier.Symbol,
      Parameters: IRVariableDeclarations(paramSymbols),
      Body: compileIRStatement(s.Statement),
    }

  case *CompoundStatement:
    var symbols []*Symbol
    for _, d := range s.Declarations {
      declaration, ok := d.(*Declaration)
      if ok {
        symbols = append(symbols, findSymbolsFromDeclaration(declaration)...)
      }
    }

    var statements []IRStatement
    for _, statement := range s.Statements {
      statements = append(statements, compileIRStatement(statement))
    }

    return &IRCompoundStatement{
      Declarations: IRVariableDeclarations(symbols),
      Statements: statements,
    }

  case *ExpressionStatement:
    switch e := s.Value.(type) {
    case *ExpressionList:
    case *BinOpExpression:
      if e.IsAssignment() {
        assignee := findIdentifierExpression(e.Left)
        symbol := assignee.Symbol

        switch r := e.Right.(type) {
        case *UnaryExpression:
          // a = *p;
          if r.Operator == "*" {
            tmp := tmpvar()
            right, decls, beforeRight := compileIRExpression(r.Value)

            statements := []IRStatement{
              &IRAssignmentStatement{
                Var: tmp,
                Expression: right,
              },
              &IRReadStatement{Dest: symbol, Src: tmp},
            }

            return &IRCompoundStatement{
              Declarations: append(IRVariableDeclarations([]*Symbol{tmp}), decls...),
              Statements: append(beforeRight, statements...),
            }
          }
        }

        right, decls, beforeRight := compileIRExpression(e.Right)

        switch left := e.Left.(type) {
        case *UnaryExpression:
          // tmp = exp
          // *left = tmp
          if left.Operator == "*" {
            tmp := tmpvar()

            statements := []IRStatement{
              &IRAssignmentStatement{
                Var: tmp,
                Expression: right,
              },
              &IRWriteStatement{Dest: symbol, Src: tmp},
            }

            return &IRCompoundStatement{
              Declarations: append(IRVariableDeclarations([]*Symbol{tmp}), decls...),
              Statements: append(beforeRight, statements...),
            }
          }

        default:
          body := &IRAssignmentStatement{
            Var: symbol,
            Expression: right,
          }

          return &IRCompoundStatement {
            Declarations: decls,
            Statements: append(beforeRight, body),
          }
        }
      }
    }

  case *IfStatement:
    conditionVar := tmpvar()

    trueLabel := label("true")
    falseLabel := label("false")
    endLabel := label("end")

    condition, decls, beforeCondition := compileIRExpression(s.Condition)

    statements := []IRStatement{
      &IRAssignmentStatement{
        Var: conditionVar,
        Expression: condition,
      },
      &IRIfStatement{
        Var: conditionVar,
        TrueLabel: trueLabel,
        FalseLabel: falseLabel,
      },
      &IRLabelStatement{ Name: trueLabel },
      compileIRStatement(s.TrueStatement),
      &IRGotoStatement{ Label: endLabel },
      &IRLabelStatement{ Name: falseLabel },
      compileIRStatement(s.FalseStatement),
      &IRLabelStatement{ Name: endLabel },
    }

    return &IRCompoundStatement{
      Declarations: append(IRVariableDeclarations([]*Symbol{conditionVar}), decls...),
      Statements: append(beforeCondition, statements...),
    }

  case *WhileStatement:
    conditionVar := &Symbol {
      Name: "#tmp",
      Type: Int(),
    }

    beginLabel := label("while-begin")
    endLabel := label("while-end")

    condition, decls, beforeCondition := compileIRExpression(s.Condition)

    statements := []IRStatement{
      &IRAssignmentStatement{
        Var: conditionVar,
        Expression: condition,
      },
      &IRLabelStatement{ Name: beginLabel },
      &IRIfStatement{
        Var: conditionVar,
        FalseLabel: endLabel,
      },
      compileIRStatement(s.Statement),
      &IRGotoStatement{ Label: beginLabel },
      &IRLabelStatement{ Name: endLabel },
    }

    return &IRCompoundStatement{
      Declarations: append(IRVariableDeclarations([]*Symbol{conditionVar}), decls...),
      Statements: append(beforeCondition, statements...),
    }

  case *ReturnStatement:
    panic("not implemented")

  default:
    pp.Println(s)
    panic("unexpected statement")
  }

  return nil
}

func IRVariableDeclarations(symbols []*Symbol) []*IRVariableDeclaration {
  var declarations []*IRVariableDeclaration
  for _, symbol := range symbols {
    declarations = append(declarations, &IRVariableDeclaration{
      Var: symbol,
    })
  }

  return declarations
}

func findSymbolsFromDeclaration(declaration *Declaration) []*Symbol {
  var symbols []*Symbol
  for _, declarator := range declaration.Declarators {
    identifier := findIdentifierExpression(declarator.Identifier)
    symbols = append(symbols, identifier.Symbol)
  }

  return symbols
}

func compileIRExpression(expression Expression) (IRExpression, []*IRVariableDeclaration, []IRStatement) {
  switch e := expression.(type) {
  case *NumberExpression:
    value, _ := strconv.Atoi(e.Value)
    return &IRNumberExpression{
      Value: value,
    }, nil, nil

  case *IdentifierExpression:
    return &IRVariableExpression{
      Var: e.Symbol,
    }, nil, nil

  case *UnaryExpression:
    if e.Operator == "&" {
      value, decls, statements := compileIRExpression(e.Value)
      v, _ := value.(*IRVariableExpression)

      return &IRAddressExpression{
        Var: v.Var,
      }, decls, statements
    }

  case *BinOpExpression:
    // return (a || b) && c
    // v;
    // if (a) {
    //   v = 1
    // } else if (b) {
    //   v = 1;
    // } else {
    //   v = 0;
    // }
    // int v;
    // if (a) {
    //   if (b) {
    //     v = 1;
    //   } else {
    //     v = 0;
    //   }
    // } else {
    //   v = 0;
    // }

    switch e.Operator {
    case "&&":
      tmp := tmpvar()

      decls := IRVariableDeclarations([]*Symbol{tmp})
      statements := []IRStatement {
        compileIRStatement(&IfStatement{
          Condition: e.Left,
          TrueStatement: &IfStatement{
            Condition: e.Right,
            TrueStatement: assignStatementBySymbol(tmp, 1),
            FalseStatement: assignStatementBySymbol(tmp, 0),
          },
          FalseStatement: assignStatementBySymbol(tmp, 0),
        }),
      }

      return &IRVariableExpression{Var: tmp}, decls, statements

    case "||":
      tmp := tmpvar()

      decls := IRVariableDeclarations([]*Symbol{tmp})
      statements := []IRStatement {
        compileIRStatement(&IfStatement{
          Condition: e.Left,
          TrueStatement: assignStatementBySymbol(tmp, 1),
          FalseStatement: &IfStatement{
            Condition: e.Right,
            TrueStatement: assignStatementBySymbol(tmp, 1),
            FalseStatement: assignStatementBySymbol(tmp, 0),
          },
        }),
      }

      return &IRVariableExpression{Var: tmp}, decls, statements
    }

    left, leftDecls, beforeLeft := compileIRExpression(e.Left)
    right, rightDecls, beforeRight := compileIRExpression(e.Right)

    return &IRBinaryExpression{
      Operator: e.Operator,
      Left: left,
      Right: right,
    }, append(leftDecls, rightDecls...), append(beforeLeft, beforeRight...)

  default:
    panic("unexpected expression")

  }

  panic("unexpected expression")
}

func assignStatementBySymbol(symbol *Symbol, value int) *ExpressionStatement {
  return &ExpressionStatement {
    Value: &BinOpExpression{
      Operator: "=",
      Left: &IdentifierExpression{ Symbol: symbol },
      Right: &NumberExpression{ Value: strconv.Itoa(value) },
    },
  }
}
