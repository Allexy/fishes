package tokenizer

import "fmt"

type TokenType uint8

// Possible token types
const (
	TokenBOF              TokenType = iota
	TokenDefault                    // ...
	TokenWord                       // abc
	TokenNumber                     // 123.0
	TokenString                     // "..."
	TokenLogic                      // true / false
	TokenOperator                   // +-*/...
	TokenOpenParen                  // (
	TokenCloseParen                 // )
	TokenOpenBracket                // [
	TokenCloseBracket               // ]
	TokenOpenBrace                  // {
	TokenCloseBrace                 // }
	TokenColon                      // :
	TokenSemicolon                  // ;
	TokenComa                       // ,
	TokenPoint                      // .
	TokenAt                         // @
	TokenAssignment                 // =
	TokenArrow                      // =>
	TokenVariable                   // $varname
	TokenComment                    // #.... or //...
	TokenMultilineComment           // /*...*/
	TokenWhiteSpace                 // any white space
	TokenEOF
)

type Token struct {
	Token      TokenType
	Text       string
	SourceName string
	Line, Col  uint32
}

func (t Token) String() string {
	var name string
	switch t.Token {
	case TokenBOF:
		name = "BOF"
	case TokenDefault:
		name = "DEFAULT"
	case TokenWord:
		name = "WORD"
	case TokenString:
		name = "STRING"
	case TokenNumber:
		name = "NUMBER"
	case TokenLogic:
		name = "LOGIC"
	case TokenOperator:
		name = "OPERATOR"
	case TokenOpenParen:
		name = "O_PAREN"
	case TokenCloseParen:
		name = "C_PAREN"
	case TokenOpenBracket:
		name = "O_BRACKET"
	case TokenCloseBracket:
		name = "C_BRACKET"
	case TokenOpenBrace:
		name = "O_BRACE"
	case TokenCloseBrace:
		name = "C_BRACE"
	case TokenColon:
		name = "COLON"
	case TokenSemicolon:
		name = "SEMICOLON"
	case TokenComa:
		name = "COMA"
	case TokenAt:
		name = "AT"
	case TokenPoint:
		name = "POINT"
	case TokenAssignment:
		name = "ASSIGNMENT"
	case TokenArrow:
		name = "ARROW"
	case TokenVariable:
		name = "VARIABLE"
	case TokenComment:
		name = "COMMENT"
	case TokenMultilineComment:
		name = "MULTILINE_COMMET"
	case TokenWhiteSpace:
		name = "WHITE_SPACE"
	case TokenEOF:
		name = "EOF"
	default:
		name = "UNKNOWN"
	}
	return fmt.Sprintf("%s(%q@%d:%d)", name, t.Text, t.Line, t.Col)
}
