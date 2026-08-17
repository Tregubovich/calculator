package parser

import (
	"errors"
	"math"
	"strconv"
	"strings"
)

const (
	DIGITS = "0123456789"
)

type Reader interface {
	Set(text string)
	Eof() bool
	Next() (byte, error)
	Test(set string) (bool, byte)
	Take(set string) (bool, byte)
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
	return p.parseExpr()
}

func (p *parser) parseExpr() (float64, error) {
	return p.parseBinOp(0)
}

type BinOp struct {
	sym   byte
	prior int
	op    func(float64, float64) (float64, error)
}

var (
	OPS    = "+-*/^"
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
			op: func(f float64, s float64) (float64, error) {
				if f < 0 {
					return 0, errors.New("base must be non-negative")
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
	ok, char := p.reader.Test("(" + DIGITS)
	if ok {
		res, err := p.parseBinOp(curBinOp + 1)
		if err != nil {
			return 0, err
		}
		for {
			if ok, _ := p.reader.Take(string(BINOPS[curBinOp].sym)); ok {
				secondTerm, err := p.parseBinOp(curBinOp + 1)
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
		next, err := p.reader.Expect(")" + OPS)
		if p.reader.Eof() || next == ')' || PRIORS[next] < BINOPS[curBinOp].prior {
			return res, nil
		}
		if err != nil {
			return 0, err
		}
	}
	return 0, errors.New("Unknown char: " + string(char))
}

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
	if ok, _ := p.reader.Test(DIGITS); ok {
		res, err := p.parseNum()
		if err != nil {
			return 0, err
		}
		return res, nil
	}
	return 0, errors.New("unreachable")
}

func (p *parser) parseNum() (float64, error) {
	res := strings.Builder{}
	for {
		ok, d := p.reader.Take("." + DIGITS)
		if !ok {
			break
		}
		res.WriteByte(d)
	}
	return strconv.ParseFloat(res.String(), 64)
}
