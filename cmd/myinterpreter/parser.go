package main

import (
	"fmt"
	"os"
)

type Parser struct {
	Tokens  []Token
	Current int
}

func (p *Parser) Parse() ([]Stmt, error) {
	var statements []Stmt
	for p.Tokens[p.Current].Type != Eof {
		nextStmt, err := p.declaration()
		if err != nil {
			return nil, err
		}
		statements = append(statements, nextStmt)
	}
	return statements, nil
}

func (p *Parser) declaration() (nextStmt Stmt, err error) {
	if p.match(Function) {
		nextStmt, err = p.funcDeclaration()
	} else if p.match(Var) {
		nextStmt, err = p.varDeclaration()
	} else {
		nextStmt, err = p.statement()
	}

	if err != nil {
		p.synchronize()
	}
	return nextStmt, err
}

// funDecl -> "fun" IDENTIFIER "(" parameters? ")" block
func (p *Parser) funcDeclaration() (Stmt, error) {
	p.Current++
	if p.previous().Type != Identifier {
		return nil, p.getError("Expect function name.")
	}
	funcName := p.previous()

	p.Current++
	if p.previous().Type != LeftParen {
		return nil, p.getError("Expect '(' after function name.")
	}

	var params []Token
	if p.Tokens[p.Current].Type != RightParen {
		p.Current++
		if p.previous().Type != Identifier {
			return nil, p.getError("Expect parameter name.")
		}

		params = []Token{p.previous()}
		for p.match(Comma) {
			p.Current++
			if p.previous().Type != Identifier {
				return nil, p.getError("Expect parameter name.")
			}
			params = append(params, p.previous())

			if len(params) > 255 {
				fmt.Fprintln(os.Stderr, p.getError("Can't have more than 255 parameters."))
			}
		}
	}

	p.Current++
	if p.previous().Type != RightParen {
		return nil, p.getError("Expect ')' after parameters.")
	}

	p.Current++
	if p.previous().Type != LeftBrace {
		return nil, p.getError("Expect '{' before function body.")
	}

	funcBody, err := p.block()
	if err != nil {
		return nil, err
	}

	return &FunctionStmt{
		name:       funcName,
		parameters: params,
		body:       funcBody,
	}, nil
}

func (p *Parser) varDeclaration() (Stmt, error) {
	p.Current++
	if p.previous().Type != Identifier {
		return nil, p.getError("Expect variable name.")
	}
	variableName := p.previous()
	var initializer Expr
	// optional initializer expr after '=' sign
	if p.match(Equal) {
		var err error
		initializer, err = p.expression()
		if err != nil {
			return nil, err
		}
	}
	p.Current++
	if p.previous().Type != Semicolon {
		return nil, p.getError("Expect ';' after variable declaration.")
	}

	return &VarStmt{varName: variableName, initializerExpression: initializer}, nil
}

func (p *Parser) statement() (Stmt, error) {
	if p.match(Print) {
		return p.printStatement()
	}
	if p.match(LeftBrace) {
		stmts, err := p.block()
		return &BlockStmt{statements: stmts}, err
	}
	if p.match(If) {
		return p.ifStatement()
	}
	if p.match(While) {
		return p.whileStatement()
	}
	if p.match(For) {
		return p.forStatement()
	}
	if p.match(Return) {
		return p.returnStatement()
	}

	return p.expressionStatement()
}

func (p *Parser) block() ([]Stmt, error) {
	var stmts []Stmt

	for (p.Current < len(p.Tokens)) && (p.Tokens[p.Current].Type != Eof) && (p.Tokens[p.Current].Type) != RightBrace {
		stmt, err := p.declaration()
		if err != nil {
			return nil, err
		}
		stmts = append(stmts, stmt)
	}

	p.Current++
	if p.previous().Type != RightBrace {
		return nil, p.getError("Expect '}' after block.")
	}
	return stmts, nil
}

func (p *Parser) returnStatement() (s Stmt, err error) {
	kyw := p.previous()
	var expr Expr
	if p.Tokens[p.Current].Type != Semicolon {
		expr, err = p.expression()
		if err != nil {
			return nil, err
		}
	}
	p.Current++
	if p.previous().Type != Semicolon {
		return nil, p.getError("Expect ';' after return value.")
	}

	return &ReturnStmt{
		keyword: kyw,
		value:   expr,
	}, nil
}

// TODO: refactor, find better implementation
// ForStmt -> "for" "(" (VarDecl | Expr ";")?  Expr? ";" Expr? ")" statement
func (p *Parser) forStatement() (Stmt, error) {
	// TODO: fix bug where errors get set when trying to parse empty for header element
	hadNoErrors := !LoxHadError
	// initialization is a variable declaration or expression statement
	// consume '('
	p.Current++
	if p.previous().Type != LeftParen {
		return nil, p.getError("Expect '(' after 'for'.")
	}

	// read initialization expression
	var initialization Stmt
	if p.match(Semicolon) {
		initialization = nil
	} else {
		if p.match(Var) {
			initialization, _ = p.varDeclaration()
		} else {
			initialization, _ = p.expressionStatement()
		}
	}

	// read condition expression
	var condition Expr
	if p.match(Semicolon) {
		condition = nil
	} else {
		condition, _ = p.expression()
		p.Current++
		if p.previous().Type != Semicolon {
			return nil, p.getError("Expect ';' after for condition.")
		}
	}

	var iteration Expr
	if p.match(RightParen) {
		iteration = nil
	} else {
		iteration, _ = p.expression()
		if iteration == nil {
			// no expression found; advance
			p.Current++
		}
		// consume ')'
		p.Current++
		if p.previous().Type != RightParen {
			return nil, p.getError("Expect ')' after for header.")
		}
	}

	if hadNoErrors {
		LoxHadError = false
	}

	body, err := p.statement()
	if err != nil {
		return nil, err
	}

	return &ForStmt{
		init:      initialization,
		condition: condition,
		iteration: iteration,
		loopBody:  body,
	}, nil
}

// WhileStmt -> "while" "(" Expr ")" statement
func (p *Parser) whileStatement() (Stmt, error) {
	p.Current++
	if p.previous().Type != LeftParen {
		return nil, p.getError("Expect '(' after 'while'.")
	}

	cond, err := p.expression()
	if err != nil {
		return nil, err
	}

	p.Current++
	if p.previous().Type != RightParen {
		return nil, p.getError("Expect ')' after while condition.")
	}

	body, err := p.statement()
	if err != nil {
		return nil, err
	}

	return &WhileStmt{
		condition: cond,
		loopBody:  body,
	}, nil
}

// IfStmt -> "if" "(" Expr ")" statement ("else" statement)?
func (p *Parser) ifStatement() (Stmt, error) {
	p.Current++
	if p.previous().Type != LeftParen {
		return nil, p.getError("Expect '(' after 'if'.")
	}

	cond, err := p.expression()
	if err != nil {
		return nil, err
	}

	p.Current++
	if p.previous().Type != RightParen {
		return nil, p.getError("Expect ')' after if condition.")
	}

	thenStmt, err := p.statement()
	if err != nil {
		return nil, err
	}

	var elseStmt Stmt
	if p.match(Else) {
		elseStmt, err = p.statement()
		if err != nil {
			return nil, err
		}
	}

	return &IfStmt{
		condition:  cond,
		thenBranch: thenStmt,
		elseBranch: elseStmt,
	}, nil
}

// PrintStmt -> "print" Expr ";"
func (p *Parser) printStatement() (Stmt, error) {
	value, err := p.expression()
	if err != nil {
		return nil, err
	}
	p.Current++
	if p.previous().Type != Semicolon {
		return nil, p.getError("Expect ';' after value.")
	}
	return &PrintStmt{value}, nil
}

// ExprStmt -> Expr ";"
func (p *Parser) expressionStatement() (Stmt, error) {
	value, err := p.expression()
	if err != nil {
		return nil, err
	}
	p.Current++
	if p.previous().Type != Semicolon {
		return nil, p.getError("Expect ';' after expression.")
	}
	return &ExpressionStmt{value}, nil
}

func (p *Parser) ParseExpr() (Expr, error) {
	return p.expression()
}

// expression -> assignment
func (p *Parser) expression() (Expr, error) {
	return p.assignment()
}

// assignment -> IDENTIFIER "=" assignment (left assoc.)
// assignment -> logicalOr
func (p *Parser) assignment() (Expr, error) {
	lvalue, err := p.logicalOr()
	if err != nil {
		return nil, err
	}

	if p.match(Equal) {
		value, err := p.assignment()
		if err != nil {
			return nil, err
		}

		if variableExpr, ok := lvalue.(*VariableExpr); ok {
			lvalueToken := variableExpr.variableName
			return &AssignExpr{
				variableName: lvalueToken,
				assignValue:  value,
			}, nil
		}

		return nil, p.getError("Invalid assignment target.")
	}

	return lvalue, nil
}

// logicalOr -> logicalAnd ("or" logicalAnd)*
func (p *Parser) logicalOr() (Expr, error) {
	leftAnd, err := p.logicalAnd()
	if err != nil {
		return nil, err
	}

	for p.match(Or) {
		orOperator := p.previous()
		rightAnd, err := p.logicalAnd()
		if err != nil {
			return nil, err
		}

		leftAnd = &LogicalExpr{
			left:     leftAnd,
			operator: orOperator,
			right:    rightAnd,
		}
	}

	return leftAnd, nil
}

// logicalAnd -> equality ("and" equality)*
func (p *Parser) logicalAnd() (Expr, error) {
	leftEquality, err := p.equality()
	if err != nil {
		return nil, err
	}

	for p.match(And) {
		andOperator := p.previous()
		rightEquality, err := p.equality()
		if err != nil {
			return nil, err
		}

		leftEquality = &LogicalExpr{
			left:     leftEquality,
			operator: andOperator,
			right:    rightEquality,
		}
	}

	return leftEquality, nil
}

// equality -> comparison ( ("==" | "!=") comparison)*
func (p *Parser) equality() (Expr, error) {
	leftCmp, err := p.comparison()
	if err != nil {
		return nil, err
	}
	for p.match(EqualEqual, BangEqual) {
		op := p.previous()
		rightCmp, err := p.comparison()
		if err != nil {
			return nil, err
		}
		leftCmp = &BinaryExpr{left: leftCmp, operator: op, right: rightCmp}
	}
	return leftCmp, nil
}

// comparison -> term ( (>|<|<=|>=) term )*
func (p *Parser) comparison() (Expr, error) {
	leftTerm, err := p.term()
	if err != nil {
		return nil, err
	}
	for p.match(Greater, Less, GreaterEqual, LessEqual) {
		op := p.previous()
		rightTerm, err := p.term()
		if err != nil {
			return nil, err
		}
		leftTerm = &BinaryExpr{left: leftTerm, operator: op, right: rightTerm}
	}
	return leftTerm, nil
}

// term -> factor ( ("+" | "-") factor )*
func (p *Parser) term() (Expr, error) {
	leftFactor, err := p.factor()
	if err != nil {
		return nil, err
	}
	for p.match(Plus, Minus) { // this "while" loop represents the * suffix in the notation
		op := p.previous()
		rightFactor, err := p.factor()
		if err != nil {
			return nil, err
		}
		leftFactor = &BinaryExpr{left: leftFactor, operator: op, right: rightFactor}
	}
	return leftFactor, nil
}

// factor -> unary ( ("*" | "/") unary )*
func (p *Parser) factor() (Expr, error) {
	leftUnary, err := p.unary()
	if err != nil {
		return nil, err
	}
	for p.match(Star, Slash) { // this "while" loop represents the * suffix in the notation
		op := p.previous()
		rightUnary, err := p.unary()
		if err != nil {
			return nil, err
		}
		leftUnary = &BinaryExpr{left: leftUnary, operator: op, right: rightUnary}
	}
	return leftUnary, nil
}

// unary -> ("-" | "!") unary
// unary -> call
func (p *Parser) unary() (Expr, error) {
	if p.match(Minus, Bang) {
		op := p.previous()
		nestedUnary, err := p.unary()
		if err != nil {
			return nil, err
		}
		return &UnaryExpr{operator: op, right: nestedUnary}, nil
	}

	return p.call()
}

// call -> primary ( "(" arguments? ")" )*
func (p *Parser) call() (Expr, error) {
	primaryExpr, err := p.primary()
	if err != nil {
		return nil, err
	}

	for p.match(LeftParen) {
		primaryExpr, err = p.finishCall(primaryExpr)
		if err != nil {
			return nil, err
		}
	}

	return primaryExpr, nil
}

func (p *Parser) finishCall(callee Expr) (Expr, error) {
	var args []Expr
	if p.Tokens[p.Current].Type != RightParen {
		firstArg, err := p.expression()
		if err != nil {
			return nil, err
		}
		args = []Expr{firstArg}

		for p.match(Comma) {
			nextArg, err := p.expression()
			if err != nil {
				return nil, err
			}
			args = append(args, nextArg)

			if len(args) > 255 {
				fmt.Fprintln(os.Stderr, p.getError("Can't have more than 255 arguments."))
			}
		}
	}

	p.Current++
	if p.previous().Type != RightParen {
		return nil, p.getError("Expect ')' after arguments.")
	}

	return &CallExpr{
		callee:       callee,
		arguments:    args,
		closingParen: p.previous(),
	}, nil
}

// primary -> NUMBER | STRING | "true" | "false" | "nil"
// primary -> "(" expression ")"
// primary -> IDENTIFIER (variable)
func (p *Parser) primary() (Expr, error) {
	if p.match(True) {
		return &LiteralExpr{value: true}, nil
	}
	if p.match(False) {
		return &LiteralExpr{value: false}, nil
	}
	if p.match(Nil) {
		return &LiteralExpr{value: nil}, nil
	}
	if p.match(Number, String) {
		return &LiteralExpr{value: p.previous().Literal}, nil
	}

	if p.match(LeftParen) {
		grouping, err := p.expression()
		if err != nil {
			return nil, err
		}
		p.Current++
		if p.previous().Type != RightParen {
			return nil, p.getError("Expect ')' after expression.")
		}

		return &GroupingExpr{expr: grouping}, nil
	}

	if p.match(Identifier) {
		return &VariableExpr{variableName: p.previous()}, nil
	}

	return nil, p.getError("Expect expression.")
}

func (p *Parser) synchronize() {
	// discard everything until finding a new statement boundary
	p.Current++
	for p.Current < len(p.Tokens) && p.Tokens[p.Current].Type != Eof {
		if p.previous().Type == Semicolon {
			return
		}

		switch p.Tokens[p.Current].Type {
		case Class, Function, Var, For, If, While, Print, Return:
			return
		}

		p.Current++
	}
}

func (p *Parser) previous() Token {
	return p.Tokens[p.Current-1]
}

func (p *Parser) match(tokenTypes ...TokenType) bool {
	if p.Tokens[p.Current].Type == Eof {
		return false
	}

	for _, tokenType := range tokenTypes {
		if p.Tokens[p.Current].Type == tokenType {
			p.Current++
			return true
		}
	}

	return false
}

func (p *Parser) getError(msg string, a ...any) error {
	var currentToken string
	var line int
	if p.Current >= len(p.Tokens) || p.Tokens[p.Current].Type == Eof {
		currentToken = "end"
		line = p.Tokens[len(p.Tokens)-1].Line
	} else {
		currentToken = "'" + p.Tokens[p.Current].Lexeme + "'"
		line = p.Tokens[p.Current].Line
	}

	LoxHadError = true
	return fmt.Errorf("[line %d] Error at %s: %s\n", line+1, currentToken, fmt.Sprintf(msg, a...))
}
