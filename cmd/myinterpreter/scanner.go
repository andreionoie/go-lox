package main

import (
	"fmt"
	"os"
)

type Scanner struct {
	Source    []rune
	Tokens    []Token
	Start     int
	Current   int
	Line      int
	HadErrors bool
}

func (s *Scanner) scanTokens() {
	for !s.isAtEnd() {
		s.Start = s.Current
		s.scan()
	}
	s.Tokens = append(s.Tokens, Token{
		Lexeme:  "",
		Literal: nil,
		Line:    s.Line,
		Type:    Eof,
	})
}

func (s *Scanner) isAtEnd() bool {
	return s.Current >= len(s.Source)
}

func (s *Scanner) scan() {
	nextRune := s.advance()
	switch nextRune {
	case '(':
		s.addToken(LeftParen)
	case ')':
		s.addToken(RightParen)
	case '{':
		s.addToken(LeftBrace)
	case '}':
		s.addToken(RightBrace)
	case ',':
		s.addToken(Comma)
	case '.':
		s.addToken(Dot)
	case '-':
		s.addToken(Minus)
	case '+':
		s.addToken(Plus)
	case ';':
		s.addToken(Semicolon)
	case '*':
		s.addToken(Star)
	case '=':
		if s.match('=') {
			s.addToken(EqualEqual)
		} else {
			s.addToken(Equal)
		}
	case '!':
		if s.match('=') {
			s.addToken(BangEqual)
		} else {
			s.addToken(Bang)
		}
	case '\n':
		s.Line++
	default:
		fmt.Fprintf(os.Stderr, "[line %d] Error: Unexpected character: %s\n", s.Line+1, string(nextRune))
		s.HadErrors = true
	}
}

func (s *Scanner) addTokenWithLiteral(tokenType TokenType, literal interface{}) {
	token := Token{
		Lexeme:  string(s.Source[s.Start:s.Current]),
		Literal: literal,
		Line:    s.Line,
		Type:    tokenType,
	}
	s.Tokens = append(s.Tokens, token)
}

func (s *Scanner) addToken(tokenType TokenType) {
	s.addTokenWithLiteral(tokenType, nil)
}

func (s *Scanner) advance() rune {
	curr := s.Current
	s.Current++
	return s.Source[curr]
}

func (s *Scanner) match(expected rune) bool {
	if s.isAtEnd() || (s.Source[s.Current] != expected) {
		return false
	}

	s.Current++
	return true
}
