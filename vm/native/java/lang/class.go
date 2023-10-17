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
	thread.PopFrame()
	if thread.CurrentFrame() != nil {
		thread.CurrentFrame().PushOperand(0)
	}
	return nil
}
