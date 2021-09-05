package fishes

type statementFor struct {
	initialization expression
	condition      expression
	iteration      expression
	body           expressions
}

func (s statementFor) evaluate() (Value, bool, error) {
	var (
		c, b   Value
		err    error
		bypass bool
	)
	if _, _, err = s.initialization.evaluate(); err != nil {
		return nil, false, err
	}
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
