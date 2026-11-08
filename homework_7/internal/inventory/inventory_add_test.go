package inventory

import (
	"go-education/homework_7/internal/items"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInventoryAddItem(t *testing.T) {
	inv, _ := testSetup(t)

	sword := items.NewWeapon("Меч", 5.0, 10, 100)
	potion := items.NewPotion("Зелье", 0.5, "+50HP", 3)

	id1 := inv.AddItem(sword)
	assert.Equal(t, 1, id1, "Первый предмет должен получить ID 1")
	assert.Equal(t, 1, len(inv.Items), "В инвентаре должен быть 1 предмет")
	assert.Equal(t, 2, inv.nextID)

	id2 := inv.AddItem(potion)
	assert.Equal(t, 2, id2, "Второй предмет должен получить ID 2")
	assert.Equal(t, 2, len(inv.Items), "В инвентаре должно быть 2 предмета")
	assert.Equal(t, 3, inv.nextID)

	assert.Equal(t, sword, inv.Items[0])
	assert.Equal(t, potion, inv.Items[1])
}
