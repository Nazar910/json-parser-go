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
			name:  "array with all primitive types",
			input: "[\"string\",1,2.5,true,false,null]",
			expected: []Token{
				{Type: OPEN_BRACKET, Value: "["},
				{Type: STRING, Value: "string"},
				{Type: COMMA, Value: ","},
				{Type: INT, Value: "1"},
				{Type: COMMA, Value: ","},
				{Type: FLOAT, Value: "2.5"},
				{Type: COMMA, Value: ","},
				{Type: BOOLEAN, Value: "true"},
				{Type: COMMA, Value: ","},
				{Type: BOOLEAN, Value: "false"},
				{Type: COMMA, Value: ","},
				{Type: NULL, Value: "null"},
				{Type: CLOSE_BRACKET, Value: "]"},
				{Type: EOF, Value: ""},
			},
		},
		{
			name:  "array with nested object and array",
			input: "[{\"foo\":\"bar\"},[1,2,3]]",
			expected: []Token{
				{Type: OPEN_BRACKET, Value: "["},
				{Type: OPEN_CURLY, Value: "{"},
				{Type: STRING, Value: "foo"},
				{Type: COLON, Value: ":"},
				{Type: STRING, Value: "bar"},
				{Type: CLOSE_CURLY, Value: "}"},
				{Type: COMMA, Value: ","},
				{Type: OPEN_BRACKET, Value: "["},
				{Type: INT, Value: "1"},
				{Type: COMMA, Value: ","},
				{Type: INT, Value: "2"},
				{Type: COMMA, Value: ","},
				{Type: INT, Value: "3"},
				{Type: CLOSE_BRACKET, Value: "]"},
				{Type: CLOSE_BRACKET, Value: "]"},
				{Type: EOF, Value: ""},
			},
		},
		{
			name:  "should ignore whitespaces",
			input: "{  \"foo\" :     1}",
			expected: []Token{
				{Type: OPEN_CURLY, Value: "{"},
				{Type: STRING, Value: "foo"},
				{Type: COLON, Value: ":"},
				{Type: INT, Value: "1"},
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
