package fishes

import "fmt"

type statementThrow struct {
	source        string
	line, col     uint32
	kind, message expression
}

func (s *statementThrow) evaluate() (Value, bool, error) {
	kind, _, err := s.kind.evaluate()
	if err != nil {
		return nil, false, err
	}
	if kind.Type() != ValueString {
		return nil, false, newRuntimeError(
			"TypeMismatch",
			s.source,
			fmt.Sprintf("expected type of kind is string, got %q", kind.Type()),
			s.line,
			s.col,
			nil,
		)
	}
	message, _, err := s.message.evaluate()
	if err != nil {
		return nil, false, err
	}
	if message.Type() != ValueString {
		return nil, false, newRuntimeError(
			"TypeMismatch",
			s.source,
			fmt.Sprintf("expected type of message is string, got %q", message.Type()),
			s.line,
			s.col,
			nil,
		)
	}
	return nil, false, newRuntimeError(
		kind.AsString(),
		s.source,
		message.AsString(),
		s.line,
		s.col,
		nil,
	)
}
