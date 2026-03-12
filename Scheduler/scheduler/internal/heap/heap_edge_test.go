package heap_test

// -----------------------------------------------------------------------
// Граничные условия — отдельный файл, не меняем heap_test.go
// -----------------------------------------------------------------------

import (
	"testing"

	"github.com/example/scheduler/internal/heap"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// -----------------------------------------------------------------------
// Push / Pop — граничные случаи
// -----------------------------------------------------------------------

// Единственный элемент: он же root, он же min, он же max.
func TestPopSingleElement(t *testing.T) {
	h := heap.New[item]()
	h.Push(entry(1, 42))

	got, ok := h.Pop()
	require.True(t, ok)
	assert.Equal(t, 1, got.Key())
	assert.Equal(t, 0, h.Len())
}

// После Pop всех элементов heap должна быть пустой и корректной для дальнейшего использования.
func TestPopUntilEmpty(t *testing.T) {
	h := heap.New[item]()
	for i := 1; i <= 5; i++ {
		h.Push(entry(i, i*10))
	}
	for h.Len() > 0 {
		_, ok := h.Pop()
		require.True(t, ok)
	}
	assert.Equal(t, 0, h.Len())

	// heap остаётся рабочей после полного опустошения
	h.Push(entry(99, 99))
	assert.Equal(t, 1, h.Len())
}

// -----------------------------------------------------------------------
// Push с дублирующимся ключом
// -----------------------------------------------------------------------

// Push с уже существующим Key() — нарушение контракта, ожидаем панику.
// Уникальность ключей — ответственность вызывающей стороны.
func TestPushDuplicateKeyPanics(t *testing.T) {
	h := heap.New[item]()
	h.Push(entry(1, 10))
	assert.Panics(t, func() {
		h.Push(entry(1, 99))
	}, "Push с дублирующимся ключом должен паниковать")
}

// -----------------------------------------------------------------------
// PopNLowest — граничные случаи
// -----------------------------------------------------------------------

// n=0: должен вернуть пустой слайс, не паниковать.
func TestPopNLowestZero(t *testing.T) {
	h := heap.New[item]()
	h.Push(entry(1, 10))
	h.Push(entry(2, 20))

	result := h.PopNLowest(0)
	assert.Empty(t, result)
	assert.Equal(t, 2, h.Len(), "heap не должна измениться при n=0")
}

// n > Len(): должен вернуть все элементы, не паниковать.
func TestPopNLowestMoreThanLen(t *testing.T) {
	h := heap.New[item]()
	h.Push(entry(1, 10))
	h.Push(entry(2, 20))

	result := h.PopNLowest(100)
	assert.Len(t, result, 2, "должны вернуться все 2 элемента")
	assert.Equal(t, 0, h.Len())
}

// n=0 на пустой heap.
func TestPopNLowestZeroOnEmpty(t *testing.T) {
	h := heap.New[item]()
	result := h.PopNLowest(0)
	assert.Empty(t, result)
}

// n>0 на пустой heap.
func TestPopNLowestOnEmpty(t *testing.T) {
	h := heap.New[item]()
	result := h.PopNLowest(5)
	assert.Empty(t, result)
}

// -----------------------------------------------------------------------
// RemoveByKey — граничные случаи
// -----------------------------------------------------------------------

// Удаление единственного элемента.
func TestRemoveByKeyLastElement(t *testing.T) {
	h := heap.New[item]()
	h.Push(entry(7, 77))

	got, ok := h.RemoveByKey(7)
	require.True(t, ok)
	assert.Equal(t, 7, got.Key())
	assert.Equal(t, 0, h.Len())
	assert.False(t, h.Contains(7))
}

// Удаление корня (максимального элемента).
func TestRemoveByKeyRoot(t *testing.T) {
	h := heap.New[item]()
	h.Push(entry(1, 10))
	h.Push(entry(2, 50)) // root
	h.Push(entry(3, 30))

	_, ok := h.RemoveByKey(2)
	require.True(t, ok)
	assert.Equal(t, 2, h.Len())

	// heap-инвариант должен сохраниться: следующий Pop — максимум оставшихся
	got, ok := h.Pop()
	require.True(t, ok)
	assert.Equal(t, 30, got.Order())
}

// -----------------------------------------------------------------------
// findMinIndex при n=1
// -----------------------------------------------------------------------

// При единственном элементе PopNLowest(1) должен вернуть его же.
func TestPopNLowestSingleElement(t *testing.T) {
	h := heap.New[item]()
	h.Push(entry(1, 42))

	result := h.PopNLowest(1)
	require.Len(t, result, 1)
	assert.Equal(t, 42, result[0].Order())
	assert.Equal(t, 0, h.Len())
}
