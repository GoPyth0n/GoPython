package object

type PyNumberMethods struct {
	Add func(PyObject, PyObject) PyObject
	Sub func(PyObject, PyObject) PyObject
	Mul func(PyObject, PyObject) PyObject
	Div func(PyObject, PyObject) PyObject
	Pow func(PyObject, PyObject) PyObject
}

type PyType struct {
	Name string
	Number *PyNumberMethods
}

var (
	IntType *PyType = &PyType{Name: "int", Number: &PyNumberMethods{
		Add: intAdd,
	}}
	FloatType *PyType = &PyType{Name: "float"}
	StrType *PyType = &PyType{Name: "str"}
	NoneType *PyType = &PyType{Name: "none"}
	CodeType *PyType = &PyType{Name: "code"}
	FunctionType *PyType = &PyType{Name: "function"}
)

type PyObject interface {
	Type() *PyType
	String() string
}