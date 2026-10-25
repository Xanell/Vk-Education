package main

// Вопрос: почему решение через append(nums[:index], nums[index+1:]...) испортит исходный слайс?
// Как я понимаю потому что слайс это ссылочный тип, и мы если положим его просто в другую переменную то мы просто скопируем ссылку на него поэтмоу все изменения затрагивают основной слайс

import (
	"fmt"
	"slices"
)

func removeElement(nums []int, index int) ([]int, error) {
	if index >= len(nums) || index < 0 {
		err := fmt.Errorf("индекс %d вне диапазона [0, %d]", index, len(nums)-1)
		return nil, err
	}
	newSlice := slices.Clone(nums)
	newSlice[index] = nums[index+1]
	return newSlice[:len(newSlice)-1], nil
}

func main() {
	nums := []int{1, 2, 2, 3, 1}
	result, err := removeElement(nums, 3)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Printf("%v %v, исходный слайс: %v\n", result, err, nums)
	}
	result, err = removeElement(nums, 5) //-> [] индекс 5 вне диапазона [0, 3)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Printf("%v %v, исходный слайс: %v\n", result, err, nums)
	}
	result, err = removeElement(nums, -1) //-> [] индекс -1 вне диапазона [0, 3)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Printf("%v %v, исходный слайс: %v\n", result, err, nums)
	}
}
