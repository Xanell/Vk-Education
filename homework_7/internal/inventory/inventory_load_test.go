package inventory

import (
	"errors"
	"go-education/homework_7/internal/items"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInventoryLoad(t *testing.T) {
	t.Run("Успешная загрузка инвентаря", func(t *testing.T) {
		inv, mock := testSetup(t)
		expected := []items.Item{items.NewWeapon("Меч", 5.0, 10, 100)}
		mock.EXPECT().
			Load().
			Return(expected, nil).
			Times(1)
		err := inv.Load()
		require.NoError(t, err)
		assert.Equal(t, expected, inv.Items)
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
	t.Run("Проверка ID после загрузки инвентаря", func(t *testing.T) {
		inv, mock := testSetup(t)
		loadedSword := items.NewWeapon("Старый меч", 5.0, 10, 100)
		loadedSword.SetId(10)
		mock.EXPECT().Load().Return([]items.Item{loadedSword}, nil)

		err := inv.Load()
		require.NoError(t, err)

		newPotion := items.NewPotion("Новое зелье", 0.5, "+50HP", 3)
		newID := inv.AddItem(newPotion)
		assert.Equal(t, 11, newID, "Новый предмет должен получить ID 11")
	})
}
