package main

import "fmt"

func datatype() {
	var a int = 10
	var b float64 = 20.5
	var c string = "Bakytzhan"
	var d bool = true

	e := 15.5

	fmt.Printf("a: %v, Type: %T\n", a, a)
	fmt.Printf("b: %v, Type: %T\n", b, b)
	fmt.Printf("c: %v, Type: %T\n", c, c)
	fmt.Printf("d: %v, Type: %T\n", d, d)
	fmt.Printf("e: %v, Type: %T\n", e, e)
}

func main() {
	datatype()
}
