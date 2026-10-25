package main

import (
	"fmt"
	"strings"
)

func isPalindrome(text string) bool {
	text = strings.ToLower(text)
	replacer := strings.NewReplacer(" ", "", ".", "", ",", "", "?", "", "!", "")
	text = replacer.Replace(text)
	textToRunes := []rune(text)
	reverseRunes := make([]rune, 0, len(textToRunes))
	for i := len(textToRunes) - 1; i >= 0; i-- {
		reverseRunes = append(reverseRunes, textToRunes[i])
	}
	for i := 0; i < len(textToRunes); i++ {
		if textToRunes[i] != reverseRunes[i] {
			return false
		}
	}
	return true
}

func main() {
	fmt.Println(isPalindrome("шалаш"))                        // true
	fmt.Println(isPalindrome("А роза упала на лапу Азора"))   // true
	fmt.Println(isPalindrome("Was it a car or a cat I saw?")) // true
	fmt.Println(isPalindrome("hello"))                        // false
}
