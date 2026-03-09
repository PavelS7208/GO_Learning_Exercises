package person

import (
	"testing"
	"unsafe"

	"github.com/stretchr/testify/assert"
)

func TestGamePersonSize(t *testing.T) {
	const expectedSize = 64
	actual := unsafe.Sizeof(person{})
	if actual != expectedSize {
		t.Fatalf("GamePerson size MUST be %d bytes, got %d. "+
			"Check field order and alignment!", expectedSize, actual)
	}
}

// Build-time assertion (падает при go build, если условие нарушено)
//var _ [64]byte = [unsafe.Sizeof(GamePerson{})]byte{} // ← компилятор выдаст ошибку

// Тест из задания курса
func TestGamePerson(t *testing.T) {

	//assertSize()

	assert.LessOrEqual(t, unsafe.Sizeof(person{}), uintptr(64))

	const x, y, z int32 = -2_000_000_000, 2_000_000_000, 0
	const name = "aaaaaaaaaaaaa_bbbbbbbbbbbbb_cccccccccccccc"
	const personType = PersonTypeBuilder
	//const gold = math.MaxInt32
	const gold = 2_000_000_000
	const mana = 1000
	const health = 1000
	const respect = 10
	const strength = 10
	const experience = 10
	const level = 10

	options := []Option{
		WithName(name),
		WithCoordinates(x, y, z),
		WithGold(gold),
		WithMana(mana),
		WithHealth(health),
		WithRespect(respect),
		WithStrength(strength),
		WithExperience(experience),
		WithLevel(level),
		WithHouse(true),
		WithWeapon(false),
		WithFamily(true),
		WithType(personType),
	}

	person, _ := NewPerson(options...)
	assert.Equal(t, name, person.Name())
	assert.Equal(t, x, person.X())
	assert.Equal(t, y, person.Y())
	assert.Equal(t, z, person.Z())
	assert.Equal(t, gold, int(person.Gold()))
	assert.Equal(t, mana, int(person.Mana()))
	assert.Equal(t, health, int(person.Health()))
	assert.Equal(t, respect, int(person.Respect()))
	assert.Equal(t, strength, int(person.Strength()))
	assert.Equal(t, experience, int(person.Experience()))
	assert.Equal(t, level, int(person.Level()))
	assert.True(t, person.HasHouse())
	assert.True(t, person.HasFamily())
	assert.False(t, person.HasWeapon())
	assert.Equal(t, personType, person.Type())
}
