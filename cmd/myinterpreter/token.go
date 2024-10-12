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

	String     = "STRING"
	Number     = "NUMBER"
	Identifier = "IDENTIFIER"

	True  = "TRUE"
	False = "FALSE"
	And   = "AND"
	Or    = "OR"

	Nil   = "NIL"
	This  = "THIS"
	Super = "SUPER"

	Function = "FUN"
	Class    = "CLASS"

	If     = "IF"
	Else   = "ELSE"
	For    = "FOR"
	While  = "WHILE"
	Return = "RETURN"

	Print = "PRINT"
	Var   = "VAR"

	Eof = "EOF"
)

var ReservedKeywords = map[string]TokenType{
	"true":  True,
	"false": False,
	"and":   And,
	"or":    Or,

	"nil":   Nil,
	"this":  This,
	"super": Super,

	"fun":   Function,
	"class": Class,

	"if":     If,
	"else":   Else,
	"for":    For,
	"while":  While,
	"return": Return,

	"print": Print,
	"var":   Var,
}

type Token struct {
	Type    TokenType
	Lexeme  string
	Literal interface{}
	Line    int
}

func (tok Token) String() string {
	return fmt.Sprintf("%s %s %s", tok.Type, tok.Lexeme, tok.GetLiteralAsString())
}

func (tok Token) GetTokenAsTerminal() string {
	switch tok.Type {
	case Number, String:
		return tok.GetLiteralAsString()
	case True, False, Nil:
		return strings.ToLower(string(tok.Type))
	default:
		panic("Unsupported terminal")
	}
}

func (tok Token) GetLiteralAsString() string {
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
	return lit
}
