package number

import "errors"

// MinNumber finds minimums value in a list of arguments. Raise error if list of arguments is empty
func MinNumber[V int32 | int64 | float32 | float64](numbers ...V) (V, error) {
	l := len(numbers)
	if l < 1 {
		return 0, errors.New("arguments are required")
	}
	m := numbers[0]
	for _, num := range numbers[1:] {
		if num < m {
			m = num
		}
	}
	return m, nil
}

// MaxNumber finds maximum value in a list of arguments. Raise error if list of arguments is empty
func MaxNumber[V int32 | int64 | float32 | float64](numbers ...V) (V, error) {
	l := len(numbers)
	if l < 1 {
		return 0, errors.New("arguments are required")
	}
	m := numbers[0]
	for _, num := range numbers[1:] {
		if num > m {
			m = num
		}
	}
	return m, nil
}
