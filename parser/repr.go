package parser

import "fmt"

/* --- String representation of AST tokens --- */

func (token EOFToken) String() string {
	return "<EOF>"
}

func (token EmptyToken) String() string {
	return "<Empty>"
}

func (token IntegerLiteralToken) String() string {
	return fmt.Sprintf("<Int:%d>", token.Value)
}

func (token FloatLiteralToken) String() string {
	return fmt.Sprintf("<Int:%v>", token.Value)
}

func (token KeywordToken) String() string {
	return fmt.Sprintf("<Keyword:%s>", token.Name)
}

func (token CharToken) String() string {
	return fmt.Sprintf("<Char:%s>", token.Name)
}

func (token IdentifierToken) String() string {
	return fmt.Sprintf("<Identifier:%s>", token.Name)
}

func (token ArgDeclToken) String() string {
	return fmt.Sprintf("<ArgDecl:%s:%s>", token.NameToken, token.TypeToken)
}

func (token ArgListToken) String() string {
	buf := ""
	for i, token := range token.ArgDecl {
		if i > 0 {
			buf += ", "
		}
		buf += token.String()
	}
	return fmt.Sprintf("<ArgList:%s>", buf)
}

func (token ParamListToken) String() string {
	buf := ""
	for i, token := range token.ParamList {
		if i > 0 {
			buf += ", "
		}
		buf += token.String()
	}
	return fmt.Sprintf("<ParamList:%s>", buf)
}

func (token FunctionDefToken) String() string {
	return fmt.Sprintf("<FunctionDef:%s:%s:%s>", token.Name.String(), token.ArgList.String(), token.Block.String())
}

func (token FunctionCallToken) String() string {
	return fmt.Sprintf("<FunctionCall:%s:%s>", token.Name.String(), token.ParamList.String())
}

func (token BinaryOperatorToken) String() string {
	return fmt.Sprintf("<BinaryOperatorToken:%s:%s:%s>", token.Left.String(), token.Operator.String(), token.Right.String())
}

func (token AssignmentToken) String() string {
	return fmt.Sprintf("<AssignmentToken:%s:%s>", token.Dest.String(), token.Expr.String())
}

func (token BlockToken) String() string {
	buf := ""
	for i, token := range token.ExprList {
		if i > 0 {
			buf += ", "
		}
		buf += token.String()
	}
	return fmt.Sprintf("<BlockToken:%s>", buf)
}
