package fishes

type statementCase struct {
	condition expression
	body      expressions
}

type statementSwitch struct {
	condition expression
	cases     []statementCase
}

func (s statementSwitch) evaluate() (Value, bool, error) {
	var (
		c   Value
		err error
	)
	if c, _, err = s.condition.evaluate(); err != nil {
		return nil, false, err
	}
	for _, sc := range s.cases {
		scc, _, err := sc.condition.evaluate()
		if err != nil {
			return nil, false, err
		}
		if c.Equals(scc) {
			res, bypass, err := sc.body.execute()
			if err != nil {
				return nil, false, err
			}
			if bypass {
				return res, true, nil
			}
		}
	}
	return predefinedNull, false, nil
}

func (s *statementSwitch) addCase(c statementCase) {
	s.cases = append(s.cases, c)
}
