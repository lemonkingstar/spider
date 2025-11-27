package patomic

import (
	"sync/atomic"
	"unsafe"
)

// AtomicValue encapsulates the underlying atomic memory operations.
// Can be used for implementing concurrent locks.
type AtomicValue struct {
	val int32
	ptr unsafe.Pointer
}

func NewAtomicValue(v int) *AtomicValue { return &AtomicValue{val: int32(v)} }
func (p *AtomicValue) GetInt() int      { return int(p.val) }

func (p *AtomicValue) LoadInt() int { return int(atomic.LoadInt32(&p.val)) }
func (p *AtomicValue) AddInt(v int) { atomic.AddInt32(&p.val, int32(v)) }
func (p *AtomicValue) Increase()    { p.AddInt(1) }
func (p *AtomicValue) Decrease()    { p.AddInt(-1) }

func (p *AtomicValue) StoreInt(v int) { atomic.StoreInt32(&p.val, int32(v)) }
func (p *AtomicValue) Set()           { p.StoreInt(1) }
func (p *AtomicValue) Unset()         { p.StoreInt(0) }
func (p *AtomicValue) IsSet() bool    { return p.LoadInt() != 0 }

func (p *AtomicValue) CompareAndSwapInt(old, new int) bool {
	return atomic.CompareAndSwapInt32(&p.val, int32(old), int32(new))
}
func (p *AtomicValue) SwapInt(new int) int { return int(atomic.SwapInt32(&p.val, int32(new))) }

func (p *AtomicValue) StorePointer(v unsafe.Pointer) { atomic.StorePointer(&p.ptr, v) }
func (p *AtomicValue) LoadPointer() unsafe.Pointer   { return atomic.LoadPointer(&p.ptr) }
