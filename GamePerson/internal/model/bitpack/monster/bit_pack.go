package monsterbitpack

import "GamePerson/internal/bitpack"

//  ------------ Гетеры промежуточного слоя из битов ---------------

func GetNameSize(packed *Packed32) uint32 {
	return bitpack.GetUIntFieldAs[uint32](packed[:], nameSizeField)
}

func GetHealth(packed *Packed32) uint32 {
	return bitpack.GetUIntFieldAs[uint32](packed[:], healthField)
}

func GetMana(packed *Packed32) uint32 {
	return bitpack.GetUIntFieldAs[uint32](packed[:], manaField)
}

func GetHouse(packed *Packed32) bool {
	return bitpack.GetBoolField(packed[:], houseField)
}

// --------------- Сеттеры промежуточного слоя из битов -----------------------

func SetSizeName(packed *Packed32, value uint32) error {
	return bitpack.SetUIntFieldAs[uint32](packed[:], nameSizeField, value)
}

func SetMana(packed *Packed32, value uint32) error {
	return bitpack.SetUIntFieldAs[uint32](packed[:], manaField, value)
}
func SetManaUnchecked(packed *Packed32, value uint32) {
	bitpack.SetUIntFieldUncheckedAs[uint32](packed[:], manaField, value)
}

func SetHealth(packed *Packed32, value uint32) error {
	return bitpack.SetUIntFieldAs[uint32](packed[:], healthField, value)
}
func SetHealthUnchecked(packed *Packed32, value uint32) {
	bitpack.SetUIntFieldUncheckedAs[uint32](packed[:], healthField, value)
}

func SetHouse(packed *Packed32, value bool) error {
	return bitpack.SetBoolField(packed[:], houseField, value)
}
func SetHouseUnchecked(packed *Packed32, value bool) {
	bitpack.SetBoolFieldUnchecked(packed[:], houseField, value)
}
