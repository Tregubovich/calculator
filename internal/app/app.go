package app

import (
	"calculator/internal/calculator"
	charreader "calculator/internal/char-reader"
	"calculator/internal/interactor"
	"calculator/internal/parser"
	"fmt"
)

func Run() {
	interactor := interactor.NewCLI()

	reader := charreader.New()
	parser := parser.New(reader)

	err := calculator.Calculate(parser, interactor)
	if err != nil {
		fmt.Println(err)
	}

}
