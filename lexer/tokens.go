package lexer

type Token struct {
	TypeToken TokenType
	Text      string
	Pos       int
}

func NewToken(typeToken TokenType, text string, pos int) *Token {
	return &Token{TypeToken: typeToken, Text: text, Pos: pos}
}

type TokenType struct {
	Name  string
	Regex string
}

func NewTokenType(name string, regex string) *TokenType {
	return &TokenType{Name: name, Regex: regex}
}

var TokenTypeList = map[string]TokenType{
	"NUMBER":    *NewTokenType("NUMBER", "[0-9]*"),
	"VARIABLE":  *NewTokenType("VARIABLE", "[a-zA-Z]+"),
	"SEMICOLON": *NewTokenType("SEMICOLON", ";"),
	"SPACE":     *NewTokenType("SPACE", " \\n\\t\\r"),
	"ASSIGN":    *NewTokenType("ASSIGN", "->"),
	"LOG":       *NewTokenType("LOG", "show"),
	"MINUS":     *NewTokenType("MINUS", "-"),
	"PLUS":      *NewTokenType("PLUS", "+"),
	"LPAREN":    *NewTokenType("LPAREN", "\\("),
	"RPAREN":    *NewTokenType("RPAREN", "\\)"),
	"EQUAL":     *NewTokenType("EQUAL", "equal"),
}
