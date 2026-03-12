package heap_test

import (
	"strings"
	"testing"

	"github.com/example/scheduler/internal/heap"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// item is a minimal heap.item implementation for testing.
// key uniquely identifies the entry; order determines heap position.
type item struct {
	key   int
	order int
}

func (i item) Key() int   { return i.key }
func (i item) Order() int { return i.order }

// entry is a convenience constructor.
func entry(key, order int) item { return item{key: key, order: order} }

// -----------------------------------------------------------------------
// Push / Pop
// -----------------------------------------------------------------------

func TestPushPop(t *testing.T) {
	h := heap.New[item]()
	h.Push(entry(1, 10))
	h.Push(entry(2, 50))
	h.Push(entry(3, 30))

	got, ok := h.Pop()
	require.True(t, ok)
	assert.Equal(t, 50, got.Order())

	got, ok = h.Pop()
	require.True(t, ok)
	assert.Equal(t, 30, got.Order())
}

func TestPopEmpty(t *testing.T) {
	h := heap.New[item]()
	_, ok := h.Pop()
	assert.False(t, ok)
}

// -----------------------------------------------------------------------
// PopNLowest
// -----------------------------------------------------------------------

func TestPopNLowest(t *testing.T) {
	h := heap.New[item]()
	for i, order := range []int{0, 25, 50, 75, 100} {
		h.Push(entry(i+1, order))
	}

	lowest := h.PopNLowest(2)
	require.Len(t, lowest, 2)

	// The next item from the heap must have a higher order than all stolen ones.
	remaining, ok := h.Pop()
	require.True(t, ok)
	for _, it := range lowest {
		assert.Less(t, it.Order(), remaining.Order())
	}
}

func TestPopNLowestHeapInvariant(t *testing.T) {
	h := heap.New[item]()
	for i := 1; i <= 8; i++ {
		h.Push(entry(i, i*10))
	}

	h.PopNLowest(3)

	prev := 100
	for h.Len() > 0 {
		got, ok := h.Pop()
		require.True(t, ok)
		assert.LessOrEqual(t, got.Order(), prev)
		prev = got.Order()
	}
}

// -----------------------------------------------------------------------
// RemoveByKey
// -----------------------------------------------------------------------

func TestRemoveByKey(t *testing.T) {
	h := heap.New[item]()
	h.Push(entry(1, 10))
	h.Push(entry(2, 20))

	got, ok := h.RemoveByKey(1)
	require.True(t, ok)
	assert.Equal(t, 1, got.Key())
	assert.False(t, h.Contains(1))
}

func TestRemoveByKeyNotFound(t *testing.T) {
	h := heap.New[item]()
	_, ok := h.RemoveByKey(999)
	assert.False(t, ok)
}

// -----------------------------------------------------------------------
// Contains
// -----------------------------------------------------------------------

func TestContains(t *testing.T) {
	h := heap.New[item]()
	h.Push(entry(42, 10))
	assert.True(t, h.Contains(42))
	assert.False(t, h.Contains(99))
}

// -----------------------------------------------------------------------
// String
// -----------------------------------------------------------------------

func TestStringEmpty(t *testing.T) {
	h := heap.New[item]()
	assert.Equal(t, "(empty)", h.String())
}

func TestStringLevels(t *testing.T) {
	h := heap.New[item]()
	for i, order := range []int{10, 20, 30, 40, 50} {
		h.Push(entry(i+1, order))
	}
	s := h.String()

	assert.Contains(t, s, "level 0:")
	assert.Contains(t, s, "level 1:")
	assert.Contains(t, s, "level 2:")

	// Root must be the max-order item.
	firstLine := strings.SplitN(s, "\n", 2)[0]
	assert.Contains(t, firstLine, "order=50")
}
