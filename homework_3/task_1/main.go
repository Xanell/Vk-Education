package main

import (
	"fmt"
	"strings"
)

func wordCount(s string) map[string]int {
	wordCounts := map[string]int{}
	words := strings.Fields(s)
	for _, word := range words {
		wordCounts[word]++
	}
	return wordCounts
}

func main() {
	fmt.Println(wordCount("go is fun go"))
	fmt.Println(wordCount("Hello hello World"))
	fmt.Println(wordCount(" мир мир труд "))
}
