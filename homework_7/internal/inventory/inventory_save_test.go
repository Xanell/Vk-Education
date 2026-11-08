package inventory

import (
	"errors"
	"go-education/homework_7/internal/items"
	"testing"

	"github.com/stretchr/testify/require"
)

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
