package interactor

import (
	"bufio"
	"os"
)

type CLI struct {
	reader bufio.Reader
	writer bufio.Writer
	errors bufio.Writer
}

func NewCLI() *CLI {
	return &CLI{
		reader: *bufio.NewReader(os.Stdin),
		writer: *bufio.NewWriter(os.Stdout),
		errors: *bufio.NewWriter(os.Stderr),
	}
}

func (c *CLI) Input() (string, error) {
	text, err := c.reader.ReadString('\n')
	if err != nil {
		return "", err
	}
	return text[:len(text)-1], nil
}

func (c *CLI) Output(text string) error {
	_, err := c.writer.WriteString(text + "\n")
	if err != nil {
		return err
	}
	err = c.writer.Flush()
	if err != nil {
		return err
	}
	return nil
}

func (c *CLI) Error(err error) {
	c.writer.WriteString("Error: " + err.Error() + "\n")
	c.writer.Flush()
}
