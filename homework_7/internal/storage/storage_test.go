package storage

import (
	"go-education/homework_7/internal/items"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupStorage(t *testing.T) (*FileStorage, string) {
	t.Helper()
	tempDir := t.TempDir()
	filePath := filepath.Join(tempDir, "test_inventory.txt")
	fs := NewFileStorage(filePath)
	return fs, filePath
}

func TestFileStorageSave(t *testing.T) {
	fs, filePath := setupStorage(t)
	weapon := items.NewWeapon("Меч", 5.0, 10, 100)
	weapon.SetId(1)

	itemsToSave := []items.Item{weapon}

	err := fs.Save(itemsToSave)
	require.NoError(t, err)

	content, _ := os.ReadFile(filePath)
	assert.Contains(t, string(content), "1|Weapon|Меч|5.00|10|100")
}

func TestFileStorageLoad(t *testing.T) {
	t.Run("Успешная загрузка из файла", func(t *testing.T) {
		fs, filePath := setupStorage(t)
		content := "1|Weapon|Меч|5.00|10|100\n2|Armor|Щит|10.00|5\n"
		err := os.WriteFile(filePath, []byte(content), 0644)
		require.NoError(t, err, "Не удалось подготовить тестовый файл")

		loadedItems, err := fs.Load()

		require.NoError(t, err, "Загрузка должна пройти без ошибок")

		require.Len(t, loadedItems, 2)

		assert.Equal(t, "Меч", loadedItems[0].Name())
		assert.Equal(t, 1, loadedItems[0].GetId())
		assert.Equal(t, "Щит", loadedItems[1].Name())
		assert.Equal(t, 2, loadedItems[1].GetId())
	})
	t.Run("Ошибка загрузки из файла", func(t *testing.T) {
		fs, _ := setupStorage(t)
		_, err := fs.Load()
		require.ErrorIs(t, err, os.ErrNotExist)
	})
	t.Run("Ошибка загрузки битые строки в файле", func(t *testing.T) {
		fs, filePath := setupStorage(t)
		content := "1|Weapon|Меч|5.00|10|100\n" +
			"2|Unknown|Кольцо|0.5|5\n" +
			"3|Weapon|Топор|abc|10|100\n"
		err := os.WriteFile(filePath, []byte(content), 0644)
		require.NoError(t, err)
		loadedItems, err := fs.Load()
		require.Len(t, loadedItems, 1)
		assert.Equal(t, "Меч", loadedItems[0].Name())
		require.Error(t, err)
		assert.Contains(t, err.Error(), "неизвестный тип")
		assert.Contains(t, err.Error(), "неверный формат для вес")
	})
}
