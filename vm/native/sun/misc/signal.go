package misc

import (
	"github.com/murakmii/gj/vm"
	"os"
	"syscall"
)

var signals = map[string]os.Signal{
	"HUP":  syscall.SIGHUP,
	"TERM": syscall.SIGTERM,
	"INT":  syscall.SIGINT,
}

func SignalFindSignal(thread *vm.Thread, args []interface{}) error {
	sig, exist := signals[vm.JavaStringToGoString(args[0].(*vm.Instance))]
	if !exist {
		sig = syscall.Signal(-1)
	}

	thread.CurrentFrame().PushOperand(int(sig.(syscall.Signal)))
	return nil
}

func SignalHandle(thread *vm.Thread, _ []interface{}) error {
	// TODO save handler number
	thread.CurrentFrame().PushOperand(int64(0))
	return nil
}
