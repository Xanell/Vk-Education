package main

import "fmt"

/*
Task
Нарисуйте шахматную доску из пробелов и `#`. Размер задаётся константой `N`.
`N = 8`:
Пример:
 # # # #
# # # #
 # # # #
# # # #
 # # # #
# # # #
 # # # #
# # # #
Проверьте на `N = 5` и `N = 1`.
*/

const N int = 8

func main() {
	for i := 1; i <= N; i++ {
		for j := 1; j <= N; j++ {
			if (i+j)%2 == 0 {
				fmt.Print("#")
			} else {
				fmt.Print(" ")
			}
		}
		fmt.Println(" ")
	}
}
