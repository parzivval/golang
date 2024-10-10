package main

import "fmt"

func add(a int, b int) int {
	return a + b
}

func swap(x, y string) (string, string) {
	return y, x
}

func quotientRemainder(dividend, divisor int) (int, int) {
	quotient := dividend / divisor
	remainder := dividend % divisor
	return quotient, remainder
}

func main() {

	sum := add(5, 10)
	fmt.Printf("Sum of 5 and 10: %d\n", sum)

	a, b := swap("Hello", "World")
	fmt.Printf("Swapped: %s, %s\n", a, b)

	quotient, remainder := quotientRemainder(25, 4)
	fmt.Printf("Quotient: %d, Remainder: %d\n", quotient, remainder)
}
