package parser

import (
	"snek/lexer"
	"snek/token"
)

type Parser struct {
	l         *lexer.Lexer
	curToken  token.Token
	peekToken token.Token
	errors    []string
}

func (p *Parser) ParseProgram() {

}
