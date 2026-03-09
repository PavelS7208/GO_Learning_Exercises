package bitpack

import (
	"math"
	"testing"
)

// ============ Критически важный тест: проверка покрытия и непересечения битов ============
// TestBitLayoutCoverage проверяет корректность раскладки битов для сложных структур:
// 1. Непересечение — поля не должны перекрываться
// 2. Полное покрытие — все биты должны быть назначены полям
// 3. Корректность границ — поля должны находиться в пределах структуры

// Тестовая структура: компактный заголовок сетевого пакета (32 бита)
var (
	VersionField  = MustNewUIntBitField(0, 3, 15)     // 4 бита: версия (0-15)
	PriorityField = MustNewUIntBitField(4, 6, 7)      // 3 бита: приоритет (0-7)
	EncryptedFlag = MustNewBoolBitField(7)            // 1 бит: шифрование
	SequenceField = MustNewUIntBitField(8, 23, 65535) // 16 бит: номер последовательности
	ReservedField = MustNewUIntBitField(24, 30, 127)  // 7 бит: зарезервировано
	LastFlag      = MustNewBoolBitField(31)           // 1 бит: флаг последнего пакета
)

// TestPacketHeaderLayout проверяет непересечение и полное покрытие битов
func TestPacketHeaderLayout(t *testing.T) {
	fields := []struct {
		name  string
		start BitPosition
		end   BitPosition
	}{
		{"Version", VersionField.Start, VersionField.End},
		{"Priority", PriorityField.Start, PriorityField.End},
		{"Encrypted", EncryptedFlag.Position, EncryptedFlag.Position},
		{"Sequence", SequenceField.Start, SequenceField.End},
		{"Reserved", ReservedField.Start, ReservedField.End},
		{"LastFlag", LastFlag.Position, LastFlag.Position},
	}

	const totalBits = 32
	coverage := make([]bool, totalBits)

	// Проверяем каждое поле и отмечаем покрытые биты
	for _, f := range fields {
		t.Run(f.name, func(t *testing.T) {
			// Валидация границ поля
			if f.start > f.end || f.end >= totalBits {
				t.Errorf("поле %s выходит за пределы структуры: [%d:%d]", f.name, f.start, f.end)
			}

			// Проверка на пересечение с уже покрытыми битами
			for pos := f.start; pos <= f.end; pos++ {
				if coverage[pos] {
					t.Errorf("бит %d уже занят другим полем (конфликт в поле %s)", pos, f.name)
				}
				coverage[pos] = true
			}
		})
	}

	// Проверка полного покрытия
	missing := 0
	for i, covered := range coverage {
		if !covered {
			t.Errorf("бит %d не назначен ни одному полю", i)
			missing++
		}
	}
	if missing > 0 {
		t.Fatalf("обнаружено %d непокрытых битов из %d", missing, totalBits)
	}

	//	t.Logf("✅ Все %d бита корректно распределены между %d полями без пересечений", totalBits, len(fields))
}

// Аналогично, го только со знаковыми полями

func TestBitLayoutCoverage(t *testing.T) {
	// Пример: компактная структура данных для метеостанции (32 бита)
	// Биты: 31..0
	// [31]       : флаг ошибки датчика
	// [30..24]   : зарезервировано (7 бит)
	// [23..16]   : давление (отклонение от нормы, знаковое, 8 бит: -128..127)
	// [15..8]    : влажность (беззнаковое, 8 бит: 0..100)
	// [7..0]     : температура (знаковое, 8 бит: -50..100)

	// Определение полей структуры
	fields := []struct {
		name     string
		start    BitPosition
		end      BitPosition
		isSigned bool
		minValue int64  // для signed
		maxValue int64  // для signed
		maxUint  uint64 // для unsigned
	}{
		{"ErrorFlag", 31, 31, false, 0, 0, 1},
		{"Reserved", 24, 30, false, 0, 0, 127},
		{"PressureDelta", 16, 23, true, -128, 127, 0},
		{"Humidity", 8, 15, false, 0, 0, 100},
		{"Temperature", 0, 7, true, -50, 100, 0},
	}

	const totalBits = 32
	coverage := make([]bool, totalBits)
	conflicts := make(map[int][]string)

	// Проверка каждого поля и построение карты покрытия
	for _, f := range fields {
		t.Run("field_validation/"+f.name, func(t *testing.T) {
			// Валидация границ поля
			if f.start > f.end {
				t.Errorf("поле %s: start (%d) > end (%d)", f.name, f.start, f.end)
			}
			if f.end >= totalBits {
				t.Errorf("поле %s выходит за пределы структуры: [%d:%d] >= %d", f.name, f.start, f.end, totalBits)
			}

			// Проверка на пересечение с уже покрытыми битами
			for pos := f.start; pos <= f.end; pos++ {
				if pos >= totalBits {
					continue
				}
				if coverage[pos] {
					conflicts[int(pos)] = append(conflicts[int(pos)], f.name)
				} else {
					coverage[pos] = true
				}
			}
		})
	}

	// Проверка пересечений
	if len(conflicts) > 0 {
		t.Error("Обнаружены пересекающиеся биты между полями:")
		for bit, fieldNames := range conflicts {
			t.Errorf("  бит %d занят полями: %v", bit, fieldNames)
		}
		t.FailNow()
	}

	// Проверка полного покрытия
	missing := 0
	for i := 0; i < totalBits; i++ {
		if !coverage[i] {
			t.Errorf("бит %d не назначен ни одному полю", i)
			missing++
		}
	}
	if missing > 0 {
		t.Fatalf("ОШИБКА: обнаружено %d непокрытых битов из %d", missing, totalBits)
	}

	// Дополнительная проверка: все поля корректно создаются
	for _, f := range fields {
		t.Run("field_creation/"+f.name, func(t *testing.T) {
			if f.isSigned {
				_, err := NewIntBitField(f.start, f.end, f.minValue, f.maxValue)
				if err != nil {
					t.Errorf("не удалось создать знаковое поле %s: %v", f.name, err)
				}
			} else {
				_, err := NewUIntBitField(f.start, f.end, f.maxUint)
				if err != nil {
					t.Errorf("не удалось создать беззнаковое поле %s: %v", f.name, err)
				}
			}
		})
	}

	//t.Logf("✅ Все %d бита корректно распределены между %d полями без пересечений", totalBits, len(fields))
}

func TestGenericGetSetUInt(t *testing.T) {
	field, _ := NewUIntBitField(0, 9, 1000)
	packed := make([]byte, 8)

	// Тест с uint32
	err := SetUIntFieldAs[uint32](packed, field, 500)
	if err != nil {
		t.Fatalf("SetUIntFieldAs failed: %v", err)
	}

	result := GetUIntFieldAs[uint32](packed, field)
	if result != 500 {
		t.Errorf("Expected 500, got %d", result)
	}

	// Тест с uint16
	err = SetUIntFieldAs[uint16](packed, field, 250)
	if err != nil {
		t.Fatalf("SetUIntFieldAs failed: %v", err)
	}

	result16 := GetUIntFieldAs[uint16](packed, field)
	if result16 != 250 {
		t.Errorf("Expected 250, got %d", result16)
	}
}

func TestGenericGetSetInt(t *testing.T) {
	field, _ := NewIntBitField(0, 9, -500, 500)
	packed := make([]byte, 8)

	// Тест с int32
	err := SetIntFieldAs[int32](packed, field, -250)
	if err != nil {
		t.Fatalf("SetIntFieldAs failed: %v", err)
	}

	result := GetIntFieldAs[int32](packed, field)
	if result != -250 {
		t.Errorf("Expected -250, got %d", result)
	}

	// Тест с int16
	err = SetIntFieldAs[int16](packed, field, 100)
	if err != nil {
		t.Fatalf("SetIntFieldAs failed: %v", err)
	}

	result16 := GetIntFieldAs[int16](packed, field)
	if result16 != 100 {
		t.Errorf("Expected 100, got %d", result16)
	}
}

// ============ Тесты для UnpackBytes и PackBytes ============

func TestUnpackBytes(t *testing.T) {
	tests := []struct {
		name  string
		input []byte
		want  BitSet64
	}{
		{
			name:  "single byte",
			input: []byte{0xFF},
			want:  BitSet64(0xFF),
		},
		{
			name:  "two bytes little-endian",
			input: []byte{0x34, 0x12},
			want:  BitSet64(0x1234),
		},
		{
			name:  "four bytes",
			input: []byte{0x78, 0x56, 0x34, 0x12},
			want:  BitSet64(0x12345678),
		},
		{
			name:  "eight bytes full",
			input: []byte{0xF0, 0xDE, 0xBC, 0x9A, 0x78, 0x56, 0x34, 0x12},
			want:  BitSet64(0x123456789ABCDEF0),
		},
		{
			name:  "all zeros",
			input: []byte{0, 0, 0, 0, 0, 0, 0, 0},
			want:  BitSet64(0),
		},
		{
			name:  "all ones",
			input: []byte{0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF},
			want:  BitSet64(math.MaxUint64),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := UnpackBytes(tt.input)
			if got != tt.want {
				t.Errorf("UnpackBytes() = 0x%X, want 0x%X", uint64(got), uint64(tt.want))
			}
		})
	}
}

func TestPackBytes(t *testing.T) {
	tests := []struct {
		name  string
		value BitSet64
		size  int
		want  []byte
	}{
		{
			name:  "single byte",
			value: BitSet64(0xFF),
			size:  1,
			want:  []byte{0xFF},
		},
		{
			name:  "two bytes little-endian",
			value: BitSet64(0x1234),
			size:  2,
			want:  []byte{0x34, 0x12},
		},
		{
			name:  "four bytes",
			value: BitSet64(0x12345678),
			size:  4,
			want:  []byte{0x78, 0x56, 0x34, 0x12},
		},
		{
			name:  "eight bytes full",
			value: BitSet64(0x123456789ABCDEF0),
			size:  8,
			want:  []byte{0xF0, 0xDE, 0xBC, 0x9A, 0x78, 0x56, 0x34, 0x12},
		},
		{
			name:  "partial pack - only use first 3 bytes",
			value: BitSet64(0x123456789ABCDEF0),
			size:  3,
			want:  []byte{0xF0, 0xDE, 0xBC},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := make([]byte, tt.size)
			PackBytes(got, tt.value)

			if len(got) != len(tt.want) {
				t.Fatalf("length mismatch: got %d, want %d", len(got), len(tt.want))
			}

			for i := range got {
				if got[i] != tt.want[i] {
					t.Errorf("byte[%d] = 0x%X, want 0x%X", i, got[i], tt.want[i])
				}
			}
		})
	}
}

func TestUnpackPackRoundTrip(t *testing.T) {
	// Проверяем, что UnpackBytes и PackBytes обратимы
	tests := []BitSet64{
		0,
		1,
		0xFF,
		0x1234,
		0x12345678,
		0x123456789ABCDEF0,
		BitSet64(math.MaxUint64),
	}

	for _, value := range tests {
		t.Run(string(rune('0'+int(value))), func(t *testing.T) {
			packed := make([]byte, 8)
			PackBytes(packed, value)
			unpacked := UnpackBytes(packed)

			if unpacked != value {
				t.Errorf("round-trip failed: original 0x%X, got 0x%X", uint64(value), uint64(unpacked))
			}
		})
	}
}

// ============ Тесты для GetUIntFieldAs ============

func TestGetUIntFieldAsVariousTypes(t *testing.T) {
	field := MustNewUIntBitField(8, 15, 255)
	packed := make([]byte, 8)

	// Устанавливаем значение 200
	SetUIntFieldUncheckedAs[uint64](packed, field, 200)

	// Проверяем разные типы
	if got := GetUIntFieldAs[uint8](packed, field); got != 200 {
		t.Errorf("GetUIntFieldAs[uint8]() = %d, want 200", got)
	}

	if got := GetUIntFieldAs[uint16](packed, field); got != 200 {
		t.Errorf("GetUIntFieldAs[uint16]() = %d, want 200", got)
	}

	if got := GetUIntFieldAs[uint32](packed, field); got != 200 {
		t.Errorf("GetUIntFieldAs[uint32]() = %d, want 200", got)
	}

	if got := GetUIntFieldAs[uint64](packed, field); got != 200 {
		t.Errorf("GetUIntFieldAs[uint64]() = %d, want 200", got)
	}
}

func TestGetUIntFieldAsEdgeCases(t *testing.T) {
	// Тест с максимальными значениями для разных типов
	tests := []struct {
		name  string
		bits  BitPosition
		value uint64
	}{
		{"uint8 max", 8, 255},
		{"uint16 max", 16, 65535},
		{"uint32 max", 32, 4294967295},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			field := MustNewUIntBitField(0, tt.bits-1, tt.value)
			packed := make([]byte, 8)

			SetUIntFieldUncheckedAs[uint64](packed, field, tt.value)
			got := GetUIntFieldAs[uint64](packed, field)

			if got != tt.value {
				t.Errorf("max value round-trip failed: got %d, want %d", got, tt.value)
			}
		})
	}
}

// ============ Тесты для GetIntFieldAs ============

func TestGetIntFieldAsVariousTypes(t *testing.T) {
	field := MustNewIntBitField(0, 15, -32768, 32767)
	packed := make([]byte, 8)

	testValues := []int64{-100, -1, 0, 1, 100}

	for _, value := range testValues {
		t.Run(string(rune('0'+int(value))), func(t *testing.T) {
			SetIntFieldUncheckedAs[int64](packed, field, value)

			if got := GetIntFieldAs[int8](packed, field); got != int8(value) && value >= math.MinInt8 && value <= math.MaxInt8 {
				t.Errorf("GetIntFieldAs[int8]() = %d, want %d", got, value)
			}

			if got := GetIntFieldAs[int16](packed, field); got != int16(value) {
				t.Errorf("GetIntFieldAs[int16]() = %d, want %d", got, value)
			}

			if got := GetIntFieldAs[int32](packed, field); got != int32(value) {
				t.Errorf("GetIntFieldAs[int32]() = %d, want %d", got, value)
			}

			if got := GetIntFieldAs[int64](packed, field); got != value {
				t.Errorf("GetIntFieldAs[int64]() = %d, want %d", got, value)
			}
		})
	}
}

func TestGetIntFieldAsNegativeValues(t *testing.T) {
	field := MustNewIntBitField(8, 15, -128, 127)
	packed := make([]byte, 8)

	testValues := []int64{-128, -100, -50, -1}

	for _, value := range testValues {
		t.Run(string(rune('0'+int(value))), func(t *testing.T) {
			SetIntFieldUncheckedAs[int64](packed, field, value)
			got := GetIntFieldAs[int64](packed, field)

			if got != value {
				t.Errorf("negative value round-trip failed: stored %d, got %d", value, got)
			}
		})
	}
}

// ============ Тесты для GetBoolField ============

func TestGetBoolField(t *testing.T) {
	field := MustNewBoolBitField(5)

	// Тест false
	packed := make([]byte, 8)
	if got := GetBoolField(packed, field); got {
		t.Error("GetBoolField() = true on zero bytes, want false")
	}

	// Тест true
	SetBoolFieldUnchecked(packed, field, true)
	if got := GetBoolField(packed, field); !got {
		t.Error("GetBoolField() = false after setting true, want true")
	}

	// Тест false после очистки
	SetBoolFieldUnchecked(packed, field, false)
	if got := GetBoolField(packed, field); got {
		t.Error("GetBoolField() = true after setting false, want false")
	}
}

func TestGetBoolFieldMultiplePositions(t *testing.T) {
	packed := make([]byte, 8)

	// Создаём несколько булевых полей
	field0 := MustNewBoolBitField(0)
	field7 := MustNewBoolBitField(7)
	field31 := MustNewBoolBitField(31)
	field63 := MustNewBoolBitField(63)

	// Устанавливаем только чётные позиции
	SetBoolFieldUnchecked(packed, field0, true)
	SetBoolFieldUnchecked(packed, field31, true)

	// Проверяем все
	if !GetBoolField(packed, field0) {
		t.Error("bit 0 should be true")
	}
	if GetBoolField(packed, field7) {
		t.Error("bit 7 should be false")
	}
	if !GetBoolField(packed, field31) {
		t.Error("bit 31 should be true")
	}
	if GetBoolField(packed, field63) {
		t.Error("bit 63 should be false")
	}
}

// ============ Тесты для SetUIntFieldAs ============

func TestSetUIntFieldAsValidation(t *testing.T) {
	field := MustNewUIntBitField(0, 7, 100)
	packed := make([]byte, 8)

	// Валидное значение
	err := SetUIntFieldAs[uint64](packed, field, 50)
	if err != nil {
		t.Errorf("SetUIntFieldAs with valid value failed: %v", err)
	}

	// Превышение max
	err = SetUIntFieldAs[uint64](packed, field, 101)
	if err == nil {
		t.Error("SetUIntFieldAs should fail when value > max")
	}
}

func TestSetUIntFieldAsPreservesOtherFields(t *testing.T) {
	field1 := MustNewUIntBitField(0, 7, 255)
	field2 := MustNewUIntBitField(8, 15, 255)
	field3 := MustNewUIntBitField(16, 23, 255)

	packed := make([]byte, 8)

	// Устанавливаем все поля
	SetUIntFieldUncheckedAs[uint64](packed, field1, 0xAA)
	SetUIntFieldUncheckedAs[uint64](packed, field2, 0xBB)
	SetUIntFieldUncheckedAs[uint64](packed, field3, 0xCC)

	// Изменяем только среднее поле
	SetUIntFieldUncheckedAs[uint64](packed, field2, 0xDD)

	// Проверяем, что другие не изменились
	if got := GetUIntFieldAs[uint64](packed, field1); got != 0xAA {
		t.Errorf("field1 changed: got 0x%X, want 0xAA", got)
	}
	if got := GetUIntFieldAs[uint64](packed, field2); got != 0xDD {
		t.Errorf("field2 not updated: got 0x%X, want 0xDD", got)
	}
	if got := GetUIntFieldAs[uint64](packed, field3); got != 0xCC {
		t.Errorf("field3 changed: got 0x%X, want 0xCC", got)
	}
}

// ============ Тесты для SetIntFieldAs ============

func TestSetIntFieldAsValidation(t *testing.T) {
	field := MustNewIntBitField(0, 7, -50, 50)
	packed := make([]byte, 8)

	// Валидные значения
	err := SetIntFieldAs[int64](packed, field, 25)
	if err != nil {
		t.Errorf("SetIntFieldAs with valid positive value failed: %v", err)
	}

	err = SetIntFieldAs[int64](packed, field, -25)
	if err != nil {
		t.Errorf("SetIntFieldAs with valid negative value failed: %v", err)
	}

	// Ниже минимума
	err = SetIntFieldAs[int64](packed, field, -51)
	if err == nil {
		t.Error("SetIntFieldAs should fail when value < min")
	}

	// Выше максимума
	err = SetIntFieldAs[int64](packed, field, 51)
	if err == nil {
		t.Error("SetIntFieldAs should fail when value > max")
	}
}

func TestSetIntFieldAsRoundTrip(t *testing.T) {
	field := MustNewIntBitField(8, 23, -32768, 32767)
	packed := make([]byte, 8)

	testValues := []int64{-32768, -10000, -100, -1, 0, 1, 100, 10000, 32767}

	for _, value := range testValues {
		t.Run(string(rune('0'+int(value))), func(t *testing.T) {
			err := SetIntFieldAs[int64](packed, field, value)
			if err != nil {
				t.Fatalf("SetIntFieldAs failed: %v", err)
			}

			got := GetIntFieldAs[int64](packed, field)
			if got != value {
				t.Errorf("round-trip failed: stored %d, got %d", value, got)
			}
		})
	}
}

// ============ Тесты для SetBoolField ============

func TestSetBoolField(t *testing.T) {
	field := MustNewBoolBitField(16)
	packed := make([]byte, 8)

	// SetBoolField не должен возвращать ошибок
	err := SetBoolField(packed, field, true)
	if err != nil {
		t.Errorf("SetBoolField(true) returned error: %v", err)
	}

	if !GetBoolField(packed, field) {
		t.Error("bit not set to true")
	}

	err = SetBoolField(packed, field, false)
	if err != nil {
		t.Errorf("SetBoolField(false) returned error: %v", err)
	}

	if GetBoolField(packed, field) {
		t.Error("bit not cleared to false")
	}
}

// ============ Тесты для deprecated функций (обратная совместимость) ============

func TestDeprecatedGetUIntField(t *testing.T) {
	field := MustNewUIntBitField(0, 15, 65535)
	packed := make([]byte, 8)

	SetUIntFieldUncheckedAs[uint64](packed, field, 12345)

	// Тест deprecated функции
	got := GetUIntField(packed, field)
	if got != 12345 {
		t.Errorf("GetUIntField() = %d, want 12345", got)
	}
}

func TestDeprecatedSetUIntField(t *testing.T) {
	field := MustNewUIntBitField(0, 15, 65535)
	packed := make([]byte, 8)

	// Тест deprecated функции
	err := SetUIntField(packed, field, 54321)
	if err != nil {
		t.Errorf("SetUIntField failed: %v", err)
	}

	got := GetUIntField(packed, field)
	if got != 54321 {
		t.Errorf("value = %d, want 54321", got)
	}
}

func TestDeprecatedGetIntField(t *testing.T) {
	field := MustNewIntBitField(0, 15, -32768, 32767)
	packed := make([]byte, 8)

	SetIntFieldUncheckedAs[int64](packed, field, -12345)

	// Тест deprecated функции
	got := GetIntField(packed, field)
	if got != -12345 {
		t.Errorf("GetIntField() = %d, want -12345", got)
	}
}

func TestDeprecatedSetIntField(t *testing.T) {
	field := MustNewIntBitField(0, 15, -32768, 32767)
	packed := make([]byte, 8)

	// Тест deprecated функции
	err := SetIntField(packed, field, -9876)
	if err != nil {
		t.Errorf("SetIntField failed: %v", err)
	}

	got := GetIntField(packed, field)
	if got != -9876 {
		t.Errorf("value = %d, want -9876", got)
	}
}

// ============ Интеграционные тесты: сложные структуры ============

func TestComplexStructurePacking(t *testing.T) {
	// Моделируем сложную структуру данных (64 бита)
	// [0-7]    uint8: version (0-255)
	// [8-15]   int8: temperature (-128..127)
	// [16-23]  uint8: humidity (0-100)
	// [24-31]  int8: pressure delta (-50..50)
	// [32]     bool: alarm
	// [33-47]  uint15: sensor ID (0-32767)
	// [48-63]  uint16: sequence (0-65535)

	versionField := MustNewUIntBitField(0, 7, 255)
	tempField := MustNewIntBitField(8, 15, -128, 127)
	humidityField := MustNewUIntBitField(16, 23, 100)
	pressureField := MustNewIntBitField(24, 31, -50, 50)
	alarmField := MustNewBoolBitField(32)
	sensorField := MustNewUIntBitField(33, 47, 32767)
	sequenceField := MustNewUIntBitField(48, 63, 65535)

	packed := make([]byte, 8)

	// Записываем значения
	SetUIntFieldUncheckedAs[uint8](packed, versionField, 3)
	SetIntFieldUncheckedAs[int8](packed, tempField, -15)
	SetUIntFieldUncheckedAs[uint8](packed, humidityField, 65)
	SetIntFieldUncheckedAs[int8](packed, pressureField, 10)
	SetBoolFieldUnchecked(packed, alarmField, true)
	SetUIntFieldUncheckedAs[uint16](packed, sensorField, 12345)
	SetUIntFieldUncheckedAs[uint16](packed, sequenceField, 54321)

	// Проверяем все значения
	if got := GetUIntFieldAs[uint8](packed, versionField); got != 3 {
		t.Errorf("version = %d, want 3", got)
	}
	if got := GetIntFieldAs[int8](packed, tempField); got != -15 {
		t.Errorf("temperature = %d, want -15", got)
	}
	if got := GetUIntFieldAs[uint8](packed, humidityField); got != 65 {
		t.Errorf("humidity = %d, want 65", got)
	}
	if got := GetIntFieldAs[int8](packed, pressureField); got != 10 {
		t.Errorf("pressure = %d, want 10", got)
	}
	if got := GetBoolField(packed, alarmField); !got {
		t.Error("alarm = false, want true")
	}
	if got := GetUIntFieldAs[uint16](packed, sensorField); got != 12345 {
		t.Errorf("sensor ID = %d, want 12345", got)
	}
	if got := GetUIntFieldAs[uint16](packed, sequenceField); got != 54321 {
		t.Errorf("sequence = %d, want 54321", got)
	}
}

func TestPackedTypesUsage(t *testing.T) {
	// Тест использования именованных типов Packed
	var p8 Packed8
	var p16 Packed16
	var p32 Packed32
	var p64 Packed64

	// Проверяем размеры
	if len(p8) != 1 {
		t.Errorf("Packed8 length = %d, want 1", len(p8))
	}
	if len(p16) != 2 {
		t.Errorf("Packed16 length = %d, want 2", len(p16))
	}
	if len(p32) != 4 {
		t.Errorf("Packed32 length = %d, want 4", len(p32))
	}
	if len(p64) != 8 {
		t.Errorf("Packed64 length = %d, want 8", len(p64))
	}

	// Тест записи/чтения с Packed64
	field := MustNewUIntBitField(0, 31, math.MaxUint32)
	SetUIntFieldUncheckedAs[uint32](p64[:], field, 0x12345678)
	got := GetUIntFieldAs[uint32](p64[:], field)

	if got != 0x12345678 {
		t.Errorf("Packed64 usage: got 0x%X, want 0x12345678", got)
	}
}
