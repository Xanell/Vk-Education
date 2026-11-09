package inventory

import (
	"go-education/homework_7/internal/items"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInventoryRemoveItemById(t *testing.T) {
	t.Run("Успешное удаление из середины", func(t *testing.T) {
		inv, _ := testSetup(t)

		inv.AddItem(items.NewWeapon("Меч", 5.0, 10, 100))
		inv.AddItem(items.NewPotion("Зелье", 0.5, "+50HP", 3))
		inv.AddItem(items.NewArmor("Щит", 10.0, 5))

		err := inv.RemoveItemById(2)

		require.NoError(t, err)
		assert.Len(t, inv.Items, 2)
		assert.Equal(t, 1, inv.Items[0].GetId())
		assert.Equal(t, 3, inv.Items[1].GetId())
		assert.Equal(t, 4, inv.nextID, "Счетчик nextID не должен меняться при удалении")
	})

	t.Run("Ошибка предмет не найден", func(t *testing.T) {
		inv, _ := testSetup(t)

		inv.AddItem(items.NewWeapon("Меч", 5.0, 10, 100)) // ID 1

		err := inv.RemoveItemById(99)

		require.ErrorIs(t, err, ErrItemNotFound)
		assert.Len(t, inv.Items, 1)
		assert.Equal(t, 1, inv.Items[0].GetId())
	})

	t.Run("Ошибка удаление из пустого инвентаря", func(t *testing.T) {
		inv, _ := testSetup(t)

		err := inv.RemoveItemById(1)

		require.ErrorIs(t, err, ErrItemNotFound)
		assert.Len(t, inv.Items, 0)
	})
}
