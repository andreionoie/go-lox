package main

import (
	"fmt"
)

type Parser struct {
	Tokens  []Token
	Current int
}

func (p *Parser) Parse() (Expr, error) {
	return p.expression()
}

func (p *Parser) expression() (Expr, error) {
	return p.equality()
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

	return nil, p.getError("Expect expression.")
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
	if p.Tokens[p.Current].Type == Eof {
		currentToken = "end"
	} else {
		currentToken = "'" + p.Tokens[p.Current].Lexeme + "'"
	}

	LoxHadError = true
	return fmt.Errorf("[line %d] Error at %s: %s\n", p.Tokens[p.Current].Line+1, currentToken, fmt.Sprintf(msg, a...))
}
