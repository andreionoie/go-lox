package main

import "fmt"

type AstInterpreter struct {
	StubExprVisitor
}

func (itp *AstInterpreter) VisitLiteralExpr(e *LiteralExpr) (result interface{}, err error) {
	if e.value == nil {
		return "nil", nil
	}

	return fmt.Sprintf("%v", e.value), nil
}
