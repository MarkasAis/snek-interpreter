package main

import (
	"snek/lexer"
	"snek/parser"
)

func main() {
	// code := `def foo():
	// x = 42 # comments!!!

	// if x > 0:
	//     print("Hello")`

	code := `
pass
pass
pass
`

	lexer0 := lexer.NewLexer(code)
	lexer0.PrintTokens()

	lexer := lexer.NewLexer(code)
	// lexer.PrintTokens()

	parser := parser.New(lexer)
	parser.ParseProgram()
}
