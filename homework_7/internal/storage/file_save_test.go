package storage

import (
	"errors"
	"go-education/homework_7/internal/items"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type failingWriter struct{}

func (f *failingWriter) Write(p []byte) (n int, err error) {
	return 0, errors.New("сбой записи: диск переполнен")
}

func TestFileStorageSave(t *testing.T) {
	t.Run("Сохранение Weapon", func(t *testing.T) {
		fs, filePath := testSetup(t)

		weapon := items.NewWeapon("Меч", 5.5, 10, 100)
		weapon.SetId(1)
		itemsToSave := []items.Item{weapon}
		err := fs.Save(itemsToSave)
		require.NoError(t, err)

		got, err := os.ReadFile(filePath)
		require.NoError(t, err)

		expected := "1|Weapon|Меч|5.50|10|100\n"
		assert.Equal(t, expected, string(got))
	})
	t.Run("Сохранение Potion", func(t *testing.T) {
		fs, filePath := testSetup(t)

		potion := items.NewPotion("Зелье", 0.5, "+50HP", 3)
		potion.SetId(1)
		itemsToSave := []items.Item{potion}
		err := fs.Save(itemsToSave)
		require.NoError(t, err)

		got, err := os.ReadFile(filePath)
		require.NoError(t, err)

		expected := "1|Potion|Зелье|0.50|+50HP|3\n"
		assert.Equal(t, expected, string(got))
	})
	t.Run("Сохранение Armor", func(t *testing.T) {
		fs, filePath := testSetup(t)

		armor := items.NewArmor("Броня", 5.0, 10)
		armor.SetId(1)
		itemsToSave := []items.Item{armor}
		err := fs.Save(itemsToSave)
		require.NoError(t, err)

		got, err := os.ReadFile(filePath)
		require.NoError(t, err)

		expected := "1|Armor|Броня|5.00|10\n"
		assert.Equal(t, expected, string(got))
	})
	t.Run("Ошибка сохранения", func(t *testing.T) {
		fs, _ := testSetup(t)

		weapon := items.NewWeapon("Меч", 5.0, 10, 100)

		err := fs.writeToFile(&failingWriter{}, []items.Item{weapon})

		require.Error(t, err)
		require.ErrorContains(t, err, "слот 0")
		require.ErrorContains(t, err, "сбой записи: диск переполнен")
	})
}
