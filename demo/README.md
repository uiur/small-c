# 最終報告
## データフロー解析
* 到達可能定義解析

データフロー解析処理、最適化処理は `optimize.go` で行っている。

```go
// optimize.go
type DataflowBlock struct {
	Name       string  // BEGIN, END
	Statements []IRStatement
	Next       []*DataflowBlock
	Prev       []*DataflowBlock
}

func Optimize(program *IRProgram) *IRProgram {
	for i, f := range program.Functions {
		statements := flatStatement(f)

    // 中間表現プログラム列をデータフローのブロックごとに分ける
		blocks := splitStatementsIntoBlocks(statements)

    // ブロックの配列からデータフローを構成
    // blockそれぞれについて, block.Nextを設定していく
		buildDataflowGraph(blocks)

    // データフローを見て不動点反復により到達可能定義解析する
    // 返り値はブロックごとに, 各シンボルの到達可能な定義文 を入れたmap
    // blockOut = (DataflowBlock -> (*Symbol -> []IRStatement))
		blockOut := searchReachingDefinitions(blocks)

    // ...
	}

	return program
}

// 不動点反復なので、状態が収束するまで地道に解析して状態を更新していくという雰囲気
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

// ひとつのプログラム点を見て状態を更新する
// 到達可能定義解析の実質的な処理
func analyzeReachingDefinition(statement IRStatement, inState BlockState) BlockState {
	switch s := statement.(type) {
	case *IRAssignmentStatement:
		inState[s.Var] = []IRStatement{s}
		symbols := extractAddressVarsFromExpression(s.Expression)
		for _, symbol := range symbols {
			inState[symbol] = append(inState[symbol], s)
		}

	case *IRReadStatement:
		inState[s.Dest] = []IRStatement{s}

  // ポインタ参照書き込みがあったら, 諦めムードにしておく
	case *IRWriteStatement:
		for symbol := range inState {
			inState[symbol] = append(inState[symbol], s)
		}

	case *IRCallStatement:
		inState[s.Dest] = []IRStatement{s}

	}

	return inState
}

```

## 最適化
* 定数畳み込み
* 無駄な命令の除去

を到達可能定義解析を用いて実装した。

```go
func Optimize(program *IRProgram) *IRProgram {
	for i, f := range program.Functions {
    // ...
		blockOut := searchReachingDefinitions(blocks)

    // 実装の都合で文ごとの到達可能定義を計算しなおしている
		allStatementState := reachingDefinitionsOfStatements(blocks, blockOut, statements)

    // 定数畳み込み
		program.Functions[i] = transformByConstantFolding(program.Functions[i], allStatementState)
    // 無駄コード除去
		program.Functions[i] = transformByDeadCodeElimination(program.Functions[i], allStatementState)
	}

	return program
}
```

### 定数畳み込み

```go
func transformByConstantFolding(f *IRFunctionDefinition, allStatementState map[IRStatement]BlockState) *IRFunctionDefinition {
	traversed := Traverse(f, func(statement IRStatement) IRStatement {
		foldConstantStatement(statement, allStatementState)
		return statement
	})

	return traversed.(*IRFunctionDefinition)
}

// 代入文ならexpressionを見て、それが定数だったら埋め込む
func foldConstantStatement(statement IRStatement, allStatementState map[IRStatement]BlockState) (bool, int) {
	switch s := statement.(type) {
	case *IRAssignmentStatement:
		isConstant, value := foldConstantExpression(s, s.Expression, allStatementState)
		if isConstant {
			s.Expression = &IRNumberExpression{Value: value}
			return true, value
		}
	}

	return false, 0
}

// 到達可能定義の情報を使って、再帰的に定数畳み込みしていく
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
			e.Left = &IRNumberExpression{Value: leftValue}
		}

		if rightIsConstant {
			e.Right = &IRNumberExpression{Value: rightValue}
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

			case "!=":
				value := 0
				if leftValue != rightValue {
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
```

### 無駄な命令の除去

```go
// 収束するまで繰り返す
// * 文を使用しているか到達可能定義を用いて見ていく
// * 消しても大丈夫そうな使われていない文を発見したら消す
//
// 最後に要らない宣言を削除
func transformByDeadCodeElimination(f *IRFunctionDefinition, allStatementState map[IRStatement]BlockState) *IRFunctionDefinition {
	changed := true
	for changed {
		changed = false

		used := make(map[IRStatement]bool)
		markAsUsed := func(s IRStatement, symbol *Symbol) {
			for _, definition := range allStatementState[s][symbol] {
				used[definition] = true
			}
		}

		Traverse(f, func(statement IRStatement) IRStatement {
			switch s := statement.(type) {
			case *IRCompoundStatement:
				used[s] = true

			case *IRAssignmentStatement:
				if s.Var.IsGlobal() {
					used[s] = true
				}

				vars := extractVarsFromExpression(s.Expression)
				for _, v := range vars {
					markAsUsed(s, v)
				}

			case *IRReadStatement:
				if s.Dest.IsGlobal() {
					used[s] = true
				}

				markAsUsed(s, s.Src)

			case *IRWriteStatement:
				markAsUsed(s, s.Src)
				markAsUsed(s, s.Dest)

			case *IRCallStatement:
				if s.Dest.IsGlobal() {
					used[s] = true
				}

				for _, argVar := range s.Vars {
					markAsUsed(s, argVar)
				}

			case *IRSystemCallStatement:
				markAsUsed(s, s.Var)

			case *IRReturnStatement:
				markAsUsed(s, s.Var)

			case *IRIfStatement:
				markAsUsed(s, s.Var)
			}

			return statement
		})

		transformed := Traverse(f, func(statement IRStatement) IRStatement {
			switch statement.(type) {
			case *IRAssignmentStatement, *IRReadStatement:
				if !used[statement] {
					changed = true
					return nil
				}
			}

			return statement
		})

		f = transformed.(*IRFunctionDefinition)
	}

	return removeUnusedVariableDeclaration(f)
}

```

## 例
簡単な例を用いて最適化を試す。

``` c
// demo/optimize_constant.sc
int main() {
  int a, b;
  int c;
  c = 3;

  a = c; // 3
  b = a + c; // 3 + 3
  print(a + b == 9);  // 3 + 6 == 9
}
```

`-optimize=false` オプションをつけて、比較用に最適化しなかった結果を出力する。
``` c
❯ ./small-c -optimize=false demo/optimize_constant.sc > demo/optimize_constant.s
❯ ./small-c demo/optimize_constant.sc > demo/optimize_constant_optimized.s
```

### 最適化前

``` sh
❯ spim -show_stats  -f demo/optimize_constant.s
Loaded: /usr/local/share/spim/exceptions.s
1
--- Summary ---
# of executed instructions
- Total:    47
- Memory:   21
- Others:   26

--- Details ---
       add	       2
      addi	       9
     addiu	       2
      addu	       1
       beq	       1
       jal	       1
        jr	       1
        lw	      12
       ori	       5
       sll	       2
        sw	       9
   syscall	       2

```

### 最適化後

```sh
❯ spim -show_stats  -f demo/optimize_constant_optimized.s
Loaded: /usr/local/share/spim/exceptions.s
1
--- Summary ---
# of executed instructions
- Total:    22
- Memory:    7
- Others:   15

--- Details ---
      addi	       3
     addiu	       2
      addu	       1
       jal	       1
        jr	       1
        lw	       4
       ori	       3
       sll	       2
        sw	       3
   syscall	       2

```

Total: 47 -> 22

このように定数畳み込みと無駄な命令を除去を組み合わせると大きな効果がある場合がある。
