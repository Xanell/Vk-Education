package utils

import (
	"fmt"
	"go-education/homework_6/internal/items"
)

func Filter[T any](items []T, predicate func(T) bool) []T {
	result := make([]T, 0, len(items))
	for _, item := range items {
		if predicate(item) {
			result = append(result, item)
		}
	}
	return result
}

func Map[T any, R any](items []T, transform func(T) R) []R {
	result := make([]R, len(items))
	for index, item := range items {
		result[index] = transform(item)
	}
	return result
}

func Find[T any](items []T, condition func(T) bool) (T, bool) {
	for _, item := range items {
		if condition(item) {
			return item, true
		}
	}
	var zero T
	return zero, false
}

func DescribeItem(i domain.Item) string {
	switch value := i.(type) {
	case nil:
		return "Предмет отсутствует"
	case *domain.Weapon:
		if value == nil {
			return "Предмет отсутствует"
		}
	case *domain.Armor:
		if value == nil {
			return "Предмет отсутствует"
		}
	case *domain.Potion:
		if value == nil {
			return "Предмет отсутствует"
		}
	default:
		return "Неизвестный предмет"
	}
	return fmt.Sprintf("[ID: %d] %s (вес: %.2f)", i.GetId(), i.Name(), i.Weight())
}
