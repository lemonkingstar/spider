package patomic

import (
	"fmt"
	"reflect"
	"sync"
	"testing"
	"unsafe"
)

func TestAtomicInt(t *testing.T) {
	// 不使用 atomic操作
	sum := 0
	wg := sync.WaitGroup{}
	for i := 0; i < 10000; i++ {
		wg.Add(1)
		go func() {
			sum += 1
			wg.Done()
		}()
	}
	wg.Wait()
	fmt.Println("sum = ", sum)

	// 使用 atomic操作
	ato := NewAtomicValue(0)
	for i := 0; i < 10000; i++ {
		wg.Add(1)
		go func() {
			ato.Increase()
			wg.Done()
		}()
	}
	wg.Wait()
	fmt.Println("ato = ", ato.GetInt())
}

func TestAtomicPointer(t *testing.T) {
	type V struct {
		i int32
		j int64
	}
	v := &V{
		i: 100,
		j: 1024,
	}

	ato := &AtomicValue{}
	ato.StorePointer(unsafe.Pointer(v))

	v2 := ato.LoadPointer()
	fmt.Println((*V)(v2))

	fmt.Println(reflect.TypeOf(v2))
}
