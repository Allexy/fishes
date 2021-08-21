package tokenizer

import (
	"bufio"
	"fmt"
	"io"
	"unicode"
)

type Tokenizer struct {
	sourceName string
	reader     io.Reader

	repeatCounter uint32

	currentToken TokenType
	currentLine  uint32
	currentCol   uint32

	tokenBegunLine uint32
	tokenBegunCol  uint32

	tokens []Token
	buffer []rune

	escapedString bool
}

func NewTokenizer(reader io.Reader, sourceName string) *Tokenizer {
	tr := &Tokenizer{
		sourceName:     sourceName,
		reader:         reader,
		repeatCounter:  0,
		currentToken:   TokenDefault,
		currentLine:    0,
		currentCol:     0,
		tokenBegunLine: 0,
		tokenBegunCol:  0,
		tokens:         make([]Token, 0, 255),
		buffer:         make([]rune, 0, 1024),
		escapedString:  false,
	}
	return tr
}

func (tr *Tokenizer) Tokenize() (TokenWalker, error) {
	// does all stuff
	bufReader := bufio.NewReader(tr.reader)
	tr.createBOF()
	for {
		r, _, err := bufReader.ReadRune()
		if err != nil {
			if err == io.EOF {
				tr.createEOF()
				break
			}
			return nil, NewTokenizerError(tr.sourceName, "failed to read source: "+err.Error(), tr.currentLine, tr.currentCol, err)
		}
		tr.countLinesAndCols(r)
		tr.doRepeat()
		for tr.repeat() {
			if err := tr.process(r); err != nil {
				return nil, err
			}
		}
	}
	walker, err := optimizeAndValidate(NewTokenWalker(tr.tokens))
	if err != nil {
		return nil, err
	}
	return walker, nil
}

// Appends token from current state
//  and resets state
func (tr *Tokenizer) createFromCurrent() {
	// Skip comments and white spaces on this stage because these may contain more than 1 character more than 1 line
	if tr.currentToken != TokenComment && tr.currentToken != TokenWhiteSpace && tr.currentToken != TokenMultilineComment {
		token := Token{
			Token:      tr.currentToken,
			Text:       string(tr.buffer),
			SourceName: tr.sourceName,
			Line:       tr.tokenBegunLine,
			Col:        tr.tokenBegunCol,
		}
		tr.tokens = append(tr.tokens, token)
	}
	tr.buffer = tr.buffer[:0]
	tr.tokenBegunLine = 0
	tr.tokenBegunCol = 0
	tr.currentToken = TokenDefault
}

func (tr *Tokenizer) createFromTypeAndRune(tt TokenType, r rune) {
	tr.beginToken(tt)
	tr.appendToBuffer(r)
	tr.createFromCurrent()
}

// Appends runt to builder
func (tr *Tokenizer) appendToBuffer(r rune) {
	tr.buffer = append(tr.buffer, r)
}

// Creates and appends token with type BOF
func (tr *Tokenizer) createBOF() {
	tr.beginToken(TokenBOF)
	tr.createFromCurrent()
	tr.currentLine = 1
}

// Creates and appends token with type EOF
func (tr *Tokenizer) createEOF() {
	if tr.currentToken != TokenDefault {
		tr.createFromCurrent()
	}
	tr.beginToken(TokenEOF)
	tr.createFromCurrent()
}

// Processes single rune
func (tr *Tokenizer) process(r rune) error {
	switch tr.currentToken {
	case TokenDefault:
		return tr.handleDefaultState(r)
	case TokenString:
		if tr.escapedString {
			switch r {
			case 't':
				tr.appendToBuffer('\t')
			case 'b':
				tr.appendToBuffer('\b')
			case 'r':
				tr.appendToBuffer('\r')
			case 'n':
				tr.appendToBuffer('\n')
			case 'f':
				tr.appendToBuffer('\f')
			default:
				tr.appendToBuffer(r)
			}
			tr.escapedString = false
		} else {
			switch r {
			case '\\':
				tr.escapedString = true
			case '"':
				tr.createFromCurrent()
			default:
				tr.appendToBuffer(r)
			}
		}
	case TokenNumber:
		if unicode.IsDigit(r) {
			tr.appendToBuffer(r)
		} else if r == '.' {
			for _, c := range tr.buffer {
				if c == '.' {
					return NewTokenizerError(tr.sourceName, "unexpected symbol \".\"", tr.currentLine, tr.currentCol, nil)
				}
			}
			tr.appendToBuffer(r)
		} else {
			tr.createFromCurrent()
			tr.doRepeat()
		}
	case TokenPoint:
		if unicode.IsDigit(r) {
			tr.currentToken = TokenNumber
		} else {
			tr.createFromCurrent()
		}
		tr.doRepeat()
	case TokenOperator:
		if len(tr.buffer) == 2 {
			tr.createFromCurrent()
			tr.doRepeat()
		} else {
			switch r {
			case '>', '<', '=', '!', '+', '-', '/', '*', '&', '|', '%':
				tr.appendToBuffer(r)
				tTxt := string(tr.buffer)
				switch tTxt {
				case "//":
					tr.currentToken = TokenComment
					tr.buffer = tr.buffer[:0]
				case "/*":
					tr.currentToken = TokenMultilineComment
					tr.buffer = tr.buffer[:0]
				}
			default:
				tr.createFromCurrent()
				tr.doRepeat()
			}
		}
	case TokenWord:
		if unicode.IsLetter(r) || unicode.IsDigit(r) || r == '_' {
			tr.appendToBuffer(r)
		} else {
			tr.createFromCurrent()
			tr.doRepeat()
		}
	case TokenVariable:
		if unicode.IsLetter(r) || unicode.IsDigit(r) || r == '_' {
			tr.appendToBuffer(r)
		} else {
			if len(tr.buffer) > 0 {
				tr.createFromCurrent()
				tr.doRepeat()
			} else {
				return NewTokenizerError(tr.sourceName, "empty identifier", tr.currentLine, tr.currentCol, nil)
			}
		}
	case TokenComment:
		tr.appendToBuffer(r)
		if r == '\n' {
			tr.createFromCurrent()
		}
	case TokenMultilineComment:
		if r == '/' && len(tr.buffer) > 2 {
			prevRune := tr.buffer[len(tr.buffer)-1]
			if prevRune == '*' {
				tr.createFromCurrent()
			}
		}
	case TokenWhiteSpace:
		if unicode.IsSpace(r) {
			tr.appendToBuffer(r)
		} else {
			tr.createFromCurrent()
			tr.doRepeat()
		}
	}
	return nil
}

// Starts new token or creates one form single rune
func (tr *Tokenizer) handleDefaultState(r rune) error {

	if unicode.IsSpace(r) {
		tr.beginToken(TokenWhiteSpace)
		tr.doRepeat()
		return nil
	}

	if unicode.IsLetter(r) || r == '_' {
		tr.beginToken(TokenWord)
		tr.appendToBuffer(r)
		return nil
	}

	if unicode.IsDigit(r) {
		tr.beginToken(TokenNumber)
		tr.appendToBuffer(r)
		return nil
	}

	switch r {
	case '(':
		tr.createFromTypeAndRune(TokenOpenParen, r)
	case ')':
		tr.createFromTypeAndRune(TokenCloseParen, r)
	case '[':
		tr.createFromTypeAndRune(TokenOpenBracket, r)
	case ']':
		tr.createFromTypeAndRune(TokenCloseBracket, r)
	case '{':
		tr.createFromTypeAndRune(TokenOpenBrace, r)
	case '}':
		tr.createFromTypeAndRune(TokenCloseBrace, r)
	case ':':
		tr.createFromTypeAndRune(TokenColon, r)
	case ';':
		tr.createFromTypeAndRune(TokenSemicolon, r)
	case ',':
		tr.createFromTypeAndRune(TokenComa, r)
	case '@':
		tr.createFromTypeAndRune(TokenAt, r)
	case '"':
		tr.beginToken(TokenString)
		// leading and terminating quote marks must not be in string
	case '$':
		tr.beginToken(TokenVariable)
		// $ is skipped
	case '.':
		// Can be beginning of number or just point
		tr.beginToken(TokenPoint)
		tr.appendToBuffer(r)
	case '#':
		tr.beginToken(TokenComment)
		// '#' sign is not needed in token's text
	case '>', '<', '=', '!', '+', '-', '/', '*', '&', '|', '%':
		tr.beginToken(TokenOperator)
		tr.appendToBuffer(r)
	default:
		return NewTokenizerError(tr.sourceName, fmt.Sprintf("unknown sumbol %q", r), tr.currentLine, tr.currentCol, nil)
	}
	return nil
}

// Initializes state
func (tr *Tokenizer) beginToken(tt TokenType) {
	tr.tokenBegunLine = tr.currentLine
	tr.tokenBegunCol = tr.currentCol
	tr.currentToken = tt
}

// Increments line number if current rune is \n otherwise increments col number
func (tr *Tokenizer) countLinesAndCols(r rune) {
	if r == '\n' {
		tr.currentCol = 0
		tr.currentLine += 1
	} else {
		tr.currentCol += 1
	}
}

// Returns true if current state requires to repeat iteration of rune processing
func (tr *Tokenizer) repeat() bool {
	if tr.repeatCounter > 0 {
		tr.repeatCounter -= 1
		return true
	}
	return false
}

// Sets current state to request to repeat iteration of single rune processing
func (tr *Tokenizer) doRepeat() {
	tr.repeatCounter += 1
}
