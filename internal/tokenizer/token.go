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

func (tt TokenType) String() string {
	switch tt {
	case TokenBOF:
		return "BOF"
	case TokenDefault:
		return "DEFAULT"
	case TokenWord:
		return "WORD"
	case TokenString:
		return "STRING"
	case TokenNumber:
		return "NUMBER"
	case TokenLogic:
		return "LOGIC"
	case TokenOperator:
		return "OPERATOR"
	case TokenOpenParen:
		return "O_PAREN"
	case TokenCloseParen:
		return "C_PAREN"
	case TokenOpenBracket:
		return "O_BRACKET"
	case TokenCloseBracket:
		return "C_BRACKET"
	case TokenOpenBrace:
		return "O_BRACE"
	case TokenCloseBrace:
		return "C_BRACE"
	case TokenColon:
		return "COLON"
	case TokenSemicolon:
		return "SEMICOLON"
	case TokenComa:
		return "COMA"
	case TokenAt:
		return "AT"
	case TokenPoint:
		return "POINT"
	case TokenAssignment:
		return "ASSIGNMENT"
	case TokenArrow:
		return "ARROW"
	case TokenVariable:
		return "VARIABLE"
	case TokenComment:
		return "COMMENT"
	case TokenMultilineComment:
		return "MULTILINE_COMMET"
	case TokenWhiteSpace:
		return "WHITE_SPACE"
	case TokenEOF:
		return "EOF"
	default:
		panic("unknown token type")
	}
}

type Token struct {
	Token      TokenType
	Text       string
	SourceName string
	Line, Col  uint32
}

func (t Token) String() string {
	return fmt.Sprintf("%s(%q@%d:%d)", t.Token, t.Text, t.Line, t.Col)
}
