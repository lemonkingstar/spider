package plock

import (
	"sync"
	"sync/atomic"
)

const (
	lockedFlag   int32 = 1
	unlockedFlag int32 = 0
)

// TryLock
// app try lock
type TryLock struct {
	in     sync.Mutex
	status *int32
}

func NewTryLock() *TryLock {
	status := unlockedFlag
	return &TryLock{
		status: &status,
	}
}

func (p *TryLock) Lock() {
	p.in.Lock()
	atomic.StoreInt32(p.status, lockedFlag)
}

func (p *TryLock) Unlock() {
	p.in.Unlock()
	atomic.StoreInt32(p.status, unlockedFlag)
}

func (p *TryLock) TryLock() bool {
	if atomic.CompareAndSwapInt32(p.status, unlockedFlag, lockedFlag) {
		p.Lock()
		return true
	}
	return false
}

func (p *TryLock) IsLocked() bool {
	if atomic.LoadInt32(p.status) == lockedFlag {
		return true
	}
	return false
}
