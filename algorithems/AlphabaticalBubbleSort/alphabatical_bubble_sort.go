package main

import "fmt"

func main() {
	var animals []string = []string{
		"pan",
		"cat",
		"chair",
		"laptop",
		"rat",
		"moose",
		"colt",
		"bike",
		"car",
		"cable",
	}

	fmt.Println("Unsorted:", animals)
	bubbleSort(animals)
	fmt.Println("Sorted:  ", animals)
}

func bubbleSort(animals []string) {
	var N int = len(animals)
	var i int
	for i = 0; i < N; i++ {
		sweep(animals)
	}
}

func sweep(animals []string) {
	var N int = len(animals)
	var firstIndex int = 0
	var secondIndex int = 1

	for secondIndex < N {
		var firstVal string = animals[firstIndex]
		var secondVal string = animals[secondIndex]

		if firstVal > secondVal {
			animals[firstIndex] = secondVal
			animals[secondIndex] = firstVal
		}

		firstIndex++
		secondIndex++
	}
}
