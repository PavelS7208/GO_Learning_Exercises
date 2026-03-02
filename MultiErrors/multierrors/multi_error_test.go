package multierrors

import (
	"errors"
	"fmt"
	"testing"
)

// Вспомогательные ошибки для тестов
var (
	errFoo = errors.New("foo error")
	errBar = errors.New("bar error")
	errBaz = errors.New("baz error")
)

// Кастомная ошибка для тестов As/Is
type customError struct {
	msg string
}

func (e *customError) Error() string { return e.msg }

// Ошибка с обёрткой для тестов %w
type wrappedError struct {
	err error
}

func (e *wrappedError) Error() string { return "wrapped: " + e.err.Error() }
func (e *wrappedError) Unwrap() error { return e.err }

// ============================================================================
// БАЗОВЫЕ ТЕСТЫ: конструкторы, Error(), Unwrap(), утилиты
// ============================================================================

func TestNew(t *testing.T) {
	t.Run("empty returns nil", func(t *testing.T) {
		if got := New(); got != nil {
			t.Errorf("New() = %v, want nil", got)
		}
	})

	t.Run("single error", func(t *testing.T) {
		got := New(errFoo)
		if got == nil || got.Len() != 1 || !errors.Is(got, errFoo) {
			t.Errorf("New(errFoo) failed: len=%d, Is(errFoo)=%v", got.Len(), errors.Is(got, errFoo))
		}
	})

	t.Run("multiple errors", func(t *testing.T) {
		got := New(errFoo, errBar, errBaz)
		if got.Len() != 3 {
			t.Errorf("Len() = %d, want 3", got.Len())
		}
	})

	t.Run("with nil errors filtered", func(t *testing.T) {
		got := New(nil, errFoo, nil, errBar)
		if got.Len() != 2 {
			t.Errorf("Len() = %d, want 2 (nils filtered)", got.Len())
		}
	})
}

func TestErrors_Error(t *testing.T) {
	tests := []struct {
		name     string
		errors   []error
		expected string
	}{
		{
			name:     "empty",
			errors:   []error{},
			expected: "",
		},
		{
			name:     "single error",
			errors:   []error{errFoo},
			expected: "1 errors occurred:\n\t* foo error\n",
		},
		{
			name:     "multiple errors",
			errors:   []error{errFoo, errBar},
			expected: "2 errors occurred:\n\t* foo error\t* bar error\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &Errors{errors: tt.errors}
			if got := e.Error(); got != tt.expected {
				t.Errorf("Error() = %q, want %q", got, tt.expected)
			}
		})
	}
}

func TestErrors_Unwrap(t *testing.T) {
	e := New(errFoo, errBar)
	unwrapped := e.Unwrap()

	if len(unwrapped) != 2 {
		t.Errorf("Unwrap() len = %d, want 2", len(unwrapped))
	}
	// Проверка, что возвращается копия, а не ссылка на внутреннее поле
	unwrapped[0] = errBaz
	if e.errors[0] == errBaz { // Сравнение намеренное: проверяем что это одно и тоже
		t.Error("Unwrap() returned reference to internal slice, not a copy")
	}
}

func TestErrorsSlice_CopySemantics(t *testing.T) {
	e := New(errFoo, errBar)
	slice := e.ErrorsSlice()

	slice[0] = errBaz
	if e.errors[0] == errBaz { // Сравнение намеренное: проверяем что это одно и тоже
		t.Error("ErrorsSlice() returned reference to internal slice")
	}
}

func TestIsEmpty_And_Len(t *testing.T) {
	t.Run("empty is empty", func(t *testing.T) {
		e := &Errors{}
		if !e.IsEmpty() || e.Len() != 0 {
			t.Error("Empty Errors should report IsEmpty=true, Len=0")
		}
	})

	t.Run("non-empty", func(t *testing.T) {
		e := New(errFoo)
		if e.IsEmpty() || e.Len() != 1 {
			t.Error("Non-empty Errors should report IsEmpty=false, Len=1")
		}
	})
}

// ============================================================================
// ТЕСТЫ ИТЕРАТОРА All() (Go 1.23+)
// ============================================================================

func TestErrors_All(t *testing.T) {
	e := New(errFoo, errBar, errBaz)

	t.Run("iterates all errors", func(t *testing.T) {
		var collected []error
		for err := range e.All() {
			collected = append(collected, err)
		}
		if len(collected) != 3 {
			t.Errorf("Collected %d errors, want 3", len(collected))
		}
	})

	t.Run("early termination with yield", func(t *testing.T) {
		count := 0
		for range e.All() {
			count++
			if count == 2 {
				break // эмуляция раннего выхода
			}
		}
		// Просто проверяем, что итерация работает корректно
		if count != 2 {
			t.Errorf("Early termination failed: count=%d", count)
		}
	})

	t.Run("empty errors", func(t *testing.T) {
		empty := &Errors{}
		count := 0
		for range empty.All() {
			count++
		}
		if count != 0 {
			t.Error("Empty Errors should yield nothing")
		}
	})
}

// ============================================================================
// ТЕСТЫ Is() — сложные сценарии
// ============================================================================

func TestErrors_Is(t *testing.T) {
	t.Run("direct match", func(t *testing.T) {
		e := New(errFoo, errBar)
		if !errors.Is(e, errFoo) || !errors.Is(e, errBar) {
			t.Error("Is() should find direct errors")
		}
		if errors.Is(e, errBaz) {
			t.Error("Is() should not find unrelated error")
		}
	})

	t.Run("wrapped error match", func(t *testing.T) {
		wrapped := fmt.Errorf("context: %w", errFoo)
		e := New(wrapped, errBar)
		if !errors.Is(e, errFoo) {
			t.Error("Is() should unwrap and find errFoo")
		}
	})

	t.Run("nested *Errors via Is", func(t *testing.T) {
		inner := New(errFoo)
		outer := New(inner, errBar)
		if !errors.Is(outer, errFoo) {
			t.Error("Is() should find errors in nested *Errors")
		}
	})

	t.Run("Is with target *Errors — subset check", func(t *testing.T) {
		// e содержит [foo, bar, baz], target — [foo, bar]
		e := New(errFoo, errBar, errBaz)
		target := New(errFoo, errBar)
		if !errors.Is(e, target) {
			t.Error("Is() should return true when e contains all errors from target")
		}

		// Обратный случай: target содержит больше, чем e
		targetMore := New(errFoo, errBar, errBaz, errors.New("extra"))
		if errors.Is(e, targetMore) {
			t.Error("Is() should return false when target has errors not in e")
		}
	})

	t.Run("nil errors skipped in Is", func(t *testing.T) {
		e := &Errors{errors: []error{nil, errFoo}}
		if !errors.Is(e, errFoo) {
			t.Error("Is() should skip nil and find errFoo")
		}
	})
}

// ============================================================================
// ТЕСТЫ As() — работа с типами через reflection
// ============================================================================

func TestErrors_As(t *testing.T) {
	t.Run("finds custom error type", func(t *testing.T) {
		custom := &customError{msg: "custom"}
		e := New(errFoo, custom, errBar)

		var got *customError
		if !errors.As(e, &got) {
			t.Error("As() should find customError")
		}
		if got != custom { // Сравнение указателей намеренное: проверяем, что As() вернул тот же экземпляр
			t.Errorf("As() got %v, want %v", got, custom)
		}
	})

	t.Run("wrapped custom error", func(t *testing.T) {
		custom := &customError{msg: "wrapped custom"}
		wrapped := fmt.Errorf("wrap: %w", custom)
		e := New(wrapped)

		var got *customError
		if !errors.As(e, &got) {
			t.Error("As() should unwrap and find customError")
		}
	})

	t.Run("nil target returns false", func(t *testing.T) {
		e := New(errFoo)
		if e.As(nil) {
			t.Error("As(nil) should return false")
		}
	})

	t.Run("non-pointer target returns false", func(t *testing.T) {
		e := New(errFoo)
		var got customError // не указатель
		if e.As(got) {
			t.Error("As() with non-pointer should return false")
		}
	})
}

// ============================================================================
// ТЕСТЫ Append() — добавление БЕЗ распаковки *Errors
// ============================================================================

func TestAppend_NoFlatten(t *testing.T) {
	t.Run("appends simple errors", func(t *testing.T) {
		e := Append(nil, errFoo, errBar)
		if e.Len() != 2 {
			t.Errorf("Len() = %d, want 2", e.Len())
		}
	})

	t.Run("preserves *Errors as single element", func(t *testing.T) {
		inner := New(errFoo, errBar)
		e := Append(nil, inner, errBaz)

		// inner должен быть одним элементом, а не распакован
		if e.Len() != 2 {
			t.Errorf("Append should preserve *Errors as single item: got len=%d", e.Len())
		}
		// Проверяем, что первый элемент — это именно *Errors, а не errFoo
		if !errors.Is(e, inner) {
			t.Error("Append should preserve *Errors identity")
		}
	})

	t.Run("nil input returns nil", func(t *testing.T) {
		if got := Append(nil); got != nil {
			t.Errorf("Append(nil) = %v, want nil", got)
		}
	})

	t.Run("method Append on existing *Errors", func(t *testing.T) {
		e := New(errFoo)
		e = e.Append(errBar, errBaz)
		if e.Len() != 3 {
			t.Errorf("Method Append failed: len=%d", e.Len())
		}
	})
}

// ============================================================================
// ТЕСТЫ Flatten() — рекурсивное сглаживание
// ============================================================================

func TestFlatten_Basic(t *testing.T) {
	t.Run("flattens nested *Errors", func(t *testing.T) {
		inner := New(errFoo, errBar)
		e := Flatten(inner, errBaz)

		// Ожидаем 3 плоских ошибки, а не 2 элемента
		if e.Len() != 3 {
			t.Errorf("Flatten should produce 3 errors, got %d", e.Len())
		}
		// Все ошибки должны быть найдены через Is
		if !errors.Is(e, errFoo) || !errors.Is(e, errBar) || !errors.Is(e, errBaz) {
			t.Error("Flatten should preserve all errors for Is() checks")
		}
	})

	t.Run("deeply nested *Errors", func(t *testing.T) {
		level3 := New(errFoo)
		level2 := New(level3, errBar)
		level1 := New(level2, errBaz)

		flat := Flatten(level1)
		if flat.Len() != 3 {
			t.Errorf("Deep flatten: got %d errors, want 3", flat.Len())
		}
	})

	t.Run("method Flatten()", func(t *testing.T) {
		inner := New(errFoo, errBar)
		e := New(inner, errBaz)
		flat := e.Flatten()

		if flat.Len() != 3 {
			t.Errorf("Method Flatten: got %d, want 3", flat.Len())
		}
	})
}

func TestFlatten_WithWrappedErrors(t *testing.T) {
	t.Run("unwraps %w wrapped *Errors", func(t *testing.T) {
		inner := New(errFoo)
		wrapped := fmt.Errorf("context: %w", inner)

		flat := Flatten(wrapped, errBar)
		if flat.Len() != 2 {
			t.Errorf("Should unwrap *Errors from key w wrapper: got %d", flat.Len())
		}
		if !errors.Is(flat, errFoo) || !errors.Is(flat, errBar) {
			t.Error("Flatten should find errors through %w wrapper")
		}
	})

	t.Run("handles mixed wrapped and plain", func(t *testing.T) {
		inner := New(errFoo)
		wrappedInner := fmt.Errorf("wrap1: %w", inner)
		plain := errBar
		wrappedPlain := fmt.Errorf("wrap2: %w", errBaz)

		flat := Flatten(wrappedInner, plain, wrappedPlain)
		if flat.Len() != 3 {
			t.Errorf("Mixed flatten: got %d, want 3", flat.Len())
		}
	})
}

func TestFlatten_NilHandling(t *testing.T) {
	t.Run("nil input returns nil", func(t *testing.T) {
		if got := Flatten(nil); got != nil {
			t.Errorf("Flatten(nil) = %v, want nil", got)
		}
	})

	t.Run("filters nils during flatten", func(t *testing.T) {
		inner := New(nil, errFoo, nil)
		flat := Flatten(inner, nil, errBar)
		if flat.Len() != 2 {
			t.Errorf("Should filter nils: got %d, want 2", flat.Len())
		}
	})
}

// ============================================================================
// ТЕСТЫ FlattenWith()
// ============================================================================

func TestFlattenWith(t *testing.T) {
	t.Run("flatten then add more", func(t *testing.T) {
		inner := New(errFoo, errBar)
		e := New(inner)
		flat := e.FlattenWith(errBaz)

		if flat.Len() != 3 {
			t.Errorf("FlattenWith: got %d, want 3", flat.Len())
		}
	})

	t.Run("chained FlattenWith", func(t *testing.T) {
		e := New(errFoo)
		e = e.FlattenWith(errBar).FlattenWith(errBaz)
		if e.Len() != 3 {
			t.Errorf("Chained FlattenWith: got %d, want 3", e.Len())
		}
	})
}

// ============================================================================
// КРИТИЧЕСКИЕ ТЕСТЫ: Append vs Flatten — идентичность выводов
// ============================================================================

func TestAppendVsFlatten_OutputIdentity(t *testing.T) {
	// Сценарий: создаём сложную вложенную структуру,
	// затем сравниваем Error() вывод после Append+Flatten vs прямого Flatten

	setupComplex := func() (error, []error) {
		// Уровень 3: базовые ошибки
		level3 := New(errFoo, errBar)
		// Уровень 2: обёртка + ошибки
		level2 := New(fmt.Errorf("wrapL2: %w", level3), errBaz)
		// Уровень 1: ещё одна обёртка
		level1 := fmt.Errorf("wrapL1: %w", level2)
		extra := []error{errors.New("extra1"), errors.New("extra2")}
		return level1, extra
	}

	t.Run("Error() output identical: Append+Flatten vs direct Flatten", func(t *testing.T) {
		root, extras := setupComplex()

		// Вариант 1: Append затем Flatten
		v1 := Append(root, extras...)
		if v1 != nil {
			v1 = v1.Flatten()
		}

		// Вариант 2: сразу Flatten
		all := append([]error{root}, extras...)
		v2 := Flatten(all[0], all[1:]...)

		// Сравниваем строковое представление
		str1, str2 := "", ""
		if v1 != nil {
			str1 = v1.Error()
		}
		if v2 != nil {
			str2 = v2.Error()
		}

		if str1 != str2 {
			t.Errorf("Error() outputs differ:\nAppend+Flatten: %q\nDirect Flatten: %q", str1, str2)
		}

		// Сравниваем количество ошибок
		len1, len2 := 0, 0
		if v1 != nil {
			len1 = v1.Len()
		}
		if v2 != nil {
			len2 = v2.Len()
		}
		if len1 != len2 {
			t.Errorf("Len mismatch: Append+Flatten=%d, Direct Flatten=%d", len1, len2)
		}
	})

	t.Run("Is() behavior identical after both approaches", func(t *testing.T) {
		root, extras := setupComplex()

		v1 := Append(root, extras...)
		if v1 != nil {
			v1 = v1.Flatten()
		}
		v2 := Flatten(root, extras...)

		testErrors := []error{errFoo, errBar, errBaz, errors.New("extra1"), errors.New("extra2"), errors.New("nonexistent")}
		for _, target := range testErrors {
			r1, r2 := false, false
			if v1 != nil {
				r1 = errors.Is(v1, target)
			}
			if v2 != nil {
				r2 = errors.Is(v2, target)
			}
			if r1 != r2 {
				t.Errorf("Is(%q) mismatch: Append+Flatten=%v, Direct Flatten=%v", target.Error(), r1, r2)
			}
		}
	})
}

// ============================================================================
// EDGE CASES: рекурсия, циклы, производительность
// ============================================================================

func TestFlatten_NoInfiniteRecursion(t *testing.T) {
	// Хотя в реальной практике циклические ссылки маловероятны,
	// проверяем, что наша рекурсия безопасна для глубоких структур

	t.Run("very deep nesting", func(t *testing.T) {
		current := New(errFoo)
		for i := 0; i < 100; i++ {
			current = New(current, errBar)
		}
		flat := Flatten(current)

		// ✅ Исправлено: 1 errFoo + 100 errBar = 101 ошибка
		if flat.Len() != 101 {
			t.Errorf("Deep nesting: got %d errors, want 101", flat.Len())
		}

		// Дополнительно: проверяем, что все ошибки на месте
		if !errors.Is(flat, errFoo) {
			t.Error("Flatten should preserve errFoo")
		}

		// Считаем количество errBar для уверенности
		barCount := 0
		for _, err := range flat.ErrorsSlice() {
			if errors.Is(err, errBar) {
				barCount++
			}
		}
		if barCount != 100 {
			t.Errorf("Expected 100 errBar instances, got %d", barCount)
		}
	})
}

func TestConcurrentSafety_Readonly(t *testing.T) {
	// Тест на то, что методы-геттеры не мутируют состояние
	e := New(errFoo, errBar, errBaz)

	t.Run("Unwrap and ErrorsSlice can be called concurrently", func(t *testing.T) {
		done := make(chan bool, 2)
		go func() {
			for i := 0; i < 100; i++ {
				_ = e.Unwrap()
			}
			done <- true
		}()
		go func() {
			for i := 0; i < 100; i++ {
				_ = e.ErrorsSlice()
			}
			done <- true
		}()
		<-done
		<-done
		// Если паники не произошло — тест пройден
	})
}

// ============================================================================
// ТЕСТЫ containsAll (внутренний метод, тестируем через Is)
// ============================================================================

func TestContainsAll_ViaIs(t *testing.T) {
	t.Run("subset returns true", func(t *testing.T) {
		superset := New(errFoo, errBar, errBaz)
		subset := New(errFoo, errBar)
		if !errors.Is(superset, subset) {
			t.Error("Is() should return true for subset check")
		}
	})

	t.Run("superset as target returns false", func(t *testing.T) {
		small := New(errFoo)
		large := New(errFoo, errBar)
		if errors.Is(small, large) {
			t.Error("Is() should return false when target has more errors")
		}
	})

	t.Run("empty target always true", func(t *testing.T) {
		e := New(errFoo)
		empty := &Errors{}
		if !errors.Is(e, empty) {
			t.Error("Is() should return true when target is empty (vacuous truth)")
		}
	})
}

// ============================================================================
// БЕНЧМАРКИ
// ============================================================================

func BenchmarkErrors_Error(b *testing.B) {
	e := New(errFoo, errBar, errBaz, errors.New("test1"), errors.New("test2"))
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = e.Error()
	}
}

func BenchmarkFlatten_Deep(b *testing.B) {
	// Создаём глубокую структуру один раз
	setup := func() *Errors {
		current := New(errFoo)
		for i := 0; i < 50; i++ {
			current = New(current, errBar)
		}
		return current
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = Flatten(setup())
	}
}

func BenchmarkIs_LargeSet(b *testing.B) {
	// Создаём набор из 100 ошибок
	errs := make([]error, 100)
	for i := range errs {
		errs[i] = fmt.Errorf("error %d", i)
	}
	e := New(errs...)
	target := errs[50]

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = errors.Is(e, target)
	}
}

func BenchmarkAll_Iterator(b *testing.B) {
	errs := make([]error, 50)
	for i := range errs {
		errs[i] = fmt.Errorf("error %d", i)
	}
	e := New(errs...)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		count := 0
		for range e.All() {
			count++
		}
		_ = count
	}
}
