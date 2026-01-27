package main

import (
	"fmt"
	"strconv"
	"strings"
)

// Number — ограничение для всех числовых типов
type Number interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64
}

type CircularQueue[T Number] struct {
	values []T
	cap    int
	front  int
	size   int
}

func NewCircularQueue[T Number](capacity int) *CircularQueue[T] {
	capacity = max(capacity, 1)
	return &CircularQueue[T]{
		values: make([]T, capacity),
		cap:    capacity,
	}
}

func (cq *CircularQueue[T]) Empty() bool {
	return cq.size == 0
}

func (cq *CircularQueue[T]) Full() bool {
	return cq.size == cq.cap
}

func (cq *CircularQueue[T]) Push(value T) bool {
	if cq.Full() {
		return false
	}
	rear := (cq.front + cq.size) % cq.cap
	cq.values[rear] = value
	cq.size++
	return true
}

func (cq *CircularQueue[T]) Pop() bool {
	if cq.Empty() {
		return false
	}
	// Необязательно обнулять, но можно для отладки
	cq.values[cq.front] = 0
	cq.front = (cq.front + 1) % cq.cap
	cq.size--
	return true
}

func (cq *CircularQueue[T]) Front() T {
	if cq.Empty() {
		return -1
	}
	return cq.values[cq.front]
}

func (cq *CircularQueue[T]) Back() T {
	if cq.Empty() {
		return -1
	}
	rear := (cq.front + cq.size - 1) % cq.cap
	return cq.values[rear]
}

func (cq *CircularQueue[T]) String() string {
	var sb strings.Builder

	// Предвыделение памяти для типичного случая
	sb.Grow(50 + cq.size*3)

	sb.WriteString("CircularQueue{size:")
	sb.WriteString(strconv.Itoa(cq.size))
	sb.WriteString(", cap:")
	sb.WriteString(strconv.Itoa(cq.cap))
	sb.WriteString(", front:")
	sb.WriteString(strconv.Itoa(cq.front))
	sb.WriteString(", data:[")

	if !cq.Empty() {
		for i := 0; i < cq.size; i++ {
			idx := (cq.front + i) % cq.cap
			//sb.WriteString(strconv.Itoa(cq.values[idx]))
			fmt.Fprintf(&sb, "%v", cq.values[idx])
			if i < cq.size-1 {
				sb.WriteString(", ")
			}
		}
	}
	sb.WriteString("]}")

	return sb.String()
}
