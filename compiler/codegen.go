package compiler

import (
	"fmt"
	"gopython/core"
)

type CodeGen struct {
	chunk *core.Chunk
}

func NewCodeGen() *CodeGen {
	return &CodeGen{
		chunk: &core.Chunk{},
	}
}

func Compile(node Node) *core.Chunk {
	cg := NewCodeGen()
	cg.visit(node)
	// Implicitly return at the end of the module
	cg.emit(core.PUSH_NULL, 0)
	cg.emit(core.RETURN_VALUE, 0)
	return cg.chunk
}

func (cg *CodeGen) emit(op core.Bytecode, arg int) {
	cg.chunk.Instrs = append(cg.chunk.Instrs, core.Instruction{Op: op, Arg: arg})
}

func (cg *CodeGen) addConstant(c core.Constant) int {
	cg.chunk.Consts = append(cg.chunk.Consts, c)
	return len(cg.chunk.Consts) - 1
}

func (cg *CodeGen) addName(name string) int {
	cg.chunk.Names = append(cg.chunk.Names, name)
	return len(cg.chunk.Names) - 1
}

func (cg *CodeGen) visit(n Node) {
	switch node := n.(type) {
	case *Module:
		for _, stmt := range node.Body {
			cg.visit(stmt)
		}

	case *AssignStmt:
		cg.visit(node.Value)
		nameIdx := cg.addName(node.Target)
		cg.emit(core.STORE_NAME, nameIdx)

	case *ExprStmt:
		cg.visit(node.Value)

	case *ReturnStmt:
		if node.Value != nil {
			cg.visit(node.Value)
		} else {
			cg.emit(core.PUSH_NULL, 0)
		}
		cg.emit(core.RETURN_VALUE, 0)

	case *IfStmt:
		cg.visit(node.Cond)
		jumpFalse := len(cg.chunk.Instrs)
		cg.emit(core.JUMP_IF_FALSE, 0) // placeholder
		for _, stmt := range node.Then {
			cg.visit(stmt)
		}
		if node.Else != nil {
			jumpEnd := len(cg.chunk.Instrs)
			cg.emit(core.JUMP_FORWARD, 0) // placeholder
			cg.chunk.Instrs[jumpFalse].Arg = len(cg.chunk.Instrs) // else starts here
			for _, stmt := range node.Else {
				cg.visit(stmt)
			}
			cg.chunk.Instrs[jumpEnd].Arg = len(cg.chunk.Instrs)
		} else {
			cg.chunk.Instrs[jumpFalse].Arg = len(cg.chunk.Instrs)
		}

	case *FunctionDef:
		funcCg := NewCodeGen()
		for _, arg := range node.Params {
			funcCg.addName(arg)
		}
		for _, stmt := range node.Body {
			funcCg.visit(stmt)
		}
		funcCg.emit(core.PUSH_NULL, 0)
		funcCg.emit(core.RETURN_VALUE, 0)

		pyCode := &core.PyCode{
			Chunk:    funcCg.chunk,
			ArgNames: append([]string(nil), node.Params...),
			ArgCount: len(node.Params),
		}
		funcConst := core.Constant{Type: core.CONST_CODE, Code: pyCode}
		constIdx := cg.addConstant(funcConst)

		cg.emit(core.LOAD_CONST, constIdx)
		cg.emit(core.MAKE_FUNCTION, 0)

		nameIdx := cg.addName(node.Name)
		cg.emit(core.STORE_NAME, nameIdx)

	case *Literal:
		// Convert node Constant to  core.Constant
		var c core.Constant
		switch node.Constant.Type {
		case core.CONST_INT:
			c = core.Constant{Type: core.CONST_INT, Int: node.Constant.Int}
		case core.CONST_FLOAT:
			c = core.Constant{Type: core.CONST_FLOAT, Float: node.Constant.Float}
		case core.CONST_STRING:
			c = core.Constant{Type: core.CONST_STRING, Str: node.Constant.Str}
		default:
			c = core.Constant{Type: core.CONST_NONE}
		}
		idx := cg.addConstant(c)
		cg.emit(core.LOAD_CONST, idx)

	case *NameExpr:
		idx := cg.addName(node.Name)
		cg.emit(core.LOAD_NAME, idx)

	case *BinOp:
		cg.visit(node.Left)
		cg.visit(node.Right)
		if isCompareOp(node.Op) {
			cg.emit(core.COMPARE_OP, CompareOpArg(node.Op))
		} else {
			cg.emit(core.BINARY_OP, BinOpArg(node.Op))
		}

	case *CallExpr:
		cg.visit(node.Func)
		for _, arg := range node.Args {
			cg.visit(arg)
		}
		cg.emit(core.CALL, len(node.Args))

	default:
		panic(fmt.Sprintf("codegen: unhandled node type %T", n))
	}
}