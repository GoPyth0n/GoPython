package object

func Equals(a PyObject, b PyObject) PyObject {
	if ia, ok1 := asInt(a); ok1 {
		if ib, ok2 := asInt(b); ok2 {
			if ib.Value.Sign() == 0 {
				panic("ZeroDivisionError: integer division by zero")
			}

			z := ia.Value.Cmp(ib.Value)
			r := z == 0


			return &PyBoolObject{Value: r}
		}
	}

	panic("not supported")
}