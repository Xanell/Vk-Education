package main

// Вопрос: почему решение через append(nums[:index], nums[index+1:]...) испортит исходный слайс?
// Как я понимаю потому что слайс это ссылочный тип, и мы если положим его просто в другую переменную то мы просто скопируем ссылку на него поэтмоу все изменения затрагивают основной слайс

import (
	"fmt"
	"slices"
)

func removeElement(nums []int, index int) ([]int, error) {
	if index >= len(nums) || index < 0 {
		err := fmt.Errorf("индекс %d вне диапазона [0, %d)", index, len(nums)-1)
		return nil, err
	}
	newSlice := slices.Clone(nums)
	return append(newSlice[:index], newSlice[index+1:]...), nil
}

func main() {
	testNums := [][]int{
		{1, 2, 2, 3, 1},
		{1, 2, 3},
		{1, 2, 3},
	}
	testIndex := []int{0, 5, -1}
	for i, num := range testNums {
		result, err := removeElement(num, testIndex[i])
		if err != nil {
			fmt.Println(err)
		} else {
			fmt.Printf("%v %v, исходный слайс: %v\n", result, err, num)
		}
	}

	// Для разбора
	newNums := []int{0, 1, 2, 3, 4}
	// под капотом наш newNums выглядит как структура трех полей {ptr, len, cap}
	index := 2
	a := newNums[:index]   // Получается A копирует ptr от newNums, len = index(2), а cap = 5 потому что наш срез начинается с 0 индекса и можно считать новый cap как cap_старого - начало_среза (5-0=5)
	b := newNums[index+1:] // B так же копирует ptr от newNums, len так же равен index(2), а вот cap же будет (5-3=2)
	res := append(a, b...) // И когда мы делаем append получается берем A = {0,1} len = 2 cap = 5, и пишем туда элементы из B = {3,4} len = 2, cap = 2, и так как мы пишем все в элемент A а у него cap = 5, и он в памяти указывает на тот же масив, то в res у нас получается правильный слайс а исходный слайс портится
	fmt.Println(newNums)
	fmt.Println(res)
}
