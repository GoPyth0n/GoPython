package object

import (
	"math"
	"math/big"
)

func asInt(v PyObject) (*PyLongObject, bool) {
	x, ok := v.(*PyLongObject)
	return x, ok
}

func asFloat(v PyObject) (*PyFloatObject, bool) {
	x, ok := v.(*PyFloatObject)
	return x, ok
}

func Add(a, b PyObject) PyObject {
	if fa, ok := asFloat(a); ok {
		if fb, ok := asFloat(b); ok {
			return &PyFloatObject{Value: fa.Value + fb.Value}
		}
		if ib, ok := asInt(b); ok {
			return &PyFloatObject{Value: fa.Value + float64(ib.Value.Int64())}
		}
	}

	if ia, ok := asInt(a); ok {
		if fb, ok := asFloat(b); ok {
			return &PyFloatObject{Value: float64(ia.Value.Int64()) + fb.Value}
		}
		if ib, ok := asInt(b); ok {
			z := new(big.Int).Add(ia.Value, ib.Value)
			return &PyLongObject{Value: z}
		}
	}

	panic("unsupported +")
}

func Sub(a, b PyObject) PyObject {
	if fa, ok := asFloat(a); ok {
		if fb, ok := asFloat(b); ok {
			return &PyFloatObject{Value: fa.Value - fb.Value}
		}
		if ib, ok := asInt(b); ok {
			return &PyFloatObject{Value: fa.Value - float64(ib.Value.Int64())}
		}
	}

	if ia, ok := asInt(a); ok {
		if fb, ok := asFloat(b); ok {
			return &PyFloatObject{Value: float64(ia.Value.Int64()) - fb.Value}
		}
		if ib, ok := asInt(b); ok {
			z := new(big.Int).Sub(ia.Value, ib.Value)
			return &PyLongObject{Value: z}
		}
	}

	panic("unsupported -")
}

func Mul(a, b PyObject) PyObject {
	if fa, ok := asFloat(a); ok {
		if fb, ok := asFloat(b); ok {
			return &PyFloatObject{Value: fa.Value * fb.Value}
		}
		if ib, ok := asInt(b); ok {
			return &PyFloatObject{Value: fa.Value * float64(ib.Value.Int64())}
		}
	}

	if ia, ok := asInt(a); ok {
		if fb, ok := asFloat(b); ok {
			return &PyFloatObject{Value: float64(ia.Value.Int64()) * fb.Value}
		}
		if ib, ok := asInt(b); ok {
			z := new(big.Int).Mul(ia.Value, ib.Value)
			return &PyLongObject{Value: z}
		}
	}

	panic("unsupported *")
}

func Div(a, b PyObject) PyObject {
	if fa, ok := asFloat(a); ok {
		if fb, ok := asFloat(b); ok {
			return &PyFloatObject{Value: fa.Value / fb.Value}
		}
		if ib, ok := asInt(b); ok {
			return &PyFloatObject{Value: fa.Value / float64(ib.Value.Int64())}
		}
	}

	if ia, ok := asInt(a); ok {
		if fb, ok := asFloat(b); ok {
			return &PyFloatObject{Value: float64(ia.Value.Int64()) / fb.Value}
		}
		if ib, ok := asInt(b); ok {
			return &PyFloatObject{
				Value: float64(ia.Value.Int64()) / float64(ib.Value.Int64()),
			}
		}
	}

	panic("unsupported /")
}

func IDiv(a, b PyObject) PyObject {
	// int // int
	if ia, ok1 := asInt(a); ok1 {
		if ib, ok2 := asInt(b); ok2 {
			if ib.Value.Sign() == 0 {
				panic("ZeroDivisionError: integer division by zero")
			}

			z := new(big.Int).Div(ia.Value, ib.Value)
			return &PyLongObject{Value: z}
		}
	}

	// float // float
	if fa, ok1 := asFloat(a); ok1 {
		if fb, ok2 := asFloat(b); ok2 {
			if fb.Value == 0 {
				panic("ZeroDivisionError: float division by zero")
			}

			return &PyFloatObject{
				Value: math.Floor(fa.Value / fb.Value),
			}
		}
	}

	// int // float
	if ia, ok1 := asInt(a); ok1 {
		if fb, ok2 := asFloat(b); ok2 {
			if fb.Value == 0 {
				panic("ZeroDivisionError: float division by zero")
			}

			return &PyFloatObject{
				Value: math.Floor(float64(ia.Value.Int64()) / fb.Value),
			}
		}
	}

	// float // int
	if fa, ok1 := asFloat(a); ok1 {
		if ib, ok2 := asInt(b); ok2 {
			if ib.Value.Sign() == 0 {
				panic("ZeroDivisionError: integer division by zero")
			}

			return &PyFloatObject{
				Value: math.Floor(fa.Value / float64(ib.Value.Int64())),
			}
		}
	}

	panic("unsupported operands for //")
}

func Pow(a, b PyObject) PyObject {
	if fa, ok := asFloat(a); ok {
		if fb, ok := asFloat(b); ok {
			return &PyFloatObject{Value: math.Pow(fa.Value, fb.Value)}
		}
		if ib, ok := asInt(b); ok {
			return &PyFloatObject{Value: math.Pow(fa.Value, float64(ib.Value.Int64()))}
		}
	}

	if ia, ok := asInt(a); ok {
		if fb, ok := asFloat(b); ok {
			return &PyFloatObject{Value: math.Pow(float64(ia.Value.Int64()), fb.Value)}
		}
		if ib, ok := asInt(b); ok {
			return &PyLongObject{
				Value: new(big.Int).Exp(ia.Value, ib.Value, nil),
			}
		}
	}

	panic("unsupported **")
}
