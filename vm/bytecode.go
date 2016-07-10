package vm

// Stack machine
// Each instruction has 0-2 arguments
// Instructions are fix-sized (24 bytes)

// Comments on the instructions go to the bottom of the file

const (
	INST_PUSH int64 = iota
	INST_POP
	INST_BINARY
	INST_INVOKE
	INST_ASSIGN
)

const ARG_NOOP int64 = 0xFFFFFFFF
const (
	ARG_OP_ADD int64 = iota
	ARG_OP_SUB
	ARG_OP_MUL
	ARG_OP_DIV
	ARG_OP_ASSIGN
)

// Base type
type Instruction struct {
	Type int64
	Arg1 int64
	Arg2 int64
}

/* --- Instruction declarations --- */

type PushInstruction struct {
	Instruction
}

type PopInstruction struct {
	Instruction
}

type BinaryInstruction struct {
	Instruction
}

type InvokeInstruction struct {
	Instruction
}

type AssignInstruction struct {
	Instruction
}

// Put value ontop of stack. Value could be anything castable to int64
// StackSize +1
func PushInst(value int64) Instruction {
	return Instruction{INST_PUSH, value, ARG_NOOP}
}

// Remove topmost value from stack
// StackSize -1
func PopInst() Instruction {
	return Instruction{INST_POP, ARG_NOOP, ARG_NOOP}
}

// Pops 2 values from the stack, compute the operation, then push the value back
// StackSize -1
func BinaryInst(op string) Instruction {
	switch op {
	case "+":
		return Instruction{INST_BINARY, ARG_OP_ADD, ARG_NOOP}
	case "-":
		return Instruction{INST_BINARY, ARG_OP_SUB, ARG_NOOP}
	case "*":
		return Instruction{INST_BINARY, ARG_OP_MUL, ARG_NOOP}
	case "/":
		return Instruction{INST_BINARY, ARG_OP_DIV, ARG_NOOP}
	}
	panic("Don't let this happen")
}

// Invoke the function with the given ID
// StackSize: -(number of arguments of function)
func InvokeInst(id int64) Instruction {
	return Instruction{INST_INVOKE, id, ARG_NOOP}
}

// Pops the topmost value from the stack and assigns it to the variable ID
// StackSize: -1
func AssignInst(id int64) Instruction {
	return Instruction{INST_ASSIGN, id, ARG_NOOP}
}
