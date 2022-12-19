package main

import "fmt"

func someFuncInAnotherFile() {
	const op = "someFuncInAnotherFile"

	fmt.Println("yo")
}

func anotherFuncInAnotherFile() {
	const op = "someFuncInAnotherFile"

	fmt.Println("yo")
}
