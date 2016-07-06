package parser

import "fmt"

/* --- String representation of AST nodes --- */

func (node EOFNode) String() string {
	return "<EOF>"
}

func (node EmptyNode) String() string {
	return "<Empty>"
}

func (node IntegerLiteralNode) String() string {
	return fmt.Sprintf("<Int:%d>", node.Value)
}

func (node FloatLiteralNode) String() string {
	return fmt.Sprintf("<Int:%v>", node.Value)
}

func (node KeywordNode) String() string {
	return fmt.Sprintf("<Keyword:%s>", node.Name)
}

func (node CharNode) String() string {
	return fmt.Sprintf("<Char:%s>", node.Name)
}

func (node IdentifierNode) String() string {
	return fmt.Sprintf("<Identifier:%s>", node.Name)
}

func (node ArgDeclNode) String() string {
	return fmt.Sprintf("<ArgDecl:%s:%s>", node.NameNode, node.TypeNode)
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

func (node FunctionDefNode) String() string {
	return fmt.Sprintf("<FunctionDef:%s:%s:%s>", node.Name.String(), node.ArgList.String(), node.Block.String())
}

func (node FunctionCallNode) String() string {
	return fmt.Sprintf("<FunctionCall:%s:%s>", node.Name.String(), node.ParamList.String())
}

func (node BinaryOperatorNode) String() string {
	return fmt.Sprintf("<BinaryOperatorNode:%s:%s:%s>", node.Left.String(), node.Operator.String(), node.Right.String())
}

func (node AssignmentNode) String() string {
	return fmt.Sprintf("<AssignmentNode:%s:%s>", node.Dest.String(), node.Expr.String())
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
