package main

import (
	"fmt"
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
	if p.match(Var) {
		nextStmt, err = p.varDeclaration()
	} else {
		nextStmt, err = p.statement()
	}

	if err != nil {
		p.synchronize()
	}
	return nextStmt, err
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

func (p *Parser) expression() (Expr, error) {
	return p.assignment()
}

// assignment -> IDENTIFIER "=" assignment (left assoc.)
// assignment -> equality
func (p *Parser) assignment() (Expr, error) {
	leftEquality, err := p.equality()
	if err != nil {
		return nil, err
	}

	if p.match(Equal) {
		value, err := p.assignment()
		if err != nil {
			return nil, err
		}

		if variableExpr, ok := leftEquality.(*VariableExpr); ok {
			lvalueToken := variableExpr.variableName
			return &AssignExpr{variableName: lvalueToken, assignValue: value}, nil
		}

		return nil, p.getError("Invalid assignment target.")
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
// unary -> primary
func (p *Parser) unary() (Expr, error) {
	if p.match(Minus, Bang) {
		op := p.previous()
		nestedUnary, err := p.unary()
		if err != nil {
			return nil, err
		}
		return &UnaryExpr{operator: op, right: nestedUnary}, nil
	}

	return p.primary()
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
