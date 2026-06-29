package object
import "strconv"
type PyLongObject struct {
	Value int
}

func (o *PyLongObject) Type() *PyType {
	return IntType
}

func (o *PyLongObject) String() string {
	return strconv.Itoa(o.Value)
}

func NewInt(v int) PyObject {
	return &PyLongObject{
		Value: v,
	}
}

func intAdd(a, b PyObject) PyObject {
	lhs := a.(*PyLongObject)
	rhs := b.(*PyLongObject)

	return NewInt(lhs.Value + rhs.Value)
}