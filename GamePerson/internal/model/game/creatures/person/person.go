package person

import (
	personbitpack "GamePerson/internal/model/bitpack/person"
	"GamePerson/internal/model/config"
	"GamePerson/internal/model/game/creatures/base/entity"
	"errors"

	"fmt"
)

// ================  GamePerson ==============================================

type PersonType uint

// GamePersonType — тип игрока
const (
	PersonTypeBuilder    PersonType = iota // 0: строитель
	PersonTypeBlacksmith                   // 1: кузнец
	PersonTypeWarrior                      // 2: воин
)

func (t PersonType) String() string {
	switch t {
	case PersonTypeBuilder:
		return "Builder"
	case PersonTypeBlacksmith:
		return "Blacksmith"
	case PersonTypeWarrior:
		return "Warrior"
	default:
		return fmt.Sprintf("PersonType(%d)", t)
	}
}

//==============================================================

// person Схема битовой упаковки в 48 битах (6 байт) описана в schema
//------ Атрибуты хранящиеся в запакованной части структуры ---------------------
//  Здоровье (0-1000)
//  Уважение (0-10)
//  Сила (0-10)
//  Опыт (0-10)
//  Уровень (0-10)
//  Мана (0-1000)
//  Тип игрока (4 варианта максимум)
//   Есть дом (булева)
//  Есть оружие (булева)
//  Есть семья (булева)

type person struct {
	name   [config.MaxNameLength]byte // 42 байта: символы имени латиницей (без указателей!)
	packed personbitpack.Packed48     // 6 байт: битовая упаковка мелких полей (см. ниже)
	gold   uint32                     // 4 байта: золото [0…2_000_000_000]
	x      int32                      // 4 байта: координата X [-2_000_000_000…2_000_000_000]
	y      int32                      // 4 байта: координата Y
	z      int32                      // 4 байта: координата Z
	// Итого: 42 + 6 + 4 + 4 + 4 + 4 = 64 байта
}

// Person представляет игрового персонажа со всеми атрибутами.
// Композиция интерфейсов обеспечивает:
//   - Боевые способности (Combatant)
//   - Управление ресурсами (Wealthy, Magical)
//   - Прогрессию (Experienced)
//   - Социальные аспекты (Reputable, PropertyOwner, FamilyMember)
type Person interface {
	entity.Combatant
	entity.Wealthy
	entity.Magical
	entity.Experienced
	entity.Reputable
	entity.PropertyOwner
	entity.FamilyMember

	Type() PersonType
	SetType(PersonType) error
}

// ------------- Конструктор -----------------------------------
// Functional Options Pattern
// Подход для создания персонажа через включение свойств

type Option func(*person) error

func NewPerson(options ...Option) (Person, error) {
	p := &person{}

	// Дефолтные значения должны быть гарантированно валидны
	// Если они невалидны — это баг конфигурации, паникуем
	mustSetDefaults(p)

	for _, option := range options {
		if err := option(p); err != nil {
			return nil, fmt.Errorf("failed to apply option: %w", err)
		}
	}
	return p, nil
}

//----------------------  Сервисные методы person -----------------------------------

func mustSetDefaults(p *person) {
	if err := p.SetType(PersonTypeBuilder); err != nil {
		panic(fmt.Sprintf("BUG: invalid default PersonType=%v: %v", PersonTypeBuilder, err))
	}
	if err := p.SetHealth(config.PersonDefaultHealth); err != nil {
		panic(fmt.Sprintf("BUG: invalid default health=%v: %v", config.PersonDefaultHealth, err))
	}
	if err := p.SetMana(config.PersonDefaultMana); err != nil {
		panic(fmt.Sprintf("BUG: invalid default mana=%v: %v", config.PersonDefaultMana, err))
	}
	if err := p.SetLevel(config.PersonDefaultLevel); err != nil {
		panic(fmt.Sprintf("BUG: invalid default level=%v: %v", config.PersonDefaultLevel, err))
	}
}

func (p *person) rawNameBytes() [config.MaxNameLength]byte {
	return p.name
}

func (p *person) String() string {
	if p == nil {
		return "<nil>"
	}

	return fmt.Sprintf(
		"GamePerson{Name: %q, PersonType: %s, Health: %d/%d, Mana: %d/%d, Level: %d, Gold: %d, House: %v, Weapon: %v, Family: %v, Pos: (%d,%d,%d)}",
		p.Name(),
		p.Type(),
		p.Health(),
		config.PersonMaxHealth,
		p.Mana(),
		config.PersonMaxMana,
		p.Level(),
		p.Gold(),
		p.HasHouse(),
		p.HasWeapon(),
		p.HasFamily(),
		p.X(),
		p.Y(),
		p.Z(),
	)
}

func (p *person) internalData() string {
	return fmt.Sprintf(
		"Raw{name: %q, packed: %v, gold: %d, coords: (%d,%d,%d)}",
		p.rawNameBytes(),
		p.packed,
		p.gold,
		p.x, p.y, p.z,
	)
}

// Validate (полная) - для проверки целостности после сериализаций
func (p *person) Validate() error {
	// Собираем ВСЕ ошибки, а не останавливаемся на первой
	var errs []error
	// Инвариант 1: длина имени соответствует буферу
	nameLen := personbitpack.GetNameSize(&p.packed)
	if nameLen > config.MaxNameLength {
		errs = append(errs, fmt.Errorf("name length %d exceeds maximum %d", nameLen, config.MaxNameLength))
	}

	// Инвариант 2: байты после имени — нулевые (защита от мусора в буфере)
	for i := nameLen; i < config.MaxNameLength; i++ {
		if p.name[i] != 0 {
			errs = append(errs, fmt.Errorf("name buffer contains garbage at position %d", i))
			break // достаточно одной ошибки для мусора
		}
	}

	// Инвариант 3: координаты в допустимом диапазоне
	if err := entity.ValidateCoordinate("X", p.x); err != nil {
		errs = append(errs, err)
	}
	if err := entity.ValidateCoordinate("Y", p.y); err != nil {
		errs = append(errs, err)

	}
	if err := entity.ValidateCoordinate("Z", p.z); err != nil {
		errs = append(errs, err)
	}

	// Инвариант 4: бизнес-лимиты (дублирующая проверка как защита от багов)
	if health := p.Health(); health > config.PersonMaxHealth {
		errs = append(errs, fmt.Errorf("health %d exceeds maximum %d", health, config.PersonMaxHealth))
	}
	if mana := p.Mana(); mana > config.PersonMaxMana {
		errs = append(errs, fmt.Errorf("mana %d exceeds maximum %d", mana, config.PersonMaxMana))
	}
	if gold := p.Gold(); gold > config.PersonMaxGold {
		errs = append(errs, fmt.Errorf("gold %d exceeds maximum %d", gold, config.PersonMaxGold))
	}

	pt := p.Type()
	if pt < PersonTypeBuilder || pt > PersonTypeWarrior {
		errs = append(errs, fmt.Errorf("invalid person type: %d", pt))
	}

	if respect := p.Respect(); respect > config.PersonMaxRespect {
		errs = append(errs, fmt.Errorf("respect %d exceeds maximum %d", respect, config.PersonMaxRespect))
	}
	if strength := p.Strength(); strength > config.PersonMaxStrength {
		errs = append(errs, fmt.Errorf("strength %d exceeds maximum %d", strength, config.PersonMaxStrength))
	}

	if exp := p.Experience(); exp > config.PersonMaxExperience {
		errs = append(errs, fmt.Errorf("exp %d exceeds maximum %d", exp, config.PersonMaxExperience))
	}

	if level := p.Level(); level > config.PersonMaxLevel {
		errs = append(errs, fmt.Errorf("level %d exceeds maximum %d", level, config.PersonMaxLevel))
	}
	if len(errs) > 0 {
		return errors.Join(errs...)
	}
	return nil
}
