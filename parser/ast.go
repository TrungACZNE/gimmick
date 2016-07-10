package parser

import . "github.com/trungaczne/gimmick/vm"

/* --- AST Token definitions ---*/

type Token interface {
	String() string
	CodeGen(builder CodeBuilder)
}

type EOFToken struct {
}

type EmptyToken struct {
}

type IntegerLiteralToken struct {
	Value int64
}

type FloatLiteralToken struct {
	Value float64
}

type KeywordToken struct {
	Name string
}

type CharToken struct {
	Name string
}

type IdentifierToken struct {
	Name string
}

type ArgDeclToken struct {
	NameToken IdentifierToken
	TypeToken IdentifierToken
}

type ArgListToken struct {
	ArgDecl []ArgDeclToken
}

type ParamListToken struct {
	ParamList []Token
}

type FunctionDefToken struct {
	Name    IdentifierToken
	ArgList ArgListToken
	Block   BlockToken
}

type FunctionCallToken struct {
	Name      IdentifierToken
	ParamList ParamListToken
}

type BinaryOperatorToken struct {
	Left     Token
	Operator CharToken
	Right    Token
}

type AssignmentToken struct {
	Dest IdentifierToken
	Expr Token
}

type BlockToken struct {
	ExprList []Token
}
