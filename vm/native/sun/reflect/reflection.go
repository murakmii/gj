package reflect

import (
	"github.com/murakmii/gj/vm"
)

func ReflectionGetCallerClassV(thread *vm.Thread, _ []interface{}) error {
	callerClassName := thread.InvokerFrame().CurrentClass().File().ThisClass()
	class, err := thread.VM().Class(callerClassName, nil)
	if err != nil {
		return err
	}

	thread.CurrentFrame().PushOperand(class.Java())
	return nil
}

func ReflectionGetClassAccessFlags(thread *vm.Thread, args []interface{}) error {
	className := args[0].(*vm.Instance).VMData().(*string)

	class, err := thread.VM().Class(*className, thread)
	if err != nil {
		return err
	}

	thread.CurrentFrame().PushOperand(int32(class.File().AccessFlag()))
	return nil
}
