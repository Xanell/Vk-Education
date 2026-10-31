package inventoryService

import (
	"go-education/homework_6/internal/items"
	"go-education/homework_6/pkg/utils"
)

type Storage interface {
	Save(inv *Inventory) error
	Load(inv *Inventory) error
}

type Inventory struct {
	NextID int
	Items  []domain.Item
	store  Storage
}

func NewInventory(store Storage) *Inventory {
	return &Inventory{
		NextID: 1,
		Items:  make([]domain.Item, 0),
		store:  store,
	}
}

func (inv *Inventory) AddItem(item domain.Item) int {
	item.SetId(inv.NextID)
	inv.Items = append(inv.Items, item)
	inv.NextID++
	return item.GetId()
}

func (inv *Inventory) Weapons() []*domain.Weapon {
	weapons := utils.Filter(inv.Items, func(item domain.Item) bool {
		_, ok := item.(*domain.Weapon)
		return ok
	})
	return utils.Map(weapons, func(item domain.Item) *domain.Weapon {
		w, ok := item.(*domain.Weapon)
		if !ok {
			return nil
		}
		return w
	})
}

func (inv *Inventory) BrokenItems() []domain.Item {
	return utils.Filter(inv.Items, func(item domain.Item) bool {
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
	return utils.Map(inv.Items, func(item domain.Item) string {
		return item.Name()
	})
}

func (inv *Inventory) FindItemById(id int) (domain.Item, bool) {
	return utils.Find(inv.Items, func(item domain.Item) bool {
		return item.GetId() == id
	})
}

func (inv *Inventory) Load() error {
	return inv.store.Load(inv)
}

func (inv *Inventory) Save() error {
	return inv.store.Save(inv)
}
