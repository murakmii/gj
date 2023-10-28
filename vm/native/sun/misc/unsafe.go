package misc

import (
	"encoding/binary"
	"fmt"
	"github.com/murakmii/gj/vm"
	"unsafe"
)

func UnsafeAllocateMemory(thread *vm.Thread, args []interface{}) error {
	thread.CurrentFrame().PushOperand(thread.VM().NativeMem().Alloc(args[1].(int64)))
	return nil
}

func UnsafeFreeMemory(thread *vm.Thread, args []interface{}) error {
	thread.VM().NativeMem().Free(args[1].(int64))
	return nil
}

func UnsafePutLongNativeMem(thread *vm.Thread, args []interface{}) error {
	binary.BigEndian.PutUint64(
		thread.VM().NativeMem().Ref(args[1].(int64)),
		uint64(args[2].(int64)),
	)
	return nil
}

func UnsafeGetByteNativeMem(thread *vm.Thread, args []interface{}) error {
	thread.CurrentFrame().PushOperand(
		int32(thread.VM().NativeMem().Ref(args[1].(int64))[0]))
	return nil
}

func UnsafeAddressSize(thread *vm.Thread, _ []interface{}) error {
	thread.CurrentFrame().PushOperand(int32(unsafe.Sizeof(uintptr(0))))
	return nil
}

func UnsafeArrayBaseOffset(thread *vm.Thread, _ []interface{}) error {
	thread.CurrentFrame().PushOperand(int32(0))
	return nil
}

func UnsafeArrayIndexScale(thread *vm.Thread, _ []interface{}) error {
	thread.CurrentFrame().PushOperand(int32(1))
	return nil
}

func UnsafeCompareAndSwapInt(thread *vm.Thread, args []interface{}) error {
	obj, ok := args[1].(*vm.Instance)
	if !ok {
		return fmt.Errorf("Unsafe.compareAndSwapInt received arg[1] is NOT instance")
	}

	fID, ok := args[2].(int64)
	if !ok {
		return fmt.Errorf("Unsafe.compareAndSwapInt received arg[2] is NOT int64")
	}

	cmp := args[3].(int32)
	set := args[4].(int32)

	result, err := obj.CompareAndSwapInt(int(fID), cmp, set)
	if err != nil {
		return err
	}

	var ret int32
	if result {
		ret = 1
	}

	thread.CurrentFrame().PushOperand(ret)
	return nil
}

func UnsafeCompareAndSwapLong(thread *vm.Thread, args []interface{}) error {
	obj, ok := args[1].(*vm.Instance)
	if !ok {
		return fmt.Errorf("Unsafe.compareAndSwapLong received arg[1] is NOT instance")
	}

	fID, ok := args[2].(int64)
	if !ok {
		return fmt.Errorf("Unsafe.compareAndSwapInt received arg[2] is NOT int64")
	}

	cmp := args[3].(int64)
	set := args[4].(int64)

	result, err := obj.CompareAndSwapLong(int(fID), cmp, set)
	if err != nil {
		return err
	}

	var ret int32
	if result {
		ret = 1
	}

	thread.CurrentFrame().PushOperand(ret)
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

	var ret int32
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

	slot, ok := fieldInstance.GetField("slot", "I").(int32)
	if !ok {
		return fmt.Errorf("can't get slot in Unsafe.objectFieldOffset")
	}

	thread.CurrentFrame().PushOperand(int64(slot))
	return nil
}

func UnsafeGetObjectVolatile(thread *vm.Thread, args []interface{}) error {
	instance := args[1].(*vm.Instance)
	offset := args[2].(int64)

	thread.CurrentFrame().PushOperand(instance.GetFieldByID(int(offset)))
	return nil
}

func UnsafeGetIntVolatile(thread *vm.Thread, args []interface{}) error {
	instance := args[1].(*vm.Instance)
	value := instance.GetFieldByID(int(args[2].(int64)))

	var result int32
	var ok bool
	if value != nil {
		result, ok = value.(int32)
		if !ok {
			return fmt.Errorf("fetched value is NOT int(%+v) in Unsafe.getIntVolatile", value)
		}
	}

	thread.CurrentFrame().PushOperand(result)
	return nil
}

func UnsafePutObjectVolatile(thread *vm.Thread, args []interface{}) error {
	instance := args[1].(*vm.Instance)
	offset := args[2].(int64)
	value := args[3]

	instance.PutFieldByID(int(offset), value)
	return nil
}
