package parser

import . "github.com/trungaczne/gimmick/vm"

/* --- VM bytecode generation routines ---*/

func (node IntegerLiteralNode) CodeGen(builder CodeBuilder) {
	builder.Push(
		Instruction{INST_PUSH, node.Value, ARG_NOOP},
	)
}

func (node FloatLiteralNode) CodeGen(builder CodeBuilder) {
	builder.Push(PushInst(int64(node.Value)))
}

func (node IdentifierNode) CodeGen(builder CodeBuilder) {
	builder.Push(PushInst(builder.Resolve(node.Name)))
}

func (node FunctionDefNode) CodeGen(builder CodeBuilder) {
	builder.DefineFunc(node.ArgList, func(scopedBuilder CodeBuilder) {
		node.Block.CodeGen(scopedBuilder)
	})
}

func (node FunctionCallNode) CodeGen(builder CodeBuilder) {
	for _, arg := range node.ParamList {
		// IMPLICATION: arguments are processed from left to right
		arg.CodeGen(builder)
	}
	builder.Push(InvokeInst(builder.Resolve(node.Name)))
}

func (node BinaryOperatorNode) CodeGen(builder CodeBuilder) {
	node.Left.CodeGen(builder)
	node.Right.CodeGen(builder)
	builder.Push(BinaryInst(node.Operator))
}

func (node AssignmentNode) CodeGen(builder CodeBuilder) {
	node.Expr.CodeGen(builder)
	id := builder.ResolveOrDefine(node.Dest)
	builder.Push(AssignInst(id))
}

func (node BlockNode) CodeGen(builder CodeBuilder) {
	for i, expr := range node.ExprList {
		expr.CodeGen(builder)
		if i != len(node.ExprList)-1 {
			// only the last expression should push to the stack
			builder.Push(PopInst())
		}
	}
}

func (node ModuleNode) CodeGen(builder CodeBuilder) {
	node.Block.CodeGen(builder)
}
