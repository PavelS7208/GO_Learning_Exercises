package entity

import (
	"GamePerson/internal/model/config"
	"fmt"
	"strings"
	"unicode"
)

// ValidateAndCopyName валидирует имя и копирует его в буфер.
// Возвращает валидированное имя и ошибку
func ValidateAndCopyName(dst []byte, name string) (string, error) {
	name = strings.TrimSpace(name)
	if name == "" {
		return "", fmt.Errorf("name cannot be empty")
	}

	// КРИТИЧНО: обрезаем СРАЗУ, до валидации
	maxLen := len(dst)
	if len(name) > maxLen {
		name = name[:maxLen]
	}

	// Теперь валидируем только обрезанную часть
	for _, r := range name {
		if r > unicode.MaxASCII || !isValidNameChar(r) {
			return "", fmt.Errorf(
				"name must contain only ASCII letters (A-Z, a-z), digits, space, underscore or dash; got %q",
				r,
			)
		}
	}

	// Очищаем буфер и копируем
	for i := range dst {
		dst[i] = 0
	}
	copy(dst, name)

	return name, nil
}

func isValidNameChar(r rune) bool {
	return (r >= 'A' && r <= 'Z') ||
		(r >= 'a' && r <= 'z') ||
		(r >= '0' && r <= '9') ||
		r == ' ' || r == '_' || r == '-'
}

// ValidateCoordinate проверяет координату на выход за глобальные границы
func ValidateCoordinate(axis string, value int32) error {
	if value < config.MinCoord || value > config.MaxCoord {
		return fmt.Errorf("%s coordinate %d out of range [%d, %d]",
			axis, value, config.MinCoord, config.MaxCoord)
	}
	return nil
}
