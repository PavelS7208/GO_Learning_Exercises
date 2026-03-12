package scheduler

import (
	"sync"
	"sync/atomic"

	"github.com/example/scheduler/internal/heap"
)

// SchedulableItem — контракт для элементов которые можно хранить в планировщике.
// Определён на стороне потребителя (pool), реализуется неявно (task).
type SchedulableItem interface {
	heap.Item
}

// minStealSize - константа для конфигуратора балансировщика.
// Сколько в node минимально Item-ов должно быть чтобы из него было разрешено забирать задачи
const minStealSize = 4

// ---------- node - Это защищенная для конкуретного доступа обертка --------------
//
//	Нода Планировщика. Позволяет добавлять и выдавать SchedulableItem (задачу) согласно приоритетам.
//	Хранит список SchedulableItem-ов в бинарной куче отсортированной в порядке уменьшения значений приоритета.
//	Позволяет на лету изменять приоритет уже добавленной задачи
//	В случае запуска балансировщика позволяет отдать половину низкоприоритетных задач для перемещения их в другую ноду
type node[T SchedulableItem] struct {
	mu   sync.Mutex
	heap heap.MaxHeap[T]
	// count mirrors heap.Len() приближённо.
	// Обновляется атомарно ПОСЛЕ освобождения mu, поэтому значение может быть
	// кратковременно stale. Подходит для эвристик (выбор жертвы при rebalance),
	// но не для проверок инвариантов.
	count atomic.Int64
}

func newNode[T SchedulableItem]() *node[T] {
	return &node[T]{
		heap: heap.New[T](),
	}
}

// Count Кол-во SchedulableItem (задач) в планировщике
// Используем Atomic counter — lock не требуется
func (n *node[T]) Count() int {
	return int(n.count.Load())
}

// add SchedulableItem в Планировщик
func (n *node[T]) add(t T) {
	n.mu.Lock()
	n.heap.Push(t)
	n.mu.Unlock()
	n.count.Add(1)
}

// get возвращает с удалением item c максимальным приоритетом
// Возвращает (zero value, false) если нода пустая
func (n *node[T]) get() (T, bool) {
	n.mu.Lock()
	t, ok := n.heap.Pop()
	n.mu.Unlock()
	if ok {
		n.count.Add(-1)
	}
	return t, ok
}

// change Меняем по идентификатору приоритет item (задачи).
// Возвращаем false - если ключ идентификатор неверный
func (n *node[T]) change(key heap.KeyType, updated T) bool {
	n.mu.Lock()
	defer n.mu.Unlock()
	_, ok := n.heap.RemoveByKey(key)
	if !ok {
		return false
	}
	n.heap.Push(updated)
	return true
}

// Сервисный метод для реализации алгоритма балансировки
// Позволяет своровать (stealing) половину низкоприоритетных задач
// Использует параметр minStealSize

func (n *node[T]) stealHalf() []T {
	n.mu.Lock()
	if n.heap.Len() < minStealSize {
		n.mu.Unlock()
		return nil
	}
	half := n.heap.Len() / 2
	stolen := n.heap.PopNLowest(half)
	n.mu.Unlock()
	n.count.Add(-int64(len(stolen)))
	return stolen
}
