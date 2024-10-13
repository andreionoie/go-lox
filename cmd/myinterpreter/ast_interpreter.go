package main

import "fmt"

type AstInterpreter struct {
	StubExprVisitor
}

func (itp *AstInterpreter) VisitBinaryExpr(e *BinaryExpr) (result interface{}, err error) {
	leftExpr, err := e.left.Accept(itp)
	rightExpr, err := e.right.Accept(itp)

	leftNumber, okLeft := leftExpr.(float64)
	rightNumber, okRight := rightExpr.(float64)
	if !(okLeft && okRight) {
		return nil, fmt.Errorf("cannot minus the expression '%v'", rightExpr)
	}

	switch e.operator.Type {
	case Star:
		return leftNumber * rightNumber, err
	case Slash:
		return leftNumber / rightNumber, err
	}
	panic("Unsupported binary operator!")
}

func (itp *AstInterpreter) VisitUnaryExpr(e *UnaryExpr) (result interface{}, err error) {
	rightExpr, err := e.right.Accept(itp)
	switch e.operator.Type {
	case Bang:
		return !isTruthy(rightExpr), err
	case Minus:
		rightNumber, ok := rightExpr.(float64)
		if !ok {
			return nil, fmt.Errorf("cannot minus the expression '%v'", rightExpr)
		}
		return -rightNumber, err
	}
	panic("Unsupported unary operator!")
}

func (itp *AstInterpreter) VisitGroupingExpr(e *GroupingExpr) (result interface{}, err error) {
	return e.expr.Accept(itp)
}

func (itp *AstInterpreter) VisitLiteralExpr(e *LiteralExpr) (result interface{}, err error) {
	return e.value, nil
}

func isTruthy(val interface{}) bool {
	if val == nil {
		return false
	}
	valAsBool, ok := val.(bool)
	if ok {
		return valAsBool
	}

	return true
}
