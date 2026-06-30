package object

type PyArithmeticMethods struct {
	Add  func(PyObject, PyObject) PyObject
	Sub  func(PyObject, PyObject) PyObject
	Mul  func(PyObject, PyObject) PyObject
	Div  func(PyObject, PyObject) PyObject
	IDiv  func(PyObject, PyObject) PyObject
	BAnd func(PyObject, PyObject) PyObject
	Pow  func(PyObject, PyObject) PyObject
}

type PyComparisonMethods struct {
	Equals func(PyObject, PyObject) PyObject
}

type PyType struct {
	Name         string
	ArithMethods *PyArithmeticMethods
	CompMethods *PyComparisonMethods
}

var (
	IntType *PyType = &PyType{Name: "int", ArithMethods: &PyArithmeticMethods{
		Add:  Add,
		Sub:  Sub,
		Mul:  Mul,
		Div: Div,
		IDiv:  IDiv,
		Pow: Pow,
		BAnd: BAnd,
	},
	CompMethods: &PyComparisonMethods{
		Equals: Equals,
	}}
	FloatType    *PyType = &PyType{Name: "float", ArithMethods: &PyArithmeticMethods{
		Add:  Add,
		Sub:  Sub,
		Mul:  Mul,
		Div: Div,
		IDiv:  IDiv,
		Pow: Pow,
		BAnd: BAnd,
	}}
	BooleanType *PyType = &PyType{Name: "bool"}
	StrType      *PyType = &PyType{Name: "str"}
	NoneType     *PyType = &PyType{Name: "none"}
	CodeType     *PyType = &PyType{Name: "code"}
	FunctionType *PyType = &PyType{Name: "function"}
)

type PyObject interface {
	Type() *PyType
	String() string
}
