package msync

import "sync"

// RWPMutex — разделяемый мьютекс с приоритетом писателей.
//
// Поля:
//   mu       — защищает счётчик readers
//   writerMutex  — удерживается пока есть хотя бы один читатель или активный писатель;
//                не даёт писателю войти, пока читатели не уйдут
//   readerMutex — «турникет»: писатель захватывает его на всё время работы,
//                блокируя новых читателей (writer priority)
//   readers  — число активных читателей

type RWPMutex struct {
	mu          sync.Mutex
	writerMutex sync.Mutex
	readerMutex sync.Mutex
	readers     int
}

func NewRWPMutex() *RWPMutex {
	return &RWPMutex{}
}

// Lock захватывает блокировку (писатель).
// Блокирует новых читателей и ждёт завершения активных.
func (m *RWPMutex) Lock() {
	m.readerMutex.Lock()
	m.writerMutex.Lock()
}

// Unlock освобождает эксклюзивную блокировку (писатель).
func (m *RWPMutex) Unlock() {
	m.writerMutex.Unlock()
	m.readerMutex.Unlock()
}

// RLock захватывает блокировку (читатель).
// Несколько читателей могут удерживать её одновременно.
// Если есть ожидающий писатель — блокируется до его завершения (writer priority).
func (m *RWPMutex) RLock() {
	// Пройти через турникет: если писатель его держит — ждём.
	m.readerMutex.Lock()
	// разрешаем другим читателям работать
	m.readerMutex.Unlock()

	m.mu.Lock()
	m.readers++
	if m.readers == 1 { //  Как только появился читать, то он блокирует появление писателя
		m.writerMutex.Lock()
	}
	m.mu.Unlock()
}

// RUnlock освобождает блокировку (читатель).
func (m *RWPMutex) RUnlock() {
	m.mu.Lock()
	m.readers--
	if m.readers == 0 { // Последний читатель при уходе разрешает входить писателю
		m.writerMutex.Unlock()
	}
	m.mu.Unlock()
}

// TryLock пытается захватить эксклюзивную блокировку (писателя) без ожидания.
// Возвращает true если блокировка захвачена, false если мьютекс занят.
func (m *RWPMutex) TryLock() bool {
	// Попытка захватить «турникет» (readerMutex).
	// Если не удалось:
	//   либо другой писатель уже держит турникет (ждёт читателей или пишет),
	//   либо писатель пишет и залочил
	// В любом случае — немедленно отказываем, не блокируясь.
	if !m.readerMutex.TryLock() {
		return false
	}
	// Раз дошли сюда, значит без ожиданий захватили readerMutex (первую часть Lock выполнили)
	// (m.readerMutex - захвачен, не забыть его освободить при выходе c false)

	// Пробуем без ожидания захватить блокировку writerMutex.
	if !m.writerMutex.TryLock() {
		m.readerMutex.Unlock() // откатить захват турникета по комментарию выше
		return false
	}
	// Если сюда попали, значит и m.readerMutex и m.writerMutex (условие на Lock() ) удалось захватить без ожидания
	return true
}

// TryRLock пытается захватить разделяемую блокировку читателя без ожидания.
// Возвращает true если блокировка захвачена, false если активен или ожидает писатель.
func (m *RWPMutex) TryRLock() bool {
	// Проверяем турникет, если писатель его держит, то он закрыт
	if !m.readerMutex.TryLock() {
		return false
	}
	// Если сюда попали, значит писателя нет (ни активного ни ожидающего), но мы залочили readerMutex, открываем его
	m.readerMutex.Unlock()

	// Добавляем читателя (если не false условие будет) под защитой
	m.mu.Lock()
	defer m.mu.Unlock()
	// Если уже есть читатели (значит кто-то до нас writerMutex блокирнул),- то нам удалось без очередей (++ читателя)
	if m.readers > 0 {
		m.readers++
		return true
	}

	// Нужна доп. проверка на то что мы первый читатель, так как пока проверяли выше уже и писатель мог прийти (TOCTOU проблема)
	if !m.writerMutex.TryLock() { // писатель шустрый оказался - false
		return false
	}
	// Итоговый true. мы были первыми читателями
	m.readers++
	return true
}
