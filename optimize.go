package main

import (
  "github.com/k0kubun/pp"
)


type DataflowBlock struct {
  Name string
  Func *IRFunctionDefinition
	Statements []IRStatement
	Next       []*DataflowBlock
}

func Optimize(program *IRProgram) *IRProgram {
  for _, f := range program.Functions {
    statements := flatStatement(f)

  	var blocks []*DataflowBlock

  	block := &DataflowBlock{}
  	for _, statement := range statements {
  		switch s := statement.(type) {
  		case *IRFunctionDefinition, *IRLabelStatement:
  			// in
        if len(block.Statements) > 0 {
    			blocks = append(blocks, block)
        }

  			block = &DataflowBlock{Statements: []IRStatement{s}}

  		case *IRIfStatement, *IRGotoStatement, *IRReturnStatement:
  			// out
  			block.Statements = append(block.Statements, s)
  			blocks = append(blocks, block)
        block = &DataflowBlock{}

  		default:
  			block.Statements = append(block.Statements, s)
  		}
  	}

    if len(block.Statements) > 0 {
      blocks = append(blocks, block)
    }

    // add info of function definition
    var definition *IRFunctionDefinition
    for _, block := range blocks {
      first := block.Statements[0]
      d, ok := first.(*IRFunctionDefinition)
      if ok {
        definition = d
      }

      block.Func = definition
    }

    beginBlock := &DataflowBlock{Name: "BEGIN"}
    beginBlock.Next = append(beginBlock.Next, blocks[0])

    endBlock := &DataflowBlock{Name: "END"}
    lastBlock := blocks[len(blocks)-1]
    lastBlock.Next = append(lastBlock.Next, endBlock)

    // add edge
    for _, block := range blocks {
      jumpStatement := block.Statements[len(block.Statements)-1]
      switch s := jumpStatement.(type) {
      case *IRGotoStatement:
        // goto label -> label block
        nextBlock := findBlockByLabel(blocks, s.Label)
        block.Next = append(block.Next, nextBlock)

      case *IRIfStatement:
        // if block -> true_label block, false_label block
        if len(s.TrueLabel) > 0 {
          trueLabelBlock := findBlockByLabel(blocks, s.TrueLabel)
          block.Next = append(block.Next, trueLabelBlock)
        }

        if len(s.FalseLabel) > 0 {
          falseLabelBlock := findBlockByLabel(blocks, s.FalseLabel)
          block.Next = append(block.Next, falseLabelBlock)
        }

      case *IRReturnStatement:
        // return block -> end block
        block.Next = append(block.Next, endBlock)

      }
    }

  	pp.Println(beginBlock)
  }

	return program
}

func findBlockByLabel(blocks []*DataflowBlock, label string) *DataflowBlock {
  for _, block := range blocks {
    inStatement := block.Statements[0]
    labelStatement, ok := inStatement.(*IRLabelStatement)
    if ok && labelStatement.Name == label {
      return block
    }
  }

  return nil
}

func findBlockByFuncName(blocks []*DataflowBlock, name string) *DataflowBlock {
  for _, block := range blocks {
    inStatement := block.Statements[0]
    funcStatement, ok := inStatement.(*IRFunctionDefinition)
    if ok && funcStatement.Var.Name == name {
      return block
    }
  }

  return nil
}

func findCallBlocks(blocks []*DataflowBlock, funcName string) []*DataflowBlock {
  var result []*DataflowBlock
  for i, block := range blocks {
    outStatement := block.Statements[len(block.Statements)-1]

    callStatement, ok := outStatement.(*IRCallStatement)
    if ok && callStatement.Func.Name == funcName {
      result = append(result, blocks[i+1])
    }
  }

  return result
}

func irStatements(program *IRProgram) []IRStatement {
	var statements []IRStatement
	for _, f := range program.Functions {
		statements = append(statements, flatStatement(f)...)
	}

	return statements
}

func flatStatement(statement IRStatement) []IRStatement {
	switch s := statement.(type) {
	case *IRFunctionDefinition:
		var statements []IRStatement

		statements = append(statements, s)

		for _, p := range s.Parameters {
			statements = append(statements, flatStatement(p)...)
		}

		statements = append(statements, flatStatement(s.Body)...)

		return statements

	case *IRCompoundStatement:
		var statements []IRStatement

		for _, d := range s.Declarations {
			statements = append(statements, flatStatement(d)...)
		}

		for _, child := range s.Statements {
			statements = append(statements, flatStatement(child)...)
		}

		return statements
	}

	return []IRStatement{statement}
}
