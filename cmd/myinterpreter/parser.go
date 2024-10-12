package main

type Parser struct {
	Tokens  []Token
	Current int
}

func (p *Parser) Parse() Expr {
	return p.expression()
}

func (p *Parser) expression() Expr {
	return p.primary()
}

// primary -> NUMBER | STRING | "true" | "false" | "nil"
// primary -> "(" expression ")"
func (p *Parser) primary() Expr {
	if p.match(True) {
		return &LiteralExpr{value: "true"}
	}
	if p.match(False) {
		return &LiteralExpr{value: "false"}
	}
	if p.match(Nil) {
		return &LiteralExpr{value: "nil"}
	}
	if p.match(Number, String) {
		return &LiteralExpr{value: p.previous().GetLiteralAsString()}
	}

	if p.match(LeftParen) {
		grouping := p.expression()
		p.Current++
		if p.previous().Type != RightParen {
			panic("Unmatched parens")
		}

		return &GroupingExpr{expr: grouping}
	}

	panic("unreachable")
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
