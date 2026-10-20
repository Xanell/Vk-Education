package main

import "fmt"

/*
Task
Даны сумма заказа и статус покупателя. Посчитайте скидку, стоимость доставки и итоговую сумму к оплате, выведите их как в примерах.
Правила:
- Статус покупателя: тип `CustomerStatus int`, константы `Regular`, `Silver`, `Gold` через `iota`.
- Скидка через `switch`: `Regular` 0%, `Silver` 5%, `Gold` 10%. Неизвестный статус: вывести `неизвестный статус`, считать без скидки.
- Доставка 299 ₽, бесплатно, если сумма после скидки от 3000 ₽. Порог и цену доставки вынести в константы.
- Сумма `float64`, вывод через `%.2f`.
Пример:
Сумма: 4500.00, статус: Gold
Скидка: 450.00
Доставка: 0.00
Итого: 4050.00
*/

type CustomerStatus int

const (
	Regular CustomerStatus = iota
	Silver
	Gold
)
const (
	DeliveryCost  = 299
	DeliveryPrice = 3000 // DeliveryPrice Как лучше назвать переменную обозначающую порог?
)

func main() {
	price := 2000.00
	status := Regular
	textStatus := ""
	
	var delivery float64
	var discountRate float64

	switch status {
	case Regular:
		discountRate = 0
		textStatus = "Regular"
	case Silver:
		discountRate = 0.05
		textStatus = "Silver"
	case Gold:
		discountRate = 0.1
		textStatus = "Gold"
	default:
		discountRate = 0
		textStatus = "Неизвестный статус"
	}
	discount := price * discountRate
	discountPrice := price - discount

	if discountPrice < DeliveryPrice {
		delivery = DeliveryCost
	}

	totalPrice := discountPrice + delivery
	fmt.Printf("Сумма %.2f, статус: %s\nСкидка: %.2f\nДоставка: %.2f\nИтого: %.2f", price, textStatus, discountRate, delivery, totalPrice)
}
