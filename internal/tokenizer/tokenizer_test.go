package tokenizer

import (
	"io"
	"strings"
	"testing"
)

// Creates and returns string reader and source name
func _mk(s string) (io.Reader, string) {
	return strings.NewReader(s), "string"
}

// Testing walker

func TestTokenWalker(t *testing.T) {
	tokens := []Token{{}, {}, {}}
	tw := NewTokenWalker(tokens)
	if tw.Size() != 3 {
		t.Errorf("Expected size 3 but got %d", tw.Size())
	}
	var step uint = 0
	for ; tw.Next(); step++ {
		tw.Move(1)
		if step > 2 {
			t.Errorf("Expected that walker iterate by 3 steps, current step: %d", step)
		}
	}
}

func TestTokenWalkerMatch(t *testing.T) {
	// Simulate typical for function definition sequence of tokens
	tokens := []Token{
		{Token: TokenBOF},
		{Token: TokenWord}, {Token: TokenWord}, {Token: TokenOpenParen}, {Token: TokenCloseParen},
		{Token: TokenOpenBrace}, {Token: TokenOpenBrace}, {Token: TokenCloseBrace}, {Token: TokenWord},
		{Token: TokenOpenBrace}, {Token: TokenCloseBrace}, {Token: TokenWord}, {Token: TokenCloseBrace},
		{Token: TokenEOF},
	}
	tw := NewTokenWalker(tokens)
	tw.Move(1)
	// Any tokens between prentices and braces are matched by TT_DEFAULT, walker takes into consederation nested
	// opening and closing prentices
	if !tw.Match(TokenWord, TokenWord, TokenOpenParen, TokenDefault, TokenCloseParen, TokenOpenBrace, TokenDefault, TokenCloseBrace) {
		t.Error("Expected that sequence is matching patter")
	}
}

// Testing words

func TestWord(t *testing.T) {
	tw, err := NewTokenizer(_mk("word")).Tokenize()
	if err != nil {
		t.Errorf("Tokenization failed with err: %v", err)
	}
	if tw.Size() != 3 {
		t.Errorf("Expected 3 tokens in result got %d", tw.Size())
	}
	token := tw.Get(1)
	if token.Token != TokenWord {
		t.Errorf("Expected token of type TT_WORD but got %q", token)
	}
	if token.Text != "word" {
		t.Errorf("Expected token text is 'word' but got %q", token)
	}
}

func TestWordWithUnderscore(t *testing.T) {
	tw, err := NewTokenizer(_mk("word_word_word")).Tokenize()
	if err != nil {
		t.Errorf("Tokenization failed with err: %v", err)
	}
	if tw.Size() != 3 {
		t.Errorf("Expected 3 tokens in result got %d", tw.Size())
	}
	token := tw.Get(1)
	if token.Token != TokenWord {
		t.Errorf("Expected token of type TT_WORD but got %q", token)
	}
	if token.Text != "word_word_word" {
		t.Errorf("Expected token text is 'word_word_word' but got %q", token)
	}
}

func TestWordWithNumber(t *testing.T) {
	tw, err := NewTokenizer(_mk("word00203word")).Tokenize()
	if err != nil {
		t.Errorf("Tokenization failed with err: %v", err)
	}
	if tw.Size() != 3 {
		t.Errorf("Expected 3 tokens in result got %d", tw.Size())
	}
	token := tw.Get(1)
	if token.Token != TokenWord {
		t.Errorf("Expected token of type TT_WORD but got %q", token)
	}
	if token.Text != "word00203word" {
		t.Errorf("Expected token text is 'word00203word' but got %q", token)
	}
}

func TestWordWithEverything(t *testing.T) {
	tw, err := NewTokenizer(_mk("word00203word_s")).Tokenize()
	if err != nil {
		t.Errorf("Tokenization failed with err: %v", err)
	}
	if tw.Size() != 3 {
		t.Errorf("Expected 3 tokens in result got %d", tw.Size())
	}
	token := tw.Get(1)
	if token.Token != TokenWord {
		t.Errorf("Expected token of type TT_WORD but got %q", token)
	}
	if token.Text != "word00203word_s" {
		t.Errorf("Expected token text is 'word00203word_s' but got %q", token)
	}
}

// Testing numerical literals

func TestNumberSimple(t *testing.T) {
	tw, err := NewTokenizer(_mk("10")).Tokenize()
	if err != nil {
		t.Errorf("Tokenization failed with err: %v", err)
	}
	if tw.Size() != 3 {
		t.Errorf("Expected 3 tokens in result got %d", tw.Size())
	}
	token := tw.Get(1)
	if token.Token != TokenNumber {
		t.Errorf("Expected token of type TT_NUMBER but got %q", token)
	}
}

func TestNumberWithPointInMiddle(t *testing.T) {
	tw, err := NewTokenizer(_mk("10.5")).Tokenize()
	if err != nil {
		t.Errorf("Tokenization failed with err: %v", err)
	}
	if tw.Size() != 3 {
		t.Errorf("Expected 3 tokens in result got %d", tw.Size())
	}
	token := tw.Get(1)
	if token.Token != TokenNumber {
		t.Errorf("Expected token of type TT_NUMBER but got %q", token)
	}
	if token.Text != "10.5" {
		t.Errorf("Expected token text is '10.0' but got %q", token)
	}
}

func TestNumberWithLeadingPoint(t *testing.T) {
	tw, err := NewTokenizer(_mk(".5")).Tokenize()
	if err != nil {
		t.Errorf("Tokenization failed with err: %v", err)
	}
	if tw.Size() != 3 {
		t.Errorf("Expected 3 tokens in result got %d", tw.Size())
	}
	token := tw.Get(1)
	if token.Token != TokenNumber {
		t.Errorf("Expected token of type TT_NUMBER but got %q", token)
	}
	if token.Text != "0.5" {
		t.Errorf("Expected token text is '0.5' but got %q", token)
	}
}

func TestNumberWithTerminatingPoint(t *testing.T) {
	tw, err := NewTokenizer(_mk("5.")).Tokenize()
	if err != nil {
		t.Errorf("Tokenization failed with err: %v", err)
	}
	if tw.Size() != 3 {
		t.Errorf("Expected 3 tokens in result got %d", tw.Size())
	}
	token := tw.Get(1)
	if token.Token != TokenNumber {
		t.Errorf("Expected token of type TT_NUMBER but got %q", token)
	}
	if token.Text != "5.0" {
		t.Errorf("Expected token text is '5.0' but got %q", token)
	}
}

func TestNegativeNumberSimple(t *testing.T) {
	tw, err := NewTokenizer(_mk("-10")).Tokenize()
	if err != nil {
		t.Errorf("Tokenization failed with err: %v", err)
	}
	if tw.Size() != 3 {
		t.Errorf("Expected 3 tokens in result got %d", tw.Size())
	}
	token := tw.Get(1)
	if token.Token != TokenNumber {
		t.Errorf("Expected token of type TT_NUMBER but got %q", token)
	}
}

func TestNegativeNumberWithPointInMiddle(t *testing.T) {
	tw, err := NewTokenizer(_mk("-10.5")).Tokenize()
	if err != nil {
		t.Errorf("Tokenization failed with err: %v", err)
	}
	if tw.Size() != 3 {
		t.Errorf("Expected 3 tokens in result got %d", tw.Size())
	}
	token := tw.Get(1)
	if token.Token != TokenNumber {
		t.Errorf("Expected token of type TT_NUMBER but got %q", token)
	}
	if token.Text != "-10.5" {
		t.Errorf("Expected token text is '-10.5' but got %q", token)
	}
}

func TestNegativeNumberWithLeadingPoint(t *testing.T) {
	tw, err := NewTokenizer(_mk("-.5")).Tokenize()
	if err != nil {
		t.Errorf("Tokenization failed with err: %v", err)
	}
	if tw.Size() != 3 {
		t.Errorf("Expected 3 tokens in result got %d", tw.Size())
	}
	token := tw.Get(1)
	if token.Token != TokenNumber {
		t.Errorf("Expected token of type TT_NUMBER but got %q", token)
	}
	if token.Text != "-0.5" {
		t.Errorf("Expected token text is '-0.5' but got %q", token)
	}
}

func TestNegativeNumberWithTerminatingPoint(t *testing.T) {
	tw, err := NewTokenizer(_mk("-5.")).Tokenize()
	if err != nil {
		t.Errorf("Tokenization failed with err: %v", err)
	}
	if tw.Size() != 3 {
		t.Errorf("Expected 3 tokens in result got %d", tw.Size())
	}
	token := tw.Get(1)
	if token.Token != TokenNumber {
		t.Errorf("Expected token of type TT_NUMBER but got %q", token)
	}
	if token.Text != "-5.0" {
		t.Errorf("Expected token text is '-5.0' but got %q", token)
	}
}

// Testing string literals

func TestString(t *testing.T) {
	tw, err := NewTokenizer(_mk("\"Test string\"")).Tokenize()
	if err != nil {
		t.Errorf("Tokenization failed with err: %v", err)
	}
	if tw.Size() != 3 {
		t.Errorf("Expected 3 tokens in result got %d", tw.Size())
	}
	token := tw.Get(1)
	if token.Token != TokenString {
		t.Errorf("Expected token of type TT_STRING but got %q", token)
	}
	if token.Text != "Test string" {
		t.Errorf("Expected token text is 'Test string' but got %q", token)
	}
}

func TestStringWithEscSeq(t *testing.T) {
	tw, err := NewTokenizer(_mk("\"Test\\nstring\"")).Tokenize()
	if err != nil {
		t.Errorf("Tokenization failed with err: %v", err)
	}
	if tw.Size() != 3 {
		t.Errorf("Expected 3 tokens in result got %d", tw.Size())
	}
	token := tw.Get(1)
	if token.Token != TokenString {
		t.Errorf("Expected token of type TT_STRING but got %q", token)
	}
	if token.Text != "Test\nstring" {
		t.Errorf("Expected token text is 'Test\\nstring' but got %q", token)
	}
}

func TestStringWithNestedQuotes(t *testing.T) {
	tw, err := NewTokenizer(_mk("\"Test \\\"string\\\"\"")).Tokenize()
	if err != nil {
		t.Errorf("Tokenization failed with err: %v", err)
	}
	if tw.Size() != 3 {
		t.Errorf("Expected 3 tokens in result got %d", tw.Size())
	}
	token := tw.Get(1)
	if token.Token != TokenString {
		t.Errorf("Expected token of type TT_STRING but got %q", token)
	}
	if token.Text != "Test \"string\"" {
		t.Errorf("Expected token text is 'Test \\\"string\\\"' but got %q", token)
	}
}

// Testing logical literals

func TestLogicalTrue(t *testing.T) {
	tw, err := NewTokenizer(_mk("true")).Tokenize()
	if err != nil {
		t.Errorf("Tokenization failed with err: %v", err)
	}
	if tw.Size() != 3 {
		t.Errorf("Expected 3 tokens in result got %d", tw.Size())
	}
	token := tw.Get(1)
	if token.Token != TokenLogic {
		t.Errorf("Expected token of type TT_LOGIC but got %q", token)
	}
	if token.Text != "true" {
		t.Errorf("Expected token text is 'true' but got %q", token)
	}
}

func TestLogicalFalse(t *testing.T) {
	tw, err := NewTokenizer(_mk("false")).Tokenize()
	if err != nil {
		t.Errorf("Tokenization failed with err: %v", err)
	}
	if tw.Size() != 3 {
		t.Errorf("Expected 3 tokens in result got %d", tw.Size())
	}
	token := tw.Get(1)
	if token.Token != TokenLogic {
		t.Errorf("Expected token of type TT_LOGIC but got %q", token)
	}
	if token.Text != "false" {
		t.Errorf("Expected token text is 'false' but got %q", token)
	}
}

// Testing operators (I'll do this latter)

// Testing syntax punctuation

func TestOpenParen(t *testing.T) {
	tw, err := NewTokenizer(_mk("(")).Tokenize()
	if err != nil {
		t.Errorf("Tokenization failed with err: %v", err)
	}
	if tw.Size() != 3 {
		t.Errorf("Expected 3 tokens in result got %d", tw.Size())
	}
	token := tw.Get(1)
	if token.Token != TokenOpenParen {
		t.Errorf("Expected token of type TT_O_PAREN but got %q", token)
	}
	if token.Text != "(" {
		t.Errorf("Expected token text is '(' but got %q", token)
	}
}

func TestCloseParen(t *testing.T) {
	tw, err := NewTokenizer(_mk(")")).Tokenize()
	if err != nil {
		t.Errorf("Tokenization failed with err: %v", err)
	}
	if tw.Size() != 3 {
		t.Errorf("Expected 3 tokens in result got %d", tw.Size())
	}
	token := tw.Get(1)
	if token.Token != TokenCloseParen {
		t.Errorf("Expected token of type TT_C_PAREN but got %q", token)
	}
	if token.Text != ")" {
		t.Errorf("Expected token text is ')' but got %q", token)
	}
}

func TestOpenBracket(t *testing.T) {
	tw, err := NewTokenizer(_mk("[")).Tokenize()
	if err != nil {
		t.Errorf("Tokenization failed with err: %v", err)
	}
	if tw.Size() != 3 {
		t.Errorf("Expected 3 tokens in result got %d", tw.Size())
	}
	token := tw.Get(1)
	if token.Token != TokenOpenBracket {
		t.Errorf("Expected token of type TT_O_BRACKET but got %q", token)
	}
	if token.Text != "[" {
		t.Errorf("Expected token text is '[' but got %q", token)
	}
}

func TestCloseBracket(t *testing.T) {
	tw, err := NewTokenizer(_mk("]")).Tokenize()
	if err != nil {
		t.Errorf("Tokenization failed with err: %v", err)
	}
	if tw.Size() != 3 {
		t.Errorf("Expected 3 tokens in result got %d", tw.Size())
	}
	token := tw.Get(1)
	if token.Token != TokenCloseBracket {
		t.Errorf("Expected token of type TT_C_BRACKET but got %q", token)
	}
	if token.Text != "]" {
		t.Errorf("Expected token text is ']' but got %q", token)
	}
}

func TestOpenBrace(t *testing.T) {
	tw, err := NewTokenizer(_mk("{")).Tokenize()
	if err != nil {
		t.Errorf("Tokenization failed with err: %v", err)
	}
	if tw.Size() != 3 {
		t.Errorf("Expected 3 tokens in result got %d", tw.Size())
	}
	token := tw.Get(1)
	if token.Token != TokenOpenBrace {
		t.Errorf("Expected token of type TT_O_BRACE but got %q", token)
	}
	if token.Text != "{" {
		t.Errorf("Expected token text is '{' but got %q", token)
	}
}

func TestCloseBrace(t *testing.T) {
	tw, err := NewTokenizer(_mk("}")).Tokenize()
	if err != nil {
		t.Errorf("Tokenization failed with err: %v", err)
	}
	if tw.Size() != 3 {
		t.Errorf("Expected 3 tokens in result got %d", tw.Size())
	}
	token := tw.Get(1)
	if token.Token != TokenCloseBrace {
		t.Errorf("Expected token of type TT_C_BRACE but got %q", token)
	}
	if token.Text != "}" {
		t.Errorf("Expected token text is '}' but got %q", token)
	}
}

func TestColon(t *testing.T) {
	tw, err := NewTokenizer(_mk(":")).Tokenize()
	if err != nil {
		t.Errorf("Tokenization failed with err: %v", err)
	}
	if tw.Size() != 3 {
		t.Errorf("Expected 3 tokens in result got %d", tw.Size())
	}
	token := tw.Get(1)
	if token.Token != TokenColon {
		t.Errorf("Expected token of type TT_COLON but got %q", token)
	}
	if token.Text != ":" {
		t.Errorf("Expected token text is ':' but got %q", token)
	}
}

func TestSemicolon(t *testing.T) {
	tw, err := NewTokenizer(_mk(";")).Tokenize()
	if err != nil {
		t.Errorf("Tokenization failed with err: %v", err)
	}
	if tw.Size() != 3 {
		t.Errorf("Expected 3 tokens in result got %d", tw.Size())
	}
	token := tw.Get(1)
	if token.Token != TokenSemicolon {
		t.Errorf("Expected token of type TT_SEMICOLON but got %q", token)
	}
	if token.Text != ";" {
		t.Errorf("Expected token text is ';' but got %q", token)
	}
}

func TestComa(t *testing.T) {
	tw, err := NewTokenizer(_mk(",")).Tokenize()
	if err != nil {
		t.Errorf("Tokenization failed with err: %v", err)
	}
	if tw.Size() != 3 {
		t.Errorf("Expected 3 tokens in result got %d", tw.Size())
	}
	token := tw.Get(1)
	if token.Token != TokenComa {
		t.Errorf("Expected token of type TT_SEMICOLON but got %q", token)
	}
	if token.Text != "," {
		t.Errorf("Expected token text is ',' but got %q", token)
	}
}

func TestPoint(t *testing.T) {
	tw, err := NewTokenizer(_mk(".")).Tokenize()
	if err != nil {
		t.Errorf("Tokenization failed with err: %v", err)
	}
	if tw.Size() != 3 {
		t.Errorf("Expected 3 tokens in result got %d", tw.Size())
	}
	token := tw.Get(1)
	if token.Token != TokenPoint {
		t.Errorf("Expected token of type TT_SEMICOLON but got %q", token)
	}
	if token.Text != "." {
		t.Errorf("Expected token text is '.' but got %q", token)
	}
}

func TestAt(t *testing.T) {
	tw, err := NewTokenizer(_mk("@")).Tokenize()
	if err != nil {
		t.Errorf("Tokenization failed with err: %v", err)
	}
	if tw.Size() != 3 {
		t.Errorf("Expected 3 tokens in result got %d", tw.Size())
	}
	token := tw.Get(1)
	if token.Token != TokenAt {
		t.Errorf("Expected token of type TT_AT but got %q", token)
	}
	if token.Text != "@" {
		t.Errorf("Expected token text is '@' but got %q", token)
	}
}

func TestAssignment(t *testing.T) {
	tw, err := NewTokenizer(_mk("=")).Tokenize()
	if err != nil {
		t.Errorf("Tokenization failed with err: %v", err)
	}
	if tw.Size() != 3 {
		t.Errorf("Expected 3 tokens in result got %d", tw.Size())
	}
	token := tw.Get(1)
	if token.Token != TokenAssignment {
		t.Errorf("Expected token of type TT_ASSIGNMENT but got %q", token)
	}
	if token.Text != "=" {
		t.Errorf("Expected token text is '=' but got %q", token)
	}
}

func TestArrow(t *testing.T) {
	tw, err := NewTokenizer(_mk("=>")).Tokenize()
	if err != nil {
		t.Errorf("Tokenization failed with err: %v", err)
	}
	if tw.Size() != 3 {
		t.Errorf("Expected 3 tokens in result got %d", tw.Size())
	}
	token := tw.Get(1)
	if token.Token != TokenArrow {
		t.Errorf("Expected token of type TT_ARROW but got %q", token)
	}
	if token.Text != "=>" {
		t.Errorf("Expected token text is '=>' but got %q", token)
	}
}

// Testing variable

func TestVariable(t *testing.T) {
	tw, err := NewTokenizer(_mk("$word")).Tokenize()
	if err != nil {
		t.Errorf("Tokenization failed with err: %v", err)
	}
	if tw.Size() != 3 {
		t.Errorf("Expected 3 tokens in result got %d", tw.Size())
	}
	token := tw.Get(1)
	if token.Token != TokenVariable {
		t.Errorf("Expected token of type TT_VARIABLE but got %q", token)
	}
	if token.Text != "word" {
		t.Errorf("Expected token text is 'word' but got %q", token)
	}
}

func TestVariableWithUnderscore(t *testing.T) {
	tw, err := NewTokenizer(_mk("$word_word_word")).Tokenize()
	if err != nil {
		t.Errorf("Tokenization failed with err: %v", err)
	}
	if tw.Size() != 3 {
		t.Errorf("Expected 3 tokens in result got %d", tw.Size())
	}
	token := tw.Get(1)
	if token.Token != TokenVariable {
		t.Errorf("Expected token of type TT_VARIABLE but got %q", token)
	}
	if token.Text != "word_word_word" {
		t.Errorf("Expected token text is 'word_word_word' but got %q", token)
	}
}

func TestVariableWithNumber(t *testing.T) {
	tw, err := NewTokenizer(_mk("$word00203word")).Tokenize()
	if err != nil {
		t.Errorf("Tokenization failed with err: %v", err)
	}
	if tw.Size() != 3 {
		t.Errorf("Expected 3 tokens in result got %d", tw.Size())
	}
	token := tw.Get(1)
	if token.Token != TokenVariable {
		t.Errorf("Expected token of type TT_VARIABLE but got %q", token)
	}
	if token.Text != "word00203word" {
		t.Errorf("Expected token text is 'word00203word' but got %q", token)
	}
}

func TestVariableWithEverything(t *testing.T) {
	tw, err := NewTokenizer(_mk("$word00203word_s")).Tokenize()
	if err != nil {
		t.Errorf("Tokenization failed with err: %v", err)
	}
	if tw.Size() != 3 {
		t.Errorf("Expected 3 tokens in result got %d", tw.Size())
	}
	token := tw.Get(1)
	if token.Token != TokenVariable {
		t.Errorf("Expected token of type TT_VARIABLE but got %q", token)
	}
	if token.Text != "word00203word_s" {
		t.Errorf("Expected token text is 'word00203word_s' but got %q", token)
	}
}
