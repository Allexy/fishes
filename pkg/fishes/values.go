package fishes

type ValueType uint8

const (
	ValueString ValueType = iota
	ValueNumber
	ValueBoolean
	ValueLambda
	ValueCollection
	ValueEmpty
)

type Value interface {
	Copy(another Value)
	AsString() string
	AsNumber() float32
	AsBoolean() bool
	AsCallable() Callable
	// Get collection accessor
	Get(key Value) Value
	// Set collection accessor
	Set(key Value, value Value)
	Type() ValueType
	// calculate values hash for comparison of values
	Hash() int32
}

type Expression interface {
	Evaluate() Value
}

type Callable interface {
	Call(args ...Value) Value
}
