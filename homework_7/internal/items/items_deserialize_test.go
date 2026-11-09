package items

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestItemDeserialize(t *testing.T) {
	for _, tt := range []struct {
		name    string
		input   string
		item    Storable
		want    Storable
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
