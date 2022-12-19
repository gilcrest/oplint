package main

import "fmt"

type helloWorld string

func (s helloWorld) Do() {
	const op = "DoNot"

	fmt.Println("Hello World")
}

func main() {
	const op = "main"

	fmt.Println("WTF")
	itWorks()
}

func itWorks() {
	const op = "itDoesNOTWork"
	fmt.Println("it works!")
}
