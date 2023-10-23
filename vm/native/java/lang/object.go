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

func ObjectWait(thread *vm.Thread, args []interface{}) error {
	// TODO: interrupt
	_, err := args[0].(*vm.Instance).Monitor().Wait(thread, int(args[1].(int64)))
	return err
}

func ObjectNotifyAll(thread *vm.Thread, args []interface{}) error {
	return args[0].(*vm.Instance).Monitor().NotifyAll(thread)
}
