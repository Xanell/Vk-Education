package main

import "fmt"

/*
Task
Вклад 100000 ₽ под 12% годовых, капитализация ежемесячная: каждый месяц прибавляется `ставка / 12` процентов от текущего баланса.
1. Выведите баланс на конец каждого месяца за год:
Пример:
месяц  1: 101000.00
месяц  2: 102010.00
месяц  3: 103030.10
месяц 12: 112682.50
2. Через сколько месяцев вклад удвоится? Используйте `for` только с условием.
Пример:
Вклад удвоится через 70 мес. Баланс: 200676.34
Вопрос: почему в реальных системах деньги не хранят во `float64`?
Ответ: проблема возникает при округлении чисел, если прибавить 0.1 к 0.2 у нас получится 0.3 но с хвостиком 04 на концею
Соответственно будет расхождение при большом числе операциий. Поэтому деньги лучше хранить в Int
Все суммы нужно хранить в минимальных единицах счисления
*/

const monthsInYear = 12

func main() {
	deposit := 100000.00
	depositDuration := 12
	depositRate := 0.12
	balance := deposit

	monthlyRate := depositRate / monthsInYear

	for i := 0; i < depositDuration; i++ {
		monthProfit := balance * monthlyRate
		balance += monthProfit
		fmt.Printf("Месяц %d: %.2f\n", i+1, balance)
	}
	monthCounter := 0
	targetBalance := deposit * 2
	balance = deposit

	for balance < targetBalance {
		monthProfit := balance * monthlyRate
		balance += monthProfit
		monthCounter++
	}
	fmt.Printf("Вклад удвоится через %d мес. Баланс %.2f\n", monthCounter, balance)
	a, b := 0.1, 0.2
	fmt.Println(a + b)
}
