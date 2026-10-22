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
	accountBalance := initial * subUnitAmount
	deposit = func(initial int64) {
		initialSubUnits := initial * subUnitAmount
		accountBalance += initialSubUnits
		fmt.Printf("Пополнение %s -> %s\n", formatedOutput(initialSubUnits), formatedOutput(accountBalance))
	}
	withDraw = func(initial int64) error {
		initialSubUnits := initial * subUnitAmount
		if accountBalance < initialSubUnits {
			err := fmt.Errorf("Снятие %s: недостаточно средств (баланс %s)\n", formatedOutput(initialSubUnits), formatedOutput(accountBalance))
			return err
		}
		accountBalance -= initialSubUnits
		fmt.Printf("Снятие %s -> %s\n", formatedOutput(initial), formatedOutput(accountBalance))
		return nil
	}
	balance = func() int64 {
		fmt.Printf("Счёт %s\n", formatedOutput(accountBalance))
		return accountBalance
	}
	return deposit, withDraw, balance
}

func main() {
	fmt.Println("-----Счет А-----")
	depositA, withDrawA, balanceA := NewAccount(1000)
	balanceA()
	depositA(10)
	err := withDrawA(10000)
	if err != nil {
		fmt.Print(err)
	}
	balanceA()
	/*
		fmt.Println("-----Счет Б-----")
		depositB, withDrawB, balanceB := NewAccount(500)
		balanceB()
		depositB(100)
		err = withDrawB(300)
		if err != nil {
			fmt.Print(err)
		}
		balanceB()
	*/
}
