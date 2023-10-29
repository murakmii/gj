package misc

import (
	"encoding/binary"
	"fmt"
	"github.com/murakmii/gj/vm"
	"unsafe"
)

func init() {
	class := "sun/misc/Unsafe"

	vm.NativeMethods.Register(class, "addressSize", "()I", func(thread *vm.Thread, args []interface{}) error {
		thread.CurrentFrame().PushOperand(int32(unsafe.Sizeof(uintptr(0))))
		return nil
	})

	vm.NativeMethods.Register(class, "allocateMemory", "(J)J", func(thread *vm.Thread, args []interface{}) error {
		thread.CurrentFrame().PushOperand(thread.VM().NativeMem().Alloc(args[1].(int64)))
		return nil
	})

	vm.NativeMethods.Register(class, "arrayBaseOffset", "(Ljava/lang/Class;)I", func(thread *vm.Thread, args []interface{}) error {
		thread.CurrentFrame().PushOperand(int32(0))
		return nil
	})

	vm.NativeMethods.Register(class, "arrayIndexScale", "(Ljava/lang/Class;)I", func(thread *vm.Thread, args []interface{}) error {
		thread.CurrentFrame().PushOperand(int32(1))
		return nil
	})

	vm.NativeMethods.Register(class, "compareAndSwapInt", "(Ljava/lang/Object;JII)Z", func(thread *vm.Thread, args []interface{}) error {
		obj := args[1].(*vm.Instance)
		fID := args[2].(int64)
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
	})

	vm.NativeMethods.Register(class, "compareAndSwapLong", "(Ljava/lang/Object;JJJ)Z", func(thread *vm.Thread, args []interface{}) error {
		obj := args[1].(*vm.Instance)
		fID := args[2].(int64)
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
	})

	vm.NativeMethods.Register(class, "compareAndSwapObject", "(Ljava/lang/Object;JLjava/lang/Object;Ljava/lang/Object;)Z", func(thread *vm.Thread, args []interface{}) error {
		obj := args[1].(*vm.Instance)
		fID := args[2].(int64)

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
	})

	vm.NativeMethods.Register(class, "freeMemory", "(J)V", func(thread *vm.Thread, args []interface{}) error {
		thread.VM().NativeMem().Free(args[1].(int64))
		return nil
	})

	vm.NativeMethods.Register(class, "getByte", "(J)B", func(thread *vm.Thread, args []interface{}) error {
		thread.CurrentFrame().PushOperand(int32(thread.VM().NativeMem().Ref(args[1].(int64))[0]))
		return nil
	})

	vm.NativeMethods.Register(class, "getIntVolatile", "(Ljava/lang/Object;J)I", func(thread *vm.Thread, args []interface{}) error {
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
	})

	vm.NativeMethods.Register(class, "getObjectVolatile", "(Ljava/lang/Object;J)Ljava/lang/Object;", func(thread *vm.Thread, args []interface{}) error {
		instance := args[1].(*vm.Instance)
		offset := args[2].(int64)

		thread.CurrentFrame().PushOperand(instance.GetFieldByID(int(offset)))
		return nil
	})

	vm.NativeMethods.Register(class, "objectFieldOffset", "(Ljava/lang/reflect/Field;)J", func(thread *vm.Thread, args []interface{}) error {
		slot, ok := args[1].(*vm.Instance).GetField("slot", "I").(int32)
		if !ok {
			return fmt.Errorf("can't get slot in Unsafe.objectFieldOffset")
		}

		thread.CurrentFrame().PushOperand(int64(slot))
		return nil
	})

	vm.NativeMethods.Register(class, "putLong", "(JJ)V", func(thread *vm.Thread, args []interface{}) error {
		binary.BigEndian.PutUint64(
			thread.VM().NativeMem().Ref(args[1].(int64)),
			uint64(args[2].(int64)),
		)
		return nil
	})

	vm.NativeMethods.Register(class, "putObjectVolatile", "(Ljava/lang/Object;JLjava/lang/Object;)V", func(thread *vm.Thread, args []interface{}) error {
		instance := args[1].(*vm.Instance)
		offset := args[2].(int64)
		value := args[3]

		instance.PutFieldByID(int(offset), value)
		return nil
	})

	vm.NativeMethods.Register(class, "registerNatives", "()V", vm.NopNativeMethod)
}
