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
func (*FunctionDef) node() {}

func (*AssignStmt) stmtNode()  {}
func (*ReturnStmt) stmtNode()  {}
func (*ExprStmt) stmtNode()    {}
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
