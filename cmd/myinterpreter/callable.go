package main

import "time"

type LoxCallable interface {
	Arity() int
	Call(itp *AstInterpreter, arguments []interface{}) (interface{}, error)
}

type ClockFunc struct{}

func (c ClockFunc) Arity() int {
	return 0
}

func (c ClockFunc) Call(itp *AstInterpreter, arguments []interface{}) (interface{}, error) {
	return float64(time.Now().Unix()), nil
}

func (c ClockFunc) String() string {
	return "<native fn>"
}
