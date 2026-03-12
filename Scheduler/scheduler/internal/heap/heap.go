package heap

import (
	"fmt"
	"strings"
)

// -------------  MaxHeap Binary Heap с добавками для скорости --------------------------
//  На вершине значение item с максимальным значением возвращаемым Order()
//  Есть еще Key  - уникальный идентификатор item-а

// KeyType тип для идентификации объектов в структуре (index lookup).
type KeyType = int

// OrderType Тип п которому определяем сортировку в структуре (heap comparison).
type OrderType = int

// Item — контракт для хранения в MaxHeap.
// Key()   — уникальный идентификатор для O(1) поиска
// Order() — значение по которому сортируем в структуре
type Item interface {
	Key() KeyType
	Order() OrderType
}

// item — приватный алиас для обратной совместимости внутри пакета

type item = Item

// MaxHeap Not thread-safe — синхронизацией должны заниматься его клиенты
type MaxHeap[T item] struct {
	items   []T
	indices map[KeyType]int // Key() -> index в items
}

func New[T item]() MaxHeap[T] {
	return MaxHeap[T]{
		items:   make([]T, 0),
		indices: make(map[int]int),
	}
}

func (h *MaxHeap[T]) Len() int {
	return len(h.items)
}

func (h *MaxHeap[T]) Contains(key KeyType) bool {
	_, ok := h.indices[key]
	return ok
}

// Push добавляем элемент в структуру.
// Паникует если элемент с таким Key() уже существует —
// уникальность ключей является обязанностью вызывающей стороны.
func (h *MaxHeap[T]) Push(item T) {
	if _, exists := h.indices[item.Key()]; exists {
		panic(fmt.Sprintf("heap: duplicate key %d", item.Key()))
	}
	h.items = append(h.items, item)
	idx := len(h.items) - 1
	h.indices[item.Key()] = idx
	h.siftUp(idx)
}

// Pop removes and returns the item with the highest Order().
// Returns (zero value, false) if the heap is empty.
func (h *MaxHeap[T]) Pop() (T, bool) {
	if len(h.items) == 0 {
		var zero T
		return zero, false
	}
	return h.removeAt(0)
}

// RemoveByKey удаляет по ключу (не значению по которому сортируем)
// Returns (zero value, false) if the key is not found.
func (h *MaxHeap[T]) RemoveByKey(key KeyType) (T, bool) {
	idx, ok := h.indices[key]
	if !ok {
		var zero T
		return zero, false
	}
	return h.removeAt(idx)
}

// PopNLowest забираем из структуры N элементов самых последних (в смысле Order) и грохаем их
// Порядок возвращаемых - произвольный, не сортированный
func (h *MaxHeap[T]) PopNLowest(n int) []T {
	result := make([]T, 0, n)
	for i := 0; i < n && len(h.items) > 0; i++ {
		minIdx := h.findMinIndex()
		item, _ := h.removeAt(minIdx)
		result = append(result, item)
	}
	return result
}

// ----  Сервисные внутренние функции, работают с индексами в слайсе хранения ---------

// findMinIndex находим индекс для элемента самого маленького в смысле  Order().
// In a max-heap the non-leaf nodes occupy indices 0..n/2-1, so all
// leaf nodes are in the range [n/2, n). The minimum is always a leaf
// because every parent is >= its children. When n==1 the single root
// is both the only leaf and the minimum, and n/2 == 0 is correct.
func (h *MaxHeap[T]) findMinIndex() int {
	n := len(h.items)
	minIdx := n / 2 // index первого листа с слайсе (или root при n==1)
	for i := minIdx + 1; i < n; i++ {
		if h.items[i].Order() < h.items[minIdx].Order() {
			minIdx = i
		}
	}
	return minIdx
}

// removeAt - удаляем элемент по индексу в слайсе хранения
func (h *MaxHeap[T]) removeAt(i int) (T, bool) {
	if len(h.items) == 0 {
		var zero T
		return zero, false
	}

	item := h.items[i]
	delete(h.indices, item.Key())

	last := len(h.items) - 1
	if i == last {
		h.items = h.items[:last]
		return item, true
	}

	h.items[i] = h.items[last]
	h.items = h.items[:last]
	h.indices[h.items[i].Key()] = i

	parent := (i - 1) / 2
	if i > 0 && h.items[i].Order() > h.items[parent].Order() {
		h.siftUp(i)
	} else {
		h.siftDown(i)
	}
	return item, true
}

// Опускание и поднимание при вставке/удалении в структуру
//
//	Тут  i - это индекс слайса хранения
func (h *MaxHeap[T]) siftUp(i int) {
	for i > 0 {
		parent := (i - 1) / 2
		if h.items[parent].Order() >= h.items[i].Order() {
			break
		}
		h.swap(i, parent)
		i = parent
	}
}

func (h *MaxHeap[T]) siftDown(i int) {
	n := len(h.items)
	for {
		largest := i
		left, right := 2*i+1, 2*i+2
		if left < n && h.items[left].Order() > h.items[largest].Order() {
			largest = left
		}
		if right < n && h.items[right].Order() > h.items[largest].Order() {
			largest = right
		}
		if largest == i {
			break
		}
		h.swap(i, largest)
		i = largest
	}
}

func (h *MaxHeap[T]) swap(i, j int) {
	h.indices[h.items[i].Key()] = j
	h.indices[h.items[j].Key()] = i
	h.items[i], h.items[j] = h.items[j], h.items[i]
}

// Вывод в удобном виде в форме иерархии и индекса в слайсе
func (h *MaxHeap[T]) String() string {
	n := len(h.items)
	if n == 0 {
		return "(empty)"
	}

	var sb strings.Builder
	level := 0
	levelStart := 0

	for levelStart < n {
		width := 1 << level // 2^level
		end := min(levelStart+width, n)

		fmt.Fprintf(&sb, "level %d:", level)
		for i := levelStart; i < end; i++ {
			fmt.Fprintf(&sb, "  [%d] key=%d order=%d", i, h.items[i].Key(), h.items[i].Order())
		}
		sb.WriteByte('\n')

		levelStart += width
		level++
	}
	return strings.TrimRight(sb.String(), "\n")
}
