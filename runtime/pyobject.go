package runtime

import (
	"fmt"
	"gopython/core"
	"math"
)

type PyObjectType int

const (
	INT PyObjectType = iota
	FLOAT
	COMPLEX

	STRING
	LIST

	CODE
	FUNCTION
	NONE
)

type PyObject struct {
	Type PyObjectType

	IntVal int
	FloatVal float64
	StrVal string
	List []*PyObject

	Func *PyFunction
	Code *core.PyCode
}

type PyNumber interface {
	Add(*PyObject) *PyObject
	Sub(*PyObject) *PyObject
	Mul(*PyObject) *PyObject
	Div(*PyObject) *PyObject
	IDiv(*PyObject) *PyObject
}

func (a *PyObject) Add(b *PyObject) *PyObject {
	if a.Type == FLOAT || b.Type == FLOAT {
		return &PyObject{
			Type: FLOAT,
			FloatVal: toFloat(a) + toFloat(b),
		}
	}
	return &PyObject{
		Type: INT,
		IntVal: a.IntVal + b.IntVal,
	}
}

func (a *PyObject) Sub(b *PyObject) *PyObject {
	if a.Type == FLOAT || b.Type == FLOAT {
		return &PyObject{
			Type: FLOAT,
			FloatVal: toFloat(a) - toFloat(b),
		}
	}
	return &PyObject{
		Type: INT,
		IntVal: a.IntVal - b.IntVal,
	}
}

func (a *PyObject) Mul(b *PyObject) *PyObject {
	if a.Type == FLOAT || b.Type == FLOAT {
		return &PyObject{
			Type: FLOAT,
			FloatVal: toFloat(a) * toFloat(b),
		}
	}
	return &PyObject{
		Type: INT,
		IntVal: a.IntVal * b.IntVal,
	}
}

func (a *PyObject) Div(b *PyObject) *PyObject {
	return &PyObject{
		Type: FLOAT,
		FloatVal: toFloat(a) / toFloat(b),
	}
}

func (a *PyObject) Pow(b *PyObject) *PyObject {
	return &PyObject{
		Type: FLOAT,
		FloatVal: math.Pow(toFloat(a), toFloat(b)),
	}
}

func (a *PyObject) IDiv(b *PyObject) *PyObject {
	return &PyObject{
		Type: INT,
		IntVal: int(math.Floor(toFloat(a) / toFloat(b))),
	}
}

type PyFunction struct {
    Name        string
    Code        *core.PyCode

    Globals     map[string]*PyObject
    Builtins    map[string]*PyObject

    Defaults    []*PyObject
    KwDefaults  map[string]*PyObject

    Closure     []*PyObject

    Doc         string
    Module      string
    Annotations map[string]*PyObject
}

func (o *PyObject) String() string {
	switch o.Type {
	case INT:
		return fmt.Sprintf("INT(%d)", o.IntVal)
	case FLOAT:
		return fmt.Sprintf("FLOAT(%f)", o.FloatVal)
	case STRING:
		return fmt.Sprintf("STRING(%s)", o.StrVal)
	case NONE:
		return "NONE"
	case FUNCTION:
		return "FUNCTION"
	default:
		return "UNKNOWN"
	}
}

func NewInt(v int) *PyObject {
	return &PyObject{
		Type: INT,
		IntVal: v,
	}
}

func NewFloat(f float64) *PyObject {
	return &PyObject{
		Type: FLOAT,
		FloatVal: f,
	}
}

func NewString(s string) *PyObject {
	return &PyObject{
		Type: FLOAT,
		StrVal: s,
	}
}


func NewFunction(fn *PyFunction) *PyObject {
	return &PyObject{
		Type: FUNCTION,
		Func: fn,
	}
}