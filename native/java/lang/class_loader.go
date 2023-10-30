package lang

import (
	"github.com/murakmii/gojiai/vm"
	"strings"
)

func init() {
	_class := "java/lang/ClassLoader"

	vm.NativeMethods.Register(_class, "registerNatives", "()V", vm.NopNativeMethod)

	findClass := func(thread *vm.Thread, args []interface{}) error {
		className := strings.ReplaceAll(args[1].(*vm.Instance).AsString(), ".", "/")
		class, err := thread.VM().Class(className, thread)
		if err != nil {
			if strings.Contains(err.Error(), "not found") {
				thread.CurrentFrame().PushOperand(nil)
				return nil
			}
			return err
		}

		thread.CurrentFrame().PushOperand(class.Java())
		return nil
	}

	vm.NativeMethods.Register(_class, "findLoadedClass0", "(Ljava/lang/String;)Ljava/lang/Class;", findClass)
	vm.NativeMethods.Register(_class, "findBootstrapClass", "(Ljava/lang/String;)Ljava/lang/Class;", findClass)
}
