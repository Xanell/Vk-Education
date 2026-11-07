package items

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"strings"
	"testing"
)

func TestItemUse(t *testing.T) {
	for _, tt := range []struct {
		name      string
		inputItem Item
		wantItem  Item
		wantMsg   string
		wantErr   error
	}{
		{
			name: "Ошибка: Сломанный меч", inputItem: &Weapon{
				name:       "Sword",
				Durability: 0,
			},
			wantItem: &Weapon{
				name:       "Sword",
				Durability: 0,
			},
			wantMsg: "",
			wantErr: ErrBroken,
		},
		{
			name: "Ошибка: Пустое зелье", inputItem: &Potion{
				name:    "Potion",
				Charges: 0,
			},
			wantItem: &Potion{
				name:    "Potion",
				Charges: 0,
			},
			wantMsg: "",
			wantErr: ErrExhausted,
		},
		{
			name: "Использовать: Целый меч", inputItem: &Weapon{
				name:       "Sword",
				Durability: 100,
			},
			wantItem: &Weapon{
				name:       "Sword",
				Durability: 99,
			},
			wantMsg: "Использован Sword",
			wantErr: nil,
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			msg, err := tt.inputItem.Use()
			require.ErrorIs(t, err, tt.wantErr)
			assert.Equal(t, tt.wantMsg, msg)
			assert.Equal(t, tt.wantItem, tt.inputItem)
		})
	}
}

func TestItemIsBroken(t *testing.T) {
	for _, tt := range []struct {
		name      string
		input     Breakable
		wantState bool
	}{
		{
			name: "Прочность равна 1", input: &Weapon{
				name:       "Sword",
				Durability: 1,
			},
			wantState: false,
		},
		{
			name: "Прочность равна 0", input: &Weapon{
				name:       "Sword",
				Durability: 0,
			},
			wantState: true,
		},
		{
			name: "Прочность равна -1", input: &Weapon{
				name:       "Sword",
				Durability: -1,
			},
			wantState: true,
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.wantState, tt.input.IsBroken())
		})
	}
}

func TestItemSerialize(t *testing.T) {
	for _, tt := range []struct {
		name      string
		inputItem Storable
		want      string
	}{
		{
			name: "Сериализация оружия",
			inputItem: &Weapon{
				ID:         1,
				name:       "Sword",
				weight:     10.0,
				Damage:     10,
				Durability: 100,
			},
			want: "1|Weapon|Sword|10.00|10|100\n",
		},
		{
			name: "Сериализация брони",
			inputItem: &Armor{
				ID:      1,
				name:    "Броня",
				weight:  5.0,
				Defense: 10,
			},
			want: "1|Armor|Броня|5.00|10\n",
		},
		{
			name: "Сериализация Зелья",
			inputItem: &Potion{
				ID:      1,
				name:    "Хилка",
				weight:  0.5,
				Effect:  "+50Hp",
				Charges: 3,
			},
			want: "1|Potion|Хилка|0.50|+50Hp|3\n",
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			err := tt.inputItem.Serialize(&buf)
			require.NoError(t, err)
			assert.Equal(t, tt.want, buf.String())
		})
	}
}

func TestItemDeserialize(t *testing.T) {
	for _, tt := range []struct {
		name    string
		input   string
		item    Storable
		want    *Weapon
		wantErr bool
	}{
		{
			name:  "Успешный парсинг предмета",
			input: "1|Weapon|Sword|10.00|10|100\n",
			item:  &Weapon{},
			want: &Weapon{
				ID:         1,
				name:       "Sword",
				weight:     10.0,
				Damage:     10,
				Durability: 100,
			},
			wantErr: false,
		},
		{
			name:    "Ошибка: не хватает полей",
			input:   "1|Weapon|",
			item:    &Weapon{},
			want:    &Weapon{},
			wantErr: true,
		},
		{
			name:    "Ошибка: неверный формат урона",
			input:   "1|Weapon|Sword|10.00|abcd|100\n",
			item:    &Weapon{},
			want:    &Weapon{},
			wantErr: true,
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			reader := strings.NewReader(tt.input)
			err := tt.item.Deserialize(reader)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			require.NoError(t, err)
			got := tt.item.(*Weapon)
			assert.Equal(t, tt.want, got)
		})
	}
}
