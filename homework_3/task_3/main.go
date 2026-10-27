package main

import "fmt"

func firstUnique(s string) rune {
	symbolsRune := map[rune]int{}
	for _, r := range s {
		symbolsRune[r]++
	}
	for _, r := range s {
		if symbolsRune[r] == 1 {
			return r
		}
	}
	return 0
}

func main() {
	fmt.Printf("%q\n", firstUnique("abacabad"))    // -> 'c'
	fmt.Printf("%q\n", firstUnique("abcabcd"))     // -> 'd'
	fmt.Printf("%q\n", firstUnique("абракадабра")) // -> 'к'
	fmt.Printf("%q\n", firstUnique("aabbcc"))      // -> '\x00'
}
