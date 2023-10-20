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

func UnsafeCompareAndSwapObject(thread *vm.Thread, args []interface{}) error {
	obj, ok := args[1].(*vm.Instance)
	if !ok {
		return fmt.Errorf("Unsafe.compareAndSwapObject received arg[1] is NOT instance")
	}

	fID, ok := args[2].(int64)
	if !ok {
		return fmt.Errorf("Unsafe.compareAndSwapObject received arg[2] is NOT int64")
	}

	var cmp *vm.Instance
	if e, ok := args[3].(*vm.Instance); ok {
		cmp = e
	}

	var set *vm.Instance
	if s, ok := args[4].(*vm.Instance); ok {
		set = s
	}

	result, err := obj.CompareAndSwap(int(fID), cmp, set)
	if err != nil {
		return err
	}

	ret := 0
	if result {
		ret = 1
	}

	thread.CurrentFrame().PushOperand(ret)
	return nil
}

func UnsafeObjectFieldOffset(thread *vm.Thread, args []interface{}) error {
	fieldInstance, ok := args[1].(*vm.Instance)
	if !ok || fieldInstance.Class().File().ThisClass() != "java/lang/reflect/Field" {
		return fmt.Errorf("Unsafe.objectFieldOffset received arg is NOT field instance")
	}

	slotName := "slot"
	slotDesc := "I"
	slot, ok := fieldInstance.GetField(&slotName, &slotDesc).(int)
	if !ok {
		return fmt.Errorf("can't get slot in Unsafe.objectFieldOffset")
	}

	thread.CurrentFrame().PushOperand(int64(slot))
	return nil
}
