package fishes

import "io"

type ValueType uint8
type MapKey int32

const (
	ValueString ValueType = iota
	ValueNumber
	ValueBoolean
	ValueLambda
	ValueList
	ValueMap
	ValueEmpty
)

type Value interface {
	AsString() string
	AsNumber() float64
	AsBoolean() bool
	AsCallable() Callable
	AsList() []Value
	AsMap() map[MapKey]Value
	MapKey() MapKey
	Type() ValueType
	Equals(v Value) bool
}

type Scenario interface {
	Run() (Value, error)
}

type Callable interface {
	Call(args ...Value) (Value, error)
}

func (vt ValueType) String() string {
	switch vt {
	case ValueString:
		return "ValueString"
	case ValueNumber:
		return "ValueNumber"
	case ValueBoolean:
		return "ValueBoolean"
	case ValueLambda:
		return "ValueLambda"
	case ValueList:
		return "ValueList"
	case ValueMap:
		return "ValueMap"
	case ValueEmpty:
		return "ValueEmpty"
	}
	panic("unknown ValueType")
}

func Compile(r io.Reader) (Scenario, error) {
	panic("Compile is not implemented")
}
