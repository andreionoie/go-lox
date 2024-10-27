package main

import "fmt"

type Environment struct {
	Enclosing *Environment
	Values    map[string]interface{}
}

func (e *Environment) Define(name string, value interface{}) {
	e.Values[name] = value
}

func (e *Environment) Assign(name Token, value interface{}) error {
	if _, ok := e.Values[name.Lexeme]; ok {
		e.Values[name.Lexeme] = value
		return nil
	}

	if e.Enclosing != nil {
		return e.Enclosing.Assign(name, value)
	}

	return fmt.Errorf("Undefined variable '" + name.Lexeme + "'.")
}

func (e *Environment) Get(name Token) (interface{}, error) {
	if val, ok := e.Values[name.Lexeme]; ok {
		return val, nil
	}

	if e.Enclosing != nil {
		return e.Enclosing.Get(name)
	}

	return nil, fmt.Errorf("undeclared")
}
