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
	return fmt.Errorf("Undefined variable '" + name.Lexeme + "'")
}

func (e *Environment) Get(name Token) (interface{}, error) {
	if e.Enclosing != nil {
		return e.Enclosing.Get(name)
	}

	if val, ok := e.Values[name.Lexeme]; ok {
		return val, nil
	}

	return nil, fmt.Errorf("undeclared")
}

func (e *Environment) Clone() Environment {
	cloned := &Environment{
		Enclosing: e.Enclosing,
		Values:    make(map[string]interface{}, len(e.Values)),
	}

	for k, v := range e.Values {
		cloned.Values[k] = v
	}

	return *cloned
}
