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

type customerStatus int

const (
	regular customerStatus = iota
	silver
	gold
)
const (
	deliveryCost          = 299
	freeDeliveryThreshold = 3000
)

func main() {
	price := 4500.00
	status := gold
	textStatus := ""

	var delivery float64
	var discountRate float64

	switch status {
	case regular:
		textStatus = "Regular"
	case silver:
		discountRate = 0.05
		textStatus = "Silver"
	case gold:
		discountRate = 0.1
		textStatus = "Gold"
	default:
		textStatus = "Неизвестный статус"
	}
	discount := price * discountRate
	discountPrice := price - discount

	if discountPrice < freeDeliveryThreshold {
		delivery = deliveryCost
	}

	totalPrice := discountPrice + delivery

	fmt.Printf("Сумма %.2f, статус: %s\n", price, textStatus)
	fmt.Printf("Скидка: %.2f\n", discount)
	fmt.Printf("Доставка: %.2f\n", delivery)
	fmt.Printf("Итого: %.2f", totalPrice)
}
