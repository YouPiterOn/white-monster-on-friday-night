package ast

type Visitor[R any] interface {
	VisitAssignment(n *Assignment) R
	VisitReturn(n *Return) R
	VisitNumberLiteral(n *NumberLiteral) R
	VisitIdentifier(n *Identifier) R
	VisitBinaryExpr(n *BinaryExpr) R
}
