package main

import "fmt"

type TokenType string

const (
	LeftParen  TokenType = "LEFT_PAREN"
	RightParen           = "RIGHT_PAREN"
	LeftBrace            = "LEFT_BRACE"
	RightBrace           = "RIGHT_BRACE"
	Comma                = "COMMA"
	Dot                  = "DOT"
	Minus                = "MINUS"
	Plus                 = "PLUS"
	Semicolon            = "SEMICOLON"
	Star                 = "STAR"
	Eof                  = "EOF"
)

type Token struct {
	Type    TokenType
	Lexeme  string
	Literal interface{}
	Line    int
}

func (tok Token) String() string {
	lit := "null"
	if tok.Literal != nil {
		lit = fmt.Sprintf("%s", tok.Literal)
	}

	return fmt.Sprintf("%s %s %s", tok.Type, tok.Lexeme, lit)
}
