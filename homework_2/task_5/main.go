package main

import "fmt"

func countBits(n uint8) int {
	count := 0
	for n > 0 {
		if n&1 == 1 {
			count++
		}
		n = n >> 1
	}
	return count
}

func isPowerOfTwo(n uint) bool {
	return n != 0 && (n&(n-1)) == 0
}

func main() {
	n := uint8(255)
	fmt.Println(countBits(n))
	fmt.Println(isPowerOfTwo(32))
}
