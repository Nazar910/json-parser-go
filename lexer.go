package main

import (
	"errors"
	"fmt"
)

type TokenType string

const (
	OPEN_CURLY    TokenType = "OPEN_CURLY"
	CLOSE_CURLY   TokenType = "CLOSE_CURLY"
	OPEN_BRACKET  TokenType = "OPEN_BRACKET"
	CLOSE_BRACKET TokenType = "CLOSE_BRACKET"
	COLON         TokenType = "COLON"
	COMMA         TokenType = "COMMA"
	INT           TokenType = "INT"
	FLOAT         TokenType = "FLOAT"
	BOOLEAN       TokenType = "BOOLEAN"
	NULL          TokenType = "NULL"
	STRING        TokenType = "STRING"
	EOF           TokenType = "EOF"
)

type Token struct {
	Type  TokenType
	Value string
}

type Lexer struct {
	input       string
	pos         int
	currentChar byte
}

func NewLexer(input string) *Lexer {
	return &Lexer{
		input:       input,
		pos:         0,
		currentChar: input[0],
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

	return Token{Type: STRING, Value: value}, nil
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
		return Token{Type: FLOAT, Value: value}, nil
	}

	value := l.input[startPos:l.pos]
	return Token{Type: INT, Value: value}, nil
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
			return Token{}, fmt.Errorf("Unexpected char '%c'", l.currentChar)
		}
		l.advance()
	}

	return Token{Type: BOOLEAN, Value: expectedString}, nil
}

func (l *Lexer) readNull() (Token, error) {
	expected := "null"

	for i := 0; i < len(expected); i++ {
		if l.currentChar != expected[i] {
			return Token{}, fmt.Errorf("Unexpected char '%c'", l.currentChar)
		}
		l.advance()
	}

	return Token{Type: NULL, Value: "null"}, nil
}

func (l *Lexer) GetNextToken() (Token, error) {
	if l.pos >= len(l.input) {
		return Token{Type: EOF, Value: ""}, nil
	}

	switch l.currentChar {
	case '{':
		l.advance()
		return Token{Type: OPEN_CURLY, Value: "{"}, nil
	case '"':
		return l.readString()
	case ':':
		l.advance()
		return Token{Type: COLON, Value: ":"}, nil
	case '}':
		l.advance()
		return Token{Type: CLOSE_CURLY, Value: "}"}, nil
	case '[':
		l.advance()
		return Token{Type: OPEN_BRACKET, Value: "["}, nil
	case ']':
		l.advance()
		return Token{Type: CLOSE_BRACKET, Value: "]"}, nil
	case ',':
		l.advance()
		return Token{Type: COMMA, Value: ","}, nil
	case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
		return l.readNumber()
	case 't', 'f':
		return l.readBoolean()
	case 'n':
		return l.readNull()
	default:
		return Token{}, fmt.Errorf("Unknown character '%c'", l.currentChar)
	}
}
