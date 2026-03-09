package bitpack

import (
	"fmt"
	"math"
)

// =================  IntBitField ===========================================
// ================ Хранение знаковых целых ================================

type IntBitField struct {
	Start BitPosition // Начальная позиция бита (включительно)
	End   BitPosition // Конечная позиция бита (включительно)
	Min   int64       // Минимальное допустимое значение
	Max   int64       // Максимальное допустимое значение
	mask  uint64      // Кэшированная маска для извлечения битов
}

// NewIntBitField Конструктор создаёт битовое поле для знаковых целых с валидацией диапазона
func NewIntBitField(start, end BitPosition, min, max int64) (IntBitField, error) {
	if start > end {
		return IntBitField{}, newStartAfterEndError(start, end)
	}
	if end >= 64 {
		return IntBitField{}, newEndOutOfRangeError(end)
	}

	width := end - start + 1
	allowedMin, allowedMax := intRangeForWidth(width)

	// Валидация пользовательского диапазона
	if min < allowedMin || max > allowedMax {
		return IntBitField{}, newValueOutOfRangeError(min, max, allowedMin, allowedMax, width)
	}
	if min > max {
		return IntBitField{}, newValueRangeInvertedError(min, max)
	}

	mask := computeMask(width)

	return IntBitField{
		Start: start,
		End:   end,
		Min:   min,
		Max:   max,
		mask:  mask,
	}, nil
}

// MustNewIntBitField создаёт битовое поле или паникует при ошибке конфигурации
// Используется ТОЛЬКО для статических конфигураций

func MustNewIntBitField(start, end BitPosition, min, max int64) IntBitField {
	bf, err := NewIntBitField(start, end, min, max)
	if err != nil {
		panic(fmt.Sprintf("FATAL: invalid static int bit field configuration [start=%d, end=%d, min=%d, max=%d]: %v",
			start, end, min, max, err))
	}
	return bf
}

// NewIntBitFieldAuto создаёт поле с полным диапазоном для заданной ширины
func NewIntBitFieldAuto(start, end BitPosition) (IntBitField, error) {
	width := end - start + 1
	minValue, maxValue := intRangeForWidth(width)
	return NewIntBitField(start, end, minValue, maxValue)
}

// Get извлекает знаковое целое из bitSet с корректным знаковым расширением.
// Использует арифметический сдвиг для восстановления знака.
func (bf IntBitField) Get(bitSet BitSet64) int64 {
	// Извлекаем биты как беззнаковое число
	raw := (uint64(bitSet) >> bf.Start) & bf.mask

	// Знаковое расширение (sign extension) через арифметический сдвиг
	width := bf.Width()
	if width == 64 {
		return int64(raw)
	}

	// Сдвигаем влево, чтобы знаковый бит стал старшим битом uint64,
	// затем арифметически сдвигаем вправо для восстановления знака
	shift := 64 - width
	return int64(raw<<shift) >> shift
}

// Update записывает знаковое значение в bitSet с валидацией диапазона
func (bf IntBitField) Update(bitSet BitSet64, value int64) (BitSet64, error) {
	if value < bf.Min || value > bf.Max {
		return bitSet, newValueOutOfRangeError(value, value, bf.Min, bf.Max, bf.Width())
	}
	return bf.UpdateUnchecked(bitSet, value), nil

}

func (bf IntBitField) UpdateUnchecked(bitSet BitSet64, value int64) BitSet64 {

	// Преобразуем в беззнаковое представление (дополнительный код)
	// При записи в битовое поле отрицательные числа автоматически
	// представляются корректно благодаря приведению к uint64
	unsignedValue := uint64(value) & bf.mask
	fieldMask := bf.mask << bf.Start
	newBits := (uint64(bitSet) & ^fieldMask) | (unsignedValue << bf.Start)
	return BitSet64(newBits)
}

// ------------- Сервисные методы --------------------------------

// Диапазон допустимых значений для знакового целого заданной ширины
func intRangeForWidth(width uint8) (min, max int64) {
	if width >= 64 {
		return math.MinInt64, math.MaxInt64
	}
	signBit := int64(1) << (width - 1)
	return -signBit, signBit - 1
}

// Width - Ширина битового поля
func (bf IntBitField) Width() uint8 {
	return bf.End - bf.Start + 1
}

// Строковое представление для отладки
func (bf IntBitField) String() string {
	return fmt.Sprintf("IntBitField[%d:%d] range=[%d,%d]", bf.Start, bf.End, bf.Min, bf.Max)
}

// =================  UIntBitField  =========================================
//  ================ Хранение беззнаковых целых =============================

type UIntBitField struct {
	Start BitPosition // Начальная позиция бита в структуре (включительно)
	End   BitPosition // Конечная позиция бита в структуре (включительно)
	Max   uint64      // Максимальное значение для хранения в структуре битов
	mask  uint64      // Кэшированная маска для производительности
}

// Конструктор NewUIntBitField создаёт битовое поле с проверкой на ошибки логики хранения

func NewUIntBitField(start, end BitPosition, max uint64) (UIntBitField, error) {
	if start > end {
		return UIntBitField{}, newStartAfterEndError(start, end)
	}
	if end >= 64 {
		return UIntBitField{}, newEndOutOfRangeError(end)
	}

	width := end - start + 1
	allowedMax := maxAllowedForWidth(width)
	if max > allowedMax {
		return UIntBitField{}, newValueOverflowError(max, allowedMax, width)
	}

	// Кэшируем маску при создании — поля иммутабельны
	mask := computeMask(width)

	return UIntBitField{
		Start: start,
		End:   end,
		Max:   max,
		mask:  mask,
	}, nil
}

// MustNewUIntBitField создаёт битовое поле или паникует при ошибке конфигурации
// используется ТОЛЬКО для статических конфигураций, проверенных на этапе разработки

func MustNewUIntBitField(start, end BitPosition, max uint64) UIntBitField {
	bf, err := NewUIntBitField(start, end, max)
	if err != nil {
		panic(fmt.Sprintf("FATAL: invalid static bit field configuration [start=%d, end=%d, max=%d]: %v",
			start, end, max, err))
	}
	return bf
}

// Get - извлекает из указанного bitSet значение целого типа согласно конфигурации поля

func (bf UIntBitField) Get(bitSet BitSet64) uint64 {
	return (uint64(bitSet) >> bf.Start) & bf.mask
}

// Update - Возвращает новое значение 64 битной карты, после изменения значения

func (bf UIntBitField) Update(bitSet BitSet64, value uint64) (BitSet64, error) {
	if value > bf.Max {
		return bitSet, newValueOverflowError(value, bf.Max, bf.Width())
	}
	return bf.UpdateUnchecked(bitSet, value), nil
}

func (bf UIntBitField) UpdateUnchecked(bitSet BitSet64, value uint64) BitSet64 {
	fieldMask := bf.mask << bf.Start
	return BitSet64((uint64(bitSet) & ^fieldMask) | ((value & bf.mask) << bf.Start))
}

// ------------- Сервисные методы --------------------------------

func (bf UIntBitField) Width() uint8 {
	return bf.End - bf.Start + 1
}

// Строковое представление для отладки
func (bf UIntBitField) String() string {
	return fmt.Sprintf("UIntBitField[%d:%d] max=%d", bf.Start, bf.End, bf.Max)
}

// =================  BoolBitField =====================================================
//  ================ Хранение булевой величины в одном бите =============================

type BoolBitField struct {
	Position BitPosition // Начальная позиция бита в структуре (включительно)
	bitMask  uint64      // кэшированная маска: 1 << Position
}

func NewBoolBitField(pos BitPosition) (BoolBitField, error) {
	if pos >= 64 {
		return BoolBitField{}, newPositionOutOfRangeError(pos)
	}
	return BoolBitField{
		Position: pos,
		bitMask:  uint64(1) << pos,
	}, nil
}

// MustNewBoolBitField создаёт булево битовое поле и паникует при ошибках
// Используется ТОЛЬКО для статических конфигураций, проверенных на этапе разработки

func MustNewBoolBitField(pos BitPosition) BoolBitField {
	bf, err := NewBoolBitField(pos)
	if err != nil {
		panic(fmt.Sprintf("FATAL: invalid static bool bit field configuration [position=%d]: %v",
			pos, err))
	}
	return bf
}

// Get извлекает значение из указанного bitSet значение bool типа из согласно конфигурации битов
func (bf BoolBitField) Get(bitSet BitSet64) bool {
	return (uint64(bitSet) & bf.bitMask) != 0
}

// Изменяет значение бита для булева поля и возвращает новое
// Ошибок быть не может (если конфигурация уже валидная), но оставлено для соответствия интерфейсу и расширения

func (bf BoolBitField) Update(bitSet BitSet64, value bool) (BitSet64, error) {
	return bf.UpdateUnchecked(bitSet, value), nil
}

func (bf BoolBitField) UpdateUnchecked(bitSet BitSet64, value bool) BitSet64 {
	if value {
		return BitSet64(uint64(bitSet) | bf.bitMask)
	}
	return BitSet64(uint64(bitSet) &^ bf.bitMask)
}

// -------   Доп методы для bool поля: установить(true), clear (false), переключить

func (bf BoolBitField) Set(bitSet BitSet64) BitSet64 {
	return BitSet64(uint64(bitSet) | bf.bitMask)
}

func (bf BoolBitField) Clear(bitSet BitSet64) BitSet64 {
	return BitSet64(uint64(bitSet) &^ bf.bitMask)
}

func (bf BoolBitField) Toggle(bitSet BitSet64) BitSet64 {
	return BitSet64(uint64(bitSet) ^ bf.bitMask)
}

// Строковое представление для отладки
func (bf BoolBitField) String() string {
	return fmt.Sprintf("\"BoolBitField[%d] ", bf.Position)
}
