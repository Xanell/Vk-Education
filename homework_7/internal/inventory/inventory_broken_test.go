package inventory

import (
	"go-education/homework_7/internal/items"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInventoryBrokenItems(t *testing.T) {
	t.Run("Успешно находит все сломанные и израсходованные предметы", func(t *testing.T) {
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

		got := inv.BrokenItems()
		expected := []items.Item{
			brokenSword,
			emptyPotion,
		}
		require.Len(t, got, 2, "Должно быть найдено ровно 2 сломанных предмета")
		assert.Equal(t, expected, got)
	})
	t.Run("Возвращает пустой список, если сломанные предметы отсутствуют", func(t *testing.T) {
		inv, _ := testSetup(t)

		inv.AddItem(items.NewWeapon("Меч", 5.0, 10, 100))
		inv.AddItem(items.NewPotion("Зелье", 0.5, "+50HP", 3))

		got := inv.BrokenItems()
		assert.Len(t, got, 0, "Сломанных предметов быть не должно")
		assert.Empty(t, got)
	})
	t.Run("Пустой инвентарь", func(t *testing.T) {
		inv, _ := testSetup(t)

		got := inv.BrokenItems()
		assert.Len(t, got, 0, "В пустом инвентаре не может быть сломанных предметов")
		assert.Empty(t, got)
	})
}
