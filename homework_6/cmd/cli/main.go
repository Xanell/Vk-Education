package main

import (
	"errors"
	"flag"
	"fmt"
	"go-education/homework_6/internal/inventory"
	domain "go-education/homework_6/internal/items"
	"go-education/homework_6/internal/storage"
	"go-education/homework_6/pkg/utils"
	"os"
)

var ErrUsage = errors.New("неизвестная команда")
var commandList = []string{
	"list",
	"add",
	"use",
}

func itemList(inv *inventoryService.Inventory) error {
	if len(inv.Items) == 0 {
		fmt.Println("Инвентарь пуст.")
		return nil
	}

	fmt.Println("--- Ваш инвентарь ---")
	for _, item := range inv.Items {
		fmt.Println(utils.DescribeItem(item))
	}
	return nil
}

func addItem(inv *inventoryService.Inventory, args []string) error {
	addCmd := flag.NewFlagSet("add", flag.ContinueOnError)
	itemType := addCmd.String("type", "", "Тип предмета (weapon, armor, potion)")
	name := addCmd.String("name", "", "Имя предмета")
	weight := addCmd.Float64("weight", 5.0, "Вес")

	damage := addCmd.Int("damage", 10, "Урон (для weapon)")
	durability := addCmd.Int("durability", 100, "Прочность (для weapon)")
	defense := addCmd.Int("defense", 5, "Защита (для armor)")
	effect := addCmd.String("effect", "", "Эффект (для potion)")
	charges := addCmd.Int("charges", 1, "Заряды (для potion)")
	if err := addCmd.Parse(args); err != nil {
		return err
	}
	if *name == "" {
		return errors.New("флаг -name обязателен")
	}
	if *itemType == "" {
		return errors.New("флаг -type обязателен")
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
	fmt.Printf("Успешно добавлен предмет %q с ID: %d\n", *name, id)
	return nil
}

func useItem(inv *inventoryService.Inventory, args []string) error {
	useCmd := flag.NewFlagSet("use", flag.ContinueOnError)
	itemId := useCmd.Int("id", 0, "ID предмета для использования")
	if err := useCmd.Parse(args); err != nil {
		return err
	}
	value, ok := inv.FindItemById(*itemId)
	if !ok {
		return fmt.Errorf("предмет с ID: %d не найден", *itemId)
	}
	msg, err := value.Use()
	if err != nil {
		return err
	}
	fmt.Println(msg)
	return nil
}

func run(args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("%w: список команд: %s", ErrUsage, commandList)
	}
	fs := storage.NewFileStorage("inventory.txt")
	inv := inventoryService.NewInventory(fs)
	_ = inv.Load()
	defer func() {
		if err := inv.Save(); err != nil {
			fmt.Fprintln(os.Stderr, "Ошибка автосохранения:", err)
		}
	}()
	switch args[0] {
	case "list":
		return itemList(inv)
	case "add":
		return addItem(inv, args[1:])
	case "use":
		return useItem(inv, args[1:])
	default:
		return fmt.Errorf("%w: неизвестная команда %s", ErrUsage, args[0])
	}
}

func exitCode(err error) int {
	switch {
	case err == nil:
		return 0
	case errors.Is(err, ErrUsage):
		return 2
	default:
		return 1
	}
}

func report(err error) int {
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
	return exitCode(err)
}

func main() {
	os.Exit(report(run(os.Args[1:])))
}
