package parser

import (
	"fmt"
	"snek/ast"
	"snek/token"
)

type (
	statementParseFn func() ast.Node
	prefixParseFn    func() ast.Node
	infixParseFn     func(ast.Node) ast.Node
)

const (
	LOWEST int = iota
	OR
	AND
	NOT
	COMPARE
	SUM
	PRODUCT
	CALL
	ATTR
)

var precedences = map[token.TokenType]int{
	token.OR:       OR,
	token.AND:      AND,
	token.NOT:      NOT,
	token.COMPARE:  COMPARE,
	token.SUM:      SUM,
	token.PRODUCT:  PRODUCT,
	token.LPAREN:   CALL,
	token.LBRACKET: CALL,
	token.DOT:      ATTR,
}

type Parser struct {
	tokens    []token.Token
	pos       int
	curToken  token.Token
	peekToken token.Token
	errors    []string

	statementFns map[token.TokenType]statementParseFn
	prefixFns    map[token.TokenType]prefixParseFn
	infixFns     map[token.TokenType]infixParseFn
}

func New(tokens []token.Token) *Parser {
	p := &Parser{
		tokens: tokens,
		pos:    -1,
		errors: []string{},

		statementFns: make(map[token.TokenType]statementParseFn),
		prefixFns:    make(map[token.TokenType]prefixParseFn),
		infixFns:     make(map[token.TokenType]infixParseFn),
	}

	p.statementFns[token.IF] = p.parseIfStatement

	p.prefixFns[token.IDENTIFIER] = p.parseIdentifierPrefix
	p.prefixFns[token.NUMBER] = p.parseNumberPrefix
	p.prefixFns[token.WHILE] = nil
	p.prefixFns[token.BREAK] = nil
	p.prefixFns[token.CONTINUE] = nil
	p.prefixFns[token.FOR] = nil
	p.prefixFns[token.DEF] = nil
	p.prefixFns[token.RETURN] = nil

	p.infixFns[token.SUM] = p.parseInfixExpression
	p.infixFns[token.PRODUCT] = p.parseInfixExpression

	p.nextToken()
	p.nextToken()
	return p
}

func (p *Parser) Parse() ast.Node {
	statements := []ast.Node{}

	for !p.curTokenIs(token.EOF) {
		if p.curTokenIs(token.NEW_LINE) {
			p.nextToken()
		} else {
			stmt := p.parseStatement()
			if stmt == nil {
				break
			}
			statements = append(statements, stmt)
		}
	}

	return &ast.BlockNode{Statements: statements}
}

func (p *Parser) parseSuite() ast.Node {
	if p.curTokenIs(token.NEW_LINE) {
		p.nextToken()
		p.expect(token.INDENT)

		statements := []ast.Node{}

		for !p.curTokenIs(token.DEDENT) {
			stmt := p.parseStatement()
			if stmt == nil {
				break
			}
			statements = append(statements, stmt)
		}

		p.expect(token.DEDENT)
		return &ast.BlockNode{Statements: statements}
	}

	stmt := p.parseStatementList()
	p.expect(token.NEW_LINE)
	return stmt
}

func (p *Parser) parseStatement() ast.Node {
	if p.isCompoundStatement() {
		return p.parseCompoundStatement()
	} else {
		stmt := p.parseStatementList()
		p.expect(token.NEW_LINE)
		return stmt
	}
}

func (p *Parser) parseSimpleStatement() ast.Node {
	return p.parseAssignmentStatement()
}

func (p *Parser) parseCompoundStatement() ast.Node {
	stmtParsingFn := p.statementFns[p.curToken.Type]
	if stmtParsingFn == nil {
		p.errors = append(p.errors, "no statement parse function for "+p.curToken.Type.String())
		return nil
	}
	return stmtParsingFn()
}

func (p *Parser) parseStatementList() ast.Node {
	statements := []ast.Node{}

	for {
		stmt := p.parseSimpleStatement()
		if stmt == nil {
			break
		}
		statements = append(statements, stmt)

		if !p.curTokenIs(token.SEMICOLON) {
			break
		}
		p.nextToken()
	}

	return &ast.BlockNode{Statements: statements}
}

func (p *Parser) isCompoundStatement() bool {
	switch p.curToken.Type {
	case token.IF, token.WHILE, token.FOR:
		return true
	}
	return false
}

func (p *Parser) parseIfStatement() ast.Node {
	return p.parseIfElifStatement(false)
}

func (p *Parser) parseIfElifStatement(isElif bool) ast.Node {
	defer untrace(trace("parseIfElifStatement"))
	stmt := &ast.IfNode{}

	if !isElif {
		p.expect(token.IF)
	} else {
		p.expect(token.ELIF)
	}

	stmt.Condition = p.parseExpression(LOWEST)

	p.expect(token.COLON)
	stmt.Body = p.parseSuite()

	if p.curTokenIs(token.ELIF) {
		stmt.Else = p.parseIfElifStatement(true)
	} else if p.curTokenIs(token.ELSE) {
		p.nextToken()
		p.expect(token.COLON)
		stmt.Else = p.parseSuite()
	}

	return stmt
}

func (p *Parser) parseAssignmentStatement() ast.Node {
	defer untrace(trace("parseAssignmentStatement"))
	stmt := p.parseExpression(LOWEST)

	if p.curTokenIs(token.ASSIGN) {
		p.nextToken()

		assignment := &ast.AssignmentNode{Target: stmt}
		assignment.Value = p.parseExpression(LOWEST)
		stmt = assignment
	}

	return stmt
}

func (p *Parser) parseExpression(precedence int) ast.Node {
	defer untrace(trace("parseExpression"))
	prefix := p.prefixFns[p.curToken.Type]
	if prefix == nil {
		p.errors = append(p.errors, "no prefix parse function for "+p.curToken.Type.String())
		return nil
	}

	leftExpr := prefix()

	for precedence < getPrecedence(p.curToken.Type) {
		infix := p.infixFns[p.curToken.Type]
		if infix == nil {
			return leftExpr
		}
		leftExpr = infix(leftExpr)
	}

	return leftExpr
}

func (p *Parser) parseIdentifierPrefix() ast.Node {
	defer p.nextToken()
	return &ast.IdentifierNode{Name: p.curToken.Literal}
}

func (p *Parser) parseNumberPrefix() ast.Node {
	defer p.nextToken()
	return &ast.NumberNode{Value: p.curToken.Literal}
}

func (p *Parser) parseInfixExpression(left ast.Node) ast.Node {
	defer untrace(trace("parseInfixExpression"))
	expression := &ast.InfixNode{
		Operator: p.curToken.Literal,
		Left:     left,
	}
	precedence := getPrecedence(p.curToken.Type)
	p.nextToken()
	expression.Right = p.parseExpression(precedence)
	return expression
}

func (p *Parser) nextToken() {
	p.curToken = p.peekToken

	p.pos++
	if p.pos >= len(p.tokens) {
		return
	}

	p.peekToken = p.tokens[p.pos]
}

func (p *Parser) curTokenIs(t token.TokenType) bool {
	return p.curToken.Type == t
}

func (p *Parser) expect(t token.TokenType) bool {
	if p.curTokenIs(t) {
		p.nextToken()
		return true
	} else {
		p.curError(t)
		return false
	}
}

func (p *Parser) curError(t token.TokenType) {
	msg := fmt.Sprintf("expected token to be %s, got %s instead", t, p.curToken.Type)
	p.errors = append(p.errors, msg)
}

func getPrecedence(tok token.TokenType) int {
	if p, ok := precedences[tok]; ok {
		return p
	}
	return LOWEST
}

func (p *Parser) Errors() []string {
	return p.errors
}
