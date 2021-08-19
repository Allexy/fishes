package ast

import (
	"fmt"

	"github.com/Allexy/fishes/internal/lang"
	"github.com/Allexy/fishes/internal/tokenizer"
	"github.com/Allexy/fishes/pkg/fishes"
)

var predefinedTrue fishes.Value = newConstantValue(fishes.ValueBoolean, lang.KwTrue)
var predefinedFalse fishes.Value = newConstantValue(fishes.ValueBoolean, lang.KwFalse)
var predefinedNull fishes.Value = newConstantValue(fishes.ValueEmpty, "")

type Scope struct {
	constants map[string]fishes.Value
	variables map[string]fishes.Value
}

func NewScope() *Scope {
	return &Scope{
		constants: make(map[string]fishes.Value),
		variables: make(map[string]fishes.Value),
	}
}

func spawnScope(parent *Scope) *Scope {
	child := &Scope{
		constants: parent.constants,
		variables: make(map[string]fishes.Value),
	}
	for k, v := range parent.variables {
		child.variables[k] = v
	}
	return child
}

func (s *Scope) Parse(walker tokenizer.TokenWalker) error {
	// walker.Move(1) is for step over BOF and EOF tokens
	for walker.Move(1); walker.Next(); walker.Move(1) {
		// named constant definition @const_name = 123;
		if walker.Match(tokenizer.TokenAt, tokenizer.TokenWord, tokenizer.TokenAssignment, tokenizer.TokenNumber, tokenizer.TokenSemicolon) {
			s.parseNamedConstant(walker.Get(1), walker.Get(3))
			walker.Move(5)
			continue
		}
		// named constant definition @const_name = "abc";
		if walker.Match(tokenizer.TokenAt, tokenizer.TokenWord, tokenizer.TokenAssignment, tokenizer.TokenString, tokenizer.TokenSemicolon) {
			s.parseNamedConstant(walker.Get(1), walker.Get(3))
			walker.Move(5)
			continue
		}
		// named constant definition @const_name = True|False;
		if walker.Match(tokenizer.TokenAt, tokenizer.TokenWord, tokenizer.TokenAssignment, tokenizer.TokenLogic, tokenizer.TokenSemicolon) {
			s.parseNamedConstant(walker.Get(1), walker.Get(3))
			walker.Move(5)
			continue
		}
		// final test
		if t := walker.Get(0); t.Token != tokenizer.TokenEOF {
			return newParserError(t.SourceName, fmt.Sprintf("Unrecognized expression starting from %q", t.Token), t.Line, t.Col, nil)
		}
	}
	return nil
}

func (s *Scope) parseNamedConstant(tWord *tokenizer.Token, tVal *tokenizer.Token) error {
	if _, exists := s.constants[tWord.Text]; exists {
		return newParserError(tWord.SourceName, fmt.Sprintf("Constant %q already defined", tWord.Text), tWord.Line, tVal.Col, nil)
	}
	switch tVal.Token {
	case tokenizer.TokenNumber:
		s.constants[tWord.Text] = newConstantValue(fishes.ValueNumber, tVal.Text)
	case tokenizer.TokenString:
		s.constants[tWord.Text] = newConstantValue(fishes.ValueString, tVal.Text)
	case tokenizer.TokenLogic:
		switch tVal.Text {
		case lang.KwTrue:
			s.constants[tWord.Text] = predefinedTrue
		case lang.KwFalse:
			s.constants[tWord.Text] = predefinedFalse
		default:
			panic(newParserError(tVal.SourceName, fmt.Sprintf("Unexpected token text %q", tVal.Text), tVal.Line, tVal.Col, nil))
		}
	default:
		panic(newParserError(tVal.SourceName, fmt.Sprintf("Unexpected token type %q", tVal.Token), tVal.Line, tVal.Col, nil))
	}
	return nil
}
