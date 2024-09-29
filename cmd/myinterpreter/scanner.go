package main

import "fmt"

type Scanner struct {
	Source  []rune
	Tokens  []Token
	Start   int
	Current int
	Line    int
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
	case '\n':
		s.Line++
	default:
		fmt.Print("unrecognized character " + string(nextRune))
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
