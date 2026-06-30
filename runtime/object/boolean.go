package object

type PyBoolObject struct {
	Value bool
}

func (o *PyBoolObject) Type() *PyType {
	return BooleanType
}

func (o *PyBoolObject) String() string {
	if o.Value == true {
		return "True"
	}
	return "False"
}