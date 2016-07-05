package parser

import (
	"fmt"
	"regexp"
	"strconv"
)

// TODO use io.File instead of String
// TODO handle Unicode

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

type TryFunc func(parser *Parser, cursor int) (Node, int, error)

func EndOfFile(parser *Parser, cursor int) (Node, int, error) {
	_, newCursor, err := parser.fetch(cursor)
	if _, ok := err.(EOFError); ok {
		return EOFNode{}, newCursor, nil
	}
	return nil, cursor, NotMatchError("EOF")
}

var EmptyFile = EndOfFile

// ConjunctionWrapper convert an array of node into a single node
type ConjunctionWrapper func(nodes []Node) Node

func Conjunction(wrapper ConjunctionWrapper, functions ...TryFunc) TryFunc {
	return func(parser *Parser, cursor int) (Node, int, error) {
		newCursor := cursor
		nodes := []Node{}
		for _, f := range functions {
			token, _cursor, err := f(parser, newCursor)
			if err != nil {
				return nil, cursor, err
			}
			newCursor = _cursor
			nodes = append(nodes, token)
		}
		return wrapper(nodes), newCursor, nil
	}
}

// Disjunction can take multiple paths, which may return different
// node types. The wrapper's job is to cast them all back 1 type
type DisjunctionWrapper func(node Node) Node

func Disjunction(wrapper DisjunctionWrapper, functions ...TryFunc) TryFunc {
	return func(parser *Parser, cursor int) (Node, int, error) {
		var bestNode Node
		bestCursor := -1
		for _, f := range functions {
			token, newCursor, err := f(parser, cursor)
			if err == nil {
				if newCursor > bestCursor {
					bestCursor = newCursor
					bestNode = token
				}
			}
		}
		if bestCursor > -1 {
			return wrapper(bestNode), bestCursor, nil
		}
		return nil, cursor, NotMatchError("Disjunction")
	}
}

func TryEmpty(parser *Parser, cursor int) (Node, int, error) {
	return EmptyNode{}, cursor, nil
}

var Empty = TryEmpty

var REG_IDENTIFIER_INITIAL = regexp.MustCompile("[_a-zA-Z]")
var REG_IDENTIFIER = regexp.MustCompile("^[_a-zA-Z0-9]+")

func Identifier(parser *Parser, cursor int) (Node, int, error) {
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

// shared by Keyword and Symbol
func tryString(p *Parser, cursor int, str string) (int, error) {
	cursor = p.findNonWhiteSpace(cursor)
	l := len(p.text)
	if cursor == -1 || cursor+len(str) > l || p.text[cursor:cursor+len(str)] != str {
		return cursor, fmt.Errorf("Can't match")
	}
	return cursor + len(str), nil
}

func Keyword(str string) TryFunc {
	return func(parser *Parser, cursor int) (Node, int, error) {
		newCursor, err := tryString(parser, cursor, str)
		if err != nil {
			return nil, newCursor, NotMatchError("Keyword")
		}
		return KeywordNode{str}, newCursor, nil
	}
}

func Symbol(str string) TryFunc {
	return func(parser *Parser, cursor int) (Node, int, error) {
		newCursor, err := tryString(parser, cursor, str)
		if err != nil {
			return nil, newCursor, NotMatchError("Symbol")
		}
		return SymbolNode{str}, newCursor, nil
	}
}

// TODO FIX THESE REGEXES
var REG_INTEGER_LITERAL = regexp.MustCompile("[0-9]*")
var REG_FLOAT_LITERAL = regexp.MustCompile("[0-9]*\\.[0-9]*")

func IntegerLiteral(parser *Parser, cursor int) (Node, int, error) {
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

func FloatLiteral(parser *Parser, cursor int) (Node, int, error) {
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
	return FloatLiteralNode{f}, cursor + len(literal), nil
}

func ArgDeclWrapper(nodes []Node) Node {
	if len(nodes) != 3 {
		panic(fmt.Sprintf("Should have 3 nodes: %v", nodes))
	}
	declName, ok1 := nodes[0].(IdentifierNode)
	declType, ok2 := nodes[2].(IdentifierNode)
	if !ok1 || !ok2 {
		panic("Typecasting failure")
	}
	return ArgDeclNode{declName, declType}
}

var ArgDecl = Conjunction(ArgDeclWrapper, Identifier, Symbol(":"), Identifier)

func Node2ArgListNode(node Node) Node {
	list := []ArgDeclNode{}
	switch arglistNode := node.(type) {
	default:
		list = nil
	case ArgListNode:
		list = append(list, arglistNode.ArgDecl...)
	case ArgDeclNode:
		list = append(list, arglistNode)
	case EmptyNode:
		// do nothing
	}
	return ArgListNode{list}
}

func ArgListWrapper(nodes []Node) Node {
	if len(nodes) != 3 {
		panic(fmt.Sprintf("Should have 3 nodes: %v", nodes))
	}
	head, ok := nodes[0].(ArgDeclNode)

	if !ok {
		panic("Typecasting failure")
	}

	newList := []ArgDeclNode{head}

	node, ok := Node2ArgListNode(nodes[2]).(ArgListNode)
	if !ok {
		panic("Programming error")
	}

	tail := node.ArgDecl
	if tail == nil {
		panic("Typecasting failure")
	}
	newList = append(newList, tail...)
	return ArgListNode{newList}
}

func ArgList(parser *Parser, cursor int) (Node, int, error) {
	// TODO this accepts (arg,) which shouldn't be the case

	return Disjunction(
		Node2ArgListNode,
		Conjunction(ArgListWrapper, ArgDecl, Symbol(","), ArgList),
		ArgDecl,
		Empty,
	)(parser, cursor)
}

func FunctionDefWrapper(nodes []Node) Node {
	if len(nodes) != 8 {
		panic(fmt.Sprintf("Should have 8 nodes: %v", nodes))
	}
	name, ok1 := nodes[1].(IdentifierNode)
	arglist, ok2 := nodes[3].(ArgListNode)
	block, ok3 := nodes[6].(BlockNode)
	if !ok1 || !ok2 || !ok3 {
		fmt.Println(ok1, ok2, ok3)
		panic("Typecasting failure")
	}
	return FunctionDefNode{name, arglist, block}
}

func FunctionDef(parser *Parser, cursor int) (Node, int, error) {
	return Conjunction(
		FunctionDefWrapper,
		Keyword("function"), Identifier,
		Symbol("("), ArgList, Symbol(")"),
		Symbol("{"),
		TryBlock,
		Symbol("}"),
	)(parser, cursor)
}

func Node2ParamListNode(node Node) Node {
	list := []Node{}
	switch paramNode := node.(type) {
	default:
		// probably should wrap Expression instead of doing this
		list = append(list, paramNode)
	case ParamListNode:
		list = append(list, paramNode.ParamList...)
	case EmptyNode:
		// do nothing
	}
	return ParamListNode{list}
}

func ParamListWrapper(nodes []Node) Node {
	if len(nodes) != 3 {
		panic(fmt.Sprintf("Should have 3 nodes: %v", nodes))
	}
	head := nodes[0]
	tail, ok := Node2ParamListNode(nodes[2]).(ParamListNode)
	if !ok {
		panic("Typecasting failure")
	}
	newList := []Node{head}
	newList = append(newList, tail.ParamList...)
	return ParamListNode{newList}
}

func ParamList(parser *Parser, cursor int) (Node, int, error) {
	return Disjunction(
		Node2ParamListNode,
		Conjunction(
			ParamListWrapper,
			Expression, Symbol(","), ParamList,
		),
		Expression,
		Empty,
	)(parser, cursor)
}

func FunctionCallWrapper(nodes []Node) Node {
	if len(nodes) != 4 {
		panic(fmt.Sprintf("Should have 4 nodes: %v", nodes))
	}
	name, ok := nodes[0].(IdentifierNode)
	paramList := nodes[2].(ParamListNode)
	if !ok {
		panic("Typecasting failure")
	}
	return FunctionCallNode{name, paramList}
}

func FunctionCall(parser *Parser, cursor int) (Node, int, error) {
	return Conjunction(
		FunctionCallWrapper,
		Identifier, Symbol("("), ParamList, Symbol(")"),
	)(parser, cursor)
}

func identity(node Node) Node {
	return node
}

func Literal(parser *Parser, cursor int) (Node, int, error) {
	return Disjunction(identity, IntegerLiteral, FloatLiteral)(parser, cursor)
}

var BinaryOperator = Disjunction(
	identity,
	Symbol("+"),
	Symbol("-"),
	Symbol("*"),
	Symbol("/"),
	Symbol("="),
)

func BinaryOperatorWrapper(nodes []Node) Node {
	if len(nodes) != 3 {
		panic(fmt.Sprintf("Should have 3 nodes: %v", nodes))
	}
	operator, ok := nodes[1].(SymbolNode)

	if !ok {
		panic("Typecasting failure")
	}
	return BinaryOperatorNode{nodes[0], operator, nodes[2]}
}

func ExpressionWrapper(nodes []Node) Node {
	if len(nodes) != 3 {
		panic(fmt.Sprintf("Should have 3 nodes: %v", nodes))
	}
	return nodes[1]
}

func Expression(parser *Parser, cursor int) (Node, int, error) {
	return Disjunction(
		identity,
		Conjunction(
			BinaryOperatorWrapper,
			GuardedExpression, BinaryOperator, Expression,
		),
		GuardedExpression,
	)(parser, cursor)
}

// prevents left recursion
func GuardedExpression(parser *Parser, cursor int) (Node, int, error) {
	return Disjunction(
		identity,
		Conjunction(ExpressionWrapper, Symbol("("), Expression, Symbol(")")),
		Literal,
		Identifier,
		FunctionDef,
		FunctionCall,
	)(parser, cursor)
}

func BlockWrapper(nodes []Node) Node {
	if len(nodes) != 2 {
		panic(fmt.Sprintf("Should have 2 nodes: %v", nodes))
	}
	tail, ok := nodes[2].(BlockNode)
	if ok {
		panic("Typecasting failure")
	}
	newList := []Node{}
	newList = append(newList, nodes[0])
	for _, token := range tail.ExprList {
		newList = append(newList, token)
	}
	return BlockNode{newList}
}

func Node2BlockNode(node Node) Node {
	list := []Node{}
	switch blockNode := node.(type) {
	default:
		list = append(list, blockNode)
	case BlockNode:
		list = append(list, blockNode.ExprList...)
	case EmptyNode:
		// do nothing
	}
	return BlockNode{list}
}

func TryBlock(parser *Parser, cursor int) (Node, int, error) {
	return Disjunction(
		Node2BlockNode,
		Conjunction(BlockWrapper, Expression, TryBlock),
		Expression,
		Empty,
	)(parser, cursor)
}
