package parser

import "fmt"

// Node types

type Node interface {
	String() string
}

type EOFNode struct{}

func (node EOFNode) String() string {
	return "<EOF>"
}

type ExpressionNode struct{}

func (node ExpressionNode) String() string {
	return "<Expr>"
}

type IntegerLiteralNode struct {
	Value int64
}

func (node IntegerLiteralNode) String() string {
	return fmt.Sprintf("<Int:%d>", node.Value)
}

type FloatLiteralNode struct {
	Value float64
}

func (node FloatLiteralNode) String() string {
	return fmt.Sprintf("<Int:%v>", node.Value)
}

type FuncDefNode struct{}

func (node FuncDefNode) String() string {
	return "<FuncDef>"
}

type FuncCallNode struct{}

func (node FuncCallNode) String() string {
	return "<FuncCall>"
}

type KeywordNode struct {
	Name string
}

func (node KeywordNode) String() string {
	return fmt.Sprintf("<Keyword:%s>", node.Name)
}

type SymbolNode struct {
	Name string
}

func (node SymbolNode) String() string {
	return fmt.Sprintf("<Symbol:%s>", node.Name)
}

type IdentifierNode struct {
	Name string
}

func (node IdentifierNode) String() string {
	return fmt.Sprintf("<Identifier:%s>", node.Name)
}

/*
// array of nodes, not to be used directly
func (nodes []Node) String() string {
	buf := "<"
	for i, node := range nodes {
		if i > 0 {
			buf += ", "
		}
		buf += node.String()
	}
	return buf + ">"
}*/

type ArgDeclNode struct {
	NameNode IdentifierNode
	TypeNode IdentifierNode
}

func (node ArgDeclNode) String() string {
	return fmt.Sprintf("<ArgDecl:%s:%s>", node.NameNode, node.TypeNode)
}

type ArgListNode struct {
	ArgDecl []ArgDeclNode
}

func (node ArgListNode) String() string {
	buf := ""
	for i, node := range node.ArgDecl {
		if i > 0 {
			buf += ", "
		}
		buf += node.String()
	}
	return fmt.Sprintf("<ArgList:%s>", buf)
}

type ParamListNode struct {
	ParamList []Node
}

func (node ParamListNode) String() string {
	buf := ""
	for i, node := range node.ParamList {
		if i > 0 {
			buf += ", "
		}
		buf += node.String()
	}
	return fmt.Sprintf("<ParamList:%s>", buf)
}

type EmptyNode struct{}

func (node EmptyNode) String() string {
	return "<Empty>"
}

type FunctionDefNode struct {
	Name    IdentifierNode
	ArgList ArgListNode
	Block   BlockNode
}

func (node FunctionDefNode) String() string {
	return fmt.Sprintf("<FunctionDef:%s:%s:%s>", node.Name.String(), node.ArgList.String(), node.Block.String())
}

type FunctionCallNode struct {
	Name      IdentifierNode
	ParamList ParamListNode
}

func (node FunctionCallNode) String() string {
	return fmt.Sprintf("<FunctionCall:%s:%s>", node.Name.String(), node.ParamList.String())
}

type BinaryOperatorNode struct {
	Left     Node
	Operator SymbolNode
	Right    Node
}

func (node BinaryOperatorNode) String() string {
	return fmt.Sprintf("<BinaryOperatorNode:%s:%s:%s>", node.Left.String(), node.Operator.String(), node.Right.String())
}

type BlockNode struct {
	ExprList []Node
}

func (node BlockNode) String() string {
	buf := ""
	for i, node := range node.ExprList {
		if i > 0 {
			buf += ", "
		}
		buf += node.String()
	}
	return fmt.Sprintf("<BlockNode:%s>", buf)
}
