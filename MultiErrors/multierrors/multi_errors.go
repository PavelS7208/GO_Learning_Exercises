package multierrors

import (
	"errors"
	"iter"
	"reflect"
	"strconv"
	"strings"
)

type Errors struct {
	errors []error
}

func New(errs ...error) *Errors {
	return Append(nil, errs...)
}

func (e *Errors) Error() string {

	sb := strings.Builder{}
	if len(e.errors) > 0 {
		sb.WriteString(strconv.Itoa(len(e.errors)))
		sb.WriteString(" errors occurred:\n")
		for _, err := range e.errors {
			sb.WriteString("\t* ")
			sb.WriteString(err.Error())
		}
		sb.WriteString("\n")
	}
	return sb.String()
}

func (e *Errors) IsEmpty() bool {
	return len(e.errors) == 0
}

func (e *Errors) Len() int {
	return len(e.errors)
}

// ErrorsSlice - обычный геттер
func (e *Errors) ErrorsSlice() []error {
	return append([]error(nil), e.errors...)
}

// All - Для новых клиентов на Go 1.23+
func (e *Errors) All() iter.Seq[error] {
	return func(yield func(error) bool) {
		for _, err := range e.errors {
			if err != nil && !yield(err) {
				return
			}
		}
	}
}

// Append - Сахар для удобства, метод структуры
func (e *Errors) Append(errs ...error) *Errors {
	return Append(e, errs...)
}

// Flatten - Сахар для удобства, метод структуры
func (e *Errors) Flatten() *Errors {
	return Flatten(e)
}

// FlattenWith - Сахар для удобства, если нужно добавить ещё ошибок при сглаживании
func (e *Errors) FlattenWith(errs ...error) *Errors {
	return Flatten(e, errs...)
}

func (e *Errors) Is(target error) bool {
	// Проверяем является ли проверяемый таргет нашим типом
	var targetErrors *Errors
	if errors.As(target, &targetErrors) {
		// Если да проверяем, содержит ли e все ошибки из target (в e может быть и больше, чем в target)
		return e.containsAll(*targetErrors)
	}
	//  Значит какой-то другой поддерживающий интерфейс ошибки
	//  Если в е - есть в списке target значит Is  истина
	for _, err := range e.errors {
		if err == nil {
			continue
		}
		if errors.Is(err, target) {
			return true
		}
	}
	return false
}

func (e *Errors) As(target interface{}) bool {
	if target == nil {
		return false
	}
	//  Должна быть проверка что target это указатель
	t := reflect.ValueOf(target)
	if t.Kind() != reflect.Ptr || t.IsNil() {
		return false
	}
	for _, err := range e.errors {
		//  Все верно, ругается предупреждением IDE - что ошибочно
		if errors.As(err, target) {
			return true
		}
	}
	return false
}

func (e *Errors) Unwrap() []error {
	return append([]error(nil), e.errors...)
}

// containsAll - Для каждой ошибки из target проверяем есть ли она в списке e
func (e *Errors) containsAll(target Errors) bool {
	for _, targetErr := range target.errors {
		found := false
		for _, err := range e.errors {
			if errors.Is(err, targetErr) {
				found = true
				break
			}
		}
		if !found {
			return false
		}
	}
	return true
}
