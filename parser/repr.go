package parser

import . "github.com/trungaczne/gimmick/vm"
import "fmt"

func NodeArrString(nodes []Node) string {
	buf := "["
	for i, node := range nodes {
		buf += node.String()
		if i != len(nodes)-1 {
			buf += ","
		}
	}
	return buf + "]"
}

func NameTypeArrString(nt []NameType) string {
	buf := "["
	for i, nametype := range nt {
		buf += nametype.Name + ":" + nametype.Type
		if i != len(nt)-1 {
			buf += ","
		}
	}
	return buf + "]"
}

func PrettyPrint(str string, indentWidth int) string {
	curIndent := 0
	buf := ""
	for _, c := range str {
		if c == '{' {
			curIndent += indentWidth
			buf += string(c) + "\n" + fmt.Sprintf(fmt.Sprintf("%%%ds", curIndent), " ")
		} else if c == '}' {
			curIndent -= indentWidth
			buf += "\n" + fmt.Sprintf(fmt.Sprintf("%%%ds", curIndent), " ") + string(c)
		} else {
			buf += string(c)
		}
	}
	return buf
}

/* --- String representation of AST nodes --- */

func (node IdentifierNode) String() string {
	return fmt.Sprintf("{Identifier:%s}", node.Name)
}

func (node IntegerLiteralNode) String() string {
	return fmt.Sprintf("{Int:%d}", node.Value)
}

func (node FloatLiteralNode) String() string {
	return fmt.Sprintf("{Float:%v}", node.Value)
}

func (node ParamListToken) String() string {
	buf := ""
	for i, node := range node.ParamList {
		if i > 0 {
			buf += ", "
		}
		buf += node.String()
	}
	return fmt.Sprintf("{ParamList:%s}", buf)
}

func (node FunctionDefNode) String() string {
	return fmt.Sprintf("{FunctionDef:%s:%s:%s}", node.Name, NameTypeArrString(node.ArgList), node.Block.String())
}

func (node FunctionCallNode) String() string {
	return fmt.Sprintf("{FunctionCall:%s:%s}", node.Name, node.ParamList)
}

func (node BinaryOperatorNode) String() string {
	return fmt.Sprintf("{BinaryOperatorNode:%s:%s:%s}", node.Left.String(), node.Operator, node.Right.String())
}

func (node AssignmentNode) String() string {
	return fmt.Sprintf("{AssignmentNode:%s:%s}", node.Dest, node.Expr.String())
}

func (node BlockNode) String() string {
	buf := ""
	for i, node := range node.ExprList {
		if i > 0 {
			buf += ", "
		}
		buf += node.String()
	}
	return fmt.Sprintf("{BlockNode:%s}", NodeArrString(node.ExprList))
}

func (node ModuleNode) String() string {
	return node.Block.String()
}
