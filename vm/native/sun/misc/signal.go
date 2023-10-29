package misc

import (
	"github.com/murakmii/gj/vm"
	"os"
	"syscall"
)

func init() {
	class := "sun/misc/Signal"

	signals := map[string]os.Signal{
		"HUP":  syscall.SIGHUP,
		"TERM": syscall.SIGTERM,
		"INT":  syscall.SIGINT,
	}

	vm.NativeMethods.Register(class, "findSignal", "(Ljava/lang/String;)I", func(thread *vm.Thread, args []interface{}) error {
		sig, exist := signals[args[0].(*vm.Instance).AsString()]
		if !exist {
			sig = syscall.Signal(-1)
		}

		thread.CurrentFrame().PushOperand(int32(sig.(syscall.Signal)))
		return nil
	})

	vm.NativeMethods.Register(class, "handle0", "(IJ)J", func(thread *vm.Thread, args []interface{}) error {
		// TODO save handler number
		thread.CurrentFrame().PushOperand(int64(0))
		return nil
	})
}
