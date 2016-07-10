package vm

import "testing"

func TestStack(t *testing.T) {
	s := InterpStack{}
	s.Push(1)
	s.Push(2)
	s.Push(3)

	var val int64
	var values []int64
	var err error

	val, err = s.Pop()
	if err != nil {
		t.Errorf("Expecting no error: %v", err)
	}
	if val != 3 {
		t.Errorf("Expecting 3: %d", val)
	}

	val, err = s.Pop()
	if err != nil {
		t.Errorf("Expecting no error: %v", err)
	}
	if val != 2 {
		t.Errorf("Expecting 2: %d", val)
	}

	val, err = s.Pop()
	if err != nil {
		t.Errorf("Expecting no error: %v", err)
	}
	if val != 1 {
		t.Errorf("Expecting 1: %d", val)
	}

	val, err = s.Pop()
	if err == nil {
		t.Errorf("Expecting error")
	}

	s.Push(4)
	s.Push(5)
	s.Push(6)

	values, err = s.Pops(5)
	if err == nil {
		t.Errorf("Expecting error")
	}

	values, err = s.Pops(3)
	if err != nil {
		t.Errorf("Expecting no error: %v", err)
	}
	if values[0] != 6 {
		t.Errorf("Wrong order")
	}
}

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
	if result != 135-9939 {
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
	if result != 1290889+89324783 {
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
	if result != 101033*(-123873) {
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
	if result != 100/3 {
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
	if val != 200 || err != nil {
		t.Error("Wrong result")
	}
}
