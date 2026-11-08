package storage

import (
	"errors"
	"go-education/homework_7/internal/items"
	"os"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type failingReader struct{}

func (f *failingReader) Read(p []byte) (n int, err error) {
	return 0, errors.New("сбой чтения: диск отключен")
}

func TestFileStorageLoad(t *testing.T) {
	t.Run("Успешная загрузка Weapon", func(t *testing.T) {
		fs, filePath := testSetup(t)
		content := "1|Weapon|Меч|5.00|10|100\n"

		err := os.WriteFile(filePath, []byte(content), 0644)
		require.NoError(t, err, "Не удалось подготовить тестовый файл")

		got, err := fs.Load()
		require.NoError(t, err, "Загрузка должна пройти без ошибок")

		require.Len(t, got, 1, "Должен загрузиться 1 предмет")

		expected := items.NewWeapon("Меч", 5.0, 10, 100)
		expected.SetId(1)

		assert.Equal(t, expected, got[0])
	})
	t.Run("Успешная загрузка Potion", func(t *testing.T) {
		fs, filePath := testSetup(t)
		content := "1|Potion|Зелье|0.50|+50HP|3\n"

		err := os.WriteFile(filePath, []byte(content), 0644)
		require.NoError(t, err, "Не удалось подготовить тестовый файл")

		got, err := fs.Load()
		require.NoError(t, err, "Загрузка должна пройти без ошибок")

		require.Len(t, got, 1, "Должен загрузиться 1 предмет")

		expected := items.NewPotion("Зелье", 0.5, "+50HP", 3)
		expected.SetId(1)

		assert.Equal(t, expected, got[0])
	})
	t.Run("Успешная загрузка Armor", func(t *testing.T) {
		fs, filePath := testSetup(t)
		content := "1|Armor|Броня|5.00|10\n"

		err := os.WriteFile(filePath, []byte(content), 0644)
		require.NoError(t, err, "Не удалось подготовить тестовый файл")

		got, err := fs.Load()
		require.NoError(t, err, "Загрузка должна пройти без ошибок")

		require.Len(t, got, 1, "Должен загрузиться 1 предмет")

		expected := items.NewArmor("Броня", 5.0, 10)
		expected.SetId(1)

		assert.Equal(t, expected, got[0])
	})
	t.Run("Ошибка загрузки из файла", func(t *testing.T) {
		fs, _ := testSetup(t)
		_, err := fs.Load()
		require.ErrorIs(t, err, os.ErrNotExist)
	})
	t.Run("Ошибка загрузки битые строки в файле", func(t *testing.T) {
		fs, filePath := testSetup(t)

		validLine := "1|Weapon|Меч|5.00|10|100"
		brokenLine := "2|Unknown|Кольцо|0.5|5"
		content := validLine + "\n" + brokenLine + "\n"

		err := os.WriteFile(filePath, []byte(content), 0644)
		require.NoError(t, err)

		got, err := fs.Load()
		require.Len(t, got, 1, "Проверяем что целая строка загрузилась")
		assert.Equal(t, "Меч", got[0].Name())

		require.Error(t, err)
		var parseErr *items.ParseItemError
		require.True(t, errors.As(err, &parseErr), "Ошибка должна содержать ParseItemError")

		assert.Equal(t, 2, parseErr.Line, "Номер строки должен быть 2")
		assert.Equal(t, brokenLine, parseErr.Value, "Текст строки не совпадает")
		assert.ErrorContains(t, parseErr.Err, "неизвестный тип предмета")
	})
	t.Run("Ошибка загрузки парсинг строки в файле", func(t *testing.T) {
		fs, filePath := testSetup(t)

		validLine := "1|Weapon|Меч|5.00|10|100"
		brokenLine := "2|Armor|Броня|abc|5"
		content := validLine + "\n" + brokenLine + "\n"

		err := os.WriteFile(filePath, []byte(content), 0644)
		require.NoError(t, err)

		got, err := fs.Load()
		require.Len(t, got, 1, "Проверяем что целая строка загрузилась")
		assert.Equal(t, "Меч", got[0].Name())

		require.Error(t, err)
		var parseErr *items.ParseItemError
		require.True(t, errors.As(err, &parseErr), "Ошибка должна содержать ParseItemError")

		assert.Equal(t, 2, parseErr.Line, "Номер строки должен быть 2")
		assert.Equal(t, brokenLine, parseErr.Value, "Текст строки не совпадает")
		assert.ErrorIs(t, parseErr.Err, strconv.ErrSyntax)
	})
	t.Run("Ошибка загрузки", func(t *testing.T) {
		fs, _ := testSetup(t)

		_, err := fs.loadFromFile(&failingReader{})

		require.Error(t, err)
		assert.ErrorContains(t, err, "сбой чтения: диск отключен")
	})
}
