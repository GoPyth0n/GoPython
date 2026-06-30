package runtime

import (
	"fmt"
	"gopython/core"
	"gopython/runtime/object"
	"math/big"
)

func OpLoadSmallInt(vm *VirtualMachine, frame *Frame, instr core.Instruction) {
	if instr.Arg >= -5 && instr.Arg <= 256 {
		frame.Stack.Push(object.IntCache[instr.Arg+5])
		return
	}
	frame.Stack.Push(&object.PyLongObject{Value: big.NewInt(int64(instr.Arg))})
}

func OpLoadConst(vm *VirtualMachine, frame *Frame, instr core.Instruction) {
	c := frame.Chunk.Consts[instr.Arg]

	switch c.Type {
	case core.CONST_INT:
		frame.Stack.Push(&object.PyLongObject{Value: &c.Int})

	case core.CONST_FLOAT:
		frame.Stack.Push(&object.PyFloatObject{Value: c.Float})

	case core.CONST_CODE:
		frame.Stack.Push(&object.PyCodeObject{
			Chunk:    c.Code.Chunk,
			ArgCount: c.Code.ArgCount,
			ArgNames: c.Code.ArgNames,
		})
	}
}

func OpStoreName(vm *VirtualMachine, frame *Frame, instr core.Instruction) {
	val := frame.Stack.Pop()
	vm.Globals[frame.Chunk.Names[instr.Arg]] = val
}

func OpStoreGlobal(vm *VirtualMachine, frame *Frame, instr core.Instruction) {
	val := frame.Stack.Pop()
	vm.Globals[frame.Chunk.Names[instr.Arg]] = val
}

func OpLoadName(vm *VirtualMachine, frame *Frame, instr core.Instruction) {
	name := frame.Chunk.Names[instr.Arg]

	if val, ok := frame.Locals[name]; ok {
		frame.Stack.Push(val)
		return
	}

	if val, ok := vm.Globals[name]; ok {
		frame.Stack.Push(val)
		return
	}

	panic("NameError: " + name)
}

func OpCompareOp(vm *VirtualMachine, frame *Frame, instr core.Instruction) {
	rhs := frame.Stack.Pop()
	lhs := frame.Stack.Pop()

	switch instr.Arg {
	case core.CMP_EQ:
		frame.Stack.Push(lhs.Type().CompMethods.Equals(lhs, rhs))
	case core.CMP_NE:
		frame.Stack.Push(lhs.Type().CompMethods.NotEquals(lhs, rhs))
	case core.CMP_LT:
		frame.Stack.Push(lhs.Type().CompMethods.Less(lhs, rhs))
	case core.CMP_GT:
		frame.Stack.Push(lhs.Type().CompMethods.Greater(lhs, rhs))
	case core.CMP_LE:
		frame.Stack.Push(lhs.Type().CompMethods.LessEqual(lhs, rhs))
	case core.CMP_GE:
		frame.Stack.Push(lhs.Type().CompMethods.GreaterEqual(lhs, rhs))
	}
}

// isTruthy is a quick type-switch rather than a proper per-type method table
// (unlike ArithMethods/CompMethods). Fine for now; worth promoting to a
// PyType.Truthy func if more types need custom truthiness rules later.
func isTruthy(obj object.PyObject) bool {
	switch v := obj.(type) {
	case *object.PyBoolObject:
		return v.Value
	case *object.PyLongObject:
		return v.Value.Sign() != 0
	case *object.PyFloatObject:
		return v.Value != 0
	case *object.PyNoneObject:
		return false
	default:
		return true
	}
}

func OpJumpIfFalse(vm *VirtualMachine, frame *Frame, instr core.Instruction) {
	val := frame.Stack.Pop()
	if !isTruthy(val) {
		frame.PC = instr.Arg
	}
}

func OpJumpForward(vm *VirtualMachine, frame *Frame, instr core.Instruction) {
	frame.PC = instr.Arg
}

func OpBinaryOp(vm *VirtualMachine, frame *Frame, instr core.Instruction) {
	b := frame.Stack.Pop()
	a := frame.Stack.Pop()

	switch instr.Arg {
	case 0:
		frame.Stack.Push(a.Type().ArithMethods.Add(a, b))
	case 1:
		frame.Stack.Push(a.Type().ArithMethods.BAnd(a, b))
	case 2:
		frame.Stack.Push(a.Type().ArithMethods.IDiv(a, b))
	case 5:
		frame.Stack.Push(a.Type().ArithMethods.Mul(a, b))
	case 8:
		frame.Stack.Push(a.Type().ArithMethods.Pow(a, b))
	case 10:
		frame.Stack.Push(a.Type().ArithMethods.Sub(a, b))
	case 11:
		frame.Stack.Push(a.Type().ArithMethods.Div(a, b))
	}
}

func OpPushNull(vm *VirtualMachine, frame *Frame, instr core.Instruction) {
	frame.Stack.Push(vm.NoneSingleton)
}

func OpReturnValue(vm *VirtualMachine, frame *Frame, instr core.Instruction) {
	var value object.PyObject

	// 1. try stack
	if len(frame.Stack.data) > 0 {
		value = frame.Stack.Pop()
	} else {
		value = frame.ReturnValue // fallback
	}

	if value == nil {
		value = vm.NoneSingleton
	}

	frame.ReturnValue = value

	// pop frame
	vm.frames = vm.frames[:len(vm.frames)-1]

	// return to caller
	if len(vm.frames) > 0 {
		caller := vm.frames[len(vm.frames)-1]
		caller.Stack.Push(value)
	} else {
		vm.returnValue = value
	}
}

func OpMakeFunction(vm *VirtualMachine, frame *Frame, instr core.Instruction) {
	codeObj := frame.Stack.Pop()

	pyCode := codeObj.(*object.PyCodeObject)

	fn := &object.PyFunctionObject{
		Name:    "<lambda>",
		Code:    *pyCode,
		Globals: vm.Globals,
	}

	frame.Stack.Push(fn)
}

func OpCall(vm *VirtualMachine, frame *Frame, instr core.Instruction) {
	argCount := instr.Arg

	args := make([]object.PyObject, argCount)

	for i := argCount - 1; i >= 0; i-- {
		args[i] = frame.Stack.Pop()
	}

	fnObj := frame.Stack.Pop()

	pyfn, ok := fnObj.(*object.PyFunctionObject)
	if !ok {
		panic("CALL expects function object")
	}

	if argCount != pyfn.Code.ArgCount {
		panic(fmt.Sprintf(
			"expected %d arguments, got %d",
			pyfn.Code.ArgCount,
			argCount,
		))
	}

	newFrame := &Frame{
		Chunk:  pyfn.Code.Chunk,
		PC:     0,
		Stack:  &PyStack{},
		Locals: make(map[string]object.PyObject),
	}

	for i, name := range pyfn.Code.ArgNames {
		newFrame.Locals[name] = args[i]
	}

	vm.frames = append(vm.frames, newFrame)
}