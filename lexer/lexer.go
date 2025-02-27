package lexer

import (
	"fmt"
	"regexp"
	"snek/token"
)

// Lexer struct for lexing input lazily
type Lexer struct {
	input  string
	pos    int // Current position in input
	line   int // Current line number
	column int // Current column number
}

// Token regex patterns
var tokenPatterns = []struct {
	regex *regexp.Regexp
	tType token.TokenType
}{
	{regexp.MustCompile(`^(\d*\.)?\d+\b`), token.NUMBER},
	{regexp.MustCompile(`^def\b`), token.DEF},
	{regexp.MustCompile(`^if\b`), token.IF},
	{regexp.MustCompile(`^else\b`), token.ELSE},
	{regexp.MustCompile(`^elif\b`), token.ELIF},
	{regexp.MustCompile(`^for\b`), token.FOR},
	{regexp.MustCompile(`^(?:in|not\s+in)\b`), token.IN},
	{regexp.MustCompile(`^while\b`), token.WHILE},
	{regexp.MustCompile(`^or\b`), token.OR},
	{regexp.MustCompile(`^and\b`), token.AND},
	{regexp.MustCompile(`^not\b`), token.NOT},
	{regexp.MustCompile(`^pass\b`), token.PASS},
	{regexp.MustCompile(`^break\b`), token.BREAK},
	{regexp.MustCompile(`^continue\b`), token.CONTINUE},
	{regexp.MustCompile(`^return\b`), token.RETURN},
	{regexp.MustCompile(`^global\b`), token.GLOBAL},
	{regexp.MustCompile(`^import\b`), token.IMPORT},
	{regexp.MustCompile(`^from\b`), token.FROM},
	{regexp.MustCompile(`^(==|!=|>=|>|<=|<)`), token.COMPARE},
	{regexp.MustCompile(`^(=|\+=|-=|\*=|/=|//=|%=)`), token.ASSIGN},
	{regexp.MustCompile(`^[+-]`), token.ADD},
	{regexp.MustCompile(`^\*\*`), token.EXP},
	{regexp.MustCompile(`^//`), token.MULT},
	{regexp.MustCompile(`^[*/%]`), token.MULT},
	{regexp.MustCompile(`^[a-zA-Z_]\w*`), token.IDENTIFIER},
	{regexp.MustCompile(`^"([^"\\]*(\\.[^"\\]*)*)"`), token.STRING}, // Matches double-quoted strings
	{regexp.MustCompile(`^'([^'\\]*(\\.[^'\\]*)*)'`), token.STRING}, // Matches single-quoted strings
	{regexp.MustCompile(`^\n[ \t]*`), token.NEW_LINE},               // Captures indentation after a newline
	{regexp.MustCompile(`^\s+`), token.IGNORE},                      // Ignore whitespace
	{regexp.MustCompile(`^#.*`), token.IGNORE},                      // Ignore comments
	{regexp.MustCompile(`^\(`), token.BRACKET_OPEN},
	{regexp.MustCompile(`^\)`), token.BRACKET_CLOSE},
	{regexp.MustCompile(`^\[`), token.SQUARE_BRACKET_OPEN},
	{regexp.MustCompile(`^\]`), token.SQUARE_BRACKET_CLOSE},
	{regexp.MustCompile(`^\{`), token.CURL_BRACE_OPEN},
	{regexp.MustCompile(`^\}`), token.CURL_BRACE_CLOSE},
	{regexp.MustCompile(`^,`), token.COMMA},
	{regexp.MustCompile(`^:`), token.COLON},
	{regexp.MustCompile(`^\.`), token.DOT},
	{regexp.MustCompile(`^\S+`), token.UNKNOWN},
}

// NewLexer initializes a lazy lexer
func NewLexer(input string) *Lexer {
	return &Lexer{input: input, pos: 0, line: 1, column: 1}
}

// NextToken retrieves the next token lazily with line/column tracking
func (l *Lexer) NextToken() token.Token {
	// If we reached the end of input, return EOF
	if l.pos >= len(l.input) {
		return token.Token{Type: token.EOF, Literal: "", Pos: l.pos, Line: l.line, Column: l.column - 1}
	}

	input := l.input[l.pos:]
	startColumn := l.column // Store column before reading token

	// Handle initial indentation on first token
	if l.pos == 0 {
		initialIndent := regexp.MustCompile(`^[ \t]+`).FindString(input)
		if initialIndent != "" {
			l.pos += len(initialIndent)
			return token.Token{Type: token.NEW_LINE, Literal: "\n" + initialIndent, Pos: l.pos, Line: l.line, Column: 1}
		}
	}

	// Loop through token patterns to find a match
	for _, pattern := range tokenPatterns {
		if match := pattern.regex.FindString(input); match != "" {
			l.pos += len(match)

			// Skip ignored tokens like whitespace and comments
			if pattern.tType == token.IGNORE {
				return l.NextToken()
			}

			// Handle new lines properly
			if pattern.tType == token.NEW_LINE {
				l.line++
				l.column = 1
			} else {
				l.column += len(match)
			}

			return token.Token{
				Type:    pattern.tType,
				Literal: match,
				Pos:     l.pos,
				Line:    l.line,
				Column:  startColumn,
			}
		}
	}

	// If no match, return an UNKNOWN token
	l.pos++
	return token.Token{
		Type:    token.UNKNOWN,
		Literal: string(input[0]),
		Pos:     l.pos - 1,
		Line:    l.line,
		Column:  startColumn,
	}
}

// Debug function to print tokens lazily with line/column info
func (l *Lexer) PrintTokens() {
	for {
		tok := l.NextToken()
		fmt.Printf("Pos: %-3d | Line: %-3d | Col: %-3d | %-20s | %q\n",
			tok.Pos, tok.Line, tok.Column, tok.Type, tok.Literal)
		if tok.Type == token.EOF {
			break
		}
	}
}
