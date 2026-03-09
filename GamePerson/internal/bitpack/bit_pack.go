// Package bitpack предоставляет эффективную упаковку данных в битовые поля.
//
// # Checked vs Unchecked API
//
// Каждая Set операция имеет две версии:
//
//   - Checked (SetXXXFieldAs): Валидирует все параметры, возвращает error.
//     Использовать для пользовательского ввода и внешних данных.
//
//   - Unchecked (SetXXXFieldUncheckedAs): Пропускает валидацию для производительности.
//     Использовать ТОЛЬКО когда данные уже валидированы или из доверенных источников.
//
// # Требования для Unchecked версий
//
// Вызывающий код ОБЯЗАН гарантировать:
//   - len(packed) > 0 и len(packed) <= 8
//   - value находится в допустимом диапазоне для поля
//
// Нарушение этих требований приводит к undefined behavior.
//
// # Пример использования
//
//   // Пользовательский ввод - используем checked
//   if err := bitpack.SetUIntFieldAs(packed[:], healthField, userInput); err != nil {
//       return fmt.Errorf("invalid health: %w", err)
//   }
//
//   // Batch обработка валидированных данных - используем unchecked
//   for _, player := range validatedPlayers {
//       bitpack.SetUIntFieldUncheckedAs(packed[:], healthField, player.Health)
//   }

package bitpack

import (
	"math"
)

// ==================== Простые утилиты ====================

func UnpackBytes(data []byte) BitSet64 {
	var result BitSet64
	for i, b := range data {
		result |= BitSet64(b) << (8 * i)
	}
	return result
}

func PackBytes(data []byte, value BitSet64) {
	for i := range data {
		data[i] = byte(value >> (8 * i))
	}
}

// ==================== Get ====================

func GetUIntFieldAs[T UnsignedInteger](packed []byte, field UIntBitField) T {
	return T(field.Get(UnpackBytes(packed)))
}

func GetIntFieldAs[T SignedInteger](packed []byte, field IntBitField) T {
	return T(field.Get(UnpackBytes(packed)))
}

func GetBoolField(packed []byte, field BoolBitField) bool {
	return field.Get(UnpackBytes(packed))
}

// ==================== Set Unchecked версии ====================

func SetUIntFieldUncheckedAs[T UnsignedInteger](packed []byte, field UIntBitField, value T) {
	bits := field.UpdateUnchecked(UnpackBytes(packed), uint64(value))
	PackBytes(packed, bits)
}

func SetIntFieldUncheckedAs[T SignedInteger](packed []byte, field IntBitField, value T) {
	bits := field.UpdateUnchecked(UnpackBytes(packed), int64(value))
	PackBytes(packed, bits)
}

func SetBoolFieldUnchecked(packed []byte, field BoolBitField, value bool) {
	bits := UnpackBytes(packed) // Сначала читаем текущее состояние
	if value {
		bits = field.Set(bits)
	} else {
		bits = field.Clear(bits)
	}
	PackBytes(packed, bits)
}

// ==================== Set Checked версии ====================

func SetUIntFieldAs[T UnsignedInteger](packed []byte, field UIntBitField, value T) error {
	if len(packed) == 0 {
		return newSliceEmptyError()
	}
	if len(packed) > 8 {
		return newSliceTooLargeError(len(packed))
	}

	bits, err := field.Update(UnpackBytes(packed), uint64(value))
	if err != nil {
		return err
	}
	PackBytes(packed, bits)
	return nil
}

func SetIntFieldAs[T SignedInteger](packed []byte, field IntBitField, value T) error {
	if len(packed) == 0 {
		return newSliceEmptyError()
	}
	if len(packed) > 8 {
		return newSliceTooLargeError(len(packed))
	}

	bits, err := field.Update(UnpackBytes(packed), int64(value))
	if err != nil {
		return err
	}
	PackBytes(packed, bits)
	return nil
}

func SetBoolField(packed []byte, field BoolBitField, value bool) error {
	if len(packed) == 0 {
		return newSliceEmptyError()
	}
	if len(packed) > 8 {
		return newSliceTooLargeError(len(packed))
	}

	SetBoolFieldUnchecked(packed, field, value)
	return nil
}

// ==================== Обратная совместимость (deprecated) ====================
// Оставляем старые функции для плавной миграции, но помечаем как deprecated

// Deprecated: Use GetUIntFieldAs[uint64] instead
func GetUIntField(packed []byte, field UIntBitField) uint64 {
	return GetUIntFieldAs[uint64](packed, field)
}

// Deprecated: Use SetUIntFieldAs[uint64] instead
func SetUIntField(packed []byte, field UIntBitField, value uint64) error {
	return SetUIntFieldAs[uint64](packed, field, value)
}

// Deprecated: Use GetIntFieldAs[int64] instead
func GetIntField(packed []byte, field IntBitField) int64 {
	return GetIntFieldAs[int64](packed, field)
}

// Deprecated: Use SetIntFieldAs[int64] instead
func SetIntField(packed []byte, field IntBitField, value int64) error {
	return SetIntFieldAs[int64](packed, field, value)
}

// / --------------- Доп методы ------------------------------------
// computeMask возвращает битовую маску для заданной ширины.
// ТРЕБОВАНИЕ: width должна быть <= 64, иначе результат некорректен.
// Вызывающая функция ОБЯЗАНА проверить width перед вызовом.
func computeMask(width uint8) uint64 {
	if width >= 64 {
		return math.MaxUint64
	}
	return (uint64(1) << width) - 1
}

func maxAllowedForWidth(width uint8) uint64 {
	if width >= 64 {
		return math.MaxUint64
	}
	return (uint64(1) << width) - 1
}
