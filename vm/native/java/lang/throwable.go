package lang

import "github.com/murakmii/gj/vm"

func ThrowableFillInStackTrace(thread *vm.Thread, args []interface{}) error {
	throwable := args[0].(*vm.Instance)
	throwable.SetVMData(thread.StackTrack())

	thread.CurrentFrame().PushOperand(throwable)
	return nil
}
