package imageutil

import (
	"strconv"
	"unicode"
)

type Scanner struct {
	input        []byte
	position     int  // current position in input
	readPosition int  // current reading position
	ch           byte // current character
}

func NewScanner(input []byte) *Scanner {
	s := &Scanner{input: input}
	s.readChar() // Initialize first Character
	return s
}

func (s *Scanner) readChar() {
	if s.readPosition >= len(s.input) {
		s.ch = 0
	} else {
		s.ch = s.input[s.readPosition]
	}
	s.position = s.readPosition
	s.readPosition++
}

func (s *Scanner) skipWhitespace() {
	for s.ch == ' ' || s.ch == '\n' || s.ch == '\r' || s.ch == '\t' {
		s.readChar()
	}
}

func (s *Scanner) NextNumber() uint8 {
	s.skipWhitespace()
	b := s.readNumber()
	num, err := strconv.Atoi(string(b))
	if err != nil {
		panic("error: " + err.Error())
	}
	return uint8(num)
}

func (s *Scanner) readNumber() []byte {
	position := s.position
	for isDigit(s.ch) {
		s.readChar()
	}
	return s.input[position:s.position]
}

func isDigit(b byte) bool {
	return unicode.IsDigit(rune(b))
}
