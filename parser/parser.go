package parser

import (
	"fmt"
	"snek/lexer"
	"snek/token"
	"strings"
)

type Parser struct {
	l         *lexer.Lexer
	curToken  token.Token
	peekToken token.Token
	errors    []string
}

func New(l *lexer.Lexer) *Parser {
	p := &Parser{
		l:      l,
		errors: []string{},
	}

	p.nextToken()
	p.nextToken()
	return p
}

func (p *Parser) ParseProgram() {
	p.parseBlock(-1)
}

func (p *Parser) parseBlock(prevIndentation int) {
	startIndentation := p.getIndentation()

	if startIndentation <= prevIndentation {
		fmt.Print("not enough indentation")
		return
	}

	statements := []int{}

	curIndentation := startIndentation
	for curIndentation == startIndentation {
		p.nextToken()
		stmt := p.parseStatement()
		// if stmt != nil {
		statements = append(statements, stmt)
		// }

		if !p.curTokenIs(token.NEW_LINE) {
			break
		}
		curIndentation = p.getIndentation()
	}

	if curIndentation > startIndentation {
		fmt.Printf("too much indentation [%d, %d]", curIndentation, startIndentation)
	}

	if len(statements) == 0 {
		fmt.Print("no statements")
	}

	fmt.Print(len(statements))
}

func (p *Parser) parseStatement() int {
	for !p.curTokenIs(token.NEW_LINE) && !p.curTokenIs(token.EOF) {
		p.nextToken()
	}
	return 1
}

func (p *Parser) getIndentation() int {
	if p.curToken.Type != token.NEW_LINE {
		fmt.Print("expected new line")
		return -1
	}

	for p.peekTokenIs(token.NEW_LINE) {
		p.nextToken()
	}

	spaceCount := strings.Count(p.curToken.Literal, " ")
	tabCount := strings.Count(p.curToken.Literal, "\t")

	if spaceCount > 0 && tabCount > 0 {
		fmt.Print("mixed indentation")
		return -1
	}

	return spaceCount + tabCount
}

func (p *Parser) nextToken() {
	p.curToken = p.peekToken
	p.peekToken = p.l.NextToken()
}

func (p *Parser) curTokenIs(t token.TokenType) bool {
	return p.curToken.Type == t
}

func (p *Parser) peekTokenIs(t token.TokenType) bool {
	return p.peekToken.Type == t
}

func (p *Parser) expectPeek(t token.TokenType) bool {
	if p.peekTokenIs(t) {
		p.nextToken()
		return true
	} else {
		p.peekError(t)
		return false
	}
}

func (p *Parser) curError(t token.TokenType) {
	msg := fmt.Sprintf("expected token to be %s, got %s instead", t, p.curToken.Type)
	p.errors = append(p.errors, msg)
}

func (p *Parser) peekError(t token.TokenType) {
	msg := fmt.Sprintf("expected next token to be %s, got %s instead", t, p.peekToken.Type)
	p.errors = append(p.errors, msg)
}
