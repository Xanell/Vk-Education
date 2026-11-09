package items

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

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
