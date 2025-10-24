package main

import (
	"errors"
	"fmt"
)

// Define a proper JSON AST (Abstract Syntax Tree)
type JSONValue interface {
	jsonValue() // marker method
	Equals(other JSONValue) bool
}

// Concrete types for each JSON type
type JSONNull struct{}
type JSONBool struct{ Value bool }
type JSONNumber struct{ Value float64 }
type JSONString struct{ Value string }
type JSONArray struct{ Elements []JSONValue }
type JSONObject struct{ Fields map[string]JSONValue }

// Implement the interface
func (JSONNull) jsonValue()   {}
func (JSONBool) jsonValue()   {}
func (JSONNumber) jsonValue() {}
func (JSONString) jsonValue() {}
func (JSONArray) jsonValue()  {}
func (JSONObject) jsonValue() {}

func (n JSONNull) Equals(other JSONValue) bool {
	_, ok := other.(JSONNull)
	return ok
}

func (n JSONBool) Equals(other JSONValue) bool {
	if otherBool, ok := other.(JSONBool); ok {
		return n.Value == otherBool.Value
	}

	return false
}

func (n JSONString) Equals(other JSONValue) bool {
	if otherString, ok := other.(JSONString); ok {
		return n.Value == otherString.Value
	}
	return false
}

func (o JSONObject) Equals(other JSONValue) bool {
	if otherObject, ok := other.(JSONObject); ok {
		if len(o.Fields) != len(otherObject.Fields) {
			return false
		}

		for key, value := range o.Fields {
			if otherValue, exists := otherObject.Fields[key]; !exists || !otherValue.Equals(value) {
				return false
			}
		}
		return true
	}
	return false
}

type Parser struct {
	lexer        Lexer
	currentToken Token
}

func NewParser(input string) Parser {
	return Parser{lexer: *NewLexer(input)}
}

func (p *Parser) Parse() (JSONValue, error) {
	currentToken, error := p.lexer.GetNextToken()
	p.currentToken = currentToken

	if error != nil {
		return JSONNull{}, error
	}

	if p.currentToken.Type == OPEN_CURLY {
		return p.object()
	}

	return JSONNull{}, fmt.Errorf("Unsupported token %v", p.currentToken)
}

func (p *Parser) object() (JSONObject, error) {
	err := p.eat(OPEN_CURLY)

	if err != nil {
		return JSONObject{}, err
	}

	key, err := p.string()

	if err != nil {
		return JSONObject{}, err
	}

	err = p.eat(COLON)

	if err != nil {
		return JSONObject{}, err
	}

	value, err := p.value()

	if err != nil {
		return JSONObject{}, err
	}

	err = p.eat(CLOSE_CURLY)

	if err != nil {
		return JSONObject{}, err
	}

	fieldsMap := make(map[string]JSONValue)
	fieldsMap[key.Value] = value

	return JSONObject{Fields: fieldsMap}, nil
}

func (p *Parser) value() (JSONValue, error) {
	if p.currentToken.Type == STRING {
		return p.string()
	}

	return JSONNull{}, errors.New("Unsupported token type")
}

func (p *Parser) string() (JSONString, error) {
	str := p.currentToken.Value
	err := p.eat(STRING)
	if err != nil {
		return JSONString{}, err
	}
	return JSONString{Value: str}, nil
}

func (p *Parser) eat(expectedType TokenType) error {
	if p.currentToken.Type == expectedType {
		nextToken, error := p.lexer.GetNextToken()
		if error != nil {
			return error
		}

		p.currentToken = nextToken
		return nil
	}
	return errors.New("Syntax error")
}
