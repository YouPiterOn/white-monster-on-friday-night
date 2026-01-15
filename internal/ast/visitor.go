package ast

type Visitor[R any] interface {
	VisitProgram(n *Program) R
	VisitDeclaration(n *Declaration) R
	VisitAssignment(n *Assignment) R
	VisitReturn(n *Return) R
	VisitIntLiteral(n *IntLiteral) R
	VisitFloatLiteral(n *FloatLiteral) R
	VisitStringLiteral(n *StringLiteral) R
	VisitBoolLiteral(n *BoolLiteral) R
	VisitNullLiteral(n *NullLiteral) R
	VisitArrayLiteral(n *ArrayLiteral) R
	VisitIdentifier(n *Identifier) R
	VisitBinaryExpr(n *BinaryExpr) R
	VisitIndexExpr(n *IndexExpr) R
	VisitParam(n *Param) R
	VisitFunction(n *Function) R
	VisitBlock(n *Block) R
	VisitCallExpr(n *CallExpr) R
	VisitMemberExpr(n *MemberExpr) R
	VisitIf(n *If) R
}
