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
	DeferStatement     = "deferNode"

	LabelStatement = "labelNode"
	JumpStatement  = "jumpNode"

	NamedTypeDeclStatement = "namedTypeDecl"
	StructureDeclStatement = "structDecl"
	TraitDeclStatement     = "traitDecl"
	ImplDeclStatement      = "implDecl"
	FunctionProtoStatement = "funcProtoDecl"
	FunctionDeclStatement  = "funcDecl"
)

type NamedType struct {
	Mutable bool      `json:"mutable"`
	Name    Token     `json:"name"`
	Owned   bool      `json:"owned"`
	Type    *TypeNode `json:"type"`
}

type BlockNode struct {
	// hm
	Statements []*ParseTreeNode `json:"statements"`
}

//
type ElseIfNode struct {
	Cond  *ExpressionNode `json:"cond"`
	Block *BlockNode      `json:"block"`
}

// "func" iden "(" args ")"
type FunctionPrototypeDeclaration struct {
	Name      Token        `json:"name"`
	Arguments []*NamedType `json:"arguments"`

	// TODO should this be set to anything by
	// default, e.g. we can inject a "void"
	// into here?
	ReturnType *TypeNode `json:"return_type"`
}

// [ FunctionPrototypeDeclaration ] "{" { Stat ";" } "}"
type FunctionDeclaration struct {
	*FunctionPrototypeDeclaration
	Body *BlockNode `json:"body"`
}

type LetStatementNode struct {
	Name  Token           `json:"name"`
	Owned bool            `json:"owned"`
	Type  *TypeNode       `json:"type"`
	Value *ExpressionNode `json:"value"`
}
type MutableStatementNode struct {
	Name  Token           `json:"name"`
	Owned bool            `json:"owned"`
	Type  *TypeNode       `json:"type"`
	Value *ExpressionNode `json:"value"`
}
type ReturnStatementNode struct {
	Value *ExpressionNode `json:"value"`
}

// CONSTR

type WhileLoopNode struct {
	Cond  *ExpressionNode `json:"cond"`
	Post  *ExpressionNode `json:"post,omitempty"`
	Block *BlockNode      `json:"block"`
}

type LoopNode struct {
	Block *BlockNode `json:"block"`
}

type IfNode struct {
	Cond    *ExpressionNode `json:"cond"`
	Block   *BlockNode      `json:"block"`
	ElseIfs []*ElseIfNode   `json:"else_ifs"`
	Else    *BlockNode      `json:"else"`
}

type DeferNode struct {
	Block     *BlockNode     `json:"block"`
	Statement *ParseTreeNode `json:"statement"`
}

// DECL

// "struct" iden { ... }
type StructureDeclaration struct {
	Name   Token        `json:"name"`
	Fields []*NamedType `json:"fields"`
}

// "trait" iden { ... }
type TraitDeclaration struct {
	Name    Token                           `json:"name"`
	Members []*FunctionPrototypeDeclaration `json:"members"`
}

// todo
type ImplDeclaration struct {
	Name      Token                  `json:"name"`
	Functions []*FunctionDeclaration `json:"functions"`
}

type LabelNode struct {
	LabelName Token `json:"label_name"`
}

type JumpNode struct {
	Location Token `json:"location"`
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
	Kind StatementType `json:"kind"`

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
	DeferNode     *DeferNode     `json:"deferNode,omitempty"`

	// JUMP STUFF
	LabelNode *LabelNode `json:"labelNode,omitempty"`
	JumpNode  *JumpNode  `json:"jumpNode,omitempty"`

	// DECL

	StructureDeclaration *StructureDeclaration `json:"structureDecl,omitempty"`
	TraitDeclaration     *TraitDeclaration     `json:"traitDecl,omitempty"`
	ImplDeclaration      *ImplDeclaration      `json:"implDecl,omitempty"`

	NamedType                    *NamedType                    `json:"namedTypeDecl,omitempty"`
	FunctionPrototypeDeclaration *FunctionPrototypeDeclaration `json:"funcProtoDecl,omitempty"`
	FunctionDeclaration          *FunctionDeclaration          `json:"funcDecl,omitempty"`
}
