package interactor

import (
	"encoding/json"
	"net/http"
)

type Data struct {
	Text string `json:"text"`
}

type WebInteractor struct {
	input  chan string
	output chan string
	errors chan error
}

func NewWebInteractor() *WebInteractor {
	return &WebInteractor{
		input:  make(chan string),
		output: make(chan string),
		errors: make(chan error),
	}
}

func (i *WebInteractor) Input() (string, error) {
	text := <-i.input
	return text, nil
}

func (i *WebInteractor) Output(text string) error {
	i.output <- text
	return nil
}

func (i *WebInteractor) Error(err error) {
	i.errors <- err
}

func (i *WebInteractor) Handler(
	w http.ResponseWriter,
	r *http.Request,
) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var data Data

	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	i.input <- data.Text

	select {
	case result := <-i.output:
		json.NewEncoder(w).Encode(Data{
			Text: result,
		})

	case err := <-i.errors:
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
}
