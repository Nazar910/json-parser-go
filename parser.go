package main

import (
	"errors"
	"fmt"
	"strconv"
)

// Define a proper JSON AST (Abstract Syntax Tree)
type JSONValue interface {
	Equals(other JSONValue) bool
}

// Concrete types for each JSON type
type JSONNull struct{}
type JSONBool struct{ Value bool }
type JSONInt struct{ Value int }
type JSONFloat struct{ Value float64 }
type JSONString struct{ Value string }
type JSONArray struct{ Elements []JSONValue }
type JSONObject struct{ Fields map[string]JSONValue }

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

func (a JSONArray) Equals(other JSONValue) bool {
	if other, ok := other.(JSONArray); ok {
		if len(a.Elements) != len(other.Elements) {
			return false
		}

		for i := range a.Elements {
			if !a.Elements[i].Equals(other.Elements[i]) {
				return false
			}
		}

		return true
	}

	return false
}

func (num JSONInt) Equals(other JSONValue) bool {
	otherInt, ok := other.(JSONInt)

	return ok && num.Value == otherInt.Value
}

func (num JSONFloat) Equals(other JSONValue) bool {
	otherFloat, ok := other.(JSONFloat)

	return ok && num.Value == otherFloat.Value
}

type Parser struct {
	lexer        Lexer
	currentToken Token
}

func NewParser(input string) Parser {
	return Parser{lexer: *NewLexer(input)}
}

func (p *Parser) Parse() (JSONValue, error) {
	currentToken, error := p.lexer.NextToken()
	p.currentToken = currentToken

	if error != nil {
		return JSONNull{}, error
	}

	if p.currentToken.Type == OpenCurly {
		return p.object()
	}

	if p.currentToken.Type == OpenBracket {
		return p.array()
	}

	return JSONNull{}, fmt.Errorf("Unsupported token %v", p.currentToken)
}

func (p *Parser) array() (JSONArray, error) {
	err := p.eat(OpenBracket)

	if err != nil {
		return JSONArray{}, err
	}

	// TODO: possible special case for empty array

	arr := JSONArray{Elements: []JSONValue{}}

	elem, err := p.value()

	if err != nil {
		return JSONArray{}, err
	}

	arr.Elements = append(arr.Elements, elem)

	for p.currentToken.Type == Comma {
		err := p.eat(Comma)

		if err != nil {
			return JSONArray{}, err
		}

		elem, err := p.value()

		if err != nil {
			return JSONArray{}, err
		}

		arr.Elements = append(arr.Elements, elem)
	}

	return arr, nil
}

func (p *Parser) object() (JSONObject, error) {
	err := p.eat(OpenCurly)

	if err != nil {
		return JSONObject{}, err
	}

	key, err := p.string()

	if err != nil {
		return JSONObject{}, err
	}

	err = p.eat(Colon)

	if err != nil {
		return JSONObject{}, err
	}

	value, err := p.value()

	if err != nil {
		return JSONObject{}, err
	}

	err = p.eat(CloseCurly)

	if err != nil {
		return JSONObject{}, err
	}

	fieldsMap := make(map[string]JSONValue)
	fieldsMap[key.Value] = value

	return JSONObject{Fields: fieldsMap}, nil
}

func (p *Parser) value() (JSONValue, error) {
	if p.currentToken.Type == String {
		return p.string()
	}

	if p.currentToken.Type == Int {
		return p.int()
	}

	if p.currentToken.Type == Float {
		return p.float()
	}

	if p.currentToken.Type == Boolean {
		return p.boolean()
	}

	if p.currentToken.Type == Null {
		return p.null()
	}

	if p.currentToken.Type == OpenCurly {
		return p.object()
	}

	if p.currentToken.Type == OpenBracket {
		return p.array()
	}

	return JSONNull{}, fmt.Errorf("Unsupported token type %s", p.currentToken.Type)
}

func (p *Parser) int() (JSONInt, error) {
	intStr := p.currentToken.Value

	p.eat(Int)

	intNum, err := strconv.ParseInt(intStr, 10, 32)

	return JSONInt{Value: int(intNum)}, err
}

func (p *Parser) float() (JSONFloat, error) {
	floatStr := p.currentToken.Value

	p.eat(Float)

	floatNum, err := strconv.ParseFloat(floatStr, 32)

	return JSONFloat{Value: floatNum}, err
}

func (p *Parser) boolean() (JSONBool, error) {
	booleanStr := p.currentToken.Value

	p.eat(Boolean)

	return JSONBool{Value: booleanStr == "true"}, nil
}

func (p *Parser) null() (JSONNull, error) {
	p.eat(Null)

	return JSONNull{}, nil
}

func (p *Parser) string() (JSONString, error) {
	str := p.currentToken.Value
	err := p.eat(String)
	if err != nil {
		return JSONString{}, err
	}
	return JSONString{Value: str}, nil
}

func (p *Parser) eat(expectedType TokenType) error {
	if p.currentToken.Type == expectedType {
		nextToken, error := p.lexer.NextToken()
		if error != nil {
			return error
		}

		p.currentToken = nextToken
		return nil
	}
	return errors.New("Syntax error")
}
