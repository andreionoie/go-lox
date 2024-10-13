package main

import "fmt"

type AstInterpreter struct {
	StubExprVisitor
}

func (itp *AstInterpreter) VisitBinaryExpr(e *BinaryExpr) (result interface{}, err error) {
	leftExpr, err := e.left.Accept(itp)
	rightExpr, err := e.right.Accept(itp)

	leftNumber, okLeftNumber := leftExpr.(float64)
	rightNumber, okRightNumber := rightExpr.(float64)

	switch e.operator.Type {
	case Star, Slash, Minus:
		if !(okLeftNumber && okRightNumber) {
			return nil, fmt.Errorf("cannot operate the non-numbers '%v' and/or '%v'", leftExpr, rightExpr)
		}

		if e.operator.Type == Star {
			return leftNumber * rightNumber, err
		} else if e.operator.Type == Slash {
			return leftNumber / rightNumber, err
		} else if e.operator.Type == Minus {
			return leftNumber - rightNumber, err
		}
	case Plus:
		if okLeftNumber && okRightNumber {
			return leftNumber + rightNumber, err
		}
		leftString, okLeft := leftExpr.(string)
		rightString, okRight := rightExpr.(string)
		if okLeft && okRight {
			return leftString + rightString, err
		}
		return nil, fmt.Errorf("cannot plus the expressions '%v' with '%v'", leftExpr, rightExpr)
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
