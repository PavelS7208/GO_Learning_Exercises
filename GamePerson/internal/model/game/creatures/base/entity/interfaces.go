package entity

// ============ Базовые интерфейсы (атомарные) ============

type Named interface {
	Name() string
}

type Positioned interface {
	X() int32
	Y() int32
	Z() int32
	Coordinates() (x, y, z int32)
}

type Movable interface {
	Positioned
	SetX(x int32) error
	SetY(y int32) error
	SetZ(z int32) error
}

type Wealthy interface {
	Gold() uint32
	SetGold(uint32) error
}

type Living interface {
	Health() uint32
	SetHealth(uint32) error
}

type Magical interface {
	Mana() uint32
	SetMana(uint32) error
}

type Experienced interface {
	Experience() uint32
	SetExperience(uint32) error
	Level() uint32
	SetLevel(uint32) error
}

type Strong interface {
	Strength() uint32
	SetStrength(uint32) error
}

type Reputable interface {
	Respect() uint32
	SetRespect(uint32) error
}

type PropertyOwner interface {
	HasHouse() bool
	SetHouse(bool) error
}

type Armed interface {
	HasWeapon() bool
	SetWeapon(bool) error
}

type FamilyMember interface {
	HasFamily() bool
	SetFamily(bool) error
}

// ============ Композитные интерфейсы ============

// Базовая сущность игрового мира

type Entity interface {
	Named
	Positioned
}

// Существо, которое может перемещаться

type Creature interface {
	Entity
	Movable
	Living
}

// Боевая единица

type Combatant interface {
	Creature
	Strong
	Armed
}

// Validatable гарантирует, что объект находится в валидном состоянии
// после создания или десериализации
type Validatable interface {
	Validate() error
}
