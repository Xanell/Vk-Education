package inventory

import (
	"fmt"
	"go-education/homework_6/internal/items"
	"go-education/homework_6/pkg/collections"
)

type Storage interface {
	Save(item []domain.Item) error
	Load() ([]domain.Item, error)
}

type Inventory struct {
	nextID int
	Items  []domain.Item
	store  Storage
}

func NewInventory(store Storage) *Inventory {
	return &Inventory{
		nextID: 1,
		Items:  make([]domain.Item, 0),
		store:  store,
	}
}

func (inv *Inventory) AddItem(item domain.Item) int {
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

func (inv *Inventory) Weapons() []*domain.Weapon {
	weapons := collections.Filter(inv.Items, func(item domain.Item) bool {
		_, ok := item.(*domain.Weapon)
		return ok
	})
	return collections.Map(weapons, func(item domain.Item) *domain.Weapon {
		w, ok := item.(*domain.Weapon)
		if !ok {
			return nil
		}
		return w
	})
}

func (inv *Inventory) BrokenItems() []domain.Item {
	return collections.Filter(inv.Items, func(item domain.Item) bool {
		switch value := item.(type) {
		case *domain.Weapon:
			return value.Durability <= 0
		case *domain.Potion:
			return value.Charges <= 0
		default:
			return false
		}
	})
}

func (inv *Inventory) ItemNames() []string {
	return collections.Map(inv.Items, func(item domain.Item) string {
		return item.Name()
	})
}

func (inv *Inventory) FindItemById(id int) (domain.Item, bool) {
	return collections.Find(inv.Items, func(item domain.Item) bool {
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
	return fmt.Errorf("ошибка при загрузке данных: %w", err)
}

func (inv *Inventory) Save() error {
	return inv.store.Save(inv.Items)
}
