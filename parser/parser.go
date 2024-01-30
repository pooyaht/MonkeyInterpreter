package parser

import (
	"errors"
	"fmt"

	"github.com/pooyaht/MonkeyInterpreter/ast"
	"github.com/pooyaht/MonkeyInterpreter/lexer"
	"github.com/pooyaht/MonkeyInterpreter/token"
)

type Parser struct {
	l         *lexer.Lexer
	currToken token.Token
	peekToken token.Token
	errors    []error
}

func New(l *lexer.Lexer) *Parser {
	p := &Parser{
		l:      l,
		errors: []error{},
	}
	p.nextToken()
	p.nextToken()
	return p
}

func (p *Parser) Errors() []error {
	return p.errors
}

func (p *Parser) nextToken() {
	p.currToken = p.peekToken
	p.peekToken = p.l.NextToken()
}

func (p *Parser) ParseProgram() *ast.Program {
	program := &ast.Program{}
	program.Statements = []ast.Statement{}
	for p.currToken.Type != token.EOF {
		stmt := p.parseStatement()

		if stmt != nil {
			program.Statements = append(program.Statements, stmt)
		}
		p.nextToken()
	}
	return program
}

func (p *Parser) parseStatement() ast.Statement {

	switch p.currToken.Type {
	case token.LET:
		return p.parseLetStatement()
	default:
		return nil
	}
}

func (p *Parser) parseLetStatement() *ast.LetStatement {
	stmt := &ast.LetStatement{Token: p.currToken}

	if !p.match(p.peekToken.Type, token.IDENT) {
		return nil
	}
	p.nextToken()

	stmt.Name = &ast.Identifier{Token: p.currToken, Value: p.currToken.Literal}

	if !p.match(p.peekToken.Type, token.ASSIGN) {
		return nil
	}
	// TODO : skipping the rvalue expression
	for !p.match(p.currToken.Type, token.SEMICOLON) {
		p.nextToken()
	}
	return stmt
}

func (p *Parser) match(token token.TokenType, expectedToken token.TokenType) bool {
	if token == expectedToken {
		return true
	} else {
		p.peekError(expectedToken)
		return false
	}
}

func (p *Parser) peekError(t token.TokenType) {
	errMsg := fmt.Sprintf("expected next token to be %s, got %s instead",
		t, p.peekToken.Type)
	err := errors.New(errMsg)
	p.errors = append(p.errors, err)
}
