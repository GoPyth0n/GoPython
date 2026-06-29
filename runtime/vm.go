package runtime

import (
	"gopython/core"
	"gopython/runtime/object"
	"fmt" 
)

type VirtualMachine struct {
	frames []*Frame
	handlers [256]InstrHandler
	Globals map[string]object.PyObject
	NoneSingleton *object.PyNoneObject
	returnValue object.PyObject
}

func NewVM() *VirtualMachine {
	vm := &VirtualMachine{
		Globals: make(map[string]object.PyObject),
		NoneSingleton: &object.PyNoneObject{},
	}

	vm.handlers[core.LOAD_SMALL_INT] = OpLoadSmallInt
	vm.handlers[core.STORE_NAME]     = OpStoreName
	vm.handlers[core.LOAD_NAME]      = OpLoadName
	vm.handlers[core.STORE_GLOBAL]   = OpStoreGlobal
	vm.handlers[core.PUSH_NULL]      = OpPushNull
	vm.handlers[core.BINARY_OP]      = OpBinaryOp
	vm.handlers[core.RETURN_VALUE] = OpReturnValue
	vm.handlers[core.LOAD_CONST] = OpLoadConst
	vm.handlers[core.CALL] = OpCall
	vm.handlers[core.MAKE_FUNCTION] = OpMakeFunction

	return vm
}

type InstrHandler func(vm *VirtualMachine, frame *Frame, instr core.Instruction)
func (vm *VirtualMachine) Run() {
	for len(vm.frames) > 0 {
		frame := vm.frames[len(vm.frames)-1]

		if frame.PC >= len(frame.Chunk.Instrs) {
			// implicit return: RETURN_VALUE already falls back to None
			// when the stack is empty, so don't force one on top of
			// whatever a nested call may have just left there.
			vm.handlers[core.RETURN_VALUE](vm, frame, core.Instruction{})
			continue
		}

		instr := frame.Chunk.Instrs[frame.PC]
		frame.PC++

		vm.handlers[instr.Op](vm, frame, instr)
	}
}

func (vm *VirtualMachine) PushFrame(chunk *core.Chunk) {
	vm.frames = append(vm.frames, &Frame{
		Chunk: chunk,
		PC: 0,
		Stack: &PyStack{},
		Locals: make(map[string]object.PyObject),
	})
}

func (vm *VirtualMachine) Dump() {
	fmt.Println("=== VM DUMP ===")

	fmt.Println("\n--- FRAME STACK ---")
	fmt.Println("Frames:", len(vm.frames))

	for i := len(vm.frames) - 1; i >= 0; i-- {
		f := vm.frames[i]

		fmt.Printf("\n[Frame %d]\n", i)
		fmt.Println("PC:", f.PC)
		if f.ReturnValue != nil {
			fmt.Println("ReturnValue:", f.ReturnValue.String())
		} else {
			fmt.Println("ReturnValue: <not set>")
		}
		fmt.Println("Instructions:", len(f.Chunk.Instrs))

		fmt.Println("Stack:")
		for j := len(f.Stack.data) - 1; j >= 0; j-- {
			fmt.Printf("  [%d] %s\n", j, f.Stack.data[j].String())
		}

		fmt.Println("Locals:")
		for k, v := range f.Locals {
			fmt.Printf("  %s = %s\n", k, v.String())
		}
	}

	fmt.Println("\n--- GLOBALS ---")
	for k, v := range vm.Globals {
		fmt.Printf("  %s = %s\n", k, v.String())
	}

	if vm.returnValue != nil {
		fmt.Println("\n--- LAST RETURN VALUE ---")
		fmt.Println(vm.returnValue.String())
	}
}