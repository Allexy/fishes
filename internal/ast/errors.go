package ast

import "fmt"

type parserError struct {
	sourceName string
	message    string
	line       uint32
	col        uint32
	cause      error
}

func newParserError(sourceName string, message string, line uint32, col uint32, previous error) parserError {
	return parserError{sourceName, message, line, col, previous}
}

func (pe parserError) Error() string {
	return fmt.Sprintf("Parse error in file %s: %s\nAt line %d; col: %d", pe.sourceName, pe.message, pe.line, pe.col)
}

func (pe parserError) Cause() error {
	return pe.cause
}
