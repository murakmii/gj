package misc

import (
	"fmt"
	"github.com/murakmii/gj/vm"
	"unsafe"
)

func UnsafeRegisterNatives(_ *vm.Thread, args []interface{}) error {
	fmt.Println("execute sum/misc/Unsafe.registerNatives")
	return nil
}

func UnsafeAddressSize(thread *vm.Thread, _ []interface{}) error {
	thread.CurrentFrame().PushOperand(int(unsafe.Sizeof(uintptr(0))))
	return nil
}

func UnsafeArrayBaseOffset(thread *vm.Thread, _ []interface{}) error {
	thread.CurrentFrame().PushOperand(0)
	return nil
}

func UnsafeArrayIndexScale(thread *vm.Thread, _ []interface{}) error {
	thread.CurrentFrame().PushOperand(1)
	return nil
}
