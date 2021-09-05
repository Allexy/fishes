package fishes

import (
	"fmt"
	"github.com/Allexy/fishes/internal/lang"
	"math"
)

type value struct {
	vType     ValueType
	vString   string
	vNumber   float64
	vBoolean  bool
	vCallable Callable
}

var predefinedTrue = &value{ValueBoolean, lang.KwTrue, math.NaN(), true, nil}
var predefinedFalse = &value{ValueBoolean, lang.KwFalse, math.NaN(), false, nil}
var predefinedNull = &value{ValueEmpty, lang.KwNull, math.NaN(), false, nil}

func fromString(v string) *value {
	return &value{ValueString, v, math.NaN(), len(v) > 0, nil}
}

func fromNumber(v float64) *value {
	return &value{ValueNumber, "", v, !math.IsNaN(v) && v != 0.0, nil}
}

func fromBoolean(v bool) *value {
	if v {
		return predefinedTrue
	}
	return predefinedFalse
}

func variable() *value {
	return &value{ValueEmpty, "", math.NaN(), false, nil}
}

func (vv *value) assign(another Value) {

	/*
		Converting only boolean values in order to test string is empty
		$someVar = "string value";
		if(!$someVar) {
			//do stuff...
		}
	*/
	switch another.Type() {
	case ValueString:
		vv.vType = ValueString
		vv.vString = another.AsString()
		vv.vNumber = math.NaN()
		vv.vBoolean = len(vv.vString) > 0
		vv.vCallable = nil
	case ValueNumber:
		vv.vType = ValueNumber
		vv.vString = ""
		vv.vNumber = another.AsNumber()
		vv.vBoolean = !math.IsNaN(vv.vNumber)
		vv.vCallable = nil
	case ValueBoolean:
		vv.vType = ValueBoolean
		vv.vString = ""
		vv.vNumber = math.NaN()
		vv.vBoolean = another.AsBoolean()
		vv.vCallable = nil
	case ValueLambda:
		vv.vType = ValueLambda
		vv.vString = ""
		vv.vNumber = math.NaN()
		vv.vCallable = another.AsCallable()
		vv.vBoolean = vv.vCallable != nil
	case ValueEmpty:
		vv.vType = ValueEmpty
		vv.vString = ""
		vv.vNumber = math.NaN()
		vv.vCallable = nil
		vv.vBoolean = false
	}
}

func (vv *value) AsString() string {
	return vv.vString
}

func (vv *value) AsNumber() float64 {
	return vv.vNumber
}

func (vv *value) AsBoolean() bool {
	return vv.vBoolean
}

func (vv *value) AsCallable() Callable {
	return vv.vCallable
}

func (vv *value) AsList() []Value {
	panic("AsList() is not implemented")
}

func (vv *value) AsMap() map[MapKey]Value {
	panic("AsMap() is not implemented")
}

func (vv *value) MapKey() MapKey {
	panic("MapKey() is not implemented")
}

func (vv *value) Type() ValueType {
	return vv.vType
}

func (vv *value) Equals(another Value) bool {
	if vv.vType != another.Type() {
		return false
	}
	switch vv.vType {
	case ValueString:
		return vv.vString == another.AsString()
	case ValueNumber:
		return compareFloat64(vv.vNumber, another.AsNumber())
	case ValueBoolean:
		return vv.vBoolean == another.AsBoolean()
	case ValueEmpty:
		return true
	case ValueList:
		panic("lists comparison is not implemented")
	case ValueMap:
		panic("maps comparison is not implemented")
	case ValueLambda:
		return vv.vCallable == another.AsCallable() // is it legal?
	}
	panic("unknown value type")
}

func (vv *value) String() string {
	switch vv.vType {
	case ValueString:
		return fmt.Sprintf("%v(%q)", vv.vType, vv.vString)
	case ValueBoolean:
		return fmt.Sprintf("%v(%v)", vv.vType, vv.vBoolean)
	case ValueNumber:
		return fmt.Sprintf("%v(%v)", vv.vType, vv.vNumber)
	default:
		return vv.vType.String()
	}
}

func (vv *value) evaluate() (Value, bool, error) {
	return vv, false, nil
}

func compareFloat64(a, b float64) bool {
	if math.IsNaN(a) || math.IsNaN(b) {
		return false
	}
	if a == b {
		return true
	}
	const smallestNormal = 2.2250738585072014e-308 // 2**-1022
	const epsilon = 1e-4
	absA, absB, diff := math.Abs(a), math.Abs(b), math.Abs(a-b)
	if absA == 0.0 || absB == 0.0 || absA+absB < smallestNormal {
		return diff < (epsilon * smallestNormal)
	}
	return diff/math.Min(absA+absB, math.MaxFloat64) < epsilon
}
