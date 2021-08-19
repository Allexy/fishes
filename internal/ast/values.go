package ast

import (
	"fmt"
	"math"
	"strconv"

	"github.com/Allexy/fishes/internal/lang"
	"github.com/Allexy/fishes/pkg/fishes"
)

// Variable value

type VariableValue struct {
	variableType   fishes.ValueType
	stringValue    string
	numericalValue float64
	booleanValue   bool
	callableValue  fishes.Callable
}

func (vv *VariableValue) Copy(another fishes.Value) {
	/*
		Converting only boolean values in order to test string is empty
		$someVar = "string value";
		if(!$someVar) {
			//do stuff...
		}
	*/
	switch another.Type() {
	case fishes.ValueString:
		vv.variableType = fishes.ValueString
		vv.stringValue = another.AsString()
		vv.numericalValue = math.NaN()
		vv.booleanValue = len(vv.stringValue) > 0
		vv.callableValue = nil
	case fishes.ValueNumber:
		vv.variableType = fishes.ValueNumber
		vv.stringValue = ""
		vv.numericalValue = another.AsNumber()
		vv.booleanValue = !math.IsNaN(vv.numericalValue)
		vv.callableValue = nil
	case fishes.ValueBoolean:
		vv.variableType = fishes.ValueBoolean
		vv.stringValue = ""
		vv.numericalValue = math.NaN()
		vv.booleanValue = another.AsBoolean()
		vv.callableValue = nil
	case fishes.ValueLambda:
		vv.variableType = fishes.ValueLambda
		vv.stringValue = ""
		vv.numericalValue = math.NaN()
		vv.callableValue = another.AsCallable()
		vv.booleanValue = vv.callableValue != nil
	case fishes.ValueCollection:
		panic("No implementation yet")
	case fishes.ValueEmpty:
		vv.variableType = fishes.ValueEmpty
		vv.stringValue = ""
		vv.numericalValue = math.NaN()
		vv.callableValue = nil
		vv.booleanValue = false
	}
}

func (vv VariableValue) AsString() string {
	return vv.stringValue
}

func (vv VariableValue) AsNumber() float64 {
	return vv.numericalValue
}

func (vv VariableValue) AsBoolean() bool {
	return vv.booleanValue
}

func (vv VariableValue) AsCallable() fishes.Callable {
	return vv.callableValue
}

func (vv VariableValue) Get(k fishes.Value) fishes.Value {
	panic("No implementation yet")
}

func (vv VariableValue) Set(k fishes.Value, v fishes.Value) {
	panic("No implementation yet")
}

func (vv VariableValue) Type() fishes.ValueType {
	return vv.variableType
}

// Constant value (consts and literals)

type constantValue struct {
	valueType      fishes.ValueType
	stringValue    string
	numericalValue float64
	booleanValue   bool
}

func newConstantValue(vt fishes.ValueType, txt string) fishes.Value {
	switch vt {
	case fishes.ValueString:
		return &constantValue{valueType: vt, stringValue: txt}
	case fishes.ValueNumber:
		val, err := strconv.ParseFloat(txt, 64)
		if err != nil {
			// this should never happen
			// but better to check
			panic(err)
		}
		return &constantValue{valueType: vt, numericalValue: val}
	case fishes.ValueBoolean:
		switch txt {
		case lang.KwTrue:
			return &constantValue{valueType: vt, booleanValue: true}
		case lang.KwFalse:
			return &constantValue{valueType: vt, booleanValue: false}
		default:
			panic(fmt.Errorf("invalid value for boolean constant %q", txt))
		}
	case fishes.ValueEmpty:
		return &constantValue{valueType: vt}
	default:
		panic(fmt.Errorf("unexpected value type: %d", vt))
	}
}

func (cv constantValue) Copy(another fishes.Value) {
	panic("Can't use Copy acceptor on constant")
}

func (cv constantValue) AsString() string {
	return cv.stringValue
}

func (cv constantValue) AsNumber() float64 {
	return cv.numericalValue
}

func (cv constantValue) AsBoolean() bool {
	return cv.booleanValue
}

func (cv constantValue) AsCallable() fishes.Callable {
	return nil
}

func (cv constantValue) Get(k fishes.Value) fishes.Value {
	panic("No implementation yet")
}

func (cv constantValue) Set(k fishes.Value, v fishes.Value) {
	panic("No implementation yet")
}

func (cv constantValue) Type() fishes.ValueType {
	return cv.valueType
}
