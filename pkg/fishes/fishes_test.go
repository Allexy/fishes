package fishes

import (
	"github.com/Allexy/fishes/internal/tokenizer"
	"github.com/stretchr/testify/require"
	"math"
	"strings"
	"testing"
)

func walker(sourceCode string) (tokenizer.TokenWalker, error) {
	return tokenizer.NewTokenizer(strings.NewReader(sourceCode), "string").Tokenize()
}

func TestValueString(t *testing.T) {
	var val Value = fromString("string")
	require.Equalf(t, ValueString, val.Type(), "expected type %q; got %q", ValueString, val.Type())
	require.Equalf(t, "string", val.AsString(), "expected value %q; got %q", "string", val.AsString())
	require.True(t, val.AsBoolean(), "expected boolean representation is true")
	require.Truef(t, math.IsNaN(val.AsNumber()), "expected NaN got %q", val.AsNumber())
	require.Nil(t, val.AsCallable(), "no callable value expected")

	val = fromString("")
	require.False(t, val.AsBoolean(), "expected boolean representation is false")
}

func TestValueNumber(t *testing.T) {
	var val Value = fromNumber(10.0)
	require.Equalf(t, ValueNumber, val.Type(), "expected type %q; got %q", ValueNumber, val.Type())
	require.Equalf(t, 10.0, val.AsNumber(), "expected value %q; got %q", "string", val.AsNumber())
	require.True(t, val.AsBoolean(), "expected boolean representation is true")
	require.Equalf(t, "", val.AsString(), "expected empty string got %q", val.AsString())
	require.Nil(t, val.AsCallable(), "no callable value expected")
	val = fromNumber(0.0)
	require.False(t, val.AsBoolean(), "expected boolean representation is false")
	val = fromNumber(math.NaN())
	require.False(t, val.AsBoolean(), "expected boolean representation is false")
}

func TestValueBoolean(t *testing.T) {
	var val Value = fromBoolean(true)
	require.Equalf(t, ValueBoolean, val.Type(), "expected type %q; got %q", ValueBoolean, val.Type())
	require.Equalf(t, "true", val.AsString(), "expected value %q; got %q", "true", val.AsString())
	require.True(t, val.AsBoolean(), "expected boolean representation is true")
	require.True(t, math.IsNaN(val.AsNumber()), "expected numerical representation is NaN")
	require.Nil(t, val.AsCallable(), "no callable value expected")

	val = fromBoolean(false)
	require.Equalf(t, ValueBoolean, val.Type(), "expected type %q; got %q", ValueBoolean, val.Type())
	require.Equalf(t, "false", val.AsString(), "expected value %q; got %q", "false", val.AsString())
	require.False(t, val.AsBoolean(), "expected boolean representation is false")
	require.True(t, math.IsNaN(val.AsNumber()), "expected numerical representation is NaN")
	require.Nil(t, val.AsCallable(), "no callable value expected")
}

func TestParseConstNumber(t *testing.T) {
	w, err := walker("const CONST_NAME=1234;")
	require.NoError(t, err, "must not be tokenization error")
	scope := newScope()
	require.NoError(t, scope.parse(w), "must not be parsing errors")
	require.Containsf(t, scope.constants, "CONST_NAME", "must contains constant %q", "CONST_NAME")
	val := scope.constants["CONST_NAME"].AsNumber()
	require.Equalf(t, 1234.0, val, "expected value is 1234.0 but got %.1f", val)
}

func TestParseConstString(t *testing.T) {
	w, err := walker("const CONST_NAME=\"abc\";")
	require.NoError(t, err, "must not be tokenization error")
	scope := newScope()
	require.NoError(t, scope.parse(w), "must not be parsing errors")
	require.Containsf(t, scope.constants, "CONST_NAME", "must contains constant %q", "CONST_NAME")
	val := scope.constants["CONST_NAME"].AsString()
	require.Equalf(t, "abc", val, "expected value is %q but got %q", "abc", val)
}

func TestParseConstTrue(t *testing.T) {
	w, err := walker("const CONST_NAME=true;")
	require.NoError(t, err, "must not be tokenization error")
	scope := newScope()
	require.NoError(t, scope.parse(w), "must not be parsing errors")
	require.Containsf(t, scope.constants, "CONST_NAME", "must contains constant %q", "CONST_NAME")
	val := scope.constants["CONST_NAME"].AsBoolean()
	require.Equalf(t, true, val, "expected value is %q but got %q", true, val)
}

func TestParseConstFalse(t *testing.T) {
	w, err := walker("const CONST_NAME=false;")
	require.NoError(t, err, "must not be tokenization error")
	scope := newScope()
	require.NoError(t, scope.parse(w), "must not be parsing errors")
	require.Containsf(t, scope.constants, "CONST_NAME", "must contains constant %q", "CONST_NAME")
	val := scope.constants["CONST_NAME"].AsBoolean()
	require.Equalf(t, false, val, "expected value is %q but got %q", false, val)
}
