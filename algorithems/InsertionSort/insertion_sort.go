package main

import "fmt"

func main() {
	var numbers []int = []int{21, 4, 2, 13, 10, 0, 19, 11, 7, 5, 23, 18, 9, 14, 6, 8, 1, 20, 17, 3, 16, 22, 24, 15, 12}
	fmt.Println("Unsorted:", numbers)
	insertionSort(numbers)
	fmt.Println("Sorted:  ", numbers)
}

func insertionSort(arr []int) []int {
	for currentIndex := 1; currentIndex < len(arr); currentIndex++ {
		temporary := arr[currentIndex]
		iterator := currentIndex
		for ; iterator > 0 && arr[iterator-1] >= temporary; iterator-- {
			arr[iterator] = arr[iterator-1]
		}
		arr[iterator] = temporary
	}
	return arr
}
