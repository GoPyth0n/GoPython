package object

import "math/big"

func BAnd(a, b PyObject) PyObject {
	ia, ok1 := asInt(a)
	ib, ok2 := asInt(b)

	if ok1 && ok2 {
		z := new(big.Int).And(ia.Value, ib.Value)
		return &PyLongObject{Value: z}
	}

	panic("unsupported & (bitwise AND only works on ints)")
}