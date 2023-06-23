package json

import (
	"bufio"
	"bytes"
	"io"

	"github.com/gpabois/gostd/ops"
	"github.com/gpabois/gostd/option"
)

const eof = rune(0)

type Scanner struct {
	r *bufio.Reader
}

func NewScanner(r io.Reader) *Scanner {
	return &Scanner{r: bufio.NewReader(r)}
}

func (s *Scanner) read() rune {
	ch, _, err := s.r.ReadRune()
	if err != nil {
		return eof
	}
	return ch
}

func (s *Scanner) rewind() {
	_ = s.r.UnreadRune()
}

func (s *Scanner) scanIdent() Token {
	var buf bytes.Buffer

	for {
		ch := s.read()

		if !isLetter(ch) {
			s.rewind()
			break
		}

		buf.WriteRune(ch)
	}

	switch buf.String() {
	case "true":
		return Token{
			typ: TOK_TRUE,
			lit: buf.String(),
		}
	case "false":
		return Token{
			typ: TOK_FALSE,
			lit: buf.String(),
		}
	case "null":
		return Token{
			typ: TOK_NULL,
			lit: buf.String(),
		}
	default:
		return Token{
			typ: TOK_INVALID,
			lit: buf.String(),
		}
	}
}

func (s *Scanner) scanNumber() Token {
	var buf bytes.Buffer
	isFraction := false

	for {
		ch := s.read()
		if isDigit(ch) {
			buf.WriteRune(ch)
		} else if ch == '.' && !isFraction {
			isFraction = true
			buf.WriteRune(ch)
		} else {
			s.rewind()
			return Token{}.Number(buf.String())
		}
	}
}

func (s *Scanner) scanWhiteSpaces() Token {
	var buf bytes.Buffer
	for {
		ch := s.read()

		if !isWhiteSpace(ch) {
			s.rewind()

			return Token{
				typ: TOK_WS,
				lit: buf.String(),
			}
		}

		buf.WriteRune(ch)
	}
}
func (s *Scanner) scanString() Token {
	var buf bytes.Buffer
	prev := rune(0)
	for {
		ch := s.read()
		// We escaped the " and ' character
		if (ch == '"' || isEscape(ch)) && isEscape(prev) {
			buf.WriteRune(ch)
		} else if ch == '"' || ch == eof {
			return Token{}.String(buf.String())
		} else {
			buf.WriteRune(ch)
		}
		prev = ch
	}

}

func (s *Scanner) Next() option.Option[Token] {
	tok := s.Scan()

	if tok.typ == TOK_EOF {
		return option.None[Token]()
	}

	return option.Some(tok)
}

func (s *Scanner) Scan() Token {
	tok := s.scan()
	for ; tok.typ == TOK_WS; tok = s.scan() {
	}

	return tok
}

func (s *Scanner) scan() Token {
	// Read character
	ch := s.read()

	if ch == '"' {
		return s.scanString()
	} else if isWhiteSpace(ch) {
		return s.scanWhiteSpaces()
	} else if ch == ':' {
		return Token{}.Colon()
	} else if ch == '{' {
		return Token{}.OpenDocument()
	} else if ch == '}' {
		return Token{}.CloseDocument()
	} else if ch == '[' {
		return Token{}.OpenArray()
	} else if ch == ']' {
		return Token{}.CloseArray()
	} else if ch == ',' {
		return Token{}.Comma()
	} else if isDigit(ch) {
		s.rewind()
		return s.scanNumber()
	} else if isLetter(ch) {
		s.rewind()
		return s.scanIdent()
	} else if ch == eof {
		return Token{}.Eof()
	} else {
		return Token{}.Invalid(string(ch))
	}
}

func isEscape(ch rune) bool {
	return ch == '\\'
}

func isLetter(ch rune) bool {
	return ops.Within(ch, 'a', 'z') || ops.Within(ch, 'A', 'Z')
}

func isWhiteSpace(ch rune) bool {
	return ch == ' ' || ch == '\t' || ch == '\n'
}

func isDigit(ch rune) bool {
	return ops.Within(ch, '0', '9')
}
