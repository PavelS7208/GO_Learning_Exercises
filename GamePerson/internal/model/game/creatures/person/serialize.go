package person

import (
	"GamePerson/internal/model/game/creatures/base/entity"
	"GamePerson/internal/model/game/creatures/base/serializer"
)

// ===== Основные конструкторы Person из внешних представлений (99% случаев — без параметров) ============

// NewFromJSON создаёт персонажа из JSON с проверкой целостности по умолчанию
func NewFromJSON(data []byte) (Person, error) {
	return newFromJSON(data, nil)
}

// NewFromXML создаёт персонажа из XML с проверкой целостности по умолчанию
func NewFromXML(data []byte) (Person, error) {
	return newFromXML(data, nil)
}

// NewFromYAML создаёт персонажа из YAML с проверкой целостности по умолчанию
func NewFromYAML(data []byte) (Person, error) {
	return newFromYAML(data, nil)
}

// ============ Расширенные конструкторы (специальные случаи) ============

// NewFromJSONWithIntegrity создаёт персонажа из JSON с кастомным проверяющим целостности
func NewFromJSONWithIntegrity(data []byte, integrity *entity.IntegrityChecker) (Person, error) {
	return newFromJSON(data, integrity)
}

// NewFromXMLWithIntegrity создаёт персонажа из XML с кастомным проверяющим целостности
func NewFromXMLWithIntegrity(data []byte, integrity *entity.IntegrityChecker) (Person, error) {
	return newFromXML(data, integrity)
}

// NewFromYAMLWithIntegrity создаёт персонажа из YAML с кастомным проверяющим целостности
func NewFromYAMLWithIntegrity(data []byte, integrity *entity.IntegrityChecker) (Person, error) {
	return newFromYAML(data, integrity)
}

// ---------------  Внутренняя реализация конструкторов -------------------------------

func newFromJSON(data []byte, integrity *entity.IntegrityChecker) (Person, error) {
	ser := NewSerializer(integrity)
	return ser.FromJSON(data)
}

func newFromXML(data []byte, integrity *entity.IntegrityChecker) (Person, error) {
	ser := NewSerializer(integrity)
	return ser.FromXML(data)
}

func newFromYAML(data []byte, integrity *entity.IntegrityChecker) (Person, error) {
	ser := NewSerializer(integrity)
	return ser.FromYAML(data)
}

// =============  Реализация сериализации/десириализации для person ====================

// personConverter — адаптер для существующих функций ToDTO/FromDTO
type personConverter struct{}

func (c personConverter) ToDTO(p Person) PersonDTO {
	return ToDTO(p) // существующая функция из dto.go
}

func (c personConverter) FromDTO(dto PersonDTO) (Person, error) {
	return FromDTO(dto) // существующая функция из dto.go
}

// NewSerializer — фабрика с типобезопасностью
func NewSerializer(integrity *entity.IntegrityChecker) *serializer.Serializer[Person, PersonDTO] {
	return serializer.New[Person, PersonDTO](personConverter{}, integrity)
}
