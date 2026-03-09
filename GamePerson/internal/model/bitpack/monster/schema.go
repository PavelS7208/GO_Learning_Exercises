package monsterbitpack

import (
	"GamePerson/internal/bitpack"
	"GamePerson/internal/model/config"
)

type Packed32 = bitpack.Packed32

const (
	bitNameStart = 0
	bitNameEnd   = 5

	bitManaStart = 6
	bitManaEnd   = 15

	bitHealthStart = 16
	bitHealthEnd   = 29

	bitHousePos = 30
)

// ВАЖНО: Все поля должны быть валидны при компиляции!
// Используем метод Must* с паникой для проверки при старте
// Изменение констант bit*Start/bit*End требует:
//  1. Обновления комментария схемы упаковки
//  2. Запуска `go test -run TestBitFieldCoverage`
//  3. Проверки отсутствия пересечений битов
var (
	nameSizeField = bitpack.MustNewUIntBitField(bitNameStart, bitNameEnd, uint64(config.MaxNameLength))
	manaField     = bitpack.MustNewUIntBitField(bitManaStart, bitManaEnd, uint64(config.MonsterMaxMana))
	healthField   = bitpack.MustNewUIntBitField(bitHealthStart, bitHealthEnd, uint64(config.MonsterMaxHealth))
	houseField    = bitpack.MustNewBoolBitField(bitHousePos)
)
