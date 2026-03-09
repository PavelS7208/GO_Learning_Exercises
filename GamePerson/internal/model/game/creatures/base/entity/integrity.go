package entity

import "fmt"

type IntegrityChecker struct{}

func NewIntegrityChecker() *IntegrityChecker {
	return &IntegrityChecker{}
}

// Check — метод БЕЗ параметров типа (дженерик не выведет тип)
func (ic *IntegrityChecker) Check(obj any, sourceFormat string) error {
	type validatable interface {
		Validate() error
	}

	if v, ok := obj.(validatable); ok {
		if err := v.Validate(); err != nil {
			return fmt.Errorf("%s deserialized object failed integrity check: %w", sourceFormat, err)
		}
	}
	return nil
}
