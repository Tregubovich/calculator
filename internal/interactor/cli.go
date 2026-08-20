package interactor

import (
	"bufio"
	"os"
	"strings"
)

type CLI struct {
	reader *bufio.Reader
	writer *bufio.Writer
}

func NewCLI() *CLI {
	return &CLI{
		reader: bufio.NewReader(os.Stdin),
		writer: bufio.NewWriter(os.Stdout),
	}
}

func (c *CLI) Input() (string, error) {
	c.writer.WriteString(">> ")
	c.writer.Flush()
	text, err := c.reader.ReadString('\n')
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(text), nil
}

func (c *CLI) print(w *bufio.Writer, text string) error {
	if _, err := w.WriteString(text + "\n"); err != nil {
		return err
	}
	return w.Flush()
}

func (c *CLI) Output(text string) error {
	return c.print(c.writer, text)
}

func (c *CLI) Error(err error) {
	c.print(c.writer, "\033[31mError: "+err.Error()+"\033[0m")
}
