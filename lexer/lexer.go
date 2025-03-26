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
	tokenList := *TokenTypeList

	for _, value := range tokenList {
		regex, err := regexp.Compile("^" + value.Regex)
		if err != nil {
			panic(err)
		}
		str := l.code[l.pos:]
		result := regex.FindString(str)
		if result != "" {
			token := NewToken(value, result, l.pos)
			if token.TypeToken != tokenList["SPACE"] {
				l.Tokens = append(l.Tokens, *token)
			}
			l.pos += len(result)
			return true
		}
	}
	panic(fmt.Sprintf("Ошибка компиляции кода на позиции %d", l.pos))
}
