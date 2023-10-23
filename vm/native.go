package vm

import (
	"fmt"
	"github.com/murakmii/gj/class_file"
)

type (
	NativeMethodFunc = func(thread *Thread, args []interface{}) error
)

var nativeMethods = make(map[string]NativeMethodFunc)

func RegisterNativeMethod(id string, method NativeMethodFunc) {
	nativeMethods[id] = method
}

func CallNativeMethod(thread *Thread, class *Class, method *class_file.MethodInfo, args []interface{}) error {
	id := class.File().ThisClass() + "/" + *method.Name() + *method.Descriptor()
	native, exist := nativeMethods[id]
	if !exist {
		return fmt.Errorf("native method not found: %s", id)
	}
	return native(thread, args)
}
