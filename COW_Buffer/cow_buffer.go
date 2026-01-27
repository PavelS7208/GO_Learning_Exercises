package main

import (
	"unsafe"
)

type COWBuffer struct {
	data   []byte
	refs   *int
	closed bool
}

func NewCOWBuffer(data []byte) COWBuffer {

	sharedData := make([]byte, len(data))
	copy(sharedData, data)
	//sharedData := data  // заменить на это если требуется шарить с переданными данными
	refCounter := 1
	return COWBuffer{
		data:   sharedData,
		refs:   &refCounter,
		closed: false,
	} // need to implement
}

func (b *COWBuffer) Clone() COWBuffer {
	if b.isClosed() || b.refs == nil || b.data == nil {
		// или паниковать, или возвращать пустой буфер
		panic("Clone() called on closed or invalid COWBuffer")
		//return COWBuffer{}
	}
	// Потоко-небезопасно, но для теста пойдет
	*b.refs++
	return COWBuffer{
		data:   b.data,
		refs:   b.refs,
		closed: false,
	}
}

func (b *COWBuffer) Close() {
	if b.isClosed() || b.refs == nil {
		return // уже закрыт или не инициализирован
	}
	// Потоко-небезопасно, но для теста пойдет
	b.closed = true
	*b.refs--
	b.data = nil
}

func (b *COWBuffer) Update(index int, value byte) bool {
	if b.isClosed() || b.data == nil || index < 0 || index >= len(b.data) {
		return false
	}

	// Если мы единственные владельцы — меняем данные напрямую
	if *b.refs == 1 {
		b.data[index] = value
		return true
	}

	// Иначе: копируем данные (Copy-on-Write)
	newData := make([]byte, len(b.data))
	copy(newData, b.data)
	newData[index] = value

	// Отвязываемся от старого счётчика
	// Создаём новый счётчик для нашей копии
	*b.refs--
	newRefCount := 1
	b.data = newData
	b.refs = &newRefCount

	return true
}

func (b *COWBuffer) String() string {
	if b.isClosed() || len(b.data) == 0 {
		return ""
	}
	//  без создания новой строки на основе data
	//  просто представляем data как набор символов
	return unsafe.String(&b.data[0], len(b.data))
}

func (b *COWBuffer) isClosed() bool {
	return b.closed
}
