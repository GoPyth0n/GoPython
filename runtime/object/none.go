package object

type PyNoneObject struct {

}

func (o *PyNoneObject) Type() *PyType {
	return NoneType
}

func (o *PyNoneObject) String() string {
	return "None"
}