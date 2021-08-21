package ast

import (
	"github.com/stretchr/testify/require"
	"strings"
	"testing"

	"github.com/Allexy/fishes/internal/tokenizer"
)

func walker(sourceCode string) (tokenizer.TokenWalker, error) {
	return tokenizer.NewTokenizer(strings.NewReader(sourceCode), "string").Tokenize()
}

func TestParseConstNumber(t *testing.T) {
	w, err := walker("@CONST=1234;")
	require.NoError(t, err, "Must not be tokenization error")
	scope := NewScope()
	require.NoError(t, scope.Parse(w), "Must not be parsing errors")
	require.Containsf(t, scope.constants, "CONST", "Must contains constant %q", "CONST")
	val := scope.constants["CONST"].AsNumber()
	require.Equalf(t, 1234.0, val, "Expected value is 1234.0 but got %.1f", val)
}

func TestParseConstString(t *testing.T) {
	w, err := walker("@CONST=\"abc\";")
	require.NoError(t, err, "Must not be tokenization error")
	scope := NewScope()
	require.NoError(t, scope.Parse(w), "Must not be parsing errors")
	require.Containsf(t, scope.constants, "CONST", "Must contains constant %q", "CONST")
	val := scope.constants["CONST"].AsString()
	require.Equalf(t, "abc", val, "Expected value is %q but got %q", "abc", val)
}

func TestParseConstTrue(t *testing.T) {
	w, err := walker("@CONST=true;")
	require.NoError(t, err, "Must not be tokenization error")
	scope := NewScope()
	require.NoError(t, scope.Parse(w), "Must not be parsing errors")
	require.Containsf(t, scope.constants, "CONST", "Must contains constant %q", "CONST")
	val := scope.constants["CONST"].AsBoolean()
	require.Equalf(t, true, val, "Expected value is %q but got %q", true, val)
}

func TestParseConstFalse(t *testing.T) {
	w, err := walker("@CONST=false;")
	require.NoError(t, err, "Must not be tokenization error")
	scope := NewScope()
	require.NoError(t, scope.Parse(w), "Must not be parsing errors")
	require.Containsf(t, scope.constants, "CONST", "Must contains constant %q", "CONST")
	val := scope.constants["CONST"].AsBoolean()
	require.Equalf(t, false, val, "Expected value is %q but got %q", false, val)
}
