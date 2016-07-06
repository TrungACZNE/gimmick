package parser

import . "github.com/trungaczne/gimmick/vm"

/* --- AST Node definitions ---*/

type Node interface {
	String() string
	CodeGen(builder CodeBuilder)
}

type EOFNode struct {
}

type EmptyNode struct {
}

type IntegerLiteralNode struct {
	Value int64
}

type FloatLiteralNode struct {
	Value float64
}

type KeywordNode struct {
	Name string
}

type CharNode struct {
	Name string
}

type IdentifierNode struct {
	Name string
}

type ArgDeclNode struct {
	NameNode IdentifierNode
	TypeNode IdentifierNode
}

type ArgListNode struct {
	ArgDecl []ArgDeclNode
}

type ParamListNode struct {
	ParamList []Node
}

type FunctionDefNode struct {
	Name    IdentifierNode
	ArgList ArgListNode
	Block   BlockNode
}

type FunctionCallNode struct {
	Name      IdentifierNode
	ParamList ParamListNode
}

type BinaryOperatorNode struct {
	Left     Node
	Operator CharNode
	Right    Node
}

type AssignmentNode struct {
	Dest IdentifierNode
	Expr Node
}

type BlockNode struct {
	ExprList []Node
}
