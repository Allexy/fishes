package main

import (
	"fmt"
	"os"

	"github.com/Allexy/fishes/internal/tokenizer"
)

func main() {
	fmt.Println("Test example")
	sourceName := "self_test.fs"
	f, err := os.Open(sourceName)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	tzr := tokenizer.NewTokenizer(f, sourceName)
	if tokens, err := tzr.Tokenize(); err == nil {
		for tokens.Next() {
			t := tokens.Get(0)
			fmt.Println(t)
			tokens.Move(1)
		}
	} else {
		panic(err)
	}

	fmt.Println("DONE")
}
