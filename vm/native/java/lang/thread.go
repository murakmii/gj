package lang

import (
	"github.com/murakmii/gj/vm"
)

func ThreadIsAlive(thread *vm.Thread, _ []interface{}) error {
	thread.CurrentFrame().PushOperand(0)
	return nil
}

func ThreadStart0(_ *vm.Thread, _ []interface{}) error {
	return nil
}

func ThreadCurrentThread(thread *vm.Thread, args []interface{}) error {
	thread.CurrentFrame().PushOperand(thread.JavaThread())
	return nil
}
