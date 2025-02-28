package main

import (
	"fmt"
	"io"
	"os"
	"snek/lexer"
	"snek/parser"
	"snek/token"
)

func main() {
	// code := `def foo():
	// x = 42 # comments!!!

	// if x > 0:
	//     print("Hello")`

	code := `
x = 10 % 20
y = 3`

	l := lexer.New(code)
	tokens := l.Tokenize()
	PrintTokens(tokens)

	if len(l.Errors()) > 0 {
		io.WriteString(os.Stdout, "Lexer Errors:\n")
		for _, msg := range l.Errors() {
			io.WriteString(os.Stdout, "\t- "+msg+"\n")
		}
		return
	}

	fmt.Println("----------")

	p := parser.New(tokens)
	ast := p.Parse()

	fmt.Println("----------")
	fmt.Println(ast.String())

	if len(p.Errors()) > 0 {
		io.WriteString(os.Stdout, "Parser Errors:\n")
		for _, msg := range p.Errors() {
			io.WriteString(os.Stdout, "\t- "+msg+"\n")
		}
		return
	}
}

func PrintTokens(tokens []token.Token) {
	fmt.Println("Pos  | Type                 | Literal")
	fmt.Println("-----------------------------------------------------")

	for _, tok := range tokens {
		fmt.Printf("%-4d | %-20s | %q\n",
			tok.Pos, tok.Type, tok.Literal)
	}
}
