package bitpack

import (
	"errors"
	"fmt"
	"math"
	"testing"
)

// ============ Тесты для UIntBitField ============

func TestNewUIntBitField(t *testing.T) {
	tests := []struct {
		name    string
		start   BitPosition
		end     BitPosition
		max     uint64
		wantErr ErrorKind
	}{
		{
			name:    "valid field",
			start:   0,
			end:     7,
			max:     255,
			wantErr: -1, // no error
		},
		{
			name:    "start after end",
			start:   10,
			end:     5,
			max:     100,
			wantErr: KindStartAfterEnd,
		},
		{
			name:    "end out of range",
			start:   60,
			end:     64,
			max:     100,
			wantErr: KindEndOutOfRange,
		},
		{
			name:    "value overflow for width",
			start:   0,
			end:     3, // 4 bits → max 15
			max:     16,
			wantErr: KindValueOverflow,
		},
		{
			name:    "full 64-bit field",
			start:   0,
			end:     63,
			max:     math.MaxUint64,
			wantErr: -1,
		},
		{
			name:    "single bit field",
			start:   31,
			end:     31,
			max:     1,
			wantErr: -1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := NewUIntBitField(tt.start, tt.end, tt.max)
			if tt.wantErr == -1 {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
				return
			}

			// Проверяем, что ошибка имеет нужный тип
			var bfErr *Error
			if !errors.As(err, &bfErr) {
				t.Fatalf("expected *Error, got %T (error: %v)", err, err)
			}

			if bfErr.Kind != tt.wantErr {
				t.Errorf("expected error kind %v, got %v", tt.wantErr, bfErr.Kind)
			}
		})
	}
}

func TestMustNewUIntBitField(t *testing.T) {
	// Тест успешного создания
	t.Run("valid field", func(t *testing.T) {
		defer func() {
			if r := recover(); r != nil {
				t.Errorf("unexpected panic: %v", r)
			}
		}()
		bf := MustNewUIntBitField(0, 7, 255)
		if bf.Start != 0 || bf.End != 7 || bf.Max != 255 {
			t.Errorf("unexpected field values: %+v", bf)
		}
	})

	// Тест паники при ошибке
	t.Run("invalid field panics", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Error("expected panic but none occurred")
			}
		}()
		_ = MustNewUIntBitField(10, 5, 100) // start > end
	})
}

func TestUIntBitFieldGet(t *testing.T) {
	tests := []struct {
		name     string
		field    UIntBitField
		bitSet   BitSet64
		expected uint64
	}{
		{
			name:     "extract middle bits",
			field:    MustNewUIntBitField(4, 7, 15),
			bitSet:   BitSet64(0b11110000),
			expected: 0b1111, // bits 4-7
		},
		{
			name:     "extract LSBs",
			field:    MustNewUIntBitField(0, 3, 15),
			bitSet:   BitSet64(0b10101010),
			expected: 0b1010,
		},
		{
			name:     "extract MSBs",
			field:    MustNewUIntBitField(60, 63, 15),
			bitSet:   BitSet64(0xF000000000000000),
			expected: 0xF,
		},
		{
			name:     "single bit extraction",
			field:    MustNewUIntBitField(5, 5, 1),
			bitSet:   BitSet64(0b100000),
			expected: 1,
		},
		{
			name:     "zero value extraction",
			field:    MustNewUIntBitField(10, 15, 63),
			bitSet:   BitSet64(0),
			expected: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.field.Get(tt.bitSet)
			if got != tt.expected {
				t.Errorf("Get() = %d, want %d", got, tt.expected)
			}
		})
	}
}

func TestUIntBitFieldUpdate(t *testing.T) {
	tests := []struct {
		name        string
		field       UIntBitField
		initial     BitSet64
		value       uint64
		want        BitSet64
		wantErrKind ErrorKind
	}{
		{
			name:        "no error marker",
			field:       MustNewUIntBitField(4, 7, 15),
			initial:     BitSet64(0b00000000),
			value:       0b1111,
			want:        BitSet64(0b11110000),
			wantErrKind: -1,
		},
		{
			name:        "update without affecting other bits",
			field:       MustNewUIntBitField(4, 7, 15),
			initial:     BitSet64(0b11111111),
			value:       0b0000,
			want:        BitSet64(0b00001111),
			wantErrKind: -1,
		},
		{
			name:        "update MSBs",
			field:       MustNewUIntBitField(60, 63, 15),
			initial:     BitSet64(0),
			value:       0xF,
			want:        BitSet64(0xF000000000000000),
			wantErrKind: -1,
		},
		{
			name:        "value overflow",
			field:       MustNewUIntBitField(0, 3, 10), // max=10, but width allows 15
			initial:     BitSet64(0),
			value:       11,
			want:        0,
			wantErrKind: KindValueOverflow,
		},
		{
			name:        "full 64-bit update",
			field:       MustNewUIntBitField(0, 63, math.MaxUint64),
			initial:     BitSet64(0),
			value:       math.MaxUint64,
			want:        BitSet64(math.MaxUint64),
			wantErrKind: -1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.field.Update(tt.initial, tt.value)

			if tt.wantErrKind != -1 {
				if err == nil {
					t.Fatal("expected error but got nil")
				}
				var bfErr *Error
				if !errors.As(err, &bfErr) {
					t.Fatalf("expected *Error, got %T (error: %v)", err, err)
				}
				if bfErr.Kind != tt.wantErrKind {
					t.Errorf("error kind = %v, want %v", bfErr.Kind, tt.wantErrKind)
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got != tt.want {
				t.Errorf("Update() = 0x%X, want 0x%X", uint64(got), uint64(tt.want))
			}
		})
	}
}

func TestUIntBitFieldWidth(t *testing.T) {
	tests := []struct {
		start, end BitPosition
		want       uint8
	}{
		{0, 0, 1},
		{0, 7, 8},
		{32, 63, 32},
		{0, 63, 64},
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("[%d:%d]", tt.start, tt.end), func(t *testing.T) {
			bf := MustNewUIntBitField(tt.start, tt.end, maxAllowedForWidth(tt.end-tt.start+1))
			if got := bf.Width(); got != tt.want {
				t.Errorf("Width() = %d, want %d", got, tt.want)
			}
		})
	}
}

func TestUIntBitFieldString(t *testing.T) {
	bf := MustNewUIntBitField(10, 15, 63)
	want := "UIntBitField[10:15] max=63"
	if got := bf.String(); got != want {
		t.Errorf("String() = %q, want %q", got, want)
	}
}

// ============ Тесты для BoolBitField ============

func TestNewBoolBitField(t *testing.T) {
	tests := []struct {
		name    string
		pos     BitPosition
		wantErr bool
	}{
		{"valid position 0", 0, false},
		{"valid position 63", 63, false},
		{"invalid position 64", 64, true},
		{"invalid position 100", 100, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := NewBoolBitField(tt.pos)
			if tt.wantErr {
				if err == nil {
					t.Error("expected error but got nil")
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
			}
		})
	}
}

func TestMustNewBoolBitField(t *testing.T) {
	t.Run("valid position", func(t *testing.T) {
		defer func() {
			if r := recover(); r != nil {
				t.Errorf("unexpected panic: %v", r)
			}
		}()
		bf := MustNewBoolBitField(31)
		if bf.Position != 31 {
			t.Errorf("Position = %d, want 31", bf.Position)
		}
	})

	t.Run("invalid position panics", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Error("expected panic but none occurred")
			}
		}()
		_ = MustNewBoolBitField(64)
	})
}

func TestBoolBitFieldGet(t *testing.T) {
	tests := []struct {
		name     string
		pos      BitPosition
		bitSet   BitSet64
		expected bool
	}{
		{"bit 0 set", 0, BitSet64(1), true},
		{"bit 0 clear", 0, BitSet64(0), false},
		{"bit 5 set", 5, BitSet64(1 << 5), true},
		{"bit 5 clear", 5, BitSet64(0), false},
		{"bit 63 set", 63, BitSet64(1 << 63), true},
		{"multiple bits set", 3, BitSet64(0b1010), true},
		{"multiple bits clear", 2, BitSet64(0b1010), false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bf := MustNewBoolBitField(tt.pos)
			got := bf.Get(tt.bitSet)
			if got != tt.expected {
				t.Errorf("Get() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestBoolBitFieldUpdate(t *testing.T) {
	tests := []struct {
		name     string
		pos      BitPosition
		initial  BitSet64
		value    bool
		expected BitSet64
	}{
		{"set bit 0", 0, 0, true, BitSet64(1)},
		{"clear bit 0", 0, BitSet64(1), false, 0},
		{"set bit 31", 31, 0, true, BitSet64(1 << 31)},
		{"clear bit 31", 31, BitSet64(1 << 31), false, 0},
		{"set bit 63", 63, 0, true, BitSet64(1 << 63)},
		{"preserve other bits when setting", 4, BitSet64(0b111), true, BitSet64(0b10111)},
		{"preserve other bits when clearing", 1, BitSet64(0b111), false, BitSet64(0b101)},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bf := MustNewBoolBitField(tt.pos)
			got, err := bf.Update(tt.initial, tt.value)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got != tt.expected {
				t.Errorf("Update() = 0x%X, want 0x%X", uint64(got), uint64(tt.expected))
			}
		})
	}
}

func TestBoolBitFieldSetClearToggle(t *testing.T) {
	bf := MustNewBoolBitField(5)
	initial := BitSet64(0)

	// Set
	set := bf.Set(initial)
	if !bf.Get(set) {
		t.Error("Set() did not set the bit")
	}

	// Clear
	cleared := bf.Clear(set)
	if bf.Get(cleared) {
		t.Error("Clear() did not clear the bit")
	}

	// Toggle twice should return to original
	toggled := bf.Toggle(initial)
	if !bf.Get(toggled) {
		t.Error("first Toggle() did not set the bit")
	}
	toggledAgain := bf.Toggle(toggled)
	if bf.Get(toggledAgain) {
		t.Error("second Toggle() did not clear the bit")
	}
}

// ============ Тесты ошибок ============

func TestErrorFormatting(t *testing.T) {
	tests := []struct {
		name string
		err  error
		want string
	}{
		{
			name: "start after end",
			err:  newStartAfterEndError(10, 5),
			want: "bit field range error: start position (10) must be <= end position (5)",
		},
		{
			name: "end out of range",
			err:  newEndOutOfRangeError(64),
			want: "bit field range error: end position (64) must be < 64",
		},
		{
			name: "value overflow",
			err:  newValueOverflowError(16, 15, 4),
			want: "value 16 exceeds capacity of 4-bit field (max allowed: 15)",
		},
		{
			name: "position out of range",
			err:  newPositionOutOfRangeError(64),
			want: "bool bit field error: position (64) must be < 64",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.err.Error()
			if got != tt.want {
				t.Errorf("Error() = %q\nwant %q", got, tt.want)
			}
		})
	}
}

// ============ Интеграционные тесты ============

func TestMultipleFieldsInSameBitSet(t *testing.T) {
	// Проверка независимой работы нескольких полей в одном битовом наборе
	year := MustNewUIntBitField(0, 6, 99)   // 7 bits for year (0-99)
	month := MustNewUIntBitField(7, 10, 12) // 4 bits for month (1-12)
	day := MustNewUIntBitField(11, 15, 31)  // 5 bits for day (1-31)
	isLeap := MustNewBoolBitField(16)       // 1 bit for leap year flag

	var bitSet BitSet64

	// Установка значений
	var err error
	bitSet, err = year.Update(bitSet, 24) // 2024
	if err != nil {
		t.Fatalf("year update failed: %v", err)
	}
	bitSet, err = month.Update(bitSet, 2) // February
	if err != nil {
		t.Fatalf("month update failed: %v", err)
	}
	bitSet, err = day.Update(bitSet, 29) // 29th
	if err != nil {
		t.Fatalf("day update failed: %v", err)
	}
	bitSet, err = isLeap.Update(bitSet, true)
	if err != nil {
		t.Fatalf("isLeap update failed: %v", err)
	}

	// Проверка извлечения
	if got := year.Get(bitSet); got != 24 {
		t.Errorf("year = %d, want 24", got)
	}
	if got := month.Get(bitSet); got != 2 {
		t.Errorf("month = %d, want 2", got)
	}
	if got := day.Get(bitSet); got != 29 {
		t.Errorf("day = %d, want 29", got)
	}
	if got := isLeap.Get(bitSet); !got {
		t.Errorf("isLeap = %v, want true", got)
	}

	// Изменение одного поля не должно влиять на другие
	bitSet, err = month.Update(bitSet, 3) // Change to March
	if err != nil {
		t.Fatalf("month update failed: %v", err)
	}
	if got := year.Get(bitSet); got != 24 {
		t.Errorf("year changed after month update: %d", got)
	}
	if got := day.Get(bitSet); got != 29 {
		t.Errorf("day changed after month update: %d", got)
	}
}

func TestEdgeCaseFullBitWidth(t *testing.T) {
	// Тестирование 64-битного поля (особый случай для маски)
	fullField := MustNewUIntBitField(0, 63, math.MaxUint64)
	initial := BitSet64(0x123456789ABCDEF0)

	// Извлечение
	got := fullField.Get(initial)
	if got != uint64(initial) {
		t.Errorf("Get() = 0x%X, want 0x%X", got, uint64(initial))
	}

	// Обновление
	newVal := uint64(0xFEDCBA9876543210)
	updated, err := fullField.Update(initial, newVal)
	if err != nil {
		t.Fatalf("Update failed: %v", err)
	}
	if uint64(updated) != newVal {
		t.Errorf("Update() = 0x%X, want 0x%X", uint64(updated), newVal)
	}
}

// ============ Дополнительные тесты для утилит ============

// ============ Тесты для IntBitField ============

func TestNewIntBitField(t *testing.T) {
	tests := []struct {
		name    string
		start   BitPosition
		end     BitPosition
		min     int64
		max     int64
		wantErr ErrorKind
	}{
		{
			name:    "valid 8-bit field full range",
			start:   0,
			end:     7,
			min:     -128,
			max:     127,
			wantErr: -1,
		},
		{
			name:    "valid 8-bit field restricted range",
			start:   0,
			end:     7,
			min:     -50,
			max:     50,
			wantErr: -1,
		},
		{
			name:    "start after end",
			start:   10,
			end:     5,
			min:     -10,
			max:     10,
			wantErr: KindStartAfterEnd,
		},
		{
			name:    "end out of range",
			start:   60,
			end:     64,
			min:     -10,
			max:     10,
			wantErr: KindEndOutOfRange,
		},
		{
			name:    "min below allowed for width",
			start:   0,
			end:     3, // 4 bits → min=-8
			min:     -9,
			max:     7,
			wantErr: KindValueOverflow,
		},
		{
			name:    "max above allowed for width",
			start:   0,
			end:     3, // 4 bits → max=7
			min:     -8,
			max:     8,
			wantErr: KindValueOverflow,
		},
		{
			name:    "min > max inverted range",
			start:   0,
			end:     7,
			min:     100,
			max:     50,
			wantErr: KindValueRangeInverted,
		},
		{
			name:    "full 64-bit field",
			start:   0,
			end:     63,
			min:     math.MinInt64,
			max:     math.MaxInt64,
			wantErr: -1,
		},
		{
			name:    "single bit field (-1 to 0)",
			start:   31,
			end:     31,
			min:     -1,
			max:     0,
			wantErr: -1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := NewIntBitField(tt.start, tt.end, tt.min, tt.max)
			if tt.wantErr == -1 {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
				return
			}

			var bfErr *Error
			if !errors.As(err, &bfErr) {
				t.Fatalf("expected *Error, got %T (error: %v)", err, err)
			}

			if bfErr.Kind != tt.wantErr {
				t.Errorf("expected error kind %v, got %v", tt.wantErr, bfErr.Kind)
			}
		})
	}
}

func TestMustNewIntBitField(t *testing.T) {
	t.Run("valid field", func(t *testing.T) {
		defer func() {
			if r := recover(); r != nil {
				t.Errorf("unexpected panic: %v", r)
			}
		}()
		bf := MustNewIntBitField(0, 7, -128, 127)
		if bf.Start != 0 || bf.End != 7 || bf.Min != -128 || bf.Max != 127 {
			t.Errorf("unexpected field values: %+v", bf)
		}
	})

	t.Run("invalid field panics", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Error("expected panic but none occurred")
			}
		}()
		_ = MustNewIntBitField(10, 5, -10, 10) // start > end
	})
}

func TestNewIntBitFieldAuto(t *testing.T) {
	tests := []struct {
		name    string
		start   BitPosition
		end     BitPosition
		wantMin int64
		wantMax int64
	}{
		{"4 bits", 0, 3, -8, 7},
		{"8 bits", 0, 7, -128, 127},
		{"16 bits", 0, 15, -32768, 32767},
		{"32 bits", 0, 31, -2147483648, 2147483647},
		{"64 bits", 0, 63, math.MinInt64, math.MaxInt64},
		{"single bit", 31, 31, -1, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bf, err := NewIntBitFieldAuto(tt.start, tt.end)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if bf.Min != tt.wantMin {
				t.Errorf("Min = %d, want %d", bf.Min, tt.wantMin)
			}
			if bf.Max != tt.wantMax {
				t.Errorf("Max = %d, want %d", bf.Max, tt.wantMax)
			}
		})
	}
}

func TestIntBitFieldGet(t *testing.T) {
	tests := []struct {
		name     string
		field    IntBitField
		bitSet   BitSet64
		expected int64
	}{
		{
			name:     "positive 4-bit value",
			field:    MustNewIntBitField(0, 3, -8, 7),
			bitSet:   BitSet64(0b0101), // 5
			expected: 5,
		},
		{
			name:     "negative 4-bit value (-3)",
			field:    MustNewIntBitField(0, 3, -8, 7),
			bitSet:   BitSet64(0b1101), // -3 in two's complement
			expected: -3,
		},
		{
			name:     "negative 8-bit value (-128)",
			field:    MustNewIntBitField(0, 7, -128, 127),
			bitSet:   BitSet64(0b10000000),
			expected: -128,
		},
		{
			name:     "negative 8-bit value (-1)",
			field:    MustNewIntBitField(0, 7, -128, 127),
			bitSet:   BitSet64(0b11111111),
			expected: -1,
		},
		{
			name:     "extract from middle bits (negative)",
			field:    MustNewIntBitField(4, 7, -8, 7),
			bitSet:   BitSet64(0b11010000), // bits 4-7 = 0b1101 (-3)
			expected: -3,
		},
		{
			name:     "extract MSBs (16-bit negative)",
			field:    MustNewIntBitField(48, 63, -32768, 32767),
			bitSet:   BitSet64(0xFFFF000000000000), // -1 in upper 16 bits
			expected: -1,
		},
		{
			name:     "zero value",
			field:    MustNewIntBitField(10, 15, -32, 31),
			bitSet:   BitSet64(0),
			expected: 0,
		},
		{
			name:     "max positive 8-bit",
			field:    MustNewIntBitField(0, 7, -128, 127),
			bitSet:   BitSet64(0b01111111),
			expected: 127,
		},
		{
			name:     "min negative 8-bit",
			field:    MustNewIntBitField(0, 7, -128, 127),
			bitSet:   BitSet64(0b10000000),
			expected: -128,
		},
		{
			name:     "64-bit full range positive",
			field:    MustNewIntBitField(0, 63, math.MinInt64, math.MaxInt64),
			bitSet:   BitSet64(0x7FFFFFFFFFFFFFFF),
			expected: math.MaxInt64,
		},
		{
			name:     "64-bit full range negative",
			field:    MustNewIntBitField(0, 63, math.MinInt64, math.MaxInt64),
			bitSet:   BitSet64(0x8000000000000000),
			expected: math.MinInt64,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.field.Get(tt.bitSet)
			if got != tt.expected {
				t.Errorf("Get() = %d, want %d (bitset=0x%X)", got, tt.expected, uint64(tt.bitSet))
			}
		})
	}
}

func TestIntBitFieldUpdate(t *testing.T) {
	tests := []struct {
		name        string
		field       IntBitField
		initial     BitSet64
		value       int64
		want        BitSet64
		wantErrKind ErrorKind
	}{
		{
			name:        "positive value within range",
			field:       MustNewIntBitField(0, 3, -8, 7),
			initial:     BitSet64(0),
			value:       5,
			want:        BitSet64(0b0101),
			wantErrKind: -1,
		},
		{
			name:        "negative value within range (-3)",
			field:       MustNewIntBitField(0, 3, -8, 7),
			initial:     BitSet64(0),
			value:       -3,
			want:        BitSet64(0b1101), // two's complement representation
			wantErrKind: -1,
		},
		{
			name:        "update without affecting other bits",
			field:       MustNewIntBitField(4, 7, -8, 7),
			initial:     BitSet64(0b00001111),
			value:       -1, // 0b1111
			want:        BitSet64(0b11111111),
			wantErrKind: -1,
		},
		{
			name:        "value below min",
			field:       MustNewIntBitField(0, 7, -50, 50),
			initial:     BitSet64(0),
			value:       -51,
			want:        0,
			wantErrKind: KindValueOverflow,
		},
		{
			name:        "value above max",
			field:       MustNewIntBitField(0, 7, -50, 50),
			initial:     BitSet64(0),
			value:       51,
			want:        0,
			wantErrKind: KindValueOverflow,
		},
		{
			name:        "max positive 8-bit",
			field:       MustNewIntBitField(0, 7, -128, 127),
			initial:     BitSet64(0),
			value:       127,
			want:        BitSet64(0b01111111),
			wantErrKind: -1,
		},
		{
			name:        "min negative 8-bit",
			field:       MustNewIntBitField(0, 7, -128, 127),
			initial:     BitSet64(0),
			value:       -128,
			want:        BitSet64(0b10000000),
			wantErrKind: -1,
		},
		{
			name:        "64-bit update positive",
			field:       MustNewIntBitField(0, 63, math.MinInt64, math.MaxInt64),
			initial:     BitSet64(0),
			value:       0x123456789ABCDEF0,
			want:        BitSet64(0x123456789ABCDEF0),
			wantErrKind: -1,
		},
		{
			name:        "64-bit update negative",
			field:       MustNewIntBitField(0, 63, math.MinInt64, math.MaxInt64),
			initial:     BitSet64(0),
			value:       -9223372036854775808, // math.MinInt64
			want:        BitSet64(0x8000000000000000),
			wantErrKind: -1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.field.Update(tt.initial, tt.value)

			if tt.wantErrKind != -1 {
				if err == nil {
					t.Fatal("expected error but got nil")
				}
				var bfErr *Error
				if !errors.As(err, &bfErr) {
					t.Fatalf("expected *Error, got %T (error: %v)", err, err)
				}
				if bfErr.Kind != tt.wantErrKind {
					t.Errorf("error kind = %v, want %v", bfErr.Kind, tt.wantErrKind)
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got != tt.want {
				t.Errorf("Update() = 0x%X, want 0x%X (value=%d)", uint64(got), uint64(tt.want), tt.value)
			}

			// Дополнительная проверка: извлечение должно вернуть исходное значение
			restored := tt.field.Get(got)
			if restored != tt.value {
				t.Errorf("round-trip failed: stored %d but got back %d", tt.value, restored)
			}
		})
	}
}

func TestIntBitFieldWidth(t *testing.T) {
	tests := []struct {
		start, end BitPosition
		want       uint8
	}{
		{0, 0, 1},
		{0, 7, 8},
		{0, 15, 16},
		{32, 63, 32},
		{0, 63, 64},
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("[%d:%d]", tt.start, tt.end), func(t *testing.T) {
			minVal, maxVal := intRangeForWidth(tt.end - tt.start + 1)
			bf := MustNewIntBitField(tt.start, tt.end, minVal, maxVal)
			if got := bf.Width(); got != tt.want {
				t.Errorf("Width() = %d, want %d", got, tt.want)
			}
		})
	}
}

func TestIntBitFieldString(t *testing.T) {
	bf := MustNewIntBitField(10, 15, -32, 31)
	want := "IntBitField[10:15] range=[-32,31]"
	if got := bf.String(); got != want {
		t.Errorf("String() = %q, want %q", got, want)
	}
}

func TestIntRangeForWidth(t *testing.T) {
	tests := []struct {
		width   uint8
		wantMin int64
		wantMax int64
	}{
		{1, -1, 0},
		{4, -8, 7},
		{8, -128, 127},
		{16, -32768, 32767},
		{32, -2147483648, 2147483647},
		{63, math.MinInt64 / 2, math.MaxInt64 / 2},
		{64, math.MinInt64, math.MaxInt64},
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("width=%d", tt.width), func(t *testing.T) {
			gotMin, gotMax := intRangeForWidth(tt.width)
			if gotMin != tt.wantMin || gotMax != tt.wantMax {
				t.Errorf("intRangeForWidth(%d) = (%d, %d), want (%d, %d)",
					tt.width, gotMin, gotMax, tt.wantMin, tt.wantMax)
			}
		})
	}
}

// ============ Интеграционные тесты для знаковых полей ============

func TestIntBitFieldRoundTrip(t *testing.T) {
	// Тест корректности полного цикла: запись → извлечение для разных значений
	field := MustNewIntBitField(8, 15, -128, 127) // 8-bit signed at offset 8

	testValues := []int64{
		-128, -100, -50, -10, -1, 0, 1, 10, 50, 100, 127,
	}

	for _, value := range testValues {
		t.Run(fmt.Sprintf("value=%d", value), func(t *testing.T) {
			// Запись
			updated, err := field.Update(0, value)
			if err != nil {
				t.Fatalf("Update failed for %d: %v", value, err)
			}

			// Извлечение
			restored := field.Get(updated)

			// Проверка
			if restored != value {
				t.Errorf("round-trip failed: stored %d but got back %d (bitset=0x%X)",
					value, restored, uint64(updated))
			}

			// Проверка, что другие биты не затронуты
			otherBits := uint64(updated) & ^(field.mask << field.Start)
			if otherBits != 0 {
				t.Errorf("unexpected bits set outside field: 0x%X", otherBits)
			}
		})
	}
}

func TestMixedSignedUnsignedFields(t *testing.T) {
	// Сценарий: структура с температурой (знаковая) и влажностью (беззнаковая)
	temperature := MustNewIntBitField(0, 7, -50, 100)      // 8 bits signed
	humidity := MustNewUIntBitField(8, 15, 100)            // 8 bits unsigned
	pressureDelta := MustNewIntBitField(16, 23, -128, 127) // 8 bits signed

	var bitSet BitSet64

	// Установка значений
	var err error
	bitSet, err = temperature.Update(bitSet, -5)
	if err != nil {
		t.Fatalf("temperature update failed: %v", err)
	}
	bitSet, err = humidity.Update(bitSet, 65)
	if err != nil {
		t.Fatalf("humidity update failed: %v", err)
	}
	bitSet, err = pressureDelta.Update(bitSet, -10)
	if err != nil {
		t.Fatalf("pressureDelta update failed: %v", err)
	}

	// Проверка извлечения
	if got := temperature.Get(bitSet); got != -5 {
		t.Errorf("temperature = %d, want -5", got)
	}
	if got := humidity.Get(bitSet); got != 65 {
		t.Errorf("humidity = %d, want 65", got)
	}
	if got := pressureDelta.Get(bitSet); got != -10 {
		t.Errorf("pressureDelta = %d, want -10", got)
	}

	// Изменение одного поля не должно влиять на другие
	bitSet, err = temperature.Update(bitSet, 25)
	if err != nil {
		t.Fatalf("temperature update failed: %v", err)
	}
	if got := humidity.Get(bitSet); got != 65 {
		t.Errorf("humidity changed after temperature update: %d", got)
	}
	if got := pressureDelta.Get(bitSet); got != -10 {
		t.Errorf("pressureDelta changed after temperature update: %d", got)
	}
}

// ============ Тест производительности знакового расширения ============

func TestSignExtensionEdgeCases(t *testing.T) {
	// Проверка корректности знакового расширения для краевых случаев
	tests := []struct {
		name     string
		width    uint8
		rawBits  uint64 // Битовое представление в поле
		expected int64
	}{
		// 4-bit field edge cases
		{"4-bit min", 4, 0b1000, -8},
		{"4-bit max", 4, 0b0111, 7},
		{"4-bit -1", 4, 0b1111, -1},
		// 8-bit field edge cases
		{"8-bit min", 8, 0b10000000, -128},
		{"8-bit max", 8, 0b01111111, 127},
		{"8-bit -1", 8, 0b11111111, -1},
		// 16-bit field edge cases
		{"16-bit min", 16, 0b10000000_00000000, -32768},
		{"16-bit max", 16, 0b01111111_11111111, 32767},
		// 32-bit field edge cases
		{"32-bit min", 32, 0x80000000, -2147483648},
		{"32-bit max", 32, 0x7FFFFFFF, 2147483647},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Создаём поле шириной tt.width бит в младших позициях
			mainVal, maxVal := intRangeForWidth(tt.width)
			field := MustNewIntBitField(0, tt.width-1, mainVal, maxVal)

			// Формируем bitset с заданными битами
			bitSet := BitSet64(tt.rawBits & computeMask(tt.width))

			// Извлекаем значение
			got := field.Get(bitSet)

			if got != tt.expected {
				t.Errorf("sign extension failed: got %d, want %d (rawBits=0x%X, width=%d)",
					got, tt.expected, tt.rawBits, tt.width)
			}
		})
	}
}

// ============ Дополнительные тесты для UIntBitField ============

func TestUIntBitFieldUpdateUnchecked(t *testing.T) {
	tests := []struct {
		name    string
		field   UIntBitField
		initial BitSet64
		value   uint64
		want    BitSet64
	}{
		{
			name:    "unchecked update in middle",
			field:   MustNewUIntBitField(8, 15, 255),
			initial: BitSet64(0xFFFF),
			value:   0xAB,
			want:    BitSet64(0xABFF),
		},
		{
			name:    "unchecked update with value exceeding max (should work without validation)",
			field:   MustNewUIntBitField(0, 7, 100),
			initial: BitSet64(0),
			value:   255, // Превышает max=100, но unchecked не валидирует
			want:    BitSet64(255),
		},
		{
			name:    "unchecked full 64-bit",
			field:   MustNewUIntBitField(0, 63, math.MaxUint64),
			initial: BitSet64(0),
			value:   math.MaxUint64,
			want:    BitSet64(math.MaxUint64),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.field.UpdateUnchecked(tt.initial, tt.value)
			if got != tt.want {
				t.Errorf("UpdateUnchecked() = 0x%X, want 0x%X", uint64(got), uint64(tt.want))
			}
		})
	}
}

func TestUIntBitFieldMaskCaching(t *testing.T) {
	// Проверяем, что маска корректно кэширована при создании
	field := MustNewUIntBitField(4, 11, 255) // 8 бит
	expectedMask := uint64(0xFF)

	if field.mask != expectedMask {
		t.Errorf("cached mask = 0x%X, want 0x%X", field.mask, expectedMask)
	}
}

func TestUIntBitFieldBoundaryValues(t *testing.T) {
	// Проверка граничных значений для разных ширин
	tests := []struct {
		width uint8
		max   uint64
	}{
		{1, 1},
		{2, 3},
		{4, 15},
		{8, 255},
		{16, 65535},
		{32, 0xFFFFFFFF},
		{63, math.MaxInt64},
		{64, math.MaxUint64},
	}

	for _, tt := range tests {
		field, err := NewUIntBitField(0, tt.width-1, tt.max)
		if err != nil {
			t.Errorf("width %d: unexpected error: %v", tt.width, err)
			continue
		}

		// Тестируем максимальное значение
		updated, err := field.Update(0, tt.max)
		if err != nil {
			t.Errorf("width %d: failed to set max value: %v", tt.width, err)
			continue
		}

		retrieved := field.Get(updated)
		if retrieved != tt.max {
			t.Errorf("width %d: round-trip failed for max: got %d, want %d", tt.width, retrieved, tt.max)
		}
	}
}

// ============ Дополнительные тесты для IntBitField ============

func TestIntBitFieldUpdateUnchecked(t *testing.T) {
	tests := []struct {
		name    string
		field   IntBitField
		initial BitSet64
		value   int64
		want    BitSet64
	}{
		{
			name:    "unchecked negative update",
			field:   MustNewIntBitField(0, 7, -128, 127),
			initial: BitSet64(0),
			value:   -42,
			want:    BitSet64(0xD6), //  complement of -42 in 8 bits
		},
		{
			name:    "unchecked update preserves other bits",
			field:   MustNewIntBitField(8, 15, -128, 127),
			initial: BitSet64(0xFF00FF),
			value:   -1,
			want:    BitSet64(0xFFFFFF),
		},
		{
			name:    "unchecked with value exceeding range (should work without validation)",
			field:   MustNewIntBitField(0, 7, -50, 50),
			initial: BitSet64(0),
			value:   127, // Превышает max=50, но unchecked не валидирует
			want:    BitSet64(127),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.field.UpdateUnchecked(tt.initial, tt.value)
			if got != tt.want {
				t.Errorf("UpdateUnchecked() = 0x%X, want 0x%X", uint64(got), uint64(tt.want))
			}
		})
	}
}

func TestIntBitFieldInvalidRanges(t *testing.T) {
	tests := []struct {
		name    string
		start   BitPosition
		end     BitPosition
		min     int64
		max     int64
		wantErr ErrorKind
	}{
		{
			name:    "min > max",
			start:   0,
			end:     7,
			min:     100,
			max:     -100,
			wantErr: KindValueRangeInverted,
		},
		{
			name:    "min below allowed range",
			start:   0,
			end:     3, // 4 бита: -8..7
			min:     -10,
			max:     7,
			wantErr: KindValueOverflow,
		},
		{
			name:    "max above allowed range",
			start:   0,
			end:     3, // 4 бита: -8..7
			min:     -8,
			max:     10,
			wantErr: KindValueOverflow,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := NewIntBitField(tt.start, tt.end, tt.min, tt.max)
			if err == nil {
				t.Fatal("expected error but got nil")
			}

			var bfErr *Error
			if !errors.As(err, &bfErr) {
				t.Fatalf("expected *Error, got %T", err)
			}

			if bfErr.Kind != tt.wantErr {
				t.Errorf("error kind = %v, want %v", bfErr.Kind, tt.wantErr)
			}
		})
	}
}

func TestIntBitFieldNegativeZeroDistinction(t *testing.T) {
	// Проверяем, что 0 и -0 обрабатываются одинаково (нет -0 в two's complement)
	field := MustNewIntBitField(0, 7, -128, 127)

	zero, err := field.Update(0, 0)
	if err != nil {
		t.Fatalf("failed to set 0: %v", err)
	}

	negZero, err := field.Update(0, -0)
	if err != nil {
		t.Fatalf("failed to set -0: %v", err)
	}

	if zero != negZero {
		t.Errorf("0 and -0 produced different bitsets: 0x%X vs 0x%X", uint64(zero), uint64(negZero))
	}

	if field.Get(zero) != 0 {
		t.Errorf("Get(0) = %d, want 0", field.Get(zero))
	}
}

// ============ Дополнительные тесты для BoolBitField ============

func TestBoolBitFieldAllPositions(t *testing.T) {
	// Проверяем, что все позиции от 0 до 63 работают корректно
	for pos := BitPosition(0); pos < 64; pos++ {
		t.Run(string(rune('0'+pos)), func(t *testing.T) {
			field := MustNewBoolBitField(pos)

			// Установка true
			setTrue := field.Set(0)
			if !field.Get(setTrue) {
				t.Errorf("pos %d: Set() didn't set bit to true", pos)
			}

			// Очистка
			cleared := field.Clear(setTrue)
			if field.Get(cleared) {
				t.Errorf("pos %d: Clear() didn't clear bit", pos)
			}

			// Toggle с false на true
			toggled := field.Toggle(cleared)
			if !field.Get(toggled) {
				t.Errorf("pos %d: Toggle() from false didn't set to true", pos)
			}

			// Toggle с true на false
			toggled = field.Toggle(toggled)
			if field.Get(toggled) {
				t.Errorf("pos %d: Toggle() from true didn't set to false", pos)
			}
		})
	}
}

func TestBoolBitFieldUpdateReturn(t *testing.T) {
	// Проверяем, что Update не возвращает ошибок (для соответствия интерфейсу)
	field := MustNewBoolBitField(5)

	_, err := field.Update(0, true)
	if err != nil {
		t.Errorf("Update(true) returned unexpected error: %v", err)
	}

	_, err = field.Update(0, false)
	if err != nil {
		t.Errorf("Update(false) returned unexpected error: %v", err)
	}
}

func TestBoolBitFieldIsolation(t *testing.T) {
	// Проверяем, что изменение одного бита не влияет на другие
	field1 := MustNewBoolBitField(0)
	field2 := MustNewBoolBitField(31)
	field3 := MustNewBoolBitField(63)

	var bitSet BitSet64

	// Устанавливаем биты
	bitSet = field1.Set(bitSet)
	bitSet = field2.Set(bitSet)
	bitSet = field3.Set(bitSet)

	// Проверяем все установлены
	if !field1.Get(bitSet) || !field2.Get(bitSet) || !field3.Get(bitSet) {
		t.Error("not all bits were set")
	}

	// Очищаем средний бит
	bitSet = field2.Clear(bitSet)

	// Проверяем, что другие не затронуты
	if !field1.Get(bitSet) {
		t.Error("bit 0 was affected by clearing bit 31")
	}
	if field2.Get(bitSet) {
		t.Error("bit 31 was not cleared")
	}
	if !field3.Get(bitSet) {
		t.Error("bit 63 was affected by clearing bit 31")
	}
}

func TestBoolBitFieldString(t *testing.T) {
	field := MustNewBoolBitField(42)
	str := field.String()

	// Проверяем, что строка содержит позицию
	expected := "\"BoolBitField[42] "
	if str != expected {
		t.Errorf("String() = %q, want %q", str, expected)
	}
}

// ============ Тесты вспомогательных функций ============

func TestComputeMask(t *testing.T) {
	tests := []struct {
		width uint8
		want  uint64
	}{
		{0, 0}, // Граничный случай
		{1, 0x1},
		{4, 0xF},
		{8, 0xFF},
		{16, 0xFFFF},
		{32, 0xFFFFFFFF},
		{63, 0x7FFFFFFFFFFFFFFF},
		{64, 0xFFFFFFFFFFFFFFFF},
		{100, 0xFFFFFFFFFFFFFFFF}, // > 64 должно дать MaxUint64
	}

	for _, tt := range tests {
		got := computeMask(tt.width)
		if got != tt.want {
			t.Errorf("computeMask(%d) = 0x%X, want 0x%X", tt.width, got, tt.want)
		}
	}
}

func TestMaxAllowedForWidth(t *testing.T) {
	tests := []struct {
		width uint8
		want  uint64
	}{
		{1, 1},
		{2, 3},
		{4, 15},
		{8, 255},
		{16, 65535},
		{32, 4294967295},
		{64, math.MaxUint64},
	}

	for _, tt := range tests {
		got := maxAllowedForWidth(tt.width)
		if got != tt.want {
			t.Errorf("maxAllowedForWidth(%d) = %d, want %d", tt.width, got, tt.want)
		}
	}
}

// ============ Стресс-тесты и граничные случаи ============

func TestMultipleFieldsNoInterference(t *testing.T) {
	// Создаём множество полей и проверяем отсутствие взаимовлияния
	field1 := MustNewUIntBitField(0, 7, 255)
	field2 := MustNewIntBitField(8, 15, -128, 127)
	field3 := MustNewBoolBitField(16)
	field4 := MustNewUIntBitField(17, 31, 32767)
	field5 := MustNewBoolBitField(32)
	field6 := MustNewIntBitField(33, 47, -16384, 16383)
	field7 := MustNewUIntBitField(48, 63, 65535)

	var bitSet BitSet64
	var err error

	// Устанавливаем все поля
	bitSet, err = field1.Update(bitSet, 0xAB)
	if err != nil {
		t.Fatalf("field1 update failed: %v", err)
	}

	bitSet, err = field2.Update(bitSet, -42)
	if err != nil {
		t.Fatalf("field2 update failed: %v", err)
	}

	bitSet = field3.Set(bitSet)

	bitSet, err = field4.Update(bitSet, 12345)
	if err != nil {
		t.Fatalf("field4 update failed: %v", err)
	}

	bitSet = field5.Clear(bitSet)

	bitSet, err = field6.Update(bitSet, -9999)
	if err != nil {
		t.Fatalf("field6 update failed: %v", err)
	}

	bitSet, err = field7.Update(bitSet, 54321)
	if err != nil {
		t.Fatalf("field7 update failed: %v", err)
	}

	// Проверяем все значения
	if got := field1.Get(bitSet); got != 0xAB {
		t.Errorf("field1 = %d, want 171", got)
	}
	if got := field2.Get(bitSet); got != -42 {
		t.Errorf("field2 = %d, want -42", got)
	}
	if got := field3.Get(bitSet); !got {
		t.Error("field3 = false, want true")
	}
	if got := field4.Get(bitSet); got != 12345 {
		t.Errorf("field4 = %d, want 12345", got)
	}
	if got := field5.Get(bitSet); got {
		t.Error("field5 = true, want false")
	}
	if got := field6.Get(bitSet); got != -9999 {
		t.Errorf("field6 = %d, want -9999", got)
	}
	if got := field7.Get(bitSet); got != 54321 {
		t.Errorf("field7 = %d, want 54321", got)
	}
}

func TestZeroWidthFieldPrevention(t *testing.T) {
	// Проверяем, что поля нулевой ширины (start == end) работают корректно
	// Это single-bit поля

	// UInt single bit
	field1, err := NewUIntBitField(5, 5, 1)
	if err != nil {
		t.Errorf("single-bit UInt field creation failed: %v", err)
	}
	if field1.Width() != 1 {
		t.Errorf("single-bit UInt width = %d, want 1", field1.Width())
	}

	// Int single bit (может хранить -1 или 0)
	field2, err := NewIntBitField(10, 10, -1, 0)
	if err != nil {
		t.Errorf("single-bit Int field creation failed: %v", err)
	}
	if field2.Width() != 1 {
		t.Errorf("single-bit Int width = %d, want 1", field2.Width())
	}
}

// ============ Тесты для проверки корректности two's complement ============

func TestTwosComplementEncoding(t *testing.T) {
	// Проверяем, что отрицательные числа корректно кодируются в дополнительном коде
	tests := []struct {
		value       int64
		width       uint8
		expectedRaw uint64 // Ожидаемое беззнаковое представление
	}{
		{-1, 4, 0b1111},
		{-2, 4, 0b1110},
		{-8, 4, 0b1000},
		{-1, 8, 0xFF},
		{-128, 8, 0x80},
		{-1, 16, 0xFFFF},
		{-32768, 16, 0x8000},
	}

	for _, tt := range tests {
		t.Run(string(rune('0'+tt.value)), func(t *testing.T) {
			minVal, maxVal := intRangeForWidth(tt.width)
			field := MustNewIntBitField(0, tt.width-1, minVal, maxVal)

			updated := field.UpdateUnchecked(0, tt.value)
			raw := uint64(updated) & computeMask(tt.width)

			if raw != tt.expectedRaw {
				t.Errorf("two's complement encoding of %d in %d bits: got 0x%X, want 0x%X",
					tt.value, tt.width, raw, tt.expectedRaw)
			}

			// Проверяем обратное преобразование
			retrieved := field.Get(updated)
			if retrieved != tt.value {
				t.Errorf("round-trip failed: encoded %d, retrieved %d", tt.value, retrieved)
			}
		})
	}
}
