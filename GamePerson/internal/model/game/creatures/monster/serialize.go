package monster

import (
	"GamePerson/internal/model/game/creatures/base/entity"
	"GamePerson/internal/model/game/creatures/base/serializer"
)

// ===== Основные конструкторы Person из внешних представлений (99% случаев — без параметров) ============

// NewFromJSON создаёт персонажа из JSON с проверкой целостности по умолчанию
func NewFromJSON(data []byte) (Monster, error) {
	return newFromJSON(data, nil)
}

// NewFromXML создаёт персонажа из XML с проверкой целостности по умолчанию
func NewFromXML(data []byte) (Monster, error) {
	return newFromXML(data, nil)
}

// NewFromYAML создаёт персонажа из YAML с проверкой целостности по умолчанию
func NewFromYAML(data []byte) (Monster, error) {
	return newFromYAML(data, nil)
}

// ============ Расширенные конструкторы (специальные случаи) ============

// NewFromJSONWithIntegrity создаёт персонажа из JSON с кастомным проверяющим целостности
func NewFromJSONWithIntegrity(data []byte, integrity *entity.IntegrityChecker) (Monster, error) {
	return newFromJSON(data, integrity)
}

// NewFromXMLWithIntegrity создаёт персонажа из XML с кастомным проверяющим целостности
func NewFromXMLWithIntegrity(data []byte, integrity *entity.IntegrityChecker) (Monster, error) {
	return newFromXML(data, integrity)
}

// NewFromYAMLWithIntegrity создаёт персонажа из YAML с кастомным проверяющим целостности
func NewFromYAMLWithIntegrity(data []byte, integrity *entity.IntegrityChecker) (Monster, error) {
	return newFromYAML(data, integrity)
}

// ---------------  Внутренняя реализация конструкторов -------------------------------

func newFromJSON(data []byte, integrity *entity.IntegrityChecker) (Monster, error) {
	ser := NewSerializer(integrity)
	return ser.FromJSON(data)
}

func newFromXML(data []byte, integrity *entity.IntegrityChecker) (Monster, error) {
	ser := NewSerializer(integrity)
	return ser.FromXML(data)
}

func newFromYAML(data []byte, integrity *entity.IntegrityChecker) (Monster, error) {
	ser := NewSerializer(integrity)
	return ser.FromYAML(data)
}

type monsterConverter struct{}

func (c monsterConverter) ToDTO(m Monster) MonsterDTO              { return ToDTO(m) }
func (c monsterConverter) FromDTO(dto MonsterDTO) (Monster, error) { return FromDTO(dto) }

func NewSerializer(integrity *entity.IntegrityChecker) *serializer.Serializer[Monster, MonsterDTO] {
	return serializer.New[Monster, MonsterDTO](monsterConverter{}, integrity)
}
