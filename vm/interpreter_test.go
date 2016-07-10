package vm

import "testing"

func TestBinaryInst(t *testing.T) {
	interp := NewInterpreter()

	// subtraction
	f := []Instruction{
		PushInst(135),
		PushInst(9939),
		BinaryInst("-"),
	}

	id := interp.AddFunc(f)
	err := interp.ExecFunc(id)
	if err != nil {
		t.Error(err)
	}

	result, err := interp.Stack.Pop()
	r := result.(int64)
	if r != 135-9939 {
		t.Error("Wrong result")
	}

	// addition
	f = []Instruction{
		PushInst(1290889),
		PushInst(89324783),
		BinaryInst("+"),
	}

	id = interp.AddFunc(f)
	err = interp.ExecFunc(id)
	if err != nil {
		t.Error(err)
	}

	result, err = interp.Stack.Pop()
	r = result.(int64)
	if r != 1290889+89324783 {
		t.Error("Wrong result")
	}

	// multiply
	f = []Instruction{
		PushInst(101033),
		PushInst(-123873),
		BinaryInst("*"),
	}

	id = interp.AddFunc(f)
	err = interp.ExecFunc(id)
	if err != nil {
		t.Error(err)
	}

	result, err = interp.Stack.Pop()
	r = result.(int64)
	if r != 101033*(-123873) {
		t.Error("Wrong result")
	}

	// divison
	f = []Instruction{
		PushInst(100),
		PushInst(3),
		BinaryInst("/"),
	}

	id = interp.AddFunc(f)
	err = interp.ExecFunc(id)
	if err != nil {
		t.Error(err)
	}

	result, err = interp.Stack.Pop()
	r = result.(int64)
	if r != 100/3 {
		t.Error("Wrong result")
	}

	f = []Instruction{
		PushInst(100),
		PushInst(0),
		BinaryInst("/"),
	}

	id = interp.AddFunc(f)
	err = interp.ExecFunc(id)
	if err == nil {
		t.Error("Expecting error")
	}
}

func TestInvokeInst(t *testing.T) {
	interp := NewInterpreter()

	childFunc := []Instruction{
		PushInst(100),
		PushInst(100),
		BinaryInst("+"),
	}

	childID := interp.AddFunc(childFunc)

	parentFunc := []Instruction{
		InvokeInst(childID),
	}

	parentID := interp.AddFunc(parentFunc)

	err := interp.ExecFunc(parentID)
	if err != nil {
		t.Error(err)
	}

	val, err := interp.Stack.Pop()
	v := val.(int64)
	if v != 200 || err != nil {
		t.Error("Wrong result")
	}
}
