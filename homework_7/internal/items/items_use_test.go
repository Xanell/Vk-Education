package items

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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
			name: "Ошибка: Сломанный меч",
			inputItem: &Weapon{
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
			name: "Ошибка: Пустое зелье",
			inputItem: &Potion{
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
			name: "Успех: Использовать целый меч",
			inputItem: &Weapon{
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
		{
			name: "Успех: Использовать полное зелье",
			inputItem: &Potion{
				name:    "Potion",
				Charges: 3,
			},
			wantItem: &Potion{
				name:    "Potion",
				Charges: 2,
			},
			wantMsg: "Использовано Potion",
			wantErr: nil,
		},
		{
			name: "Успех: Экипировать броню",
			inputItem: &Armor{
				name:    "Shield",
				Defense: 10,
			},
			wantItem: &Armor{
				name:    "Shield",
				Defense: 10,
			},
			wantMsg: "Экипировать Shield",
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
