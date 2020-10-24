package main

import (
	"fmt"
	"container/list"
)


func main() {

	mylist := list.New()
	mylist.PushBack(1)
	mylist.PushFront(2)
	
	for element := mylist.Front(); element != nil; element = element.Next() {
		
		if element.Value != 1 {
			mylist.Remove(element)
		}
	}
	
	for element := mylist.Front(); element != nil; element = element.Next() {
		
		fmt.Println(element.Value)
	}

}