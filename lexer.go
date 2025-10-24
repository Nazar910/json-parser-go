package main

import "errors"

type TokenType string

const (
	OPEN_CURLY  TokenType = "OPEN_CURLY"
	CLOSE_CURLY TokenType = "CLOSE_CURLY"
	COLON       TokenType = "COLON"
	STRING      TokenType = "STRING"
	EOF         TokenType = "EOF"
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
	default:
		return Token{}, errors.New("unknown character")
	}
}
