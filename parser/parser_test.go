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
	tryFunc = Conjunction(testWrapper, tryFunc, EndOfFile)
	parser := NewParser(text)
	tokens, _, err := tryFunc(parser, 0)
	fmt.Println(tokens)
	if err != nil {
		t.Errorf("Should not fail: %s(\"%s\") - %s", funcName, text, err)
	}
}

func fail(t *testing.T, funcName string, tryFunc TryFunc, text string) {
	tryFunc = Conjunction(testWrapper, tryFunc, EndOfFile)
	parser := NewParser(text)
	tokens, _, err := tryFunc(parser, 0)
	if err == nil {
		t.Errorf("Should not succeed: %s - %v", funcName, tokens)
	}
}

func TestGeneral(t *testing.T) {
	pass(t, "EndOfFile", EndOfFile, "")
	fail(t, "EndOfFile", EndOfFile, ";")

	pass(t, "IntegerLiteral", IntegerLiteral, "1230498")
	fail(t, "IntegerLiteral", IntegerLiteral, "XD")

	pass(t, "FloatLiteral", FloatLiteral, "100.00")
	pass(t, "FloatLiteral", FloatLiteral, ".02")

	pass(t, "Keyword", Keyword("function"), "function")
	fail(t, "Char", Char("function"), "fart()")

	pass(t, "Disjunction", Disjunction(
		identity,
		Keyword("function"),
		Keyword("compile"),
	),
		"compile",
	)

	pass(t, "ArgList", ArgList, "name:string")
	pass(t, "ArgList", ArgList, "name:string, age:int")
	pass(t, "ArgList", ArgList, "")
	fail(t, "ArgList", ArgList, ",name:string, age:int")

	pass(t, "FunctionDef", FunctionDef, "function myfunc(){}")
	pass(t, "FunctionDef", FunctionDef, "function myfunc(name: hello, hi:there){}")
	pass(t, "FunctionDef", FunctionDef, "function myfunc(name: e){}")
	pass(t, "FunctionDef", FunctionDef, "function myfunc(name: e,){}") // this shouldn't pass btw
	fail(t, "FunctionDef", FunctionDef, "function myfunc(){")
	fail(t, "FunctionDef", FunctionDef, "function myfunc){")
	fail(t, "FunctionDef", FunctionDef, "function myfunc(,name: e){}")

	pass(t, "ParamList", ParamList, "a, b, c")
	pass(t, "ParamList", ParamList, "var_iable")
	pass(t, "ParamList", ParamList, "")

	pass(t, "Expression", Expression, "myfunc(100, 200)")
	pass(t, "Expression", Expression, "myfunc()")
	pass(t, "Expression", Expression, "1.3")
	pass(t, "Expression", Expression, "x = 100")
	fail(t, "Expression", Expression, ";myfunc()")

	pass(t, "Expression", Expression, "name = myfunc(100, 200) + 588 * (x + 2)")
	pass(t, "Expression", Expression, "1 + 1")
	fail(t, "Expression", Expression, "(100 + 200) = myfunc")
}
