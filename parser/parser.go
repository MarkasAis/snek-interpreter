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
	PREFIX
	EXP
	ATTR
)

var precedences = map[token.TokenType]int{
	token.OR:       OR,
	token.AND:      AND,
	token.NOT:      NOT,
	token.COMPARE:  COMPARE,
	token.SUM:      SUM,
	token.PRODUCT:  PRODUCT,
	token.EXP:      EXP,
	token.LPAREN:   ATTR,
	token.LBRACKET: ATTR,
	token.DOT:      ATTR,
}

type ParseError struct {
	Value string
}

func (e *ParseError) Error() string { return e.Value }

// Reference: https://docs.python.org/3/reference/grammar.html

type Parser struct {
	tokens    []token.Token
	pos       int
	curToken  token.Token
	peekToken token.Token

	simpleStatementFns  map[token.TokenType]statementParseFn
	compundStatementFns map[token.TokenType]statementParseFn
	prefixFns           map[token.TokenType]prefixParseFn
	infixFns            map[token.TokenType]infixParseFn
}

func New(tokens []token.Token) *Parser {
	p := &Parser{
		tokens: tokens,
		pos:    -1,

		simpleStatementFns:  make(map[token.TokenType]statementParseFn),
		compundStatementFns: make(map[token.TokenType]statementParseFn),
		prefixFns:           make(map[token.TokenType]prefixParseFn),
		infixFns:            make(map[token.TokenType]infixParseFn),
	}

	p.simpleStatementFns[token.PASS] = p.parseControlStatement
	p.simpleStatementFns[token.BREAK] = p.parseControlStatement
	p.simpleStatementFns[token.CONTINUE] = p.parseControlStatement
	p.simpleStatementFns[token.RETURN] = p.parseReturnStatement
	p.simpleStatementFns[token.IMPORT] = nil
	p.simpleStatementFns[token.GLOBAL] = nil // TODO: add nonlocal

	p.compundStatementFns[token.DEF] = p.parseFunctionDef
	p.compundStatementFns[token.IF] = p.parseIfStatement
	p.compundStatementFns[token.FOR] = p.parseForStatement
	p.compundStatementFns[token.WHILE] = p.parseWhileStatement

	p.prefixFns[token.IDENTIFIER] = p.parseIdentifierPrefix
	p.prefixFns[token.NUMBER] = p.parseNumberPrefix
	p.prefixFns[token.LPAREN] = p.parseGroupPrefix
	p.prefixFns[token.SUM] = p.parseExpressionPrefix

	p.infixFns[token.OR] = p.parseExpressionInfix
	p.infixFns[token.AND] = p.parseExpressionInfix
	p.infixFns[token.COMPARE] = p.parseExpressionInfix
	p.infixFns[token.SUM] = p.parseExpressionInfix
	p.infixFns[token.PRODUCT] = p.parseExpressionInfix
	p.infixFns[token.EXP] = p.parseExpressionInfix
	p.infixFns[token.DOT] = p.parseExpressionInfix
	p.infixFns[token.LPAREN] = p.parseCallInfix
	p.infixFns[token.LBRACKET] = p.parseSlicesInfix

	p.nextToken()
	p.nextToken()
	return p
}

func (p *Parser) ParseFile() (ast.Node, error) {
	defer untrace(trace("file"))
	return p.parseStatements(token.EOF)
}

func (p *Parser) parseBlock() (ast.Node, error) {
	defer untrace(trace("block"))
	if p.curTokenIs(token.NEW_LINE) {
		p.nextToken()

		if err := p.expect(token.INDENT); err != nil {
			return nil, err
		}

		res, err := p.parseStatements(token.DEDENT)
		if err != nil {
			return res, err
		}

		if err := p.expect(token.DEDENT); err != nil {
			return res, err
		}

		return res, nil
	}

	return p.parseSimpleStatements()
}

func (p *Parser) parseStatements(endToken token.TokenType) (ast.Node, error) {
	defer untrace(trace("statements"))
	block := &ast.BlockNode{Statements: []ast.Node{}}

	for !p.curTokenIs(endToken) {
		stmt, err := p.parseStatement()
		if err != nil {
			return block, err
		}
		block.Statements = append(block.Statements, stmt)

	}

	return block, nil
}

func (p *Parser) parseStatement() (ast.Node, error) {
	defer untrace(trace("statement"))
	if p.isCompoundStatement() {
		return p.parseCompoundStatement()
	} else {
		return p.parseSimpleStatements()
	}
}

func (p *Parser) parseSimpleStatements() (ast.Node, error) {
	defer untrace(trace("simpleStatements"))
	block := &ast.BlockNode{Statements: []ast.Node{}}

	for !p.curTokenIs(token.NEW_LINE) {
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

	p.nextToken()

	if len(block.Statements) == 0 {
		return block, &ParseError{Value: "empty simple statements"}
	}

	return block, nil
}

func (p *Parser) parseSimpleStatement() (ast.Node, error) {
	defer untrace(trace("simpleStatement"))
	stmtParsingFn := p.simpleStatementFns[p.curToken.Type]
	if stmtParsingFn != nil {
		return stmtParsingFn()
	}

	startPos := p.pos
	if res, err := p.parseAssignmentStatement(); err == nil {
		return res, nil
	}

	p.setPos(startPos)
	return p.parseExpressions()
}

func (p *Parser) parseExpressions() (ast.Node, error) {
	defer untrace(trace("expressions"))
	n := &ast.ExpressionsNode{}

	for {
		startPos := p.pos
		res, err := p.parseExpression(LOWEST)
		if err != nil {
			if len(n.Expressions) == 0 { // Must not be empty
				return n, err
			} else {
				p.setPos(startPos)
				break
			}
		}

		n.Expressions = append(n.Expressions, res)

		if p.curTokenIs(token.COMMA) {
			p.nextToken()
		}
	}

	return n, nil
}

// Simple statement parsers

func (p *Parser) parseControlStatement() (ast.Node, error) {
	defer untrace(trace("controlStatement"))
	stmt := &ast.ControlNode{Type: p.curToken.Literal}
	p.nextToken()
	return stmt, nil
}

func (p *Parser) parseReturnStatement() (ast.Node, error) {
	defer untrace(trace("returnStatement"))

	if err := p.expect(token.RETURN); err != nil {
		return nil, err
	}

	stmt := &ast.ReturnNode{}
	res, err := p.parseExpression(LOWEST)
	stmt.Value = res
	if err != nil {
		return stmt, err
	}

	return stmt, nil
}

// Compound statement parsers

func (p *Parser) parseFunctionDef() (ast.Node, error) {
	if err := p.expect(token.DEF); err != nil {
		return nil, err
	}

	stmt := &ast.FunctionDefNode{}
	res, err := p.parseIdentifierPrefix()
	stmt.Name = res
	if err != nil {
		return stmt, err
	}

	if err := p.expect(token.LPAREN); err != nil {
		return stmt, err
	}

	params, err := p.parseParams()
	stmt.Params = params
	if err != nil {
		return stmt, err
	}

	if err := p.expect(token.RPAREN); err != nil {
		return stmt, err
	}

	if err := p.expect(token.COLON); err != nil {
		return stmt, err
	}

	res, err = p.parseBlock()
	stmt.Body = res
	if err != nil {
		return stmt, err
	}

	return stmt, nil
}

func (p *Parser) parseParams() ([]ast.Node, error) {
	params := []ast.Node{}
	requireDefault := false

	for !p.curTokenIs(token.RPAREN) {
		res, err := p.parseParam(requireDefault)
		if err != nil {
			return params, err
		}

		params = append(params, res)
		if res.DefaultValue != nil {
			requireDefault = true
		}

		if p.curTokenIs(token.COMMA) {
			p.nextToken()
		} else {
			break
		}
	}

	return params, nil
}

func (p *Parser) parseParam(requireDefault bool) (*ast.ParamNode, error) {
	n := &ast.ParamNode{}
	res, err := p.parseIdentifierPrefix()
	n.Name = res
	if err != nil {
		return n, err
	}

	if requireDefault && !p.curTokenIs(token.ASSIGN) {
		return n, p.curError(token.ASSIGN)
	}

	if p.curTokenIs(token.ASSIGN) {
		p.nextToken()
		res, err := p.parseExpression(LOWEST)
		n.DefaultValue = res
		if err != nil {
			return n, err
		}
	}

	return n, nil
}

func (p *Parser) parseCompoundStatement() (ast.Node, error) {
	defer untrace(trace("compoundStatement"))
	stmtParsingFn := p.compundStatementFns[p.curToken.Type]
	if stmtParsingFn == nil {
		return nil, &ParseError{Value: fmt.Sprintf("no statement parse function for %s", p.curToken.Type)}
	}
	return stmtParsingFn()
}

func (p *Parser) parseIfStatement() (ast.Node, error) {
	defer untrace(trace("ifStatement"))
	return p.parseIfElifStatement(false)
}

func (p *Parser) parseIfElifStatement(isElif bool) (ast.Node, error) {
	defer untrace(trace("ifElifStatement"))
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

	res, err = p.parseBlock()
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
		res, err := p.parseElseBlock()
		stmt.Else = res
		if err != nil {
			return stmt, err
		}
	}

	return stmt, nil
}

func (p *Parser) parseWhileStatement() (ast.Node, error) {
	defer untrace(trace("whileStatement"))
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

	res, err = p.parseBlock()
	stmt.Body = res
	if err != nil {
		return stmt, err
	}

	if p.curTokenIs(token.ELSE) {
		res, err := p.parseElseBlock()
		stmt.Else = res
		if err != nil {
			return stmt, err
		}
	}

	return stmt, nil
}

func (p *Parser) parseElseBlock() (ast.Node, error) {
	defer untrace(trace("elseBlock"))
	if err := p.expect(token.ELSE); err != nil {
		return nil, err
	}

	if err := p.expect(token.COLON); err != nil {
		return nil, err
	}

	return p.parseBlock()
}

func (p *Parser) parseForStatement() (ast.Node, error) {
	defer untrace(trace("forStatement"))
	if err := p.expect(token.FOR); err != nil {
		return nil, err
	}

	stmt := &ast.ForNode{}
	res, err := p.parseTargets()
	stmt.Targets = res
	if err != nil {
		return stmt, err
	}

	if err := p.expect(token.IN); err != nil {
		return nil, err
	}

	res, err = p.parseExpression(LOWEST)
	stmt.Values = res
	if err != nil {
		return stmt, err
	}

	if err := p.expect(token.COLON); err != nil {
		return nil, err
	}

	res, err = p.parseBlock()
	stmt.Body = res
	if err != nil {
		return stmt, err
	}

	if p.curTokenIs(token.ELSE) {
		res, err := p.parseElseBlock()
		stmt.Else = res
		if err != nil {
			return stmt, err
		}
	}

	return stmt, nil
}

func (p *Parser) parseTargets() (ast.Node, error) {
	defer untrace(trace("targets"))
	return p.parseExpression(LOWEST) // TODO: implement real version
}

func (p *Parser) parseAssignmentStatement() (ast.Node, error) {
	defer untrace(trace("assignmentStatement"))
	stmt, err := p.parseExpression(LOWEST)
	if err != nil {
		return stmt, err
	}

	if !p.curTokenIs(token.ASSIGN) {
		return stmt, p.curError(token.ASSIGN)
	}

	assignment := &ast.AssignmentNode{Target: stmt, Operator: p.curToken.Literal}
	p.nextToken()

	res, err := p.parseExpression(LOWEST)
	assignment.Value = res
	if err != nil {
		return stmt, err
	}

	stmt = assignment

	return stmt, nil
}

func (p *Parser) parseExpression(precedence int) (ast.Node, error) {
	defer untrace(trace("expression"))

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
	defer untrace(trace("identifierPrefix"))
	if !p.curTokenIs(token.IDENTIFIER) {
		return nil, p.curError(token.IDENTIFIER)
	}

	defer p.nextToken()
	return &ast.IdentifierNode{Name: p.curToken.Literal}, nil
}

func (p *Parser) parseNumberPrefix() (ast.Node, error) {
	defer untrace(trace("numberPrefix"))
	if !p.curTokenIs(token.NUMBER) {
		return nil, p.curError(token.NUMBER)
	}

	defer p.nextToken()
	return &ast.NumberNode{Value: p.curToken.Literal}, nil
}

func (p *Parser) parseExpressionPrefix() (ast.Node, error) {
	defer untrace(trace("expressionPrefix"))
	expression := &ast.PrefixNode{
		Operator: p.curToken.Literal,
	}
	p.nextToken()

	res, err := p.parseExpression(PREFIX)
	expression.Right = res
	if err != nil {
		return expression, err
	}

	return expression, nil
}

func (p *Parser) parseGroupPrefix() (ast.Node, error) {
	defer untrace(trace("groupPrefix"))
	if !p.curTokenIs(token.LPAREN) {
		return nil, p.curError(token.LPAREN)
	}

	p.nextToken()
	res, err := p.parseExpression(LOWEST)
	if err != nil {
		return res, err
	}

	if err := p.expect(token.RPAREN); err != nil {
		return res, err
	}

	return res, nil
}

func (p *Parser) parseExpressionInfix(left ast.Node) (ast.Node, error) {
	defer untrace(trace("expressionInfix"))
	expression := &ast.InfixNode{
		Operator: p.curToken.Literal,
		Left:     left,
	}
	precedence := getPrecedence(p.curToken.Type)
	if p.curTokenIs(token.EXP) {
		precedence -= 1
	}

	p.nextToken()

	res, err := p.parseExpression(precedence)
	if err != nil {
		return expression, err
	}

	expression.Right = res
	return expression, nil
}

func (p *Parser) parseCallInfix(left ast.Node) (ast.Node, error) {
	defer untrace(trace("callInfix"))
	expression := &ast.CallNode{
		Function: left,
	}

	p.nextToken()

	if !p.curTokenIs(token.RPAREN) {
		res, err := p.parseArgs()
		expression.Args = res
		if err != nil {
			return expression, err
		}
	}

	if err := p.expect(token.RPAREN); err != nil {
		return expression, err
	}

	return expression, nil
}

func (p *Parser) parseArgs() ([]ast.Node, error) {
	defer untrace(trace("args"))
	args := []ast.Node{}

	for !p.curTokenIs(token.RPAREN) {
		res, err := p.parseExpression(LOWEST)
		if err != nil {
			return args, err
		}

		args = append(args, res)

		if p.curTokenIs(token.COMMA) {
			p.nextToken()
		} else {
			break
		}
	}

	return args, nil
}

func (p *Parser) parseSlicesInfix(left ast.Node) (ast.Node, error) { // TODO: support [a:b:c]
	defer untrace(trace("slicesInfix"))
	n := &ast.SliceNode{Left: left}

	p.nextToken()

	res, err := p.parseExpression(LOWEST)
	n.Index = res
	if err != nil {
		return n, err
	}

	if err := p.expect(token.RBRACKET); err != nil {
		return n, err
	}

	return n, nil
}

func (p *Parser) nextToken() {
	p.curToken = p.peekToken

	p.pos++
	if p.pos >= len(p.tokens) {
		return
	}

	p.peekToken = p.tokens[p.pos]
}

func (p *Parser) setPos(index int) {
	p.pos = index
	p.curToken = p.tokens[p.pos-1]
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

func (p *Parser) isCompoundStatement() bool {
	_, ok := p.compundStatementFns[p.curToken.Type]
	return ok
}
