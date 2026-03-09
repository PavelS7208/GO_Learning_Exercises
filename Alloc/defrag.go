package main

import (
	"fmt"
	"sort"
	"unsafe"
)

// blockInfo хранит метаданные о блоке для дефрагментации
type blockInfo struct {
	origIdx   int // Индекс в исходном списке рабочих pointers
	headerPos int // Позиция байта-заголовка в нашем heap (где хранится размер)
	dataPos   int // Позиция начала данных в нашем heap
	length    int // Длина (размер в байтах) блока аллоцированных данных в нашем heap (из заголовка)
}

// Defragment упаковывает блоки памяти плотно друг к другу.
// Устраняет дефрагментацию между блоками, переносит (изменяет значения указателей)
// Формат блока: [1 байт заголовка: длина][N байт данных (до 256 байт)]
// Алгоритм двухпроходный: сначала расчёт позиций что куда переносить, потом копирование.
// Аргументы: memory - наша heap, pointers - указатели на аллоцированные участки(блоки) в нашем heap
func Defragment(memory []byte, pointers []unsafe.Pointer) {
	if len(memory) == 0 || len(pointers) == 0 {
		return
	}
	// Начало нашего heap (адрес нулевого байта)
	basePtr := uintptr(unsafe.Pointer(&memory[0]))
	// Статистика аллокаций
	blocks := make([]blockInfo, 0, len(pointers))

	// Собираем блоки и проверяем все что можем
	for i, ptr := range pointers {
		if ptr == nil {
			continue
		}

		// Конвертируем указатель в индекс от начала нашего heap (безопасно, т.к. память то наша)
		// При проверках panic так как учебная задачка
		absPos := int(uintptr(ptr) - basePtr)
		if absPos <= 0 || absPos >= len(memory) {
			panic(fmt.Sprintf("invalid pointer at index %d: pos=%d", i, absPos))
		}

		headerPos := absPos - 1
		length := int(memory[headerPos]) // размер хранится в байте (uint8), т.е. макс. 255 байт длина блока

		// Валидация границ блока, что у нас не передался мусор в ptr
		if absPos+length > len(memory) {
			panic(fmt.Sprintf("block at %d extends beyond memory: %d+%d > %d",
				absPos, absPos, length, len(memory)))
		}

		blocks = append(blocks, blockInfo{
			origIdx:   i,
			headerPos: headerPos,
			dataPos:   absPos,
			length:    length,
		})
	}

	// Ничего не нужно дефрагментировать, указателей для перемещения нет
	if len(blocks) == 0 {
		return
	}

	// Сортируем нашу статистику по исходной позиции для упрощения порядка перестановок при дефрагментации
	sort.Slice(blocks, func(i, j int) bool {
		return blocks[i].dataPos < blocks[j].dataPos
	})

	// Находим что куда нужно переместить (без самого пока переноса)
	// Используем -1 как значение "не перемещается"
	const NotNeedToMoved int = -1

	newPos := make([]int, len(memory))
	for i := range newPos {
		newPos[i] = NotNeedToMoved
	}

	currentPos := 1 // первая позиция данных (после заголовка первого блока)

	for _, block := range blocks {
		// Проверяем: если блок уже на нужной позиции — пропускаем, если нет - метим что будем его двигать
		if block.dataPos != currentPos {
			newPos[block.dataPos] = currentPos
		}
		currentPos += 1 + block.length // заголовок + данные
	}

	// А теперь двигаем туда куда нужно сами данные в heap
	// Выделяем буфер ОДИН РАЗ перед циклом — максимум 255 байт + 1 на всякий случай
	tempBuf := make([]byte, 256)

	for _, block := range blocks {
		targetDataPos := newPos[block.dataPos]
		if targetDataPos == NotNeedToMoved { // блок не переносится
			continue
		}
		// Копируем текущий блок в буфер
		buf := tempBuf[:block.length]
		copy(buf, memory[block.dataPos:block.dataPos+block.length])

		// Записываем заголовок и данные на новую позицию и меняем указатель
		//  В идеале - транзакция на эти операции, чтобы было атомарно
		memory[targetDataPos-1] = byte(block.length)
		copy(memory[targetDataPos:targetDataPos+block.length], buf)
		offset := uintptr(targetDataPos)
		pointers[block.origIdx] = unsafe.Pointer(basePtr + offset)
	}

	// Если нужно обнулить оставшуюся heap после последнего блока:
	// currentPos указывает за последний блок
	//if currentPos < len(memory) {
	//     clear(memory[currentPos:])
	// }

}
