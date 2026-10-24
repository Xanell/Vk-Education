package main

// Вопрос: почему ходов ровно 2ⁿ − 1?
// Получается что самый нижний большой диск движется ровно 1 раз
// А все верхние диски делают каждый раз в 2 раза больше
// для n = 3, самый нижний диск сделает 1 ход, диск над ним сделает 2 шага и самый маленький диск сделает в 2 раза больше предыдущего т.е 4, равно 7 шагов

import (
	"fmt"
	"time"
)

func Hanoi(n int, from, to, via string, move func(from, to string)) {
	if n == 0 {
		return
	}
	Hanoi(n-1, from, via, to, move)
	move(from, to)
	Hanoi(n-1, via, to, from, move)
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

func runTimerAndTask(n int, from, to, via string, move func(from, to string)) {
	timer := timeTracker()
	defer timer()
	Hanoi(n, from, to, via, move)
}

func main() {
	towerA := "A"
	towerB := "B"
	towerC := "C"

	// Задача для 3 дисков
	fmt.Println("---Задача для 3 дисков---")
	disks := 3
	stepA := counter()
	move := func(from, to string) {
		fmt.Printf("%d: %s -> %s\n", stepA(), from, to)
	}
	runTimerAndTask(disks, towerA, towerC, towerB, move)

	// Задача для 20 дисков
	fmt.Println("---Задача для 20 дисков---")
	disks = 20
	stepB := counter()
	lastStep := 0
	move = func(from, to string) {
		lastStep = stepB()
	}
	runTimerAndTask(disks, towerA, towerC, towerB, move)
	fmt.Printf("%d ходов\n", lastStep)
}
