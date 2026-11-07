package inventory

import (
	"errors"
	"go-education/homework_7/internal/inventory/mock"
	"go-education/homework_7/internal/items"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func testSetup(t *testing.T) (*Inventory, *inventory_mock.MockStorage) {
	t.Helper()
	ctrl := gomock.NewController(t)
	mockStorage := inventory_mock.NewMockStorage(ctrl)
	inv := NewInventory(mockStorage)
	return inv, mockStorage
}

func TestInventoryAddItem(t *testing.T) {
	inv, _ := testSetup(t)

	sword := items.NewWeapon("Меч", 5.0, 10, 100)
	potion := items.NewPotion("Зелье", 0.5, "+50HP", 3)

	id1 := inv.AddItem(sword)
	assert.Equal(t, 1, id1, "Первый предмет должен получить ID 1")
	assert.Equal(t, 1, len(inv.Items), "В инвентаре должен быть 1 предмет")
	assert.Equal(t, 2, inv.nextID)

	id2 := inv.AddItem(potion)
	assert.Equal(t, 2, id2, "Второй предмет должен получить ID 2")
	assert.Equal(t, 2, len(inv.Items), "В инвентаре должно быть 2 предмета")
	assert.Equal(t, 3, inv.nextID)

	assert.Equal(t, sword, inv.Items[0])
	assert.Equal(t, potion, inv.Items[1])
}

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

	t.Run("Ошибка: предмет не найден", func(t *testing.T) {
		inv, _ := testSetup(t)

		inv.AddItem(items.NewWeapon("Меч", 5.0, 10, 100)) // ID 1

		err := inv.RemoveItemById(99)

		require.ErrorIs(t, err, ErrItemNotFound)
		assert.Len(t, inv.Items, 1)
		assert.Equal(t, 1, inv.Items[0].GetId())
	})

	t.Run("Ошибка: удаление из пустого инвентаря", func(t *testing.T) {
		inv, _ := testSetup(t)

		err := inv.RemoveItemById(1)

		require.ErrorIs(t, err, ErrItemNotFound)
		assert.Len(t, inv.Items, 0)
	})
}

func TestInventoryLoad(t *testing.T) {
	t.Run("Успешная загрузка инвентаря", func(t *testing.T) {
		inv, mock := testSetup(t)
		expectedItem := []items.Item{items.NewWeapon("Меч", 5.0, 10, 100)}
		mock.EXPECT().
			Load().
			Return(expectedItem, nil).
			Times(1)
		err := inv.Load()
		require.NoError(t, err)
		assert.Equal(t, expectedItem, inv.Items)
	})
	t.Run("Ошибка при загрузке инвентаря", func(t *testing.T) {
		inv, mock := testSetup(t)
		wantErr := errors.New("ошибка записи на диск")
		mock.EXPECT().
			Load().
			Return(nil, wantErr).
			Times(1)
		err := inv.Load()
		require.ErrorIs(t, err, wantErr)
	})
	t.Run("Загрузка пустого файла очищает инвентарь", func(t *testing.T) {
		inv, mock := testSetup(t)
		inv.AddItem(items.NewWeapon("Старый меч", 5.0, 10, 100))
		mock.EXPECT().Load().Return([]items.Item{}, nil)
		err := inv.Load()
		require.NoError(t, err)
		assert.Len(t, inv.Items, 0, "Инвентарь должен очиститься при загрузке пустого файла")
		assert.Equal(t, 1, inv.nextID, "nextID должен сброситься на 1")
	})
}

func TestInventorySave(t *testing.T) {
	t.Run("Успешное сохранение инвентаря", func(t *testing.T) {
		inv, mock := testSetup(t)
		weapon := items.NewWeapon("Меч", 5.0, 10, 100)
		inv.AddItem(weapon)
		mock.EXPECT().
			Save(inv.Items).
			Return(nil).
			Times(1)
		err := inv.Save()
		require.NoError(t, err)
	})
	t.Run("Ошибка при сохранении инвентаря", func(t *testing.T) {
		inv, mock := testSetup(t)
		wantErr := errors.New("ошибка записи на диск")
		mock.EXPECT().
			Save(inv.Items).
			Return(wantErr).
			Times(1)
		err := inv.Save()
		require.ErrorIs(t, err, wantErr)
	})
}
