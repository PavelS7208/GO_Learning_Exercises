package personbitpack

import (
	"GamePerson/internal/bitpack"
	"GamePerson/internal/model/config"
)

// ================= Схема битовой упаковки для персонажа ========================
//  В 48 битах (6 байт) храним:

type Packed48 = bitpack.Packed48

// Биты 0- 5: длина имени (6 бит, макс 63 → достаточно для 42)
// Биты 6- 9: уважение (4 бита, 0-15 → покрывает 0-10)
// Биты 10-13: сила (4 бита, 0-15 → покрывает 0-10)
// Биты 14-17: опыт (4 бита, 0-15 → покрывает 0-10)
// Биты 18-21: уровень (4 бита, 0-15 → покрывает 0-10)
// Биты 22-23: тип игрока (2 бита, 0-3 → покрывает 3 варианта)
// Бит 24: есть дом (1 бит)
// Бит 25: есть оружие (1 бит)
// Бит 26: есть семья (1 бит)
// Биты 27-36: мана (10 бит, 0-1023 → покрывает 0-1000)
// Биты 37-46: здоровье (10 бит, 0-1023 → покрывает 0-1000)
// Бит 47: резерв (всегда 0)

const (
	bitNameStart = 0
	bitNameEnd   = 5

	bitRespectStart = 6
	bitRespectEnd   = 9

	bitStrengthStart = 10
	bitStrengthEnd   = 13

	bitExperienceStart = 14
	bitExperienceEnd   = 17

	bitLevelStart = 18
	bitLevelEnd   = 21

	bitTypeStart = 22
	bitTypeEnd   = 23

	bitHousePos  = 24 // 1 bit
	bitWeaponPos = 25 // 1 bit
	bitFamilyPos = 26 // 1 bit

	bitManaStart = 27
	bitManaEnd   = 36

	bitHealthStart = 37
	bitHealthEnd   = 46
)

// Битовые поля инициализируются при загрузке пакета.
// MustNew* функции паникуют при ошибках конфигурации.
//
// ВАЖНО: Паника здесь - это ПРАВИЛЬНОЕ поведение.
// Если эти константы невалидны, это баг программиста,
// а не runtime ошибка. Программа должна упасть немедленно
// при запуске, а не работать с поврежденными данными.
var (
	nameSizeField   = bitpack.MustNewUIntBitField(bitNameStart, bitNameEnd, uint64(config.MaxNameLength))
	manaField       = bitpack.MustNewUIntBitField(bitManaStart, bitManaEnd, uint64(config.PersonMaxMana))
	healthField     = bitpack.MustNewUIntBitField(bitHealthStart, bitHealthEnd, uint64(config.PersonMaxHealth))
	respectField    = bitpack.MustNewUIntBitField(bitRespectStart, bitRespectEnd, uint64(config.PersonMaxRespect))
	strengthField   = bitpack.MustNewUIntBitField(bitStrengthStart, bitStrengthEnd, uint64(config.PersonMaxStrength))
	experienceField = bitpack.MustNewUIntBitField(bitExperienceStart, bitExperienceEnd, uint64(config.PersonMaxExperience))
	levelField      = bitpack.MustNewUIntBitField(bitLevelStart, bitLevelEnd, uint64(config.PersonMaxLevel))
	typeField       = bitpack.MustNewUIntBitField(bitTypeStart, bitTypeEnd, uint64(config.PersonMaxTypeIndex))
	houseField      = bitpack.MustNewBoolBitField(bitHousePos)
	weaponField     = bitpack.MustNewBoolBitField(bitWeaponPos)
	familyField     = bitpack.MustNewBoolBitField(bitFamilyPos)
)
