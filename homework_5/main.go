package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

var (
	ErrBroken    = errors.New("предмет сломан")
	ErrExhausted = errors.New("отсутствуют заряды зелья")
)

type Item interface {
	Name() string
	Weight() float64
	Use() (string, error)
	SetId(id int)
	GetId() int
}

type Storable interface {
	Serialize(w io.Writer) error
	Deserialize(r io.Reader) error
}

type Breakable interface {
	IsBroken() bool
}

type ParseItemError struct {
	Line  int
	Value string
	Err   error
}

func (e *ParseItemError) Error() string {
	return fmt.Sprintf("ошибка в строке %d (данные: %q): %v", e.Line, e.Value, e.Err)
}

func (e *ParseItemError) Unwrap() error {
	return e.Err
}

type Weapon struct {
	ID         int
	name       string
	weight     float64
	Damage     int
	Durability int
}

func (w *Weapon) IsBroken() bool {
	return w.Durability <= 0
}
func (w *Weapon) Name() string {
	return w.name
}
func (w *Weapon) Weight() float64 {
	return w.weight
}
func (w *Weapon) Use() (string, error) {
	if w.IsBroken() {
		return "", ErrBroken
	}
	w.Durability--
	return "Использован " + w.name, nil
}

func (w *Weapon) SetId(id int) {
	w.ID = id
}

func (w *Weapon) GetId() int {
	return w.ID
}

func (w *Weapon) Serialize(writer io.Writer) error {
	_, err := fmt.Fprintf(writer, "%d|Weapon|%s|%.2f|%d|%d\n", w.ID, w.name, w.weight, w.Damage, w.Durability)
	if err != nil {
		return fmt.Errorf("не удалось записать оружие: %q %w", w.name, err)
	}
	return nil
}

func (w *Weapon) Deserialize(r io.Reader) error {
	data, err := io.ReadAll(r)
	if err != nil {
		return err
	}
	line := string(data)
	lineSlice := strings.Split(line, "|")
	if len(lineSlice) != 6 {
		return fmt.Errorf("ожидалось 6 полей, получено %d", len(lineSlice))
	}
	w.ID, err = strconv.Atoi(lineSlice[0])
	if err != nil {
		return fmt.Errorf("неверный формат для ID: %w", err)
	}
	w.name = lineSlice[2]
	weight, err := strconv.ParseFloat(lineSlice[3], 64)
	if err != nil {
		return fmt.Errorf("неверный формат для веса: %w", err)
	}
	w.weight = weight
	damage, err := strconv.Atoi(lineSlice[4])
	if err != nil {
		return fmt.Errorf("неверный формат для урона: %w", err)
	}
	w.Damage = damage
	durability, err := strconv.Atoi(lineSlice[5])
	if err != nil {
		return fmt.Errorf("неверный формат для прочности: %w", err)
	}
	w.Durability = durability
	return nil
}

var _ Item = (*Weapon)(nil)

type Armor struct {
	ID      int
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
func (a *Armor) Use() (string, error) {
	return "Экипировать " + a.name, nil
}

func (a *Armor) SetId(id int) {
	a.ID = id
}

func (a *Armor) GetId() int {
	return a.ID
}

func (a *Armor) Serialize(writer io.Writer) error {
	_, err := fmt.Fprintf(writer, "%d|Armor|%s|%.2f|%d\n", a.ID, a.name, a.weight, a.Defense)
	if err != nil {
		return fmt.Errorf("не удалось записать броню: %q, %w", a.name, err)
	}
	return nil
}

func (a *Armor) Deserialize(r io.Reader) error {
	data, err := io.ReadAll(r)
	if err != nil {
		return err
	}
	line := string(data)
	lineSlice := strings.Split(line, "|")
	if len(lineSlice) != 5 {
		return fmt.Errorf("ожидалось 4 поля, получено %d", len(lineSlice))
	}
	a.ID, err = strconv.Atoi(lineSlice[0])
	if err != nil {
		return fmt.Errorf("неверный формат для ID: %w", err)
	}
	a.name = lineSlice[2]
	weight, err := strconv.ParseFloat(lineSlice[3], 64)
	if err != nil {
		return fmt.Errorf("неверный формат для веса: %w", err)
	}
	a.weight = weight
	defense, err := strconv.Atoi(lineSlice[4])
	if err != nil {
		return fmt.Errorf("неверный формат защиты: %w", err)
	}
	a.Defense = defense
	return nil
}

var _ Item = (*Armor)(nil)

type Potion struct {
	ID      int
	name    string
	weight  float64
	Effect  string
	Charges int
}

func (p *Potion) IsBroken() bool {
	return p.Charges <= 0
}
func (p *Potion) Name() string {
	return p.name
}
func (p *Potion) Weight() float64 {
	return p.weight
}
func (p *Potion) Use() (string, error) {
	if p.IsBroken() {
		return "", ErrExhausted
	}
	p.Charges--
	return "Использовано " + p.name, nil
}

func (p *Potion) SetId(id int) {
	p.ID = id
}

func (p *Potion) GetId() int {
	return p.ID
}

func (p *Potion) Serialize(writer io.Writer) error {
	_, err := fmt.Fprintf(writer, "%d|Potion|%s|%.2f|%s|%d\n", p.ID, p.name, p.weight, p.Effect, p.Charges)
	if err != nil {
		return fmt.Errorf("не удалось записать зелье: %q %w", p.name, err)
	}
	return nil
}

func (p *Potion) Deserialize(r io.Reader) error {
	data, err := io.ReadAll(r)
	if err != nil {
		return err
	}
	line := string(data)
	lineSlice := strings.Split(line, "|")

	if len(lineSlice) != 6 {
		return fmt.Errorf("ожидалось 6 полей, получено %d", len(lineSlice))
	}
	p.ID, err = strconv.Atoi(lineSlice[0])
	if err != nil {
		return fmt.Errorf("неверный формат для ID: %w", err)
	}
	p.name = lineSlice[2]
	weight, err := strconv.ParseFloat(lineSlice[3], 64)
	if err != nil {
		return fmt.Errorf("неверный формат для веса: %w", err)
	}
	p.weight = weight
	p.Effect = lineSlice[4]
	charges, err := strconv.Atoi(lineSlice[5])
	if err != nil {
		return fmt.Errorf("неверный формат для зарядов: %w", err)
	}
	p.Charges = charges
	return nil
}

var _ Item = (*Potion)(nil)

func DescribeItem(i Item) string {
	switch value := i.(type) {
	case nil:
		return "Предмет отсутствует"
	case *Weapon:
		if value == nil {
			return "Предмет отсутствует"
		}
	case *Armor:
		if value == nil {
			return "Предмет отсутствует"
		}
	case *Potion:
		if value == nil {
			return "Предмет отсутствует"
		}
	default:
		return "Неизвестный предмет"
	}
	return fmt.Sprintf("%s (вес: %.2f)", i.Name(), i.Weight())
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
	nextID int
	Items  []Item
}

func (inv *Inventory) AddItem(item Item) int {
	item.SetId(inv.nextID)
	inv.Items = append(inv.Items, item)
	inv.nextID++
	return item.GetId()
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
		switch value := item.(type) {
		case *Weapon:
			return value.Durability <= 0
		case *Potion:
			return value.Charges <= 0
		default:
			return false
		}
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

func (inv *Inventory) Save(w io.Writer) error {
	var errs []error
	for i, item := range inv.Items {
		if storableItem, ok := item.(Storable); ok {
			if err := storableItem.Serialize(w); err != nil {
				errs = append(errs, fmt.Errorf("слот %d (предмет %q): %w", i, item.Name(), err))
			}
		}
	}
	return errors.Join(errs...)
}

func (inv *Inventory) Load(r io.Reader) error {
	var errs []error
	scanner := bufio.NewScanner(r)
	lineNum := 0
	for scanner.Scan() {
		lineNum++
		line := scanner.Text()
		lineSlice := strings.Split(line, "|")

		if len(lineSlice) < 2 {
			err := fmt.Errorf("недостаточно данных: строка не содержит тип предмета")
			parseErr := &ParseItemError{
				Line:  lineNum,
				Value: line,
				Err:   err,
			}
			errs = append(errs, parseErr)
			continue
		}
		switch lineSlice[1] {
		case "Weapon":
			w := &Weapon{}
			err := w.Deserialize(strings.NewReader(line))
			if err != nil {
				parseErr := &ParseItemError{
					Line:  lineNum,
					Value: line,
					Err:   err,
				}
				errs = append(errs, parseErr)
			} else {
				inv.Items = append(inv.Items, w)
				if w.GetId() >= inv.nextID {
					inv.nextID = w.ID + 1
				}
			}
		case "Armor":
			a := &Armor{}
			err := a.Deserialize(strings.NewReader(line))
			if err != nil {
				parseErr := &ParseItemError{
					Line:  lineNum,
					Value: line,
					Err:   err,
				}
				errs = append(errs, parseErr)
			} else {
				inv.Items = append(inv.Items, a)
				if a.GetId() >= inv.nextID {
					inv.nextID = a.ID + 1
				}
			}
		case "Potion":
			p := &Potion{}
			err := p.Deserialize(strings.NewReader(line))
			if err != nil {
				parseErr := &ParseItemError{
					Line:  lineNum,
					Value: line,
					Err:   err,
				}
				errs = append(errs, parseErr)
			} else {
				inv.Items = append(inv.Items, p)
				if p.GetId() >= inv.nextID {
					inv.nextID = p.ID + 1
				}
			}
		default:
			err := fmt.Errorf("неизвестный тип предмета %s", lineSlice[1])
			parseErr := &ParseItemError{
				Line:  lineNum,
				Value: line,
				Err:   err,
			}
			errs = append(errs, parseErr)
		}
	}
	if err := scanner.Err(); err != nil {
		errs = append(errs, fmt.Errorf("ошибка чтения из файла: %w", err))
	}
	return errors.Join(errs...)
}

func saveToFile(inventory *Inventory, fileName string) (err error) {
	file, err := os.Create(fileName)
	if err != nil {
		return fmt.Errorf("ошибка сохранения в файл: %s", err)
	}
	defer func() {
		closeErr := file.Close()
		if closeErr != nil {
			err = errors.Join(err, fmt.Errorf("ошибка при закрытии файла: %w", closeErr))
		}
	}()
	return inventory.Save(file)
}

func loadFromFile(inventory *Inventory, fileName string) error {
	file, err := os.Open(fileName)
	if err != nil {
		return fmt.Errorf("ошибка загрузки из файла: %w", err)
	}
	defer file.Close()
	return inventory.Load(file)
}

func main() {
	inventory := Inventory{
		nextID: 1,
	}

	sword := &Weapon{name: "Меч", weight: 5.0, Damage: 10, Durability: 5}
	armor := &Armor{name: "Щит", weight: 10.0, Defense: 5}
	healPotion := &Potion{name: "Лечебное зелье", weight: 0.5, Effect: "+50 Hp", Charges: 3}
	bow := &Weapon{name: "Сломанный лук", weight: 3.5, Damage: 1, Durability: 0}

	inventory.AddItem(sword)
	inventory.AddItem(armor)
	inventory.AddItem(healPotion)
	inventory.AddItem(bow)
	
	fmt.Println("---До Use---")
	fmt.Printf("Состояние меча: %d, Заряды зелья: %d\n", sword.Durability, healPotion.Charges)
	sword.Use()
	healPotion.Use()
	fmt.Println("---После Use---")
	fmt.Printf("Состояние меча: %d, Заряды зелья: %d\n", sword.Durability, healPotion.Charges)
	fmt.Println("---Use предмета с разделением ошибки через ErrorIs---")
	res, err := bow.Use()
	if errors.Is(err, ErrBroken) {
		fmt.Println("Предмет сломан")
	} else if err != nil {
		fmt.Println("Непредвиденная ошибка: ", err)
	} else {
		fmt.Println(res)
	}

	newInv := Inventory{}

	fmt.Println("---Файла ещё нет---")
	err = loadFromFile(&newInv, "fileNotExist.txt")
	if errors.Is(err, os.ErrNotExist) {
		fmt.Println("Отсутствует фаил", err)
	}

	fmt.Println("---Нет прав на чтение или запись---")
	err = loadFromFile(&newInv, "C:\\\\Windows\\\\System32\\\\config\\\\SAM")
	if errors.Is(err, os.ErrPermission) {
		fmt.Println("Отсутствие прав на чтение/запись", err)
	}

	fmt.Println("---Запись в Файл c одной битой строкой/Загрузка пустого файла---")
	newInv = Inventory{}
	err = loadFromFile(&newInv, "inventory.txt")
	if err != nil {
		var parseItemError *ParseItemError
		if errors.As(err, &parseItemError) {
			fmt.Println("Внимание найдена ошибка")
			fmt.Println(parseItemError)
			fmt.Println("Успешно было загружено предметов:", len(newInv.Items))
		} else {
			fmt.Println("Непредвиденная ошибка: ", err)
		}
	} else {
		if len(newInv.Items) > 0 {
			fmt.Println("Успешная загрузка, загружено: ", len(newInv.Items))
		} else {
			fmt.Println("Успешная загрузка. Инвентарь пуст.")
		}
	}

	// Ошибка новое содержимое короче старого решается через os.Create
}
