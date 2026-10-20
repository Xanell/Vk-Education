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
Ответ: возможно возникает проблема с округлением при расчетах
*/

func main() {
	deposit := 100000.00
	depositDuration := 12
	depositRate := 0.12
	expectedDeposit := deposit

	for i := 0; i < depositDuration; i++ {
		depositProfit := expectedDeposit * (1 + depositRate/12)
		expectedDeposit = depositProfit
		fmt.Printf("Месяц %d: %.2f\n", i+1, expectedDeposit)
	}
	monthCounter := 0
	targetDeposit := deposit * 2
	expectedDeposit = deposit

	for expectedDeposit < targetDeposit {
		depositProfit := expectedDeposit * (1 + depositRate/12)
		expectedDeposit = depositProfit
		monthCounter++
	}
	fmt.Printf("Вклад удвоится через %d мес. Баланс %.2f", monthCounter, expectedDeposit)
}
