package multierrors

// Append добавляет ошибки в список без распаковки вложенных структур.
// Если передан *Errors, он вставляется как единый объект (сохраняется иерархия).
func Append(err error, errs ...error) *Errors {
	result := &Errors{}

	// Простая логика добавления: если не nil, то кладем в слайс
	add := func(e error) {
		if e != nil {
			result.errors = append(result.errors, e)
		}
	}

	add(err)
	for _, e := range errs {
		add(e)
	}

	if len(result.errors) == 0 {
		return nil
	}
	return result
}
