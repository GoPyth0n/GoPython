package object

func compareNumbers(a, b PyObject) (float64, float64, bool) {
	if ia, ok := asInt(a); ok {
		if ib, ok := asInt(b); ok {
			return float64(ia.Value.Int64()), float64(ib.Value.Int64()), true
		}
		if fb, ok := asFloat(b); ok {
			return float64(ia.Value.Int64()), fb.Value, true
		}
	}

	if fa, ok := asFloat(a); ok {
		if ib, ok := asInt(b); ok {
			return fa.Value, float64(ib.Value.Int64()), true
		}
		if fb, ok := asFloat(b); ok {
			return fa.Value, fb.Value, true
		}
	}

	return 0, 0, false
}

func Equals(a, b PyObject) PyObject {
	x, y, ok := compareNumbers(a, b)
	if !ok {
		panic("unsupported ==")
	}
	return &PyBoolObject{Value: x == y}
}

func NotEquals(a, b PyObject) PyObject {
	x, y, ok := compareNumbers(a, b)
	if !ok {
		panic("unsupported !=")
	}
	return &PyBoolObject{Value: x != y}
}

func Less(a, b PyObject) PyObject {
	x, y, ok := compareNumbers(a, b)
	if !ok {
		panic("unsupported <")
	}
	return &PyBoolObject{Value: x < y}
}

func Greater(a, b PyObject) PyObject {
	x, y, ok := compareNumbers(a, b)
	if !ok {
		panic("unsupported >")
	}
	return &PyBoolObject{Value: x > y}
}

func LessEqual(a, b PyObject) PyObject {
	x, y, ok := compareNumbers(a, b)
	if !ok {
		panic("unsupported <=")
	}
	return &PyBoolObject{Value: x <= y}
}

func GreaterEqual(a, b PyObject) PyObject {
	x, y, ok := compareNumbers(a, b)
	if !ok {
		panic("unsupported >=")
	}
	return &PyBoolObject{Value: x >= y}
}