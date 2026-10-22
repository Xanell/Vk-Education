package main

import "fmt"

/*
Переведите секунды в дни, часы, минуты и секунды.
Нулевые части не выводим, для нуля выводим `0с`.
Считайте через целочисленное деление `/` и остаток `%`, пакет `time` не используйте.
Пример:
93784 -> 1д 2ч 3м 4с
86405 -> 1д 5с
3600  -> 1ч
59    -> 59с
0     -> 0с
*/

func main() {
	seconds := 180
	if seconds == 0 {
		fmt.Printf("%dс", seconds)
	}
	days := seconds / 86400
	if days != 0 {
		fmt.Printf("%dд ", days)
	}
	seconds %= 86400

	hours := seconds / 3600
	if hours != 0 {
		fmt.Printf("%dч ", hours)
	}
	seconds %= 3600

	minutes := seconds / 60
	if minutes != 0 {
		fmt.Printf("%dм ", minutes)
	}
	seconds %= 60
	if seconds != 0 {
		fmt.Printf("%dс ", seconds)
	}
}
