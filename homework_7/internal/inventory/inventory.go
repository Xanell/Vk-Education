package inventory

import (
	"fmt"
	"go-education/homework_7/internal/items"
	"go-education/homework_7/pkg/collections"
)

type Storage interface {
	Save(item []items.Item) error
	Load() ([]items.Item, error)
}

type Inventory struct {
	nextID int
	Items  []items.Item
	store  Storage
}

func NewInventory(store Storage) *Inventory {
	return &Inventory{
		nextID: 1,
		Items:  make([]items.Item, 0),
		store:  store,
	}
}

func (inv *Inventory) AddItem(item items.Item) int {
	item.SetId(inv.nextID)
	inv.Items = append(inv.Items, item)
	inv.nextID++
	return item.GetId()
}

func (inv *Inventory) RemoveItemById(id int) error {
	for i, item := range inv.Items {
		if item.GetId() == id {
			inv.Items = append(inv.Items[:i], inv.Items[i+1:]...)
			return nil
		}
	}
	return fmt.Errorf("предмет с ID %d не найден", id)
}

func (inv *Inventory) Weapons() []*items.Weapon {
	weapons := collections.Filter(inv.Items, func(item items.Item) bool {
		_, ok := item.(*items.Weapon)
		return ok
	})
	return collections.Map(weapons, func(item items.Item) *items.Weapon {
		w, ok := item.(*items.Weapon)
		if !ok {
			return nil
		}
		return w
	})
}

func (inv *Inventory) BrokenItems() []items.Item {
	return collections.Filter(inv.Items, func(item items.Item) bool {
		switch value := item.(type) {
		case *items.Weapon:
			return value.Durability <= 0
		case *items.Potion:
			return value.Charges <= 0
		default:
			return false
		}
	})
}

func (inv *Inventory) ItemNames() []string {
	return collections.Map(inv.Items, func(item items.Item) string {
		return item.Name()
	})
}

func (inv *Inventory) FindItemById(id int) (items.Item, bool) {
	return collections.Find(inv.Items, func(item items.Item) bool {
		return item.GetId() == id
	})
}

func (inv *Inventory) Load() error {
	items, err := inv.store.Load()
	if items != nil {
		inv.Items = items
	}
	inv.nextID = 1
	for _, item := range inv.Items {
		if item.GetId() >= inv.nextID {
			inv.nextID = item.GetId() + 1
		}
	}
	if err != nil {
		return fmt.Errorf("ошибка при загрузке данных: %w", err)
	}
	return nil
}

func (inv *Inventory) Save() error {
	return inv.store.Save(inv.Items)
}
