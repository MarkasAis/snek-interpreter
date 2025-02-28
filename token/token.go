package token

type TokenType int

const (
	UNKNOWN TokenType = iota
	IGNORE
	NUMBER
	IDENTIFIER
	STRING
	DEF
	LPAREN
	RPAREN
	LBRACKET
	RBRACKET
	LBRACE
	RBRACE
	COMMA
	COLON
	SEMICOLON
	DOT
	NEW_LINE
	ASSIGN
	IF
	ELSE
	ELIF
	FOR
	IN
	WHILE
	OR
	AND
	NOT
	COMPARE
	SUM
	PRODUCT
	EXP
	PASS
	RETURN
	BREAK
	CONTINUE
	GLOBAL
	IMPORT
	FROM
	INDENT
	DEDENT
	EOF
)

type Token struct {
	Type    TokenType
	Literal string
	Pos     int
}

func (t TokenType) String() string {
	switch t {
	case NUMBER:
		return "NUMBER"
	case DEF:
		return "DEF"
	case IF:
		return "IF"
	case ELSE:
		return "ELSE"
	case ELIF:
		return "ELIF"
	case FOR:
		return "FOR"
	case IN:
		return "IN"
	case WHILE:
		return "WHILE"
	case OR:
		return "OR"
	case AND:
		return "AND"
	case NOT:
		return "NOT"
	case PASS:
		return "PASS"
	case BREAK:
		return "BREAK"
	case CONTINUE:
		return "CONTINUE"
	case RETURN:
		return "RETURN"
	case GLOBAL:
		return "GLOBAL"
	case IMPORT:
		return "IMPORT"
	case FROM:
		return "FROM"
	case COMPARE:
		return "COMPARE"
	case ASSIGN:
		return "ASSIGN"
	case SUM:
		return "SUM"
	case EXP:
		return "EXP"
	case PRODUCT:
		return "PRODUCT"
	case IDENTIFIER:
		return "IDENTIFIER"
	case STRING:
		return "STRING"
	case NEW_LINE:
		return "NEW_LINE"
	case IGNORE:
		return "IGNORE"
	case LPAREN:
		return "LPAREN"
	case RPAREN:
		return "RPAREN"
	case LBRACKET:
		return "LBRACKET"
	case RBRACKET:
		return "RBRACKET"
	case LBRACE:
		return "LBRACE"
	case RBRACE:
		return "RBRACE"
	case COMMA:
		return "COMMA"
	case COLON:
		return "COLON"
	case SEMICOLON:
		return "SEMICOLON"
	case DOT:
		return "DOT"
	case UNKNOWN:
		return "UNKNOWN"
	case INDENT:
		return "INDENT"
	case DEDENT:
		return "DEDENT"
	case EOF:
		return "EOF"
	default:
		return "UNKNOWN_TYPE"
	}
}
