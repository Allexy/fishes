package fishes

type statementWhile struct {
	condition expression
	body      expressions
}

func (s statementWhile) evaluate() (Value, bool, error) {
	var (
		c, b   Value
		err    error
		bypass bool
	)
	for {
		c, _, err = s.condition.evaluate()
		if err != nil {
			return nil, false, err
		}
		if c.AsBoolean() {
			b, bypass, err = s.body.execute()
			if err != nil {
				return nil, false, err
			}
			if bypass {
				return b, true, nil
			}
			continue
		}
		return predefinedNull, false, nil
	}
}
