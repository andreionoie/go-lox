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
	case '<':
		if s.match('=') {
			s.addToken(LessEqual)
		} else {
			s.addToken(Less)
		}
	case '>':
		if s.match('=') {
			s.addToken(GreaterEqual)
		} else {
			s.addToken(Greater)
		}
	case '/':
		if s.match('/') {
			for !s.isAtEnd() && s.Source[s.Current] != '\n' {
				s.Current++
			}
		} else {
			s.addToken(Slash)
		}
	case '"':
		s.string()
	case ' ', '\t':
		// noop
	case '\n':
		s.Line++
	default:
		s.logError("Unexpected character: %c", nextRune)
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
	s.Current++
	return s.Source[s.Current-1]
}

func (s *Scanner) match(expected rune) bool {
	if s.isAtEnd() || (s.Source[s.Current] != expected) {
		return false
	}

	s.Current++
	return true
}

func (s *Scanner) string() {
	for !s.isAtEnd() && s.Source[s.Current] != '"' {
		if s.Source[s.Current] == '\n' {
			s.Line++
		}
		s.Current++
	}

	if s.isAtEnd() {
		s.logError("Unterminated string.")
		return
	}
	s.Current++
	s.addTokenWithLiteral(String, string(s.Source[s.Start+1:s.Current-1]))
}

func (s *Scanner) logError(msg string, a ...any) {
	fmtString := fmt.Sprintf("[line %d] Error: %s\n", s.Line+1, msg)
	fmt.Fprintf(os.Stderr, fmtString, a...)
	s.HadErrors = true
}
