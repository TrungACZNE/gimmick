package utils

import "testing"

func TestStack(t *testing.T) {
	s := Stack{}
	s.Push(1)
	s.Push(2)
	s.Push(3)

	var val interface{}
	var values []interface{}
	var err error

	val, err = s.Pop()
	vi := val.(int)
	if err != nil {
		t.Errorf("Expecting no error: %v", err)
	}
	if vi != 3 {
		t.Errorf("Expecting 3: %d", vi)
	}

	val, err = s.Pop()
	vi = val.(int)
	if err != nil {
		t.Errorf("Expecting no error: %v", err)
	}
	if vi != 2 {
		t.Errorf("Expecting 2: %d", vi)
	}

	val, err = s.Pop()
	vi = val.(int)
	if err != nil {
		t.Errorf("Expecting no error: %v", err)
	}
	if vi != 1 {
		t.Errorf("Expecting 1: %d", vi)
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
