package main

import (
	"LessonsTask/Errors/MultiErrors/multierrors"
	"context"
	"errors"
	"fmt"
	"os"
)

func main() {

	exampleAppend()
	fmt.Println()

	exampleIs()
	fmt.Println()

	exampleAs()
	fmt.Println()

	exampleFlatten()
	fmt.Println()

	// Пример с итераторами (только для Go где оно уже есть)
	exampleIter()
}

// ---------- Пример Создание и вывод ошибок -----------------------
func exampleAppend() {

	err1 := errors.New("не удалось подключиться к БД")
	err2 := errors.New("таймаут кэша")
	err3 := errors.New("ошибка валидации email")

	// Собираем ошибки в контейнер (nil игнорируется)
	mErr := multierrors.Append(err1, err2, nil, err3)
	wrappedGroup := fmt.Errorf("registration failed: %w", mErr)
	if mErr != nil {
		fmt.Println(wrappedGroup)
	}

	mErr = multierrors.Append(wrappedGroup, errors.New("ошибка еще одна"))
	// Оборачиваем группу в контекст
	if mErr != nil {
		fmt.Println(mErr)
	}
}

// ------------ Пример Проверка ошибок через errors.Is() --------------
func exampleIs() {

	var (
		ErrNotFound = errors.New("resource not found")
		ErrTimeout  = errors.New("connection timeout")
	)

	// Эмулируем ошибку с контекстом (делаем враппер %w)
	wrapped := fmt.Errorf("check user: %w", ErrNotFound)

	// Собираем нашу мульти-ошибку
	mErr := multierrors.Append(wrapped, ErrTimeout)

	// errors.Is автоматически "заглянет" внутрь обёртки wrapped
	if errors.Is(mErr, ErrNotFound) {
		fmt.Println("Нашли ошибку - Ресурс не найден")
	}

	// Проверяем ошибку, которой точно нет
	if errors.Is(mErr, context.Canceled) {
		fmt.Println("Нашли ошибку - Контекст отменен")
	} else {
		fmt.Println("Не нашли ошибку - Контекст отменен")
	}
}

// ------------- Пример Извлечение типа ошибки через errors.As() ----
func exampleAs() {

	// Создаем список с разными ошибками
	mErr := multierrors.Append(
		errors.New("system error"),
		// Пример что можно использовать любой тип
		&ValidationError{Field: "email"},
	)

	// Пытаемся найти ValidationError
	var ve *ValidationError
	if errors.As(mErr, &ve) {
		fmt.Printf("Найдена ошибка валидации поля: %s\n", ve.Field)
	}

	// Можно проверить и отсутствие нужного типа
	var ce *customError
	if !errors.As(mErr, &ce) {
		fmt.Println("Ошибка customError не найдена")
	}
}

// ------- Пример Flatten. Раскручиваем вложенности в плоский список ошибок -----
func exampleFlatten() {

	fmt.Println("Пример с flatten")
	// Группа ошибок
	validationErrs := multierrors.New(
		errors.New("empty name"),
		errors.New("invalid age"),
	)
	if validationErrs == nil {
		fmt.Println("Ошибка при создании validationErrs")
		os.Exit(1)
	}
	// Оборачиваем группу в контекст
	wrappedGroup := fmt.Errorf("registration failed: %w", validationErrs)

	all1 := multierrors.Append(wrappedGroup, errors.New("email taken"))
	if all1 == nil {
		fmt.Println("Ошибка создания all1")
		os.Exit(1)
	}
	fmt.Println("Печатаем созданную ошибку с иерархией внутри:")
	fmt.Println(all1)

	// а теперь тоже самое, но в плоскую структуру
	all2 := multierrors.Flatten(wrappedGroup, errors.New("email taken"))
	if all2 == nil {
		fmt.Println("Ошибка создания all2")
		os.Exit(1)
	}
	fmt.Println("Печатаем созданную ошибку в виде плоского списка:")
	fmt.Println(all2)

	// Или даже так
	fmt.Println("Печатаем созданную ошибку после преобразования в плоский список:")
	fmt.Println(all1.Flatten())
}

// ----- Пример Итерации ----
func exampleIter() {

	mErr := multierrors.New(
		errors.New("step 1 failed"),
		errors.New("step 2 failed"),
	)
	if mErr == nil {
		fmt.Println("Ошибка при создании mErr")
		os.Exit(1)
	}
	fmt.Println(mErr)
	// Итерируемся по списку ошибок
	for err := range mErr.All() {
		fmt.Printf(" %v\n", err)
	}
}

// Для примеров кое что нужно вне main

type ValidationError struct {
	Field string
}

func (e *ValidationError) Error() string {
	return "validation failed: " + e.Field
}

type customError struct {
	msg string
}

func (e customError) Error() string {
	return e.msg
}
