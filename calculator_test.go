package main

import (
	"fmt"
	"testing"
)

func TestIntDeclare(t *testing.T) {
	script := "int a;"
	lex := &SimpleLexer{}
	cal := &Calculator{}
	reader := lex.tokenize(script)
	node, err := cal.intDeclare(reader)
	if err != nil {
		t.Fatal(err)
	}
	DumpAST(node, "")
	fmt.Println("-----------------------")

	script = "int b = 2;"
	lex = &SimpleLexer{}
	cal = &Calculator{}
	node, err = cal.intDeclare(lex.tokenize(script))
	if err != nil {
		t.Fatal(err)
	}
	DumpAST(node, "")
	fmt.Println("-----------------------")

	script = "int c = 2+3;"
	lex = &SimpleLexer{}
	cal = &Calculator{}
	node, err = cal.intDeclare(lex.tokenize(script))
	if err != nil {
		t.Fatal(err)
	}
	DumpAST(node, "")
	fmt.Println("-----------------------")

	script = "int a = 3+b;"
	lex = &SimpleLexer{}
	cal = &Calculator{}
	node, err = cal.intDeclare(lex.tokenize(script))
	if err != nil {
		t.Fatal(err)
	}
	DumpAST(node, "")
	fmt.Println("-----------------------")

	script = "int a = 3*a-2+b;"
	lex = &SimpleLexer{}
	cal = &Calculator{}
	node, err = cal.intDeclare(lex.tokenize(script))
	if err != nil {
		t.Fatal(err)
	}
	DumpAST(node, "")
	fmt.Println("-----------------------")
}

func TestCalculator(t *testing.T) {
	// testing basic expression
	cal := &Calculator{}
	script := "2+3*5"
	fmt.Println("\ncalculating: " + script)
	if err := cal.execute(script); err != nil {
		t.Fatal(err)
	}
	fmt.Println("-------------------------")

	script = "1+2+3*5"
	fmt.Println("\ncalculating: " + script)
	if err := cal.execute(script); err != nil {
		t.Fatal(err)
	}
	fmt.Println("-------------------------")

	script = "1*2+3*5"
	fmt.Println("\ncalculating: " + script)
	if err := cal.execute(script); err != nil {
		t.Fatal(err)
	}
	fmt.Println("-------------------------")

	script = "2*(3+(2*2))"
	fmt.Println("\ncalculating: " + script)
	if err := cal.execute(script); err != nil {
		t.Fatal(err)
	}
	fmt.Println("-------------------------")

	script = "1*2+3*2+"
	fmt.Println("\ncalculating: " + script)
	if err := cal.execute(script); err != nil {
		t.Fatal(err)
	}
	fmt.Println("-------------------------")
}
