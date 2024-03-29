package tokenizer

import "fmt"

type TokenizerError struct {
	fileName string
	message  string
	line     uint32
	col      uint32
	cause    error
}

func NewTokenizerError(fileName string, message string, line uint32, col uint32, previous error) TokenizerError {
	return TokenizerError{fileName, message, line, col, previous}
}

func (te TokenizerError) Error() string {
	return fmt.Sprintf("Error in file %s: %s\nAt line %d; col: %d", te.fileName, te.message, te.line, te.col)
}

func (te TokenizerError) Cause() error {
	return te.cause
}
