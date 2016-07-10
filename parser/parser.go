package parser

import (
	"fmt"
	"regexp"
	"strconv"
)

// TODO use io.File instead of String
// TODO handle Unicode (haha nice joke)

type Parser struct {
	text string
}

func NewParser(text string) *Parser {
	return &Parser{text: text}
}

func isWhiteSpace(c byte) bool {
	return c == ' ' || c == '\n' || c == '\t'
}

type EOFError int

func (err EOFError) Error() string {
	return fmt.Sprintf("Unexpected EOF")
}

type NotMatchError string

func (err NotMatchError) Error() string {
	return fmt.Sprintf("Could not match type: " + string(err))
}

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

type TryFunc func(parser *Parser, cursor int) (Token, int, error)

func EndOfFile(parser *Parser, cursor int) (Token, int, error) {
	_, newCursor, err := parser.fetch(cursor)
	if _, ok := err.(EOFError); ok {
		return EOFToken{}, newCursor, nil
	}
	return nil, cursor, NotMatchError("EOF")
}

var EmptyFile = EndOfFile

// ConjunctionWrapper convert an array of token into a single token
type ConjunctionWrapper func(tokens []Token) Token

func Conjunction(wrapper ConjunctionWrapper, functions ...TryFunc) TryFunc {
	return func(parser *Parser, cursor int) (Token, int, error) {
		newCursor := cursor
		tokens := []Token{}
		for _, f := range functions {
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

// Disjunction can take multiple paths, which may return different
// token types. The wrapper's job is to cast them all back 1 type
type DisjunctionWrapper func(token Token) Token

func Disjunction(wrapper DisjunctionWrapper, functions ...TryFunc) TryFunc {
	return func(parser *Parser, cursor int) (Token, int, error) {
		var bestToken Token
		bestCursor := -1
		for _, f := range functions {
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
		return nil, cursor, NotMatchError("Disjunction")
	}
}

func TryEmpty(parser *Parser, cursor int) (Token, int, error) {
	return EmptyToken{}, cursor, nil
}

var Empty = TryEmpty

var REG_IDENTIFIER_INITIAL = regexp.MustCompile("[_a-zA-Z]")
var REG_IDENTIFIER = regexp.MustCompile("^[_a-zA-Z0-9]+")

func Identifier(parser *Parser, cursor int) (Token, int, error) {
	first, newCursor, err := parser.fetch(cursor)
	if err != nil || !REG_IDENTIFIER_INITIAL.Match([]byte{first}) {
		return nil, cursor, NotMatchError("Identifier")
	}
	if newCursor == len(parser.text) {
		return IdentifierToken{string(first)}, newCursor, nil
	}
	rest := REG_IDENTIFIER.FindString(parser.text[newCursor:len(parser.text)])
	return IdentifierToken{string(first) + rest}, newCursor + len(rest), nil
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

func Keyword(str string) TryFunc {
	return func(parser *Parser, cursor int) (Token, int, error) {
		newCursor, err := tryString(parser, cursor, str)
		if err != nil {
			return nil, newCursor, NotMatchError("Keyword")
		}
		return KeywordToken{str}, newCursor, nil
	}
}

func Char(str string) TryFunc {
	return func(parser *Parser, cursor int) (Token, int, error) {
		newCursor, err := tryString(parser, cursor, str)
		if err != nil {
			return nil, newCursor, NotMatchError("Char")
		}
		return CharToken{str}, newCursor, nil
	}
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
	return IntegerLiteralToken{int64(i)}, cursor + len(literal), nil
}

func FloatLiteral(parser *Parser, cursor int) (Token, int, error) {
	cursor = parser.findNonWhiteSpace(cursor)
	if cursor == -1 || cursor >= len(parser.text) {
		return nil, cursor, NotMatchError("IntegerLiteral")
	}
	literal := REG_FLOAT_LITERAL.FindString(parser.text)
	if literal == "" {
		return nil, cursor, NotMatchError("FloatLiteral")
	}
	f, err := strconv.ParseFloat(literal, 64)
	if err != nil {
		panic("Logic error")
	}
	return FloatLiteralToken{f}, cursor + len(literal), nil
}

func ArgDeclWrapper(tokens []Token) Token {
	if len(tokens) != 3 {
		panic(fmt.Sprintf("Should have 3 tokens: %v", tokens))
	}
	declName, ok1 := tokens[0].(IdentifierToken)
	declType, ok2 := tokens[2].(IdentifierToken)
	if !ok1 || !ok2 {
		panic("Typecasting failure")
	}
	return ArgDeclToken{declName, declType}
}

var ArgDecl = Conjunction(ArgDeclWrapper, Identifier, Char(":"), Identifier)

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

func ArgListWrapper(tokens []Token) Token {
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

	return Disjunction(
		Token2ArgListToken,
		Conjunction(ArgListWrapper, ArgDecl, Char(","), ArgList),
		ArgDecl,
		Empty,
	)(parser, cursor)
}

func FunctionDefWrapper(tokens []Token) Token {
	if len(tokens) != 8 {
		panic(fmt.Sprintf("Should have 8 tokens: %v", tokens))
	}
	name, ok1 := tokens[1].(IdentifierToken)
	arglist, ok2 := tokens[3].(ArgListToken)
	block, ok3 := tokens[6].(BlockToken)
	if !ok1 || !ok2 || !ok3 {
		fmt.Println(ok1, ok2, ok3)
		panic("Typecasting failure")
	}
	return FunctionDefToken{name, arglist, block}
}

func FunctionDef(parser *Parser, cursor int) (Token, int, error) {
	return Conjunction(
		FunctionDefWrapper,
		Keyword("function"), Identifier,
		Char("("), ArgList, Char(")"),
		Char("{"),
		TryBlock,
		Char("}"),
	)(parser, cursor)
}

func Token2ParamListToken(token Token) Token {
	list := []Token{}
	switch paramToken := token.(type) {
	default:
		// probably should wrap Expression instead of doing this
		list = append(list, paramToken)
	case ParamListToken:
		list = append(list, paramToken.ParamList...)
	case EmptyToken:
		// do nothing
	}
	return ParamListToken{list}
}

func ParamListWrapper(tokens []Token) Token {
	if len(tokens) != 3 {
		panic(fmt.Sprintf("Should have 3 tokens: %v", tokens))
	}
	head := tokens[0]
	tail, ok := Token2ParamListToken(tokens[2]).(ParamListToken)
	if !ok {
		panic("Typecasting failure")
	}
	newList := []Token{head}
	newList = append(newList, tail.ParamList...)
	return ParamListToken{newList}
}

func ParamList(parser *Parser, cursor int) (Token, int, error) {
	return Disjunction(
		Token2ParamListToken,
		Conjunction(
			ParamListWrapper,
			Expression, Char(","), ParamList,
		),
		Expression,
		Empty,
	)(parser, cursor)
}

func FunctionCallWrapper(tokens []Token) Token {
	if len(tokens) != 4 {
		panic(fmt.Sprintf("Should have 4 tokens: %v", tokens))
	}
	name, ok := tokens[0].(IdentifierToken)
	paramList := tokens[2].(ParamListToken)
	if !ok {
		panic("Typecasting failure")
	}
	return FunctionCallToken{name, paramList}
}

func FunctionCall(parser *Parser, cursor int) (Token, int, error) {
	return Conjunction(
		FunctionCallWrapper,
		Identifier, Char("("), ParamList, Char(")"),
	)(parser, cursor)
}

func identity(token Token) Token {
	return token
}

func Literal(parser *Parser, cursor int) (Token, int, error) {
	return Disjunction(identity, IntegerLiteral, FloatLiteral)(parser, cursor)
}

var BinaryOperator = Disjunction(
	identity,
	Char("+"),
	Char("-"),
	Char("*"),
	Char("/"),
)

func BinaryOperatorWrapper(tokens []Token) Token {
	if len(tokens) != 3 {
		panic(fmt.Sprintf("Should have 3 tokens: %v", tokens))
	}
	operator, ok := tokens[1].(CharToken)

	if !ok {
		panic("Typecasting failure")
	}
	return BinaryOperatorToken{tokens[0], operator, tokens[2]}
}

func BracketExpressionWrapper(tokens []Token) Token {
	if len(tokens) != 3 {
		panic(fmt.Sprintf("Should have 3 tokens: %v", tokens))
	}
	return tokens[1]
}

func Expression(parser *Parser, cursor int) (Token, int, error) {
	return Disjunction(
		identity,
		Conjunction(
			BinaryOperatorWrapper,
			GuardedExpression, BinaryOperator, Expression,
		),
		GuardedExpression,
	)(parser, cursor)
}

func AssignmentWrapper(tokens []Token) Token {
	if len(tokens) != 3 {
		panic(fmt.Sprintf("Should have 3 tokens: %v", tokens))
	}
	id, ok := tokens[0].(IdentifierToken)
	if !ok {
		panic("Typecasting failure")
	}

	return AssignmentToken{id, tokens[2]}
}

// prevents left recursion
func GuardedExpression(parser *Parser, cursor int) (Token, int, error) {
	return Disjunction(
		identity,
		Conjunction(BracketExpressionWrapper, Char("("), Expression, Char(")")),
		Conjunction(AssignmentWrapper, Identifier, Char("="), Expression),
		Literal,
		Identifier,
		FunctionDef,
		FunctionCall,
	)(parser, cursor)
}

func BlockWrapper(tokens []Token) Token {
	if len(tokens) != 2 {
		panic(fmt.Sprintf("Should have 2 tokens: %v", tokens))
	}
	tail, ok := tokens[2].(BlockToken)
	if ok {
		panic("Typecasting failure")
	}
	newList := []Token{}
	newList = append(newList, tokens[0])
	for _, token := range tail.ExprList {
		newList = append(newList, token)
	}
	return BlockToken{newList}
}

func Token2BlockToken(token Token) Token {
	list := []Token{}
	switch blockToken := token.(type) {
	default:
		list = append(list, blockToken)
	case BlockToken:
		list = append(list, blockToken.ExprList...)
	case EmptyToken:
		// do nothing
	}
	return BlockToken{list}
}

func TryBlock(parser *Parser, cursor int) (Token, int, error) {
	return Disjunction(
		Token2BlockToken,
		Conjunction(BlockWrapper, Expression, TryBlock),
		Expression,
		Empty,
	)(parser, cursor)
}
