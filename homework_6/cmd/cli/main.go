package main

import (
	"errors"
	"flag"
	"fmt"
	"go-education/homework_6/internal/inventory"
	"go-education/homework_6/internal/items"
	"go-education/homework_6/internal/storage"
	"io"
	"os"
)

var (
	ErrUsage          = errors.New("неизвестная команда")
	ErrFlagIsNotExist = errors.New("отсутствует флаг")
)

var commandList = []string{
	"list",
	"add",
	"use",
}

func itemList(inv *inventory.Inventory, stdout io.Writer) error {
	if len(inv.Items) == 0 {
		fmt.Fprintln(stdout, "Инвентарь пуст.")
		return nil
	}

	fmt.Fprintln(stdout, "--- Ваш инвентарь ---")
	for _, item := range inv.Items {
		fmt.Fprintln(stdout, domain.DescribeItem(item))
	}
	return nil
}

func addItem(inv *inventory.Inventory, args []string, stdout, stderr io.Writer) error {
	addCmd := flag.NewFlagSet("add", flag.ContinueOnError)
	itemType := addCmd.String("type", "", "Тип предмета (weapon, armor, potion)")
	name := addCmd.String("name", "", "Имя предмета")
	weight := addCmd.Float64("weight", 0.0, "Вес")

	damage := addCmd.Int("damage", 10, "Урон (для weapon)")
	durability := addCmd.Int("durability", 100, "Прочность (для weapon)")
	defense := addCmd.Int("defense", 5, "Защита (для armor)")
	effect := addCmd.String("effect", "", "Эффект (для potion)")
	charges := addCmd.Int("charges", 1, "Заряды (для potion)")
	addCmd.SetOutput(stderr)
	addCmd.Usage = func() {
		fmt.Fprintln(stderr, "для использования: add [-type ItemType] [-name NAME] [-weight WEIGHT]")
		addCmd.PrintDefaults()
	}
	if err := addCmd.Parse(args); err != nil {
		return err
	}
	if *itemType == "" {
		return fmt.Errorf("ошибка: %w флаг -type обязателен", ErrFlagIsNotExist)
	}
	if *name == "" {
		return fmt.Errorf("ошибка: %w флаг -name обязателен", ErrFlagIsNotExist)
	}
	if *weight == 0.0 {
		return fmt.Errorf("ошибка: %w флаг -weight обязателен", ErrFlagIsNotExist)
	}
	var item domain.Item

	switch *itemType {
	case "weapon":
		item = domain.NewWeapon(*name, *weight, *damage, *durability)
	case "armor":
		item = domain.NewArmor(*name, *weight, *defense)
	case "potion":
		item = domain.NewPotion(*name, *weight, *effect, *charges)
	default:
		return fmt.Errorf("неизвестный тип предмета: %s", *itemType)
	}
	id := inv.AddItem(item)
	fmt.Fprintf(stdout, "Успешно добавлен предмет %q с ID: %d\n", *name, id)
	return nil
}

func useItem(inv *inventory.Inventory, args []string, stdout, stderr io.Writer) error {
	useCmd := flag.NewFlagSet("use", flag.ContinueOnError)
	itemId := useCmd.Int("id", 0, "ID предмета для использования")
	useCmd.SetOutput(stderr)
	useCmd.Usage = func() {
		fmt.Fprintln(stderr, "для использования: use [-id ID]")
		useCmd.PrintDefaults()
	}
	if err := useCmd.Parse(args); err != nil {
		return err
	}
	if *itemId == 0 {
		return fmt.Errorf("ошибка: %w флаг -id обязателен", ErrFlagIsNotExist)
	}
	value, ok := inv.FindItemById(*itemId)
	if !ok {
		return fmt.Errorf("предмет с ID: %d не найден", *itemId)
	}
	msg, err := value.Use()
	if err != nil {
		return err
	}
	fmt.Fprintln(stdout, msg)
	return nil
}

func deleteItem(inv *inventory.Inventory, args []string, stdout, stderr io.Writer) error {
	deleteCmd := flag.NewFlagSet("delete", flag.ContinueOnError)
	itemId := deleteCmd.Int("id", 0, "ID предмета для удаления")
	deleteCmd.SetOutput(stderr)
	deleteCmd.Usage = func() {
		fmt.Fprintln(stderr, "для использования: delete [-id ID]")
		deleteCmd.PrintDefaults()
	}
	if err := deleteCmd.Parse(args); err != nil {
		return err
	}
	if *itemId == 0 {
		return fmt.Errorf("ошибка: %w флаг -id обязателен", ErrFlagIsNotExist)
	}
	err := inv.RemoveItemById(*itemId)
	if err != nil {
		return err
	}
	fmt.Fprintf(stdout, "Предмет с ID %d успешно удален из инвентаря.\n", *itemId)
	return nil
}

func run(args []string, stdout, stderr io.Writer) error {
	if len(args) == 0 {
		return fmt.Errorf("%w: список команд: %s", ErrUsage, commandList)
	}
	fs := storage.NewFileStorage("inventory.txt")
	inv := inventory.NewInventory(fs)
	_ = inv.Load()
	defer func() {
		if err := inv.Save(); err != nil {
			fmt.Fprintln(stderr, "Ошибка автосохранения:", err)
		}
	}()
	switch args[0] {
	case "list":
		return itemList(inv, stdout)
	case "add":
		return addItem(inv, args[1:], stdout, stderr)
	case "use":
		return useItem(inv, args[1:], stdout, stderr)
	case "delete":
		return deleteItem(inv, args[1:], stdout, stderr)
	default:
		return fmt.Errorf("%w: неизвестная команда %s", ErrUsage, args[0])
	}
}

func exitCode(err error) int {
	switch {
	case err == nil:
		return 0
	case errors.Is(err, ErrUsage) || errors.Is(err, ErrFlagIsNotExist):
		return 2
	default:
		return 1
	}
}

func report(err error, stderr io.Writer) int {
	if err != nil {
		fmt.Fprintln(stderr, err)
	}
	return exitCode(err)
}

func main() {
	os.Exit(report(run(os.Args[1:], os.Stdout, os.Stderr), os.Stderr))
}
