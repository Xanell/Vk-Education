package main

import (
	"fmt"
	"strings"
)

func areAnagrams(s1, s2 string) bool {
	s1 = strings.ToLower(s1)
	s2 = strings.ToLower(s2)
	s1 = strings.ReplaceAll(s1, " ", "")
	s2 = strings.ReplaceAll(s2, " ", "")
	if len(s1) != len(s2) {
		return false
	}
	symbolCounts := map[rune]int{}
	for _, symbol := range s1 {
		symbolCounts[symbol]++
	}

	for _, symbol := range s2 {
		symbolCounts[symbol]--
		if symbolCounts[symbol] < 0 {
			return false
		}
	}
	return true
}

func main() {
	fmt.Println(areAnagrams("listen ", "silent"))
	fmt.Println(areAnagrams("Eleven plus two", "Twelve plus one"))
	fmt.Println(areAnagrams("Апельсин", "спаниель"))
	fmt.Println(areAnagrams("Hello", "World"))
}
