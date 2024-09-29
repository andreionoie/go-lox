package main

import "fmt"

type TokenType string

const (
	LeftParen  TokenType = "LEFT_PAREN"
	RightParen           = "RIGHT_PAREN"
	Eof                  = "EOF"
)

type Token struct {
	Type    TokenType
	Lexeme  string
	Literal interface{}
	Line    int
}

func (tok Token) String() string {
	return fmt.Sprintf("%s %s %s", tok.Type, tok.Lexeme, tok.Literal)
}
