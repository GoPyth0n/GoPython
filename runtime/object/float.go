package object

import (
	"strconv"
)

type PyFloatObject struct {
	Value float64
}

func (o *PyFloatObject) Type() *PyType {
	return FloatType
}

func (o *PyFloatObject) String() string {
	return strconv.FormatFloat(o.Value, 'g', -1, 64)
}

func NewFloat(v float64) PyObject {
	return &PyFloatObject{Value: v}
}
