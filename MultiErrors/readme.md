Опишу попозже эту домашку до конца


Позволяет объединятьошибки в контейнер

```go
  err1 := errors.New("не удалось подключиться к БД")
	err2 := errors.New("таймаут кэша")
	err3 := errors.New("ошибка валидации email")

	// Собираем ошибки в контейнер (nil игнорируется)
	mErr := multierrors.Append(err1, err2, nil, err3)
  // Оборачиваем группу в контекст
  wrappedGroup := fmt.Errorf("registration failed: %w", mErr)
	if mErr != nil {
		fmt.Println(wrappedGroup)
	}
  // И к ней еще добавляем
	mErr = multierrors.Append(wrappedGroup, errors.New("ошибка еще одна"))
	if mErr != nil {
		fmt.Println(mErr)
	}
```
