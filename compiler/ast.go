package compiler

import (
	"gopython/core"
)

type Node interface {
	node()
}

type Stmt interface {
	Node
	stmtNode()
}

type Module struct {
	Body []Stmt
}

type AssignStmt struct {
	Target string
	Value  Expr
}

type IfStmt struct {
	Cond Expr
	Then []Stmt
	Else []Stmt
}

type ReturnStmt struct {
	Value Expr
}

type ExprStmt struct {
	Value Expr
}

type FunctionDef struct {
	Name   string
	Params []string
	Body   []Stmt
}

func (*Module) node()      {}
func (*AssignStmt) node()  {}
func (*ReturnStmt) node()  {}
func (*ExprStmt) node()    {}
func (*IfStmt) node()    {}
func (*FunctionDef) node() {}

func (*AssignStmt) stmtNode()  {}
func (*ReturnStmt) stmtNode()  {}
func (*ExprStmt) stmtNode()    {}
func (*IfStmt) stmtNode()    {}
func (*FunctionDef) stmtNode() {}

type Expr interface {
	Node
	exprNode()
}

type Literal struct {
	core.Constant
}

type NameExpr struct {
	Name string
}

type BinOp struct {
	Op    string
	Left  Expr
	Right Expr
}

type CallExpr struct {
	Func Expr
	Args []Expr
}

func (*Literal) node()  {}
func (*NameExpr) node() {}
func (*BinOp) node()    {}
func (*CallExpr) node() {}

func (*Literal) exprNode()  {}
func (*NameExpr) exprNode() {}
func (*BinOp) exprNode()    {}
func (*CallExpr) exprNode() {}

func BinOpArg(op string) int {
	switch op {
	case "+":
		return 0
	case "-":
		return 10
	case "*":
		return 5
	case "/":
		return 11
	case "//":
		return 2
	case "&":
		return 1
	case "**":
		return 8
	}
	panic("compiler: unknown operator " + op)
}

func isCompareOp(op string) bool {
	switch op {
	case "==", "!=", "<", ">", "<=", ">=":
		return true
	}
	return false
}

func CompareOpArg(op string) int {
	switch op {
	case "==":
		return core.CMP_EQ
	case "!=":
		return core.CMP_NE
	case "<":
		return core.CMP_LT
	case ">":
		return core.CMP_GT
	case "<=":
		return core.CMP_LE
	case ">=":
		return core.CMP_GE
	}
	panic("compiler: unknown comparison operator " + op)
}