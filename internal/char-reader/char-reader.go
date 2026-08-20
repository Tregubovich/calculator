package char_reader

import (
	"errors"
	"io"
	"unicode"
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
	c.skip()
	if !c.Eof() {
		return c.text[c.cur], nil
	}
	return 0, io.EOF
}

func (c *charReader) Test(set string) (bool, byte) {
	c.skip()
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

func (c *charReader) TakeSeq(seq string) bool {
	c.skip()
	if c.cur+len(seq) <= len(c.text) && c.text[c.cur:c.cur+len(seq)] == seq {
		c.cur += len(seq)
		return true
	}
	return false
}

func (c *charReader) Expect(set string) (byte, error) {
	ok, char := c.Test(set)
	if !ok {
		return char, errors.New("Expected: " + set + ", got " + c.curChar())
	}
	return char, nil
}

func (c *charReader) skip() {
	for c.cur < len(c.text) && unicode.IsSpace(rune(c.text[c.cur])) {
		c.cur++
	}
}

func (c *charReader) curChar() string {
	if c.cur >= len(c.text) {
		return "EOF"
	}
	return string(c.text[c.cur])
}
