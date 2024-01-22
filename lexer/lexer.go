package lexer

import (
	"github.com/pooyaht/MonkeyInterpreter/token"
)

type Lexer struct {
	input    string
	position int
}

func New(input string) Lexer {
	return Lexer{input: input}
}

func (l *Lexer) NextToken() token.Token {
	l.skipWhitespace()
	var tok token.Token
	ch := l.peek()
	switch ch {
	case '=':
		if l.peekNext() == '=' {
			l.advence()
			tok = token.Token{
				Type:    token.EQ,
				Literal: string(ch) + string(l.peek()),
			}
		} else {
			tok = newToken(token.ASSIGN, ch)
		}
	case '!':
		if l.peekNext() == '=' {
			l.advence()
			tok = token.Token{
				Type:    token.NOT_EQ,
				Literal: string(ch) + string(l.peek()),
			}
		} else {
			tok = newToken(token.BANG, ch)
		}
	case ';':
		tok = newToken(token.SEMICOLON, ch)
	case '(':
		tok = newToken(token.LPAREN, ch)
	case ')':
		tok = newToken(token.RPAREN, ch)
	case ',':
		tok = newToken(token.COMMA, ch)
	case '+':
		tok = newToken(token.PLUS, ch)
	case '-':
		tok = newToken(token.MINUS, ch)
	case '*':
		tok = newToken(token.ASTERISK, ch)
	case '/':
		if l.peekNext() == '/' {
			for l.peek() != '\n' && !l.isAtEnd() {
				l.advence()
			}
			return l.NextToken()
		} else {
			tok = newToken(token.SLASH, ch)
		}
	case '<':
		tok = newToken(token.LT, ch)
	case '>':
		tok = newToken(token.GT, ch)
	case '{':
		tok = newToken(token.LBRACE, ch)
	case '}':
		tok = newToken(token.RBRACE, ch)
	case 0:
		tok.Literal = ""
		tok.Type = token.EOF
	default:
		if isLetter(ch) {
			tok.Literal = l.readIdentifier()
			tok.Type = token.LookupIdent(tok.Literal)
			return tok
		} else if isDigit(ch) {
			tok.Literal = l.readNumber()
			tok.Type = token.INT
			return tok
		} else {
			tok = newToken(token.ILLEGAL, ch)
		}
	}
	l.advence()
	return tok
}

func (l *Lexer) readIdentifier() string {
	position := l.position
	for isLetter(l.peek()) {
		l.advence()
	}
	return l.input[position:l.position]
}

func (l *Lexer) readNumber() string {
	position := l.position
	for isDigit(l.peek()) {
		l.advence()
	}
	return l.input[position:l.position]
}

func (l *Lexer) skipWhitespace() {
	ch := l.peek()
	for ch == ' ' || ch == '\t' || ch == '\n' || ch == '\r' {
		l.advence()
		ch = l.peek()
	}
}

func (l *Lexer) advence() {
	if l.isAtEnd() {
		return
	}
	l.position++
}

func (l *Lexer) peek() byte {
	if l.isAtEnd() {
		return 0
	}
	return l.input[l.position]
}

func (l *Lexer) peekNext() byte {
	if l.position+1 >= len(l.input) {
		return 0
	}
	return l.input[l.position+1]
}

func (l *Lexer) isAtEnd() bool {
	return l.position >= len(l.input)
}

func isLetter(ch byte) bool {
	return 'a' <= ch && ch <= 'z' || 'A' <= ch && ch <= 'Z' || ch == '_'
}

func isDigit(ch byte) bool {
	return '0' <= ch && ch <= '9'
}

func newToken(tokenType token.TokenType, ch byte) token.Token {
	return token.Token{Type: tokenType, Literal: string(ch)}
}
