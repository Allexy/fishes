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
		switch token.Type {
		case TokenDefault:
			return nil, NewTokenizerError(token.SourceName, fmt.Sprintf("invalid token %v", token), token.Line, token.Col, nil)
		case TokenNumber:
			filterNumbersText(token)
			if err := isInvalidNumber(token); err != nil {
				return nil, NewTokenizerError(token.SourceName, fmt.Sprintf("invalid numerical literal %q", token.Text), token.Line, token.Col, err)
			}
		case TokenOperator:
			if isInvalidOperator(token, previous, next) {
				return nil, NewTokenizerError(token.SourceName, fmt.Sprintf("invalid operator %q", token.Text), token.Line, token.Col, nil)
			}
			switch token.Text {
			case lang.OpArrow:
				token.Type = TokenArrow
			case lang.OpAssign:
				token.Type = TokenAssignment
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
							Type:       TokenNumber,
							Text:       newText,
							SourceName: token.SourceName,
							Line:       token.Line,
							Col:        token.Col,
						}
						if err := isInvalidNumber(&replacement); err != nil {
							return nil, NewTokenizerError(token.SourceName, fmt.Sprintf("invalid numerical literal %q", token.Text), token.Line, token.Col, err)
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
				token.Type = TokenLogic
			}
		}
		optimized = append(optimized, *token)
		tw.Move(1)
	}

	tw.Clear()

	return NewTokenWalker(optimized), nil
}

func filterNumbersText(t *Token) {
	for t.Text[0] == '0' && len(t.Text) > 1 {
		t.Text = t.Text[1:len(t.Text)]
	}
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
	switch previous.Type {
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
		if p == nil || p.Type != TokenVariable {
			return false
		}
	case lang.OpIncrement, lang.OpDecrement:
		// This kinds stands before or after variable
		if (p == nil || p.Type != TokenVariable) && (n == nil || n.Type != TokenVariable) {
			return false
		}
	}
	return true
}

func isInvalidNumber(c *Token) error {
	v, err := strconv.ParseFloat(c.Text, 64)
	if err == nil {
		c.ParsedNumber = v
	}
	return err
}
