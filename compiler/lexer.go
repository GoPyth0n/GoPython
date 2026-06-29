// lexer.go
package compiler

import (
	"fmt"
	"strings"
	"unicode"
)

type TokenType int

const (
	TOKEN_EOF TokenType = iota
	TOKEN_NEWLINE
	TOKEN_INDENT
	TOKEN_DEDENT

	TOKEN_NAME
	TOKEN_NUMBER
	TOKEN_STRING

	TOKEN_DEF
	TOKEN_RETURN

	TOKEN_ASSIGN // =
	TOKEN_PLUS   // +
	TOKEN_MINUS  // -
	TOKEN_STAR   // *
	TOKEN_DSTAR  // **
	TOKEN_SLASH  // /
	TOKEN_DSLASH // //
	TOKEN_LPAREN // (
	TOKEN_RPAREN // )
	TOKEN_COLON  // :
	TOKEN_COMMA  // ,
)

var keywords = map[string]TokenType{
	"def":    TOKEN_DEF,
	"return": TOKEN_RETURN,
}

type Token struct {
	Type TokenType
	Lit  string
	Line int
}

type Lexer struct {
	src     []rune
	pos     int
	line    int
	indents []int
	tokens  []Token
}

func NewLexer(src string) *Lexer {
	// normalize line endings so \r doesn't need special-casing everywhere
	src = strings.ReplaceAll(src, "\r\n", "\n")
	return &Lexer{
		src:     []rune(src),
		line:    1,
		indents: []int{0},
	}
}

func (l *Lexer) peek() rune {
	if l.pos >= len(l.src) {
		return 0
	}
	return l.src[l.pos]
}

func (l *Lexer) peekAt(off int) rune {
	if l.pos+off >= len(l.src) {
		return 0
	}
	return l.src[l.pos+off]
}

func (l *Lexer) advance() rune {
	r := l.src[l.pos]
	l.pos++
	return r
}

func (l *Lexer) emit(t TokenType, lit string) {
	l.tokens = append(l.tokens, Token{Type: t, Lit: lit, Line: l.line})
}

func (l *Lexer) Tokenize() []Token {
	for l.pos < len(l.src) {
		l.lexLogicalLine()
	}
	for len(l.indents) > 1 {
		l.indents = l.indents[:len(l.indents)-1]
		l.emit(TOKEN_DEDENT, "")
	}
	l.emit(TOKEN_EOF, "")
	return l.tokens
}

func (l *Lexer) lexLogicalLine() {
	indent := 0
	for l.peek() == ' ' || l.peek() == '\t' {
		if l.advance() == '\t' {
			indent += 8 - (indent % 8)
		} else {
			indent++
		}
	}

	switch l.peek() {
	case '\n':
		l.advance()
		l.line++
		return
	case '#':
		l.skipComment()
		return
	case 0:
		return
	}

	cur := l.indents[len(l.indents)-1]
	if indent > cur {
		l.indents = append(l.indents, indent)
		l.emit(TOKEN_INDENT, "")
	} else if indent < cur {
		for len(l.indents) > 1 && l.indents[len(l.indents)-1] > indent {
			l.indents = l.indents[:len(l.indents)-1]
			l.emit(TOKEN_DEDENT, "")
		}
		if l.indents[len(l.indents)-1] != indent {
			panic(fmt.Sprintf("IndentationError: line %d: unindent does not match any outer indentation level", l.line))
		}
	}

	for {
		c := l.peek()
		if c == 0 || c == '\n' {
			break
		}
		if c == ' ' || c == '\t' {
			l.advance()
			continue
		}
		if c == '#' {
			l.skipComment()
			break
		}
		l.lexToken()
	}

	startLine := l.line
	if l.peek() == '\n' {
		l.advance()
		l.line++
	}
	l.tokens = append(l.tokens, Token{Type: TOKEN_NEWLINE, Line: startLine})
}

func (l *Lexer) skipComment() {
	for l.peek() != '\n' && l.peek() != 0 {
		l.advance()
	}
}

func (l *Lexer) lexToken() {
	c := l.peek()
	switch {
	case unicode.IsLetter(c) || c == '_':
		l.lexName()
	case unicode.IsDigit(c):
		l.lexNumber()
	case c == '"' || c == '\'':
		l.lexString()
	default:
		l.lexOperator()
	}
}

func (l *Lexer) lexName() {
	start := l.pos
	for unicode.IsLetter(l.peek()) || unicode.IsDigit(l.peek()) || l.peek() == '_' {
		l.advance()
	}
	lit := string(l.src[start:l.pos])
	if tt, ok := keywords[lit]; ok {
		l.emit(tt, lit)
		return
	}
	l.emit(TOKEN_NAME, lit)
}

func (l *Lexer) lexNumber() {
	start := l.pos
	for unicode.IsDigit(l.peek()) {
		l.advance()
	}
	if l.peek() == '.' && unicode.IsDigit(l.peekAt(1)) {
		l.advance()
		for unicode.IsDigit(l.peek()) {
			l.advance()
		}
	}
	l.emit(TOKEN_NUMBER, string(l.src[start:l.pos]))
}

func (l *Lexer) lexString() {
	quote := l.advance()
	start := l.pos
	for l.peek() != quote {
		if l.peek() == 0 || l.peek() == '\n' {
			panic(fmt.Sprintf("SyntaxError: line %d: unterminated string literal", l.line))
		}
		l.advance()
	}
	lit := string(l.src[start:l.pos])
	l.advance() // closing quote
	l.emit(TOKEN_STRING, lit)
}

func (l *Lexer) lexOperator() {
	c := l.advance()
	switch c {
	case '=':
		l.emit(TOKEN_ASSIGN, "=")
	case '+':
		l.emit(TOKEN_PLUS, "+")
	case '-':
		l.emit(TOKEN_MINUS, "-")
	case '*':
		if l.peek() == '*' {
			l.advance()
			l.emit(TOKEN_DSTAR, "**")
		} else {
			l.emit(TOKEN_STAR, "*")
		}
	case '/':
		if l.peek() == '/' {
			l.advance()
			l.emit(TOKEN_DSLASH, "//")
		} else {
			l.emit(TOKEN_SLASH, "/")
		}
	case '(':
		l.emit(TOKEN_LPAREN, "(")
	case ')':
		l.emit(TOKEN_RPAREN, ")")
	case ':':
		l.emit(TOKEN_COLON, ":")
	case ',':
		l.emit(TOKEN_COMMA, ",")
	default:
		panic(fmt.Sprintf("SyntaxError: line %d: unexpected character %q", l.line, string(c)))
	}
}
