package main

import (
	"fmt"
	"time"
)

func Hanoi(n int, from, to, via string, move func(from, to string)) {
	if n == 1 {
		move(from, to)
	} else {
		Hanoi(n-1, from, via, to, move)
		move(from, to)
		Hanoi(n-1, via, to, from, move)
	}
}

func counter() func() int {
	value := 0
	return func() int {
		value++
		return value
	}
}

func timeTracker() func() {
	started := time.Now()
	return func() {
		fmt.Printf("Времени прошло %s\n", time.Since(started).Round(time.Millisecond))
	}
}

func main() {
	towerA := "A"
	towerB := "B"
	towerC := "C"

	step := counter()
	timer := timeTracker()

	defer timer()
	disks := 20
	if disks > 5 {
		move := func(from, to string) {
			fmt.Printf("%d ходов \n", step())
		}
		Hanoi(disks, towerA, towerC, towerB, move)
	} else {
		move := func(from, to string) {
			fmt.Printf("%d: %s -> %s\n", step(), from, to)
		}
		Hanoi(disks, towerA, towerC, towerB, move)
	}
}
