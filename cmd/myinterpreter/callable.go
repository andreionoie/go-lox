package main

import (
	"fmt"
	"time"
)

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

type LoxFunction struct {
	declaration *FunctionStmt
}

func (lf LoxFunction) Arity() int {
	return len(lf.declaration.parameters)
}

func (lf LoxFunction) Call(itp *AstInterpreter, arguments []interface{}) (interface{}, error) {
	previousEnv := itp.env
	itp.env = &Environment{
		Enclosing: itp.Globals,
		Values:    make(map[string]interface{}),
	}
	defer func() { itp.env = previousEnv }()

	for i := 0; i < len(arguments); i++ {
		itp.env.Define(lf.declaration.parameters[i].Lexeme, arguments[i])
	}

	for _, funcStmt := range lf.declaration.body {
		funcStmt.Accept(itp)
	}

	return nil, nil
}

func (lf LoxFunction) String() string {
	return fmt.Sprintf("<fn %s>", lf.declaration.name.Lexeme)
}
