package main

import "fmt"

type Environment struct {
	Values map[string]interface{}
}

func (e *Environment) define(name string, value interface{}) {
	e.Values[name] = value
}

func (e *Environment) get(name Token) (interface{}, error) {
	if val, ok := e.Values[name.Lexeme]; ok {
		return val, nil
	}

	return nil, fmt.Errorf("undeclared")
}
