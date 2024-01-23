package main

import (
	"fmt"
	"testing"
)

func TestLex(t *testing.T) {
	lexer := &SimpleLexer{}
	script := "int age = 45;"
	reader := lexer.tokenize(script)
	lexer.dump(reader)
	fmt.Println("------------------")

	lexer = &SimpleLexer{}
	script = "age >= 410;"
	reader = lexer.tokenize(script)
	lexer.dump(reader)
	fmt.Println("------------------")

	lexer = &SimpleLexer{}
	script = "intA = 67;"
	reader = lexer.tokenize(script)
	lexer.dump(reader)
	fmt.Println("------------------")

	lexer = &SimpleLexer{}
	script = "1+2*3/6 == 21;"
	reader = lexer.tokenize(script)
	lexer.dump(reader)
	fmt.Println("------------------")

	lexer = &SimpleLexer{}
	script = "1 * 4 * (1+2);"
	reader = lexer.tokenize(script)
	lexer.dump(reader)
}
