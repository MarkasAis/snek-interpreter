package lexer

import (
	"regexp"
	"snek/token"
)

type Lexer struct {
	input       string
	pos         int
	indentStack []int // Tracks indentation levels
	startOfLine bool  // Tracks if we're at the start of a line
	tokens      []token.Token
	errors      []string
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
	{regexp.MustCompile(`^"([^"\\]*(\\.[^"\\]*)*)"`), token.STRING},
	{regexp.MustCompile(`^'([^'\\]*(\\.[^'\\]*)*)'`), token.STRING},
	{regexp.MustCompile(`^\s*\\\n`), token.IGNORE},
	{regexp.MustCompile(`^\n`), token.NEW_LINE},
	{regexp.MustCompile(`^\s+`), token.IGNORE},
	{regexp.MustCompile(`^#.*`), token.IGNORE},
	{regexp.MustCompile(`^\(`), token.BRACKET_OPEN},
	{regexp.MustCompile(`^\)`), token.BRACKET_CLOSE},
	{regexp.MustCompile(`^\[`), token.SQUARE_BRACKET_OPEN},
	{regexp.MustCompile(`^\]`), token.SQUARE_BRACKET_CLOSE},
	{regexp.MustCompile(`^\{`), token.CURL_BRACE_OPEN},
	{regexp.MustCompile(`^\}`), token.CURL_BRACE_CLOSE},
	{regexp.MustCompile(`^,`), token.COMMA},
	{regexp.MustCompile(`^:`), token.COLON},
	{regexp.MustCompile(`^;`), token.SEMICOLON},
	{regexp.MustCompile(`^\.`), token.DOT},
	{regexp.MustCompile(`^\S+`), token.UNKNOWN},
}

func New(input string) *Lexer {
	return &Lexer{
		input:       input,
		pos:         0,
		indentStack: []int{0},
		startOfLine: true,
	}
}

func (l *Lexer) Tokenize() []token.Token {
	l.tokens = []token.Token{}

	for l.pos < len(l.input) {
		l.tokenizeNext()
	}

	// Ensure all dedents are closed at EOF
	if len(l.indentStack) > 1 {
		l.tokens = append(l.tokens, token.Token{Type: token.NEW_LINE, Literal: "", Pos: l.pos})
	}

	for len(l.indentStack) > 1 {
		l.tokens = append(l.tokens, token.Token{Type: token.DEDENT, Literal: "", Pos: l.pos})
		l.indentStack = l.indentStack[:len(l.indentStack)-1]
	}

	l.tokens = append(l.tokens, token.Token{Type: token.EOF, Literal: "", Pos: l.pos})
	return l.tokens
}

func (l *Lexer) tokenizeNext() {
	input := l.input[l.pos:]

	// Handle indentation if we're at the start of a line
	if l.startOfLine {
		// Check if line is empty or only contains comment
		if match := regexp.MustCompile(`^\s*(#.*)?(\n|$)`).FindString(input); match != "" {
			l.pos += len(match)
			return
		}

		// Capture indentation
		indentation := regexp.MustCompile(`^[ \t]*`).FindString(input)
		indentLevel := len(indentation)

		// Check indentation changes
		lastIndent := l.indentStack[len(l.indentStack)-1]
		if indentLevel > lastIndent {
			l.indentStack = append(l.indentStack, indentLevel)
			l.tokens = append(l.tokens, token.Token{Type: token.INDENT, Literal: indentation, Pos: l.pos})
		} else if indentLevel < lastIndent {
			for len(l.indentStack) > 1 && indentLevel < l.indentStack[len(l.indentStack)-1] {
				l.indentStack = l.indentStack[:len(l.indentStack)-1]
				l.tokens = append(l.tokens, token.Token{Type: token.DEDENT, Literal: "", Pos: l.pos})
			}
			if l.indentStack[len(l.indentStack)-1] != indentLevel {
				l.errors = append(l.errors, "unindent does not match any outer indentation level")
			}
		}

		l.pos += indentLevel
		l.startOfLine = false
		return
	}

	// Match token patterns
	for _, pattern := range tokenPatterns {
		if match := pattern.regex.FindString(input); match != "" {
			tokenLength := len(match)

			// Skip ignored tokens
			if pattern.tType == token.IGNORE {
				l.pos += tokenLength
				return
			}

			l.tokens = append(l.tokens, token.Token{
				Type:    pattern.tType,
				Literal: match,
				Pos:     l.pos,
			})

			if pattern.tType == token.NEW_LINE {
				l.startOfLine = true
			}

			l.pos += tokenLength
			return
		}
	}

	// Unknown token handling
	l.tokens = append(l.tokens, token.Token{
		Type:    token.UNKNOWN,
		Literal: string(input[0]),
		Pos:     l.pos,
	})
	l.pos++
}

func (p *Lexer) Errors() []string {
	return p.errors
}
