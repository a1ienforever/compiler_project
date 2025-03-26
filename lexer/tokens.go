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

var TokenTypeList = &map[string]TokenType{
	// Типы данных
	"INT":  *NewTokenType("int", "int"),
	"DOUB": *NewTokenType("double", "double"),
	"VAR":  *NewTokenType("var", "var"),

	// RegExp для типов данных
	"VARIABLE": *NewTokenType("variable", "[a-zA-Z]*"),
	"INTEGER":  *NewTokenType("INTEGER", `\d+`),
	"DOUBLE":   *NewTokenType("DOUBLE", `\d+\.\d+`),
	// Арифметические операции
	"ASSIGN": *NewTokenType("ASSIGN", "="),
	"MINUS":  *NewTokenType("MINUS", "-"),
	"PLUS":   *NewTokenType("PLUS", "+"),
	"LPAREN": *NewTokenType("LPAREN", "\\("),
	"RPAREN": *NewTokenType("RPAREN", "\\)"),
	// Логические операции
	"EQUAL":    *NewTokenType("EQUAL", "equal"),
	"NONEQUAL": *NewTokenType("NONEQUAL", "non-equal"),
	"MORE":     *NewTokenType("MORE", "more"),
	"LESS":     *NewTokenType("LESS", "less"),

	"SEMICOLON": *NewTokenType("SEMICOLON", ";"),
	"SPACE":     *NewTokenType("SPACE", "[ \\n\\t\\r]"),
	"LOG":       *NewTokenType("LOG", "show"),
}
