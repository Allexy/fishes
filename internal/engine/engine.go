package engine

import (
	"errors"
	"os"
	"strings"

	"github.com/Allexy/fishes/internal/ast"
	"github.com/Allexy/fishes/internal/tokenizer"
	"github.com/Allexy/fishes/pkg/fishes"
)

func ReadString(sourceCode string) (fishes.Expression, error) {
	walker, tokenizerErr := tokenizer.NewTokenizer(strings.NewReader(sourceCode), "string").Tokenize()
	if tokenizerErr != nil {
		return nil, tokenizerErr
	}
	scope := ast.NewScope()
	parserErr := scope.Parse(walker)
	if parserErr != nil {
		return nil, parserErr
	}
	return nil, errors.New("not implemented")
}

func ReadFile(fileName string) (fishes.Expression, error) {
	f, fileErr := os.Open(fileName)
	if fileErr != nil {
		return nil, fileErr
	}
	defer f.Close()
	walker, tokenizerErr := tokenizer.NewTokenizer(f, fileName).Tokenize()
	if tokenizerErr != nil {
		return nil, tokenizerErr
	}
	scope := ast.NewScope()
	parserErr := scope.Parse(walker)
	if parserErr != nil {
		return nil, parserErr
	}
	return nil, errors.New("not implemented")
}
