package msync

import (
	"sync/atomic"
	"testing"
	"time"
)

func assertTrue(t *testing.T, val bool, msg string) {
	t.Helper()
	if !val {
		t.Errorf("expected true: %s", msg)
	}
}

func assertFalse(t *testing.T, val bool, msg string) {
	t.Helper()
	if val {
		t.Errorf("expected false: %s", msg)
	}
}

// TryLock на свободном мьютексе — должен захватить.
func TestTryLockFree(t *testing.T) {
	var m RWPMutex
	assertTrue(t, m.TryLock(), "TryLock on free mutex must succeed")
	m.Unlock()
}

// TryLock пока активен писатель — должен вернуть false.
func TestTryLockWhileWriter(t *testing.T) {
	var m RWPMutex
	m.Lock()
	assertFalse(t, m.TryLock(), "TryLock while writer active must fail")
	m.Unlock()
}

// TryLock пока есть активные читатели — должен вернуть false.
func TestTryLockWhileReaders(t *testing.T) {
	var m RWPMutex
	m.RLock()
	assertFalse(t, m.TryLock(), "TryLock while readers active must fail")
	m.RUnlock()
}

// TryLock пока другой писатель ожидает (держит readerMu) — должен вернуть false.
func TestTryLockWhileWriterWaiting(t *testing.T) {
	var m RWPMutex
	m.RLock() // читатель внутри

	var writerWaiting atomic.Bool
	go func() {
		writerWaiting.Store(true)
		m.Lock() // писатель занимает readerMu и ждёт читателя
		m.Unlock()
	}()

	// Дать горутине занять readerMu
	time.Sleep(50 * time.Millisecond)
	assertTrue(t, writerWaiting.Load(), "writer must be waiting")

	assertFalse(t, m.TryLock(), "TryLock while writer waiting must fail")
	m.RUnlock()
}

// TryRLock на свободном мьютексе — должен захватить.
func TestTryRLockFree(t *testing.T) {
	var m RWPMutex
	assertTrue(t, m.TryRLock(), "TryRLock on free mutex must succeed")
	m.RUnlock()
}

// TryRLock пока есть другие читатели — должен захватить (разделяемый доступ).
func TestTryRLockWithReaders(t *testing.T) {
	var m RWPMutex
	m.RLock()
	assertTrue(t, m.TryRLock(), "TryRLock while readers active must succeed")
	m.RUnlock()
	m.RUnlock()
}

// TryRLock пока активен писатель — должен вернуть false.
func TestTryRLockWhileWriter(t *testing.T) {
	var m RWPMutex
	m.Lock()
	assertFalse(t, m.TryRLock(), "TryRLock while writer active must fail")
	m.Unlock()
}

// TryRLock пока писатель ожидает — должен вернуть false (writer priority).
func TestTryRLockWhileWriterWaiting(t *testing.T) {
	var m RWPMutex
	m.RLock() // читатель внутри

	var writerWaiting atomic.Bool
	go func() {
		writerWaiting.Store(true)
		m.Lock() // писатель занимает readerMu и ждёт читателя
		m.Unlock()
	}()

	time.Sleep(50 * time.Millisecond)
	assertTrue(t, writerWaiting.Load(), "writer must be waiting")

	// Writer priority: новый читатель не должен обогнать ожидающего писателя.
	assertFalse(t, m.TryRLock(), "TryRLock while writer waiting must fail (writer priority)")
	m.RUnlock()
}

// После TryLock → Unlock мьютекс снова свободен.
func TestTryLockRelease(t *testing.T) {
	var m RWPMutex
	assertTrue(t, m.TryLock(), "first TryLock must succeed")
	m.Unlock()
	assertTrue(t, m.TryRLock(), "TryRLock after Unlock must succeed")
	m.RUnlock()
}
