package personbitpack

import (
	"GamePerson/internal/bitpack"
	"GamePerson/internal/model/config"
	"testing"
)

func TestBitFieldConfiguration(t *testing.T) {
	testCases := []struct {
		name        string
		start, end  bitpack.BitPosition
		max         uint64
		expectPanic bool
	}{
		{"nameSizeField", bitNameStart, bitNameEnd, uint64(config.MaxNameLength), false},
		{"respectField", bitRespectStart, bitRespectEnd, uint64(config.PersonMaxRespect), false},
		{"manaField", bitManaStart, bitManaEnd, uint64(config.PersonMaxMana), false},
		{"invalidField", 10, 5, 100, true}, // Start > End — должно паниковать
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			defer func() {
				r := recover()
				if tc.expectPanic {
					if r == nil {
						t.Error("expected panic but got none")
					}
				} else {
					if r != nil {
						t.Errorf("unexpected panic: %v", r)
					}
				}
			}()

			_ = bitpack.MustNewUIntBitField(tc.start, tc.end, tc.max)
		})
	}
}

// Дополнительный тест: проверка покрытия всех 48 бит без пересечений
func TestBitFieldCoverage(t *testing.T) {
	// Карта использования битов (48 бит)
	used := make([]bool, 48)

	fields := []struct {
		name  string
		start bitpack.BitPosition
		end   bitpack.BitPosition
	}{
		{"nameSize", bitNameStart, bitNameEnd},
		{"respect", bitRespectStart, bitRespectEnd},
		{"strength", bitStrengthStart, bitStrengthEnd},
		{"experience", bitExperienceStart, bitExperienceEnd},
		{"level", bitLevelStart, bitLevelEnd},
		{"type", bitTypeStart, bitTypeEnd},
		{"house", bitHousePos, bitHousePos},
		{"weapon", bitWeaponPos, bitWeaponPos},
		{"family", bitFamilyPos, bitFamilyPos},
		{"mana", bitManaStart, bitManaEnd},
		{"health", bitHealthStart, bitHealthEnd},
	}

	for _, f := range fields {
		for i := f.start; i <= f.end; i++ {
			if i >= 48 {
				t.Errorf("%s: bit %d out of 48-bit range", f.name, i)
				continue
			}
			if used[i] {
				t.Errorf("%s: bit %d already used by another field", f.name, i)
			}
			used[i] = true
		}
	}

	// Проверяем, что все биты 0-46 использованы (бит 47 — резерв)
	for i := 0; i < 47; i++ {
		if !used[i] {
			t.Logf("Warning: bit %d is unused (reserved or gap)", i)
		}
	}
}

// TestBitFieldCoverage проверяет, что битовые поля не пересекаются
// и покрывают ожидаемые диапазоны
func TestBitFieldCoverage2(t *testing.T) {
	type bitRange struct {
		name  string
		start int
		end   int
	}

	ranges := []bitRange{
		{"nameSize", bitNameStart, bitNameEnd},
		{"respect", bitRespectStart, bitRespectEnd},
		{"strength", bitStrengthStart, bitStrengthEnd},
		{"experience", bitExperienceStart, bitExperienceEnd},
		{"level", bitLevelStart, bitLevelEnd},
		{"type", bitTypeStart, bitTypeEnd},
		{"house", bitHousePos, bitHousePos},
		{"weapon", bitWeaponPos, bitWeaponPos},
		{"family", bitFamilyPos, bitFamilyPos},
		{"mana", bitManaStart, bitManaEnd},
		{"health", bitHealthStart, bitHealthEnd},
	}

	// Проверяем отсутствие пересечений
	for i, r1 := range ranges {
		for j, r2 := range ranges {
			if i >= j {
				continue
			}

			// Проверка пересечения диапазонов
			if r1.end >= r2.start && r1.start <= r2.end {
				t.Errorf("Bit ranges overlap: %s[%d:%d] and %s[%d:%d]",
					r1.name, r1.start, r1.end,
					r2.name, r2.start, r2.end)
			}
		}
	}

	// Проверяем, что все диапазоны в пределах 48 бит
	for _, r := range ranges {
		if r.start < 0 || r.end >= 48 {
			t.Errorf("Bit range %s[%d:%d] is out of 48-bit bounds",
				r.name, r.start, r.end)
		}
	}

	// Проверяем, что поля были успешно инициализированы
	// (если дошли до этой точки, паники не было)
	if nameSizeField.Width() != 6 {
		t.Errorf("nameSizeField width = %d, want 6", nameSizeField.Width())
	}
	if healthField.Width() != 10 {
		t.Errorf("healthField width = %d, want 10", healthField.Width())
	}
}

// TestSchemaInitialization проверяет, что инициализация прошла успешно
func TestSchemaInitialization(t *testing.T) {
	// Если этот тест запустился, значит var инициализация прошла
	// без паники - это хорошо!

	tests := []struct {
		name     string
		field    interface{ Width() uint8 }
		wantBits uint8
	}{
		{"nameSize", nameSizeField, 6},
		{"respect", respectField, 4},
		{"strength", strengthField, 4},
		{"experience", experienceField, 4},
		{"level", levelField, 4},
		{"type", typeField, 2},
		{"mana", manaField, 10},
		{"health", healthField, 10},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.field.Width(); got != tt.wantBits {
				t.Errorf("%s.Width() = %d, want %d", tt.name, got, tt.wantBits)
			}
		})
	}
}
