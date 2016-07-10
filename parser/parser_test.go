package parser

import (
	"fmt"
	"testing"
)

func testWrapper(tokens []Token) Token {
	if len(tokens) != 2 {
		panic(fmt.Sprintf("Should have 2 tokens: %v", tokens))
	}
	return tokens[0]
}

func pass(t *testing.T, funcName string, tryFunc TryFunc, text string) {
	tryFunc = MatchAll(testWrapper, tryFunc, EndOfFile)
	parser := NewParser(text)
	token, _, err := tryFunc(parser, 0)
	node, ok := token.(Node)
	if ok {
		fmt.Println(PrettyPrint(node.String(), 4))
	} else {
		fmt.Println(token)
	}
	if err != nil {
		t.Errorf("Should not fail: %s(\"%s\") - %s", funcName, text, err)
	}
}

func fail(t *testing.T, funcName string, tryFunc TryFunc, text string) {
	tryFunc = MatchAll(testWrapper, tryFunc, EndOfFile)
	parser := NewParser(text)
	tokens, _, err := tryFunc(parser, 0)
	if err == nil {
		t.Errorf("Should not succeed: %s - %v", funcName, tokens)
	}
}

func TestToken(t *testing.T) {
	pass(t, "EndOfFile", EndOfFile, "")
	fail(t, "EndOfFile", EndOfFile, ";")
	pass(t, "Keyword", keyword("def"), "def")
	fail(t, "Char", char("def"), "fart()")

	pass(t, "MatchOneOf", MatchOneOf(
		identity,
		keyword("def"),
		keyword("compile"),
	),
		"compile",
	)

	pass(t, "ParamList", ParamList, "a, b, c")
	pass(t, "ParamList", ParamList, "var_iable")
	pass(t, "ParamList", ParamList, "")

	pass(t, "ArgList", ArgList, "name:string")
	pass(t, "ArgList", ArgList, "name:string, age:int")
	pass(t, "ArgList", ArgList, "")
	fail(t, "ArgList", ArgList, ",name:string, age:int")
}

func TestNode(t *testing.T) {
	pass(t, "IntegerLiteral", IntegerLiteral, "1230498")
	fail(t, "IntegerLiteral", IntegerLiteral, "XD")

	pass(t, "FloatLiteral", FloatLiteral, "100.00")
	pass(t, "FloatLiteral", FloatLiteral, ".02")

	pass(t, "FunctionDef", FunctionDef, "def myfunc(){}")
	pass(t, "FunctionDef", FunctionDef, "def myfunc(name: hello, hi:there){}")
	pass(t, "FunctionDef", FunctionDef, "def myfunc(name: e){}")
	pass(t, "FunctionDef", FunctionDef, "def myfunc(name: e,){}") // this shouldn't pass btw
	fail(t, "FunctionDef", FunctionDef, "def myfunc(){")
	fail(t, "FunctionDef", FunctionDef, "def myfunc){")
	fail(t, "FunctionDef", FunctionDef, "def myfunc(,name: e){}")

	pass(t, "Expression", Expression, "myfunc(100, 200)")
	pass(t, "Expression", Expression, "myfunc()")
	pass(t, "Expression", Expression, "1.3")
	pass(t, "Expression", Expression, "x = 100")
	fail(t, "Expression", Expression, ";myfunc()")

	pass(t, "Expression", Expression, "name = myfunc(100, 200) + 588 * (x + 2)")
	pass(t, "Expression", Expression, "1 + 1")
	fail(t, "Expression", Expression, "(100 + 200) = myfunc")

	pass(t, "Block", Block, "100+200 300+400 x=400 y=500 * 80 * (90 + z) ")

	pass(t, "Block", Block, `
def main() {
	do_something(x, y)
}

def do_something(x : int, y:int) {
	x = 100 + 200
	y
}
`)
}
