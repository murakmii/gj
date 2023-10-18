package lang

import (
	"fmt"
	"github.com/murakmii/gj/vm"
	"unsafe"
)

func ObjectRegisterNatives(thread *vm.Thread, args []interface{}) error {
	fmt.Println("execute java/lang/Object.registerNatives")
	return nil
}

func ObjectHashCode(thread *vm.Thread, args []interface{}) error {
	instance, ok := args[0].(*vm.Instance)
	if !ok {
		return fmt.Errorf("arg for Object.hashCode is NOT instance")
	}

	thread.CurrentFrame().PushOperand(int(uintptr(unsafe.Pointer(instance))))
	return nil
}
