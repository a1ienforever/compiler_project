package lexer

import (
	"fmt"
	"regexp"
)

type Lexer struct {
	code   string
	pos    int
	Tokens []Token
}

func NewLexer(code string) *Lexer {
	return &Lexer{code: code}
}

func (l *Lexer) LexerAnalysis() *[]Token {
	for l.nextToken() {
	}

	return &l.Tokens
}

func (l *Lexer) nextToken() bool {
	if l.pos >= len(l.code) {
		return false
	}

	for _, tokenType := range TokenTypesOrdered {
		regex, err := regexp.Compile("^" + tokenType.Regex)
		if err != nil {
			panic(err)
		}
		str := l.code[l.pos:]
		match := regex.FindString(str)
		if match != "" {
			token := NewToken(tokenType, match, l.pos)
			if token.TypeToken.Name != "WHITESPACE" {
				l.Tokens = append(l.Tokens, *token)
			}
			l.pos += len(match)
			return true
		}
	}

	panic(fmt.Sprintf("Ошибка компиляции кода на позиции %d", l.pos)) // Пройтись дебагером
}
