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
	seconds := 86405
	days := seconds / 86400
	seconds %= 86400

	hours := seconds / 3600
	seconds %= 3600

	minutes := seconds / 60
	seconds %= 60

	fmt.Printf("%dд. %dч. %dм. %dс.\n", days, hours, minutes, seconds)
}
