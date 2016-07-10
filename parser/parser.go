package parser

import (
	"fmt"
	"regexp"
	"strconv"

	. "github.com/trungaczne/gimmick/vm"
)

// TODO use io.File instead of String
// TODO handle Unicode (haha nice joke)

type Parser struct {
	text string
}

func NewParser(text string) *Parser {
	return &Parser{text: text}
}

/* --- Errors --- */

type EOFError int

func (err EOFError) Error() string {
	return fmt.Sprintf("Unexpected EOF")
}

type NotMatchError string

func (err NotMatchError) Error() string {
	return fmt.Sprintf("Could not match type: " + string(err))
}

func isWhiteSpace(c byte) bool {
	return c == ' ' || c == '\n' || c == '\t'
}

/* --- Parser routines --- */
// finds the next non whitespace character position
// returns -1 if none could be found
func (p *Parser) findNonWhiteSpace(cursor int) int {
	if cursor >= len(p.text) {
		return -1
	}
	for ; cursor < len(p.text) && isWhiteSpace(p.text[cursor]); cursor += 1 {
	}
	if cursor >= len(p.text) || isWhiteSpace(p.text[cursor]) {
		return -1
	}
	return cursor
}

// fetch returns the next non-whitespace character in the text
func (p *Parser) fetch(cursor int) (byte, int, error) {
	pos := p.findNonWhiteSpace(cursor)
	if pos == -1 {
		return '?', cursor, EOFError(0)
	}
	return p.text[pos], pos + 1, nil
}

/* --- Routines to build matchers --- */
type TryFunc func(parser *Parser, cursor int) (Token, int, error)

// MatchAllWrapper convert an array of token into a single token
type MatchAllWrapper func(tokens []Token) Token

func MatchAll(wrapper MatchAllWrapper, defs ...TryFunc) TryFunc {
	return func(parser *Parser, cursor int) (Token, int, error) {
		newCursor := cursor
		tokens := []Token{}
		for _, f := range defs {
			token, _cursor, err := f(parser, newCursor)
			if err != nil {
				return nil, cursor, err
			}
			newCursor = _cursor
			tokens = append(tokens, token)
		}
		return wrapper(tokens), newCursor, nil
	}
}

// MatchOneOf can take multiple paths, which may return different
// token types. The wrapper's job is to cast them all back 1 type
type MatchOneOfWrapper func(token Token) Token

func MatchOneOf(wrapper MatchOneOfWrapper, defs ...TryFunc) TryFunc {
	return func(parser *Parser, cursor int) (Token, int, error) {
		var bestToken Token
		bestCursor := -1
		for _, f := range defs {
			token, newCursor, err := f(parser, cursor)
			if err == nil {
				if newCursor > bestCursor {
					bestCursor = newCursor
					bestToken = token
				}
			}
		}
		if bestCursor > -1 {
			return wrapper(bestToken), bestCursor, nil
		}
		return nil, cursor, NotMatchError("MatchOneOf")
	}
}

var REG_IDENTIFIER_INITIAL = regexp.MustCompile("[_a-zA-Z]")
var REG_IDENTIFIER = regexp.MustCompile("^[_a-zA-Z0-9]+")

func Identifier(parser *Parser, cursor int) (Token, int, error) {
	first, newCursor, err := parser.fetch(cursor)
	if err != nil || !REG_IDENTIFIER_INITIAL.Match([]byte{first}) {
		return nil, cursor, NotMatchError("Identifier")
	}
	if newCursor == len(parser.text) {
		return IdentifierNode{string(first)}, newCursor, nil
	}
	rest := REG_IDENTIFIER.FindString(parser.text[newCursor:len(parser.text)])
	return IdentifierNode{string(first) + rest}, newCursor + len(rest), nil
}

// shared by Keyword and Char
func tryString(p *Parser, cursor int, str string) (int, error) {
	cursor = p.findNonWhiteSpace(cursor)
	l := len(p.text)
	if cursor == -1 || cursor+len(str) > l || p.text[cursor:cursor+len(str)] != str {
		return cursor, fmt.Errorf("Can't match")
	}
	return cursor + len(str), nil
}

func keyword(str string) TryFunc {
	return func(parser *Parser, cursor int) (Token, int, error) {
		newCursor, err := tryString(parser, cursor, str)
		if err != nil {
			return nil, newCursor, NotMatchError("keyword")
		}
		return KeywordToken{str}, newCursor, nil
	}
}

func char(str string) TryFunc {
	return func(parser *Parser, cursor int) (Token, int, error) {
		newCursor, err := tryString(parser, cursor, str)
		if err != nil {
			return nil, newCursor, NotMatchError("char")
		}
		return CharToken{str}, newCursor, nil
	}
}

func AsArgDecl(tokens []Token) Token {
	if len(tokens) != 3 {
		panic(fmt.Sprintf("Should have 3 tokens: %v", tokens))
	}
	declName, ok1 := tokens[0].(IdentifierNode)
	declType, ok2 := tokens[2].(IdentifierNode)
	if !ok1 || !ok2 {
		panic("Typecasting failure")
	}
	return ArgDeclToken{declName, declType}
}

var ArgDecl = MatchAll(AsArgDecl, Identifier, char(":"), Identifier)

func Token2ArgListToken(token Token) Token {
	list := []ArgDeclToken{}
	switch arglistToken := token.(type) {
	default:
		list = nil
	case ArgListToken:
		list = append(list, arglistToken.ArgDecl...)
	case ArgDeclToken:
		list = append(list, arglistToken)
	case EmptyToken:
		// do nothing
	}
	return ArgListToken{list}
}

func AsArgList(tokens []Token) Token {
	if len(tokens) != 3 {
		panic(fmt.Sprintf("Should have 3 tokens: %v", tokens))
	}
	head, ok := tokens[0].(ArgDeclToken)

	if !ok {
		panic("Typecasting failure")
	}

	newList := []ArgDeclToken{head}

	token, ok := Token2ArgListToken(tokens[2]).(ArgListToken)
	if !ok {
		panic("Programming error")
	}

	tail := token.ArgDecl
	if tail == nil {
		panic("Typecasting failure")
	}
	newList = append(newList, tail...)
	return ArgListToken{newList}
}

func ArgList(parser *Parser, cursor int) (Token, int, error) {
	// TODO this accepts (arg,) which shouldn't be the case

	return MatchOneOf(
		Token2ArgListToken,
		MatchAll(AsArgList, ArgDecl, char(","), ArgList),
		ArgDecl,
		EmptyExpression,
	)(parser, cursor)
}

func AsFunctionDef(tokens []Token) Token {
	if len(tokens) != 8 {
		panic(fmt.Sprintf("Should have 8 tokens: %v", tokens))
	}
	name, ok1 := tokens[1].(IdentifierNode)
	arglist, ok2 := tokens[3].(ArgListToken)
	block, ok3 := tokens[6].(BlockNode)
	if !ok1 || !ok2 || !ok3 {
		fmt.Println(ok1, ok2, ok3)
		panic("Typecasting failure")
	}

	nametype := []NameType{}
	for _, v := range arglist.ArgDecl {
		nametype = append(nametype, NameType{v.NameToken.Name, v.TypeToken.Name})
	}
	return FunctionDefNode{name.Name, nametype, block}
}

func Token2ParamListToken(token Token) Token {
	list := []Node{}
	switch paramToken := token.(type) {
	default:
		// probably should wrap Expression instead of doing this
		node, ok := paramToken.(Node)
		if !ok {
			panic("Typecasting failure")
		}
		list = append(list, node)
	case ParamListToken:
		list = append(list, paramToken.ParamList...)
	case EmptyToken:
		// do nothing
	}
	return ParamListToken{list}
}

func AsParamList(tokens []Token) Token {
	if len(tokens) != 3 {
		panic(fmt.Sprintf("Should have 3 tokens: %v", tokens))
	}
	head, ok1 := tokens[0].(Node)
	tail, ok2 := Token2ParamListToken(tokens[2]).(ParamListToken)
	if !ok1 || !ok2 {
		panic("Typecasting failure")
	}
	newList := []Node{head}
	newList = append(newList, tail.ParamList...)
	return ParamListToken{newList}
}

func ParamList(parser *Parser, cursor int) (Token, int, error) {
	return MatchOneOf(
		Token2ParamListToken,
		MatchAll(
			AsParamList,
			Expression, char(","), ParamList,
		),
		Expression,
		EmptyExpression,
	)(parser, cursor)
}

func AsFunctionCall(tokens []Token) Token {
	if len(tokens) != 4 {
		panic(fmt.Sprintf("Should have 4 tokens: %v", tokens))
	}
	name, ok := tokens[0].(IdentifierNode)
	paramList := tokens[2].(ParamListToken)
	if !ok {
		panic("Typecasting failure")
	}
	return FunctionCallNode{name.Name, paramList.ParamList}
}

func identity(token Token) Token {
	return token
}

func AsBinaryOperator(tokens []Token) Token {
	if len(tokens) != 3 {
		panic(fmt.Sprintf("Should have 3 tokens: %v", tokens))
	}
	left, ok1 := tokens[0].(Node)
	operator, ok2 := tokens[1].(CharToken)
	right, ok3 := tokens[2].(Node)

	if !ok1 || !ok2 || !ok3 {
		panic("Typecasting failure")
	}
	return BinaryOperatorNode{left, operator.Name, right}
}

func AsBracketExpression(tokens []Token) Token {
	if len(tokens) != 3 {
		panic(fmt.Sprintf("Should have 3 tokens: %v", tokens))
	}
	return tokens[1]
}

func AsAssignment(tokens []Token) Token {
	if len(tokens) != 3 {
		panic(fmt.Sprintf("Should have 3 tokens: %v", tokens))
	}
	id, ok1 := tokens[0].(IdentifierNode)
	expr, ok2 := tokens[2].(Node)
	if !ok1 || !ok2 {
		panic("Typecasting failure")
	}

	return AssignmentNode{id.Name, expr}
}

func AsBlock(tokens []Token) Token {
	if len(tokens) != 2 {
		panic(fmt.Sprintf("Should have 2 tokens: %v", tokens))
	}
	head, ok1 := tokens[0].(Node)
	tail, ok2 := tokens[1].(BlockNode)
	if !ok1 || !ok2 {
		panic("Typecasting failure")
	}
	newList := []Node{}
	newList = append(newList, head)
	newList = append(newList, tail.ExprList...)
	return BlockNode{newList}
}

func Token2BlockNode(token Token) Token {
	list := []Node{}
	switch blockToken := token.(type) {
	default:
		panic("Typecasting failure")
	case BlockNode:
		list = append(list, blockToken.ExprList...)
	case Node:
		list = append(list, blockToken)
	case EmptyToken:
		// do nothing
	}
	return BlockNode{list}
}

/* --- Keywords --- */

var KEYWORD_DEF = keyword("def")

/* --- Matchers --- */

func EndOfFile(parser *Parser, cursor int) (Token, int, error) {
	_, newCursor, err := parser.fetch(cursor)
	if _, ok := err.(EOFError); ok {
		return EOFToken{}, newCursor, nil
	}
	return nil, cursor, NotMatchError("EOF")
}

func EmptyExpression(parser *Parser, cursor int) (Token, int, error) {
	return EmptyToken{}, cursor, nil
}

func Block(parser *Parser, cursor int) (Token, int, error) {
	return MatchOneOf(
		Token2BlockNode,
		MatchAll(AsBlock, Expression, Block),
		Expression,
		EmptyExpression,
	)(parser, cursor)
}

func Expression(parser *Parser, cursor int) (Token, int, error) {
	return MatchOneOf(
		identity,
		MatchAll(
			AsBinaryOperator,
			GuardedExpression, BinaryOperator, Expression,
		),
		GuardedExpression,
	)(parser, cursor)
}

// prevents left recursion
func GuardedExpression(parser *Parser, cursor int) (Token, int, error) {
	return MatchOneOf(
		identity,
		MatchAll(AsBracketExpression, char("("), Expression, char(")")),
		MatchAll(AsAssignment, Identifier, char("="), Expression),
		Literal,
		Identifier,
		FunctionDef,
		FunctionCall,
	)(parser, cursor)
}

var BinaryOperator = MatchOneOf(
	identity,
	char("+"),
	char("-"),
	char("*"),
	char("/"),
)

func FunctionCall(parser *Parser, cursor int) (Token, int, error) {
	return MatchAll(
		AsFunctionCall,
		Identifier, char("("), ParamList, char(")"),
	)(parser, cursor)
}

func Literal(parser *Parser, cursor int) (Token, int, error) {
	return MatchOneOf(identity, IntegerLiteral, FloatLiteral)(parser, cursor)
}

func FunctionDef(parser *Parser, cursor int) (Token, int, error) {
	return MatchAll(
		AsFunctionDef,
		KEYWORD_DEF, Identifier,
		char("("), ArgList, char(")"),
		char("{"),
		Block,
		char("}"),
	)(parser, cursor)
}

// TODO FIX THESE REGEXES
var REG_INTEGER_LITERAL = regexp.MustCompile("[0-9]*")
var REG_FLOAT_LITERAL = regexp.MustCompile("[0-9]*\\.[0-9]*")

func IntegerLiteral(parser *Parser, cursor int) (Token, int, error) {
	cursor = parser.findNonWhiteSpace(cursor)
	if cursor == -1 || cursor >= len(parser.text) {
		return nil, cursor, NotMatchError("IntegerLiteral")
	}
	literal := REG_INTEGER_LITERAL.FindString(parser.text[cursor:len(parser.text)])
	if literal == "" {
		return nil, cursor, NotMatchError("IntegerLiteral")
	}
	i, err := strconv.Atoi(literal)
	if err != nil {
		panic("Logic error")
	}
	return IntegerLiteralNode{int64(i)}, cursor + len(literal), nil
}

func FloatLiteral(parser *Parser, cursor int) (Token, int, error) {
	cursor = parser.findNonWhiteSpace(cursor)
	if cursor == -1 || cursor >= len(parser.text) {
		return nil, cursor, NotMatchError("FloatLiteral")
	}
	literal := REG_FLOAT_LITERAL.FindString(parser.text)
	if literal == "" {
		return nil, cursor, NotMatchError("FloatLiteral")
	}
	f, err := strconv.ParseFloat(literal, 64)
	if err != nil {
		panic("Logic error")
	}
	return FloatLiteralNode{f}, cursor + len(literal), nil
}

// matcher aliases
var EmptyFile = EndOfFile
