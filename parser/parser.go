package parser

import (
	"fmt"
	"snek/ast"
	"snek/token"
)

type Parser struct {
	tokens    []token.Token
	pos       int
	curToken  token.Token
	peekToken token.Token
	errors    []string
}

func New(tokens []token.Token) *Parser {
	p := &Parser{
		tokens: tokens,
		pos:    -1,
		errors: []string{},
	}

	p.nextToken()
	p.nextToken()
	return p
}

func (p *Parser) Parse() ast.Node {
	return p.parseAssignmentStmt()
}

func (p *Parser) parseAssignmentStmt() ast.Node {
	// Expect an identifier (the left-hand side target)
	if !p.curTokenIs(token.IDENTIFIER) {
		p.curError(token.IDENTIFIER)
		return nil
	}

	// Create a node for the left-hand side
	target := &ast.IdentifierNode{Name: p.curToken.Literal}
	p.nextToken() // Consume the identifier

	// Ensure we have an "=" token
	if !p.expect(token.ASSIGN) {
		return nil
	}

	// Parse the right-hand side (expression)
	value := p.parseExpression()
	if value == nil {
		return nil
	}

	// Return an assignment node
	return &ast.AssignmentNode{Target: target, Value: value}
}

// parseExpression is a stub for now (expand later)
func (p *Parser) parseExpression() ast.Node {
	// For now, just return an identifier or number
	if p.curTokenIs(token.IDENTIFIER) {
		node := &ast.IdentifierNode{Name: p.curToken.Literal}
		p.nextToken()
		return node
	} else if p.curTokenIs(token.NUMBER) {
		node := &ast.NumberNode{Value: p.curToken.Literal}
		p.nextToken()
		return node
	}

	// Unexpected token in expression
	p.curError(token.IDENTIFIER) // Expecting an expression
	return nil
}

// func (p *Parser) parseSuite() ast.Node {
// 	defer untrace(trace("suite"))
// 	// Case 1: Single-line suite → stmt_list NEWLINE
// 	if !p.curTokenIs(token.NEW_LINE) {
// 		stmtList := p.parseStmtList()
// 		if !p.expect(token.NEW_LINE) {
// 			return nil
// 		}
// 		return stmtList
// 	}

// 	// Case 2: Multi-line suite → NEWLINE INDENT statement+ DEDENT
// 	p.nextToken() // Consume NEW_LINE

// 	// Check for INDENT (or return nil for empty suite)
// 	if !p.curTokenIs(token.INDENT) {
// 		return nil
// 	}
// 	p.nextToken() // Consume INDENT

// 	var statements []ast.Node

// 	// Parse statements inside the block
// 	for !p.curTokenIs(token.DEDENT) && !p.curTokenIs(token.EOF) {
// 		// Skip blank lines and comments
// 		// if p.curTokenIs(token.NEW_LINE) || p.curTokenIs(token.COMMENT) {
// 		// 	p.nextToken()
// 		// 	continue
// 		// }

// 		// Parse a statement
// 		stmt := p.parseStatement()
// 		if stmt != nil {
// 			statements = append(statements, stmt)
// 		}
// 	}

// 	// Ensure `DEDENT` appears at the end of the block
// 	if !p.expect(token.DEDENT) {
// 		return nil
// 	}

// 	// Return a placeholder node for now (replace with an actual AST node)
// 	return nil //&ast.BlockNode{Statements: statements}
// }

// // stmt_list ::= simple_stmt (";" simple_stmt)* [";"]
// func (p *Parser) parseStmtList() ast.Node {
// 	defer untrace(trace("statementList"))
// 	var statements []ast.Node

// 	stmt := p.parseSimpleStmt()
// 	if stmt != nil {
// 		statements = append(statements, stmt)
// 	}

// 	// Handle optional semicolon-separated simple statements
// 	for p.curTokenIs(token.SEMICOLON) {
// 		p.nextToken() // Consume `;`
// 		if p.curTokenIs(token.NEW_LINE) {
// 			break // Allow trailing semicolon
// 		}
// 		stmt := p.parseSimpleStmt()
// 		if stmt != nil {
// 			statements = append(statements, stmt)
// 		}
// 	}

// 	return nil // &node.StatementListNode{Statements: statements}
// }

// // statement ::= stmt_list NEWLINE | compound_stmt
// func (p *Parser) parseStatement() ast.Node {
// 	defer untrace(trace("statement"))
// 	if p.isCompoundStmt() {
// 		return p.parseCompoundStmt()
// 	}

// 	stmtList := p.parseStmtList()
// 	if !p.expect(token.NEW_LINE) {
// 		return nil
// 	}
// 	return stmtList
// }

// func (p *Parser) parseSimpleStmt() ast.Node {
// 	defer untrace(trace("simpleStatement"))
// 	// For now, just consume a token and return a stub node
// 	defer p.nextToken()
// 	return nil //&node.SimpleStmtNode{}
// }

// // Placeholder: Parses `compound_stmt`
// func (p *Parser) parseCompoundStmt() ast.Node {
// 	defer untrace(trace("compundStatement"))
// 	// For now, just consume a token and return a stub node
// 	defer p.nextToken()
// 	return nil //&node.CompoundStmtNode{}
// }

// // Determines if current token starts a compound statement
// func (p *Parser) isCompoundStmt() bool {
// 	switch p.curToken.Type {
// 	case token.IF, token.FOR, token.WHILE, token.DEF:
// 		return true
// 	}
// 	return false
// }

// func (p *Parser) parseSuite() ast.Node {
// 	statements := []ast.Node{}

// 	for !p.curTokenIs(token.DEDENT) && !p.curTokenIs(token.EOF) {
// 		if p.curTokenIs(token.INDENT) {
// 			p.errors = append(p.errors, "too much indentation")
// 			return nil
// 		}

// 		stmt := p.parseStatement()
// 		statements = append(statements, stmt)
// 	}

// 	if len(statements) == 0 {
// 		p.errors = append(p.errors, "no statements")
// 		return nil
// 	}

// 	return nil
// }

// func (p *Parser) parseStatement() ast.Node {
// 	switch p.curToken.Type {
// 	case token.IF:
// 		node := &ast.IfStatement{}

// 		for !p.curTokenIs(token.COLON) {
// 			p.nextToken()
// 		}
// 		p.nextToken()
// 		node.Consequence = p.parseBlock(indentation)

// 		if p.peekTokenIs(token.ELIF) {

// 		}
// 	}

// 	for !p.curTokenIs(token.NEW_LINE) && !p.curTokenIs(token.EOF) {
// 		p.nextToken()
// 	}
// 	return 1
// }

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

func (p *Parser) peekTokenIs(t token.TokenType) bool {
	return p.peekToken.Type == t
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
