package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

type Item interface {
	Name() string
	Weight() float64
	Use() string
}

type Storable interface {
	Serialize(w io.Writer)
	Deserialize(r io.Reader)
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

func (w *Weapon) Serialize(writer io.Writer) {
	_, err := fmt.Fprintf(writer, "Weapon|%s|%.2f|%d|%d\n", w.name, w.weight, w.Damage, w.Durability)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Ошибка записи: %s", err)
	}
}

func (w *Weapon) Deserialize(r io.Reader) {
	scanner := bufio.NewScanner(r)
	if scanner.Scan() {
		line := scanner.Text()
		lineSlice := strings.Split(line, "|")

		if len(lineSlice) != 5 {
			return
		}

		w.name = lineSlice[1]
		w.weight, _ = strconv.ParseFloat(lineSlice[2], 64)
		w.Damage, _ = strconv.Atoi(lineSlice[3])
		w.Durability, _ = strconv.Atoi(lineSlice[4])
	}
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

func (a *Armor) Serialize(writer io.Writer) {
	_, err := fmt.Fprintf(writer, "Armor|%s|%.2f|%d\n", a.name, a.weight, a.Defense)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Ошибка записи: %s", err)
	}
}

func (a *Armor) Deserialize(r io.Reader) {
	scanner := bufio.NewScanner(r)
	if scanner.Scan() {
		line := scanner.Text()
		lineSlice := strings.Split(line, "|")
		if len(lineSlice) != 4 {
			return
		}
		a.name = lineSlice[1]
		a.weight, _ = strconv.ParseFloat(lineSlice[2], 64)
		a.Defense, _ = strconv.Atoi(lineSlice[3])
	}
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

func (inv *Inventory) Save(w io.Writer) {
	for _, item := range inv.Items {
		if storableItem, ok := item.(Storable); ok {
			storableItem.Serialize(w)
		}
	}
}

func (inv *Inventory) Load(r io.Reader) {
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		line := scanner.Text()
		lineSlice := strings.Split(line, "|")
		switch lineSlice[0] {
		case "Weapon":
			w := &Weapon{}
			w.Deserialize(strings.NewReader(line))
			inv.AddItem(w)

		case "Armor":
			a := &Armor{}
			a.Deserialize(strings.NewReader(line))
			inv.AddItem(a)
		}
	}
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

	fmt.Println("---Сохранение в файл---")
	file, _ := os.Create("inventory.txt")
	inventory.Save(file)
	defer file.Close()

	fmt.Println("---Чтение из файла---")
	newInv := &Inventory{}
	file2, _ := os.Open("inventory.txt")
	newInv.Load(file2)
	defer file.Close()

	newItems := newInv.ItemNames()
	for _, newItem := range newItems {
		fmt.Printf("%s\n", newItem)
	}
}
