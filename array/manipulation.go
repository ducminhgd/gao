package array

import "github.com/ducminhgd/gao/number"

// PrependN adds an element to the beginning of the array.
//
// Parameters:
// - array: the array to prepend the element to.
// - element: the element to add to the array.
// Returns the updated array.
func PrependN[T number.Number](array []T, element T) []T {
	array = append(array, 0)
	copy(array[1:], array)
	array[0] = element
	return array
}

// PrependS adds an element to the beginning of the string array.
//
// Parameters:
// - array: the string array to prepend the element to.
// - element: the string element to add to the array.
// Returns the updated string array.
func PrependS(array []string, element string) []string {
	array = append(array, "")
	copy(array[1:], array)
	array[0] = element
	return array
}
