package items

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

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
