package config

// Базовые лимиты

// Глобальные координатные лимиты (единые для всех)
const (
	MinCoord int32 = -2_000_000_000
	MaxCoord int32 = 2_000_000_000
)

const (
	MaxNameLength = 42
)

// Лимиты персонажей
const (
	PersonMaxGold       uint32 = 2_000_000_000
	PersonMaxHealth     uint32 = 1000
	PersonMaxMana       uint32 = 1000
	PersonMaxRespect    uint32 = 10
	PersonMaxStrength   uint32 = 10
	PersonMaxExperience uint32 = 10
	PersonMaxLevel      uint32 = 10
	PersonMaxTypeIndex  uint32 = 3
)

// Лимиты монстров
const (
	MonsterMaxHealth uint32 = 10000
	MonsterMaxGold   uint32 = 10000
	MonsterMaxMana   uint32 = 1000
)

const (
	PersonDefaultHealth uint32 = 100
	PersonDefaultMana   uint32 = 10
	PersonDefaultLevel  uint32 = 1
)

const (
	MonsterDefaultHealth uint32 = 300
	MonsterDefaultMana   uint32 = 100
)
