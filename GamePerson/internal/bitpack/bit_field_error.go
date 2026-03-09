package bitpack

import (
	"fmt"
)

// ==================== Единый тип ошибки ===================================

type Error struct {
	Kind    ErrorKind
	Details errorDetails
}

type ErrorKind int

const (
	KindStartAfterEnd ErrorKind = iota
	KindEndOutOfRange
	KindValueOverflow
	KindValueUnderflow
	KindValueRangeInverted
	KindPositionOutOfRange
	KindSliceEmpty
	KindSliceTooLarge
)

type errorDetails struct {
	Start       BitPosition
	End         BitPosition
	Position    BitPosition
	Value       uint64
	AllowedMax  uint64
	BitWidth    uint8
	MinValue    int64
	MaxValue    int64
	SignedValue int64
	SliceLength int
}

func (e *Error) Error() string {
	switch e.Kind {
	case KindStartAfterEnd:
		return fmt.Sprintf("bit field range error: start position (%d) must be <= end position (%d)",
			e.Details.Start, e.Details.End)
	case KindEndOutOfRange:
		return fmt.Sprintf("bit field range error: end position (%d) must be < 64", e.Details.End)
	case KindValueOverflow:
		if e.Details.SignedValue != 0 {
			return fmt.Sprintf("value %d exceeds maximum %d for %d-bit field",
				e.Details.SignedValue, e.Details.MaxValue, e.Details.BitWidth)
		}
		return fmt.Sprintf("value %d exceeds capacity of %d-bit field (max allowed: %d)",
			e.Details.Value, e.Details.BitWidth, e.Details.AllowedMax)
	case KindValueUnderflow:
		return fmt.Sprintf("value %d is less than minimum %d for %d-bit field",
			e.Details.SignedValue, e.Details.MinValue, e.Details.BitWidth)
	case KindPositionOutOfRange:
		return fmt.Sprintf("bool bit field error: position (%d) must be < 64", e.Details.Position)
	case KindValueRangeInverted:
		return fmt.Sprintf("value range inverted: min (%d) > max (%d)",
			e.Details.MinValue, e.Details.MaxValue)
	case KindSliceEmpty:
		return "packed slice is empty"
	case KindSliceTooLarge:
		return fmt.Sprintf("packed slice too large: %d bytes (max 8)", e.Details.SliceLength)
	default:
		return "unknown bit field error"
	}
}

// Is позволяет использовать errors.Is() для проверки типа ошибки
func (e *Error) Is(target error) bool {
	t, ok := target.(*Error)
	if !ok {
		return false
	}
	return e.Kind == t.Kind
}

// Для использования с errors.Is()
var (
	ErrStartAfterEnd      = &Error{Kind: KindStartAfterEnd}
	ErrEndOutOfRange      = &Error{Kind: KindEndOutOfRange}
	ErrValueOverflow      = &Error{Kind: KindValueOverflow}
	ErrValueUnderflow     = &Error{Kind: KindValueUnderflow}
	ErrValueRangeInverted = &Error{Kind: KindValueRangeInverted}
	ErrPositionOutOfRange = &Error{Kind: KindPositionOutOfRange}
	ErrSliceEmpty         = &Error{Kind: KindSliceEmpty}
	ErrSliceTooLarge      = &Error{Kind: KindSliceTooLarge}
)

// Вспомогательные конструкторы
func newStartAfterEndError(start, end BitPosition) error {
	return &Error{
		Kind:    KindStartAfterEnd,
		Details: errorDetails{Start: start, End: end},
	}
}

func newEndOutOfRangeError(end BitPosition) error {
	return &Error{
		Kind:    KindEndOutOfRange,
		Details: errorDetails{End: end},
	}
}

func newValueOverflowError(value, allowedMax uint64, width uint8) error {
	return &Error{
		Kind: KindValueOverflow,
		Details: errorDetails{
			Value: value, AllowedMax: allowedMax, BitWidth: width,
		},
	}
}

func newPositionOutOfRangeError(pos BitPosition) error {
	return &Error{
		Kind:    KindPositionOutOfRange,
		Details: errorDetails{Position: pos},
	}
}

func newValueOutOfRangeError(valueMin, valueMax, allowedMin, allowedMax int64, width uint8) error {
	kind := KindValueOverflow
	if valueMin < allowedMin {
		kind = KindValueUnderflow
	}

	return &Error{
		Kind: kind,
		Details: errorDetails{
			SignedValue: valueMin,
			MinValue:    allowedMin,
			MaxValue:    allowedMax,
			BitWidth:    width,
		},
	}
}

func newValueRangeInvertedError(min, max int64) error {
	return &Error{
		Kind: KindValueRangeInverted,
		Details: errorDetails{
			MinValue: min,
			MaxValue: max,
		},
	}
}

func newSliceEmptyError() error {
	return &Error{
		Kind: KindSliceEmpty,
	}
}

func newSliceTooLargeError(length int) error {
	return &Error{
		Kind:    KindSliceTooLarge,
		Details: errorDetails{SliceLength: length},
	}
}
