package inventory

import (
	"go-education/homework_7/internal/items"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInventoryWeapons(t *testing.T) {
	t.Run("Успешная фильтрация по категории: оружие", func(t *testing.T) {
		inv, _ := testSetup(t)

		sword := items.NewWeapon("Меч", 5.0, 10, 100)
		potion := items.NewPotion("Зелье", 0.5, "+50HP", 3)
		armor := items.NewArmor("Броня", 5.0, 10)

		inv.AddItem(sword)
		inv.AddItem(potion)
		inv.AddItem(armor)

		got := inv.Weapons()

		assert.Len(t, got, 1, "После фильтрации по категории оружие должен получиться 1 предмет sword")
		assert.Equal(t, []*items.Weapon{sword}, got)
	})
	t.Run("Нет оружия в инвентаре", func(t *testing.T) {
		inv, _ := testSetup(t)

		inv.AddItem(items.NewPotion("Зелье", 0.5, "+50HP", 3))
		inv.AddItem(items.NewArmor("Броня", 5.0, 10))

		got := inv.Weapons()
		assert.Len(t, got, 0, "Список оружия должен быть пустым")
		assert.Empty(t, got)
	})
	t.Run("Пустой инвентарь", func(t *testing.T) {
		inv, _ := testSetup(t)

		got := inv.Weapons()
		assert.Len(t, got, 0, "В пустом инвентаре не может быть оружия")
		assert.Empty(t, got)
	})
}
