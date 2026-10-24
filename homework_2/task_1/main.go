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

func Calculate(a float64, op string, b float64) (float64, error) {
	operation, ok := operations[op]
	if !ok {
		return 0, fmt.Errorf("ошибка: неизвестная операция %q", op)
	}
	return operation(a, b)
}

func main() {
	var result float64
	var err error
	a := 10.0
	b := 5.0
	c := 0.0
	result, err = Calculate(a, "+", b)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Printf("%.2f * %.2f = %.2f\n", a, b, result)
	}
	result, err = Calculate(a, "-", b)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Printf("%.2f * %.2f = %.2f\n", a, b, result)
	}
	result, err = Calculate(a, "/", b)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Printf("%.2f * %.2f = %.2f\n", a, b, result)
	}
	result, err = Calculate(a, "*", b)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Printf("%.2f * %.2f = %.2f\n", a, b, result)
	}
	result, err = Calculate(a, "/", c)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Printf("%.2f / %.2f = %.2f\n", a, c, result)
	}
	result, err = Calculate(a, "%", b)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Printf("%.2f %.2f = %.2f\n", a, b, result)
	}
}
