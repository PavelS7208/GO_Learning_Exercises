package monster

import (
	monsterbitpack "GamePerson/internal/model/bitpack/monster"
	"GamePerson/internal/model/config"
	"GamePerson/internal/model/game/creatures/base/entity"
	"fmt"
)

// -----------------   Геттеры --------------------------------------

func (m *monster) Name() string {
	length := monsterbitpack.GetNameSize(&m.packed)
	return string(m.name[:length])
}

func (m *monster) Gold() uint32                       { return m.gold }
func (m *monster) X() int32                           { return m.x }
func (m *monster) Y() int32                           { return m.y }
func (m *monster) Z() int32                           { return m.z }
func (m *monster) Coordinates() (int32, int32, int32) { return m.x, m.y, m.z }

func (m *monster) Mana() uint32   { return monsterbitpack.GetMana(&m.packed) }
func (m *monster) Health() uint32 { return monsterbitpack.GetHealth(&m.packed) }
func (m *monster) HasHouse() bool { return monsterbitpack.GetHouse(&m.packed) }

// ------------- Сеттеры для простых полей --------------------

func (m *monster) SetX(x int32) error {
	if err := entity.ValidateCoordinate("X", x); err != nil {
		return err
	}
	m.x = x
	return nil
}

func (m *monster) SetY(y int32) error {
	if err := entity.ValidateCoordinate("Y", y); err != nil {
		return err
	}
	m.y = y
	return nil
}

func (m *monster) SetZ(z int32) error {
	if err := entity.ValidateCoordinate("Z", z); err != nil {
		return err
	}
	m.z = z
	return nil
}

func (m *monster) SetGold(gold uint32) error {
	if gold > config.MonsterMaxGold {
		return fmt.Errorf("gold %d exceeds maximum %d", gold, config.MonsterMaxGold)
	}
	m.gold = gold
	return nil
}

// Сеттеры для строкового поля (с сохранением в битовой карте длины)

func (m *monster) SetName(name string) error {
	validated, err := entity.ValidateAndCopyName(m.name[:], name)
	if err != nil {
		return err
	}
	// Обновляем длину в битовом поле
	if err = monsterbitpack.SetSizeName(&m.packed, uint32(len(validated))); err != nil {
		return fmt.Errorf("failed to update name length: %w", err)
	}
	return nil
}

// ------------- Сеттеры для битово-упакованных полей целых --------------------

func (m *monster) SetMana(mana uint32) error {
	// Проверка на логическую корректность
	if mana > config.MonsterMaxMana {
		return fmt.Errorf("mana %d exceeds maximum %d", mana, config.MonsterMaxMana)
	}
	// ГАРАНТИЯ: значение валидно
	monsterbitpack.SetManaUnchecked(&m.packed, mana)
	return nil
}

func (m *monster) SetHealth(health uint32) error {
	if health > config.MonsterMaxHealth {
		return fmt.Errorf("health %d exceeds maximum %d", health, config.MonsterMaxHealth)
	}
	// ГАРАНТИЯ: значение валидно
	monsterbitpack.SetHealthUnchecked(&m.packed, health)
	return nil
}

func (m *monster) SetHouse(has bool) error { return monsterbitpack.SetHouse(&m.packed, has) }
