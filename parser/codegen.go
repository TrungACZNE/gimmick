package parser

import . "github.com/trungaczne/gimmick/vm"

/* --- VM bytecode generation routines ---*/

// every expression *must* push to the stack

func (token EOFToken) CodeGen(builder CodeBuilder) {
}

func (token IntegerLiteralToken) CodeGen(builder CodeBuilder) {
	builder.Push(
		Instruction{INST_PUSH, token.Value, ARG_NOOP},
	)
}

func (token FloatLiteralToken) CodeGen(builder CodeBuilder) {
	builder.Push(
		Instruction{INST_PUSH, int64(token.Value), ARG_NOOP},
	)
}

func (token KeywordToken) CodeGen(builder CodeBuilder) {
	// TODO fix this
	panic("This shouldn't be reached")
}

func (token CharToken) CodeGen(builder CodeBuilder) {
	panic("This shouldn't be reached")
}

func (token IdentifierToken) CodeGen(builder CodeBuilder) {
	builder.Push(
		Instruction{INST_PUSH, builder.Resolve(token.Name), ARG_NOOP},
	)
}

func (token ArgDeclToken) CodeGen(builder CodeBuilder) {
	panic("This shouldn't be reached")
}

func (token ArgListToken) CodeGen(builder CodeBuilder) {
	panic("This shouldn't be reached")
}

func (token ParamListToken) CodeGen(builder CodeBuilder) {
	panic("This shouldn't be reached")
}

func (token EmptyToken) CodeGen(builder CodeBuilder) {
}

func (token FunctionDefToken) CodeGen(builder CodeBuilder) {
	signature := []NameType{}
	for _, arg := range token.ArgList.ArgDecl {
		signature = append(signature, NameType{
			Name: arg.NameToken.Name,
			Type: arg.TypeToken.Name,
		})

	}
	builder.DefineFunc(signature, func(scopedBuilder CodeBuilder) {
		token.Block.CodeGen(scopedBuilder)
	})
}

func (token FunctionCallToken) CodeGen(builder CodeBuilder) {
	for _, arg := range token.ParamList.ParamList {
		// IMPLICATION: arguments are processed from left to right
		arg.CodeGen(builder)
	}
	builder.Push(
		Instruction{INST_INVOKE, builder.Resolve(token.Name.Name), ARG_NOOP},
	)
}

func (token BinaryOperatorToken) CodeGen(builder CodeBuilder) {
	token.Left.CodeGen(builder)
	token.Right.CodeGen(builder)
	var op int64
	switch token.Operator.Name {
	case "+":
		op = ARG_OP_ADD
	case "-":
		op = ARG_OP_SUB
	case "*":
		op = ARG_OP_MUL
	case "/":
		op = ARG_OP_DIV
	default:
		panic("This shouldn't happen")
	}
	builder.Push(Instruction{INST_BINARY, op, ARG_NOOP})
}

func (token AssignmentToken) CodeGen(builder CodeBuilder) {
	id := builder.ResolveOrDefine(token.Dest.Name)
	builder.Push(
		Instruction{INST_ASSIGN, id, ARG_NOOP},
	)
}

func (token BlockToken) CodeGen(builder CodeBuilder) {
	for i, child := range token.ExprList {
		child.CodeGen(builder)
		if i != len(token.ExprList)-1 {
			// only the last expression should push to the stack
			builder.Push(
				Instruction{INST_POP, ARG_NOOP, ARG_NOOP},
			)
		}
	}
}
