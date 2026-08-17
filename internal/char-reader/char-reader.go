package char_reader

import (
	"errors"
	"io"
)

type charReader struct {
	text string
	cur  int
}

func New() *charReader {
	return &charReader{}
}

func (c *charReader) Set(text string) {
	c.text = text
	c.cur = 0
}

func (c *charReader) Eof() bool {
	if c.cur >= len(c.text) {
		return true
	}
	return false
}

func (c *charReader) Next() (byte, error) {
	if !c.Eof() {
		return c.text[c.cur], nil
	}
	return 0, io.EOF
}

func (c *charReader) Test(set string) (bool, byte) {
	if c.Eof() {
		return false, 0
	}
	for i := range set {
		if c.text[c.cur] == set[i] {
			return true, set[i]
		}
	}
	return false, c.text[c.cur]
}

func (c *charReader) Take(set string) (bool, byte) {
	ok, char := c.Test(set)
	if ok {
		c.cur++
	}
	return ok, char
}

func (c *charReader) Expect(set string) (byte, error) {
	ok, char := c.Test(set)
	if !ok {
		return char, errors.New("Expected: " + set + ", got " + c.curChar())
	}
	return char, nil
}

func (c *charReader) curChar() string {
	if c.cur >= len(c.text) {
		return "EOF"
	}
	return string(c.text[c.cur])
}
