package person

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPersonJSONSerialization(t *testing.T) {
	// Создаем Person
	p, err := NewPerson(
		WithName("TestWarrior"),
		WithType(PersonTypeWarrior),
		WithHealth(800),
		WithMana(500),
		WithLevel(5),
		WithGold(1000),
		WithStrength(8),
		WithRespect(7),
		WithExperience(6),
		WithWeapon(true),
		WithHouse(true),
		WithFamily(false),
		WithCoordinates(10, 20, 30),
	)
	require.NoError(t, err)

	serializer := NewSerializer()

	// Сериализация в JSON
	jsonData, err := serializer.ToJSON(p)
	require.NoError(t, err)
	require.NotEmpty(t, jsonData)

	// Десериализация из JSON
	pFromJSON, err := serializer.FromJSON(jsonData)
	require.NoError(t, err)
	require.NotNil(t, pFromJSON)

	assertPersonEqual(t, p, pFromJSON)
}

func TestPersonXMLSerialization(t *testing.T) {
	// Создаем Person
	p, err := NewPerson(
		WithName("TestBuilder"),
		WithType(PersonTypeBuilder),
		WithHealth(600),
		WithMana(200),
		WithLevel(3),
		WithGold(500),
		WithStrength(5),
		WithRespect(4),
		WithExperience(3),
		WithWeapon(false),
		WithHouse(true),
		WithFamily(true),
		WithCoordinates(-10, -20, -30),
	)
	require.NoError(t, err)

	serializer := NewSerializer()

	// Сериализация в XML
	xmlData, err := serializer.ToXML(p)
	require.NoError(t, err)
	require.NotEmpty(t, xmlData)

	// Десериализация из XML
	pFromXML, err := serializer.FromXML(xmlData)
	require.NoError(t, err)
	require.NotNil(t, pFromXML)

	assertPersonEqual(t, p, pFromXML)

}

func TestPersonYAMLSerialization(t *testing.T) {
	// Создаем Person
	p, err := NewPerson(
		WithName("TestBlacksmith"),
		WithType(PersonTypeBlacksmith),
		WithHealth(700),
		WithMana(300),
		WithLevel(4),
		WithGold(2000),
		WithStrength(9),
		WithRespect(8),
		WithExperience(7),
		WithWeapon(true),
		WithHouse(false),
		WithFamily(true),
		WithCoordinates(100, 200, 300),
	)
	require.NoError(t, err)

	serializer := NewSerializer()

	// Сериализация в YAML
	yamlData, err := serializer.ToYAML(p)
	require.NoError(t, err)
	require.NotEmpty(t, yamlData)

	// Десериализация из YAML
	pFromYAML, err := serializer.FromYAML(yamlData)
	require.NoError(t, err)
	require.NotNil(t, pFromYAML)

	assertPersonEqual(t, p, pFromYAML)

}

func TestPersonRoundTripAllFormats(t *testing.T) {
	// Создаем Person со всеми полями
	original, err := NewPerson(
		WithName("RoundTripTest"),
		WithType(PersonTypeWarrior),
		WithHealth(999),
		WithMana(999),
		WithLevel(10),
		WithGold(1000000),
		WithStrength(10),
		WithRespect(10),
		WithExperience(10),
		WithWeapon(true),
		WithHouse(true),
		WithFamily(true),
		WithCoordinates(12345, -54321, 99999),
	)
	require.NoError(t, err)

	serializer := NewSerializer()

	// Тест JSON
	jsonData, err := serializer.ToJSON(original)
	require.NoError(t, err)
	fromJSON, err := serializer.FromJSON(jsonData)
	require.NoError(t, err)
	assertPersonEqual(t, original, fromJSON)

	// Тест XML
	xmlData, err := serializer.ToXML(original)
	require.NoError(t, err)
	fromXML, err := serializer.FromXML(xmlData)
	require.NoError(t, err)
	assertPersonEqual(t, original, fromXML)

	// Тест YAML
	yamlData, err := serializer.ToYAML(original)
	require.NoError(t, err)
	fromYAML, err := serializer.FromYAML(yamlData)
	require.NoError(t, err)
	assertPersonEqual(t, original, fromYAML)
}

// assertPersonEqual проверяет равенство всех полей Person
func assertPersonEqual(t *testing.T, expected, actual Person) {
	assert.Equal(t, expected.Name(), actual.Name(), "Name mismatch")
	assert.Equal(t, expected.Type(), actual.Type(), "Type mismatch")
	assert.Equal(t, expected.Health(), actual.Health(), "Health mismatch")
	assert.Equal(t, expected.Mana(), actual.Mana(), "Mana mismatch")
	assert.Equal(t, expected.Level(), actual.Level(), "Level mismatch")
	assert.Equal(t, expected.Gold(), actual.Gold(), "Gold mismatch")
	assert.Equal(t, expected.Strength(), actual.Strength(), "Strength mismatch")
	assert.Equal(t, expected.Respect(), actual.Respect(), "Respect mismatch")
	assert.Equal(t, expected.Experience(), actual.Experience(), "Experience mismatch")
	assert.Equal(t, expected.HasWeapon(), actual.HasWeapon(), "HasWeapon mismatch")
	assert.Equal(t, expected.HasHouse(), actual.HasHouse(), "HasHouse mismatch")
	assert.Equal(t, expected.HasFamily(), actual.HasFamily(), "HasFamily mismatch")
	assert.Equal(t, expected.X(), actual.X(), "X mismatch")
	assert.Equal(t, expected.Y(), actual.Y(), "Y mismatch")
	assert.Equal(t, expected.Z(), actual.Z(), "Z mismatch")
}
