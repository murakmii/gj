package lang

import (
	"github.com/murakmii/gojiai/vm"
	"time"
)

func init() {
	class := "java/lang/Thread"

	vm.NativeMethods.Register(class, "currentThread", "()Ljava/lang/Thread;", func(thread *vm.Thread, args []interface{}) error {
		thread.CurrentFrame().PushOperand(thread.JavaThread())
		return nil
	})

	vm.NativeMethods.Register(class, "isAlive", "()Z", func(caller *vm.Thread, args []interface{}) error {
		thread := args[0].(*vm.Instance).AsThread()
		var alive int32
		if thread != nil && thread.IsAlive() {
			alive = 1
		}

		caller.CurrentFrame().PushOperand(alive)
		return nil
	})

	vm.NativeMethods.Register(class, "registerNatives", "()V", vm.NopNativeMethod)
	vm.NativeMethods.Register(class, "setPriority0", "(I)V", vm.NopNativeMethod)

	vm.NativeMethods.Register(class, "sleep", "(J)V", func(thread *vm.Thread, args []interface{}) error {
		time.Sleep(time.Millisecond * time.Duration(args[0].(int64)))
		return nil
	})

	vm.NativeMethods.Register(class, "start0", "()V", func(caller *vm.Thread, args []interface{}) error {
		java := args[0].(*vm.Instance)

		daemon := java.GetField("daemon", "Z").(int32)
		name := java.GetField("name", "Ljava/lang/String;").(*vm.Instance)

		thread := vm.NewThread(caller.VM(), name.AsString(), false, daemon == 1)
		thread.SetJavaThread(java)

		java.ToBeThread(thread)
		java.PutField("threadStatus", "I", int32(0x04))

		class, method := java.Class().ResolveMethod("run", "()V")
		thread.VM().Executor().Start(thread, vm.NewFrame(class, method).SetLocal(0, java))

		return nil
	})
}
