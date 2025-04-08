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
	"VAR":  *NewTokenType("VAR", "var"),

	// RegExp для типов данных
	"SHOW":     *NewTokenType("SHOW", "show"),
	"VARIABLE": *NewTokenType("VARIABLE", "[a-zA-Z_][a-zA-Z0-9_]*"),
	"INTEGER":  *NewTokenType("INTEGER", `\d+`),
	"DOUBLE":   *NewTokenType("DOUBLE", `\d+\.\d+`),
	// Арифметические операции
	"ASSIGN": *NewTokenType("ASSIGN", "="),
	"MINUS":  *NewTokenType("MINUS", "-"),
	"PLUS":   *NewTokenType("PLUS", "\\+"),
	"LPAREN": *NewTokenType("LPAREN", "\\("),
	"RPAREN": *NewTokenType("RPAREN", "\\)"),
	// Логические операции
	"IF":       *NewTokenType("IF", "if"),
	"EQUAL":    *NewTokenType("EQUAL", "equal"),
	"NONEQUAL": *NewTokenType("NONEQUAL", "non-equal"),
	"MORE":     *NewTokenType("MORE", "more"),
	"LESS":     *NewTokenType("LESS", "less"),

	"SEMICOLON":  *NewTokenType("SEMICOLON", ";"),
	"WHITESPACE": *NewTokenType("WHITESPACE", "[ \n\t\r]+"),
}

var TokenTypesOrdered = []TokenType{
	// Ключевые слова
	*NewTokenType("IF", "if"),
	*NewTokenType("VAR", "var"),
	*NewTokenType("INT", "int"),
	*NewTokenType("DOUB", "double"),
	*NewTokenType("SHOW", "show"),

	// Логические операторы
	*NewTokenType("EQUAL", "equal"),
	*NewTokenType("NONEQUAL", "non-equal"),
	*NewTokenType("MORE", "more"),
	*NewTokenType("LESS", "less"),

	// Литералы
	*NewTokenType("DOUBLE", `\d+\.\d+`),
	*NewTokenType("INTEGER", `\d+`),
	*NewTokenType("VARIABLE", `[a-zA-Z_][a-zA-Z0-9_]*`), // важно ставить после ключевых слов

	// Арифметические операторы
	*NewTokenType("ASSIGN", "="),
	*NewTokenType("PLUS", `\+`),
	*NewTokenType("MINUS", `-`),

	// Скобки и разделители
	*NewTokenType("LPAREN", `\(`),
	*NewTokenType("RPAREN", `\)`),
	*NewTokenType("SEMICOLON", ";"),

	// Пробелы (последним, чтобы можно было игнорировать)
	*NewTokenType("WHITESPACE", `[ \n\t\r]+`),
}
