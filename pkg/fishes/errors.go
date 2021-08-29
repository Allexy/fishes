package fishes

import "fmt"

type parseError struct {
	sourceName string
	message    string
	line       uint32
	col        uint32
	cause      error
}

func newParseError(sourceName string, message string, line uint32, col uint32, previous error) parseError {
	return parseError{sourceName, message, line, col, previous}
}

func (e parseError) Error() string {
	return fmt.Sprintf("parse error in %s: %s; line %d; col: %d", e.sourceName, e.message, e.line, e.col)
}

func (e parseError) Cause() error {
	return e.cause
}

type RuntimeError struct {
	errorType  string
	sourceName string
	message    string
	line       uint32
	col        uint32
	cause      error
}

func newRuntimeError(errorType string, sourceName string, message string, line uint32, col uint32, previous error) RuntimeError {
	return RuntimeError{errorType, sourceName, message, line, col, previous}
}

func (e RuntimeError) Error() string {
	return fmt.Sprintf("parse error in %s: %s; line %d; col: %d", e.sourceName, e.message, e.line, e.col)
}

func (e RuntimeError) Cause() error {
	return e.cause
}

func (e RuntimeError) join(err error) {
	e.cause = err
}
