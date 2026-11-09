package inventory

import (
	"go-education/homework_7/internal/items"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInventoryFindItemById(t *testing.T) {
	t.Run("Успешное нахождение предмета по ID", func(t *testing.T) {
		inv, _ := testSetup(t)

		goodSword := items.NewWeapon("Целый меч", 5.0, 10, 50)
		goodPotion := items.NewPotion("Полное зелье", 0.5, "+50HP", 3)

		inv.AddItem(goodSword)
		potionID := inv.AddItem(goodPotion)

		got, ok := inv.FindItemById(potionID)
		require.True(t, ok, "Метод должен был вернуть true (предмет найден)")
		assert.Equal(t, goodPotion, got, "Найденный предмет не совпадает с ожидаемым")
	})
	t.Run("Ошибка предмет не найден", func(t *testing.T) {
		inv, _ := testSetup(t)

		inv.AddItem(items.NewWeapon("Целый меч", 5.0, 10, 50))

		got, ok := inv.FindItemById(99)
		assert.False(t, ok, "Флаг ok должен быть false")
		assert.Nil(t, got, "Предмет должен быть nil, так как он не найден")
	})
	t.Run("Поиск в пустом инвентаре", func(t *testing.T) {
		inv, _ := testSetup(t)

		got, ok := inv.FindItemById(1)
		assert.False(t, ok)
		assert.Nil(t, got)
	})
}
