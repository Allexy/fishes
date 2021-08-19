package ast

import (
	"math"
	"strings"
	"testing"

	"github.com/Allexy/fishes/internal/tokenizer"
)

func walker(sourceCode string) (tokenizer.TokenWalker, error) {
	return tokenizer.NewTokenizer(strings.NewReader(sourceCode), "string").Tokenize()
}

func TestParseConstNumber(t *testing.T) {
	w, err := walker("@CONST=1234;")
	if err != nil {
		t.Error(err)
		return
	}
	scope := NewScope()
	err = scope.Parse(w)
	if err != nil {
		t.Error(err)
		return
	}
	c, exists := scope.constants["CONST"]
	if !exists {
		t.Errorf("Expected constant %q in scopes constants", "CONST")
		return
	}
	if math.Abs(c.AsNumber()-1234.0) > 0.0 {
		t.Errorf("Expected value is 1234.0 but got %.1f", c.AsNumber())
	}
}

func TestParseConstString(t *testing.T) {
	w, err := walker("@CONST=\"abc\";")
	if err != nil {
		t.Error(err)
		return
	}
	scope := NewScope()
	err = scope.Parse(w)
	if err != nil {
		t.Error(err)
		return
	}
	c, exists := scope.constants["CONST"]
	if !exists {
		t.Errorf("Expected constant %q in scopes constants", "CONST")
		return
	}
	if c.AsString() != "abc" {
		t.Errorf("Expected value is %q but got %q", "abc", c.AsString())
	}
}

func TestParseConstTrue(t *testing.T) {
	w, err := walker("@CONST=true;")
	if err != nil {
		t.Error(err)
		return
	}
	scope := NewScope()
	err = scope.Parse(w)
	if err != nil {
		t.Error(err)
		return
	}
	c, exists := scope.constants["CONST"]
	if !exists {
		t.Errorf("Expected constant %q in scopes constants", "CONST")
		return
	}
	if c.AsBoolean() != true {
		t.Errorf("Expected value is %v but got %v", true, c.AsBoolean())
	}
}

func TestParseConstFalse(t *testing.T) {
	w, err := walker("@CONST=false;")
	if err != nil {
		t.Error(err)
		return
	}
	scope := NewScope()
	err = scope.Parse(w)
	if err != nil {
		t.Error(err)
		return
	}
	c, exists := scope.constants["CONST"]
	if !exists {
		t.Errorf("Expected constant %q in scopes constants", "CONST")
		return
	}
	if c.AsBoolean() != false {
		t.Errorf("Expected value is %v but got %v", false, c.AsBoolean())
	}
}
