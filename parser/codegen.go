package parser

import . "github.com/trungaczne/gimmick/vm"

/* --- VM bytecode generation routines ---*/

// every expression *must* push to the stack

func (token IntegerLiteralNode) CodeGen(builder CodeBuilder) {
	builder.Push(
		Instruction{INST_PUSH, token.Value, ARG_NOOP},
	)
}

func (token FloatLiteralNode) CodeGen(builder CodeBuilder) {
	builder.Push(
		Instruction{INST_PUSH, int64(token.Value), ARG_NOOP},
	)
}

func (token IdentifierNode) CodeGen(builder CodeBuilder) {
	builder.Push(
		Instruction{INST_PUSH, builder.Resolve(token.Name), ARG_NOOP},
	)
}

func (token FunctionDefNode) CodeGen(builder CodeBuilder) {
	builder.DefineFunc(token.ArgList, func(scopedBuilder CodeBuilder) {
		token.Block.CodeGen(scopedBuilder)
	})
}

func (token FunctionCallNode) CodeGen(builder CodeBuilder) {
	for _, arg := range token.ParamList {
		// IMPLICATION: arguments are processed from left to right
		arg.CodeGen(builder)
	}
	builder.Push(
		Instruction{INST_INVOKE, builder.Resolve(token.Name), ARG_NOOP},
	)
}

func (token BinaryOperatorNode) CodeGen(builder CodeBuilder) {
	token.Left.CodeGen(builder)
	token.Right.CodeGen(builder)
	var op int64
	switch token.Operator {
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

func (token AssignmentNode) CodeGen(builder CodeBuilder) {
	id := builder.ResolveOrDefine(token.Dest)
	builder.Push(
		Instruction{INST_ASSIGN, id, ARG_NOOP},
	)
}

func (token BlockNode) CodeGen(builder CodeBuilder) {
	for i, expr := range token.ExprList {
		expr.CodeGen(builder)
		if i != len(token.ExprList)-1 {
			// only the last expression should push to the stack
			builder.Push(
				Instruction{INST_POP, ARG_NOOP, ARG_NOOP},
			)
		}
	}
}
