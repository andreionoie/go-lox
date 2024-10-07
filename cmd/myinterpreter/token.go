package main

import (
	"fmt"
	"strconv"
	"strings"
)

type TokenType string

const (
	LeftParen    TokenType = "LEFT_PAREN"
	RightParen             = "RIGHT_PAREN"
	LeftBrace              = "LEFT_BRACE"
	RightBrace             = "RIGHT_BRACE"
	Comma                  = "COMMA"
	Dot                    = "DOT"
	Minus                  = "MINUS"
	Plus                   = "PLUS"
	Semicolon              = "SEMICOLON"
	Star                   = "STAR"
	Equal                  = "EQUAL"
	EqualEqual             = "EQUAL_EQUAL"
	Bang                   = "BANG"
	BangEqual              = "BANG_EQUAL"
	Less                   = "LESS"
	LessEqual              = "LESS_EQUAL"
	Greater                = "GREATER"
	GreaterEqual           = "GREATER_EQUAL"
	Slash                  = "SLASH"
	String                 = "STRING"
	Number                 = "NUMBER"
	Eof                    = "EOF"
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
		switch v := tok.Literal.(type) {
		case float64:
			// return smallest number of digits necessary
			lit = strconv.FormatFloat(v, 'f', -1, 64)
			if !strings.Contains(lit, ".") {
				lit += ".0"
			}
		case string:
			// Use the string as is
			lit = v
		default:
			// Fallback for other types
			lit = fmt.Sprintf("%v", v)
		}
	}

	return fmt.Sprintf("%s %s %s", tok.Type, tok.Lexeme, lit)
}
