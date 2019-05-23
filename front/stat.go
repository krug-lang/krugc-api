package front

import "fmt"

// StatementType ...
type StatementType string

// ...
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

	TypeAliasStatement     = "typeAliasDecl"
	TraitDeclStatement     = "traitDecl"
	ImplDeclStatement      = "implDecl"
	FunctionProtoStatement = "funcProtoDecl"
	FunctionDeclStatement  = "funcDecl"
	StructureDeclStatement = "structDecl"
)

// NamedType ...
type NamedType struct {
	Mutable bool            `json:"mutable"`
	Name    Token           `json:"name"`
	Owned   bool            `json:"owned"`
	Type    *ExpressionNode `json:"type"`
}

// BlockNode ...
type BlockNode struct {
	// hm
	Statements []*ParseTreeNode `json:"statements"`
}

// ElseIfNode ...
type ElseIfNode struct {
	Cond  *ExpressionNode `json:"cond"`
	Block *BlockNode      `json:"block"`
}

// FunctionPrototypeDeclaration ...
// "func" iden "(" args ")"
type FunctionPrototypeDeclaration struct {
	Name       Token           `json:"name"`
	Arguments  []*NamedType    `json:"arguments"`
	ReturnType *ExpressionNode `json:"return_type"`
}

// FunctionDeclaration ...
// [ FunctionPrototypeDeclaration ] "{" { Stat ";" } "}"
type FunctionDeclaration struct {
	*FunctionPrototypeDeclaration
	Body *BlockNode `json:"body"`
}

// TypeAliasNode ...
// type name = type;
type TypeAliasNode struct {
	Name Token           `json:"name"`
	Type *ExpressionNode `json:"type"`
}

func (t *TypeAliasNode) String() string {
	return fmt.Sprintf("type %s = %s", t.Name.Value, t.Type.Kind)
}

// LetStatementNode ...
type LetStatementNode struct {
	Name  Token           `json:"name"`
	Owned bool            `json:"owned"`
	Type  *ExpressionNode `json:"type"`
	Value *ExpressionNode `json:"value"`
}

// MutableStatementNode ...
type MutableStatementNode struct {
	Name  Token           `json:"name"`
	Owned bool            `json:"owned"`
	Type  *ExpressionNode `json:"type"`
	Value *ExpressionNode `json:"value"`
}

// ReturnStatementNode ...
type ReturnStatementNode struct {
	Value *ExpressionNode `json:"value"`
}

// CONSTR

// WhileLoopNode ...
type WhileLoopNode struct {
	Cond  *ExpressionNode `json:"cond"`
	Post  *ExpressionNode `json:"post,omitempty"`
	Block *BlockNode      `json:"block"`
}

// LoopNode ...
type LoopNode struct {
	Block *BlockNode `json:"block"`
}

// IfNode ...
type IfNode struct {
	Cond    *ExpressionNode `json:"cond"`
	Block   *BlockNode      `json:"block"`
	ElseIfs []*ElseIfNode   `json:"else_ifs"`
	Else    *BlockNode      `json:"else"`
}

// DeferNode ...
type DeferNode struct {
	Block     *BlockNode     `json:"block"`
	Statement *ParseTreeNode `json:"statement"`
}

// DECL

// StructureDeclaration ...
// see StructureTypeNode
type StructureDeclaration struct {
	Structure *StructureTypeNode
}

// TraitDeclaration ...
// "trait" iden { ... }
type TraitDeclaration struct {
	Name    Token                           `json:"name"`
	Members []*FunctionPrototypeDeclaration `json:"members"`
}

// ImplDeclaration ...
type ImplDeclaration struct {
	Name      Token                  `json:"name"`
	Functions []*FunctionDeclaration `json:"functions"`
}

// LabelNode ...
type LabelNode struct {
	LabelName Token `json:"label_name"`
}

// JumpNode ...
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

	TypeAliasNode           *TypeAliasNode        `json:"namedType,omitempty"`
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

	TraitDeclaration *TraitDeclaration `json:"traitDecl,omitempty"`
	ImplDeclaration  *ImplDeclaration  `json:"implDecl,omitempty"`

	NamedType                    *NamedType                    `json:"namedTypeDecl,omitempty"`
	FunctionPrototypeDeclaration *FunctionPrototypeDeclaration `json:"funcProtoDecl,omitempty"`
	FunctionDeclaration          *FunctionDeclaration          `json:"funcDecl,omitempty"`
	StructureDeclaration         *StructureDeclaration         `json:"structDecl,omitempty"`
}
