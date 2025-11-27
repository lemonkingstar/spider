package patomic

import (
	"runtime"
	"time"
)

type SimpleLock interface {
	Lock()
	Unlock()
	LockWithTimeout(time.Duration) bool
	TryLock() bool
}

type lock struct {
	ato *AtomicValue
}

func NewLock() SimpleLock {
	return &lock{ato: NewAtomicValue(0)}
}

func (l *lock) Lock() {
	for !l.TryLock() {
		runtime.Gosched()
	}
}

func (l *lock) LockWithTimeout(timeout time.Duration) bool {
	deadline := time.Now().Add(timeout)
	for !l.TryLock() {
		if time.Now().After(deadline) {
			return false
		}
		runtime.Gosched()
	}
	return true
}

func (l *lock) Unlock() {
	l.ato.StoreInt(0)
}

func (l *lock) TryLock() bool {
	return l.ato.CompareAndSwapInt(0, 1)
}
