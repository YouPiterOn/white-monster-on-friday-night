package ast

type Visitor[R any] interface {
	VisitProgram(n *Program) R
	VisitDeclaration(n *Declaration) R
	VisitAssignment(n *Assignment) R
	VisitReturn(n *Return) R
	VisitIntLiteral(n *IntLiteral) R
	VisitBoolLiteral(n *BoolLiteral) R
	VisitNullLiteral(n *NullLiteral) R
	VisitIdentifier(n *Identifier) R
	VisitBinaryExpr(n *BinaryExpr) R
	VisitParam(n *Param) R
	VisitFunction(n *Function) R
	VisitBlock(n *Block) R
	VisitCallExpr(n *CallExpr) R
	VisitIf(n *If) R
}
