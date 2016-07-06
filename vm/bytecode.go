package vm

// Stack machine
// Each instruction has 0-2 arguments
// Instructions are fix-sized

const (
	INST_PUSH = iota
	INST_POP
	INST_BINARY
	INST_FUNC_CALL
	INST_FUNC_DEF
)

type Instruction struct {
}
