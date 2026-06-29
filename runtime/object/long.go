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
	if c := GetCachedInt(v); c != nil {
		return c
	}
	return &PyLongObject{Value: v}
}

func asInt(obj PyObject) (*PyLongObject, bool) {
	val, ok := obj.(*PyLongObject)
	return val, ok
}

func intAdd(a, b PyObject) PyObject {
    lhs, ok1 := asInt(a)
    rhs, ok2 := asInt(b)

	if !ok1 || !ok2 {
		panic("unsupported operand")
	}

    return NewInt(lhs.Value + rhs.Value)
}

func intSub(a, b PyObject) PyObject {
    lhs, ok1 := asInt(a)
    rhs, ok2 := asInt(b)

	if !ok1 || !ok2 {
		panic("unsupported operand")
	}

    return NewInt(lhs.Value - rhs.Value)
}

func intMul(a, b PyObject) PyObject {
    lhs, ok1 := asInt(a)
    rhs, ok2 := asInt(b)

	if !ok1 || !ok2 {
		panic("unsupported operand")
	}

    return NewInt(lhs.Value * rhs.Value)
}

func intDiv(a, b PyObject) PyObject {
    lhs, ok1 := asInt(a)
    rhs, ok2 := asInt(b)

	if !ok1 || !ok2 {
		panic("unsupported operand")
	}

    return NewInt(lhs.Value / rhs.Value)
}

func intbAnd(a, b PyObject) PyObject {
    lhs, ok1 := asInt(a)
    rhs, ok2 := asInt(b)

	if !ok1 || !ok2 {
		panic("unsupported operand")
	}

    return NewInt(lhs.Value & rhs.Value)
}