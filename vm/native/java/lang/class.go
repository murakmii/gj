package lang

import (
	"fmt"
	"github.com/murakmii/gj/vm"
)

func ClassRegisterNatives(thread *vm.Thread, args []interface{}) error {
	fmt.Println("execute java/lang/Class.registerNatives")
	return nil
}

func ClassDesiredAssertionStatus0(thread *vm.Thread, args []interface{}) error {
	// return false
	thread.CurrentFrame().PushOperand(0)
	return nil
}

func ClassGetPrimitiveClass(thread *vm.Thread, _ []interface{}) error {
	// TODO: generate class instance
	// return null
	thread.CurrentFrame().PushOperand(nil)
	return nil
}
