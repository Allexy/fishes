package fishes

import (
	"fmt"
	"github.com/Allexy/fishes/internal/lang"
	"github.com/Allexy/fishes/internal/tokenizer"
)

type scope struct {
	constants   map[string]Value
	variables   map[string]Value
	expressions expressions
}

func newScope() *scope {
	return &scope{
		constants:   make(map[string]Value),
		variables:   make(map[string]Value),
		expressions: newExpressions(),
	}
}

func spawnScope(parent *scope) *scope {
	child := &scope{
		constants:   make(map[string]Value),
		variables:   make(map[string]Value),
		expressions: newExpressions(),
	}
	for k, v := range parent.constants {
		child.constants[k] = v
	}
	for k, v := range parent.variables {
		child.variables[k] = v
	}
	return child
}

func (s *scope) parse(walker tokenizer.TokenWalker) error {
	// walker.Move(1) is for step over EOF token
	walker.Move(1)
	if err := s.parseBlock(walker); err != nil {
		return err
	}
	// Root block always ends with EOF
	if t := walker.Get(0); t.Type != tokenizer.TokenEOF {
		return newParseError(t.SourceName, fmt.Sprintf("unrecognized expression starting from %q", t.Type), t.Line, t.Col, nil)
	}
	// todo: should return instance of Expression
	return nil
}

func (s *scope) parseBlock(walker tokenizer.TokenWalker) error {
	for walker.Next() {
		switch {
		// check exit condition
		case walker.OneOf(tokenizer.TokenEOF, tokenizer.TokenCloseBrace):
			return nil
		// entrance in nested block "{..."
		case walker.Match(tokenizer.TokenOpenBrace):
			// step over "{"
			walker.Move(1)
			nested := spawnScope(s)
			if err := nested.parseBlock(walker); err != nil {
				return err
			}
			// does a nested block end with closing brace "...}"?
			if err := s.expect(walker, tokenizer.TokenCloseBrace); err != nil {
				return err
			}
			s.expressions.merge(nested.expressions)
		// statement if(...){... or switch(...){... or while(...){... or for(...){...
		case walker.Match(tokenizer.TokenWord, tokenizer.TokenOpenParen, tokenizer.TokenAny, tokenizer.TokenCloseParen, tokenizer.TokenOpenBrace):
			if err := s.parseStatement(walker); err != nil {
				return err
			}
			// does statement end with closing brace "...}"?
			if err := s.expect(walker, tokenizer.TokenCloseBrace); err != nil {
				return err
			}
		// statement try {...
		case walker.Match(tokenizer.TokenWord, tokenizer.TokenOpenBrace):
			if err := s.parseStatement(walker); err != nil {
				return err
			}
			// does statement end with closing brace "...}"?
			if err := s.expect(walker, tokenizer.TokenCloseBrace); err != nil {
				return err
			}
		// statement throw(...);
		case walker.Match(tokenizer.TokenWord, tokenizer.TokenOpenParen, tokenizer.TokenAny, tokenizer.TokenCloseParen) && walker.Get(0).Text == lang.KwThrow:
			if err := s.parseThrow(walker); err != nil {
				return err
			}
			// does statement end with semicolon ";"?
			if err := s.expect(walker, tokenizer.TokenSemicolon); err != nil {
				return err
			}
		// named constant definition @const_name = 123 or @const_name = "abc" or @const_name = True|False
		//  after tokenizer.TokenAssignment expects tokenizer.TokenNumber or tokenizer.TokenString or tokenizer.TokenLogic
		case walker.Match(tokenizer.TokenAt, tokenizer.TokenWord, tokenizer.TokenAssignment):
			if err := s.parseNamedConstant(walker.Get(1), walker.Get(3)); err != nil {
				return err
			}
			walker.Move(4)
			// does constant definition end with semicolon ";"?
			if err := s.expect(walker, tokenizer.TokenSemicolon); err != nil {
				return err
			}
		// statement return (aliased with "=")
		case (walker.OneOf(tokenizer.TokenWord) && walker.Get(0).Text == lang.KwReturn) || walker.OneOf(tokenizer.TokenAssignment):
			if err := s.parseReturn(walker); err != nil {
				return err
			}
			// does statement end with semicolon ";"?
			if err := s.expect(walker, tokenizer.TokenSemicolon); err != nil {
				return err
			}
		default:
			expr, err := s.parseExpression(walker)
			if err != nil {
				return err
			}
			// does expression end with semicolon ";"?
			if err := s.expect(walker, tokenizer.TokenSemicolon); err != nil {
				return err
			}
			s.expressions.add(expr)
		}
	}
	return nil
}

func (s *scope) parseStatement(walker tokenizer.TokenWalker) error {
	keyword := walker.Get(0)
	switch keyword.Text {
	case lang.KwIf:
		return s.parseStatementIf(walker)
	case lang.KwSwitch:
		return s.parseStatementSwitch(walker)
	case lang.KwWhile:
		return s.parseStatementWhile(walker)
	case lang.KwFor:
		return s.parseStatementFor(walker)
	case lang.KwTry:
		return s.parseStatementTry(walker)
	}
	return newParseError(
		keyword.SourceName,
		fmt.Sprintf("invalid statement: %q is unknown", keyword.Text),
		keyword.Line,
		keyword.Col,
		nil,
	)
}

func (s *scope) parseStatementIf(walker tokenizer.TokenWalker) error {
	// step over "if("
	walker.Move(2)
	condition, err := s.parseExpression(walker)
	if err != nil {
		return err
	}
	if err := s.expect(walker, tokenizer.TokenCloseParen); err != nil {
		return err
	}
	// step over "{"
	walker.Move(1)
	ifBody := spawnScope(s)
	if err := ifBody.parseBlock(walker); err != nil {
		return err
	}
	if err := s.expect(walker, tokenizer.TokenCloseBrace); err != nil {
		return err
	}
	s.expressions.add(&statementIf{condition, ifBody.expressions})
	return nil
}

func (s *scope) parseStatementSwitch(walker tokenizer.TokenWalker) error {
	// step over "switch("
	walker.Move(2)
	sc, err := s.parseExpression(walker)
	if err != nil {
		return err
	}
	if err := s.expect(walker, tokenizer.TokenCloseParen); err != nil {
		return err
	}
	stmt := &statementSwitch{sc, make([]statementCase, 0, 25)}
	// step over "{"
	walker.Move(1)
	for !walker.OneOf(tokenizer.TokenEOF, tokenizer.TokenCloseBrace) {
		// case(...){...
		if walker.Match(tokenizer.TokenWord, tokenizer.TokenOpenParen, tokenizer.TokenAny, tokenizer.TokenCloseParen, tokenizer.TokenOpenBrace) && walker.Get(0).Text == lang.KwCase {
			// step over "case("
			walker.Move(2)
			cc, err := s.parseExpression(walker)
			if err != nil {
				return err
			}
			if err := s.expect(walker, tokenizer.TokenCloseParen); err != nil {
				return err
			}
			// step over "{"
			walker.Move(1)
			caseBody := spawnScope(s)
			if err := caseBody.parseBlock(walker); err != nil {
				return err
			}
			if err := s.expect(walker, tokenizer.TokenCloseBrace); err != nil {
				return err
			}
			stmt.addCase(statementCase{cc, caseBody.expressions})
			continue
		}
		token := walker.Get(0)
		return newParseError(
			token.SourceName,
			"invalid statement: expected \"case(...){...\"",
			token.Line,
			token.Col,
			nil,
		)
	}
	s.expressions.add(stmt)
	return nil
}

func (s *scope) parseStatementWhile(walker tokenizer.TokenWalker) error {
	// step over "while("
	walker.Move(2)
	condition, err := s.parseExpression(walker)
	if err != nil {
		return err
	}
	if err := s.expect(walker, tokenizer.TokenCloseParen); err != nil {
		return err
	}
	// step over "{"
	walker.Move(1)
	whileBody := spawnScope(s)
	if err := whileBody.parseBlock(walker); err != nil {
		return err
	}
	if err := s.expect(walker, tokenizer.TokenCloseBrace); err != nil {
		return err
	}
	s.expressions.add(&statementWhile{condition, whileBody.expressions})
	return nil
}

func (s *scope) parseStatementFor(walker tokenizer.TokenWalker) error {
	// step over "for("
	walker.Move(2)
	initialization, err := s.parseExpression(walker)
	if err != nil {
		return err
	}
	if err := s.expect(walker, tokenizer.TokenSemicolon); err != nil {
		return err
	}
	condition, err := s.parseExpression(walker)
	if err != nil {
		return err
	}
	if err := s.expect(walker, tokenizer.TokenSemicolon); err != nil {
		return err
	}
	iteration, err := s.parseExpression(walker)
	if err != nil {
		return err
	}
	if err := s.expect(walker, tokenizer.TokenCloseParen); err != nil {
		return err
	}
	// step over "{"
	walker.Move(1)
	forBody := spawnScope(s)
	if err := forBody.parseBlock(walker); err != nil {
		return err
	}
	if err := s.expect(walker, tokenizer.TokenCloseBrace); err != nil {
		return err
	}
	s.expressions.add(&statementFor{initialization, condition, iteration, forBody.expressions})
	return nil
}

func (s *scope) parseStatementTry(walker tokenizer.TokenWalker) error {
	// step over try
	walker.Move(1)
	tryBody := spawnScope(s)
	if err := tryBody.parseBlock(walker); err != nil {
		return err
	}
	// does statement end with closing brace "...}"?
	if err := s.expect(walker, tokenizer.TokenCloseBrace); err != nil {
		return err
	}
	// statement catch($kind, $message) {...
	if !walker.Match(tokenizer.TokenWord, tokenizer.TokenOpenParen, tokenizer.TokenVariable, tokenizer.TokenComa, tokenizer.TokenVariable, tokenizer.TokenCloseParen, tokenizer.TokenOpenBrace) || walker.Get(0).Text != lang.KwCatch {
		token := walker.Get(0)
		return newParseError(
			token.SourceName,
			"invalid statement: expected \"catch($kind, $message) {\"",
			token.Line,
			token.Col,
			nil,
		)
	}
	varType, varMessage := variable(), variable()
	catchBody := spawnScope(s)
	tkVarType := walker.Get(2)
	tkVarMessage := walker.Get(4)
	// step over "catch($kind, $message) {"
	walker.Move(7)
	catchBody.variables[tkVarType.Text] = varType
	catchBody.variables[tkVarMessage.Text] = varMessage
	if err := catchBody.parseBlock(walker); err != nil {
		return err
	}
	s.expressions.add(&statementTryCatch{tryBody.expressions, varType, varMessage, catchBody.expressions})
	return nil
}

func (s *scope) parseThrow(walker tokenizer.TokenWalker) error {
	// step over "throw("
	tkn := walker.Get(0)
	walker.Move(2)
	kind, err := s.parseExpression(walker)
	if err != nil {
		return err
	}
	if err := s.expect(walker, tokenizer.TokenComa); err != nil {
		return err
	}
	message, err := s.parseExpression(walker)
	if err != nil {
		return err
	}
	s.expressions.add(&statementThrow{tkn.SourceName, tkn.Line, tkn.Col, kind, message})
	return nil
}

func (s *scope) parseNamedConstant(tWord *tokenizer.Token, tVal *tokenizer.Token) error {
	if _, exists := s.constants[tWord.Text]; exists {
		return newParseError(tWord.SourceName, fmt.Sprintf("constant %q already defined", tWord.Text), tWord.Line, tVal.Col, nil)
	}
	switch tVal.Type {
	case tokenizer.TokenNumber:
		s.constants[tWord.Text] = fromNumber(tVal.ParsedNumber)
	case tokenizer.TokenString:
		s.constants[tWord.Text] = fromString(tVal.Text)
	case tokenizer.TokenLogic:
		switch tVal.Text {
		case lang.KwTrue:
			s.constants[tWord.Text] = predefinedTrue
		case lang.KwFalse:
			s.constants[tWord.Text] = predefinedFalse
		default:
			panic(newParseError(tVal.SourceName, fmt.Sprintf("unexpected token text %q", tVal.Text), tVal.Line, tVal.Col, nil))
		}
	default:
		panic(newParseError(tVal.SourceName, fmt.Sprintf("unexpected token type %q", tVal.Type), tVal.Line, tVal.Col, nil))
	}
	return nil
}

func (s *scope) parseReturn(walker tokenizer.TokenWalker) error {
	// step over "return" or "="
	walker.Move(1)
	expr, err := s.parseExpression(walker)
	if err != nil {
		return err
	}
	s.expressions.add(&statementReturn{expr})
	return nil
}

func (s *scope) parseExpression(walker tokenizer.TokenWalker) (expression, error) {
	panic("parseExpression() is not implemented")
}

func (s *scope) expect(walker tokenizer.TokenWalker, tokenType tokenizer.TokenType) error {
	token := walker.Get(0)
	if token.Type != tokenType {
		return newParseError(
			token.SourceName,
			fmt.Sprintf("expected token %q got %q", tokenType, token.Type),
			token.Line,
			token.Col,
			nil,
		)
	}
	walker.Move(1)
	return nil
}
