package main

import (
	"fmt"
	"strings"
)

func isPalindrome(s string) bool {
	s = strings.ToLower(s)
	replacer := strings.NewReplacer(".", "", " ", "", "!", "", ",", "")
	s = replacer.Replace(s)
	fmt.Println([]rune(s))
	return false
}

func main() {
	//fmt.Println(isPalindrome("шалаш")) // true
	//fmt.Println(isPalindrome("А роза упала на лапу Азора"))   // true
	fmt.Println(isPalindrome("Was it a car or a cat I saw?")) // true
	//fmt.Println(isPalindrome("hello")) // false
}
