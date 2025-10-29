package main

import (
	"testing"
)

func TestParser(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected JSONValue
	}{
		{
			name:     "basic foo bar test",
			input:    "{\"foo\":\"bar\"}",
			expected: JSONObject{Fields: map[string]JSONValue{"foo": JSONString{Value: "bar"}}},
		},
		{
			name:  "array with all types",
			input: "[\"string\",1,2.5,true,false,null]",
			expected: JSONArray{Elements: []JSONValue{
				JSONString{Value: "string"},
				JSONInt{Value: 1},
				JSONFloat{Value: 2.5},
				JSONBool{Value: true},
				JSONBool{Value: false},
				JSONNull{},
			}},
		},
		{
			name:  "array with nested object and array",
			input: "[{\"foo\":\"bar\"},[1,2,3]]",
			expected: JSONArray{Elements: []JSONValue{
				JSONObject{Fields: map[string]JSONValue{"foo": JSONString{Value: "bar"}}},
				JSONArray{Elements: []JSONValue{JSONInt{Value: 1}, JSONInt{Value: 2}, JSONInt{Value: 3}}},
			}},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			parser := NewParser(test.input)

			actual, error := parser.Parse()

			if error != nil {
				t.Errorf("unexpected error %s", error)
			}

			if !actual.Equals(test.expected) {
				t.Errorf("Expected %v, got %v", test.expected, actual)
			}
		})
	}
}
