package main

import "fmt"

type DfaState int

const (
	DfaState_Initial DfaState = iota
	DfaState_Id
	DfaState_Int1
	DfaState_Int2
	DfaState_Int3

	DfaState_IntLiteral
	DfaState_SemiColon
	DfaState_LeftParen
	DfaState_RightParen
	DfaState_Assignment
	DfaState_GT
	DfaState_GE
	DfaState_LT
	DfaState_LE
	DfaState_Plus
	DfaState_Minus
	DfaState_Star
	DfaState_Slash
	DfaState_EQ
	DfaState_PlusEQ
	DfaState_MinusEQ
	DfaState_StarEQ
	DfaState_SlashEQ
)

type Token struct {
	Text string
	Type TokenType
}

type TokenType int

const (
	TokenType_INIT TokenType = iota
	TokenType_INT
	TokenType_DOUBLE
	TokenType_STRING
	TokenType_BOOL
	TokenType_NIL

	TokenType_GT
	TokenType_GE
	TokenType_LT
	TokenType_LE
	TokenType_EQ
	TokenType_Assign
	TokenType_Plus
	TokenType_PlusPlus
	TokenType_PlusEQ
	TokenType_Minus
	TokenType_MinusMinus
	TokenType_MinusEQ
	TokenType_Star
	TokenType_StarEQ
	TokenType_Slash
	TokenType_SlashEQ

	TokenType_LeftParen
	TokenType_RightParen
	TokenType_SemiColon

	TokenType_Identifier
	TokenType_IntLiteral
)

var tokenTypeStringMap = map[TokenType]string{
	TokenType_INT:    "INT",
	TokenType_DOUBLE: "DOUBLE",
	TokenType_STRING: "STRING",
	TokenType_BOOL:   "BOOL",
	TokenType_NIL:    "NIL",

	TokenType_GT:         "GT",         // >
	TokenType_GE:         "GE",         // >=
	TokenType_LT:         "LT",         // <
	TokenType_LE:         "LE",         // >=
	TokenType_EQ:         "EQ",         // ==
	TokenType_Assign:     "Assign",     // =
	TokenType_Plus:       "Plus",       // +
	TokenType_PlusPlus:   "PlusPlus",   // ++
	TokenType_PlusEQ:     "PlusEQ",     // -=
	TokenType_Minus:      "Minus",      // -
	TokenType_MinusMinus: "MinusMinus", // --
	TokenType_MinusEQ:    "MinusEQ",    // -=
	TokenType_Star:       "Star",       // *
	TokenType_StarEQ:     "StarEQ",     // *=
	TokenType_Slash:      "Slash",      // /
	TokenType_SlashEQ:    "SlashEQ",    // /=

	TokenType_LeftParen:  "LeftParen",  // (
	TokenType_RightParen: "RightParen", // )
	TokenType_SemiColon:  "SemiColon",  // ;

	TokenType_Identifier: "Identifier",
	TokenType_IntLiteral: "IntLiteral",
}

func (t TokenType) String() string {
	if v, ok := tokenTypeStringMap[t]; ok {
		return v
	}

	return "Unknown TokenType"
}

type SimpleLexer struct {
	tokenText []rune
	tokens    []Token
	token     Token
}

// tokenize parse string to generate Token, this is a finite state machine.
func (l *SimpleLexer) tokenize(script string) *TokenReader {
	l.tokenText = []rune{}
	state := DfaState_Initial

	for _, ch := range script {
		switch state {
		case DfaState_Initial:
			state = l.initToken(ch) // Reconfirm subsequent status
		case DfaState_Id:
			if isAlpha(ch) || isDigit(ch) { // keep Identifire state
				l.tokenText = append(l.tokenText, ch)
			} else {
				state = l.initToken(ch) // exit Identifier state
			}
		case DfaState_Int1:
			if ch == 'n' { // in
				state = DfaState_Int2
				l.tokenText = append(l.tokenText, ch)
			} else if isDigit(ch) || isAlpha(ch) {
				state = DfaState_Id
				l.tokenText = append(l.tokenText, ch)
			} else {
				state = l.initToken(ch)
			}
		case DfaState_Int2:
			if ch == 't' { // int
				state = DfaState_Int3
				l.tokenText = append(l.tokenText, ch)
			} else if isDigit(ch) || isAlpha(ch) {
				state = DfaState_Id
				l.tokenText = append(l.tokenText, ch)
			} else {
				state = l.initToken(ch)
			}
		case DfaState_Int3:
			if isBlank(ch) {
				l.token.Type = TokenType_INT
				state = l.initToken(ch)
			} else {
				state = DfaState_Id
				l.tokenText = append(l.tokenText, ch)
			}
		case DfaState_GT:
			if ch == '=' { // change state to GE
				state = DfaState_GE
				l.token.Type = TokenType_GE
				l.tokenText = append(l.tokenText, ch)
			} else {
				state = l.initToken(ch) // exit Identifier state
			}
		case DfaState_Assignment:
			if ch == '=' { // change state to GE
				state = DfaState_EQ
				l.token.Type = TokenType_EQ
				l.tokenText = append(l.tokenText, ch)
			} else {
				state = l.initToken(ch) // exit Identifier state
			}
		case DfaState_Plus:
			if ch == '=' { // change state to GE
				state = DfaState_PlusEQ
				l.token.Type = TokenType_PlusEQ
				l.tokenText = append(l.tokenText, ch)
			} else {
				state = l.initToken(ch) // exit Identifier state
			}
		case DfaState_Minus:
			if ch == '=' { // change state to GE
				state = DfaState_MinusEQ
				l.token.Type = TokenType_MinusEQ
				l.tokenText = append(l.tokenText, ch)
			} else {
				state = l.initToken(ch) // exit Identifier state
			}
		case DfaState_Star:
			if ch == '=' { // change state to GE
				state = DfaState_StarEQ
				l.token.Type = TokenType_StarEQ
				l.tokenText = append(l.tokenText, ch)
			} else {
				state = l.initToken(ch) // exit Identifier state
			}
		case DfaState_Slash:
			if ch == '=' { // change state to GE
				state = DfaState_SlashEQ
				l.token.Type = TokenType_SlashEQ
				l.tokenText = append(l.tokenText, ch)
			} else {
				state = l.initToken(ch) // exit Identifier state
			}
		case DfaState_IntLiteral:
			if isDigit(ch) {
				l.tokenText = append(l.tokenText, ch)
			} else {
				state = l.initToken(ch)
			}
		default:
			state = l.initToken(ch)
		}
	}

	if len(l.tokenText) > 0 {
		l.initToken(' ')
	}

	return NewTokenReader(l.tokens)
}

// initToken finite state machine enters initial state.
func (l *SimpleLexer) initToken(ch rune) DfaState {
	if len(l.tokenText) > 0 {
		l.token.Text = string(l.tokenText)
		l.tokens = append(l.tokens, l.token)

		l.token = Token{}
		l.tokenText = []rune{}
	}

	newState := DfaState_Initial
	if isAlpha(ch) { // the first char is alpha
		if ch == 'i' {
			newState = DfaState_Int1
		} else {
			newState = DfaState_Id
		}
		l.token.Type = TokenType_Identifier
		l.tokenText = append(l.tokenText, ch)
	} else if isDigit(ch) { // the first char is digit
		newState = DfaState_IntLiteral
		l.token.Type = TokenType_IntLiteral
		l.tokenText = append(l.tokenText, ch)
	} else if ch == '>' { // the first char is GT
		newState = DfaState_GT
		l.token.Type = TokenType_GT
		l.tokenText = append(l.tokenText, ch)
	} else if ch == '<' { // the first char is GT
		newState = DfaState_LT
		l.token.Type = TokenType_LT
		l.tokenText = append(l.tokenText, ch)
	} else if ch == '+' {
		newState = DfaState_Plus
		l.token.Type = TokenType_Plus
		l.tokenText = append(l.tokenText, ch)
	} else if ch == '-' {
		newState = DfaState_Minus
		l.token.Type = TokenType_Minus
		l.tokenText = append(l.tokenText, ch)
	} else if ch == '*' {
		newState = DfaState_Star
		l.token.Type = TokenType_Star
		l.tokenText = append(l.tokenText, ch)
	} else if ch == '/' {
		newState = DfaState_Slash
		l.token.Type = TokenType_Slash
		l.tokenText = append(l.tokenText, ch)
	} else if ch == '=' {
		newState = DfaState_Assignment
		l.token.Type = TokenType_Assign
		l.tokenText = append(l.tokenText, ch)
	} else if ch == ';' {
		newState = DfaState_SemiColon
		l.token.Type = TokenType_SemiColon
		l.tokenText = append(l.tokenText, ch)
	} else if ch == '(' {
		newState = DfaState_LeftParen
		l.token.Type = TokenType_LeftParen
		l.tokenText = append(l.tokenText, ch)
	} else if ch == ')' {
		newState = DfaState_RightParen
		l.token.Type = TokenType_RightParen
		l.tokenText = append(l.tokenText, ch)
	} else {
		newState = DfaState_Initial // skip all unknown patterns
	}

	return newState
}

// dump print all tokens
func (l *SimpleLexer) dump(reader *TokenReader) {
	fmt.Printf("%s\t%s\n", "text", "type")
	var token *Token
	for {
		if token = reader.Read(); token == nil {
			break
		} else {
			fmt.Printf("%s\t%v\n", token.Text, token.Type)
		}
	}
}

// TokenReader Token stream, encapsulates a token list.
type TokenReader struct {
	tokens   []Token
	position int
}

func NewTokenReader(tokens []Token) *TokenReader {
	return &TokenReader{tokens: tokens}
}

// Read return Token in stream and remove it from stream, return nil if the stream is empty.
func (r *TokenReader) Read() *Token {
	if r.position < len(r.tokens) {
		p := r.position
		r.position++
		return &r.tokens[p]
	}
	return nil
}

// Peek return Token in stream but does not remove it from stream, return null if the stream is empty.
func (r *TokenReader) Peek() *Token {
	if r.position < len(r.tokens) {
		return &r.tokens[r.position]
	}
	return nil
}

// UnRead Token stream goes back one position and restores the original token.
func (r *TokenReader) UnRead() {
	if r.position > 0 {
		r.position--
	}
}

// GetPosition get current reading position of Token stream.
func (s *TokenReader) GetPosition() int {
	return s.position
}

// setPosition set current reading position of Token stream.
func (r *TokenReader) setPosition(position int) {
	r.position = position
}

func isAlpha(ch int32) bool {
	return (ch >= 'a' && ch <= 'z') || (ch >= 'A' && ch <= 'Z')
}

func isDigit(ch int32) bool {
	return ch >= '0' && ch <= '9'
}

func isBlank(ch int32) bool {
	return ch == ' ' || ch == '\t' || ch == '\n'
}
