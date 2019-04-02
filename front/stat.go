package front

type StatementType string

const (
	LetStatement        StatementType = "letStat"
	MutableStatement                  = "mutStat"
	ReturnStatement                   = "retStat"
	BlockStatement                    = "blockStat"
	NextStatement                     = "nextStat"
	BreakStatement                    = "breakStat"
	ExpressionStatement               = "exprStat"

	WhileLoopStatement = "whileNode"
	LoopStatement      = "loopNode"
	ElseIfStatement    = "elseIfNode"
	IfStatement        = "ifNode"

	NamedTypeDeclStatement = "namedTypeDecl"
	StructureDeclStatement = "structDecl"
	TraitDeclStatement     = "traitDecl"
	ImplDeclStatement      = "implDecl"
	FunctionProtoStatement = "funcProtoDecl"
	FunctionDeclStatement  = "funcDecl"
)

type NamedType struct {
	Name Token
	Type *TypeNode
}

type BlockNode struct {
	// hm
	Statements []*ParseTreeNode
}

//
type ElseIfNode struct {
	Cond  *ExpressionNode
	Block *BlockNode
}

// "func" iden "(" args ")"
type FunctionPrototypeDeclaration struct {
	Name      Token
	Arguments []*NamedType

	// TODO should this be set to anything by
	// default, e.g. we can inject a "void"
	// into here?
	ReturnType *TypeNode
}

// [ FunctionPrototypeDeclaration ] "{" { Stat ";" } "}"
type FunctionDeclaration struct {
	*FunctionPrototypeDeclaration
	Body *BlockNode
}

type LetStatementNode struct {
	Name  Token
	Type  *TypeNode
	Value *ExpressionNode
}
type MutableStatementNode struct {
	Name  Token
	Type  *TypeNode
	Value *ExpressionNode
}
type ReturnStatementNode struct {
	Value *ExpressionNode
}

// CONSTR

type WhileLoopNode struct {
	Cond  *ExpressionNode
	Post  *ExpressionNode
	Block *BlockNode
}

type LoopNode struct {
	Block *BlockNode
}

type IfNode struct {
	Cond    *ExpressionNode
	Block   *BlockNode
	ElseIfs []*ElseIfNode
	Else    *BlockNode
}

// DECL

// "struct" iden { ... }
type StructureDeclaration struct {
	Name   Token
	Fields []*NamedType
}

// "trait" iden { ... }
type TraitDeclaration struct {
	Name    Token
	Members []*FunctionPrototypeDeclaration
}

// todo
type ImplDeclaration struct {
	Name      Token
	Functions []*FunctionDeclaration
}

// ParseTreeNode is a big jumbo node containing all of the
// node combinations.
//
// We did use inheritance here, but this doesn't serialize into
// JSON unless I implement the JSON Serialisation for each node,
// which is quite tedious. So instead, I'm opting for the C-like
// union approach, though Go doesn't support unions so this will be
// a relatively large struct.
type ParseTreeNode struct {
	Kind StatementType

	LetStatementNode        *LetStatementNode     `json:"letStatement,omitempty"`
	MutableStatementNode    *MutableStatementNode `json:"mutStatement,omitempty"`
	ReturnStatementNode     *ReturnStatementNode  `json:"retStatement,omitempty"`
	ExpressionStatementNode *ExpressionNode       `json:"exprStatement,omitempty"`

	// CONSTR

	WhileLoopNode *WhileLoopNode `json:"whileNode,omitempty"`
	LoopNode      *LoopNode      `json:"loopNode,omitempty"`
	ElseIfNode    *ElseIfNode    `json:"elseIfNode,omitempty"`
	BlockNode     *BlockNode     `json:"blockNode,omitempty"`
	IfNode        *IfNode        `json:"ifNode,omitempty"`

	// DECL

	StructureDeclaration *StructureDeclaration `json:"structureDecl,omitempty"`
	TraitDeclaration     *TraitDeclaration     `json:"traitDecl,omitempty"`
	ImplDeclaration      *ImplDeclaration      `json:"implDecl,omitempty"`

	NamedType                    *NamedType                    `json:"namedTypeDecl,omitempty"`
	FunctionPrototypeDeclaration *FunctionPrototypeDeclaration `json:"funcProtoDecl,omitempty"`
	FunctionDeclaration          *FunctionDeclaration          `json:"funcDecl,omitempty"`
}
