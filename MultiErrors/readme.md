## MultiErrors

За основу взято домашнее задание с курса продвинутый GO. Немного расширил ТЗ к заданию.


### Append/New
Позволяет объединять ошибки в контейнер

```go
err1 := errors.New("не удалось подключиться к БД")
err2 := errors.New("таймаут кэша")
err3 := errors.New("ошибка валидации email")

// Собираем ошибки в контейнер (nil игнорируется)
mErr := multierrors.Append(err1, err2, nil, err3)
// Оборачиваем группу в контекст
wrappedGroup := fmt.Errorf("registration failed: %w", mErr)
fmt.Println(wrappedGroup)
// И к ней еще добавляем что-то
mErr = multierrors.Append(wrappedGroup, errors.New("ошибка еще одна"))
fmt.Println(mErr)
```

### Flatten
Позволяет перевести вложенную структуру ошибок в плоскую


**После Append:**
```
Errors
├─ err1
├─ Errors
│  ├─ err2
│  └─ err3
├─ err4
└─ Errors
   └─ err5
```

**После `Flatten()`:**
```
Errors
├─ err1
├─ err2
├─ err3
├─ err4
└─ err5
```

```go
	validationErrs := multierrors.New(
		errors.New("empty name"),
		errors.New("invalid age"),
	)
	
	// Оборачиваем группу в контекст для примера
	wrappedGroup := fmt.Errorf("registration failed: %w", validationErrs)

	all1 := multierrors.Append(wrappedGroup, errors.New("email taken"))
	fmt.Println("Печатаем созданную ошибку с иерархией внутри:")
	fmt.Println(all1)

	// а теперь делаем тоже самое, но в плоскую структуру
	all2 := multierrors.Flatten(wrappedGroup, errors.New("email taken"))
	fmt.Println("Печатаем созданную ошибку в виде плоского списка:")
	fmt.Println(all2)

	// Или даже так
	fmt.Println("Печатаем созданную ошибку после преобразования в плоский список:")
	fmt.Println(all1.Flatten())
```

### Is / As

Поддерживаются операции с семантикой Is / As с поиском по всему списку ошибок в иерархической структуре
```go
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
```

```go
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
```


### Тесты

Часть тестов написано в ручную, часть сгенерированы по промптам. Закрыта большая часть функционала




