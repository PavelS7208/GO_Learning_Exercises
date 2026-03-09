package monster

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMonsterJSONSerialization(t *testing.T) {
	// Создаем Monster
	m, err := NewMonster(
		WithName("TestDragon"),
		WithHealth(800),
		WithMana(500),
		WithGold(10000),
		WithHouse(true),
		WithCoordinates(10, 20, 30),
	)
	require.NoError(t, err)

	serializer := NewSerializer()

	// Сериализация в JSON
	jsonData, err := serializer.ToJSON(m)
	require.NoError(t, err)
	require.NotEmpty(t, jsonData)

	// Десериализация из JSON
	mFromJSON, err := serializer.FromJSON(jsonData)
	require.NoError(t, err)
	require.NotNil(t, mFromJSON)

	assertMonsterEqual(t, m, mFromJSON)
}

func TestMonsterXMLSerialization(t *testing.T) {
	// Создаем Monster
	m, err := NewMonster(
		WithName("TestGoblin"),
		WithHealth(200),
		WithMana(100),
		WithGold(50),
		WithHouse(false),
		WithCoordinates(-10, -20, -30),
	)
	require.NoError(t, err)

	serializer := NewSerializer()

	// Сериализация в XML
	xmlData, err := serializer.ToXML(m)
	require.NoError(t, err)
	require.NotEmpty(t, xmlData)

	// Десериализация из XML
	mFromXML, err := serializer.FromXML(xmlData)
	require.NoError(t, err)
	require.NotNil(t, mFromXML)

	assertMonsterEqual(t, m, mFromXML)
}

func TestMonsterYAMLSerialization(t *testing.T) {
	// Создаем Monster
	m, err := NewMonster(
		WithName("TestTroll"),
		WithHealth(600),
		WithMana(300),
		WithGold(5000),
		WithHouse(true),
		WithCoordinates(100, 200, 300),
	)
	require.NoError(t, err)

	serializer := NewSerializer()

	// Сериализация в YAML
	yamlData, err := serializer.ToYAML(m)
	require.NoError(t, err)
	require.NotEmpty(t, yamlData)

	// Десериализация из YAML
	mFromYAML, err := serializer.FromYAML(yamlData)
	require.NoError(t, err)
	require.NotNil(t, mFromYAML)

	// Проверка всех полей
	assertMonsterEqual(t, m, mFromYAML)

}

func TestMonsterRoundTripAllFormats(t *testing.T) {
	// Создаем Monster со всеми полями
	original, err := NewMonster(
		WithName("RoundTripMonster"),
		WithHealth(999),
		WithMana(999),
		WithGold(1000000),
		WithHouse(true),
		WithCoordinates(12345, -54321, 99999),
	)
	require.NoError(t, err)

	serializer := NewSerializer()

	// Тест JSON
	jsonData, err := serializer.ToJSON(original)
	require.NoError(t, err)
	fromJSON, err := serializer.FromJSON(jsonData)
	require.NoError(t, err)
	assertMonsterEqual(t, original, fromJSON)

	// Тест XML
	xmlData, err := serializer.ToXML(original)
	require.NoError(t, err)
	fromXML, err := serializer.FromXML(xmlData)
	require.NoError(t, err)
	assertMonsterEqual(t, original, fromXML)

	// Тест YAML
	yamlData, err := serializer.ToYAML(original)
	require.NoError(t, err)
	fromYAML, err := serializer.FromYAML(yamlData)
	require.NoError(t, err)
	assertMonsterEqual(t, original, fromYAML)
}

// assertMonsterEqual проверяет равенство всех полей Monster
func assertMonsterEqual(t *testing.T, expected, actual Monster) {
	assert.Equal(t, expected.Name(), actual.Name(), "Name mismatch")
	assert.Equal(t, expected.Health(), actual.Health(), "Health mismatch")
	assert.Equal(t, expected.Mana(), actual.Mana(), "Mana mismatch")
	assert.Equal(t, expected.Gold(), actual.Gold(), "Gold mismatch")
	assert.Equal(t, expected.HasHouse(), actual.HasHouse(), "HasHouse mismatch")
	assert.Equal(t, expected.X(), actual.X(), "X mismatch")
	assert.Equal(t, expected.Y(), actual.Y(), "Y mismatch")
	assert.Equal(t, expected.Z(), actual.Z(), "Z mismatch")
}
