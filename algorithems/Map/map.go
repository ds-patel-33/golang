package main

import "fmt"

type Service interface{
	SayHi()
}

type MyService struct{}
func (s MyService) SayHi() {
	fmt.Println("Hi")
}

type SecondService struct{}
func (s SecondService) SayHi() {
	fmt.Println("Hello From the 2nd Service")
}

func main() {
	fmt.Println("Go Maps Tutorial")

	interfaceMap := make(map[string]Service)
	
	interfaceMap["SERVICE-ID-1"] = MyService{}
	interfaceMap["SERVICE-ID-2"] = SecondService{}

	
	interfaceMap["SERVICE-ID-2"].SayHi()

}