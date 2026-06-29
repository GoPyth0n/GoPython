package runtime

import "fmt"
import (
	"gopython/runtime/object"
)

type PyStack struct {
	data []object.PyObject
}

func (stack *PyStack) Push(obj object.PyObject) {
	stack.data = append(stack.data, obj)
}

func (s *PyStack) Pop() object.PyObject {
	if len(s.data) == 0 {
		panic("stack underflow")
	}

	n := len(s.data)
	obj := s.data[n-1]
	s.data = s.data[:n-1]
	return obj
}
func (s *PyStack) String() string {
	out := "Stack:\n"
	for i := len(s.data) - 1; i >= 0; i-- {
		obj := s.data[i]
		out += fmt.Sprintf("  [%d] %s\n", i, obj.String())
	}
	return out
}