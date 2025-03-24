package lexer

import (
	"fmt"
	"regexp"
)

type Lexer struct {
	code   string
	pos    int
	tokens []Token
}

func NewLexer(code string) *Lexer {
	return &Lexer{code: code}
}

func (l *Lexer) LexerAnalysis() *[]Token {
	fmt.Println("Lexer Analysis")
	fmt.Println(l.code)
	for l.nextToken() {
		fmt.Println(*TokenTypeList)
	}
	tokenList := *TokenTypeList

	for i := 0; i < len(tokenList); i++ {

	}
	return &l.tokens
}

func (l *Lexer) nextToken() bool {
	if l.pos >= len(l.code) {

		return false
	}
	tokenList := *TokenTypeList

	for _, value := range tokenList {
		regex, err := regexp.Compile(`^` + value.Regex)
		if err == nil {
			result := regex.MatchString(l.code[l.pos:])
			fmt.Println(result, l.code, value)
		}
	}

	return true

}
