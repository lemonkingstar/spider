package putil

import (
	"github.com/mohae/deepcopy"
)

// Clone 深拷贝
// Usage: obj := Clone(input)
// obj.(StructObj)
func Clone(input interface{}) interface{} {
	return deepcopy.Copy(input)
}
