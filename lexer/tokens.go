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
	"INT":     *NewTokenType("int", "int"),
	"DOUB":    *NewTokenType("double", "double"),
	"VAR":     *NewTokenType("VAR", "var"),
	"STR":     *NewTokenType("string", "string"),
	"BOOLEAN": *NewTokenType("boolean", "boolean"),
	"TRUE":    *NewTokenType("TRUE", "true"),
	"FALSE":   *NewTokenType("FALSE", "false"),

	// RegExp для типов данных
	"SHOW":     *NewTokenType("show", "show"),
	"VARIABLE": *NewTokenType("VARIABLE", "[a-zA-Z_][a-zA-Z0-9_]*"),
	"INTEGER":  *NewTokenType("INTEGER", `\d+`),
	"DOUBLE":   *NewTokenType("DOUBLE", `\d+\.\d+`),
	"STRING":   *NewTokenType("STRING", `'[^']*'`),

	// Арифметические операции
	"ASSIGN":   *NewTokenType("ASSIGN", "="),
	"MINUS":    *NewTokenType("MINUS", "-"),
	"PLUS":     *NewTokenType("PLUS", "\\+"),
	"MULTIPLY": *NewTokenType("MULTIPLY", "\\*"),
	"DIVIDE":   *NewTokenType("DIVIDE", "/"),
	"LPAREN":   *NewTokenType("LPAREN", "\\("),
	"RPAREN":   *NewTokenType("RPAREN", "\\)"),
	"LBRACE":   *NewTokenType("LBRACE", "{"),
	"RBRACE":   *NewTokenType("RBRACE", "}"),
	// Логические операции
	"IF":       *NewTokenType("if", "if"),
	"ELSE":     *NewTokenType("else", "else"),
	"EQUAL":    *NewTokenType("EQUAL", "equal"),
	"WHILE":    *NewTokenType("while", "while"),
	"NONEQUAL": *NewTokenType("NONEQUAL", "non-equal"),
	"MORE":     *NewTokenType("MORE", "more"),
	"LESS":     *NewTokenType("LESS", "less"),
	"AND":      *NewTokenType("AND", "and"),
	"OR":       *NewTokenType("OR", "or"),

	"SEMICOLON":  *NewTokenType("SEMICOLON", ";"),
	"WHITESPACE": *NewTokenType("WHITESPACE", "[ \n\t\r]+"),
}

var TokenTypesOrdered = []TokenType{
	// Ключевые слова
	*NewTokenType("if", "if"),
	*NewTokenType("else", "else"),
	*NewTokenType("while", "while"),
	//*NewTokenType("VAR", "var"),
	*NewTokenType("int", "int"),
	*NewTokenType("double", "double"),
	*NewTokenType("show", "show"),
	*NewTokenType("string", "string"),
	*NewTokenType("boolean", "boolean"),
	*NewTokenType("TRUE", "true"),
	*NewTokenType("FALSE", "false"),

	// Логические операторы
	*NewTokenType("EQUAL", "equal"),
	*NewTokenType("NONEQUAL", "non-equal"),
	*NewTokenType("MORE", "more"),
	*NewTokenType("LESS", "less"),
	*NewTokenType("AND", "and"),
	*NewTokenType("OR", "or"),

	// Литералы
	*NewTokenType("VARIABLE", `[a-zA-Z_][a-zA-Z0-9_]*`), // важно ставить после ключевых слов
	*NewTokenType("DOUBLE", `\d+\.\d+`),
	*NewTokenType("STRING", "'[^']*'"),
	*NewTokenType("INTEGER", `\d+`),

	// Арифметические операторы
	*NewTokenType("ASSIGN", "="),
	*NewTokenType("PLUS", `\+`),
	*NewTokenType("MINUS", "-"),
	*NewTokenType("MULTIPLY", "\\*"),
	*NewTokenType("DIVIDE", "/"),

	// Скобки и разделители
	*NewTokenType("LPAREN", `\(`),
	*NewTokenType("RPAREN", `\)`),
	*NewTokenType("LBRACE", "{"),
	*NewTokenType("RBRACE", "}"),
	*NewTokenType("SEMICOLON", ";"),

	// Пробелы (последним, чтобы можно было игнорировать)
	*NewTokenType("WHITESPACE", `[ \n\t\r]+`),
}
