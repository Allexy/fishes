package fishes

type expression interface {
	// evaluate returns computed Value and boolean flag which is showing that result must be bypassed to
	//  parental scope
	evaluate() (Value, bool, error)
}

type expressions struct {
	expressions []expression
}

func newExpressions() expressions {
	return expressions{
		expressions: make([]expression, 0, 100),
	}
}

func (e expressions) execute() (Value, bool, error) {
	if len(e.expressions) == 0 {
		return predefinedNull, false, nil
	}
	var (
		result Value
		err    error
		bypass bool
	)
	for _, v := range e.expressions {
		result, bypass, err = v.evaluate()
		if err != nil {
			return predefinedNull, false, err
		}
		if bypass {
			return result, true, nil
		}
	}
	return result, false, nil
}

func (e *expressions) add(expr expression) {
	e.expressions = append(e.expressions, expr)
}

func (e *expressions) merge(a expressions) {
	for _, expr := range a.expressions {
		e.add(expr)
	}
}
