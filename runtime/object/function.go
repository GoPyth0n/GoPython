package object

type PyFunctionObject struct {
	Name string
	Globals map[string]PyObject

	Code PyCodeObject
}
 
func (o *PyFunctionObject) Type() *PyType {
	return FunctionType
}

func (o *PyFunctionObject) String() string {
	return o.Name
}
