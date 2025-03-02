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
	code := `
x[1+2][2] -= 3`

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
	ast, err := p.ParseFile()

	fmt.Println("----------")

	if err != nil {
		io.WriteString(os.Stdout, "Parse Error: "+err.Error()+"\n")
		return
	}

	fmt.Println(ast.String())
}

func PrintTokens(tokens []token.Token) {
	fmt.Println("Pos  | Type                 | Literal")
	fmt.Println("-----------------------------------------------------")

	for _, tok := range tokens {
		fmt.Printf("%-4d | %-20s | %q\n",
			tok.Pos, tok.Type, tok.Literal)
	}
}
