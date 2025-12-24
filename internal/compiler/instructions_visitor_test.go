package compiler

import (
	"testing"

	"youpiteron.dev/white-monster-on-friday-night/internal/ast"
	"youpiteron.dev/white-monster-on-friday-night/internal/common"
	"youpiteron.dev/white-monster-on-friday-night/internal/lexer"
)

func makeSourcePos(offset, line, col, length int) *common.SourcePos {
	return &common.SourcePos{
		BasePos: common.BasePos{
			Offset: offset,
			Line:   line,
			Column: col,
		},
		Length: length,
	}
}

func makeIntLiteral(value int, isStatement bool, offset, line, col int) *ast.IntLiteral {
	return &ast.IntLiteral{
		Value:       value,
		PosAt:       makeSourcePos(offset, line, col, len(string(rune(value)))),
		IsStatement: isStatement,
	}
}

func makeIdentifier(name string, isStatement bool, offset, line, col int) *ast.Identifier {
	return &ast.Identifier{
		Name:        name,
		PosAt:       makeSourcePos(offset, line, col, len(name)),
		IsStatement: isStatement,
	}
}

func makeBinaryExpr(left ast.Expression, operator lexer.OperatorSubkind, right ast.Expression, isStatement bool, offset, line, col int) *ast.BinaryExpr {
	return &ast.BinaryExpr{
		Left:        left,
		Operator:    operator,
		Right:       right,
		PosAt:       makeSourcePos(offset, line, col, 1),
		IsStatement: isStatement,
	}
}

func makeCallExpr(identifier *ast.Identifier, arguments []ast.Expression, offset, line, col int) *ast.CallExpr {
	return &ast.CallExpr{
		Identifier: *identifier,
		Arguments:  arguments,
		PosAt:      makeSourcePos(offset, line, col, 1),
	}
}

// ---------- Statement Expression Optimization Tests ----------

func TestVisitBinaryExpr_StatementExpression_Optimized(t *testing.T) {
	// Test: 1 + 1 (as statement - should be optimized away)
	left := makeIntLiteral(1, true, 0, 1, 1)
	right := makeIntLiteral(1, true, 4, 1, 5)
	binaryExpr := makeBinaryExpr(left, lexer.OperatorPlus, right, true, 2, 1, 3)

	visitor := NewInstructionsVisitor()
	visitor.EnterModuleContext()

	binaryExpr.Visit(visitor)

	moduleContext := CastModuleContext(visitor.context)
	instructions := moduleContext.instructions

	if len(instructions) != 0 {
		t.Errorf("expected 0 instructions for optimized statement expression, got %d", len(instructions))
		for i, instr := range instructions {
			t.Errorf("  instruction %d: %s", i, instr.String())
		}
	}
}

/*
*
Test `println(1)` as a statement expression.

- Should NOT be optimized
- Should generate instructions to load println, load constant 1, and call
*/
func TestVisitCallExpr_StatementExpression_NotOptimized(t *testing.T) {
	identifier := makeIdentifier("println", false, 0, 1, 1)
	arg := makeIntLiteral(1, false, 8, 1, 9)
	callExpr := makeCallExpr(identifier, []ast.Expression{arg}, 7, 1, 8)

	visitor := NewInstructionsVisitor()
	visitor.EnterModuleContext()

	callExpr.Visit(visitor)

	moduleContext := CastModuleContext(visitor.context)
	instructions := moduleContext.instructions

	if len(instructions) == 0 {
		t.Error("expected instructions for call expression with side effects, got 0")
	}

	if len(instructions) < 3 {
		t.Errorf("expected at least 3 instructions for call expression, got %d", len(instructions))
	}

	hasCall := false
	for _, instr := range instructions {
		if instr.OpCode == CALL {
			hasCall = true
			break
		}
	}
	if !hasCall {
		t.Error("expected CALL instruction for function call")
	}
}

/*
*
Test `1 + println(1)` as a statement expression.

- `1 + ` part should be optimized away
- `println(1)` part should stay and generate instructions to load println, load constant 1, and call
*/
func TestVisitBinaryExpr_StatementExpression_WithSideEffects(t *testing.T) {
	left := makeIntLiteral(1, true, 0, 1, 1)
	callIdentifier := makeIdentifier("println", false, 4, 1, 5)
	callArg := makeIntLiteral(1, false, 12, 1, 13)
	callExpr := makeCallExpr(callIdentifier, []ast.Expression{callArg}, 11, 1, 12)
	binaryExpr := makeBinaryExpr(left, lexer.OperatorPlus, callExpr, true, 2, 1, 3)

	visitor := NewInstructionsVisitor()
	visitor.EnterModuleContext()

	binaryExpr.Visit(visitor)

	moduleContext := CastModuleContext(visitor.context)
	instructions := moduleContext.instructions

	hasAddInt := false
	hasCall := false
	for _, instr := range instructions {
		if instr.OpCode == ADD_INT {
			hasAddInt = true
		}
		if instr.OpCode == CALL {
			hasCall = true
		}
	}

	if hasAddInt {
		t.Error("expected ADD_INT instruction to be optimized away for statement expression")
	}

	if !hasCall {
		t.Error("expected CALL instruction for function call with side effects")
	}
}

/*
*
Test `42` as a statement expression.

- Should be optimized away
*/
func TestVisitIntLiteral_StatementExpression_Optimized(t *testing.T) {
	intLiteral := makeIntLiteral(42, true, 0, 1, 1)

	visitor := NewInstructionsVisitor()
	visitor.EnterModuleContext()

	result := intLiteral.Visit(visitor)

	if result != nil {
		t.Error("expected nil result for optimized statement expression")
	}

	moduleContext := CastModuleContext(visitor.context)
	instructions := moduleContext.instructions

	if len(instructions) != 0 {
		t.Errorf("expected 0 instructions for optimized statement expression, got %d", len(instructions))
	}
}

/*
*
Test `42` as an expression.

- Should NOT be optimized
- Should generate instructions to load constant 42
*/
func TestVisitIntLiteral_Expression_NotOptimized(t *testing.T) {
	intLiteral := makeIntLiteral(42, false, 0, 1, 1)

	visitor := NewInstructionsVisitor()
	visitor.EnterModuleContext()

	result := intLiteral.Visit(visitor)

	if result == nil {
		t.Fatal("expected non-nil result for expression")
	}

	visitResult, ok := CastVisitExprResult(result)
	if !ok {
		t.Fatal("expected VisitExprResult")
	}

	if !visitResult.TypeOf.IsEqual(TypeInt()) {
		t.Errorf("expected type VAL_INT, got %s", visitResult.TypeOf)
	}

	moduleContext := CastModuleContext(visitor.context)
	instructions := moduleContext.instructions

	if len(instructions) != 1 {
		t.Errorf("expected 1 instruction for expression, got %d", len(instructions))
	}

	if len(instructions) > 0 && instructions[0].OpCode != LOAD_CONST {
		t.Errorf("expected LOAD_CONST instruction, got %s", instructions[0].OpCode)
	}
}
