package main

import "fmt"

// Константа для перевода в копейки
const subUnitAmount = 100

func formatedOutput(initial int64) string {
	initialRub := initial / subUnitAmount
	initialSubUnit := initial % subUnitAmount
	return fmt.Sprintf("%d.%02d", initialRub, initialSubUnit)
}

func NewAccount(initial int64) (deposit func(int64), withDraw func(int64) error, balance func() int64) {
	accountBalance := initial
	deposit = func(initial int64) {
		accountBalance += initial
	}
	withDraw = func(initial int64) error {
		if accountBalance < initial {
			err := fmt.Errorf("недостаточно средств (баланс %s)", formatedOutput(accountBalance))
			return err
		}
		accountBalance -= initial
		return nil
	}
	balance = func() int64 {
		return accountBalance
	}
	return deposit, withDraw, balance
}

func main() {
	var moneyToDeposit int64  // Переменная для пополения баланса в копейках
	var moneyToWithdraw int64 // Переменная для снятия баланса в копейках
	var err error
	fmt.Println("-----Счет А-----")
	depositA, withDrawA, balanceA := NewAccount(100000)
	fmt.Printf("Счёт А: %s\n", formatedOutput(balanceA()))
	moneyToDeposit = 10000
	depositA(moneyToDeposit)
	fmt.Printf("Пополнение %s -> %s\n", formatedOutput(moneyToDeposit), formatedOutput(balanceA()))
	moneyToWithdraw = 200000
	err = withDrawA(moneyToWithdraw)
	if err != nil {
		fmt.Printf("Снятие %s %v\n", formatedOutput(moneyToWithdraw), err)
	} else {
		fmt.Printf("Снятие %s -> %s\n", formatedOutput(moneyToWithdraw), formatedOutput(balanceA()))
	}
	moneyToWithdraw = 35000
	err = withDrawA(moneyToWithdraw)
	if err != nil {
		fmt.Printf("Снятие %s %v\n", formatedOutput(moneyToWithdraw), err)
	} else {
		fmt.Printf("Снятие %s -> %s\n", formatedOutput(moneyToWithdraw), formatedOutput(balanceA()))
	}

	fmt.Println("-----Счет B-----")
	depositB, withDrawB, balanceB := NewAccount(1000000)
	fmt.Printf("Счёт B: %s\n", formatedOutput(balanceB()))
	moneyToDeposit = 10000
	depositB(moneyToDeposit)
	fmt.Printf("Пополнение %s -> %s\n", formatedOutput(moneyToDeposit), formatedOutput(balanceB()))
	moneyToWithdraw = 200000
	err = withDrawB(moneyToWithdraw)
	if err != nil {
		fmt.Printf("Снятие %s %v\n", formatedOutput(moneyToWithdraw), err)
	} else {
		fmt.Printf("Снятие %s -> %s\n", formatedOutput(moneyToWithdraw), formatedOutput(balanceB()))
	}
	moneyToWithdraw = 35000
	err = withDrawB(moneyToWithdraw)
	if err != nil {
		fmt.Printf("Снятие %s %v\n", formatedOutput(moneyToWithdraw), err)
	} else {
		fmt.Printf("Снятие %s -> %s\n", formatedOutput(moneyToWithdraw), formatedOutput(balanceB()))
	}
}
