// parser.go
package compiler

import (
	"fmt"
	"strconv"
	"strings"

	"gopython/core"
)

type Parser struct {
	tokens []Token
	pos    int
}

func NewParser(tokens []Token) *Parser {
	return &Parser{tokens: tokens}
}

func Parse(src string) *Module {
	tokens := NewLexer(src).Tokenize()
	return NewParser(tokens).ParseModule()
}

func (p *Parser) cur() Token { return p.tokens[p.pos] }

func (p *Parser) check(t TokenType) bool { return p.cur().Type == t }

func (p *Parser) advance() Token {
	t := p.tokens[p.pos]
	if p.pos < len(p.tokens)-1 {
		p.pos++
	}
	return t
}

func (p *Parser) expect(t TokenType, what string) Token {
	if !p.check(t) {
		panic(fmt.Sprintf("SyntaxError: line %d: expected %s, got %q", p.cur().Line, what, p.cur().Lit))
	}
	return p.advance()
}

func (p *Parser) skipNewlines() {
	for p.check(TOKEN_NEWLINE) {
		p.advance()
	}
}

func (p *Parser) peekIsAssign() bool {
	return p.pos+1 < len(p.tokens) && p.tokens[p.pos+1].Type == TOKEN_ASSIGN
}


func (p *Parser) ParseModule() *Module {
	mod := &Module{}
	p.skipNewlines()
	for !p.check(TOKEN_EOF) {
		mod.Body = append(mod.Body, p.parseStmt())
		p.skipNewlines()
	}
	return mod
}

func (p *Parser) parseStmt() Stmt {
	if p.check(TOKEN_DEF) {
		return p.parseFuncDef()
	}
	return p.parseSimpleStmt()
}

func (p *Parser) parseBlock() []Stmt {
	p.expect(TOKEN_NEWLINE, "newline before block")
	p.expect(TOKEN_INDENT, "indented block")

	var body []Stmt
	p.skipNewlines()
	for !p.check(TOKEN_DEDENT) && !p.check(TOKEN_EOF) {
		body = append(body, p.parseStmt())
		p.skipNewlines()
	}
	p.expect(TOKEN_DEDENT, "end of block")
	return body
}

func (p *Parser) parseFuncDef() Stmt {
	p.advance() // 'def'
	name := p.expect(TOKEN_NAME, "function name").Lit

	p.expect(TOKEN_LPAREN, "(")
	var params []string
	if !p.check(TOKEN_RPAREN) {
		params = append(params, p.expect(TOKEN_NAME, "parameter name").Lit)
		for p.check(TOKEN_COMMA) {
			p.advance()
			params = append(params, p.expect(TOKEN_NAME, "parameter name").Lit)
		}
	}
	p.expect(TOKEN_RPAREN, ")")
	p.expect(TOKEN_COLON, ":")

	return &FunctionDef{Name: name, Params: params, Body: p.parseBlock()}
}

func (p *Parser) parseSimpleStmt() Stmt {
	var stmt Stmt

	switch {
	case p.check(TOKEN_RETURN):
		p.advance()
		var val Expr
		if !p.check(TOKEN_NEWLINE) && !p.check(TOKEN_EOF) {
			val = p.parseExpr()
		}
		stmt = &ReturnStmt{Value: val}

	case p.check(TOKEN_NAME) && p.peekIsAssign():
		name := p.advance().Lit
		p.expect(TOKEN_ASSIGN, "=")
		stmt = &AssignStmt{Target: name, Value: p.parseExpr()}

	default:
		stmt = &ExprStmt{Value: p.parseExpr()}
	}

	if !p.check(TOKEN_EOF) {
		p.expect(TOKEN_NEWLINE, "newline")
	}
	return stmt
}

func (p *Parser) parseExpr() Expr {
	left := p.parseTerm()
	for p.check(TOKEN_PLUS) || p.check(TOKEN_MINUS) {
		op := p.advance().Lit
		left = &BinOp{Op: op, Left: left, Right: p.parseTerm()}
	}
	return left
}

func (p *Parser) parseTerm() Expr {
	left := p.parseUnary()
	for p.check(TOKEN_STAR) || p.check(TOKEN_SLASH) || p.check(TOKEN_DSLASH) {
		op := p.advance().Lit
		left = &BinOp{Op: op, Left: left, Right: p.parseUnary()}
	}
	return left
}

func (p *Parser) parseUnary() Expr {
	if p.check(TOKEN_MINUS) {
		p.advance()
		return negate(p.parseUnary())
	}
	return p.parsePower()
}

func (p *Parser) parsePower() Expr {
	left := p.parseCall()
	if p.check(TOKEN_DSTAR) {
		p.advance()
		return &BinOp{Op: "**", Left: left, Right: p.parseUnary()}
	}
	return left
}

func negate(e Expr) Expr {
	return &BinOp{Op: "-", Left: &Literal{core.Constant{Type: core.CONST_INT, Int: 0}}, Right: e}
}

func (p *Parser) parseCall() Expr {
	expr := p.parseAtom()
	for p.check(TOKEN_LPAREN) {
		p.advance()
		var args []Expr
		if !p.check(TOKEN_RPAREN) {
			args = append(args, p.parseExpr())
			for p.check(TOKEN_COMMA) {
				p.advance()
				args = append(args, p.parseExpr())
			}
		}
		p.expect(TOKEN_RPAREN, ")")
		expr = &CallExpr{Func: expr, Args: args}
	}
	return expr
}

func (p *Parser) parseAtom() Expr {
	tok := p.cur()
	switch tok.Type {
	case TOKEN_NUMBER:
		p.advance()
		if strings.Contains(tok.Lit, ".") {
			f, err := strconv.ParseFloat(tok.Lit, 64)
			if err != nil {
				panic("SyntaxError: invalid float literal " + tok.Lit)
			}
			return &Literal{core.Constant{Type: core.CONST_FLOAT, Float: f}}
		}
		n, err := strconv.Atoi(tok.Lit)
		if err != nil {
			panic("SyntaxError: invalid int literal " + tok.Lit)
		}
		return &Literal{core.Constant{Type: core.CONST_INT, Int: n}}

	case TOKEN_STRING:
		p.advance()
		return &Literal{core.Constant{Type: core.CONST_STRING, Str: tok.Lit}}

	case TOKEN_NAME:
		p.advance()
		return &NameExpr{Name: tok.Lit}

	case TOKEN_LPAREN:
		p.advance()
		expr := p.parseExpr()
		p.expect(TOKEN_RPAREN, ")")
		return expr
	}
	panic(fmt.Sprintf("SyntaxError: line %d: unexpected token %q", tok.Line, tok.Lit))
}
