package lang

import (
	"fmt"
	"github.com/murakmii/gj/vm"
)

func ThreadRegisterNatives(thread *vm.Thread, args []interface{}) error {
	fmt.Println("execute java/lang/Thread.registerNatives")
	return nil
}

func ThreadCurrentThread(thread *vm.Thread, args []interface{}) error {
	thread.CurrentFrame().PushOperand(thread.JavaThread())
	return nil
}

func ThreadSetPriority0(thread *vm.Thread, args []interface{}) error {
	// Does NOT support priority changing.
	return nil
}
