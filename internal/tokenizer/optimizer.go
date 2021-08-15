package tokenizer

import (
	"errors"
	"fmt"
	"strconv"
)

func optimizeAndValidate(tw TokenWalker) (TokenWalker, error) {

	if tw.Size() < 3 {
		return nil, errors.New("too few tokens in walker")
	}

	optimized := make([]Token, 0, 1024)

	for tw.HasNext() {
		token := tw.Get(0)
		previous := tw.Get(-1)
		next := tw.Get(1)
		switch token.Token {
		case TT_DEFAULT:
			return nil, NewTokenizerError(token.SourceName, fmt.Sprintf("Invalid token %v", token), token.Line, token.Col)
		case TT_NUMBER:
			_filterNumbersText(token)
			if _isInvalidNumber(token) {
				return nil, NewTokenizerError(token.SourceName, fmt.Sprintf("Invalid numerical literal %q", token.Text), token.Line, token.Col)
			}
		case TT_OPERATOR:
			if _isInvalidOperator(token, previous, next) {
				return nil, NewTokenizerError(token.SourceName, fmt.Sprintf("Invalid operator %q", token.Text), token.Line, token.Col)
			}
			switch token.Text {
			case KW_BIND:
				token.Token = TT_ARROW
			case OP_ASSIGN:
				token.Token = TT_ASSIGNMENT
			case OP_PLUS, OP_MINUS:
				// need to check if next token is numerical literal, it may be negative number
				if tw.Match(TT_OPERATOR, TT_NUMBER) {
					if _isNegativeNumberDetected(previous) {
						_filterNumbersText(next)
						var newText string
						if token.Text == OP_PLUS {
							newText = next.Text
						} else {
							newText = token.Text + next.Text
						}
						replacement := Token{
							Token:      TT_NUMBER,
							Text:       newText,
							SourceName: token.SourceName,
							Line:       token.Line,
							Col:        token.Col,
						}
						if _isInvalidNumber(&replacement) {
							return nil, NewTokenizerError(token.SourceName, fmt.Sprintf("Invalid numerical literal %q", token.Text), token.Line, token.Col)
						}
						optimized = append(optimized, replacement)
						tw.Move(2) // Step over operator and up comming number
						continue
					}
				}
			default:
			}
		case TT_WORD:
			switch token.Text {
			case KW_TRUE, KW_FALSE:
				token.Token = TT_LOGIC
			}
		}
		optimized = append(optimized, *token)
		tw.Move(1)
	}

	tw.Clear()

	return NewTokenWalker(optimized), nil
}

func _filterNumbersText(t *Token) {
	if t.Text[0] == '.' {
		t.Text = "0" + t.Text
	} else if t.Text[len(t.Text)-1] == '.' {
		t.Text += "0"
	}
}

func _isNegativeNumberDetected(previous *Token) bool {
	if previous == nil {
		return true
	}
	switch previous.Token {
	// means that current token is part of arithmetic expression
	case TT_NUMBER, TT_VARIABLE, TT_C_PAREN:
		return false
	}
	return true
}

func _isInvalidOperator(c *Token, p *Token, n *Token) bool {
	switch c.Text {
	case KW_BIND, OP_PLUS, OP_MINUS, OP_DIV, OP_MULT, OP_MOD, OP_ASSIGN, OP_EQUALS, OP_NOT_EQUALS,
		OP_GT, OP_GTE, OP_LT, OP_LTE, OP_NOT, OP_AND, OP_OR, OP_INC, OP_DEC, OP_PLUS_ASSIGN,
		OP_MINUS_ASSIGN, OP_DIV_ASSIGN, OP_MULT_ASSIGN, OP_MOD_ASSIGN:
		return false
	}
	switch c.Text {
	case OP_PLUS_ASSIGN, OP_MINUS_ASSIGN, OP_DIV_ASSIGN, OP_MULT_ASSIGN, OP_MOD_ASSIGN:
		// These kinds of operators must follow by variable
		if p == nil || p.Token != TT_VARIABLE {
			return false
		}
	case OP_INC, OP_DEC:
		// This kinds stands before or after variable
		if (p == nil || p.Token != TT_VARIABLE) && (n == nil || n.Token != TT_VARIABLE) {
			return false
		}
	}
	return true
}

func _isInvalidNumber(c *Token) bool {
	_, ok := strconv.ParseFloat(c.Text, 64)
	return ok != nil
}
