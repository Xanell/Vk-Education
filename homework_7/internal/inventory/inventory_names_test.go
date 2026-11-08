package inventory

import (
	"go-education/homework_7/internal/items"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInventoryItemNames(t *testing.T) {
	t.Run("Успешный вывод наименований предметов в инвентаре", func(t *testing.T) {
		inv, _ := testSetup(t)

		brokenSword := items.NewWeapon("Сломанный меч", 5.0, 10, 0)
		goodSword := items.NewWeapon("Целый меч", 5.0, 10, 50)
		emptyPotion := items.NewPotion("Пустое зелье", 0.5, "+50HP", 0)
		goodPotion := items.NewPotion("Полное зелье", 0.5, "+50HP", 3)
		armor := items.NewArmor("Броня", 10.0, 5)

		inv.AddItem(brokenSword)
		inv.AddItem(goodSword)
		inv.AddItem(emptyPotion)
		inv.AddItem(goodPotion)
		inv.AddItem(armor)

		expected := []string{
			"Сломанный меч",
			"Целый меч",
			"Пустое зелье",
			"Полное зелье",
			"Броня",
		}
		got := inv.ItemNames()
		assert.Equal(t, expected, got)
	})
	t.Run("Пустой инвентарь возвращает пустой список", func(t *testing.T) {
		inv, _ := testSetup(t)

		got := inv.ItemNames()

		assert.Empty(t, got, "В пустом инвентаре не должно быть имён")
	})
}
