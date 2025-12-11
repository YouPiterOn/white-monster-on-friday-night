package ast

type Visitor[R any] interface {
	VisitProgram(n *Program) R
	VisitDeclaration(n *Declaration) R
	VisitAssignment(n *Assignment) R
	VisitReturn(n *Return) R
	VisitNumberLiteral(n *NumberLiteral) R
	VisitIdentifier(n *Identifier) R
	VisitBinaryExpr(n *BinaryExpr) R
	VisitParam(n *Param) R
	VisitFunction(n *Function) R
	VisitBlock(n *Block) R
	VisitCallExpr(n *CallExpr) R
}
