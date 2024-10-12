package main

type Parser struct {
	Tokens  []Token
	Current int
}

func (p *Parser) Parse() Expr {
	return p.expression()
}

func (p *Parser) expression() Expr {
	return p.term()
}

// term -> factor ( ("+" | "-") factor )*
func (p *Parser) term() Expr {
	leftFactor := p.factor()
	for p.match(Plus, Minus) { // this "while" loop represents the * suffix in the notation
		op := p.previous()
		rightFactor := p.factor()
		leftFactor = &BinaryExpr{left: leftFactor, operator: op, right: rightFactor}
	}
	return leftFactor
}

// factor -> unary ( ("*" | "/") unary )*
func (p *Parser) factor() Expr {
	leftUnary := p.unary()
	for p.match(Star, Slash) { // this "while" loop represents the * suffix in the notation
		op := p.previous()
		rightUnary := p.unary()
		leftUnary = &BinaryExpr{left: leftUnary, operator: op, right: rightUnary}
	}
	return leftUnary
}

// unary -> ("-" | "!") unary
// unary -> primary
func (p *Parser) unary() Expr {
	if p.match(Minus, Bang) {
		op := p.previous()
		nestedUnary := p.unary()
		return &UnaryExpr{operator: op, right: nestedUnary}
	}

	return p.primary()
}

// primary -> NUMBER | STRING | "true" | "false" | "nil"
// primary -> "(" expression ")"
func (p *Parser) primary() Expr {
	if p.match(Number, String, True, False, Nil) {
		return &LiteralExpr{value: p.previous().GetTokenAsTerminal()}
	}

	if p.match(LeftParen) {
		grouping := p.expression()
		p.Current++
		if p.previous().Type != RightParen {
			panic("TODO: implement error handling")
		}

		return &GroupingExpr{expr: grouping}
	}

	panic("TODO: implement error handling")
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
