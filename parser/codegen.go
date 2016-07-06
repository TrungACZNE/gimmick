package parser

import . "github.com/trungaczne/gimmick/vm"

/* --- VM bytecode generation routines ---*/

// every expression *must* push to the stack

func (node EOFNode) CodeGen(builder CodeBuilder) {
}

func (node IntegerLiteralNode) CodeGen(builder CodeBuilder) {
	builder.Push(
		Instruction{INST_PUSH, node.Value, ARG_NOOP},
	)
}

func (node FloatLiteralNode) CodeGen(builder CodeBuilder) {
	builder.Push(
		Instruction{INST_PUSH, int64(node.Value), ARG_NOOP},
	)
}

func (node KeywordNode) CodeGen(builder CodeBuilder) {
	// TODO fix this
	panic("This shouldn't be reached")
}

func (node CharNode) CodeGen(builder CodeBuilder) {
	panic("This shouldn't be reached")
}

func (node IdentifierNode) CodeGen(builder CodeBuilder) {
	builder.Push(
		Instruction{INST_PUSH, builder.Resolve(node.Name), ARG_NOOP},
	)
}

func (node ArgDeclNode) CodeGen(builder CodeBuilder) {
	panic("This shouldn't be reached")
}

func (node ArgListNode) CodeGen(builder CodeBuilder) {
	panic("This shouldn't be reached")
}

func (node ParamListNode) CodeGen(builder CodeBuilder) {
	panic("This shouldn't be reached")
}

func (node EmptyNode) CodeGen(builder CodeBuilder) {
}

func (node FunctionDefNode) CodeGen(builder CodeBuilder) {
	signature := []NameType{}
	for _, arg := range node.ArgList.ArgDecl {
		signature = append(signature, NameType{
			Name: arg.NameNode.Name,
			Type: arg.TypeNode.Name,
		})

	}
	builder.DefineFunc(signature, func(scopedBuilder CodeBuilder) {
		node.Block.CodeGen(scopedBuilder)
	})
}

func (node FunctionCallNode) CodeGen(builder CodeBuilder) {
	for _, arg := range node.ParamList.ParamList {
		// IMPLICATION: arguments are processed from left to right
		arg.CodeGen(builder)
	}
	builder.Push(
		Instruction{INST_INVOKE, builder.Resolve(node.Name.Name), ARG_NOOP},
	)
}

func (node BinaryOperatorNode) CodeGen(builder CodeBuilder) {
	node.Left.CodeGen(builder)
	node.Right.CodeGen(builder)
	var op int64
	switch node.Operator.Name {
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

func (node AssignmentNode) CodeGen(builder CodeBuilder) {
	id := builder.ResolveOrDefine(node.Dest.Name)
	builder.Push(
		Instruction{INST_ASSIGN, id, ARG_NOOP},
	)
}

func (node BlockNode) CodeGen(builder CodeBuilder) {
	for i, child := range node.ExprList {
		child.CodeGen(builder)
		if i != len(node.ExprList)-1 {
			// only the last expression should push to the stack
			builder.Push(
				Instruction{INST_POP, ARG_NOOP, ARG_NOOP},
			)
		}
	}
}
