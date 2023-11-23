package number

import "errors"

// MinNumber finds and returns the minimum number from a variadic list of numbers.
//
// It takes in one or more numbers of type int32, int64, float32, or float64.
// It returns the minimum number found from the given list of numbers and an error.
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

// MaxNumber returns the maximum value from a list of numbers.
//
// The function takes a variadic parameter of type V, which can be int32, int64, float32, or float64.
// It returns the maximum value found and an error if the list is empty.
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
