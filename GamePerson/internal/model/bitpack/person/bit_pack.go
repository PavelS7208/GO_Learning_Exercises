package personbitpack

import "GamePerson/internal/bitpack"

//  ------------ Гетеры промежуточного слоя из битов ---------------

func GetNameSize(packed *Packed48) uint32 {
	return bitpack.GetUIntFieldAs[uint32](packed[:], nameSizeField)
}

func GetMana(packed *Packed48) uint32 {
	return bitpack.GetUIntFieldAs[uint32](packed[:], manaField)
}

func GetHealth(packed *Packed48) uint32 {
	return bitpack.GetUIntFieldAs[uint32](packed[:], healthField)
}

func GetRespect(packed *Packed48) uint32 {
	return bitpack.GetUIntFieldAs[uint32](packed[:], respectField)
}
func GetStrength(packed *Packed48) uint32 {
	return bitpack.GetUIntFieldAs[uint32](packed[:], strengthField)
}
func GetExperience(packed *Packed48) uint32 {
	return bitpack.GetUIntFieldAs[uint32](packed[:], experienceField)
}
func GetLevel(packed *Packed48) uint32 {
	return bitpack.GetUIntFieldAs[uint32](packed[:], levelField)
}
func GetType(packed *Packed48) uint32 {
	return bitpack.GetUIntFieldAs[uint32](packed[:], typeField)
}

func GetHouse(packed *Packed48) bool {
	return bitpack.GetBoolField(packed[:], houseField)
}

func GetWeapon(packed *Packed48) bool {
	return bitpack.GetBoolField(packed[:], weaponField)
}
func GetFamily(packed *Packed48) bool {
	return bitpack.GetBoolField(packed[:], familyField)
}

// --------------- Сеттеры промежуточного слоя из битов (с проверкой и без ) -----------------------

func SetSizeName(packed *Packed48, value uint32) error {
	return bitpack.SetUIntFieldAs[uint32](packed[:], nameSizeField, value)
}

func SetMana(packed *Packed48, value uint32) error {
	return bitpack.SetUIntFieldAs[uint32](packed[:], manaField, value)
}
func SetManaUnchecked(packed *Packed48, value uint32) {
	bitpack.SetUIntFieldUncheckedAs[uint32](packed[:], manaField, value)
}

func SetHealth(packed *Packed48, value uint32) error {
	return bitpack.SetUIntFieldAs[uint32](packed[:], healthField, value)
}
func SetHealthUnchecked(packed *Packed48, value uint32) {
	bitpack.SetUIntFieldUncheckedAs[uint32](packed[:], healthField, value)
}

func SetStrength(packed *Packed48, value uint32) error {
	return bitpack.SetUIntFieldAs[uint32](packed[:], strengthField, value)
}
func SetStrengthUnchecked(packed *Packed48, value uint32) {
	bitpack.SetUIntFieldUncheckedAs[uint32](packed[:], strengthField, value)
}

func SetExperience(packed *Packed48, value uint32) error {
	return bitpack.SetUIntFieldAs[uint32](packed[:], experienceField, value)
}
func SetExperienceUnchecked(packed *Packed48, value uint32) {
	bitpack.SetUIntFieldUncheckedAs[uint32](packed[:], experienceField, value)
}

func SetRespect(packed *Packed48, value uint32) error {
	return bitpack.SetUIntFieldAs[uint32](packed[:], respectField, value)
}
func SetRespectUnchecked(packed *Packed48, value uint32) {
	bitpack.SetUIntFieldUncheckedAs[uint32](packed[:], respectField, value)
}

func SetLevel(packed *Packed48, value uint32) error {
	return bitpack.SetUIntFieldAs[uint32](packed[:], levelField, value)
}
func SetLevelUnchecked(packed *Packed48, value uint32) {
	bitpack.SetUIntFieldUncheckedAs[uint32](packed[:], levelField, value)
}

func SetType(packed *Packed48, value uint32) error {
	return bitpack.SetUIntFieldAs[uint32](packed[:], typeField, value)
}
func SetTypeUnchecked(packed *Packed48, value uint32) {
	bitpack.SetUIntFieldUncheckedAs[uint32](packed[:], typeField, value)
}

func SetHouse(packed *Packed48, value bool) error {
	return bitpack.SetBoolField(packed[:], houseField, value)
}
func SetHouseUnchecked(packed *Packed48, value bool) {
	bitpack.SetBoolFieldUnchecked(packed[:], houseField, value)
}

func SetWeapon(packed *Packed48, value bool) error {
	return bitpack.SetBoolField(packed[:], weaponField, value)
}
func SetWeaponUnchecked(packed *Packed48, value bool) {
	bitpack.SetBoolFieldUnchecked(packed[:], weaponField, value)
}

func SetFamily(packed *Packed48, value bool) error {
	return bitpack.SetBoolField(packed[:], familyField, value)
}
func SetFamilyUnchecked(packed *Packed48, value bool) {
	bitpack.SetBoolFieldUnchecked(packed[:], familyField, value)
}
