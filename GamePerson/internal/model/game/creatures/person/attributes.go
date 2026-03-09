package person

import (
	personbitpack "GamePerson/internal/model/bitpack/person"
	"GamePerson/internal/model/config"
	"GamePerson/internal/model/game/creatures/base/entity"
	"fmt"
)

// -----------------   Геттеры --------------------------------------

func (p *person) Name() string {
	length := personbitpack.GetNameSize(&p.packed)
	return string(p.name[:length])
}

func (p *person) Gold() uint32 { return p.gold }
func (p *person) X() int32     { return p.x }
func (p *person) Y() int32     { return p.y }
func (p *person) Z() int32     { return p.z }
func (p *person) Coordinates() (int32, int32, int32) {
	return p.x, p.y, p.z
}

func (p *person) Mana() uint32       { return personbitpack.GetMana(&p.packed) }
func (p *person) Respect() uint32    { return personbitpack.GetRespect(&p.packed) }
func (p *person) Health() uint32     { return personbitpack.GetHealth(&p.packed) }
func (p *person) Strength() uint32   { return personbitpack.GetStrength(&p.packed) }
func (p *person) Experience() uint32 { return personbitpack.GetExperience(&p.packed) }
func (p *person) Level() uint32      { return personbitpack.GetLevel(&p.packed) }
func (p *person) Type() PersonType   { return PersonType(personbitpack.GetType(&p.packed)) }

func (p *person) HasHouse() bool  { return personbitpack.GetHouse(&p.packed) }
func (p *person) HasWeapon() bool { return personbitpack.GetWeapon(&p.packed) }
func (p *person) HasFamily() bool { return personbitpack.GetFamily(&p.packed) }

// Сеттеры для строковой (с сохранением в битовой карте длины)

func (p *person) SetName(name string) error {
	validated, err := entity.ValidateAndCopyName(p.name[:], name)
	if err != nil {
		return err
	}
	// Обновляем длину в битовом поле
	if err = personbitpack.SetSizeName(&p.packed, uint32(len(validated))); err != nil {
		return fmt.Errorf("failed to update name length: %w", err)
	}
	return nil
}

// ------------- Сеттеры для простых полей --------------------

func (p *person) SetX(x int32) error {
	if err := entity.ValidateCoordinate("X", x); err != nil {
		return err
	}
	p.x = x
	return nil
}

func (p *person) SetY(y int32) error {
	if err := entity.ValidateCoordinate("Y", y); err != nil {
		return err
	}
	p.y = y
	return nil
}

func (p *person) SetZ(z int32) error {
	if err := entity.ValidateCoordinate("Z", z); err != nil {
		return err
	}
	p.z = z
	return nil
}

func (p *person) SetGold(gold uint32) error {
	if gold > config.PersonMaxGold {
		return fmt.Errorf("gold %d exceeds maximum %d", gold, config.PersonMaxGold)
	}
	p.gold = gold
	return nil
}

// ------------- Сеттеры для битово-упакованных полей целых --------------------
//  Проверки на бизнес слое, далее Unchecked версии

func (p *person) SetMana(mana uint32) error {
	// Проверка на логическую корректность
	if mana > config.PersonMaxMana {
		return fmt.Errorf("mana %d exceeds maximum %d", mana, config.PersonMaxMana)
	}
	// ГАРАНТИЯ: значение валидно
	personbitpack.SetManaUnchecked(&p.packed, mana)
	return nil
}

func (p *person) SetHealth(health uint32) error {
	if health > config.PersonMaxHealth {
		return fmt.Errorf("health %d exceeds maximum %d", health, config.PersonMaxHealth)
	}
	// ГАРАНТИЯ: значение валидно
	personbitpack.SetHealthUnchecked(&p.packed, health)
	return nil
}

func (p *person) SetStrength(strength uint32) error {
	if strength > config.PersonMaxStrength {
		return fmt.Errorf("strength %d exceeds maximum %d", strength, config.PersonMaxStrength)
	}
	// ГАРАНТИЯ: значение валидно
	personbitpack.SetStrengthUnchecked(&p.packed, strength)
	return nil
}

func (p *person) SetRespect(respect uint32) error {
	if respect > config.PersonMaxRespect {
		return fmt.Errorf("respect %d exceeds maximum %d", respect, config.PersonMaxRespect)
	}
	// ГАРАНТИЯ: значение валидно
	personbitpack.SetRespectUnchecked(&p.packed, respect)
	return nil
}

func (p *person) SetExperience(exp uint32) error {
	if exp > config.PersonMaxExperience {
		return fmt.Errorf("experience %d exceeds maximum %d", exp, config.PersonMaxExperience)
	}
	// ГАРАНТИЯ: значение валидно
	personbitpack.SetExperienceUnchecked(&p.packed, exp)
	return nil
}

func (p *person) SetLevel(level uint32) error {
	if level > config.PersonMaxLevel {
		return fmt.Errorf("level %d exceeds maximum %d", level, config.PersonMaxLevel)
	}
	// ГАРАНТИЯ: значение валидно
	personbitpack.SetLevelUnchecked(&p.packed, level)
	return nil
}

func (p *person) SetType(pt PersonType) error {
	if pt < PersonTypeBuilder || pt > PersonTypeWarrior {
		return fmt.Errorf("invalid person type: %d", pt)
	}
	// ГАРАНТИЯ: значение валидно
	personbitpack.SetTypeUnchecked(&p.packed, uint32(pt))
	return nil
}

func (p *person) SetHouse(has bool) error  { return personbitpack.SetHouse(&p.packed, has) }
func (p *person) SetWeapon(has bool) error { return personbitpack.SetWeapon(&p.packed, has) }
func (p *person) SetFamily(has bool) error { return personbitpack.SetFamily(&p.packed, has) }
