package lang

import (
	"github.com/murakmii/gj/vm"
	"time"
)

func ThreadCurrentThread(thread *vm.Thread, args []interface{}) error {
	thread.CurrentFrame().PushOperand(thread.JavaThread())
	return nil
}

func ThreadIsAlive(callerThread *vm.Thread, args []interface{}) error {
	thread := args[0].(*vm.Instance).AsThread()
	var alive int32
	if thread != nil && thread.IsAlive() {
		alive = 1
	}

	callerThread.CurrentFrame().PushOperand(alive)
	return nil
}

func ThreadStart0(callerThread *vm.Thread, args []interface{}) error {
	java := args[0].(*vm.Instance)

	daemon := java.GetField("daemon", "Z").(int32)
	name := java.GetField("name", "Ljava/lang/String;").(*vm.Instance)

	thread := vm.NewThread(callerThread.VM(), name.AsString(), false, daemon == 1)
	thread.SetJavaThread(java)

	java.ToBeThread(thread)
	java.PutField("threadStatus", "I", int32(0x04))

	class, method := java.Class().ResolveMethod("run", "()V")
	thread.VM().Executor().Start(thread, vm.NewFrame(class, method).SetLocal(0, java))

	return nil
}

func ThreadSleep(_ *vm.Thread, args []interface{}) error {
	time.Sleep(time.Millisecond * time.Duration(args[0].(int64)))
	return nil
}
