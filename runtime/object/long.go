package object

import "math/big"
type PyLongObject struct {
	Value *big.Int
}

func (o *PyLongObject) Type() *PyType {
	return IntType
}

func (o *PyLongObject) String() string {
	return o.Value.String()
}

func NewInt(v int) PyObject {
	if c := GetCachedInt(v); c != nil {
		return c
	}
	return &PyLongObject{Value: big.NewInt(int64(v))}
}