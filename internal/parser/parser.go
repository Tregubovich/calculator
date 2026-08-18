package parser

import (
	"errors"
	"math"
	"strconv"
	"strings"
)

const (
	DIGITS     = "0123456789"
	OPERATIONS = "+-*/^"
	LETTERS    = "pelsc"
)

type Reader interface {
	Set(text string)
	Eof() bool
	Next() (byte, error)
	Test(set string) (bool, byte)
	Take(set string) (bool, byte)
	TakeSeq(seq string) bool
	Expect(set string) (byte, error)
}

type parser struct {
	reader Reader
}

func New(reader Reader) *parser {
	return &parser{
		reader: reader,
	}
}

func (p *parser) Parse(text string) (float64, error) {
	p.reader.Set(text)
	res, err := p.parseExpr()
	if err != nil {
		return 0, err
	}
	if !p.reader.Eof() {
		char, _ := p.reader.Next()
		return 0, errors.New("Expected EOF, got " + string(char))
	}
	return res, nil
}

func (p *parser) parseExpr() (float64, error) {
	return p.parseBinOp(0)
}

type BinOp struct {
	sym   byte
	prior int
	right bool
	op    func(float64, float64) (float64, error)
}

var (
	PRIORS = map[byte]int{
		'+': 1,
		'-': 2,
		'*': 10,
		'/': 11,
		'^': 20,
	}
	BINOPS = []*BinOp{
		{
			sym:   '+',
			prior: 1,
			op: func(f float64, s float64) (float64, error) {
				return f + s, nil
			},
		},
		{
			sym:   '-',
			prior: 2,
			op: func(f float64, s float64) (float64, error) {
				return f - s, nil
			},
		},
		{
			sym:   '*',
			prior: 10,
			op: func(f float64, s float64) (float64, error) {
				return f * s, nil
			},
		},
		{
			sym:   '/',
			prior: 11,
			op: func(f float64, s float64) (float64, error) {
				if s == 0 {
					return 0, errors.New("division by zero")
				}
				return f / s, nil
			},
		},
		{
			sym:   '^',
			prior: 20,
			right: true,
			op: func(f float64, s float64) (float64, error) {
				if f < 0 && float64(int(s)) != s {
					return 0, errors.New("with non-integer power base must be non-negative")
				}
				if f == 0 && s < 0 {
					return 0, errors.New("can't raise 0 to a negative power")
				}
				return math.Pow(f, s), nil
			},
		},
	}
)

func (p *parser) parseBinOp(curBinOp int) (float64, error) {
	if curBinOp >= len(BINOPS) {
		return p.parseTerm()
	}
	ok, _ := p.reader.Test("-" + "(" + DIGITS + LETTERS)
	if ok {
		res, err := p.parseBinOp(curBinOp + 1)
		if err != nil {
			return 0, err
		}
		for {
			if ok, _ := p.reader.Take(string(BINOPS[curBinOp].sym)); ok {
				var secondTerm float64
				var err error
				if BINOPS[curBinOp].right {
					secondTerm, err = p.parseBinOp(curBinOp)
				} else {
					secondTerm, err = p.parseBinOp(curBinOp + 1)
				}
				if err != nil {
					return 0, err
				}
				res, err = BINOPS[curBinOp].op(res, secondTerm)
				if err != nil {
					return 0, err
				}
			} else {
				break
			}
		}
		next, err := p.reader.Expect(OPERATIONS)
		if err != nil && !p.reader.Eof() && next != ')' {
			return 0, err
		}
		if next == ')' || PRIORS[next] < BINOPS[curBinOp].prior {
			return res, nil
		}

	}
	if p.reader.Eof() {
		return 0, errors.New("got EOF")
	}
	unknown, err := p.reader.Next()
	if err != nil {
		return 0, err
	}
	return 0, errors.New("Unknown char: " + string(unknown))
}

var (
	CONSTANTS = map[string]float64{
		"pi":  math.Pi,
		"e":   math.E,
		"phi": math.Phi,
	}

	FUNCTIONS = map[string]func(float64) (float64, error){
		"log2": func(x float64) (float64, error) {
			return math.Log2(x), nil
		},
		"log10": func(x float64) (float64, error) {
			return math.Log10(x), nil
		},
		"ln": func(x float64) (float64, error) {
			return math.Log(x), nil
		},
		"sqrt": func(x float64) (float64, error) {
			return math.Sqrt(x), nil
		},
		"cbrt": func(x float64) (float64, error) {
			return math.Cbrt(x), nil
		},
	}
)

func (p *parser) parseTerm() (float64, error) {
	if ok, _ := p.reader.Take("("); ok {
		res, err := p.parseExpr()
		if err != nil {
			return 0, err
		}
		if _, err := p.reader.Expect(")"); err != nil {
			return 0, err
		}
		p.reader.Take(")")
		return res, nil
	}
	if ok, _ := p.reader.Take("-"); ok {
		res, err := p.parseTerm()
		if err != nil {
			return 0, err
		}
		return -res, nil
	}
	if ok, _ := p.reader.Test(DIGITS); ok {
		res, err := p.parseNum()
		if err != nil {
			return 0, err
		}
		return res, nil
	}
	if ok, _ := p.reader.Test(LETTERS); ok {
		for name, value := range CONSTANTS {
			if p.reader.TakeSeq(name) {
				return value, nil
			}
		}

		for name, op := range FUNCTIONS {
			if p.reader.TakeSeq(name) {
				_, err := p.reader.Expect("(")
				if err != nil {
					return 0, err
				}
				p.reader.Take("(")
				res, err := p.parseExpr()
				if err != nil {
					return 0, err
				}
				if _, err := p.reader.Expect(")"); err != nil {
					return 0, err
				}
				p.reader.Take(")")
				return op(res)
			}
		}
	}
	return 0, errors.New("unknown letters")
}

func (p *parser) parseNum() (float64, error) {
	res := strings.Builder{}
	if ok, _ := p.reader.Take("-"); ok {
		res.WriteByte('-')
	}
	for {
		ok, d := p.reader.Take("." + DIGITS)
		if !ok {
			break
		}
		res.WriteByte(d)
	}
	return strconv.ParseFloat(res.String(), 64)
}
