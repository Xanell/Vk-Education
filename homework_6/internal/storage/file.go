package storage

import (
	"bufio"
	"errors"
	"fmt"
	"go-education/homework_6/internal/inventory"
	"go-education/homework_6/internal/items"
	"io"
	"os"
	"strings"
)

type FileStorage struct {
	fileName string
}

func NewFileStorage(fileName string) *FileStorage {
	if fileName == "" {
		fileName = "inventory.txt"
	}
	return &FileStorage{
		fileName: fileName,
	}
}

func (fs *FileStorage) Save(inv *inventoryService.Inventory) (err error) {
	file, err := os.Create(fs.fileName)
	if err != nil {
		return fmt.Errorf("ошибка сохранения в файл: %w", err)
	}
	defer func() {
		closeErr := file.Close()
		if closeErr != nil {
			err = errors.Join(err, fmt.Errorf("ошибка при закрытии файла: %w", closeErr))
		}
	}()
	return fs.writeToFile(file, inv.Items)
}

func (fs *FileStorage) writeToFile(w io.Writer, items []domain.Item) error {
	var errs []error
	for i, item := range items {
		if storableItem, ok := item.(domain.Storable); ok {
			if err := storableItem.Serialize(w); err != nil {
				errs = append(errs, fmt.Errorf("слот %d (предмет %q): %w", i, item.Name(), err))
			}
		}
	}
	return errors.Join(errs...)
}

func (fs *FileStorage) loadFromFile(r io.Reader, inv *inventoryService.Inventory) error {
	var errs []error
	existingIDs := make(map[int]bool)
	inv.Items = make([]domain.Item, 0)
	inv.NextID = 1
	scanner := bufio.NewScanner(r)
	lineNum := 0
	for scanner.Scan() {
		var item domain.Item
		var err error

		lineNum++
		line := scanner.Text()
		lineSlice := strings.Split(line, "|")

		if len(lineSlice) < 2 {
			errs = append(errs, domain.NewParseItemError(lineNum, line, fmt.Errorf("недостаточно данных: строка не содержит тип предмета")))
			continue
		}

		switch lineSlice[1] {
		case "Weapon":
			weapon := &domain.Weapon{}
			err = weapon.Deserialize(strings.NewReader(line))
			item = weapon
		case "Armor":
			armor := &domain.Armor{}
			err = armor.Deserialize(strings.NewReader(line))
			item = armor
		case "Potion":
			potion := &domain.Potion{}
			err = potion.Deserialize(strings.NewReader(line))
			item = potion
		default:
			err = fmt.Errorf("неизвестный тип предмета %s", lineSlice[1])
		}
		if err != nil {
			errs = append(errs, domain.NewParseItemError(lineNum, line, err))
			continue
		}
		if existingIDs[item.GetId()] {
			errs = append(errs, domain.NewParseItemError(lineNum, line, fmt.Errorf("дубликат ID %d", item.GetId())))
			continue
		}
		existingIDs[item.GetId()] = true
		inv.Items = append(inv.Items, item)
		if item.GetId() >= inv.NextID {
			inv.NextID = item.GetId() + 1
		}
	}
	if err := scanner.Err(); err != nil {
		errs = append(errs, fmt.Errorf("ошибка чтения из файла: %w", err))
	}
	return errors.Join(errs...)
}

func (fs *FileStorage) Load(inv *inventoryService.Inventory) error {
	file, err := os.Open(fs.fileName)
	if err != nil {
		return fmt.Errorf("ошибка загрузки из файла: %w", err)
	}
	defer file.Close()
	return fs.loadFromFile(file, inv)
}
