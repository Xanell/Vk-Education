package main

import "fmt"

func Hanoi(n int, from, to, via string, move func(from, to string)) {

}

func main() {
	disks := 3
	towerA := "A"
	towerB := "B"
	towerC := "C"
	move := func(from, to string) {}
	fmt.Println(Hanoi(disks, towerA, towerC, towerB, move(towerA, towerC)))
}
