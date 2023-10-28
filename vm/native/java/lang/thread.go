package lang

import (
	"github.com/murakmii/gj/vm"
	"time"
)

func ThreadCurrentThread(thread *vm.Thread, args []interface{}) error {
	thread.CurrentFrame().PushOperand(thread.JavaThread())
	return nil
}

func ThreadIsAlive(thread *vm.Thread, args []interface{}) error {
	targetThread := args[0].(*vm.Instance).VMData()
	var alive int32
	if targetThread != nil && targetThread.(*vm.Thread).IsAlive() {
		alive = 1
	}

	thread.CurrentFrame().PushOperand(alive)
	return nil
}

func ThreadStart0(thread *vm.Thread, args []interface{}) error {
	jThread := args[0].(*vm.Instance)

	daemon := jThread.GetField("daemon", "Z").(int32)
	name := jThread.GetField("name", "Ljava/lang/String;").(*vm.Instance)

	newThread := vm.NewThread(thread.VM(), name.GetCharArrayField("value"), false, daemon == 1)

	jThread.SetVMData(newThread)
	newThread.SetJavaThread(jThread)

	class, method := jThread.Class().ResolveMethod("run", "()V")

	thread.VM().Executor().Start(
		newThread,
		vm.NewFrame(class, method).SetLocal(0, jThread),
	)

	return nil
}

func ThreadSleep(thread *vm.Thread, args []interface{}) error {
	sleep := args[0].(int64)
	time.Sleep(time.Millisecond * time.Duration(sleep))
	return nil
}

func ThreadSetNativeName(_ *vm.Thread, args []interface{}) error {
	args[0].(*vm.Instance).VMData().(*vm.Thread).
		SetName(args[1].(*vm.Instance).GetCharArrayField("value"))

	return nil
}
