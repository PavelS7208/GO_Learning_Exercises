package msync

import (
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestRWPMutexWithWriter(t *testing.T) {
	var mutex RWPMutex
	mutex.Lock() // writer

	var mutualExclusionWithWriter atomic.Bool
	mutualExclusionWithWriter.Store(true)
	var mutualExclusionWithReader atomic.Bool
	mutualExclusionWithReader.Store(true)

	go func() {
		mutex.Lock() // another writer
		mutualExclusionWithWriter.Store(false)
	}()

	go func() {
		mutex.RLock() // another reader
		mutualExclusionWithReader.Store(false)
	}()

	time.Sleep(time.Second)
	assert.True(t, mutualExclusionWithWriter.Load())
	assert.True(t, mutualExclusionWithReader.Load())
}

func TestRWPMutexWithReaders(t *testing.T) {
	var mutex RWPMutex
	mutex.RLock() // reader

	var mutualExclusionWithWriter atomic.Bool
	mutualExclusionWithWriter.Store(true)

	go func() {
		mutex.Lock() // another writer
		mutualExclusionWithWriter.Store(false)
	}()

	time.Sleep(time.Second)
	assert.True(t, mutualExclusionWithWriter.Load())
}

func TestRWPMutexMultipleReaders(t *testing.T) {
	var mutex RWPMutex
	mutex.RLock() // reader

	var readersCount atomic.Int32
	readersCount.Add(1)

	go func() {
		mutex.RLock() // another reader
		readersCount.Add(1)
	}()

	go func() {
		mutex.RLock() // another reader
		readersCount.Add(1)
	}()

	time.Sleep(time.Second)
	assert.Equal(t, int32(3), readersCount.Load())
}

func TestRWPMutexWithWriterPriority(t *testing.T) {
	var mutex RWPMutex
	mutex.RLock() // reader

	var mutualExclusionWithWriter atomic.Bool
	mutualExclusionWithWriter.Store(true)
	var readersCount atomic.Int32
	readersCount.Add(1)

	go func() {
		mutex.Lock() // another writer is waiting for reader
		mutualExclusionWithWriter.Store(false)
	}()

	time.Sleep(time.Second)

	go func() {
		mutex.RLock() // another reader is waiting for a higher priority writer
		readersCount.Add(1)
	}()

	go func() {
		mutex.RLock() // another reader is waiting for a higher priority writer
		readersCount.Add(1)
	}()

	time.Sleep(time.Second)

	assert.True(t, mutualExclusionWithWriter.Load())
	assert.Equal(t, int32(1), readersCount.Load())
}
