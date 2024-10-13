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
	case Star, Slash, Minus, Greater, GreaterEqual, Less, LessEqual:
		if !(okLeftNumber && okRightNumber) {
			return nil, fmt.Errorf("cannot operate the non-numbers '%v' and/or '%v'", leftExpr, rightExpr)
		}

		switch e.operator.Type {
		case Star:
			return leftNumber * rightNumber, err
		case Slash:
			return leftNumber / rightNumber, err
		case Minus:
			return leftNumber - rightNumber, err
		case Greater:
			return leftNumber > rightNumber, err
		case GreaterEqual:
			return leftNumber >= rightNumber, err
		case Less:
			return leftNumber < rightNumber, err
		case LessEqual:
			return leftNumber <= rightNumber, err
		}
		panic("unreachable")
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
	case EqualEqual:
		return leftExpr == rightExpr, err
	case BangEqual:
		return leftExpr != rightExpr, err
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
			return nil, fmt.Errorf("Operand must be a number.")
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
