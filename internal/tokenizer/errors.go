package tokenizer

import "fmt"

type TokenizerError struct {
	fileName string
	message  string
	line     uint32
	col      uint32
}

func NewTokenizerError(fileName string, message string, line uint32, col uint32) *TokenizerError {
	return &TokenizerError{fileName, message, line, col}
}

func (te TokenizerError) Error() string {
	return fmt.Sprintf("Error in file %s: %s\nAt line %d; col: %d", te.fileName, te.message, te.line, te.col)
}
