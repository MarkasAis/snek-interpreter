package main

import "snek/lexer"

func main() {
	code := `def foo():
    x = 42
    if x > 0:
        print("Hello")`

	lexer := lexer.NewLexer(code)
	lexer.PrintTokens()
}
