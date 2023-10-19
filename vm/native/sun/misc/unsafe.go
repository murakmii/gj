package misc

import (
	"fmt"
	"github.com/murakmii/gj/vm"
	"unsafe"
)

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

func UnsafeObjectFieldOffset(thread *vm.Thread, _ []interface{}) error {
	thread.DumpFrameStack(true)
	return fmt.Errorf("Unsafe.objectFieldOffset")
}
