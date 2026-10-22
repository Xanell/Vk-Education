package main

import (
	"errors"
	"fmt"
)

type Operation func(a, b float64) (float64, error)

func addition(a, b float64) (float64, error) {
	return a + b, nil
}

func subtraction(a, b float64) (float64, error) {
	return a - b, nil
}

func multiplication(a, b float64) (float64, error) {
	return a * b, nil
}

func division(a, b float64) (float64, error) {
	if b == 0 {
		return 0, errors.New("ошибка: деление на ноль")
	}
	return a / b, nil
}

var operations = map[string]Operation{
	"+": addition,
	"-": subtraction,
	"*": multiplication,
	"/": division,
}

func Calculate(a float64, op Operation, b float64) (float64, error) {
	if op == nil {
		// Как правильно форматировать вывод если op кастомного Типа?
		err := fmt.Errorf("ошибка: неизвестная операция ", op)
		return 0, err
	}
	result, err := op(a, b)
	return result, err
}

func main() {
	a := 10.0
	b := 5.0
	//op := "%"
	result, err := Calculate(a, operations["%"], b)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("%.2f", result)
}
