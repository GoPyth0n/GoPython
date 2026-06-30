package object

type PyArithmeticMethods struct {
	Add  func(PyObject, PyObject) PyObject
	Sub  func(PyObject, PyObject) PyObject
	Mul  func(PyObject, PyObject) PyObject
	Div  func(PyObject, PyObject) PyObject
	IDiv func(PyObject, PyObject) PyObject
	BAnd func(PyObject, PyObject) PyObject
	Pow  func(PyObject, PyObject) PyObject
}

type PyComparisonMethods struct {
	Equals       func(PyObject, PyObject) PyObject
	NotEquals    func(PyObject, PyObject) PyObject
	Less         func(PyObject, PyObject) PyObject
	Greater      func(PyObject, PyObject) PyObject
	LessEqual    func(PyObject, PyObject) PyObject
	GreaterEqual func(PyObject, PyObject) PyObject
}

var defaultCompMethods = &PyComparisonMethods{
	Equals:       Equals,
	NotEquals:    NotEquals,
	Less:         Less,
	Greater:      Greater,
	LessEqual:    LessEqual,
	GreaterEqual: GreaterEqual,
}

type PyType struct {
	Name         string
	ArithMethods *PyArithmeticMethods
	CompMethods  *PyComparisonMethods
}

var (
	IntType *PyType = &PyType{Name: "int", ArithMethods: &PyArithmeticMethods{
		Add:  Add,
		Sub:  Sub,
		Mul:  Mul,
		Div:  Div,
		IDiv: IDiv,
		Pow:  Pow,
		BAnd: BAnd,
	},
		CompMethods: defaultCompMethods}
	FloatType *PyType = &PyType{Name: "float", ArithMethods: &PyArithmeticMethods{
		Add:  Add,
		Sub:  Sub,
		Mul:  Mul,
		Div:  Div,
		IDiv: IDiv,
		Pow:  Pow,
		BAnd: BAnd,
	},
		CompMethods: defaultCompMethods}
	BooleanType  *PyType = &PyType{Name: "bool", CompMethods: defaultCompMethods}
	StrType      *PyType = &PyType{Name: "str"}
	NoneType     *PyType = &PyType{Name: "none"}
	CodeType     *PyType = &PyType{Name: "code"}
	FunctionType *PyType = &PyType{Name: "function"}
)

type PyObject interface {
	Type() *PyType
	String() string
}