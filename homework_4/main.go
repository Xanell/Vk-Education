package main

import (
	"fmt"
	"io"
)

type Item interface {
	Name() string
	Weight() float64
	Use() string
}

type Weapon struct {
	name       string
	weight     float64
	Damage     int
	Durability int
}

func (w *Weapon) Name() string {
	return w.name
}
func (w *Weapon) Weight() float64 {
	return w.weight
}
func (w *Weapon) Use() string {
	return "Использован " + w.name
}

var _ Item = (*Weapon)(nil)

type Armor struct {
	name    string
	weight  float64
	Defense int
}

func (a *Armor) Name() string {
	return a.name
}

func (a *Armor) Weight() float64 {
	return a.weight
}

func (a *Armor) Use() string {
	return "Экипировать " + a.name
}

var _ Item = (*Armor)(nil)

type Potion struct {
	name   string
	weight float64
	Effect string
	Charge int
}

func (p *Potion) Name() string {
	return p.name
}
func (p *Potion) Weight() float64 {
	return p.weight
}
func (p *Potion) Use() string {
	return "Использовано " + p.name
}

var _ Item = (*Potion)(nil)

func DescribeItem(i Item) string {
	if i == nil {
		return "Предмет отсутствует"
	}
	name := i.Name()
	weight := i.Weight()
	return fmt.Sprintf("%s (вес: %.2f)", name, weight)
}

func Filter[T any](items []T, predicate func(T) bool) []T {
	result := make([]T, 0, len(items))
	for _, item := range items {
		if predicate(item) {
			result = append(result, item)
		}
	}
	return result
}

func Map[T any, R any](items []T, transform func(T) R) []R {
	result := make([]R, len(items))
	for index, item := range items {
		result[index] = transform(item)
	}
	return result
}

func Find[T any](items []T, condition func(T) bool) (T, bool) {
	for _, item := range items {
		if condition(item) {
			return item, true
		}
	}
	var zero T
	return zero, false
}

type Inventory struct {
	Items []Item
}

func (inv *Inventory) AddItem(item Item) {
	inv.Items = append(inv.Items, item)
}

func (inv *Inventory) Weapons() []*Weapon {
	weapons := Filter(inv.Items, func(item Item) bool {
		_, ok := item.(*Weapon)
		return ok
	})
	return Map(weapons, func(item Item) *Weapon {
		w, ok := item.(*Weapon)
		if !ok {
			return nil
		}
		return w
	})
}

func (inv *Inventory) BrokenItems() []Item {
	return Filter(inv.Items, func(item Item) bool {
		w, ok := item.(*Weapon)
		if !ok {
			return false
		}
		if w.Durability <= 0 {
			return true
		}
		p, ok := item.(*Potion)
		if !ok {
			return false
		}
		if p.Charge <= 0 {
			return true
		}
		return false
	})
}

func (inv *Inventory) ItemNames() []string {
	return Map(inv.Items, func(item Item) string {
		return item.Name()
	})
}

func (inv *Inventory) FindItemByName(name string) (Item, bool) {
	return Find(inv.Items, func(item Item) bool {
		return item.Name() == name
	})
}

// TODO: Бонус: реализуйте интефейс Storable для Weapon и Armor:
// TODO: - Weapon: формат "Weapon|Name|Damage|Durability|Weight"
// TODO: - Armor: формат "Armor|Name|Defense|Weight"

type Storable interface {
	Serialize(w io.Writer)
	Deserialize(r io.Reader)
}

func (inv *Inventory) Save(w io.Writer) {
	// TODO: Бонус: сделайте сохранение/загрузку инвентаря в/из файла
}

func (inv *Inventory) Load(r io.Reader) {
	// TODO: Бонус: сделайте сохранение/загрузку инвентаря в/из файла
}

func main() {
	inventory := Inventory{}

	sword := &Weapon{name: "Меч", weight: 5.0, Damage: 10, Durability: 5}
	armor := &Armor{name: "Щит", weight: 10.0, Defense: 5}
	healPotion := &Potion{name: "Лечебное зелье", weight: 0.5, Effect: "+50 Hp", Charge: 3}
	bow := &Weapon{name: "Сломанный лук", weight: 3.5, Damage: 1, Durability: 0}

	inventory.AddItem(sword)
	inventory.AddItem(armor)
	inventory.AddItem(healPotion)
	inventory.AddItem(bow)

	fmt.Println("---Содержимое инвентаря---")
	for _, item := range inventory.Items {
		fmt.Printf("%s\n", DescribeItem(item))
	}

	fmt.Println("---Use предмета с выводом в консоль---")
	fmt.Println(sword.Use())
	fmt.Println(healPotion.Use())

	fmt.Println("---DescribeItem с предметом---")
	fmt.Println(DescribeItem(sword))

	fmt.Println("---DescribeItem c nil---")
	sword2 := &Weapon{}
	fmt.Println(DescribeItem(sword2))

	fmt.Println("---DescribeItem с nil-указателем")
	//var nilSword *Weapon
	//fmt.Println(DescribeItem(nilSword))

	fmt.Println("--Результат вывода Weapons---")
	weapons := inventory.Weapons()
	for _, weapon := range weapons {
		fmt.Printf("Название: %s Урон: %d Вес: %.2f Прочность: %d\n", weapon.name, weapon.Damage, weapon.weight, weapon.Durability)
	}

	fmt.Println("---Результат вывода BrokenItems---")
	brokenItems := inventory.BrokenItems()
	for _, brokenItem := range brokenItems {
		fmt.Printf("%s\n", DescribeItem(brokenItem))
	}

	fmt.Println("---Результат ItemNames---")
	itemNames := inventory.ItemNames()
	for _, itemName := range itemNames {
		fmt.Printf("%s\n", itemName)
	}

	fmt.Println("---Результат вывода FindItemByName---")
	item, found := inventory.FindItemByName("Щит")
	if !found {
		fmt.Println("Предмет не найден")
	} else {
		fmt.Printf(DescribeItem(item))
	}

	// TODO: Бонус: сделайте сохранение инвентаря в файл и загрузку инвентаря из файла
}
