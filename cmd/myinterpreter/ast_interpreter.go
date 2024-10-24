package main

import (
	"fmt"
	"os"
)

type AstInterpreter struct {
	StubExprVisitor
	StubStmtVisitor
}

func (itp *AstInterpreter) Interpret(stmts []Stmt) {
	for _, stmt := range stmts {
		_, err := stmt.Accept(itp)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			LoxHadRuntimeError = true
			return
		}
	}
}

func (itp *AstInterpreter) InterpretExpr(e Expr) {
	result, err := e.Accept(itp)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		LoxHadRuntimeError = true
	} else {
		if result == nil {
			fmt.Println("nil")
		} else {
			fmt.Println(result)
		}
	}
}

func (itp *AstInterpreter) VisitPrintStmt(s *PrintStmt) (result interface{}, err error) {
	result, err = s.expression.Accept(itp)
	if err == nil {
		if result == nil {
			fmt.Println("nil")
		} else {
			fmt.Println(result)
		}
	}
	return nil, err
}

func (itp *AstInterpreter) VisitExpressionStmt(s *ExpressionStmt) (result interface{}, err error) {
	_, err = s.expression.Accept(itp)
	return nil, err
}

func (itp *AstInterpreter) VisitBinaryExpr(e *BinaryExpr) (result interface{}, err error) {
	leftExpr, err := e.left.Accept(itp)
	rightExpr, err := e.right.Accept(itp)

	leftNumber, okLeftNumber := leftExpr.(float64)
	rightNumber, okRightNumber := rightExpr.(float64)

	switch e.operator.Type {
	case Star, Slash, Minus, Greater, GreaterEqual, Less, LessEqual:
		if !(okLeftNumber && okRightNumber) {
			return nil, fmt.Errorf("Operands must be numbers.")
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
		return nil, fmt.Errorf("Operands must be two numbers or two strings")
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
