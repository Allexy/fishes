package tokenizer

import (
	"github.com/Allexy/fishes/internal/lang"
	"github.com/stretchr/testify/require"
	"io"
	"strings"
	"testing"
)

// reader creates and returns strings.Reader and source name
func reader(s string) (io.Reader, string) {
	return strings.NewReader(s), "string"
}

// Testing walker

func TestTokenWalker(t *testing.T) {
	tw := NewTokenWalker([]Token{{}, {}, {}})
	require.Equalf(t, 3, tw.Size(), "Walker must contain 3 items got %d", tw.Size())
	var step = 0
	for ; tw.Next(); step++ {
		tw.Move(1)
		require.Lessf(t, step, 3, "expected that walker iterate by 2 steps, current step: %d", step)
	}
}

func TestTokenWalkerMatch(t *testing.T) {
	// Simulate typical function definition sequence of tokens
	tokens := []Token{
		{Type: TokenBOF},
		{Type: TokenWord}, {Type: TokenWord}, {Type: TokenOpenParen}, {Type: TokenCloseParen},
		{Type: TokenOpenBrace}, {Type: TokenOpenBrace}, {Type: TokenCloseBrace}, {Type: TokenWord},
		{Type: TokenOpenBrace}, {Type: TokenCloseBrace}, {Type: TokenWord}, {Type: TokenCloseBrace},
		{Type: TokenEOF},
	}
	// Any tokens between parentheses and curly braces are matched by TT_DEFAULT, walker takes into consideration nested
	// opening and closing parentheses
	tw := NewTokenWalker(tokens)
	require.False(
		t,
		tw.Match(TokenWord, TokenWord, TokenOpenParen, TokenDefault, TokenCloseParen, TokenOpenBrace, TokenDefault, TokenCloseBrace),
		"When steps counter of TokenWalker points to TT_BOF it must not match the pattern",
	)
	tw.Move(1)
	require.True(
		t,
		tw.Match(TokenWord, TokenWord, TokenOpenParen, TokenDefault, TokenCloseParen, TokenOpenBrace, TokenDefault, TokenCloseBrace),
		"When steps counter of TokenWalker points to TT_WORD it must match the pattern",
	)
}

// Testing words

func TestBOF(t *testing.T) {
	tw, err := NewTokenizer(reader("some_thing")).Tokenize()
	require.NoError(t, err, "Must not be tokenization error")
	require.Equalf(t, 3, tw.Size(), "Walker must contain 3 items got %d", tw.Size())
	token := tw.Get(0)
	require.Equalf(t, TokenBOF, token.Type, "expected token type %v got %v", TokenBOF, token.Type)
}

func TestWord(t *testing.T) {
	tw, err := NewTokenizer(reader("word")).Tokenize()
	require.NoError(t, err, "Must not be tokenization error")
	require.Equalf(t, 3, tw.Size(), "Walker must contain 3 items got %d", tw.Size())
	token := tw.Get(1)
	require.Equalf(t, TokenWord, token.Type, "expected token type %v got %v", TokenWord, token.Type)
	require.Equalf(t, "word", token.Text, "expected token text is %q got %q", "word", token.Text)
}

func TestWordWithUnderscore(t *testing.T) {
	tw, err := NewTokenizer(reader("word_word_word")).Tokenize()
	require.NoError(t, err, "Must not be tokenization error")
	require.Equalf(t, 3, tw.Size(), "Walker must contain 3 items got %d", tw.Size())
	token := tw.Get(1)
	require.Equalf(t, TokenWord, token.Type, "expected token type %v got %v", TokenWord, token.Type)
	require.Equalf(t, "word_word_word", token.Text, "expected token text is %q got %q", "word_word_word", token.Text)
}

func TestWordWithNumber(t *testing.T) {
	tw, err := NewTokenizer(reader("word00203word")).Tokenize()
	require.NoError(t, err, "Must not be tokenization error")
	require.Equalf(t, 3, tw.Size(), "Walker must contain 3 items got %d", tw.Size())
	token := tw.Get(1)
	require.Equalf(t, TokenWord, token.Type, "expected token type %v got %v", TokenWord, token.Type)
	require.Equalf(t, "word00203word", token.Text, "expected token text is %q got %q", "word00203word", token.Text)
}

func TestWordStartingFromNumber(t *testing.T) {
	tw, err := NewTokenizer(reader("203word")).Tokenize()
	require.NoError(t, err, "Must not be tokenization error")
	require.Equalf(t, 4, tw.Size(), "Walker must contain 4 items got %d", tw.Size())
	var token *Token
	token = tw.Get(1)
	require.Equalf(t, TokenNumber, token.Type, "expected token type %v got %v", TokenNumber, token.Type)
	require.Equalf(t, "203", token.Text, "expected token text is %q got %q", "203", token.Text)
	token = tw.Get(2)
	require.Equalf(t, TokenWord, token.Type, "expected token type %v got %v", TokenWord, token.Type)
	require.Equalf(t, "word", token.Text, "expected token text is %q got %q", "word", token.Text)
}

func TestWordWithEverything(t *testing.T) {
	tw, err := NewTokenizer(reader("word00203word_s")).Tokenize()
	require.NoError(t, err, "Must not be tokenization error")
	require.Equalf(t, 3, tw.Size(), "Walker must contain 3 items got %d", tw.Size())
	token := tw.Get(1)
	require.Equalf(t, TokenWord, token.Type, "expected token type %v got %v", TokenWord, token.Type)
	require.Equalf(t, "word00203word_s", token.Text, "expected token text is %q got %q", "word00203word_s", token.Text)
}

// Testing numerical literals

func TestNumberSimple(t *testing.T) {
	tw, err := NewTokenizer(reader("10")).Tokenize()
	require.NoError(t, err, "Must not be tokenization error")
	require.Equalf(t, 3, tw.Size(), "Walker must contain 3 items got %d", tw.Size())
	token := tw.Get(1)
	require.Equalf(t, TokenNumber, token.Type, "expected token type %v got %v", TokenNumber, token.Type)
	require.Equalf(t, "10", token.Text, "expected token text is %q got %q", "10", token.Text)
}

func TestNumberWithPointInMiddle(t *testing.T) {
	tw, err := NewTokenizer(reader("10.5")).Tokenize()
	require.NoError(t, err, "Must not be tokenization error")
	require.Equalf(t, 3, tw.Size(), "Walker must contain 3 items got %d", tw.Size())
	token := tw.Get(1)
	require.Equalf(t, TokenNumber, token.Type, "expected token type %v got %v", TokenNumber, token.Type)
	require.Equalf(t, "10.5", token.Text, "expected token text is %q got %q", "10.5", token.Text)
}

func TestNumberWithLeadingPoint(t *testing.T) {
	tw, err := NewTokenizer(reader(".5")).Tokenize()
	require.NoError(t, err, "Must not be tokenization error")
	require.Equalf(t, 3, tw.Size(), "Walker must contain 3 items got %d", tw.Size())
	token := tw.Get(1)
	require.Equalf(t, TokenNumber, token.Type, "expected token type %v got %v", TokenNumber, token.Type)
	require.Equalf(t, "0.5", token.Text, "expected token text is %q got %q", "0.5", token.Text)
}

func TestNumberWithTerminatingPoint(t *testing.T) {
	tw, err := NewTokenizer(reader("5.")).Tokenize()
	require.NoError(t, err, "Must not be tokenization error")
	require.Equalf(t, 3, tw.Size(), "Walker must contain 3 items got %d", tw.Size())
	token := tw.Get(1)
	require.Equalf(t, TokenNumber, token.Type, "expected token type %v got %v", TokenNumber, token.Type)
	require.Equalf(t, "5.0", token.Text, "expected token text is %q got %q", "5.0", token.Text)
}

func TestNegativeNumberSimple(t *testing.T) {
	tw, err := NewTokenizer(reader("-10")).Tokenize()
	require.NoError(t, err, "Must not be tokenization error")
	require.Equalf(t, 3, tw.Size(), "Walker must contain 3 items got %d", tw.Size())
	token := tw.Get(1)
	require.Equalf(t, TokenNumber, token.Type, "expected token type %v got %v", TokenNumber, token.Type)
	require.Equalf(t, "-10", token.Text, "expected token text is %q got %q", "-10", token.Text)
}

func TestNegativeNumberWithPointInMiddle(t *testing.T) {
	tw, err := NewTokenizer(reader("-10.5")).Tokenize()
	require.NoError(t, err, "Must not be tokenization error")
	require.Equalf(t, 3, tw.Size(), "Walker must contain 3 items got %d", tw.Size())
	token := tw.Get(1)
	require.Equalf(t, TokenNumber, token.Type, "expected token type %v got %v", TokenNumber, token.Type)
	require.Equalf(t, "-10.5", token.Text, "expected token text is %q got %q", "-10.5", token.Text)
}

func TestNegativeNumberWithLeadingPoint(t *testing.T) {
	tw, err := NewTokenizer(reader("-.5")).Tokenize()
	require.NoError(t, err, "Must not be tokenization error")
	require.Equalf(t, 3, tw.Size(), "Walker must contain 3 items got %d", tw.Size())
	token := tw.Get(1)
	require.Equalf(t, TokenNumber, token.Type, "expected token type %v got %v", TokenNumber, token.Type)
	require.Equalf(t, "-0.5", token.Text, "expected token text is %q got %q", "-0.5", token.Text)
}

func TestNegativeNumberWithTerminatingPoint(t *testing.T) {
	tw, err := NewTokenizer(reader("-5.")).Tokenize()
	require.NoError(t, err, "Must not be tokenization error")
	require.Equalf(t, 3, tw.Size(), "Walker must contain 3 items got %d", tw.Size())
	token := tw.Get(1)
	require.Equalf(t, TokenNumber, token.Type, "expected token type %v got %v", TokenNumber, token.Type)
	require.Equalf(t, "-5.0", token.Text, "expected token text is %q got %q", "-5.0", token.Text)
}

func TestNumberStartedFromZero(t *testing.T) {
	tw, err := NewTokenizer(reader("002")).Tokenize()
	require.NoError(t, err, "Must not be tokenization error")
	require.Equalf(t, 3, tw.Size(), "Walker must contain 3 items got %d", tw.Size())
	token := tw.Get(1)
	require.Equalf(t, TokenNumber, token.Type, "expected token type %v got %v", TokenNumber, token.Type)
	require.Equalf(t, "2", token.Text, "expected token text is %q got %q", "2", token.Text)
}

func TestNumberContainsOnlyZeros(t *testing.T) {
	tw, err := NewTokenizer(reader("00000")).Tokenize()
	require.NoError(t, err, "Must not be tokenization error")
	require.Equalf(t, 3, tw.Size(), "Walker must contain 3 items got %d", tw.Size())
	token := tw.Get(1)
	require.Equalf(t, TokenNumber, token.Type, "expected token type %v got %v", TokenNumber, token.Type)
	require.Equalf(t, "0", token.Text, "expected token text is %q got %q", "0", token.Text)
}

func TestNumberStartsWithZerosAndPoint(t *testing.T) {
	tw, err := NewTokenizer(reader("00000.01")).Tokenize()
	require.NoError(t, err, "Must not be tokenization error")
	require.Equalf(t, 3, tw.Size(), "Walker must contain 3 items got %d", tw.Size())
	token := tw.Get(1)
	require.Equalf(t, TokenNumber, token.Type, "expected token type %v got %v", TokenNumber, token.Type)
	require.Equalf(t, "0.01", token.Text, "expected token text is %q got %q", "0.01", token.Text)
}

func TestNegativeZero(t *testing.T) {
	tw, err := NewTokenizer(reader("-00000")).Tokenize()
	require.NoError(t, err, "Must not be tokenization error")
	require.Equalf(t, 3, tw.Size(), "Walker must contain 3 items got %d", tw.Size())
	token := tw.Get(1)
	require.Equalf(t, TokenNumber, token.Type, "expected token type %v got %v", TokenNumber, token.Type)
	require.Equalf(t, "-0", token.Text, "expected token text is %q got %q", "-0", token.Text)
}

// Testing string literals

func TestString(t *testing.T) {
	tw, err := NewTokenizer(reader("\"Test string\"")).Tokenize()
	require.NoError(t, err, "Must not be tokenization error")
	require.Equalf(t, 3, tw.Size(), "Walker must contain 3 items got %d", tw.Size())
	token := tw.Get(1)
	require.Equalf(t, TokenString, token.Type, "expected token type %v got %v", TokenString, token.Type)
	require.Equalf(t, "Test string", token.Text, "expected token text is %q got %q", "Test string", token.Text)
}

func TestStringWithEscSeq(t *testing.T) {
	tw, err := NewTokenizer(reader("\"Test\\nstring\"")).Tokenize()
	require.NoError(t, err, "Must not be tokenization error")
	require.Equalf(t, 3, tw.Size(), "Walker must contain 3 items got %d", tw.Size())
	token := tw.Get(1)
	require.Equalf(t, TokenString, token.Type, "expected token type %v got %v", TokenString, token.Type)
	require.Equalf(t, "Test\nstring", token.Text, "expected token text is %q got %q", "Test\nstring", token.Text)
}

func TestStringWithNestedQuotes(t *testing.T) {
	tw, err := NewTokenizer(reader("\"Test \\\"string\\\"\"")).Tokenize()
	require.NoError(t, err, "Must not be tokenization error")
	require.Equalf(t, 3, tw.Size(), "Walker must contain 3 items got %d", tw.Size())
	token := tw.Get(1)
	require.Equalf(t, TokenString, token.Type, "expected token type %v got %v", TokenString, token.Type)
	require.Equalf(t, "Test \"string\"", token.Text, "expected token text is %q got %q", "Test \"string\"", token.Text)
}

// Testing logical literals

func TestLogicalTrue(t *testing.T) {
	tw, err := NewTokenizer(reader("true")).Tokenize()
	require.NoError(t, err, "Must not be tokenization error")
	require.Equalf(t, 3, tw.Size(), "Walker must contain 3 items got %d", tw.Size())
	token := tw.Get(1)
	require.Equalf(t, TokenLogic, token.Type, "expected token type %v got %v", TokenLogic, token.Type)
	require.Equalf(t, "true", token.Text, "expected token text is %q got %q", "true", token.Text)
}

func TestLogicalFalse(t *testing.T) {
	tw, err := NewTokenizer(reader("false")).Tokenize()
	require.NoError(t, err, "Must not be tokenization error")
	require.Equalf(t, 3, tw.Size(), "Walker must contain 3 items got %d", tw.Size())
	token := tw.Get(1)
	require.Equalf(t, TokenLogic, token.Type, "expected token type %v got %v", TokenLogic, token.Type)
	require.Equalf(t, "false", token.Text, "expected token text is %q got %q", "false", token.Text)
}

// Testing operators

func TestOpPlus(t *testing.T) {
	tw, err := NewTokenizer(reader("$a + $b")).Tokenize()
	require.NoError(t, err, "Must not be tokenization error")
	require.Equalf(t, 5, tw.Size(), "Walker must contain 3 items got %d", tw.Size())

	var token *Token
	token = tw.Get(1)
	require.Equalf(t, TokenVariable, token.Type, "expected token type %v got %v", TokenVariable, token.Type)
	require.Equalf(t, "a", token.Text, "expected token text is %q got %q", "a", token.Text)

	token = tw.Get(2)
	require.Equalf(t, TokenOperator, token.Type, "expected token type %v got %v", TokenOperator, token.Type)
	require.Equalf(t, lang.OpPlus, token.Text, "expected token text is %q got %q", lang.OpPlus, token.Text)

	token = tw.Get(3)
	require.Equalf(t, TokenVariable, token.Type, "expected token type %v got %v", TokenVariable, token.Type)
	require.Equalf(t, "b", token.Text, "expected token text is %q got %q", "b", token.Text)
}

func TestOpMinus(t *testing.T) {
	tw, err := NewTokenizer(reader("$a - $b")).Tokenize()
	require.NoError(t, err, "Must not be tokenization error")
	require.Equalf(t, 5, tw.Size(), "Walker must contain 3 items got %d", tw.Size())

	var token *Token
	token = tw.Get(1)
	require.Equalf(t, TokenVariable, token.Type, "expected token type %v got %v", TokenVariable, token.Type)
	require.Equalf(t, "a", token.Text, "expected token text is %q got %q", "a", token.Text)

	token = tw.Get(2)
	require.Equalf(t, TokenOperator, token.Type, "expected token type %v got %v", TokenOperator, token.Type)
	require.Equalf(t, lang.OpMinus, token.Text, "expected token text is %q got %q", lang.OpMinus, token.Text)

	token = tw.Get(3)
	require.Equalf(t, TokenVariable, token.Type, "expected token type %v got %v", TokenVariable, token.Type)
	require.Equalf(t, "b", token.Text, "expected token text is %q got %q", "b", token.Text)
}

func TestOpDivision(t *testing.T) {
	tw, err := NewTokenizer(reader("$a / $b")).Tokenize()
	require.NoError(t, err, "Must not be tokenization error")
	require.Equalf(t, 5, tw.Size(), "Walker must contain 3 items got %d", tw.Size())

	var token *Token
	token = tw.Get(1)
	require.Equalf(t, TokenVariable, token.Type, "expected token type %v got %v", TokenVariable, token.Type)
	require.Equalf(t, "a", token.Text, "expected token text is %q got %q", "a", token.Text)

	token = tw.Get(2)
	require.Equalf(t, TokenOperator, token.Type, "expected token type %v got %v", TokenOperator, token.Type)
	require.Equalf(t, lang.OpDivision, token.Text, "expected token text is %q got %q", lang.OpDivision, token.Text)

	token = tw.Get(3)
	require.Equalf(t, TokenVariable, token.Type, "expected token type %v got %v", TokenVariable, token.Type)
	require.Equalf(t, "b", token.Text, "expected token text is %q got %q", "b", token.Text)
}

func TestOpMultiply(t *testing.T) {
	tw, err := NewTokenizer(reader("$a * $b")).Tokenize()
	require.NoError(t, err, "Must not be tokenization error")
	require.Equalf(t, 5, tw.Size(), "Walker must contain 3 items got %d", tw.Size())

	var token *Token
	token = tw.Get(1)
	require.Equalf(t, TokenVariable, token.Type, "expected token type %v got %v", TokenVariable, token.Type)
	require.Equalf(t, "a", token.Text, "expected token text is %q got %q", "a", token.Text)

	token = tw.Get(2)
	require.Equalf(t, TokenOperator, token.Type, "expected token type %v got %v", TokenOperator, token.Type)
	require.Equalf(t, lang.OpMultiply, token.Text, "expected token text is %q got %q", lang.OpMultiply, token.Text)

	token = tw.Get(3)
	require.Equalf(t, TokenVariable, token.Type, "expected token type %v got %v", TokenVariable, token.Type)
	require.Equalf(t, "b", token.Text, "expected token text is %q got %q", "b", token.Text)
}

func TestOpModulo(t *testing.T) {
	tw, err := NewTokenizer(reader("$a % $b")).Tokenize()
	require.NoError(t, err, "Must not be tokenization error")
	require.Equalf(t, 5, tw.Size(), "Walker must contain 3 items got %d", tw.Size())

	var token *Token
	token = tw.Get(1)
	require.Equalf(t, TokenVariable, token.Type, "expected token type %v got %v", TokenVariable, token.Type)
	require.Equalf(t, "a", token.Text, "expected token text is %q got %q", "a", token.Text)

	token = tw.Get(2)
	require.Equalf(t, TokenOperator, token.Type, "expected token type %v got %v", TokenOperator, token.Type)
	require.Equalf(t, lang.OpModulo, token.Text, "expected token text is %q got %q", lang.OpModulo, token.Text)

	token = tw.Get(3)
	require.Equalf(t, TokenVariable, token.Type, "expected token type %v got %v", TokenVariable, token.Type)
	require.Equalf(t, "b", token.Text, "expected token text is %q got %q", "b", token.Text)
}

func TestOpAssign(t *testing.T) {
	tw, err := NewTokenizer(reader("$a = $b")).Tokenize()
	require.NoError(t, err, "Must not be tokenization error")
	require.Equalf(t, 5, tw.Size(), "Walker must contain 3 items got %d", tw.Size())

	var token *Token
	token = tw.Get(1)
	require.Equalf(t, TokenVariable, token.Type, "expected token type %v got %v", TokenVariable, token.Type)
	require.Equalf(t, "a", token.Text, "expected token text is %q got %q", "a", token.Text)

	token = tw.Get(2)
	require.Equalf(t, TokenAssignment, token.Type, "expected token type %v got %v", TokenAssignment, token.Type)
	require.Equalf(t, lang.OpAssign, token.Text, "expected token text is %q got %q", lang.OpAssign, token.Text)

	token = tw.Get(3)
	require.Equalf(t, TokenVariable, token.Type, "expected token type %v got %v", TokenVariable, token.Type)
	require.Equalf(t, "b", token.Text, "expected token text is %q got %q", "b", token.Text)
}

func TestOpEquals(t *testing.T) {
	tw, err := NewTokenizer(reader("$a == $b")).Tokenize()
	require.NoError(t, err, "Must not be tokenization error")
	require.Equalf(t, 5, tw.Size(), "Walker must contain 3 items got %d", tw.Size())

	var token *Token
	token = tw.Get(1)
	require.Equalf(t, TokenVariable, token.Type, "expected token type %v got %v", TokenVariable, token.Type)
	require.Equalf(t, "a", token.Text, "expected token text is %q got %q", "a", token.Text)

	token = tw.Get(2)
	require.Equalf(t, TokenOperator, token.Type, "expected token type %v got %v", TokenOperator, token.Type)
	require.Equalf(t, lang.OpEquals, token.Text, "expected token text is %q got %q", lang.OpEquals, token.Text)

	token = tw.Get(3)
	require.Equalf(t, TokenVariable, token.Type, "expected token type %v got %v", TokenVariable, token.Type)
	require.Equalf(t, "b", token.Text, "expected token text is %q got %q", "b", token.Text)
}

func TestOpNotEquals(t *testing.T) {
	tw, err := NewTokenizer(reader("$a != $b")).Tokenize()
	require.NoError(t, err, "Must not be tokenization error")
	require.Equalf(t, 5, tw.Size(), "Walker must contain 3 items got %d", tw.Size())

	var token *Token
	token = tw.Get(1)
	require.Equalf(t, TokenVariable, token.Type, "expected token type %v got %v", TokenVariable, token.Type)
	require.Equalf(t, "a", token.Text, "expected token text is %q got %q", "a", token.Text)

	token = tw.Get(2)
	require.Equalf(t, TokenOperator, token.Type, "expected token type %v got %v", TokenOperator, token.Type)
	require.Equalf(t, lang.OpNotEquals, token.Text, "expected token text is %q got %q", lang.OpNotEquals, token.Text)

	token = tw.Get(3)
	require.Equalf(t, TokenVariable, token.Type, "expected token type %v got %v", TokenVariable, token.Type)
	require.Equalf(t, "b", token.Text, "expected token text is %q got %q", "b", token.Text)
}

func TestOpGreaterThan(t *testing.T) {
	tw, err := NewTokenizer(reader("$a > $b")).Tokenize()
	require.NoError(t, err, "Must not be tokenization error")
	require.Equalf(t, 5, tw.Size(), "Walker must contain 3 items got %d", tw.Size())

	var token *Token
	token = tw.Get(1)
	require.Equalf(t, TokenVariable, token.Type, "expected token type %v got %v", TokenVariable, token.Type)
	require.Equalf(t, "a", token.Text, "expected token text is %q got %q", "a", token.Text)

	token = tw.Get(2)
	require.Equalf(t, TokenOperator, token.Type, "expected token type %v got %v", TokenOperator, token.Type)
	require.Equalf(t, lang.OpGreaterThan, token.Text, "expected token text is %q got %q", lang.OpGreaterThan, token.Text)

	token = tw.Get(3)
	require.Equalf(t, TokenVariable, token.Type, "expected token type %v got %v", TokenVariable, token.Type)
	require.Equalf(t, "b", token.Text, "expected token text is %q got %q", "b", token.Text)
}

func TestOpGreaterThanOrEquals(t *testing.T) {
	tw, err := NewTokenizer(reader("$a >= $b")).Tokenize()
	require.NoError(t, err, "Must not be tokenization error")
	require.Equalf(t, 5, tw.Size(), "Walker must contain 3 items got %d", tw.Size())

	var token *Token
	token = tw.Get(1)
	require.Equalf(t, TokenVariable, token.Type, "expected token type %v got %v", TokenVariable, token.Type)
	require.Equalf(t, "a", token.Text, "expected token text is %q got %q", "a", token.Text)

	token = tw.Get(2)
	require.Equalf(t, TokenOperator, token.Type, "expected token type %v got %v", TokenOperator, token.Type)
	require.Equalf(t, lang.OpGreaterThanOrEquals, token.Text, "expected token text is %q got %q", lang.OpGreaterThanOrEquals, token.Text)

	token = tw.Get(3)
	require.Equalf(t, TokenVariable, token.Type, "expected token type %v got %v", TokenVariable, token.Type)
	require.Equalf(t, "b", token.Text, "expected token text is %q got %q", "b", token.Text)
}

func TestOpArrow(t *testing.T) {
	tw, err := NewTokenizer(reader("$a => $b")).Tokenize()
	require.NoError(t, err, "Must not be tokenization error")
	require.Equalf(t, 5, tw.Size(), "Walker must contain 3 items got %d", tw.Size())

	var token *Token
	token = tw.Get(1)
	require.Equalf(t, TokenVariable, token.Type, "expected token type %v got %v", TokenVariable, token.Type)
	require.Equalf(t, "a", token.Text, "expected token text is %q got %q", "a", token.Text)

	token = tw.Get(2)
	require.Equalf(t, TokenArrow, token.Type, "expected token type %v got %v", TokenArrow, token.Type)
	require.Equalf(t, lang.OpArrow, token.Text, "expected token text is %q got %q", lang.OpArrow, token.Text)

	token = tw.Get(3)
	require.Equalf(t, TokenVariable, token.Type, "expected token type %v got %v", TokenVariable, token.Type)
	require.Equalf(t, "b", token.Text, "expected token text is %q got %q", "b", token.Text)
}

func TestOpLesserThan(t *testing.T) {
	tw, err := NewTokenizer(reader("$a < $b")).Tokenize()
	require.NoError(t, err, "Must not be tokenization error")
	require.Equalf(t, 5, tw.Size(), "Walker must contain 3 items got %d", tw.Size())

	var token *Token
	token = tw.Get(1)
	require.Equalf(t, TokenVariable, token.Type, "expected token type %v got %v", TokenVariable, token.Type)
	require.Equalf(t, "a", token.Text, "expected token text is %q got %q", "a", token.Text)

	token = tw.Get(2)
	require.Equalf(t, TokenOperator, token.Type, "expected token type %v got %v", TokenOperator, token.Type)
	require.Equalf(t, lang.OpLesserThan, token.Text, "expected token text is %q got %q", lang.OpLesserThan, token.Text)

	token = tw.Get(3)
	require.Equalf(t, TokenVariable, token.Type, "expected token type %v got %v", TokenVariable, token.Type)
	require.Equalf(t, "b", token.Text, "expected token text is %q got %q", "b", token.Text)
}

func TestOpLesserThanOrEquals(t *testing.T) {
	tw, err := NewTokenizer(reader("$a <= $b")).Tokenize()
	require.NoError(t, err, "Must not be tokenization error")
	require.Equalf(t, 5, tw.Size(), "Walker must contain 3 items got %d", tw.Size())

	var token *Token
	token = tw.Get(1)
	require.Equalf(t, TokenVariable, token.Type, "expected token type %v got %v", TokenVariable, token.Type)
	require.Equalf(t, "a", token.Text, "expected token text is %q got %q", "a", token.Text)

	token = tw.Get(2)
	require.Equalf(t, TokenOperator, token.Type, "expected token type %v got %v", TokenOperator, token.Type)
	require.Equalf(t, lang.OpLesserThanOrEquals, token.Text, "expected token text is %q got %q", lang.OpLesserThanOrEquals, token.Text)

	token = tw.Get(3)
	require.Equalf(t, TokenVariable, token.Type, "expected token type %v got %v", TokenVariable, token.Type)
	require.Equalf(t, "b", token.Text, "expected token text is %q got %q", "b", token.Text)
}

func TestOpNot(t *testing.T) {
	tw, err := NewTokenizer(reader("= !$b")).Tokenize()
	require.NoError(t, err, "Must not be tokenization error")
	require.Equalf(t, 5, tw.Size(), "Walker must contain 3 items got %d", tw.Size())

	var token *Token
	token = tw.Get(1)
	require.Equalf(t, TokenAssignment, token.Type, "expected token type %v got %v", TokenAssignment, token.Type)
	require.Equalf(t, lang.OpAssign, token.Text, "expected token text is %q got %q", lang.OpAssign, token.Text)

	token = tw.Get(2)
	require.Equalf(t, TokenOperator, token.Type, "expected token type %v got %v", TokenOperator, token.Type)
	require.Equalf(t, lang.OpNot, token.Text, "expected token text is %q got %q", lang.OpNot, token.Text)

	token = tw.Get(3)
	require.Equalf(t, TokenVariable, token.Type, "expected token type %v got %v", TokenVariable, token.Type)
	require.Equalf(t, "b", token.Text, "expected token text is %q got %q", "b", token.Text)
}

func TestOpAnd(t *testing.T) {
	tw, err := NewTokenizer(reader("$a && $b")).Tokenize()
	require.NoError(t, err, "Must not be tokenization error")
	require.Equalf(t, 5, tw.Size(), "Walker must contain 3 items got %d", tw.Size())

	var token *Token
	token = tw.Get(1)
	require.Equalf(t, TokenVariable, token.Type, "expected token type %v got %v", TokenVariable, token.Type)
	require.Equalf(t, "a", token.Text, "expected token text is %q got %q", "a", token.Text)

	token = tw.Get(2)
	require.Equalf(t, TokenOperator, token.Type, "expected token type %v got %v", TokenOperator, token.Type)
	require.Equalf(t, lang.OpAnd, token.Text, "expected token text is %q got %q", lang.OpAnd, token.Text)

	token = tw.Get(3)
	require.Equalf(t, TokenVariable, token.Type, "expected token type %v got %v", TokenVariable, token.Type)
	require.Equalf(t, "b", token.Text, "expected token text is %q got %q", "b", token.Text)
}

func TestOpOr(t *testing.T) {
	tw, err := NewTokenizer(reader("$a || $b")).Tokenize()
	require.NoError(t, err, "Must not be tokenization error")
	require.Equalf(t, 5, tw.Size(), "Walker must contain 3 items got %d", tw.Size())

	var token *Token
	token = tw.Get(1)
	require.Equalf(t, TokenVariable, token.Type, "expected token type %v got %v", TokenVariable, token.Type)
	require.Equalf(t, "a", token.Text, "expected token text is %q got %q", "a", token.Text)

	token = tw.Get(2)
	require.Equalf(t, TokenOperator, token.Type, "expected token type %v got %v", TokenOperator, token.Type)
	require.Equalf(t, lang.OpOr, token.Text, "expected token text is %q got %q", lang.OpOr, token.Text)

	token = tw.Get(3)
	require.Equalf(t, TokenVariable, token.Type, "expected token type %v got %v", TokenVariable, token.Type)
	require.Equalf(t, "b", token.Text, "expected token text is %q got %q", "b", token.Text)
}

func TestOpIncrement(t *testing.T) {
	tw, err := NewTokenizer(reader("$a ++")).Tokenize()
	require.NoError(t, err, "Must not be tokenization error")
	require.Equalf(t, 4, tw.Size(), "Walker must contain 3 items got %d", tw.Size())

	var token *Token
	token = tw.Get(1)
	require.Equalf(t, TokenVariable, token.Type, "expected token type %v got %v", TokenVariable, token.Type)
	require.Equalf(t, "a", token.Text, "expected token text is %q got %q", "a", token.Text)

	token = tw.Get(2)
	require.Equalf(t, TokenOperator, token.Type, "expected token type %v got %v", TokenOperator, token.Type)
	require.Equalf(t, lang.OpIncrement, token.Text, "expected token text is %q got %q", lang.OpIncrement, token.Text)
}

func TestOpDecrement(t *testing.T) {
	tw, err := NewTokenizer(reader("$a --")).Tokenize()
	require.NoError(t, err, "Must not be tokenization error")
	require.Equalf(t, 4, tw.Size(), "Walker must contain 3 items got %d", tw.Size())

	var token *Token
	token = tw.Get(1)
	require.Equalf(t, TokenVariable, token.Type, "expected token type %v got %v", TokenVariable, token.Type)
	require.Equalf(t, "a", token.Text, "expected token text is %q got %q", "a", token.Text)

	token = tw.Get(2)
	require.Equalf(t, TokenOperator, token.Type, "expected token type %v got %v", TokenOperator, token.Type)
	require.Equalf(t, lang.OpDecrement, token.Text, "expected token text is %q got %q", lang.OpDecrement, token.Text)
}

func TestOpPlusAssign(t *testing.T) {
	tw, err := NewTokenizer(reader("$a += $b")).Tokenize()
	require.NoError(t, err, "Must not be tokenization error")
	require.Equalf(t, 5, tw.Size(), "Walker must contain 3 items got %d", tw.Size())

	var token *Token
	token = tw.Get(1)
	require.Equalf(t, TokenVariable, token.Type, "expected token type %v got %v", TokenVariable, token.Type)
	require.Equalf(t, "a", token.Text, "expected token text is %q got %q", "a", token.Text)

	token = tw.Get(2)
	require.Equalf(t, TokenOperator, token.Type, "expected token type %v got %v", TokenOperator, token.Type)
	require.Equalf(t, lang.OpPlusAssign, token.Text, "expected token text is %q got %q", lang.OpPlusAssign, token.Text)

	token = tw.Get(3)
	require.Equalf(t, TokenVariable, token.Type, "expected token type %v got %v", TokenVariable, token.Type)
	require.Equalf(t, "b", token.Text, "expected token text is %q got %q", "b", token.Text)
}

func TestOpMinusAssign(t *testing.T) {
	tw, err := NewTokenizer(reader("$a -= $b")).Tokenize()
	require.NoError(t, err, "Must not be tokenization error")
	require.Equalf(t, 5, tw.Size(), "Walker must contain 3 items got %d", tw.Size())

	var token *Token
	token = tw.Get(1)
	require.Equalf(t, TokenVariable, token.Type, "expected token type %v got %v", TokenVariable, token.Type)
	require.Equalf(t, "a", token.Text, "expected token text is %q got %q", "a", token.Text)

	token = tw.Get(2)
	require.Equalf(t, TokenOperator, token.Type, "expected token type %v got %v", TokenOperator, token.Type)
	require.Equalf(t, lang.OpMinusAssign, token.Text, "expected token text is %q got %q", lang.OpMinusAssign, token.Text)

	token = tw.Get(3)
	require.Equalf(t, TokenVariable, token.Type, "expected token type %v got %v", TokenVariable, token.Type)
	require.Equalf(t, "b", token.Text, "expected token text is %q got %q", "b", token.Text)
}

func TestOpDivideAssign(t *testing.T) {
	tw, err := NewTokenizer(reader("$a /= $b")).Tokenize()
	require.NoError(t, err, "Must not be tokenization error")
	require.Equalf(t, 5, tw.Size(), "Walker must contain 3 items got %d", tw.Size())

	var token *Token
	token = tw.Get(1)
	require.Equalf(t, TokenVariable, token.Type, "expected token type %v got %v", TokenVariable, token.Type)
	require.Equalf(t, "a", token.Text, "expected token text is %q got %q", "a", token.Text)

	token = tw.Get(2)
	require.Equalf(t, TokenOperator, token.Type, "expected token type %v got %v", TokenOperator, token.Type)
	require.Equalf(t, lang.OpDivideAssign, token.Text, "expected token text is %q got %q", lang.OpDivideAssign, token.Text)

	token = tw.Get(3)
	require.Equalf(t, TokenVariable, token.Type, "expected token type %v got %v", TokenVariable, token.Type)
	require.Equalf(t, "b", token.Text, "expected token text is %q got %q", "b", token.Text)
}

func TestOpMultiplyAssign(t *testing.T) {
	tw, err := NewTokenizer(reader("$a *= $b")).Tokenize()
	require.NoError(t, err, "Must not be tokenization error")
	require.Equalf(t, 5, tw.Size(), "Walker must contain 3 items got %d", tw.Size())

	var token *Token
	token = tw.Get(1)
	require.Equalf(t, TokenVariable, token.Type, "expected token type %v got %v", TokenVariable, token.Type)
	require.Equalf(t, "a", token.Text, "expected token text is %q got %q", "a", token.Text)

	token = tw.Get(2)
	require.Equalf(t, TokenOperator, token.Type, "expected token type %v got %v", TokenOperator, token.Type)
	require.Equalf(t, lang.OpMultiplyAssign, token.Text, "expected token text is %q got %q", lang.OpMultiplyAssign, token.Text)

	token = tw.Get(3)
	require.Equalf(t, TokenVariable, token.Type, "expected token type %v got %v", TokenVariable, token.Type)
	require.Equalf(t, "b", token.Text, "expected token text is %q got %q", "b", token.Text)
}

func TestOpModuloAssign(t *testing.T) {
	tw, err := NewTokenizer(reader("$a %= $b")).Tokenize()
	require.NoError(t, err, "Must not be tokenization error")
	require.Equalf(t, 5, tw.Size(), "Walker must contain 3 items got %d", tw.Size())

	var token *Token
	token = tw.Get(1)
	require.Equalf(t, TokenVariable, token.Type, "expected token type %v got %v", TokenVariable, token.Type)
	require.Equalf(t, "a", token.Text, "expected token text is %q got %q", "a", token.Text)

	token = tw.Get(2)
	require.Equalf(t, TokenOperator, token.Type, "expected token type %v got %v", TokenOperator, token.Type)
	require.Equalf(t, lang.OpModuloAssign, token.Text, "expected token text is %q got %q", lang.OpModuloAssign, token.Text)

	token = tw.Get(3)
	require.Equalf(t, TokenVariable, token.Type, "expected token type %v got %v", TokenVariable, token.Type)
	require.Equalf(t, "b", token.Text, "expected token text is %q got %q", "b", token.Text)
}

// Testing syntax punctuation

func TestOpenParen(t *testing.T) {
	tw, err := NewTokenizer(reader("(")).Tokenize()
	require.NoError(t, err, "Must not be tokenization error")
	require.Equalf(t, 3, tw.Size(), "Walker must contain 3 items got %d", tw.Size())
	token := tw.Get(1)
	require.Equalf(t, TokenOpenParen, token.Type, "expected token type %v got %v", TokenOpenParen, token.Type)
	require.Equalf(t, "(", token.Text, "expected token text is %q got %q", "(", token.Text)
}

func TestCloseParen(t *testing.T) {
	tw, err := NewTokenizer(reader(")")).Tokenize()
	require.NoError(t, err, "Must not be tokenization error")
	require.Equalf(t, 3, tw.Size(), "Walker must contain 3 items got %d", tw.Size())
	token := tw.Get(1)
	require.Equalf(t, TokenCloseParen, token.Type, "expected token type %v got %v", TokenCloseParen, token.Type)
	require.Equalf(t, ")", token.Text, "expected token text is %q got %q", ")", token.Text)
}

func TestOpenBracket(t *testing.T) {
	tw, err := NewTokenizer(reader("[")).Tokenize()
	require.NoError(t, err, "Must not be tokenization error")
	require.Equalf(t, 3, tw.Size(), "Walker must contain 3 items got %d", tw.Size())
	token := tw.Get(1)
	require.Equalf(t, TokenOpenBracket, token.Type, "expected token type %v got %v", TokenOpenBracket, token.Type)
	require.Equalf(t, "[", token.Text, "expected token text is %q got %q", "[", token.Text)
}

func TestCloseBracket(t *testing.T) {
	tw, err := NewTokenizer(reader("]")).Tokenize()
	require.NoError(t, err, "Must not be tokenization error")
	require.Equalf(t, 3, tw.Size(), "Walker must contain 3 items got %d", tw.Size())
	token := tw.Get(1)
	require.Equalf(t, TokenCloseBracket, token.Type, "expected token type %v got %v", TokenCloseBracket, token.Type)
	require.Equalf(t, "]", token.Text, "expected token text is %q got %q", "]", token.Text)
}

func TestOpenBrace(t *testing.T) {
	tw, err := NewTokenizer(reader("{")).Tokenize()
	require.NoError(t, err, "Must not be tokenization error")
	require.Equalf(t, 3, tw.Size(), "Walker must contain 3 items got %d", tw.Size())
	token := tw.Get(1)
	require.Equalf(t, TokenOpenBrace, token.Type, "expected token type %v got %v", TokenOpenBrace, token.Type)
	require.Equalf(t, "{", token.Text, "expected token text is %q got %q", "{", token.Text)
}

func TestCloseBrace(t *testing.T) {
	tw, err := NewTokenizer(reader("}")).Tokenize()
	require.NoError(t, err, "Must not be tokenization error")
	require.Equalf(t, 3, tw.Size(), "Walker must contain 3 items got %d", tw.Size())
	token := tw.Get(1)
	require.Equalf(t, TokenCloseBrace, token.Type, "expected token type %v got %v", TokenCloseBrace, token.Type)
	require.Equalf(t, "}", token.Text, "expected token text is %q got %q", "}", token.Text)
}

func TestColon(t *testing.T) {
	tw, err := NewTokenizer(reader(":")).Tokenize()
	require.NoError(t, err, "Must not be tokenization error")
	require.Equalf(t, 3, tw.Size(), "Walker must contain 3 items got %d", tw.Size())
	token := tw.Get(1)
	require.Equalf(t, TokenColon, token.Type, "expected token type %v got %v", TokenColon, token.Type)
	require.Equalf(t, ":", token.Text, "expected token text is %q got %q", ":", token.Text)
}

func TestSemicolon(t *testing.T) {
	tw, err := NewTokenizer(reader(";")).Tokenize()
	require.NoError(t, err, "Must not be tokenization error")
	require.Equalf(t, 3, tw.Size(), "Walker must contain 3 items got %d", tw.Size())
	token := tw.Get(1)
	require.Equalf(t, TokenSemicolon, token.Type, "expected token type %v got %v", TokenSemicolon, token.Type)
	require.Equalf(t, ";", token.Text, "expected token text is %q got %q", ";", token.Text)
}

func TestComa(t *testing.T) {
	tw, err := NewTokenizer(reader(",")).Tokenize()
	require.NoError(t, err, "Must not be tokenization error")
	require.Equalf(t, 3, tw.Size(), "Walker must contain 3 items got %d", tw.Size())
	token := tw.Get(1)
	require.Equalf(t, TokenComa, token.Type, "expected token type %v got %v", TokenComa, token.Type)
	require.Equalf(t, ",", token.Text, "expected token text is %q got %q", ",", token.Text)
}

func TestPoint(t *testing.T) {
	tw, err := NewTokenizer(reader(".")).Tokenize()
	require.NoError(t, err, "Must not be tokenization error")
	require.Equalf(t, 3, tw.Size(), "Walker must contain 3 items got %d", tw.Size())
	token := tw.Get(1)
	require.Equalf(t, TokenPoint, token.Type, "expected token type %v got %v", TokenPoint, token.Type)
	require.Equalf(t, ".", token.Text, "expected token text is %q got %q", ".", token.Text)
}

func TestAt(t *testing.T) {
	tw, err := NewTokenizer(reader("@")).Tokenize()
	require.NoError(t, err, "Must not be tokenization error")
	require.Equalf(t, 3, tw.Size(), "Walker must contain 3 items got %d", tw.Size())
	token := tw.Get(1)
	require.Equalf(t, TokenAt, token.Type, "expected token type %v got %v", TokenAt, token.Type)
	require.Equalf(t, "@", token.Text, "expected token text is %q got %q", "@", token.Text)
}

func TestAssignment(t *testing.T) {
	tw, err := NewTokenizer(reader("=")).Tokenize()
	require.NoError(t, err, "Must not be tokenization error")
	require.Equalf(t, 3, tw.Size(), "Walker must contain 3 items got %d", tw.Size())
	token := tw.Get(1)
	require.Equalf(t, TokenAssignment, token.Type, "expected token type %v got %v", TokenAssignment, token.Type)
	require.Equalf(t, "=", token.Text, "expected token text is %q got %q", "=", token.Text)
}

func TestArrow(t *testing.T) {
	tw, err := NewTokenizer(reader("=>")).Tokenize()
	require.NoError(t, err, "Must not be tokenization error")
	require.Equalf(t, 3, tw.Size(), "Walker must contain 3 items got %d", tw.Size())
	token := tw.Get(1)
	require.Equalf(t, TokenArrow, token.Type, "expected token type %v got %v", TokenArrow, token.Type)
	require.Equalf(t, "=>", token.Text, "expected token text is %q got %q", "=>", token.Text)
}

// Testing variable

func TestVariable(t *testing.T) {
	tw, err := NewTokenizer(reader("$word")).Tokenize()
	require.NoError(t, err, "Must not be tokenization error")
	require.Equalf(t, 3, tw.Size(), "Walker must contain 3 items got %d", tw.Size())
	token := tw.Get(1)
	require.Equalf(t, TokenVariable, token.Type, "expected token type %v got %v", TokenVariable, token.Type)
	require.Equalf(t, "word", token.Text, "expected token text is %q got %q", "word", token.Text)
}

func TestVariableWithUnderscore(t *testing.T) {
	tw, err := NewTokenizer(reader("$word_word_word")).Tokenize()
	require.NoError(t, err, "Must not be tokenization error")
	require.Equalf(t, 3, tw.Size(), "Walker must contain 3 items got %d", tw.Size())
	token := tw.Get(1)
	require.Equalf(t, TokenVariable, token.Type, "expected token type %v got %v", TokenVariable, token.Type)
	require.Equalf(t, "word_word_word", token.Text, "expected token text is %q got %q", "word_word_word", token.Text)
}

func TestVariableWithNumber(t *testing.T) {
	tw, err := NewTokenizer(reader("$word00203word")).Tokenize()
	require.NoError(t, err, "Must not be tokenization error")
	require.Equalf(t, 3, tw.Size(), "Walker must contain 3 items got %d", tw.Size())
	token := tw.Get(1)
	require.Equalf(t, TokenVariable, token.Type, "expected token type %v got %v", TokenVariable, token.Type)
	require.Equalf(t, "word00203word", token.Text, "expected token text is %q got %q", "word00203word", token.Text)
}

func TestVariableStartingFromNumber(t *testing.T) {
	tw, err := NewTokenizer(reader("$203word")).Tokenize()
	require.NoError(t, err, "Must not be tokenization error")
	require.Equalf(t, 3, tw.Size(), "Walker must contain 3 items got %d", tw.Size())
	token := tw.Get(1)
	require.Equalf(t, TokenVariable, token.Type, "expected token type %v got %v", TokenVariable, token.Type)
	require.Equalf(t, "203word", token.Text, "expected token text is %q got %q", "203word", token.Text)
}

func TestVariableWithEverything(t *testing.T) {
	tw, err := NewTokenizer(reader("$word00203word_s")).Tokenize()
	require.NoError(t, err, "Must not be tokenization error")
	require.Equalf(t, 3, tw.Size(), "Walker must contain 3 items got %d", tw.Size())
	token := tw.Get(1)
	require.Equalf(t, TokenVariable, token.Type, "expected token type %v got %v", TokenVariable, token.Type)
	require.Equalf(t, "word00203word_s", token.Text, "expected token text is %q got %q", "word00203word_s", token.Text)
}

func TestEOF(t *testing.T) {
	tw, err := NewTokenizer(reader("some_thing")).Tokenize()
	require.NoError(t, err, "Must not be tokenization error")
	require.Equalf(t, 3, tw.Size(), "Walker must contain 3 items got %d", tw.Size())
	token := tw.Get(2)
	require.Equalf(t, TokenEOF, token.Type, "expected token type %v got %v", TokenEOF, token.Type)
}
