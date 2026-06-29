package core

type Bytecode int
type ConstType int

const (
	LOAD_SMALL_INT Bytecode = iota
	STORE_NAME
	LOAD_NAME
	STORE_GLOBAL

	PUSH_NULL

	BINARY_OP

	RETURN_VALUE

	LOAD_CONST
	CALL

	MAKE_FUNCTION
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
	Chunk *Chunk
	ArgCount int
    ArgNames []string
}

type Constant struct {
	Type ConstType
	Int int
	Float float64
	Str string
	Code *PyCode
}


type Chunk struct {
	Instrs []Instruction
	Names  []string
	Consts []Constant
}