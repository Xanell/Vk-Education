package main

import (
	"bytes"
	"flag"
	"go-education/homework_7/internal/inventory"
	inventory_mock "go-education/homework_7/internal/inventory/mock"
	"go-education/homework_7/internal/items"
	"go-education/homework_7/internal/storage"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func setupCLI(t *testing.T) (store inventory.Storage, filePath string, stdout *bytes.Buffer, stderr *bytes.Buffer) {
	t.Helper()

	tempDir := t.TempDir()
	filePath = filepath.Join(tempDir, "test_inventory.txt")
	store = storage.NewFileStorage(filePath)
	stdout = new(bytes.Buffer)
	stderr = new(bytes.Buffer)

	return store, filePath, stdout, stderr
}

func TestItemList(t *testing.T) {
	t.Run("Успешный вывод пустого инвентаря", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		mockStore := inventory_mock.NewMockStorage(ctrl)

		var out, errOut bytes.Buffer

		mockStore.EXPECT().Load().Return([]items.Item{}, nil).Times(1)

		err := run([]string{"list"}, &out, &errOut, mockStore)

		require.NoError(t, err)
		assert.Equal(t, "Инвентарь пуст.\n", out.String())
		assert.Empty(t, errOut.String())
	})

	t.Run("Успешный вывод инвентаря с предметами", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		mockStore := inventory_mock.NewMockStorage(ctrl)

		var out, errOut bytes.Buffer

		sword := items.NewWeapon("Меч", 5.0, 10, 100)
		sword.SetId(1)

		mockStore.EXPECT().Load().Return([]items.Item{sword}, nil).Times(1)

		err := run([]string{"list"}, &out, &errOut, mockStore)

		require.NoError(t, err)
		assert.Contains(t, out.String(), "--- Ваш инвентарь ---")
		assert.Contains(t, out.String(), "Меч")
	})
}

func TestAddItem(t *testing.T) {
	tests := []struct {
		name            string
		args            []string
		wantErr         error
		wantOut         string
		wantFileContent string
	}{
		{
			name:            "Успешное добавление оружия",
			args:            []string{"add", "-type", "weapon", "-name", "Экскалибур", "-weight", "5.0"},
			wantErr:         nil,
			wantOut:         "Успешно добавлен и сохранен предмет \"Экскалибур\" с ID: 1\n",
			wantFileContent: "1|Weapon|Экскалибур|5.00|10|100\n",
		},
		{
			name:            "Ошибка отсутствует обязательный флаг",
			args:            []string{"add", "-type", "weapon"},
			wantErr:         ErrFlagIsNotExist,
			wantOut:         "",
			wantFileContent: "",
		},
		{
			name:            "Ошибка неизвестный флаг",
			args:            []string{"add", "-unknown", "weapon"},
			wantErr:         ErrBadFlag,
			wantOut:         "",
			wantFileContent: "",
		},
		{
			name:            "Запрос справки -h",
			args:            []string{"add", "-h"},
			wantErr:         flag.ErrHelp,
			wantOut:         "",
			wantFileContent: "",
		},
		{
			name:            "Ошибка лишние аргументы",
			args:            []string{"add", "-type", "weapon", "-name", "Меч", "лишнее_слово"},
			wantErr:         ErrArgumentsCount,
			wantOut:         "",
			wantFileContent: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			store, filePath, out, errOut := setupCLI(t)
			err := run(tt.args, out, errOut, store)

			if tt.wantErr != nil {
				require.ErrorIs(t, err, tt.wantErr)
				assert.Empty(t, out.String(), "При ошибке stdout должен быть пуст")

				_, statErr := os.Stat(filePath)
				assert.True(t, os.IsNotExist(statErr), "Файл не должен создаваться при ошибке")
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.wantOut, out.String())

				savedData, readErr := os.ReadFile(filePath)
				require.NoError(t, readErr)
				assert.Equal(t, tt.wantFileContent, string(savedData), "Предмет не сохранился в файл")
			}
		})
	}
}

func TestUseItem(t *testing.T) {
	tests := []struct {
		name            string
		fileContent     string
		args            []string
		wantErr         error
		wantOut         string
		wantFileContent string
	}{
		{
			name:            "Успешное использование меча",
			fileContent:     "1|Weapon|Экскалибур|5.00|10|100\n",
			args:            []string{"use", "-id", "1"},
			wantErr:         nil,
			wantOut:         "Использован Экскалибур\n",
			wantFileContent: "1|Weapon|Экскалибур|5.00|10|99\n",
		},
		{
			name:            "Ошибка сломанный меч",
			fileContent:     "1|Weapon|Ржавый меч|5.00|10|0\n",
			args:            []string{"use", "-id", "1"},
			wantErr:         items.ErrBroken,
			wantOut:         "",
			wantFileContent: "",
		},
		{
			name:            "Ошибка не указан ID",
			fileContent:     "",
			args:            []string{"use"},
			wantErr:         ErrFlagIsNotExist,
			wantOut:         "",
			wantFileContent: "",
		},
		{
			name:            "Ошибка предмет не найден",
			fileContent:     "",
			args:            []string{"use", "-id", "99"},
			wantErr:         inventory.ErrItemNotFound,
			wantOut:         "",
			wantFileContent: "",
		},
		{
			name:            "Ошибка неизвестный флаг",
			fileContent:     "",
			args:            []string{"use", "-unknown"},
			wantErr:         ErrBadFlag,
			wantOut:         "",
			wantFileContent: "",
		},
		{
			name:            "Запрос справки -h",
			fileContent:     "",
			args:            []string{"use", "-h"},
			wantErr:         flag.ErrHelp,
			wantOut:         "",
			wantFileContent: "",
		},
		{
			name:            "Ошибка лишние аргументы",
			fileContent:     "",
			args:            []string{"use", "-id", "1", "Weapon"},
			wantErr:         ErrArgumentsCount,
			wantOut:         "",
			wantFileContent: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			store, filePath, out, errOut := setupCLI(t)

			if tt.fileContent != "" {
				err := os.WriteFile(filePath, []byte(tt.fileContent), 0644)
				require.NoError(t, err)
			}
			err := run(tt.args, out, errOut, store)

			if tt.wantErr != nil {
				require.ErrorIs(t, err, tt.wantErr)
				if tt.fileContent != "" {
					savedData, _ := os.ReadFile(filePath)
					assert.Equal(t, tt.fileContent, string(savedData), "Файл не должен меняться при ошибке")
				}
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.wantOut, out.String())

				savedData, readErr := os.ReadFile(filePath)
				require.NoError(t, readErr)
				assert.Equal(t, tt.wantFileContent, string(savedData), "Состояние файла после выполнения не совпадает")
			}
		})
	}
}

func TestDeleteItem(t *testing.T) {
	tests := []struct {
		name            string
		fileContent     string
		args            []string
		wantErr         error
		wantOut         string
		wantFileContent string
	}{
		{
			name:            "Успешное удаление одного предмета из двух",
			fileContent:     "1|Weapon|Меч|5.00|10|100\n2|Armor|Щит|10.00|5\n",
			args:            []string{"delete", "-id", "1"},
			wantErr:         nil,
			wantOut:         "Предмет с ID 1 успешно удален из инвентаря.\n",
			wantFileContent: "2|Armor|Щит|10.00|5\n",
		},
		{
			name:            "Ошибка предмет не найден",
			fileContent:     "",
			args:            []string{"delete", "-id", "99"},
			wantErr:         inventory.ErrItemNotFound,
			wantOut:         "",
			wantFileContent: "",
		},
		{
			name:            "Ошибка не указан ID",
			fileContent:     "",
			args:            []string{"delete"},
			wantErr:         ErrFlagIsNotExist,
			wantOut:         "",
			wantFileContent: "",
		},
		{
			name:            "Ошибка неизвестный флаг",
			fileContent:     "",
			args:            []string{"delete", "-unknown"},
			wantErr:         ErrBadFlag,
			wantOut:         "",
			wantFileContent: "",
		},
		{
			name:            "Запрос справки -h",
			fileContent:     "",
			args:            []string{"delete", "-h"},
			wantErr:         flag.ErrHelp,
			wantOut:         "",
			wantFileContent: "",
		},
		{
			name:            "Ошибка лишние аргументы",
			fileContent:     "",
			args:            []string{"delete", "-id", "1", "Weapon"},
			wantErr:         ErrArgumentsCount,
			wantOut:         "",
			wantFileContent: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			store, filePath, out, errOut := setupCLI(t)

			if tt.fileContent != "" {
				err := os.WriteFile(filePath, []byte(tt.fileContent), 0644)
				require.NoError(t, err)
			}
			err := run(tt.args, out, errOut, store)

			if tt.wantErr != nil {
				require.ErrorIs(t, err, tt.wantErr)
				if tt.fileContent != "" {
					savedData, _ := os.ReadFile(filePath)
					assert.Equal(t, tt.fileContent, string(savedData), "Файл не должен меняться при ошибке")
				}
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.wantOut, out.String())

				savedData, readErr := os.ReadFile(filePath)
				require.NoError(t, readErr)
				assert.Equal(t, tt.wantFileContent, string(savedData), "Состояние файла после выполнения не совпадает")
			}
		})
	}
}
