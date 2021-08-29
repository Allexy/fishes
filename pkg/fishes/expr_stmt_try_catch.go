package fishes

import "reflect"

type statementTryCatch struct {
	tryBody    expressions
	varType    *value
	varMessage *value
	catchBody  expressions
}

func (s *statementTryCatch) evaluate() (Value, bool, error) {
	valTry, bypassTry, errTry := s.tryBody.execute()
	if errTry != nil {
		return s.evaluateCatch(errTry)
	}
	if bypassTry {
		return valTry, true, nil
	}
	return predefinedNull, false, nil
}

func (s *statementTryCatch) evaluateCatch(err error) (Value, bool, error) {
	var errorType string
	if runtime, ok := err.(RuntimeError); ok {
		errorType = runtime.errorType
	} else {
		errorType = reflect.TypeOf(err).Name()
		if errorType == "" {
			errorType = "error"
		}
	}
	s.varType.assign(fromString(errorType))
	s.varMessage.assign(fromString(err.Error()))
	val, bypass, catchErr := s.catchBody.execute()
	if catchErr != nil {
		if runtime, ok := catchErr.(RuntimeError); ok {
			runtime.join(err)
			return nil, false, runtime
		}
		return nil, false, catchErr
	}
	if bypass {
		return val, true, nil
	}
	return predefinedNull, false, nil
}
