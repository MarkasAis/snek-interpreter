package main

import (
	"io"
	"os"
	"snek/lexer"
)

func main() {
	// code := `def foo():
	// x = 42 # comments!!!

	// if x > 0:
	//     print("Hello")`

	code := `
def foo(x):
	return x2 + \
	1

x = 10

# what
y = 10`

	lexer := lexer.NewLexer(code)
	lexer.Tokenize()
	lexer.PrintTokens()

	for _, msg := range lexer.Errors() {
		io.WriteString(os.Stdout, "Lexer error: "+msg+"\n")
	}

	// if match := regexp.MustCompile(`^\s*(#.*)?(\n|$)`).FindString("  # comment\ntest\n"); match != "" {
	// 	fmt.Print(match)
	// }

	// lexer := lexer.NewLexer(code)
	// lexer.PrintTokens()

	// parser := parser.New(lexer)
	// parser.ParseProgram()
}
