package multierrors

import "errors"

// Flatten — сервисный метод для плоского объединения ошибок.
// В отличие от Append, он рекурсивно раскрывает вложенные *Errors оставляя error.
func Flatten(err error, errs ...error) *Errors {
	result := &Errors{}

	// Если на входе (err) *Errors, то разворачиваем все что в нем есть в линейный список []error (рекурсивно) и этот список добавляем к result.errors

	// Если на входе (err) обертка (Wrap %w) над *Errors, то делаем unwrap и аналогично как выше

	// Если на входе (err) сам error - то  result.errors - append err

	// Бежим по errs c аналогичным алгоритмом: если просто error - аппенд, если *Errors (или обертка) сначала в плоский, потом аппенд

	// возращаем результат

	// Рекурсивная функция для сбора ошибок в плоский список
	var collect func(error)
	collect = func(e error) {
		if e == nil {
			return
		}

		// Прямая проверка нашего типа — работает всегда, без зависимости от Go-версии
		if multi, ok := e.(*Errors); ok {
			for _, inner := range multi.errors {
				collect(inner)
			}
			return
		}

		// Для сторонних ошибок с Unwrap() — оставляем errors.As
		var nested *Errors
		if errors.As(e, &nested) {
			for _, inner := range nested.errors {
				collect(inner)
			}
			return
		}
		// Базовый error
		result.errors = append(result.errors, e)
	}

	// Обрабатываем первый аргумент (err)
	collect(err)
	// Бежим по аргументам (errs) с тем же алгоритмом
	for _, e := range errs {
		collect(e)
	}
	// Возвращаем nil, если ошибок не осталось, иначе — результат
	if len(result.errors) == 0 {
		return nil
	}
	return result
}
