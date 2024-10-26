package main

import "fmt"

type Environment struct {
	Values map[string]interface{}
}

func (e *Environment) define(name string, value interface{}) {
	e.Values[name] = value
}

func (e *Environment) assign(name Token, value interface{}) error {
	if _, ok := e.Values[name.Lexeme]; ok {
		e.Values[name.Lexeme] = value
		return nil
	}

	return fmt.Errorf("Undefined variable '" + name.Lexeme + "'")
}

func (e *Environment) get(name Token) (interface{}, error) {
	if val, ok := e.Values[name.Lexeme]; ok {
		return val, nil
	}

	return nil, fmt.Errorf("undeclared")
}
