package main

import "fmt"

func main() {
	var numbers []int = []int{21, 4, 2, 13, 10, 0, 19, 11, 7, 5, 23, 18, 9, 14, 6, 8, 1, 20, 17, 3, 16, 22, 24, 15, 12}
	fmt.Println("Unsorted:", numbers)
	var result []int = mergeSort(numbers)
	fmt.Println("Sorted:  ", result)
}

func mergeSort(data []int) []int {
	if len(data) <= 1 {
		return data
	}
	middle := len(data) / 2
	left := mergeSort(data[:middle])
	right := mergeSort(data[middle:])

	return merge(left, right)
}

func merge(left, right []int) []int {
	result := make([]int, len(left)+len(right))
	for i := 0; len(left) > 0 || len(right) > 0; i++ {
		if len(left) > 0 && len(right) > 0 {
			if left[0] < right[0] {
				result[i] = left[0]
				left = left[1:]
			} else {
				result[i] = right[0]
				right = right[1:]
			}
		} else if len(left) > 0 {
			result[i] = left[0]
			left = left[1:]
		} else if len(right) > 0 {
			result[i] = right[0]
			right = right[1:]
		}
	}
	return result
}
