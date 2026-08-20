package app

import (
	"calculator/internal/calculator"
	charreader "calculator/internal/char-reader"
	"calculator/internal/interactor"
	"calculator/internal/parser"
	"fmt"
	"net/http"
	"sync"
)

func RunWithCLI() {
	interactor := interactor.NewCLI()

	reader := charreader.New()
	parser := parser.New(reader)

	err := calculator.Calculate(parser, interactor)
	if err != nil {
		fmt.Println(err)
	}
}

func RunWithWebInteractor() {
	interactor := interactor.NewWebInteractor()

	reader := charreader.New()
	parser := parser.New(reader)

	wg := &sync.WaitGroup{}
	wg.Go(func() {
		_ = calculator.Calculate(parser, interactor)
	})

	mux := http.NewServeMux()

	mux.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("ok"))
	})
	mux.HandleFunc("/eval", interactor.Handler)

	fmt.Println("Starting server on port: 7070")
	err := http.ListenAndServe(":7070", mux)
	if err != nil {
		fmt.Println(err)
	}
}
