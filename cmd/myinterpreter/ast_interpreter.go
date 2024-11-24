package main

import (
	"fmt"
	"os"
	"strconv"
)

type AstInterpreter struct {
	StubExprVisitor
	StubStmtVisitor
	Globals *Environment
	env     *Environment
}

func NewInterpreter() *AstInterpreter {
	initialEnv := &Environment{
		Values: make(map[string]interface{}),
	}

	initialEnv.Define("clock", ClockFunc{})

	return &AstInterpreter{
		Globals: initialEnv,
		env:     initialEnv,
	}
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

func (itp *AstInterpreter) VisitFunctionStmt(s *FunctionStmt) (result interface{}, err error) {
	loxFunc := &LoxFunction{
		declaration: s,
		closure:     itp.env, // store in memory the environment (hierarchy) that was active on function declaration
	}
	itp.env.Define(s.name.Lexeme, loxFunc)
	return nil, nil
}

func (itp *AstInterpreter) VisitReturnStmt(s *ReturnStmt) (result interface{}, err error) {
	var returnValue interface{} // empty return value defaults to nil
	if s.value != nil {
		returnValue, err = s.value.Accept(itp)
		if err != nil {
			return nil, err
		}
	}

	return nil, &ReturnUnwindCallstack{returnValue}
}

// TODO: fix useless result for statements (remove)
func (itp *AstInterpreter) VisitForStmt(s *ForStmt) (result interface{}, err error) {
	if s.init != nil {
		_, err = s.init.Accept(itp)
		if err != nil {
			return nil, err
		}
	}

	var condition Expr
	if s.condition != nil {
		condition = s.condition
	} else {
		condition = &LiteralExpr{value: true}
	}

	condResult, err := condition.Accept(itp)
	if err != nil {
		return nil, err
	}

	for isTruthy(condResult) {
		_, err = s.loopBody.Accept(itp)
		if err != nil {
			return nil, err
		}

		if s.iteration != nil {
			_, err = s.iteration.Accept(itp)
			if err != nil {
				return nil, err
			}
		}

		condResult, err = condition.Accept(itp)
		if err != nil {
			return nil, err
		}
	}

	return nil, nil
}

func (itp *AstInterpreter) VisitWhileStmt(s *WhileStmt) (result interface{}, err error) {
	condResult, err := s.condition.Accept(itp)
	if err != nil {
		return nil, err
	}

	for isTruthy(condResult) {
		_, err = s.loopBody.Accept(itp)
		if err != nil {
			return nil, err
		}

		condResult, err = s.condition.Accept(itp)
		if err != nil {
			return nil, err
		}
	}

	return nil, nil
}

func (itp *AstInterpreter) VisitIfStmt(s *IfStmt) (result interface{}, err error) {
	condResult, err := s.condition.Accept(itp)
	if err != nil {
		return nil, err
	}

	if isTruthy(condResult) {
		return s.thenBranch.Accept(itp)
	} else if s.elseBranch != nil {
		return s.elseBranch.Accept(itp)
	}

	// if cond was false and no else branch, noop
	return nil, nil
}

func (itp *AstInterpreter) VisitBlockStmt(s *BlockStmt) (result interface{}, err error) {
	previousEnv := itp.env
	itp.env = &Environment{
		Enclosing: previousEnv,
		Values:    make(map[string]interface{}),
	}
	defer func() { itp.env = previousEnv }()

	for _, stmt := range s.statements {
		result, err = stmt.Accept(itp)
		if err != nil {
			return nil, err
		}
	}

	return nil, err
}

func (itp *AstInterpreter) VisitVarStmt(s *VarStmt) (result interface{}, err error) {
	var varValue interface{}
	if s.initializerExpression != nil {
		varValue, err = s.initializerExpression.Accept(itp)
	}

	itp.env.Define(s.varName.Lexeme, varValue)
	return nil, err
}

func (itp *AstInterpreter) VisitPrintStmt(s *PrintStmt) (result interface{}, err error) {
	result, err = s.expression.Accept(itp)
	if err == nil {
		fmt.Println(loxStringify(result))
	}
	return nil, err
}

func (itp *AstInterpreter) VisitExpressionStmt(s *ExpressionStmt) (result interface{}, err error) {
	result, err = s.expression.Accept(itp)
	return result, err
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

func (itp *AstInterpreter) VisitCallExpr(e *CallExpr) (result interface{}, err error) {
	callee, err := e.callee.Accept(itp)
	if err != nil {
		return nil, err
	}
	var args []interface{}
	for _, arg := range e.arguments {
		argResult, err := arg.Accept(itp)
		if err != nil {
			return nil, err
		}

		args = append(args, argResult)
	}

	if function, ok := callee.(LoxCallable); ok {
		if function.Arity() != len(args) {
			return nil, fmt.Errorf("Expected %d arguments but got %d.", function.Arity(), len(args))
		}

		callResult, err := function.Call(itp, args)
		if err != nil {
			return nil, err
		}

		return callResult, nil
	}

	return nil, fmt.Errorf("Can only call functions and classes.")
}

func (itp *AstInterpreter) VisitLogicalExpr(e *LogicalExpr) (result interface{}, err error) {
	leftResult, err := e.left.Accept(itp)
	if err != nil {
		return nil, err
	}

	// short-circuit
	if e.operator.Type == And {
		if !isTruthy(leftResult) {
			return false, nil
		}
	} else if e.operator.Type == Or {
		if isTruthy(leftResult) {
			return leftResult, nil
		}
	} else {
		panic("Unsupported logical operator " + e.operator.Lexeme)
	}

	return e.right.Accept(itp)
}

func (itp *AstInterpreter) VisitAssignExpr(e *AssignExpr) (result interface{}, err error) {
	result, err = e.assignValue.Accept(itp)
	assignErr := itp.env.Assign(e.variableName, result)
	if assignErr != nil {
		return nil, assignErr
	}

	return result, err
}

func (itp *AstInterpreter) VisitVariableExpr(e *VariableExpr) (result interface{}, err error) {
	return itp.env.Get(e.variableName)
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

func loxStringify(value interface{}) string {
	switch val := value.(type) {
	case nil:
		return "nil"
	case float64:
		return strconv.FormatFloat(val, 'f', -1, 64)
	default:
		return fmt.Sprintf("%v", val)
	}
}
