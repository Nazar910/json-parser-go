package main

import "testing"

func assertToken(t *testing.T, actual Token, expected Token) {
	if actual.Type != expected.Type {
		t.Errorf("Expected token type to be %s but got %s", expected.Type, actual.Type)
	}

	if actual.Value != expected.Value {
		t.Errorf("Expected token Value to be %s but got %s", expected.Value, actual.Value)
	}
}

func TestLexerGetNextToken(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []Token
	}{
		{
			name:  "basic foo bar test",
			input: "{\"foo\":\"bar\"}",
			expected: []Token{
				{Type: OPEN_CURLY, Value: "{"},
				{Type: STRING, Value: "foo"},
				{Type: COLON, Value: ":"},
				{Type: STRING, Value: "bar"},
				{Type: CLOSE_CURLY, Value: "}"},
				{Type: EOF, Value: ""},
			},
		},
		{
			name:  "array with all types",
			input: "[\"string\",1,2.0,true,false,null]",
			expected: []Token{
				{Type: OPEN_CURLY, Value: "{"},
				{Type: STRING, Value: "foo"},
				{Type: COLON, Value: ":"},
				{Type: STRING, Value: "bar"},
				{Type: CLOSE_CURLY, Value: "}"},
				{Type: EOF, Value: ""},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			lexer := NewLexer(test.input)
			for _, expectedToken := range test.expected {
				token, err := lexer.GetNextToken()
				if err != nil {
					t.Errorf("Unexpected error: %s", err)
				}

				assertToken(t, token, expectedToken)
			}
		})
	}
}
