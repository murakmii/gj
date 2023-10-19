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

	thread.CurrentFrame().PushOperand(vm.NewInstance(class).SetVMData(
		thread.CurrentFrame().CurrentClass().File().ThisClass(),
	))
	return nil
}
