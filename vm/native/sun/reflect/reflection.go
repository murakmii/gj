package reflect

import (
	"github.com/murakmii/gj/vm"
)

func ReflectionGetCallerClassV(thread *vm.Thread, _ []interface{}) error {
	callerClassName := thread.InvokerFrame().CurrentClass().File().ThisClass()
	thread.CurrentFrame().PushOperand(thread.VM().JavaClass(&callerClassName))
	return nil
}

func ReflectionGetClassAccessFlags(thread *vm.Thread, args []interface{}) error {
	className := args[0].(*vm.Instance).VMData().(*string)

	class, err := thread.VM().FindClass(className)
	if err != nil {
		return err
	}

	thread.CurrentFrame().PushOperand(int(class.File().AccessFlag()))
	return nil
}
