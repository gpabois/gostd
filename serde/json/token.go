package json

type TokenType int

type Token struct {
	typ TokenType
	lit string
}

func (tok Token) OpenDocument() Token {
	return Token{
		typ: TOK_OPEN_DOCUMENT,
		lit: "{",
	}
}

func (tok Token) CloseDocument() Token {
	return Token{
		typ: TOK_CLOSE_DOCUMENT,
		lit: "}",
	}
}

func (tok Token) OpenArray() Token {
	return Token{
		typ: TOK_OPEN_ARRAY,
		lit: "[",
	}
}

func (tok Token) CloseArray() Token {
	return Token{
		typ: TOK_CLOSE_ARRAY,
		lit: "]",
	}
}

func (tok Token) String(str string) Token {
	return Token{
		typ: TOK_STRING,
		lit: str,
	}
}

func (tok Token) Colon() Token {
	return Token{
		typ: TOK_COLON,
		lit: ":",
	}
}

func (tok Token) Ws(str string) Token {
	return Token{
		typ: TOK_WS,
		lit: str,
	}
}

func (tok Token) Invalid(str string) Token {
	return Token{
		typ: TOK_INVALID,
		lit: str,
	}
}

func (tok Token) Eof() Token {
	return Token{
		typ: TOK_EOF,
		lit: "[EOF]",
	}
}

func (tok Token) Comma() Token {
	return Token{
		typ: TOK_COMMA,
		lit: ",",
	}
}

func (tok Token) Number(number string) Token {
	return Token{
		typ: TOK_NUMBER,
		lit: number,
	}
}

func (tok Token) ToString() string {
	return tok.lit
}

const (
	TOK_EOF = iota
	TOK_INVALID
	TOK_WS

	TOK_OPEN_DOCUMENT
	TOK_CLOSE_DOCUMENT
	TOK_COMMA
	TOK_COLON

	TOK_STRING
	TOK_TRUE
	TOK_FALSE
	TOK_NULL

	TOK_NUMBER

	TOK_OPEN_ARRAY
	TOK_CLOSE_ARRAY
)
