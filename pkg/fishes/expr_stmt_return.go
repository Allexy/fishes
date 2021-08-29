package fishes

type statementReturn struct {
	tobeReturned expression
}

func (s statementReturn) evaluate() (Value, bool, error) {
	v, _, err := s.tobeReturned.evaluate()
	if err != nil {
		return nil, false, err
	}
	return v, true, nil
}
