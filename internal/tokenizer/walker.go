package tokenizer

type TokenWalker interface {
	Next() bool
	Get(offset int) *Token
	CanMove(step int) bool
	Move(step int)
	Match(tokens ...TokenType) bool
	Size() int
	Clear()
}

type walker struct {
	tokens   []Token
	position int
}

func NewTokenWalker(tokens []Token) TokenWalker {
	return &walker{tokens: tokens, position: 0}
}

func (w walker) Next() bool {
	return w.position < len(w.tokens)
}

func (w walker) Get(offset int) *Token {
	if w.CanMove(offset) {
		return &w.tokens[w.position+offset]
	}
	return nil
}

func (w walker) CanMove(step int) bool {
	position := w.position + step
	return position > -1 && position < len(w.tokens)
}

func (w *walker) Move(step int) {
	w.position += step
}

func (w walker) Match(tokens ...TokenType) bool {
	patternLength := len(tokens)
	if patternLength == 0 {
		return false
	}
	if !w.CanMove(patternLength) {
		return false
	}
	var (
		argIndex     int  = 0
		curIndex     int  = w.position - 1
		opensCounter int  = 0
		skipping     bool = false
	)
	for {
		if argIndex == patternLength {
			return true
		}
		curIndex += 1
		if curIndex >= len(w.tokens) {
			return false
		}
		current := w.tokens[curIndex].Token
		needed := tokens[argIndex]
		if skipping {
			opensCounter += getCounterIncerement(needed, current)
			if current == needed {
				// opensCounter != -1 <= initial opening did not take into consideration
				if needCountOpens(needed) && opensCounter != -1 {
					continue
				}
				skipping = false
				argIndex += 1
			}
			continue
		}
		argIndex += 1
		if needed == TokenDefault {
			skipping = true
			curIndex -= 1 // retry current token
			opensCounter = 0
			continue
		}
		if needed != current {
			return false
		}
	}
}

func (w walker) Size() int {
	return len(w.tokens)
}

func (w *walker) Clear() {
	w.tokens = nil
}

func needCountOpens(needed TokenType) bool {
	switch needed {
	case TokenCloseBrace, TokenCloseBracket, TokenCloseParen:
		return true
	}
	return false
}

func getCounterIncerement(needed TokenType, tt TokenType) int {
	switch needed {
	case TokenCloseParen:
		switch tt {
		case TokenOpenParen:
			return 1
		case TokenCloseParen:
			return -1
		}
		return 0
	case TokenCloseBrace:
		switch tt {
		case TokenOpenBrace:
			return 1
		case TokenCloseBrace:
			return -1
		}
		return 0
	case TokenCloseBracket:
		switch tt {
		case TokenOpenBracket:
			return 1
		case TokenCloseBracket:
			return -1
		}
		return 0
	}
	return 0
}
