package bson

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"errors"
	"io"
	"math"
	"regexp"
	"time"

	"github.com/gpabois/gostd/result"
)

type Parser struct {
	r *bufio.Reader
}

func NewParser(r io.Reader) *Parser {
	return &Parser{r: bufio.NewReader(r)}
}

func (p *Parser) Parse() result.Result[Bson] {
	return p.parseRoot()
}

func (p *Parser) read() byte {
	b, _ := p.r.ReadByte()
	return b
}

func (p *Parser) rewind() {
	_ = p.r.UnreadByte()
}

// Parse root
func (p *Parser) parseRoot() result.Result[Bson] {
	res := p.parseDocument()
	if res.HasFailed() {
		return result.Result[Bson]{}.Failed(res.UnwrapError())
	}
	return result.Success(Bson{
		Document: res.Expect(),
	})
}

// Parse document
func (p *Parser) parseDocument() result.Result[Document] {
	_, err := p.r.ReadBytes(4)

	if err != nil {
		return result.Result[Document]{}.Failed(err)
	}

	var elements []Element

	for {
		b := p.read()
		if b == 0x00 {
			return result.Success(Document{
				Elements: elements,
			})
		}

		p.rewind()
		res := p.parseElement()
		if res.HasFailed() {
			return result.Result[Document]{}.Failed(res.UnwrapError())
		}
		elements = append(elements, res.Expect())
	}
}

// Parse array
func (p *Parser) parseArray() result.Result[Array] {
	_, err := p.r.ReadBytes(4)

	if err != nil {
		return result.Result[Array]{}.Failed(err)
	}

	var elements []Element

	for {
		b := p.read()
		if b == 0x00 {
			return result.Success(Array{
				Elements: elements,
			})
		}

		p.rewind()
		res := p.parseElement()
		if res.HasFailed() {
			return result.Result[Array]{}.Failed(res.UnwrapError())
		}
		elements = append(elements, res.Expect())
	}
}

// Parse string
func (p *Parser) parseString() result.Result[string] {
	var buf bytes.Buffer
	for {
		b := p.read()
		buf.WriteByte(b)
		if b == 0 {
			return result.Success(buf.String())
		}
	}
}

// Parse boolean
func (p *Parser) parseBoolean() result.Result[bool] {
	b := p.read()

	if b > 0 {
		return result.Success(true)
	} else {
		return result.Success(false)
	}
}

// Parse datetime
func (p *Parser) parseDatetime() result.Result[time.Time] {
	b, _ := p.r.ReadBytes(8)
	bits := binary.LittleEndian.Uint64(b)
	t := time.Unix(int64(bits), 0)
	return result.Success(t)
}

// Parse double
func (p *Parser) parseFloat64() result.Result[float64] {
	b, _ := p.r.ReadBytes(8)
	bits := binary.LittleEndian.Uint64(b)
	f := math.Float64frombits(bits)
	return result.Success(f)
}

func (p *Parser) parseInt32() result.Result[int32] {
	b, _ := p.r.ReadBytes(4)
	bits := binary.LittleEndian.Uint32(b)
	return result.Success(int32(bits))
}

func (p *Parser) parseInt64() result.Result[int64] {
	b, _ := p.r.ReadBytes(8)
	bits := binary.LittleEndian.Uint64(b)
	return result.Success(int64(bits))
}

// Parse element
func (p *Parser) parseElement() result.Result[Element] {
	typ := p.read()
	keyRes := p.parseString()
	if keyRes.HasFailed() {
		return result.Result[Element]{}.Failed(keyRes.UnwrapError())
	}

	switch typ {
	// Double (float64)
	case 0x01:
		res := p.parseFloat64()
		if res.HasFailed() {
			return result.Result[Element]{}.Failed(res.UnwrapError())
		}
		return result.Success(Element{
			key:   keyRes.Expect(),
			value: Value{}.Float64(res.Expect()),
		})
	// String
	case 0x02:
		res := p.parseString()
		if res.HasFailed() {
			return result.Result[Element]{}.Failed(res.UnwrapError())
		}
		return result.Success(Element{
			key:   keyRes.Expect(),
			value: Value{}.String(res.Expect()),
		})
	// Document
	case 0x03:
		res := p.parseDocument()
		if res.HasFailed() {
			return result.Result[Element]{}.Failed(res.UnwrapError())
		}
		return result.Success(Element{
			key:   keyRes.Expect(),
			value: Value{}.Document(res.Expect()),
		})
	// Array
	case 0x04:
		res := p.parseArray()
		if res.HasFailed() {
			return result.Result[Element]{}.Failed(res.UnwrapError())
		}
		return result.Success(Element{
			key:   keyRes.Expect(),
			value: Value{}.Array(res.Expect()),
		})
	// Binary
	case 0x05:
		return result.Result[Element]{}.Failed(errors.New("not implemented"))
	// Undefined
	// DEPRECATED
	case 0x06:
		return result.Result[Element]{}.Failed(errors.New("deprecated"))
	// ObjectID (12 bytes)
	// Not implemented
	case 0x07:
		return result.Result[Element]{}.Failed(errors.New("not implemented"))
	// Boolean
	case 0x08:
		res := p.parseBoolean()
		if res.HasFailed() {
			return result.Result[Element]{}.Failed(res.UnwrapError())
		}
		return result.Success(Element{
			key:   keyRes.Expect(),
			value: Value{}.Boolean(res.Expect()),
		})
	// UTC Datetime
	case 0x09:
		res := p.parseDatetime()
		if res.HasFailed() {
			return result.Result[Element]{}.Failed(res.UnwrapError())
		}
		return result.Success(Element{
			key:   keyRes.Expect(),
			value: Value{}.Time(res.Expect()),
		})
	// Null value
	case 0x0A:
		return result.Success(Element{
			key:   keyRes.Expect(),
			value: Value{}.Null(),
		})
	// Regex value
	case 0x0B:
		patternRes := p.parseString()
		if patternRes.HasFailed() {
			return result.Result[Element]{}.Failed(patternRes.UnwrapError())
		}
		optionsRes := p.parseString()
		if optionsRes.HasFailed() {
			return result.Result[Element]{}.Failed(optionsRes.UnwrapError())
		}
		regex, err := regexp.Compile(patternRes.Expect())
		if err == nil {
			return result.Result[Element]{}.Failed(err)
		}
		return result.Success(Element{
			key:   keyRes.Expect(),
			value: Value{}.Regex(regex),
		})
	// DBPointer
	// DEPRECATED
	case 0x0C:
		return result.Result[Element]{}.Failed(errors.New("deprecated"))
	// Javascript code, treated like string
	case 0x0D:
		res := p.parseString()
		if res.HasFailed() {
			return result.Result[Element]{}.Failed(res.UnwrapError())
		}
		return result.Success(Element{
			key:   keyRes.Expect(),
			value: Value{}.String(res.Expect()),
		})
	// Symbol, Javascript code with scope
	// DEPRECATED
	case 0x0E, 0x0F:
		return result.Result[Element]{}.Failed(errors.New("deprecated"))
	// Int32
	case 0x10:
		res := p.parseInt32()
		if res.HasFailed() {
			return result.Result[Element]{}.Failed(res.UnwrapError())
		}
		return result.Success(Element{
			key:   keyRes.Expect(),
			value: Value{}.Integer(int(res.Expect())),
		})
	// Timestamp
	case 0x11:
		res := p.parseDatetime()
		if res.HasFailed() {
			return result.Result[Element]{}.Failed(res.UnwrapError())
		}
		return result.Success(Element{
			key:   keyRes.Expect(),
			value: Value{}.Time(res.Expect()),
		})
	// Int64
	case 0x12:
		res := p.parseInt64()
		if res.HasFailed() {
			return result.Result[Element]{}.Failed(res.UnwrapError())
		}
		return result.Success(Element{
			key:   keyRes.Expect(),
			value: Value{}.Integer(int(res.Expect())),
		})
	// Decimal (float128)
	// Not implemented
	case 0x13:
		return result.Result[Element]{}.Failed(errors.New("not implemented"))
	// Min, Max Keys
	case 0xFF, 0x7F:
		return result.Success(Element{})
	default:
		return result.Result[Element]{}.Failed(errors.New("invalid"))

	}
}
