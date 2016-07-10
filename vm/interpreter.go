package vm

import "fmt"

type Function struct {
	Inst []Instruction
}

type CallStack struct {
	FuncID int64
	PC     int64
}

type GimmickInterpreter struct {
	// Code ...
	Func      []*Function
	CallStack []*CallStack

	/// ... and data
	Stack InterpStack
	Heap  []int64
}

type InterpStack struct {
	Value []int64
}

func (stack *InterpStack) Push(value int64) {
	stack.Value = append(stack.Value, value)
}

func (stack *InterpStack) Pop() (int64, error) {
	if len(stack.Value) == 0 {
		return -1, fmt.Errorf("Not enough values to pop from stack")
	}
	l := len(stack.Value) - 1
	val := stack.Value[l]
	stack.Value = stack.Value[0:l]
	return val, nil
}

func (stack *InterpStack) Pops(num int64) ([]int64, error) {
	l := int64(len(stack.Value) - 1)
	if l+1 < num {
		return []int64{}, fmt.Errorf("Not enough values to pop from stack")
	}
	val := make([]int64, num)
	for i := int64(0); i < num; i++ {
		val[i] = stack.Value[l-i]
	}
	stack.Value = stack.Value[0 : l-num+1]
	return val, nil
}

func NewInterpreter() *GimmickInterpreter {
	// zero values are sufficient
	return &GimmickInterpreter{}
}

// Name is stripped by the CodeBuilder, there's only ID
func (interp *GimmickInterpreter) AddFunc(instructions []Instruction) int64 {
	newFunc := &Function{instructions}
	interp.Func = append(interp.Func, newFunc)
	id := int64(len(interp.Func) - 1)
	return id
}

func (interp *GimmickInterpreter) ExecFunc(id int64) error {
	if id < 0 || id >= int64(len(interp.Func)) {
		return fmt.Errorf("Invalid function ID to execute")
	}

	callstack := &CallStack{id, 0}
	interp.CallStack = append(interp.CallStack, callstack)
	return interp.Start()
}

func (interp *GimmickInterpreter) LastCallStack() *CallStack {
	return interp.CallStack[len(interp.CallStack)-1]
}

func (interp *GimmickInterpreter) Start() error {
	for {
		if len(interp.CallStack) == 0 {
			// Done!
			return nil
		}

		curStack := interp.LastCallStack()
		if curStack.FuncID < 0 || curStack.FuncID >= int64(len(interp.Func)) {
			return fmt.Errorf("FuncID out of bound: %v", curStack.FuncID)
		}
		curFunc := interp.Func[curStack.FuncID]

		if curStack.PC < 0 {
			return fmt.Errorf("PC out of bound: %v", curStack.PC)
		}

		if curStack.PC >= int64(len(curFunc.Inst)) {
			// nothing else to execute in this stack, return
			// TODO make an explicit instruction for returning
			interp.CallStack = interp.CallStack[0 : len(interp.CallStack)-1]
			continue
		}

		inst := curFunc.Inst[curStack.PC]
		curStack.PC += 1

		// Yay!
		err := interp.Exec(inst)
		if err != nil {
			return fmt.Errorf("Inst %v failed: %v", inst, err)
		}
	}
}

func (interp *GimmickInterpreter) Exec(inst Instruction) error {
	switch inst.Type {
	case INST_PUSH:
		interp.Stack.Push(inst.Arg1)
		return nil
	case INST_POP:
		_, err := interp.Stack.Pop()
		return err
	case INST_ASSIGN:
		return interp.ExecAssign(inst)
	case INST_BINARY:
		return interp.ExecBinary(inst)
	case INST_INVOKE:
		return interp.ExecInvoke(inst)
	}
	return nil
}

func (interp *GimmickInterpreter) ExecBinary(inst Instruction) error {
	// for A <op> B, we expect the interpreter to push A then B
	// thus B will be popped first, followed by A
	vals, err := interp.Stack.Pops(2)
	if err != nil {
		return err
	}
	left := vals[1]
	right := vals[0]
	// only do int64 math for now, the rest will require some bytecode changes
	switch inst.Arg1 {
	case ARG_OP_ADD:
		interp.Stack.Push(left + right)
		return nil
	case ARG_OP_SUB:
		interp.Stack.Push(left - right)
		return nil
	case ARG_OP_MUL:
		interp.Stack.Push(left * right)
		return nil
	case ARG_OP_DIV:
		if right == 0 {
			return fmt.Errorf("Division by zero")
		}
		interp.Stack.Push(left / right)
		return nil
	default:
		return fmt.Errorf("Bad bytecode")
	}
}

func (interp *GimmickInterpreter) ExecInvoke(inst Instruction) error {
	callstack := CallStack{inst.Arg1, 0}
	interp.CallStack = append(interp.CallStack, &callstack)
	return nil
}

func (interp *GimmickInterpreter) ExecAssign(inst Instruction) error {
	return nil
}
