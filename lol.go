package main

import (
	"fmt"
	"lol/scan"
)

func main() {
	scanner := scan.NewScanner("/**/let x = oo")
	tokens, e := scanner.Scan()
	fmt.Printf("%v\n", tokens)
	fmt.Printf("%v\n", e)
}
