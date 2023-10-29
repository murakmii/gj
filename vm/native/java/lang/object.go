package lang

import (
	"github.com/murakmii/gojiai/vm"
)

func init() {
	class := "java/lang/Object"

	vm.NativeMethods.Register(class, "clone", "()Ljava/lang/Object;", func(thread *vm.Thread, args []interface{}) error {
		thread.CurrentFrame().PushOperand(args[0].(*vm.Instance).Clone())
		return nil
	})

	vm.NativeMethods.Register(class, "getClass", "()Ljava/lang/Class;", func(thread *vm.Thread, args []interface{}) error {
		thread.CurrentFrame().PushOperand(args[0].(*vm.Instance).Class().Java())
		return nil
	})

	vm.NativeMethods.Register(class, "hashCode", "()I", func(thread *vm.Thread, args []interface{}) error {
		thread.CurrentFrame().PushOperand(args[0].(*vm.Instance).HashCode())
		return nil
	})

	vm.NativeMethods.Register(class, "notifyAll", "()V", func(thread *vm.Thread, args []interface{}) error {
		return args[0].(*vm.Instance).Monitor().NotifyAll(thread)
	})

	vm.NativeMethods.Register(class, "registerNatives", "()V", vm.NopNativeMethod)

	vm.NativeMethods.Register(class, "wait", "(J)V", func(thread *vm.Thread, args []interface{}) error {
		// TODO: interrupt
		_, err := args[0].(*vm.Instance).Monitor().Wait(thread, int(args[1].(int64)))
		return err
	})
}
