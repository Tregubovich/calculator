package math

import "errors"

func Factorial(x float64) (float64, error) {
	if x < 0 {
		return 0, errors.New("factorial: negative number")
	}
	if float64(int(x)) != x {
		return 0, errors.New("factorial: float number")
	}
	if x == 0 {
		return 1, nil
	}
	res, err := Factorial(x - 1)
	if err != nil {
		return 0, err
	}
	return res * x, nil
}
