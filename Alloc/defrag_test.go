package main

import (
	"testing"
	"unsafe"
)

func TestDefragment(t *testing.T) {
	t.Run("empty inputs", func(t *testing.T) {
		// Пустая память
		var memory []byte
		pointers := []unsafe.Pointer{unsafe.Pointer(&memory)}
		Defragment(memory, pointers) // не должно упасть

		// Пустые указатели
		memory = make([]byte, 100)
		pointers = []unsafe.Pointer{}
		Defragment(memory, pointers) // не должно упасть
	})

	t.Run("all nil pointers", func(t *testing.T) {
		memory := make([]byte, 100)
		pointers := []unsafe.Pointer{nil, nil, nil}
		Defragment(memory, pointers)

		// Все указатели должны остаться nil
		for i, p := range pointers {
			if p != nil {
				t.Errorf("pointer %d should be nil", i)
			}
		}
	})

	t.Run("skip nil pointers but process valid ones", func(t *testing.T) {
		memory := make([]byte, 100)
		basePtr := uintptr(unsafe.Pointer(&memory[0]))

		// Создаем блок
		memory[0] = 5
		for i := 1; i <= 5; i++ {
			memory[i] = byte(i)
		}
		validPtr := unsafe.Pointer(&memory[1])

		pointers := []unsafe.Pointer{nil, validPtr, nil}

		Defragment(memory, pointers)

		// Проверяем что nil остались nil, а валидный указатель не изменился
		if pointers[0] != nil || pointers[2] != nil {
			t.Error("nil pointers changed")
		}
		if uintptr(pointers[1])-basePtr != 1 {
			t.Error("valid pointer should remain at position 1")
		}
	})

	t.Run("already compacted - no movement", func(t *testing.T) {
		memory := make([]byte, 100)
		basePtr := uintptr(unsafe.Pointer(&memory[0]))

		// Блоки уже плотно упакованы
		// Блок 1: заголовок в 0, данные с 1 (3 байта)
		memory[0] = 3
		memory[1], memory[2], memory[3] = 10, 11, 12
		ptr1 := unsafe.Pointer(&memory[1])

		// Блок 2: заголовок в 4, данные с 5 (2 байта)
		memory[4] = 2
		memory[5], memory[6] = 20, 21
		ptr2 := unsafe.Pointer(&memory[5])

		pointers := []unsafe.Pointer{ptr2, ptr1} // неотсортированные

		Defragment(memory, pointers)

		// Проверяем что указатели не изменились
		if uintptr(pointers[0])-basePtr != 5 {
			t.Errorf("ptr2 changed from 5 to %d", uintptr(pointers[0])-basePtr)
		}
		if uintptr(pointers[1])-basePtr != 1 {
			t.Errorf("ptr1 changed from 1 to %d", uintptr(pointers[1])-basePtr)
		}

		// Проверяем что данные не испортились
		if memory[1] != 10 || memory[2] != 11 || memory[3] != 12 {
			t.Error("block 1 data corrupted")
		}
		if memory[5] != 20 || memory[6] != 21 {
			t.Error("block 2 data corrupted")
		}
	})

	t.Run("single block already at correct position", func(t *testing.T) {
		memory := make([]byte, 100)
		//basePtr := uintptr(unsafe.Pointer(&memory[0]))

		memory[0] = 5
		for i := 1; i <= 5; i++ {
			memory[i] = byte(i)
		}
		originalPtr := unsafe.Pointer(&memory[1])

		pointers := []unsafe.Pointer{originalPtr}

		Defragment(memory, pointers)

		if pointers[0] != originalPtr {
			t.Error("single block should not move when already at position 1")
		}
	})

	t.Run("single block needs moving from non-optimal position", func(t *testing.T) {
		memory := make([]byte, 100)
		basePtr := uintptr(unsafe.Pointer(&memory[0]))

		// Блок на позиции 10 (данные с 10, заголовок в 9)
		memory[9] = 3
		memory[10], memory[11], memory[12] = 1, 2, 3
		originalPtr := unsafe.Pointer(&memory[10])

		pointers := []unsafe.Pointer{originalPtr}

		Defragment(memory, pointers)

		newPos := uintptr(pointers[0]) - basePtr
		if newPos != 1 {
			t.Errorf("block should move to position 1, got %d", newPos)
		}

		// Проверяем заголовок на новом месте
		if memory[0] != 3 {
			t.Error("header not set correctly at new position")
		}
		// Проверяем данные
		if memory[1] != 1 || memory[2] != 2 || memory[3] != 3 {
			t.Error("data not copied correctly")
		}
	})

	t.Run("two blocks with gap between them", func(t *testing.T) {
		memory := make([]byte, 100)
		basePtr := uintptr(unsafe.Pointer(&memory[0]))

		// Блок 1 на позиции 1 (оптимально)
		memory[0] = 2
		memory[1], memory[2] = 10, 11
		ptr1 := unsafe.Pointer(&memory[1])

		// Блок 2 на позиции 10 (с разрывом)
		memory[9] = 3
		memory[10], memory[11], memory[12] = 20, 21, 22
		ptr2 := unsafe.Pointer(&memory[10])

		pointers := []unsafe.Pointer{ptr2, ptr1} // неотсортированные

		Defragment(memory, pointers)

		// Блок 1 должен остаться на 1
		if uintptr(pointers[1])-basePtr != 1 {
			t.Error("block 1 moved incorrectly")
		}

		// Блок 2 должен переместиться на 1+2+1=4
		newPos2 := uintptr(pointers[0]) - basePtr
		if newPos2 != 4 {
			t.Errorf("block 2 should be at position 4, got %d", newPos2)
		}

		// Проверяем заголовок блока 2 на новом месте
		if memory[3] != 3 {
			t.Error("block 2 header wrong at new position")
		}
		// Проверяем данные блока 2
		if memory[4] != 20 || memory[5] != 21 || memory[6] != 22 {
			t.Error("block 2 data corrupted")
		}
	})

	t.Run("three blocks all need moving", func(t *testing.T) {
		memory := make([]byte, 100)
		basePtr := uintptr(unsafe.Pointer(&memory[0]))

		// Все блоки не на своих местах
		// Блок 1: данные с 20
		memory[19] = 2
		memory[20], memory[21] = 1, 2
		ptr1 := unsafe.Pointer(&memory[20])

		// Блок 2: данные с 30
		memory[29] = 1
		memory[30] = 3
		ptr2 := unsafe.Pointer(&memory[30])

		// Блок 3: данные с 40
		memory[39] = 3
		memory[40], memory[41], memory[42] = 4, 5, 6
		ptr3 := unsafe.Pointer(&memory[40])

		pointers := []unsafe.Pointer{ptr1, ptr2, ptr3}

		Defragment(memory, pointers)

		// Ожидаемые позиции:
		// Блок 1 (2 байта): позиция 1
		// Блок 2 (1 байт): позиция 1+2+1=4
		// Блок 3 (3 байта): позиция 4+1+1=6

		checkPos := func(ptr unsafe.Pointer, expected int, msg string) {
			if uintptr(ptr)-basePtr != uintptr(expected) {
				t.Errorf("%s: expected %d, got %d", msg, expected, uintptr(ptr)-basePtr)
			}
		}

		checkPos(pointers[0], 1, "block 1")
		checkPos(pointers[1], 4, "block 2")
		checkPos(pointers[2], 6, "block 3")

		// Проверяем данные
		if memory[1] != 1 || memory[2] != 2 {
			t.Error("block 1 data wrong")
		}
		if memory[4] != 3 {
			t.Error("block 2 data wrong")
		}
		if memory[6] != 4 || memory[7] != 5 || memory[8] != 6 {
			t.Error("block 3 data wrong")
		}
	})

	t.Run("blocks in reverse order", func(t *testing.T) {
		memory := make([]byte, 100)
		basePtr := uintptr(unsafe.Pointer(&memory[0]))

		// Блоки расположены в обратном порядке (большие адреса сначала)
		// Блок 3: данные с 30
		memory[29] = 3
		memory[30], memory[31], memory[32] = 1, 2, 3
		ptr3 := unsafe.Pointer(&memory[30])

		// Блок 2: данные с 20
		memory[19] = 2
		memory[20], memory[21] = 4, 5
		ptr2 := unsafe.Pointer(&memory[20])

		// Блок 1: данные с 10
		memory[9] = 1
		memory[10] = 6
		ptr1 := unsafe.Pointer(&memory[10])

		pointers := []unsafe.Pointer{ptr1, ptr2, ptr3}

		Defragment(memory, pointers)

		// Проверяем что порядок сохранился в соответствии с исходными размерами
		// Но теперь они должны быть упакованы по возрастанию адресов

		// Должны быть упакованы: блок 1 (1 байт), блок 2 (2 байта), блок 3 (3 байта)
		// Позиции: блок1:1, блок2:1+1+1=3, блок3:3+2+1=6
		pos1 := uintptr(pointers[0]) - basePtr
		pos2 := uintptr(pointers[1]) - basePtr
		pos3 := uintptr(pointers[2]) - basePtr

		if pos1 != 1 || pos2 != 3 || pos3 != 6 {
			t.Errorf("wrong positions: %d, %d, %d", pos1, pos2, pos3)
		}

		// Проверяем что данные на месте
		if memory[1] != 6 { // блок 1
			t.Error("block 1 data wrong")
		}
		if memory[3] != 4 || memory[4] != 5 { // блок 2
			t.Error("block 2 data wrong")
		}
		if memory[6] != 1 || memory[7] != 2 || memory[8] != 3 { // блок 3
			t.Error("block 3 data wrong")
		}
	})

	t.Run("max size blocks (255 bytes)", func(t *testing.T) {
		memory := make([]byte, 1000)
		basePtr := uintptr(unsafe.Pointer(&memory[0]))

		// Создаем блок максимального размера
		memory[0] = 255
		for i := 1; i <= 255; i++ {
			memory[i] = byte(i % 256)
		}
		ptr1 := unsafe.Pointer(&memory[1])

		// Еще один блок с разрывом
		memory[300] = 255
		for i := 301; i <= 555; i++ {
			memory[i] = byte((i - 300) % 256)
		}
		ptr2 := unsafe.Pointer(&memory[301])

		pointers := []unsafe.Pointer{ptr1, ptr2}

		Defragment(memory, pointers)

		// Второй блок должен переместиться сразу за первым
		// Позиция для второго блока: 1 (данные первого) + 255 + 1 (заголовок) = 257
		newPos2 := uintptr(pointers[1]) - basePtr
		if newPos2 != 257 {
			t.Errorf("max size block should move to position 257, got %d", newPos2)
		}

		// Проверяем что заголовок на месте
		if memory[256] != 255 {
			t.Error("header for second block wrong")
		}
	})

	t.Run("multiple blocks with some already optimal", func(t *testing.T) {
		memory := make([]byte, 100)
		basePtr := uintptr(unsafe.Pointer(&memory[0]))

		// Блок 1 на оптимальной позиции (1)
		memory[0] = 2
		memory[1], memory[2] = 1, 2
		ptr1 := unsafe.Pointer(&memory[1])

		// Блок 2 с небольшим разрывом (позиция 5, оптимально было бы 4)
		memory[4] = 1
		memory[5] = 3
		ptr2 := unsafe.Pointer(&memory[5])

		// Блок 3 с большим разрывом
		memory[20] = 3
		memory[21], memory[22], memory[23] = 4, 5, 6
		ptr3 := unsafe.Pointer(&memory[21])

		pointers := []unsafe.Pointer{ptr1, ptr2, ptr3}

		Defragment(memory, pointers)

		// Блок 1 должен остаться
		if uintptr(pointers[0])-basePtr != 1 {
			t.Error("optimal block moved")
		}

		// Блок 2 должен переместиться на 1+2+1=4
		if uintptr(pointers[1])-basePtr != 4 {
			t.Error("block 2 not moved to correct position")
		}

		// Блок 3 должен переместиться на 4+1+1=6
		if uintptr(pointers[2])-basePtr != 6 {
			t.Error("block 3 not moved to correct position")
		}
	})

	t.Run("zero-length blocks not possible due to header", func(t *testing.T) {
		// В текущей реализации блок с длиной 0 невозможен,
		// но проверим что код корректен если такой появится
		memory := make([]byte, 100)
		basePtr := uintptr(unsafe.Pointer(&memory[0]))

		memory[0] = 0 // заголовок с длиной 0
		ptr := unsafe.Pointer(&memory[1])

		pointers := []unsafe.Pointer{ptr}

		Defragment(memory, pointers)

		// Блок с длиной 0 должен остаться на месте (позиция 1)
		if uintptr(pointers[0])-basePtr != 1 {
			t.Error("zero-length block moved")
		}
	})

	t.Run("pointers array modified correctly", func(t *testing.T) {
		memory := make([]byte, 100)
		basePtr := uintptr(unsafe.Pointer(&memory[0]))

		// Создаем блоки
		memory[0] = 2
		memory[1], memory[2] = 1, 2
		ptr1 := unsafe.Pointer(&memory[1])

		memory[9] = 3
		memory[10], memory[11], memory[12] = 3, 4, 5
		ptr2 := unsafe.Pointer(&memory[10])

		originalPtrs := []unsafe.Pointer{ptr1, ptr2}
		pointers := make([]unsafe.Pointer, len(originalPtrs))
		copy(pointers, originalPtrs)

		Defragment(memory, pointers)

		// Проверяем что указатели изменились (кроме первого)
		if pointers[0] == originalPtrs[0] && pointers[1] == originalPtrs[1] {
			t.Error("pointers array not updated")
		}

		// Проверяем что указатели указывают в нашу память
		for i, p := range pointers {
			pos := uintptr(p) - basePtr
			if pos < 0 || pos >= uintptr(len(memory)) {
				t.Errorf("pointer %d points outside memory", i)
			}
		}
	})
}

// Тесты на паники (негативные сценарии)
func TestDefragmentPanics(t *testing.T) {
	t.Run("panic on invalid pointer - before memory start", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Error("should panic on pointer before memory start")
			}
		}()

		memory := make([]byte, 100)
		// Создаем указатель на 1 байт до начала memory
		invalidPtr := unsafe.Pointer(uintptr(unsafe.Pointer(&memory[0])) - 1)
		pointers := []unsafe.Pointer{invalidPtr}

		Defragment(memory, pointers)
	})

	t.Run("panic on invalid pointer - after memory end", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Error("should panic on pointer after memory end")
			}
		}()

		memory := make([]byte, 100)
		// Создаем указатель за пределами memory
		invalidPtr := unsafe.Pointer(uintptr(unsafe.Pointer(&memory[99])) + 2)
		pointers := []unsafe.Pointer{invalidPtr}

		Defragment(memory, pointers)
	})

	t.Run("panic on block extending beyond memory", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Error("should panic on block extending beyond memory")
			}
		}()

		memory := make([]byte, 100)
		basePtr := uintptr(unsafe.Pointer(&memory[0]))

		// Заголовок говорит что блок длиной 10, но он расположен в конце
		memory[95] = 10                     // заголовок
		ptr := unsafe.Pointer(basePtr + 96) // данные с 96 по 106, но memory только до 99

		pointers := []unsafe.Pointer{ptr}

		Defragment(memory, pointers)
	})

	t.Run("panic on pointer at position 0", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Error("should panic on pointer at position 0")
			}
		}()

		memory := make([]byte, 100)
		basePtr := uintptr(unsafe.Pointer(&memory[0]))

		// Указатель на позицию 0 (заголовок, не данные)
		ptr := unsafe.Pointer(basePtr)

		pointers := []unsafe.Pointer{ptr}

		Defragment(memory, pointers)
	})
}
