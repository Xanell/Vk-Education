package main

import "fmt"

/*
Task
Выведите таблицу перевода из Цельсия в Фаренгейт от -20 до 40 °C с шагом 10. Формула: F = C * 9 / 5 + 32.
Рядом подпись через `switch` без выражения: ниже 0 `мороз`, ниже 15 `холодно`, ниже 25 `тепло`, иначе `жарко`.
Пример:
-20°C =  -4.0°F  мороз
-10°C =  14.0°F  мороз
0°C =  32.0°F  холодно
10°C =  50.0°F  холодно
20°C =  68.0°F  тепло
30°C =  86.0°F  жарко
40°C = 104.0°F  жарко
*/

func main() {
	for celsius := -20; celsius <= 40; celsius += 10 {
		fahrenheit := float64(celsius)*1.8 + 32
		switch {
		case celsius < 0:
			fmt.Printf("%d°C = %.1f°F мороз\n", celsius, fahrenheit)
		case celsius < 15:
			fmt.Printf("%d°C = %.1f°F холодно\n", celsius, fahrenheit)
		case celsius < 25:
			fmt.Printf("%d°C = %.1f°F тепло\n", celsius, fahrenheit)
		default:
			fmt.Printf("%d°C = %.1f°F жарко\n", celsius, fahrenheit)
		}
	}
}
