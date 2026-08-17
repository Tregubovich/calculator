package calculator

import "strconv"

type (
	Interactor interface {
		Input() (string, error)
		Output(text string) error
		Error(err error)
	}
	Parser interface {
		Parse(text string) (float64, error)
	}
)

func Calculate(parser Parser, interactor Interactor) error {
	for {
		text, err := interactor.Input()
		if err != nil {
			return err
		}
		if text == "" {
			break
		}

		res, parseErr := parser.Parse(text)
		if parseErr != nil {
			interactor.Error(parseErr)
		}
		err = interactor.Output(strconv.FormatFloat(res, 'f', -1, 64))
		if err != nil {
			return err
		}
	}
	return nil
}
