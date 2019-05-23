package front

type ExpressionType string

const (
	BuiltinExpression     ExpressionType = "builtinExpr"
	VariableExpression                   = "variableExpr"
	LiteralExpression                    = "literalExpr"
	UnaryExpression                      = "unaryExpr"
	BinaryExpression                     = "binaryExpr"
	Grouping                             = "groupingExpr"
	IndexExpression                      = "indexExpr"
	CallExpression                       = "callExpr"
	PathExpression                       = "pathExpr"
	AssignStatement                      = "assignExpr"
	ConstantExpression                   = "constExpr"
	LambdaExpression                     = "lambdaExpr"
	ListExpression                       = "listExpr"
	InitializerExpression                = "initExpr"
	TypeExpression                       = "typeExpr"
)

type LambdaExpressionNode struct {
	Proto *FunctionPrototypeDeclaration
	Body  *BlockNode
}
type BuiltinExpressionNode struct {
	Name string
	Iden *VariableReferenceNode
	Args []*ExpressionNode
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
type ExprList struct {
	Values []*ExpressionNode
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

type InitializerKind string

const (
	InitStructure InitializerKind = "init-struct"
	InitTuple                     = "init-tuple"
	InitArray                     = "init-array"
)

type InitializerExpressionNode struct {
	Kind InitializerKind

	// only set when Kind == InitStructure
	// this is the Structure name usually.
	LHand Token

	Values []*ExpressionNode
}

type ExpressionNode struct {
	Kind ExpressionType

	LambdaExpressionNode      *LambdaExpressionNode      `json:"lambdaExpr,omitempty"`
	BuiltinExpressionNode     *BuiltinExpressionNode     `json:"builtinExpr,omitempty"`
	VariableExpressionNode    *VariableExpressionNode    `json:"variableExpr,omitempty"`
	LiteralExpressionNode     *LiteralExpressionNode     `json:"literalExpr,omitempty"`
	UnaryExpressionNode       *UnaryExpressionNode       `json:"unaryExpr,omitempty"`
	BinaryExpressionNode      *BinaryExpressionNode      `json:"binaryExpr,omitempty"`
	ExprList                  *ExprList                  `json:"exprList,omitempty"`
	GroupingNode              *GroupingNode              `json:"groupingExpr,omitempty"`
	IndexExpressionNode       *IndexExpressionNode       `json:"indexExpr,omitempty"`
	CallExpressionNode        *CallExpressionNode        `json:"callExpr,omitempty"`
	PathExpressionNode        *PathExpressionNode        `json:"pathExpr,omitempty"`
	ConstantNode              *ConstantNode              `json:"constExpr,omitempty"`
	AssignStatementNode       *AssignStatementNode       `json:"assignExpr,omitempty"`
	InitializerExpressionNode *InitializerExpressionNode `json:"initExpr,omitempty"`
	TypeExpressionNode        *TypeNode                  `json:"typeExpr,omitEmpty"`
}
