package runtime
import (
	"gopython/core"
	"gopython/runtime/object"
)
type Frame struct {
	Chunk *core.Chunk
	PC int
	Stack *PyStack
	ReturnValue object.PyObject
	Locals map[string]object.PyObject
}