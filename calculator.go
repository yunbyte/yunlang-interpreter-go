package main

import (
	"fmt"
	"strconv"
)

type ASTNodeType int

const (
	ASTNodeType_Program ASTNodeType = iota //program entry, the root node

	ASTNodeType_IntDeclaration // integer variable declaration
	ASTNodeType_Expression     // expression statement, followed by a semicolon
	ASTNodeType_Assignment     // assignment statement

	ASTNodeType_Primary        // basic expression
	ASTNodeType_Multiplicative // multiplication expression
	ASTNodeType_Additive       // addition expression

	ASTNodeType_Identifier // identifier
	ASTNodeType_IntLiteral // integer literal
)

var astNodeTypeStringMap = map[ASTNodeType]string{
	ASTNodeType_Program:        "Program",
	ASTNodeType_IntDeclaration: "IntDeclaration",
	ASTNodeType_Expression:     "Expression",
	ASTNodeType_Assignment:     "Assignment",
	ASTNodeType_Primary:        "Primary",
	ASTNodeType_Multiplicative: "Multiplicative",
	ASTNodeType_Additive:       "Additive",
	ASTNodeType_Identifier:     "Identifier",
	ASTNodeType_IntLiteral:     "IntLiteral",
}

func (t ASTNodeType) String() string {
	if s, ok := astNodeTypeStringMap[t]; ok {
		return s
	}
	return "Unknown ASTNodeType"
}

type ASTNode interface {
	GetParent() ASTNode
	GetChildren() []ASTNode
	GetType() ASTNodeType
	GetText() string
	AddChild(ASTNode)
}

type SimpleASTNode struct {
	parent   ASTNode
	nodeType ASTNodeType
	children []ASTNode
	text     string
}

func (n *SimpleASTNode) GetChildren() []ASTNode {
	return n.children
}

func (n *SimpleASTNode) GetParent() ASTNode {
	return n.parent
}

func (n *SimpleASTNode) GetType() ASTNodeType {
	return n.nodeType
}

func (n *SimpleASTNode) GetText() string {
	return n.text
}

func (n *SimpleASTNode) AddChild(child ASTNode) {
	n.children = append(n.children, child)
}

func NewASTNode(nodeType ASTNodeType, text string) ASTNode {
	return &SimpleASTNode{nodeType: nodeType, text: text}
}

// DumpAST Print tree struct for ASTNode
func DumpAST(node ASTNode, indent string) {
	fmt.Println(indent, node.GetType(), " ", node.GetText())
	for _, child := range node.GetChildren() {
		DumpAST(child, indent+"\t")
	}
}

type Calculator struct {
}

// execute execute script and print all AST node and calculating process.
func (c *Calculator) execute(script string) error {
	tree, err := c.parse(script)
	if err != nil {
		return err
	}
	DumpAST(tree, "")
	c.evaluate(tree, "")

	return nil
}

// parse parse script and return root AST node.
func (c *Calculator) parse(code string) (ASTNode, error) {
	lexer := SimpleLexer{}
	tokens := lexer.tokenize(code)
	return c.prog(tokens)
}

func (c *Calculator) evaluate(node ASTNode, indent string) int {
	var result int = 0
	fmt.Printf("%sCalculating: %s\n", indent, node.GetType())

	switch node.GetType() {
	case ASTNodeType_Program:
		for _, child := range node.GetChildren() {
			result = c.evaluate(child, indent+"\t")
		}
	case ASTNodeType_Additive:
		child1 := node.GetChildren()[0]
		child2 := node.GetChildren()[1]
		value1 := c.evaluate(child1, indent+"\t")
		value2 := c.evaluate(child2, indent+"\t")
		if node.GetText() == "+" {
			result = value1 + value2
		} else {
			result = value1 - value2
		}
	case ASTNodeType_Multiplicative:
		child1 := node.GetChildren()[0]
		child2 := node.GetChildren()[1]
		value1 := c.evaluate(child1, indent+"\t")
		value2 := c.evaluate(child2, indent+"\t")
		if node.GetText() == "*" {
			result = value1 * value2
		} else {
			result = value1 / value2
		}
	case ASTNodeType_IntLiteral:
		result, _ = strconv.Atoi(node.GetText())
	}

	fmt.Printf("%sResult: %d\n", indent, result)
	return result
}

// prog syntax analysis: root node
func (c *Calculator) prog(tokens *TokenReader) (ASTNode, error) {
	node := NewASTNode(ASTNodeType_Program, "Calculator")
	child, err := c.additive(tokens)
	if err != nil {
		return nil, err
	}
	if node != nil {
		node.AddChild(child)
	}

	return node, nil
}

// intDeclare syntax analysis: integer variable declaration statement
//
// intDeclaration : Int Identifier ('=' additiveExpression)?;
// int a;
// int b = 2;
// int c = 2+3;
func (c *Calculator) intDeclare(tokens *TokenReader) (ASTNode, error) {
	var node ASTNode
	token := tokens.Peek()                           // pre-read
	if token != nil && token.Type == TokenType_INT { // match int
		tokens.Read() // remove int
		token = tokens.Peek()
		if token != nil && token.Type == TokenType_Identifier { // match Identifier
			token = tokens.Read() // remove Identifier
			// create a ASTNode, type is IntDeclaration and text is Identifier value
			node = NewASTNode(ASTNodeType_IntDeclaration, token.Text)

			token = tokens.Peek()
			if token != nil && token.Type == TokenType_Assign { // match =
				tokens.Read()
				child, err := c.additive(tokens)
				if err != nil {
					return nil, err
				}
				if child != nil {
					node.AddChild(child)
				} else {
					return nil, fmt.Errorf("syntax analysis intDeclare error: invalid variable initialization, expecting an expression.")
				}
			}
		} else {
			return nil, fmt.Errorf("syntax analysis intDeclare error: variable name expected.")
		}

		if node != nil {
			token = tokens.Peek()
			if token != nil && token.Type == TokenType_SemiColon {
				tokens.Read()
			} else {
				return nil, fmt.Errorf("syntax analysis intDeclare error: invalid statement, expecting semicolon.")
			}
		}
	}

	return node, nil
}

// additive syntax analysis: addition expression
//
// additiveExpression
// : multiplicativeExpression
// | multiplicativeExpression Plus additiveExpression
// ;
func (c *Calculator) additive(tokens *TokenReader) (ASTNode, error) {
	var node ASTNode
	child1, err := c.multiplicative(tokens) // parse the first node
	if err != nil {
		return nil, err
	}
	node = child1 // return the first node if no second node
	token := tokens.Peek()
	if child1 != nil && token != nil {
		if token.Type == TokenType_Plus || token.Type == TokenType_Minus { // match + or -
			token = tokens.Read()
			child2, err := c.additive(tokens)
			if err != nil {
				return nil, err
			}
			if child2 != nil {
				node = NewASTNode(ASTNodeType_Additive, token.Text)
				node.AddChild(child1)
				node.AddChild(child2)
			} else {
				return nil, fmt.Errorf("syntax analysis additive error: invalid additive expression, expecting the right part.")
			}
		}
	}

	return node, nil
}

// multiplicative  Syntax analysis: addition expression
//
//	multiplicativeExpression
//	: IntLiteral
//	| IntLiteral Star multiplicativeExpression
//	;
func (c *Calculator) multiplicative(tokens *TokenReader) (ASTNode, error) {
	var node ASTNode
	child1, err := c.primary(tokens) // parse the first node
	if err != nil {
		return nil, err
	}
	node = child1
	token := tokens.Peek()
	if child1 != nil && token != nil {
		if token.Type == TokenType_Star || token.Type == TokenType_Slash {
			token = tokens.Read()
			child2, err := c.multiplicative(tokens)
			if err != nil {
				return nil, err
			}
			if child2 != nil {
				node = NewASTNode(ASTNodeType_Multiplicative, token.Text)
				node.AddChild(child1)
				node.AddChild(child2)
			} else {
				return nil, fmt.Errorf("multiplicative syntax analysis error: invalid multiplicative expression, expecting the right part.")
			}
		}
	}

	return node, nil
}

// primary syntax analysis: basic expression
// IntLiteral
func (c *Calculator) primary(tokens *TokenReader) (ASTNode, error) {
	var node ASTNode
	var err error
	token := tokens.Peek()
	if token != nil {
		if token.Type == TokenType_IntLiteral {
			token = tokens.Read()
			node = NewASTNode(ASTNodeType_IntLiteral, token.Text)
		} else if token.Type == TokenType_Identifier {
			token = tokens.Read()
			node = NewASTNode(ASTNodeType_Identifier, token.Text)
		} else if token.Type == TokenType_LeftParen {
			token = tokens.Read()
			node, err = c.additive(tokens)
			if err != nil {
				return nil, err
			}
			if node != nil {
				token = tokens.Peek()
				if token != nil && token.Type == TokenType_RightParen {
					tokens.Read()
				} else {
					return nil, fmt.Errorf("primary syntax analysis error: expecting right parenthesis.")
				}

			}
		} else {
			return nil, fmt.Errorf("primary syntax analysis error: invalid multiplicative expression, expecting the right part.")
		}
	}

	return node, nil
}
