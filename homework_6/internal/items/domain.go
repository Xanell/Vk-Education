package domain

import (
	"errors"
	"fmt"
	"io"
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

func NewParseItemError(line int, value string, err error) *ParseItemError {
	return &ParseItemError{Line: line, Value: value, Err: err}
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

func NewWeapon(name string, weight float64, damage, durability int) *Weapon {
	return &Weapon{
		name:       name,
		weight:     weight,
		Damage:     damage,
		Durability: durability,
	}
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
		return fmt.Errorf("не удалось прочитать данные из потока: %w", err)
	}
	line := string(data)
	lineSlice := strings.Split(line, "|")
	w.ID, err = parseInt(lineSlice[0], "ID")
	if err != nil {
		return err
	}
	w.name = lineSlice[2]
	weight, err := parseFloat(lineSlice[3], "вес")
	if err != nil {
		return err
	}
	w.weight = weight
	damage, err := parseInt(lineSlice[4], "урон")
	if err != nil {
		return err
	}
	w.Damage = damage
	durability, err := parseInt(lineSlice[5], "прочность")
	if err != nil {
		return err
	}
	w.Durability = durability
	return nil
}

var _ Item = (*Weapon)(nil)
var _ Storable = (*Weapon)(nil)

type Armor struct {
	ID      int
	name    string
	weight  float64
	Defense int
}

func NewArmor(name string, weight float64, defense int) *Armor {
	return &Armor{
		name:    name,
		weight:  weight,
		Defense: defense,
	}
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
		return fmt.Errorf("не удалось прочитать данные из потока: %w", err)
	}
	line := string(data)
	lineSlice := strings.Split(line, "|")
	if len(lineSlice) != 5 {
		return fmt.Errorf("ожидалось 5 полей, получено %d", len(lineSlice))
	}
	a.ID, err = parseInt(lineSlice[0], "ID")
	if err != nil {
		return err
	}
	a.name = lineSlice[2]
	weight, err := parseFloat(lineSlice[3], "вес")
	if err != nil {
		return err
	}
	a.weight = weight
	defense, err := parseInt(lineSlice[4], "защита")
	if err != nil {
		return err
	}
	a.Defense = defense
	return nil
}

var _ Item = (*Armor)(nil)
var _ Storable = (*Armor)(nil)

type Potion struct {
	ID      int
	name    string
	weight  float64
	Effect  string
	Charges int
}

func NewPotion(name string, weight float64, effect string, charges int) *Potion {
	return &Potion{
		name:    name,
		weight:  weight,
		Effect:  effect,
		Charges: charges,
	}
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
	p.ID, err = parseInt(lineSlice[0], "ID")
	if err != nil {
		return err
	}
	p.name = lineSlice[2]
	weight, err := parseFloat(lineSlice[3], "вес")
	if err != nil {
		return err
	}
	p.weight = weight
	p.Effect = lineSlice[4]
	charges, err := parseInt(lineSlice[5], "заряд")
	if err != nil {
		return err
	}
	p.Charges = charges
	return nil
}

var _ Item = (*Potion)(nil)
var _ Storable = (*Potion)(nil)

type RingOfPanic struct {
	ID          int
	name        string
	weight      float64
	Description string
}

func NewRingOfPanic(name string, weight float64, description string) *RingOfPanic {
	return &RingOfPanic{
		name:        name,
		weight:      weight,
		Description: description,
	}
}
func (r *RingOfPanic) Name() string {
	return r.name
}

func (r *RingOfPanic) Weight() float64 {
	return r.weight
}

func (r *RingOfPanic) Use() (string, error) {
	panic("Кольцо паники сработало!")
}

func (r *RingOfPanic) SetId(id int) {
	r.ID = id
}
func (r *RingOfPanic) GetId() int {
	return r.ID
}

var _ Item = (*RingOfPanic)(nil)

func parseInt(s, fieldName string) (int, error) {
	value, err := strconv.Atoi(s)
	if err != nil {
		return 0, fmt.Errorf("неверный формат для %s: %w", fieldName, err)
	}
	return value, nil
}

func parseFloat(s, fieldName string) (float64, error) {
	value, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return 0.0, fmt.Errorf("неверный формат для %s: %w", fieldName, err)
	}
	return value, nil
}
