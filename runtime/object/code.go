package object

import "gopython/core"

type PyCodeObject struct {
	Chunk *core.Chunk
	ArgCount int
    ArgNames []string
}

func (o *PyCodeObject) Type() *PyType {
	return CodeType
}

func (o *PyCodeObject) String() string {
	return "code object"
}
