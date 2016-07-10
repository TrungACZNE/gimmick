package parser

import . "github.com/trungaczne/gimmick/vm"

type Token interface {
}

type Node interface {
	Token
	String() string
	CodeGen(builder CodeBuilder)
}

/* --- Tokens ---*/

type EOFToken struct {
}

type EmptyToken struct {
}

type KeywordToken struct {
	Name string
}

type CharToken struct {
	Name string
}

type ArgDeclToken struct {
	NameToken IdentifierNode
	TypeToken IdentifierNode
}

type ArgListToken struct {
	ArgDecl []ArgDeclToken
}

type ParamListToken struct {
	ParamList []Node
}

/* --- Nodes ---*/

type IntegerLiteralNode struct {
	Value int64
}

type FloatLiteralNode struct {
	Value float64
}

type IdentifierNode struct {
	Name string
}

type FunctionDefNode struct {
	Name    string
	ArgList []NameType
	Block   BlockNode
}

type FunctionCallNode struct {
	Name      string
	ParamList []Node
}

type BinaryOperatorNode struct {
	Left     Node
	Operator string
	Right    Node
}

type AssignmentNode struct {
	Dest string
	Expr Node
}

type BlockNode struct {
	ExprList []Node
}
