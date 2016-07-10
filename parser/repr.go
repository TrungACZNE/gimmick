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

/* --- String representation of AST tokens --- */

func (token IdentifierNode) String() string {
	return fmt.Sprintf("{Identifier:%s}", token.Name)
}

func (token IntegerLiteralNode) String() string {
	return fmt.Sprintf("{Int:%d}", token.Value)
}

func (token FloatLiteralNode) String() string {
	return fmt.Sprintf("{Float:%v}", token.Value)
}

func (token ParamListToken) String() string {
	buf := ""
	for i, token := range token.ParamList {
		if i > 0 {
			buf += ", "
		}
		buf += token.String()
	}
	return fmt.Sprintf("{ParamList:%s}", buf)
}

func (token FunctionDefNode) String() string {
	return fmt.Sprintf("{FunctionDef:%s:%s:%s}", token.Name, NameTypeArrString(token.ArgList), token.Block.String())
}

func (token FunctionCallNode) String() string {
	return fmt.Sprintf("{FunctionCall:%s:%s}", token.Name, token.ParamList)
}

func (token BinaryOperatorNode) String() string {
	return fmt.Sprintf("{BinaryOperatorNode:%s:%s:%s}", token.Left.String(), token.Operator, token.Right.String())
}

func (token AssignmentNode) String() string {
	return fmt.Sprintf("{AssignmentNode:%s:%s}", token.Dest, token.Expr.String())
}

func (token BlockNode) String() string {
	buf := ""
	for i, token := range token.ExprList {
		if i > 0 {
			buf += ", "
		}
		buf += token.String()
	}
	return fmt.Sprintf("{BlockNode:%s}", NodeArrString(token.ExprList))
}

func (token ModuleNode) String() string {
	return token.Block.String()
}
