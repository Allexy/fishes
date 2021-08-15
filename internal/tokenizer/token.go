package tokenizer

import "fmt"

type TokenType uint8

type Token struct {
	Token      TokenType
	Text       string
	SourceName string
	Line, Col  uint32
}

func (t Token) String() string {
	var name string
	switch t.Token {
	case TT_BOF:
		name = "BOF"
	case TT_DEFAULT:
		name = "DEFAULT"
	case TT_WORD:
		name = "WORD"
	case TT_STRING:
		name = "STRING"
	case TT_NUMBER:
		name = "NUMBER"
	case TT_LOGIC:
		name = "LOGIC"
	case TT_OPERATOR:
		name = "OPERATOR"
	case TT_O_PAREN:
		name = "O_PAREN"
	case TT_C_PAREN:
		name = "C_PAREN"
	case TT_O_BRACKET:
		name = "O_BRACKET"
	case TT_C_BRACKET:
		name = "C_BRACKET"
	case TT_O_BRACE:
		name = "O_BRACE"
	case TT_C_BRACE:
		name = "C_BRACE"
	case TT_COLON:
		name = "COLON"
	case TT_SEMICOLON:
		name = "SEMICOLON"
	case TT_COMA:
		name = "COMA"
	case TT_AT:
		name = "AT"
	case TT_POINT:
		name = "POINT"
	case TT_ASSIGNMENT:
		name = "ASSIGNMENT"
	case TT_ARROW:
		name = "ARROW"
	case TT_VARIABLE:
		name = "VARIABLE"
	case TT_COMMENT:
		name = "COMMENT"
	case TT_MULTILINE_COMMENT:
		name = "MULTILINE_COMMET"
	case TT_WHITE_SPACE:
		name = "WHITE_SPACE"
	case TT_EOF:
		name = "EOF"
	default:
		name = "UNKNOWN"
	}
	return fmt.Sprintf("%s(%q@%d:%d)", name, t.Text, t.Line, t.Col)
}
