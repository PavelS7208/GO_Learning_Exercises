package monster

import (
	monsterbitpack "GamePerson/internal/model/bitpack/monster"
	"GamePerson/internal/model/config"
	"GamePerson/internal/model/game/creatures/base/entity"
	"errors"
	"fmt"
)

// GameMonster Схема битовой упаковки в 32 битах (4 байт) описана в schema
type monster struct {
	name   [config.MaxNameLength]byte // 42 байта: символы имени латиницей (без указателей!)
	_      [2]byte                    // ← явный паддинг для выравнивания
	packed monsterbitpack.Packed32    // 4 байт: битовая упаковка мелких полей (см. ниже)
	gold   uint32                     // 4 байта: золото [0…2_000_000_000]
	x      int32                      // 4 байта: координата X [-2_000_000_000…2_000_000_000]
	y      int32                      // 4 байта: координата Y
	z      int32                      // 4 байта: координата Z
	// Итого: 42 + 2+ 4 + 4 + 4 + 4 + 4 = 64 байта
}

// Monster представляет игрового персонажа со всеми атрибутами.
// Композиция интерфейсов

type Monster interface {
	entity.Creature
	entity.Wealthy
	entity.Magical
	entity.PropertyOwner
}

// ------------- Конструктор -----------------------------------
// Functional Options Pattern
// Подход для создания персонажа через включение свойств

type Option func(*monster) error

func NewMonster(options ...Option) (Monster, error) {
	m := &monster{}

	// Дефолтные значения должны быть гарантированно валидны
	// Если они невалидны — это баг конфигурации, паникуем
	mustSetDefaults(m)

	for _, option := range options {
		if err := option(m); err != nil {
			return nil, fmt.Errorf("failed to apply option: %w", err)
		}
	}
	return m, nil
}

func mustSetDefaults(m *monster) {
	mdh := config.MonsterDefaultHealth
	if err := m.SetHealth(mdh); err != nil {
		panic(fmt.Sprintf("BUG: invalid default health=%v: %v", mdh, err))
	}
	mdm := config.MonsterDefaultMana
	if err := m.SetMana(mdm); err != nil {
		panic(fmt.Sprintf("BUG: invalid default mana=%v: %v", mdm, err))
	}
}

func (m *monster) String() string {
	if m == nil {
		return "<nil>"
	}

	return fmt.Sprintf(
		"GameMonster{Name: %q, Health: %d/%d, Mana: %d/%d, Gold: %d, House: %v, Pos: (%d,%d,%d)}",
		m.Name(),
		m.Health(),
		config.MonsterMaxHealth,
		m.Mana(),
		config.MonsterMaxMana,
		m.Gold(),
		m.HasHouse(),
		m.X(),
		m.Y(),
		m.Z(),
	)
}

func (m *monster) Validate() error {
	// Собираем ВСЕ ошибки, а не останавливаемся на первой
	var errs []error
	// Инвариант 1: длина имени соответствует буферу
	nameLen := monsterbitpack.GetNameSize(&m.packed)
	if nameLen > config.MaxNameLength {
		errs = append(errs, fmt.Errorf("name length %d exceeds maximum %d", nameLen, config.MaxNameLength))
	}

	// Инвариант 2: байты после имени — нулевые (защита от мусора в буфере)
	for i := nameLen; i < config.MaxNameLength; i++ {
		if m.name[i] != 0 {
			errs = append(errs, fmt.Errorf("name buffer contains non-zero garbage at position %d", i))
			break
		}
	}

	// Инвариант 3: координаты в допустимом диапазоне
	if err := entity.ValidateCoordinate("X", m.x); err != nil {
		errs = append(errs, fmt.Errorf("invalid X coordinate: %w", err))
	}
	if err := entity.ValidateCoordinate("Y", m.y); err != nil {
		errs = append(errs, fmt.Errorf("invalid Y coordinate: %w", err))
	}
	if err := entity.ValidateCoordinate("Z", m.z); err != nil {
		errs = append(errs, fmt.Errorf("invalid Z coordinate: %w", err))
	}
	// Инвариант 4: бизнес-лимиты (дублирующая проверка как защита от багов)
	if health := m.Health(); health > config.MonsterMaxHealth {
		errs = append(errs, fmt.Errorf("health %d exceeds maximum %d", health, config.MonsterMaxHealth))
	}
	if mana := m.Mana(); mana > config.MonsterMaxMana {
		errs = append(errs, fmt.Errorf("mana %d exceeds maximum %d", mana, config.MonsterMaxMana))
	}
	if gold := m.Gold(); gold > config.MonsterMaxGold {
		errs = append(errs, fmt.Errorf("gold %d exceeds maximum %d", gold, config.MonsterMaxGold))
	}
	if len(errs) > 0 {
		return errors.Join(errs...)
	}
	return nil
}
