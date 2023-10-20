package lang

import (
	"fmt"
	"github.com/murakmii/gj/vm"
	"unsafe"
)

func ObjectHashCode(thread *vm.Thread, args []interface{}) error {
	instance, ok := args[0].(*vm.Instance)
	if !ok {
		return fmt.Errorf("arg for Object.hashCode is NOT instance")
	}

	thread.CurrentFrame().PushOperand(int(uintptr(unsafe.Pointer(instance))))
	return nil
}

func ObjectGetClass(thread *vm.Thread, args []interface{}) error {
	instance := args[0].(*vm.Instance)
	className := instance.Class().File().ThisClass()

	thread.CurrentFrame().PushOperand(
		vm.NewInstance(thread.VM().JavaLangClassClass()).SetVMData(&className))

	return nil
}
