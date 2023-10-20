package reflect

import (
	"fmt"
	"github.com/murakmii/gj/vm"
)

func ReflectionGetCallerClassV(thread *vm.Thread, _ []interface{}) error {
	className := "java/lang/Class"
	class, state, err := thread.VM().FindInitializedClass(&className, thread)
	if err != nil {
		return err
	}
	if state == vm.FailedInitialization {
		return fmt.Errorf("failed initialization of class class in Reflection.getCallerClass")
	}

	callerClassName := thread.InvokerFrame().CurrentClass().File().ThisClass()
	thread.CurrentFrame().PushOperand(vm.NewInstance(class).SetVMData(&callerClassName))
	return nil
}

func ReflectionGetClassAccessFlags(thread *vm.Thread, args []interface{}) error {
	className := args[0].(*vm.Instance).VMData().(*string)
	fmt.Printf("get access flags for %s\n", *className)

	class, err := thread.VM().FindClass(className)
	if err != nil {
		return err
	}

	thread.CurrentFrame().PushOperand(int(class.File().AccessFlag()))
	return nil
}
