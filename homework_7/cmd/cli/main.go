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
	ErrBadFlag        = errors.New("неизвестный флаг")
	ErrArgumentsCount = errors.New("неверное кол-во аргументов")
)

var commandList = []string{
	"list",
	"add",
	"use",
	"delete",
}

func itemList(inv *inventory.Inventory, args []string, stdout, stderr io.Writer) error {
	listCmd := flag.NewFlagSet("list", flag.ContinueOnError)
	listCmd.SetOutput(stderr)

	if err := listCmd.Parse(args); err != nil {
		if errors.Is(err, flag.ErrHelp) {
			return flag.ErrHelp
		}
		return fmt.Errorf("%w: %v", ErrBadFlag, err)
	}
	if listCmd.NArg() > 0 {
		return fmt.Errorf("%w: команда list не принимает аргументов, получено: %v", ErrArgumentsCount, listCmd.Args())
	}

	if len(inv.Items) == 0 {
		fmt.Fprintln(stdout, "Инвентарь пуст.")
		return nil
	}

	fmt.Fprintln(stdout, "--- Ваш инвентарь ---")
	for _, item := range inv.Items {
		fmt.Fprintln(stdout, items.DescribeItem(item))
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
		if errors.Is(err, flag.ErrHelp) {
			return flag.ErrHelp
		}
		return fmt.Errorf("%w: %v", ErrBadFlag, err)
	}
	if addCmd.NArg() > 0 {
		return fmt.Errorf("%w: лишние аргументы: %v", ErrArgumentsCount, addCmd.Args())
	}
	if *itemType == "" {
		return fmt.Errorf("ошибка: %w -type обязателен", ErrFlagIsNotExist)
	}
	if *name == "" {
		return fmt.Errorf("ошибка: %w -name обязателен", ErrFlagIsNotExist)
	}
	if *weight == 0.0 {
		return fmt.Errorf("ошибка: %w -weight обязателен", ErrFlagIsNotExist)
	}
	var item items.Item

	switch *itemType {
	case "weapon":
		item = items.NewWeapon(*name, *weight, *damage, *durability)
	case "armor":
		item = items.NewArmor(*name, *weight, *defense)
	case "potion":
		item = items.NewPotion(*name, *weight, *effect, *charges)
	default:
		return fmt.Errorf("неизвестный тип предмета: %s", *itemType)
	}
	id := inv.AddItem(item)
	if err := inv.Save(); err != nil {
		return fmt.Errorf("не удалось сохранить изменения после добавления: %w", err)
	}
	fmt.Fprintf(stdout, "Успешно добавлен и сохранен предмет %q с ID: %d\n", *name, id)
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
		if errors.Is(err, flag.ErrHelp) {
			return flag.ErrHelp
		}
		return fmt.Errorf("%w: %v", ErrBadFlag, err)
	}
	if useCmd.NArg() > 0 {
		return fmt.Errorf("%w: лишние аргументы: %v", ErrArgumentsCount, useCmd.Args())
	}
	if *itemId == 0 {
		return fmt.Errorf("ошибка: %w -id обязателен", ErrFlagIsNotExist)
	}
	value, ok := inv.FindItemById(*itemId)
	if !ok {
		return fmt.Errorf("предмет с ID: %d не найден", *itemId)
	}
	msg, err := value.Use()
	if err != nil {
		return fmt.Errorf("ошибка при использовании: %w ", err)
	}
	if err = inv.Save(); err != nil {
		return fmt.Errorf("не удалось сохранить изменения после использования: %w", err)
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
		if errors.Is(err, flag.ErrHelp) {
			return flag.ErrHelp
		}
		return fmt.Errorf("%w: %v", ErrBadFlag, err)
	}
	if deleteCmd.NArg() > 0 {
		return fmt.Errorf("%w: лишние аргументы: %v", ErrArgumentsCount, deleteCmd.Args())
	}
	if *itemId == 0 {
		return fmt.Errorf("ошибка: %w -id обязателен", ErrFlagIsNotExist)
	}
	err := inv.RemoveItemById(*itemId)
	if err != nil {
		return err
	}
	if err = inv.Save(); err != nil {
		return fmt.Errorf("не удалось сохранить изменения после удаления: %w", err)
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
	if loadErr := inv.Load(); loadErr != nil && !errors.Is(loadErr, os.ErrNotExist) {
		return fmt.Errorf("не удалось загрузить инвентарь: %w", loadErr)
	}
	switch args[0] {
	case "list":
		return itemList(inv, args[1:], stdout, stderr)
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
	case errors.Is(err, ErrUsage) || errors.Is(err, ErrFlagIsNotExist) || errors.Is(err, ErrBadFlag) || errors.Is(err, ErrArgumentsCount):
		return 2
	default:
		return 1
	}
}

func report(err error, stderr io.Writer) int {
	if err != nil {
		if !errors.Is(err, ErrBadFlag) && !errors.Is(err, flag.ErrHelp) {
			fmt.Fprintln(stderr, err)
		}
	}
	return exitCode(err)
}

func main() {
	os.Exit(report(run(os.Args[1:], os.Stdout, os.Stderr), os.Stderr))
}
