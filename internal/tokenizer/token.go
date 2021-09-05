package tokenizer

import "fmt"

type TokenType uint8

// Possible token types
const (
	TokenBOF              TokenType = iota
	TokenAny                        // ...
	TokenWord                       // abc
	TokenConst                      // const keyword
	TokenNull                       // null keyword
	TokenIf                         // if keyword
	TokenLogicTrue                  // true keyword
	TokenLogicFalse                 // false keyword
	TokenTry                        // try keyword
	TokenCatch                      // catch keyword
	TokenThrow                      // throw keyword
	TokenFor                        // for keyword
	TokenWhile                      // while keyword
	TokenSwitch                     // switch keyword
	TokenCase                       // case keyword
	TokenReturn                     // return keyword
	TokenNumber                     // 123.0
	TokenString                     // "..."
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
		return "T_BOF"
	case TokenAny:
		return "T_ANY"
	case TokenWord:
		return "T_WORD"
	case TokenConst:
		return "T_CONST"
	case TokenNull:
		return "T_NULL"
	case TokenIf:
		return "T_IF"
	case TokenLogicTrue:
		return "T_TRUE"
	case TokenLogicFalse:
		return "T_FALSE"
	case TokenTry:
		return "T_TRY"
	case TokenCatch:
		return "T_CATCH"
	case TokenThrow:
		return "T_THROW"
	case TokenFor:
		return "T_FOR"
	case TokenWhile:
		return "T_WHILE"
	case TokenSwitch:
		return "T_SWITCH"
	case TokenCase:
		return "T_CASE"
	case TokenReturn:
		return "T_RETURN"
	case TokenString:
		return "T_STRING"
	case TokenNumber:
		return "T_NUMBER"
	case TokenOperator:
		return "T_OPERATOR"
	case TokenOpenParen:
		return "T_OPENING_PAREN"
	case TokenCloseParen:
		return "T_CLOSING_PAREN"
	case TokenOpenBracket:
		return "T_OPENING_BRACKET"
	case TokenCloseBracket:
		return "T_CLOSING_BRACKET"
	case TokenOpenBrace:
		return "T_OPENING_BRACE"
	case TokenCloseBrace:
		return "T_CLOSING_BRACE"
	case TokenColon:
		return "T_COLON"
	case TokenSemicolon:
		return "T_SEMICOLON"
	case TokenComa:
		return "T_COMA"
	case TokenAt:
		return "T_AT"
	case TokenPoint:
		return "T_POINT"
	case TokenAssignment:
		return "T_ASSIGNMENT"
	case TokenArrow:
		return "T_ARROW"
	case TokenVariable:
		return "T_VARIABLE"
	case TokenComment:
		return "T_COMMENT"
	case TokenMultilineComment:
		return "T_MULTILINE_COMMENT"
	case TokenWhiteSpace:
		return "T_WHITE_SPACE"
	case TokenEOF:
		return "T_EOF"
	default:
		panic("unknown token type")
	}
}

type Token struct {
	Type         TokenType
	Text         string
	SourceName   string
	Line, Col    uint32
	ParsedNumber float64
}

func (t Token) String() string {
	return fmt.Sprintf("%s(%q@%d:%d)", t.Type, t.Text, t.Line, t.Col)
}
