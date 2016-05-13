package main


type DataflowBlock struct {
  Name string
  Func *IRFunctionDefinition
	Statements []IRStatement
	Next       []*DataflowBlock
  Prev       []*DataflowBlock
}

func (block *DataflowBlock) AddEdge(another *DataflowBlock) {
  block.Next = append(block.Next, another)
  another.Prev = append(another.Prev, block)
}

type BlockState map[*Symbol][]IRStatement

func (state BlockState) Equal(anotherState BlockState) bool {
  for symbol, statements := range state {
    if len(state[symbol]) != len(anotherState[symbol]) {
      return false
    }

    for i := range statements {
      if state[symbol][i] != anotherState[symbol][i] {
        return false
      }
    }
  }

  return true
}

func Optimize(program *IRProgram) *IRProgram {
  for _, f := range program.Functions {
    statements := flatStatement(f)
    blocks := splitStatemetsIntoBlocks(statements)

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

    buildDataflowGraph(blocks)
    searchReachingDefinitions(blocks)
    blockOut := searchReachingDefinitions(blocks)
    allStatementState := reachingDefinitionsOfStatements(blocks, blockOut, statements)
    foldConstant(statements, allStatementState)
  }

	return program
}

func searchReachingDefinitions(blocks []*DataflowBlock) map[*DataflowBlock]BlockState {
  blockOut := make(map[*DataflowBlock]BlockState)

  changed := true
  for changed {
    changed = false

    for _, block := range blocks {
      inState := analyzeBlock(blockOut, block)
      if !inState.Equal(blockOut[block]) {
        changed = true
      }


      blockOut[block] = inState
    }
  }

  return blockOut
}

func reachingDefinitionsOfStatements(blocks []*DataflowBlock, blockOut map[*DataflowBlock]BlockState, statements []IRStatement) map[IRStatement]BlockState {
  allStatementState := make(map[IRStatement]BlockState)
  for _, block := range blocks {
    inState := BlockState{}
    for _, prevBlock := range block.Prev {
      for key, value := range blockOut[prevBlock] {
        inState[key] = append(inState[key], value...)
      }
    }

    for _, statement := range statements {
      allStatementState[statement] = inState
      inState = analyzeReachingDefinition(statement, inState)
    }
  }

  return allStatementState
}

func foldConstant(statements []IRStatement, allStatementState map[IRStatement]BlockState) {
  for _, statement := range statements {
    foldConstantStatement(statement, allStatementState)
  }
}

func foldConstantStatement(statement IRStatement, allStatementState map[IRStatement]BlockState) (bool, int) {
  switch s := statement.(type) {
  case *IRAssignmentStatement:
    isConstant, value := foldConstantExpression(s, s.Expression, allStatementState)
    if isConstant {
      s.Expression = &IRNumberExpression{ Value: value }
      return true, value
    }
  }

  return false, 0
}

func foldConstantExpression(statement IRStatement, expression IRExpression, allStatementState map[IRStatement]BlockState) (bool, int) {
  switch e := expression.(type) {
  case *IRNumberExpression:
    return true, e.Value

  case *IRVariableExpression:
    symbol := e.Var
    definitionOfVar := allStatementState[statement][symbol]
    if len(definitionOfVar) == 1 && definitionOfVar[0] != statement {
      return foldConstantStatement(definitionOfVar[0], allStatementState)
    }

    return false, 0

  case *IRBinaryExpression:
    leftIsConstant, leftValue := foldConstantExpression(statement, e.Left, allStatementState)
    rightIsConstant, rightValue := foldConstantExpression(statement, e.Right, allStatementState)

    if leftIsConstant {
      e.Left = &IRNumberExpression{ Value: leftValue }
    }

    if rightIsConstant {
      e.Right = &IRNumberExpression{ Value: rightValue }
    }

    if leftIsConstant && rightIsConstant {
      switch e.Operator {
      case "+":
        return true, leftValue + rightValue

      case "-":
        return true, leftValue - rightValue

      case "*":
        return true, leftValue * rightValue

      case "/":
        return true, leftValue / rightValue

      case "<":
        value := 0
        if leftValue < rightValue {
          value = 1
        }
        return true, value
      case ">":
        value := 0
        if leftValue > rightValue {
          value = 1
        }
        return true, value

      case "<=":
        value := 0
        if leftValue <= rightValue {
          value = 1
        }
        return true, value

      case ">=":
        value := 0
        if leftValue >= rightValue {
          value = 1
        }
        return true, value

      case "==":
        value := 0
        if leftValue == rightValue {
          value = 1
        }
        return true, value
      }

      panic("unexpected operator: " + e.Operator)
    }


    return false, 0
  }

  return false, 0
}

func analyzeBlock(blockOut map[*DataflowBlock]BlockState, block *DataflowBlock) BlockState {
  inState := BlockState{}
  for _, prevBlock := range block.Prev {
    for key, statements := range blockOut[prevBlock] {
      for _, statement := range statements {
        found := false
        for _, v := range inState[key] {
          if v == statement {
            found = true
            break
          }
        }

        if !found {
          inState[key] = append(inState[key], statement)
        }
      }
    }
  }

  for _, statement := range block.Statements {
    inState = analyzeReachingDefinition(statement, inState)
  }

  return inState
}

func analyzeReachingDefinition(statement IRStatement, inState BlockState) BlockState {
  switch s := statement.(type) {
  case *IRAssignmentStatement:
    inState[s.Var] = []IRStatement{s}

  case *IRWriteStatement:
    inState[s.Dest] = []IRStatement{s}

  case *IRReadStatement:
    inState[s.Dest] = []IRStatement{s}

  case *IRCallStatement:
    inState[s.Dest] = []IRStatement{s}
  }

  return inState
}

func splitStatemetsIntoBlocks(statements []IRStatement) []*DataflowBlock {
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

  return blocks
}

func buildDataflowGraph(blocks []*DataflowBlock) *DataflowBlock {
  beginBlock := &DataflowBlock{Name: "BEGIN"}
  beginBlock.Next = append(beginBlock.Next, blocks[0])

  endBlock := &DataflowBlock{Name: "END"}
  lastBlock := blocks[len(blocks)-1]
  lastBlock.Next = append(lastBlock.Next, endBlock)

  for i, block := range blocks {
    lastStatement := block.Statements[len(block.Statements)-1]
    switch s := lastStatement.(type) {
    case *IRGotoStatement:
      // goto label -> label block
      nextBlock := findBlockByLabel(blocks, s.Label)
      block.AddEdge(nextBlock)

    case *IRIfStatement:
      // if block -> true_label block, false_label block
      if len(s.TrueLabel) > 0 {
        trueLabelBlock := findBlockByLabel(blocks, s.TrueLabel)
        block.AddEdge(trueLabelBlock)
      }

      if len(s.FalseLabel) > 0 {
        falseLabelBlock := findBlockByLabel(blocks, s.FalseLabel)
        block.AddEdge(falseLabelBlock)
      }

      if len(s.TrueLabel) == 0 || len(s.FalseLabel) == 0 {
        if i < len(blocks) - 1 {
          nextBlock := blocks[i+1]
          block.AddEdge(nextBlock)
        }
      }

    case *IRReturnStatement:
      // return block -> end block
      block.AddEdge(endBlock)

    default:
      if i < len(blocks) - 1 {
        nextBlock := blocks[i+1]
        block.AddEdge(nextBlock)
      }
    }
  }

  return beginBlock
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
