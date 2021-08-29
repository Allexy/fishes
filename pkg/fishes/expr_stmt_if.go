package fishes

type statementIf struct {
	condition expression
	body      expressions
}

func (s statementIf) evaluate() (Value, bool, error) {
	cv, _, err := s.condition.evaluate()
	if err != nil {
		return nil, false, err
	}
	if cv.AsBoolean() {
		return s.body.execute()
	}
	return predefinedNull, false, err
}
