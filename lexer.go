package main

import (
	"errors"
	"fmt"
	"unicode"
)

type TokenType string

const (
	OpenCurly    TokenType = "OPEN_CURLY"
	CloseCurly   TokenType = "CLOSE_CURLY"
	OpenBracket  TokenType = "OPEN_BRACKET"
	CloseBracket TokenType = "CLOSE_BRACKET"
	Colon        TokenType = "COLON"
	Comma        TokenType = "COMMA"
	Int          TokenType = "INT"
	Float        TokenType = "FLOAT"
	Boolean      TokenType = "BOOLEAN"
	Null         TokenType = "NULL"
	String       TokenType = "STRING"
	Eof          TokenType = "EOF"
)

type Token struct {
	Type  TokenType
	Value string
}

// Lexer tokenizes JSON input into tokens
type Lexer struct {
	input       string
	pos         int
	currentChar byte
}

func NewLexer(input string) *Lexer {
	var ch byte
	if len(input) > 0 {
		ch = input[0]
	}
	return &Lexer{
		input:       input,
		pos:         0,
		currentChar: ch,
	}
}

func (l *Lexer) advance() {
	l.pos++
	if l.pos >= len(l.input) {
		l.currentChar = 0 // EOF
	} else {
		l.currentChar = l.input[l.pos]
	}
}

func (l *Lexer) readString() (Token, error) {
	// Skip opening quote
	l.advance()

	startPos := l.pos

	for l.currentChar != 0 && l.currentChar != '"' {
		l.advance()
	}

	if l.currentChar != '"' {
		return Token{}, errors.New("unterminated string: missing closing quote")
	}

	value := l.input[startPos:l.pos]
	l.advance() // Skip closing quote

	return Token{Type: String, Value: value}, nil
}

func (l *Lexer) readNumber() (Token, error) {
	startPos := l.pos

	for l.currentChar >= '0' && l.currentChar <= '9' {
		l.advance()
	}

	if l.currentChar == '.' {
		l.advance()

		for l.currentChar >= '0' && l.currentChar <= '9' {
			l.advance()
		}

		value := l.input[startPos:l.pos]
		return Token{Type: Float, Value: value}, nil
	}

	value := l.input[startPos:l.pos]
	return Token{Type: Int, Value: value}, nil
}

func (l *Lexer) readBoolean() (Token, error) {
	var expectedString string
	switch l.currentChar {
	case 't':
		expectedString = "true"
	case 'f':
		expectedString = "false"
	}

	for i := 0; i < len(expectedString); i++ {
		if l.currentChar != expectedString[i] {
			return Token{}, fmt.Errorf("unexpected char '%c'", l.currentChar)
		}
		l.advance()
	}

	return Token{Type: Boolean, Value: expectedString}, nil
}

func (l *Lexer) readNull() (Token, error) {
	expected := "null"

	for i := 0; i < len(expected); i++ {
		if l.currentChar != expected[i] {
			return Token{}, fmt.Errorf("unexpected char '%c'", l.currentChar)
		}
		l.advance()
	}

	return Token{Type: Null, Value: "null"}, nil
}

// NextToken returns the next token from the input
func (l *Lexer) NextToken() (Token, error) {
	for l.pos < len(l.input) && unicode.IsSpace(rune(l.currentChar)) {
		l.advance()
	}

	if l.pos >= len(l.input) {
		return Token{Type: Eof, Value: ""}, nil
	}

	switch l.currentChar {
	case '{':
		l.advance()
		return Token{Type: OpenCurly, Value: "{"}, nil
	case '"':
		return l.readString()
	case ':':
		l.advance()
		return Token{Type: Colon, Value: ":"}, nil
	case '}':
		l.advance()
		return Token{Type: CloseCurly, Value: "}"}, nil
	case '[':
		l.advance()
		return Token{Type: OpenBracket, Value: "["}, nil
	case ']':
		l.advance()
		return Token{Type: CloseBracket, Value: "]"}, nil
	case ',':
		l.advance()
		return Token{Type: Comma, Value: ","}, nil
	case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
		return l.readNumber()
	case 't', 'f':
		return l.readBoolean()
	case 'n':
		return l.readNull()
	default:
		return Token{}, fmt.Errorf("unknown character '%c'", l.currentChar)
	}
}
