package runtime

func toInt(o *PyObject) int {
	if o.Type == INT {
		return o.IntVal
	}
	return int(o.IntVal)
}

func toFloat(o *PyObject) float64 {
	if o.Type == FLOAT {
		return o.FloatVal
	}
	return float64(o.IntVal)
}