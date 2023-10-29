package vm

import (
	"github.com/murakmii/gojiai/class_file"
	"sync"
)

type (
	NativeMethodFunc     func(thread *Thread, args []interface{}) error
	NativeMethodRegistry struct {
		lock  *sync.Mutex
		cache map[*class_file.MethodInfo]NativeMethodFunc
		reg   map[string]NativeMethodFunc
	}
)

var (
	NopNativeMethod NativeMethodFunc = func(_ *Thread, _ []interface{}) error { return nil }
	NativeMethods                    = &NativeMethodRegistry{
		lock:  &sync.Mutex{},
		cache: make(map[*class_file.MethodInfo]NativeMethodFunc),
		reg:   make(map[string]NativeMethodFunc),
	}
)

func (registry NativeMethodRegistry) Register(class, method, desc string, f NativeMethodFunc) {
	// Native methods indexed by class name, method name and method descriptor
	key := class + "/" + method + desc

	registry.reg[key] = f
}

func (registry NativeMethodRegistry) Resolve(class string, method *class_file.MethodInfo) NativeMethodFunc {
	key := class + "/" + *(method.Name()) + method.Descriptor().String()
	return registry.reg[key]
}
