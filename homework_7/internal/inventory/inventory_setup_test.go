package inventory

import (
	inventory_mock "go-education/homework_7/internal/inventory/mock"
	"testing"

	"go.uber.org/mock/gomock"
)

func testSetup(t *testing.T) (*Inventory, *inventory_mock.MockStorage) {
	t.Helper()
	ctrl := gomock.NewController(t)
	mockStorage := inventory_mock.NewMockStorage(ctrl)
	inv := NewInventory(mockStorage)
	return inv, mockStorage
}
