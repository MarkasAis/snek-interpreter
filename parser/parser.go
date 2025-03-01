package parser

import (
	"fmt"
	"snek/ast"
	"snek/token"
)

type (
	statementParseFn func() (ast.Node, error)
	prefixParseFn    func() (ast.Node, error)
	infixParseFn     func(ast.Node) (ast.Node, error)
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

type ParseError struct {
	Value string
}

func (e *ParseError) Error() string { return e.Value }

type Parser struct {
	tokens    []token.Token
	pos       int
	curToken  token.Token
	peekToken token.Token

	statementFns map[token.TokenType]statementParseFn
	prefixFns    map[token.TokenType]prefixParseFn
	infixFns     map[token.TokenType]infixParseFn
}

func New(tokens []token.Token) *Parser {
	p := &Parser{
		tokens: tokens,
		pos:    -1,

		statementFns: make(map[token.TokenType]statementParseFn),
		prefixFns:    make(map[token.TokenType]prefixParseFn),
		infixFns:     make(map[token.TokenType]infixParseFn),
	}

	p.statementFns[token.IF] = p.parseIfStatement
	p.statementFns[token.WHILE] = p.parseWhileStatement
	p.statementFns[token.BREAK] = nil
	p.statementFns[token.CONTINUE] = nil
	p.statementFns[token.FOR] = nil
	p.statementFns[token.DEF] = nil
	p.statementFns[token.RETURN] = nil

	p.prefixFns[token.IDENTIFIER] = p.parseIdentifierPrefix
	p.prefixFns[token.NUMBER] = p.parseNumberPrefix

	p.infixFns[token.SUM] = p.parseInfixExpression
	p.infixFns[token.PRODUCT] = p.parseInfixExpression

	p.nextToken()
	p.nextToken()
	return p
}

func (p *Parser) Parse() (ast.Node, error) {
	block := &ast.BlockNode{Statements: []ast.Node{}}

	for !p.curTokenIs(token.EOF) {
		if p.curTokenIs(token.NEW_LINE) {
			p.nextToken()
		} else {
			stmt, err := p.parseStatement()
			if err != nil {
				return block, err
			}
			block.Statements = append(block.Statements, stmt)
		}
	}

	return block, nil
}

func (p *Parser) parseSuite() (ast.Node, error) {
	defer untrace(trace("parseSuite"))

	if p.curTokenIs(token.NEW_LINE) {
		p.nextToken()
		if err := p.expect(token.INDENT); err != nil {
			return nil, err
		}

		block := &ast.BlockNode{Statements: []ast.Node{}}

		for !p.curTokenIs(token.DEDENT) {
			stmt, err := p.parseStatement()
			if err != nil {
				return block, err
			}
			block.Statements = append(block.Statements, stmt)
		}

		if err := p.expect(token.DEDENT); err != nil {
			return block, err
		}
		return block, nil
	}

	stmt, err := p.parseStatementList()
	if err != nil {
		return stmt, err
	}

	if err := p.expect(token.NEW_LINE); err != nil {
		return stmt, err
	}

	return stmt, nil
}

func (p *Parser) parseStatement() (ast.Node, error) {
	defer untrace(trace("parseStatement"))
	if p.isCompoundStatement() {
		return p.parseCompoundStatement()
	} else {
		stmt, err := p.parseStatementList()
		if err != nil {
			return stmt, err
		}
		if err := p.expect(token.NEW_LINE); err != nil {
			return stmt, err
		}
		return stmt, nil
	}
}

func (p *Parser) parseSimpleStatement() (ast.Node, error) {
	defer untrace(trace("parseSimpleStatement"))
	switch p.curToken.Type {
	case token.PASS, token.BREAK, token.CONTINUE:
		stmt := ast.ControlNode{Type: p.curToken.Literal}
		p.nextToken()
		return &stmt, nil
	}
	return p.parseAssignmentStatement()
}

func (p *Parser) parseCompoundStatement() (ast.Node, error) {
	defer untrace(trace("parseCompoundStatement"))
	stmtParsingFn := p.statementFns[p.curToken.Type]
	if stmtParsingFn == nil {
		return nil, &ParseError{Value: fmt.Sprintf("no statement parse function for %s", p.curToken.Type)}
	}
	return stmtParsingFn()
}

func (p *Parser) parseStatementList() (ast.Node, error) {
	defer untrace(trace("parseStatementList"))
	block := &ast.BlockNode{Statements: []ast.Node{}}

	for {
		stmt, err := p.parseSimpleStatement()
		if err != nil {
			return block, err
		}
		block.Statements = append(block.Statements, stmt)

		if !p.curTokenIs(token.SEMICOLON) {
			break
		}
		p.nextToken()
	}

	return block, nil
}

func (p *Parser) isCompoundStatement() bool {
	switch p.curToken.Type {
	case token.IF, token.WHILE, token.FOR:
		return true
	}
	return false
}

func (p *Parser) parseIfStatement() (ast.Node, error) {
	return p.parseIfElifStatement(false)
}

func (p *Parser) parseIfElifStatement(isElif bool) (ast.Node, error) {
	defer untrace(trace("parseIfElifStatement"))
	stmt := &ast.IfNode{}

	startToken := token.IF
	if isElif {
		startToken = token.ELIF
	}

	if err := p.expect(startToken); err != nil {
		return stmt, err
	}

	res, err := p.parseExpression(LOWEST)
	stmt.Condition = res
	if err != nil {
		return stmt, err
	}

	if err := p.expect(token.COLON); err != nil {
		return stmt, err
	}

	res, err = p.parseSuite()
	stmt.Body = res
	if err != nil {
		return stmt, err
	}

	if p.curTokenIs(token.ELIF) {
		res, err := p.parseIfElifStatement(true)
		stmt.Else = res
		if err != nil {
			return stmt, err
		}
	} else if p.curTokenIs(token.ELSE) {
		p.nextToken()
		if err := p.expect(token.COLON); err != nil {
			return stmt, err
		}
		res, err := p.parseSuite()
		stmt.Else = res
		if err != nil {
			return stmt, err
		}
	}

	return stmt, nil
}

func (p *Parser) parseWhileStatement() (ast.Node, error) {
	defer untrace(trace("parseWhileStatement"))
	stmt := &ast.WhileNode{}

	if err := p.expect(token.WHILE); err != nil {
		return stmt, err
	}

	res, err := p.parseExpression(LOWEST)
	stmt.Condition = res
	if err != nil {
		return stmt, err
	}

	if err := p.expect(token.COLON); err != nil {
		return stmt, err
	}

	res, err = p.parseSuite()
	stmt.Body = res
	if err != nil {
		return stmt, err
	}

	if p.curTokenIs(token.ELSE) {
		p.nextToken()
		if err := p.expect(token.COLON); err != nil {
			return stmt, err
		}

		res, err := p.parseSuite()
		stmt.Else = res
		if err != nil {
			return stmt, err
		}
	}

	return stmt, nil
}

func (p *Parser) parseAssignmentStatement() (ast.Node, error) {
	defer untrace(trace("parseAssignmentStatement"))
	stmt, err := p.parseExpression(LOWEST)
	if err != nil {
		return stmt, err
	}

	if p.curTokenIs(token.ASSIGN) {
		p.nextToken()

		assignment := &ast.AssignmentNode{Target: stmt}
		res, err := p.parseExpression(LOWEST)
		assignment.Value = res
		if err != nil {
			return stmt, err
		}

		stmt = assignment
	}

	return stmt, nil
}

func (p *Parser) parseExpression(precedence int) (ast.Node, error) {
	defer untrace(trace("parseExpression"))
	prefix := p.prefixFns[p.curToken.Type]
	if prefix == nil {
		return nil, &ParseError{Value: fmt.Sprintf("no prefix parse function for %s", p.curToken.Type.String())}
	}

	leftExpr, err := prefix()
	if err != nil {
		return leftExpr, err
	}

	for precedence < getPrecedence(p.curToken.Type) {
		infix := p.infixFns[p.curToken.Type]
		if infix == nil {
			return leftExpr, nil
		}
		res, err := infix(leftExpr)
		leftExpr = res
		if err != nil {
			return leftExpr, err
		}
	}

	return leftExpr, nil
}

func (p *Parser) parseIdentifierPrefix() (ast.Node, error) {
	defer p.nextToken()
	return &ast.IdentifierNode{Name: p.curToken.Literal}, nil
}

func (p *Parser) parseNumberPrefix() (ast.Node, error) {
	defer p.nextToken()
	return &ast.NumberNode{Value: p.curToken.Literal}, nil
}

func (p *Parser) parseInfixExpression(left ast.Node) (ast.Node, error) {
	defer untrace(trace("parseInfixExpression"))
	expression := &ast.InfixNode{
		Operator: p.curToken.Literal,
		Left:     left,
	}
	precedence := getPrecedence(p.curToken.Type)
	p.nextToken()

	res, err := p.parseExpression(precedence)
	if err != nil {
		return expression, err
	}

	expression.Right = res
	return expression, nil
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

func (p *Parser) expect(t token.TokenType) error {
	if p.curTokenIs(t) {
		p.nextToken()
		return nil
	} else {
		return p.curError(t)
	}
}

func (p *Parser) curError(t token.TokenType) error {
	return &ParseError{Value: fmt.Sprintf("expected token to be %s, got %s instead", t, p.curToken.Type)}
}

func getPrecedence(tok token.TokenType) int {
	if p, ok := precedences[tok]; ok {
		return p
	}
	return LOWEST
}
