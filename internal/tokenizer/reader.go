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

	tokens  []Token
	builder []rune

	escapedString bool
}

func NewTokenizer(reader io.Reader, sourceName string) *Tokenizer {
	tr := &Tokenizer{
		sourceName:     sourceName,
		reader:         reader,
		repeatCounter:  0,
		currentToken:   TT_DEFAULT,
		currentLine:    0,
		currentCol:     0,
		tokenBegunLine: 0,
		tokenBegunCol:  0,
		tokens:         make([]Token, 0, 255),
		builder:        make([]rune, 0, 1024),
		escapedString:  false,
	}
	return tr
}

func (tr *Tokenizer) Tokenize() (TokenWalker, error) {
	// does all stuff
	bufReader := bufio.NewReader(tr.reader)
	tr._bof()
	for {
		r, _, err := bufReader.ReadRune()
		if err != nil {
			if err == io.EOF {
				tr._eof()
				break
			}
			return nil, NewTokenizerError(tr.sourceName, "Failed to read source: "+err.Error(), tr.currentLine, tr.currentCol)
		}
		tr._countLinesAndCols(r)
		tr._doRepeat()
		for tr._needRepeat() {
			if err := tr._process(r); err != nil {
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
func (tr *Tokenizer) _createFromCurrent() {
	// Skip comments and white spaces on this stage because these may contain more than 1 character more than 1 line
	if tr.currentToken != TT_COMMENT && tr.currentToken != TT_WHITE_SPACE && tr.currentToken != TT_MULTILINE_COMMENT {
		token := Token{
			Token:      tr.currentToken,
			Text:       string(tr.builder),
			SourceName: tr.sourceName,
			Line:       tr.tokenBegunLine,
			Col:        tr.tokenBegunCol,
		}
		tr.tokens = append(tr.tokens, token)
	}
	tr.builder = tr.builder[:0]
	tr.tokenBegunLine = 0
	tr.tokenBegunCol = 0
	tr.currentToken = TT_DEFAULT
}

func (tr *Tokenizer) _createFromTypeAndRune(tt TokenType, r rune) {
	tr._begin(tt)
	tr._put(r)
	tr._createFromCurrent()
}

// Appends runt to builder
func (tr *Tokenizer) _put(r rune) {
	tr.builder = append(tr.builder, r)
}

// Creates and appends token with type BOF
func (tr *Tokenizer) _bof() {
	tr._begin(TT_BOF)
	tr._createFromCurrent()
	tr.currentLine = 1
}

// Creates and appends token with type EOF
func (tr *Tokenizer) _eof() {
	if tr.currentToken != TT_DEFAULT {
		tr._createFromCurrent()
	}
	tr._begin(TT_EOF)
	tr._createFromCurrent()
}

// Processes single rune
func (tr *Tokenizer) _process(r rune) error {
	switch tr.currentToken {
	case TT_DEFAULT:
		return tr._processDefState(r)
	case TT_STRING:
		if tr.escapedString {
			switch r {
			case 't':
				tr._put('\t')
			case 'b':
				tr._put('\b')
			case 'r':
				tr._put('\r')
			case 'n':
				tr._put('\n')
			case 'f':
				tr._put('\f')
			default:
				tr._put(r)
			}
			tr.escapedString = false
		} else {
			switch r {
			case '\\':
				tr.escapedString = true
			case '"':
				tr._createFromCurrent()
			default:
				tr._put(r)
			}
		}
	case TT_NUMBER:
		if unicode.IsDigit(r) {
			tr._put(r)
		} else if r == '.' {
			for _, c := range tr.builder {
				if c == '.' {
					return NewTokenizerError(tr.sourceName, "Unexpected symbol \".\"", tr.currentLine, tr.currentCol)
				}
			}
			tr._put(r)
		} else {
			tr._createFromCurrent()
			tr._doRepeat()
		}
	case TT_POINT:
		if unicode.IsDigit(r) {
			tr.currentToken = TT_NUMBER
		} else {
			tr._createFromCurrent()
		}
		tr._doRepeat()
	case TT_OPERATOR:
		if len(tr.builder) == 2 {
			tr._createFromCurrent()
			tr._doRepeat()
		} else {
			switch r {
			case '>', '<', '=', '!', '+', '-', '/', '*', '&', '|', '%':
				tr._put(r)
				tTxt := string(tr.builder)
				switch tTxt {
				case "//":
					tr.currentToken = TT_COMMENT
					tr.builder = tr.builder[:0]
				case "/*":
					tr.currentToken = TT_MULTILINE_COMMENT
					tr.builder = tr.builder[:0]
				}
			default:
				tr._createFromCurrent()
				tr._doRepeat()
			}
		}
	case TT_WORD:
		if unicode.IsLetter(r) || unicode.IsDigit(r) || r == '_' {
			tr._put(r)
		} else {
			tr._createFromCurrent()
			tr._doRepeat()
		}
	case TT_VARIABLE:
		if unicode.IsLetter(r) || unicode.IsDigit(r) || r == '_' {
			tr._put(r)
		} else {
			if len(tr.builder) > 0 {
				tr._createFromCurrent()
				tr._doRepeat()
			} else {
				return NewTokenizerError(tr.sourceName, "Empty identifier", tr.currentLine, tr.currentCol)
			}
		}
	case TT_COMMENT:
		tr._put(r)
		if r == '\n' {
			tr._createFromCurrent()
		}
	case TT_MULTILINE_COMMENT:
		if r == '/' && len(tr.builder) > 2 {
			prevRune := tr.builder[len(tr.builder)-1]
			if prevRune == '*' {
				tr._createFromCurrent()
			}
		}
	case TT_WHITE_SPACE:
		if unicode.IsSpace(r) {
			tr._put(r)
		} else {
			tr._createFromCurrent()
			tr._doRepeat()
		}
	}
	return nil
}

// Starts new token or creates one form single rune
func (tr *Tokenizer) _processDefState(r rune) error {

	if unicode.IsSpace(r) {
		tr._begin(TT_WHITE_SPACE)
		tr._doRepeat()
		return nil
	}

	if unicode.IsLetter(r) || r == '_' {
		tr._begin(TT_WORD)
		tr._put(r)
		return nil
	}

	if unicode.IsDigit(r) {
		tr._begin(TT_NUMBER)
		tr._put(r)
		return nil
	}

	switch r {
	case '(':
		tr._createFromTypeAndRune(TT_O_PAREN, r)
	case ')':
		tr._createFromTypeAndRune(TT_C_PAREN, r)
	case '[':
		tr._createFromTypeAndRune(TT_O_BRACKET, r)
	case ']':
		tr._createFromTypeAndRune(TT_C_BRACKET, r)
	case '{':
		tr._createFromTypeAndRune(TT_O_BRACE, r)
	case '}':
		tr._createFromTypeAndRune(TT_C_BRACE, r)
	case ':':
		tr._createFromTypeAndRune(TT_COLON, r)
	case ';':
		tr._createFromTypeAndRune(TT_SEMICOLON, r)
	case ',':
		tr._createFromTypeAndRune(TT_COMA, r)
	case '@':
		tr._createFromTypeAndRune(TT_AT, r)
	case '"':
		tr._begin(TT_STRING)
		// leading and terminating quote marks must not be in string
	case '$':
		tr._begin(TT_VARIABLE)
		// $ is skipped
	case '.':
		// Can be beginning of number or just point
		tr._begin(TT_POINT)
		tr._put(r)
	case '#':
		tr._begin(TT_COMMENT)
		// '#' sign is not needed in token's text
	case '>', '<', '=', '!', '+', '-', '/', '*', '&', '|', '%':
		tr._begin(TT_OPERATOR)
		tr._put(r)
	default:
		return NewTokenizerError(tr.sourceName, fmt.Sprintf("Unknown sumbol %q", r), tr.currentLine, tr.currentCol)
	}
	return nil
}

// Initializes state
func (tr *Tokenizer) _begin(tt TokenType) {
	tr.tokenBegunLine = tr.currentLine
	tr.tokenBegunCol = tr.currentCol
	tr.currentToken = tt
}

// Increments line number if current rune is \n otherwise increments col number
func (tr *Tokenizer) _countLinesAndCols(r rune) {
	if r == '\n' {
		tr.currentCol = 0
		tr.currentLine += 1
	} else {
		tr.currentCol += 1
	}
}

// Returns true if current state requires to repeat iteration of rune processing
func (tr *Tokenizer) _needRepeat() bool {
	if tr.repeatCounter > 0 {
		tr.repeatCounter -= 1
		return true
	}
	return false
}

// Sets current state to request to repeat iteration of single rune processing
func (tr *Tokenizer) _doRepeat() {
	tr.repeatCounter += 1
}
