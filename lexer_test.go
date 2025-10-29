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
				{Type: OpenCurly, Value: "{"},
				{Type: String, Value: "foo"},
				{Type: Colon, Value: ":"},
				{Type: String, Value: "bar"},
				{Type: CloseCurly, Value: "}"},
				{Type: Eof, Value: ""},
			},
		},
		{
			name:  "array with all primitive types",
			input: "[\"string\",1,2.5,true,false,null]",
			expected: []Token{
				{Type: OpenBracket, Value: "["},
				{Type: String, Value: "string"},
				{Type: Comma, Value: ","},
				{Type: Int, Value: "1"},
				{Type: Comma, Value: ","},
				{Type: Float, Value: "2.5"},
				{Type: Comma, Value: ","},
				{Type: Boolean, Value: "true"},
				{Type: Comma, Value: ","},
				{Type: Boolean, Value: "false"},
				{Type: Comma, Value: ","},
				{Type: Null, Value: "null"},
				{Type: CloseBracket, Value: "]"},
				{Type: Eof, Value: ""},
			},
		},
		{
			name:  "array with nested object and array",
			input: "[{\"foo\":\"bar\"},[1,2,3]]",
			expected: []Token{
				{Type: OpenBracket, Value: "["},
				{Type: OpenCurly, Value: "{"},
				{Type: String, Value: "foo"},
				{Type: Colon, Value: ":"},
				{Type: String, Value: "bar"},
				{Type: CloseCurly, Value: "}"},
				{Type: Comma, Value: ","},
				{Type: OpenBracket, Value: "["},
				{Type: Int, Value: "1"},
				{Type: Comma, Value: ","},
				{Type: Int, Value: "2"},
				{Type: Comma, Value: ","},
				{Type: Int, Value: "3"},
				{Type: CloseBracket, Value: "]"},
				{Type: CloseBracket, Value: "]"},
				{Type: Eof, Value: ""},
			},
		},
		{
			name:  "should ignore whitespaces",
			input: "{  \"foo\" :     1}",
			expected: []Token{
				{Type: OpenCurly, Value: "{"},
				{Type: String, Value: "foo"},
				{Type: Colon, Value: ":"},
				{Type: Int, Value: "1"},
				{Type: CloseCurly, Value: "}"},
				{Type: Eof, Value: ""},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			lexer := NewLexer(test.input)
			for _, expectedToken := range test.expected {
				token, err := lexer.NextToken()
				if err != nil {
					t.Errorf("Unexpected error: %s", err)
				}

				assertToken(t, token, expectedToken)
			}
		})
	}
}
