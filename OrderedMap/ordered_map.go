package main

import (
	"cmp"
	"fmt"
	"strings"
)

type OrderedMap[K cmp.Ordered, V any] struct {
	root *node[K, V]
	size int
}

type node[K cmp.Ordered, V any] struct {
	left  *node[K, V]
	right *node[K, V]
	key   K
	value V
}

func (n *node[K, V]) get(key K) (V, bool) {
	if n == nil {
		var zero V
		return zero, false
	}
	switch cmp.Compare(key, n.key) {
	case -1:
		return n.left.get(key)
	case 1:
		return n.right.get(key)
	default:
		return n.value, true
	}
}

func (n *node[K, V]) insert(key K, value V) (*node[K, V], bool) {
	if n == nil {
		return &node[K, V]{key: key, value: value}, true
	}
	var inserted bool
	switch cmp.Compare(key, n.key) {
	case -1:
		n.left, inserted = n.left.insert(key, value)
		return n, inserted
	case 1:
		n.right, inserted = n.right.insert(key, value)
		return n, inserted
	default: // key == n.key
		n.value = value
		return n, false // no new node
	}
}

func (n *node[K, V]) delete(key K) (*node[K, V], bool) {
	if n == nil {
		return nil, false
	}
	var deleted bool
	switch cmp.Compare(key, n.key) {
	case -1: //  Меньше - ушли влево по дереву
		n.left, deleted = n.left.delete(key)
		return n, deleted
	case 1: //  Больше - ушли вправо по дереву
		n.right, deleted = n.right.delete(key)
		return n, deleted
	default: //  Нашли что удалять
		if n.left == nil { // Нет одного из детей, значит простой случай
			return n.right, true
		}
		if n.right == nil {
			return n.left, true
		}
		// Два ребенка - ищем по алгоритму и перестраиваем дерево внизу
		minNode := n.right.min()
		n.key = minNode.key
		n.value = minNode.value
		n.right, _ = n.right.delete(minNode.key) // we know it exists
		return n, true
	}
}

func (n *node[K, V]) min() *node[K, V] {
	cur := n
	for cur.left != nil {
		cur = cur.left
	}
	return cur
}

func (n *node[K, V]) contains(key K) bool {
	if n == nil {
		return false
	}
	switch cmp.Compare(key, n.key) {
	case -1:
		return n.left.contains(key)
	case 1:
		return n.right.contains(key)
	default: //  Если равно
		return true
	}
}

func (n *node[K, V]) forEachInOrder(fn func(key K, value V)) {
	if n == nil {
		return
	}
	n.left.forEachInOrder(fn)  //  Сначала по ключам меньшим
	fn(n.key, n.value)         // По ключу
	n.right.forEachInOrder(fn) // По ключам большим
}

func NewOrderedMap[K cmp.Ordered, V any]() *OrderedMap[K, V] {
	return &OrderedMap[K, V]{root: nil}
}

func (m *OrderedMap[K, V]) Get(key K) (V, bool) {
	return m.root.get(key)
}

func (m *OrderedMap[K, V]) Insert(key K, value V) bool {
	var inserted bool
	m.root, inserted = m.root.insert(key, value)
	if inserted {
		m.size++
	}
	return inserted // полезно знать, был ли это update или insert
}

func (m *OrderedMap[K, V]) Delete(key K) bool {
	var deleted bool
	m.root, deleted = m.root.delete(key)
	if deleted {
		m.size--
	}
	return deleted
}

func (m *OrderedMap[K, V]) Contains(key K) bool {
	return m.root.contains(key)
}

func (m *OrderedMap[K, V]) Size() int {
	return m.size
}

func (m *OrderedMap[K, V]) Empty() bool {
	return m.size == 0
}

func (m *OrderedMap[K, V]) ForEachInOrder(action func(K, V)) {
	if m == nil {
		return
	}
	m.root.forEachInOrder(action)
}

func (m *OrderedMap[K, V]) String() string {
	if m == nil || m.size == 0 {
		return "[]"
	}

	var b strings.Builder
	b.WriteByte('[')
	first := true
	m.ForEachInOrder(func(key K, value V) {
		if !first {
			b.WriteByte(' ')
		}
		fmt.Fprintf(&b, "(%v,%v)", key, value)
		first = false
	})
	b.WriteByte(']')
	return b.String()
}
