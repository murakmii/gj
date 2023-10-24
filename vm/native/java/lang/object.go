package lang

import (
	"fmt"
	"github.com/murakmii/gj/vm"
	"unsafe"
)

func ObjectClone(thread *vm.Thread, args []interface{}) error {
	thread.CurrentFrame().PushOperand(args[0].(*vm.Instance).Clone())
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

func ObjectGetClass(thread *vm.Thread, args []interface{}) error {
	thread.CurrentFrame().PushOperand(args[0].(*vm.Instance).Class().Java())
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
