package core

import "math/big"

type Bytecode byte
type ConstType int

const (
	LOAD_SMALL_INT Bytecode = iota
	STORE_NAME
	LOAD_NAME
	STORE_GLOBAL

	PUSH_NULL

	COMPARE_OP
	BINARY_OP

	JUMP_IF_FALSE
	JUMP_FORWARD

	RETURN_VALUE

	LOAD_CONST
	CALL

	MAKE_FUNCTION
)

// COMPARE_OP arg codes. CMP_EQ keeps the value (88) already used elsewhere in
// the codebase before this was wired up; the rest are just adjacent values,
// not an attempt to bit-match real CPython's COMPARE_OP oparg packing.
const (
	CMP_EQ = 88
	CMP_NE = 89
	CMP_LT = 90
	CMP_GT = 91
	CMP_LE = 92
	CMP_GE = 93
)

const (
	CONST_INT ConstType = iota
	CONST_FLOAT
	CONST_STRING
	CONST_CODE
	CONST_NONE
)

type Instruction struct {
	Op  Bytecode
	Arg int
}

type PyCode struct {
	Chunk    *Chunk
	ArgCount int
	ArgNames []string
}

type Constant struct {
	Type  ConstType
	Int   big.Int
	Float float64
	Str   string
	Code  *PyCode
}

type Chunk struct {
	Instrs []Instruction
	Names  []string
	Consts []Constant
}