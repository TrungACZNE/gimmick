package utils

import "fmt"

type Stack struct {
	Value []interface{}
}

func (stack *Stack) Push(value interface{}) {
	stack.Value = append(stack.Value, value)
}

func (stack *Stack) Pop() (interface{}, error) {
	if len(stack.Value) == 0 {
		return -1, fmt.Errorf("Not enough values to pop from stack")
	}
	l := len(stack.Value) - 1
	val := stack.Value[l]
	stack.Value = stack.Value[0:l]
	return val, nil
}

func (stack *Stack) Pops(num int64) ([]interface{}, error) {
	l := int64(len(stack.Value) - 1)
	if l+1 < num {
		return []interface{}{}, fmt.Errorf("Not enough values to pop from stack")
	}
	val := make([]interface{}, num)
	for i := int64(0); i < num; i++ {
		val[i] = stack.Value[l-i]
	}
	stack.Value = stack.Value[0 : l-num+1]
	return val, nil
}
