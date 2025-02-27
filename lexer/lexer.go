package lexer

import (
	"fmt"
	"regexp"
	"snek/token"
)

type Lexer struct {
	input      string
	pos        int  // Current position in input
	line       int  // Current line number
	column     int  // Current column number
	firstToken bool // Flag to track if we are on the first token
}

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

func NewLexer(input string) *Lexer {
	return &Lexer{input: input, pos: 0, line: 0, column: 0, firstToken: true}
}

// NextToken retrieves the next token lazily with **correct `Pos` tracking**
func (l *Lexer) NextToken() token.Token {
	// If we reached the end of input, return EOF
	if l.pos >= len(l.input) {
		return token.Token{Type: token.EOF, Literal: "", Pos: l.pos, Line: l.line, Column: l.column}
	}

	input := l.input[l.pos:]

	startPos := l.pos
	startLine := l.line
	startColumn := l.column

	// Ensure the lexer always starts with a NEW_LINE token (with indentation if present)
	if l.firstToken {
		l.firstToken = false
		initialIndent := regexp.MustCompile(`^[ \t]*`).FindString(input)
		startPos := l.pos
		l.pos += len(initialIndent)
		l.column = len(initialIndent)
		return token.Token{
			Type:    token.NEW_LINE,
			Literal: initialIndent,
			Pos:     startPos,
			Line:    startLine,
			Column:  startColumn,
		}
	}

	for _, pattern := range tokenPatterns {
		if match := pattern.regex.FindString(input); match != "" {
			tokenLength := len(match)

			// Skip ignored tokens like whitespace and comments
			if pattern.tType == token.IGNORE {
				l.pos += tokenLength
				l.column += tokenLength
				return l.NextToken()
			}

			// Handle new lines
			if pattern.tType == token.NEW_LINE {
				l.line++
				l.column = tokenLength - 1
				l.pos += tokenLength
				return token.Token{
					Type:    pattern.tType,
					Literal: match,
					Pos:     startPos,
					Line:    startLine,
					Column:  startColumn,
				}
			}

			l.pos += tokenLength
			l.column += tokenLength

			return token.Token{
				Type:    pattern.tType,
				Literal: match,
				Pos:     startPos,
				Line:    startLine,
				Column:  startColumn,
			}
		}
	}

	// If no match, return an UNKNOWN token
	startPos = l.pos
	l.pos++
	l.column++
	return token.Token{
		Type:    token.UNKNOWN,
		Literal: string(input[0]),
		Pos:     startPos,
		Line:    l.line,
		Column:  startColumn,
	}
}

func (l *Lexer) PrintTokens() {
	fmt.Println("Pos  | Line | Col  | Type                 | Literal")
	fmt.Println("-----------------------------------------------------")

	for {
		tok := l.NextToken()
		fmt.Printf("%-4d | %-4d | %-4d | %-20s | %q\n",
			tok.Pos, tok.Line, tok.Column, tok.Type, tok.Literal)
		if tok.Type == token.EOF {
			break
		}
	}
}
