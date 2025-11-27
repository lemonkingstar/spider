package patomic

import (
	"reflect"
	"sync"
	"testing"
	"unsafe"
)

func TestAtomicInt(t *testing.T) {
	// not use atomic
	sum := 0
	wg := sync.WaitGroup{}
	for i := 0; i < 10000; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			sum += 1
		}()
	}
	wg.Wait()
	t.Log("sum = ", sum)

	// use atomic
	ato := NewAtomicValue(0)
	for i := 0; i < 10000; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			ato.Increase()
		}()
	}
	wg.Wait()
	t.Log("sum = ", ato.GetInt())
}

func TestAtomicPointer(t *testing.T) {
	type V struct {
		i int32
		j int64
	}
	val := &V{i: 100, j: 1024}
	ato := &AtomicValue{}
	ato.StorePointer(unsafe.Pointer(val))

	val2 := ato.LoadPointer()
	t.Log((*V)(val2))
	t.Log(reflect.TypeOf(val2))
}

func TestAtomicLock(t *testing.T) {
	// not use atomic
	locked := false
	wg := sync.WaitGroup{}
	for i := 0; i < 10000; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			if locked {
				return
			}
			locked = true
			t.Log("first locked.")
		}()
	}
	wg.Wait()

	// use atomic
	ato := NewAtomicValue(0)
	for i := 0; i < 10000; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			if !ato.CompareAndSwapInt(0, 1) {
				return
			}
			t.Log("second locked.")
		}()
	}
	wg.Wait()

	ato.Unset()
	for i := 0; i < 10000; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			if ato.IsSet() {
				return
			}
			ato.Set()
			t.Log("third locked.")
		}()
	}
	wg.Wait()
}
