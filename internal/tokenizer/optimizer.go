package tokenizer

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/Allexy/fishes/internal/lang"
)

func optimizeAndValidate(tw TokenWalker) (TokenWalker, error) {

	if tw.Size() < 3 {
		return nil, errors.New("too few tokens in walker")
	}

	optimized := make([]Token, 0, 1024)

	for tw.Next() {
		token := tw.Get(0)
		previous := tw.Get(-1)
		next := tw.Get(1)
		switch token.Token {
		case TokenDefault:
			return nil, NewTokenizerError(token.SourceName, fmt.Sprintf("Invalid token %v", token), token.Line, token.Col, nil)
		case TokenNumber:
			filterNumbersText(token)
			if isInvalidNumber(token) {
				return nil, NewTokenizerError(token.SourceName, fmt.Sprintf("Invalid numerical literal %q", token.Text), token.Line, token.Col, nil)
			}
		case TokenOperator:
			if isInvalidOperator(token, previous, next) {
				return nil, NewTokenizerError(token.SourceName, fmt.Sprintf("Invalid operator %q", token.Text), token.Line, token.Col, nil)
			}
			switch token.Text {
			case lang.OpArrow:
				token.Token = TokenArrow
			case lang.OpAssign:
				token.Token = TokenAssignment
			case lang.OpPlus, lang.OpMinus:
				// need to check if next token is numerical literal, it may be negative number
				if tw.Match(TokenOperator, TokenNumber) {
					if isNegativeNumberDetected(previous) {
						filterNumbersText(next)
						var newText string
						if token.Text == lang.OpPlus {
							newText = next.Text
						} else {
							newText = token.Text + next.Text
						}
						replacement := Token{
							Token:      TokenNumber,
							Text:       newText,
							SourceName: token.SourceName,
							Line:       token.Line,
							Col:        token.Col,
						}
						if isInvalidNumber(&replacement) {
							return nil, NewTokenizerError(token.SourceName, fmt.Sprintf("Invalid numerical literal %q", token.Text), token.Line, token.Col, nil)
						}
						optimized = append(optimized, replacement)
						tw.Move(2) // Step over operator and up comming number
						continue
					}
				}
			default:
			}
		case TokenWord:
			switch token.Text {
			case lang.KwTrue, lang.KwFalse:
				token.Token = TokenLogic
			}
		}
		optimized = append(optimized, *token)
		tw.Move(1)
	}

	tw.Clear()

	return NewTokenWalker(optimized), nil
}

func filterNumbersText(t *Token) {
	if t.Text[0] == '.' {
		t.Text = "0" + t.Text
	} else if t.Text[len(t.Text)-1] == '.' {
		t.Text += "0"
	}
}

func isNegativeNumberDetected(previous *Token) bool {
	if previous == nil {
		return true
	}
	switch previous.Token {
	// means that current token is part of arithmetic expression
	case TokenNumber, TokenVariable, TokenCloseParen:
		return false
	}
	return true
}

func isInvalidOperator(c *Token, p *Token, n *Token) bool {
	switch c.Text {
	case lang.OpArrow, lang.OpPlus, lang.OpMinus, lang.OpDivision, lang.OpMultiply, lang.OpModulo, lang.OpAssign, lang.OpEquals,
		lang.OpNotEquals, lang.OpGreaterThan, lang.OpGreaterThanOrEquals, lang.OpLesserThan, lang.OpLesserThanOrEquals, lang.OpNot,
		lang.OpAnd, lang.OpOr, lang.OpIncrement, lang.OpDecrement, lang.OpPlusAssign, lang.OpMinusAssign, lang.OpDivideAssign,
		lang.OpMultiplyAssign, lang.OpModuloAssign:
		return false
	}
	switch c.Text {
	case lang.OpPlusAssign, lang.OpMinusAssign, lang.OpDivideAssign, lang.OpMultiplyAssign, lang.OpModuloAssign:
		// These kinds of operators must follow by variable
		if p == nil || p.Token != TokenVariable {
			return false
		}
	case lang.OpIncrement, lang.OpDecrement:
		// This kinds stands before or after variable
		if (p == nil || p.Token != TokenVariable) && (n == nil || n.Token != TokenVariable) {
			return false
		}
	}
	return true
}

func isInvalidNumber(c *Token) bool {
	_, ok := strconv.ParseFloat(c.Text, 64)
	return ok != nil
}
