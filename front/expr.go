package front

type ExpressionType int

const (
	BuiltinExpression ExpressionType = iota
	VariableExpression
	LiteralExpression
	UnaryExpression
	BinaryExpression
	Grouping
	IndexExpression
	CallExpression
	PathExpression
	AssignStatement
	ConstantExpression
)

type BuiltinExpressionNode struct {
	Name string
	Type *TypeNode
}
type VariableExpressionNode struct {
	Name string
}
type LiteralExpressionNode struct {
	Value string
}
type UnaryExpressionNode struct {
	Operator string
	Value    *ExpressionNode
}
type BinaryExpressionNode struct {
	LHand    *ExpressionNode
	Operator string
	RHand    *ExpressionNode
}
type GroupingNode struct {
	Value *ExpressionNode
}
type IndexExpressionNode struct {
	Left  *ExpressionNode
	Value *ExpressionNode
}
type CallExpressionNode struct {
	Left   *ExpressionNode
	Params []*ExpressionNode
}
type PathExpressionNode struct {
	Values []*ExpressionNode
}
type AssignStatementNode struct {
	LHand *ExpressionNode
	Op    string
	RHand *ExpressionNode
}

type ExpressionNode struct {
	Kind ExpressionType

	BuiltinExpressionNode  *BuiltinExpressionNode
	VariableExpressionNode *VariableExpressionNode
	LiteralExpressionNode  *LiteralExpressionNode
	UnaryExpressionNode    *UnaryExpressionNode
	BinaryExpressionNode   *BinaryExpressionNode
	GroupingNode           *GroupingNode
	IndexExpressionNode    *IndexExpressionNode
	CallExpressionNode     *CallExpressionNode
	PathExpressionNode     *PathExpressionNode
	ConstantNode           *ConstantNode
	AssignStatementNode    *AssignStatementNode
}
