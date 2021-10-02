package utils

import "sync"

type Locker interface {
	ReadLock()
	ReadUnlock()
	WriteLock()
	WriteUnlock()
}

type NullLocker struct {}

func (NullLocker) ReadLock() {}
func (NullLocker) ReadUnlock() {}
func (NullLocker) WriteLock() {}
func (NullLocker) WriteUnlock() {}

type MutexLock struct {
	mu sync.RWMutex
}

func (m *MutexLock) ReadLock() { m.mu.RLock() }
func (m *MutexLock) ReadUnlock() { m.mu.RUnlock() }
func (m *MutexLock) WriteLock() { m.mu.Lock() }
func (m *MutexLock) WriteUnlock() { m.mu.Unlock() }

