package serializer

import (
	"GamePerson/internal/model/game/creatures/base/entity"
	"encoding/json"
	"encoding/xml"
	"fmt"

	"gopkg.in/yaml.v3"
)

// Converter — минимальный интерфейс преобразования сущность ↔ DTO
type Converter[E any, D any] interface {
	ToDTO(entity E) D
	FromDTO(dto D) (E, error)
}

// Serializer — обобщённый сериализатор для ЛЮБОЙ сущности
type Serializer[E any, D any] struct {
	converter Converter[E, D]
	integrity *entity.IntegrityChecker
}

func New[E any, D any](converter Converter[E, D], integrity *entity.IntegrityChecker) *Serializer[E, D] {
	if integrity == nil {
		integrity = entity.NewIntegrityChecker()
	}
	return &Serializer[E, D]{converter: converter, integrity: integrity}
}

// Вспомогательная функция десериализации (без дублирования)
func (s *Serializer[E, D]) deserialize(
	data []byte,
	unmarshal func([]byte, interface{}) error,
	formatName string,
) (E, error) {
	var dto D
	if err := unmarshal(data, &dto); err != nil {
		var zero E
		return zero, fmt.Errorf("failed to unmarshal from %s: %w", formatName, err)
	}

	e, err := s.converter.FromDTO(dto)
	if err != nil {
		var zero E
		return zero, fmt.Errorf("%s deserialization failed: %w", formatName, err)
	}

	// Единая точка валидации для ВСЕХ сущностей
	if err := s.integrity.Check(e, formatName); err != nil {
		var zero E
		return zero, err
	}

	return e, nil
}

// Публичные методы — без дублирования логики

func (s *Serializer[E, D]) ToJSON(entity E) ([]byte, error) {
	dto := s.converter.ToDTO(entity)
	data, err := json.MarshalIndent(dto, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("failed to marshal to JSON: %w", err)
	}
	return data, nil
}

func (s *Serializer[E, D]) FromJSON(data []byte) (E, error) {
	return s.deserialize(data, json.Unmarshal, "JSON")
}

func (s *Serializer[E, D]) ToXML(entity E) ([]byte, error) {
	dto := s.converter.ToDTO(entity)
	data, err := xml.MarshalIndent(dto, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("failed to marshal to XML: %w", err)
	}
	return append([]byte(xml.Header), data...), nil
}

func (s *Serializer[E, D]) FromXML(data []byte) (E, error) {
	return s.deserialize(data, xml.Unmarshal, "XML")
}

func (s *Serializer[E, D]) ToYAML(entity E) ([]byte, error) {
	dto := s.converter.ToDTO(entity)
	data, err := yaml.Marshal(dto)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal to YAML: %w", err)
	}
	return data, nil
}

func (s *Serializer[E, D]) FromYAML(data []byte) (E, error) {
	return s.deserialize(data, yaml.Unmarshal, "YAML")
}
