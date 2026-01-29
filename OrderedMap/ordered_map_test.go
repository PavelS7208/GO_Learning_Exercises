package main

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBalunCoursesOrderedMap(t *testing.T) {
	data := NewOrderedMap[int, int]()
	assert.Zero(t, data.Size())

	data.Insert(10, 10)
	data.Insert(5, 5)
	data.Insert(15, 15)
	data.Insert(2, 2)
	data.Insert(4, 4)
	data.Insert(12, 12)
	data.Insert(14, 14)

	assert.Equal(t, 7, data.Size())
	assert.True(t, data.Contains(4))
	assert.True(t, data.Contains(12))
	assert.False(t, data.Contains(3))
	assert.False(t, data.Contains(13))

	var keys []int
	expectedKeys := []int{2, 4, 5, 10, 12, 14, 15}
	data.ForEachInOrder(func(key, _ int) {
		keys = append(keys, key)
	})

	assert.True(t, reflect.DeepEqual(expectedKeys, keys))

	data.Delete(15)
	data.Delete(14)
	data.Delete(2)

	assert.Equal(t, 4, data.Size())
	assert.True(t, data.Contains(4))
	assert.True(t, data.Contains(12))
	assert.False(t, data.Contains(2))
	assert.False(t, data.Contains(14))

	keys = nil
	expectedKeys = []int{4, 5, 10, 12}
	data.ForEachInOrder(func(key, _ int) {
		keys = append(keys, key)
	})

	assert.True(t, reflect.DeepEqual(expectedKeys, keys))
}

// Тесты мои
func TestNewOrderedMap(t *testing.T) {
	m := NewOrderedMap[int, string]()

	if m == nil {
		t.Error("NewOrderedMap returned nil")
	}
	if m.Size() != 0 {
		t.Errorf("Expected size 0, got %d", m.Size())
	}
	if m.root != nil {
		t.Error("New map should have nil root")
	}
}

func TestOrderedMap_InsertAndGet(t *testing.T) {
	tests := []struct {
		name     string
		inserts  []struct{ k, v string }
		key      string
		wantVal  string
		wantOk   bool
		wantSize int
	}{
		{
			name:     "empty map",
			key:      "key",
			wantOk:   false,
			wantSize: 0,
		},
		{
			name: "single insert",
			inserts: []struct{ k, v string }{
				{"a", "value1"},
			},
			key:      "a",
			wantVal:  "value1",
			wantOk:   true,
			wantSize: 1,
		},
		{
			name: "multiple inserts",
			inserts: []struct{ k, v string }{
				{"c", "val3"},
				{"a", "val1"},
				{"b", "val2"},
			},
			key:      "b",
			wantVal:  "val2",
			wantOk:   true,
			wantSize: 3,
		},
		{
			name: "update existing key",
			inserts: []struct{ k, v string }{
				{"key", "old"},
				{"key", "new"},
			},
			key:      "key",
			wantVal:  "new",
			wantOk:   true,
			wantSize: 1,
		},
		{
			name: "key not found",
			inserts: []struct{ k, v string }{
				{"a", "val1"},
				{"c", "val3"},
			},
			key:      "b",
			wantOk:   false,
			wantSize: 2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := NewOrderedMap[string, string]()
			for _, ins := range tt.inserts {
				m.Insert(ins.k, ins.v)
			}

			if m.Size() != tt.wantSize {
				t.Errorf("Size() = %d, want %d", m.Size(), tt.wantSize)
			}

			gotVal, gotOk := m.Get(tt.key)
			if gotOk != tt.wantOk {
				t.Errorf("Get() ok = %v, want %v", gotOk, tt.wantOk)
			}
			if gotOk && gotVal != tt.wantVal {
				t.Errorf("Get() value = %v, want %v", gotVal, tt.wantVal)
			}
		})
	}
}

func TestOrderedMap_Contains(t *testing.T) {
	m := NewOrderedMap[int, bool]()

	keys := []int{5, 3, 7, 2, 4, 6, 8}
	for _, k := range keys {
		m.Insert(k, true)
	}

	// Test existing keys
	for _, k := range keys {
		if !m.Contains(k) {
			t.Errorf("Contains(%d) = false, want true", k)
		}
	}

	// Test non-existing keys
	nonKeys := []int{-1, 0, 1, 9, 10}
	for _, k := range nonKeys {
		if m.Contains(k) {
			t.Errorf("Contains(%d) = true, want false", k)
		}
	}
}

func TestOrderedMap_Delete(t *testing.T) {
	tests := []struct {
		name        string
		setup       func() *OrderedMap[int, string]
		keyToDelete int
		wantDeleted bool
		wantSize    int
		checkOrder  []int
	}{
		{
			name: "delete empty map",
			setup: func() *OrderedMap[int, string] {
				return NewOrderedMap[int, string]()
			},
			keyToDelete: 1,
			wantDeleted: false,
			wantSize:    0,
		},
		{
			name: "delete leaf node",
			setup: func() *OrderedMap[int, string] {
				m := NewOrderedMap[int, string]()
				m.Insert(5, "five")
				m.Insert(3, "three")
				m.Insert(7, "seven")
				return m
			},
			keyToDelete: 3,
			wantDeleted: true,
			wantSize:    2,
			checkOrder:  []int{5, 7},
		},
		{
			name: "delete node with one child",
			setup: func() *OrderedMap[int, string] {
				m := NewOrderedMap[int, string]()
				m.Insert(5, "five")
				m.Insert(3, "three")
				m.Insert(2, "two") // left child of 3
				return m
			},
			keyToDelete: 3,
			wantDeleted: true,
			wantSize:    2,
			checkOrder:  []int{2, 5},
		},
		{
			name: "delete node with two children",
			setup: func() *OrderedMap[int, string] {
				m := NewOrderedMap[int, string]()
				m.Insert(5, "five")
				m.Insert(3, "three")
				m.Insert(7, "seven")
				m.Insert(6, "six")
				m.Insert(8, "eight")
				return m
			},
			keyToDelete: 7,
			wantDeleted: true,
			wantSize:    4,
			checkOrder:  []int{3, 5, 6, 8},
		},
		{
			name: "delete root with two children",
			setup: func() *OrderedMap[int, string] {
				m := NewOrderedMap[int, string]()
				m.Insert(5, "five")
				m.Insert(3, "three")
				m.Insert(7, "seven")
				return m
			},
			keyToDelete: 5,
			wantDeleted: true,
			wantSize:    2,
			checkOrder:  []int{3, 7},
		},
		{
			name: "delete non-existent key",
			setup: func() *OrderedMap[int, string] {
				m := NewOrderedMap[int, string]()
				m.Insert(1, "one")
				m.Insert(2, "two")
				return m
			},
			keyToDelete: 3,
			wantDeleted: false,
			wantSize:    2,
			checkOrder:  []int{1, 2},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := tt.setup()
			deleted := m.Delete(tt.keyToDelete)

			if deleted != tt.wantDeleted {
				t.Errorf("Erase() = %v, want %v", deleted, tt.wantDeleted)
			}

			if m.Size() != tt.wantSize {
				t.Errorf("Size() after erase = %d, want %d", m.Size(), tt.wantSize)
			}

			// Verify key is actually deleted
			if deleted && m.Contains(tt.keyToDelete) {
				t.Errorf("Key %d still exists after erase", tt.keyToDelete)
			}

			// Verify order if specified
			if tt.checkOrder != nil {
				var keys []int
				m.ForEachInOrder(func(k int, v string) {
					keys = append(keys, k)
				})

				if len(keys) != len(tt.checkOrder) {
					t.Errorf("Order length = %d, want %d", len(keys), len(tt.checkOrder))
				} else {
					for i := range keys {
						if keys[i] != tt.checkOrder[i] {
							t.Errorf("Order[%d] = %d, want %d", i, keys[i], tt.checkOrder[i])
							break
						}
					}
				}
			}
		})
	}
}

func TestOrderedMap_ForEachInOrder(t *testing.T) {
	m := NewOrderedMap[int, string]()

	// Insert in random order
	insertOrder := []int{5, 2, 8, 1, 3, 7, 9, 4, 6}
	for _, k := range insertOrder {
		m.Insert(k, fmt.Sprintf("val%d", k))
	}

	// Should iterate in sorted order
	expectedOrder := []int{1, 2, 3, 4, 5, 6, 7, 8, 9}
	var actualOrder []int

	m.ForEachInOrder(func(k int, v string) {
		actualOrder = append(actualOrder, k)
		expectedVal := fmt.Sprintf("val%d", k)
		if v != expectedVal {
			t.Errorf("ForEachInOrder: value for key %d = %s, want %s", k, v, expectedVal)
		}
	})

	if len(actualOrder) != len(expectedOrder) {
		t.Errorf("ForEachInOrder visited %d items, want %d", len(actualOrder), len(expectedOrder))
	}

	for i := range actualOrder {
		if actualOrder[i] != expectedOrder[i] {
			t.Errorf("ForEachInOrder order[%d] = %d, want %d", i, actualOrder[i], expectedOrder[i])
			break
		}
	}

	// Test with empty map
	emptyMap := NewOrderedMap[string, int]()
	count := 0
	emptyMap.ForEachInOrder(func(k string, v int) {
		count++
	})
	if count != 0 {
		t.Errorf("ForEachInOrder on empty map called callback %d times", count)
	}
}

func TestOrderedMap_String(t *testing.T) {
	tests := []struct {
		name     string
		setup    func() *OrderedMap[int, string]
		expected string
	}{
		{
			name: "empty map",
			setup: func() *OrderedMap[int, string] {
				return NewOrderedMap[int, string]()
			},
			expected: "[]",
		},
		{
			name: "single element",
			setup: func() *OrderedMap[int, string] {
				m := NewOrderedMap[int, string]()
				m.Insert(1, "one")
				return m
			},
			expected: "[(1,one)]",
		},
		{
			name: "multiple elements sorted",
			setup: func() *OrderedMap[int, string] {
				m := NewOrderedMap[int, string]()
				m.Insert(3, "three")
				m.Insert(1, "one")
				m.Insert(2, "two")
				return m
			},
			expected: "[(1,one) (2,two) (3,three)]",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := tt.setup()
			result := m.String()
			if result != tt.expected {
				t.Errorf("String() = %q, want %q", result, tt.expected)
			}
		})
	}
}

func TestOrderedMap_DegenerateCase(t *testing.T) {
	// Test worst-case scenario (sorted insert creates linked list)
	m := NewOrderedMap[int, string]()

	// Insert in sorted order (worst case for BST)
	for i := 1; i <= 100; i++ {
		m.Insert(i, fmt.Sprintf("val%d", i))
	}

	// Verify all values are accessible
	for i := 1; i <= 100; i++ {
		val, ok := m.Get(i)
		if !ok {
			t.Errorf("Get(%d) failed after sequential insert", i)
		}
		if val != fmt.Sprintf("val%d", i) {
			t.Errorf("Get(%d) = %s, want val%d", i, val, i)
		}
	}

	// Verify size
	if m.Size() != 100 {
		t.Errorf("Size() = %d, want 100", m.Size())
	}

	// Verify order
	expected := 1
	m.ForEachInOrder(func(k int, v string) {
		if k != expected {
			t.Errorf("Order mismatch: got %d, want %d", k, expected)
		}
		expected++
	})
}

func TestOrderedMap_ComplexTypes(t *testing.T) {
	// Test 1: Complex value types with simple keys
	type complexData struct {
		value  string
		nested struct {
			id int
		}
	}
	m := NewOrderedMap[string, complexData]()

	data1 := complexData{value: "first", nested: struct{ id int }{id: 1}}
	data2 := complexData{value: "second", nested: struct{ id int }{id: 2}}

	m.Insert("key1", data1)
	m.Insert("key2", data2)

	val, ok := m.Get("key1")
	if !ok || val.value != "first" || val.nested.id != 1 {
		t.Errorf("Get complex value = %+v, %v, want proper struct, true", val, ok)
	}

	// Test 2: Different ordered key types

	// String keys
	t.Run("string keys", func(t *testing.T) {
		m := NewOrderedMap[string, int]()
		m.Insert("zebra", 1)
		m.Insert("apple", 2)

		if !m.Contains("apple") {
			t.Error("string key not found")
		}

		// Check order
		var keys []string
		m.ForEachInOrder(func(k string, v int) {
			keys = append(keys, k)
		})

		if len(keys) != 2 || keys[0] != "apple" || keys[1] != "zebra" {
			t.Errorf("string keys order = %v, want [apple zebra]", keys)
		}
	})

	// Float keys
	t.Run("float keys", func(t *testing.T) {
		m := NewOrderedMap[float64, string]()
		m.Insert(3.14, "pi")
		m.Insert(2.71, "e")
		m.Insert(1.41, "sqrt2")

		val, ok := m.Get(2.71)
		if !ok || val != "e" {
			t.Errorf("float key Get = %v, %v, want e, true", val, ok)
		}

		// Check order
		var keys []float64
		m.ForEachInOrder(func(k float64, v string) {
			keys = append(keys, k)
		})

		expected := []float64{1.41, 2.71, 3.14}
		for i, k := range keys {
			if k != expected[i] {
				t.Errorf("float order[%d] = %v, want %v", i, k, expected[i])
			}
		}
	})

	// Int keys with negative values
	t.Run("int keys with negatives", func(t *testing.T) {
		m := NewOrderedMap[int, string]()
		m.Insert(100, "hundred")
		m.Insert(-50, "minus fifty")
		m.Insert(0, "zero")
		m.Insert(-100, "minus hundred")

		// Check order
		keys := []int{}
		m.ForEachInOrder(func(k int, v string) {
			keys = append(keys, k)
		})

		expected := []int{-100, -50, 0, 100}
		for i, k := range keys {
			if k != expected[i] {
				t.Errorf("int order[%d] = %d, want %d", i, k, expected[i])
			}
		}
	})

	// Byte/rune keys
	t.Run("rune keys", func(t *testing.T) {
		m := NewOrderedMap[rune, string]()
		m.Insert('z', "last")
		m.Insert('a', "first")
		m.Insert('m', "middle")

		keys := []rune{}
		m.ForEachInOrder(func(k rune, v string) {
			keys = append(keys, k)
		})

		expected := []rune{'a', 'm', 'z'}
		for i, k := range keys {
			if k != expected[i] {
				t.Errorf("rune order[%d] = %c, want %c", i, k, expected[i])
			}
		}
	})
}

func TestOrderedMap_MixedOperations(t *testing.T) {
	m := NewOrderedMap[int, int]()

	// Phase 1: Insert many elements
	for i := 0; i < 50; i++ {
		m.Insert(i, i*10)
	}

	if m.Size() != 50 {
		t.Errorf("After insert phase: size = %d, want 50", m.Size())
	}

	// Phase 2: Update some elements
	for i := 0; i < 50; i += 2 {
		m.Insert(i, i*100) // Update even numbers
	}

	// Phase 3: Delete some elements
	for i := 10; i < 20; i++ {
		m.Delete(i)
	}

	if m.Size() != 40 {
		t.Errorf("After delete phase: size = %d, want 40", m.Size())
	}

	// Verify remaining elements
	for i := 0; i < 50; i++ {
		val, ok := m.Get(i)
		if i >= 10 && i < 20 {
			// Deleted range
			if ok {
				t.Errorf("Key %d should be deleted but found with value %d", i, val)
			}
		} else {
			// Should exist
			if !ok {
				t.Errorf("Key %d not found but should exist", i)
			}
			expected := i * 10
			if i%2 == 0 {
				expected = i * 100 // Updated value
			}
			if val != expected {
				t.Errorf("Key %d = %d, want %d", i, val, expected)
			}
		}
	}

	// Verify order is still correct
	prev := -1
	m.ForEachInOrder(func(k, v int) {
		if k <= prev {
			t.Errorf("Order violation: %d comes after %d", k, prev)
		}
		prev = k
	})
}

func TestOrderedMap_ZeroValueKeys(t *testing.T) {
	m := NewOrderedMap[int, string]()

	// Test with zero value key
	m.Insert(0, "zero")
	val, ok := m.Get(0)
	if !ok || val != "zero" {
		t.Errorf("Zero key: got %v, %v, want zero, true", val, ok)
	}

	// Test with other keys including zero
	m.Insert(-1, "minus one")
	m.Insert(1, "one")

	// Verify all three
	testCases := []struct {
		key      int
		expected string
	}{
		{-1, "minus one"},
		{0, "zero"},
		{1, "one"},
	}

	for _, tc := range testCases {
		val, ok := m.Get(tc.key)
		if !ok || val != tc.expected {
			t.Errorf("Key %d: got %v, %v, want %s, true", tc.key, val, ok, tc.expected)
		}
	}

	// Verify order
	expectedOrder := []int{-1, 0, 1}
	actualOrder := []int{}
	m.ForEachInOrder(func(k int, v string) {
		actualOrder = append(actualOrder, k)
	})

	for i := range expectedOrder {
		if actualOrder[i] != expectedOrder[i] {
			t.Errorf("Order[%d] = %d, want %d", i, actualOrder[i], expectedOrder[i])
		}
	}
}

// Benchmark tests
func BenchmarkOrderedMap_Insert(b *testing.B) {
	for n := 0; n < b.N; n++ {
		m := NewOrderedMap[int, int]()
		for i := 0; i < 1000; i++ {
			m.Insert(i, i*2)
		}
	}
}

func BenchmarkOrderedMap_Get(b *testing.B) {
	m := NewOrderedMap[int, int]()
	for i := 0; i < 1000; i++ {
		m.Insert(i, i*2)
	}

	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		m.Get(n % 1000)
	}
}

func BenchmarkOrderedMap_ForEachInOrder(b *testing.B) {
	m := NewOrderedMap[int, int]()
	for i := 0; i < 1000; i++ {
		m.Insert(i, i*2)
	}

	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		sum := 0
		m.ForEachInOrder(func(k, v int) {
			sum += v
		})
		_ = sum
	}
}
