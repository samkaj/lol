package main

import (
	"fmt"
	"lol/scan"
)

func main() {
	scanner := scan.NewScanner("let x = 10")
	tokens, errors := scanner.Scan()
	fmt.Printf("%v\n", tokens)
	fmt.Printf("%v\n", errors)
}
